package protocol

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol/pool"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type repgenFollowerPhase string

const (
	repgenFollowerPhaseUnknown     repgenFollowerPhase = "unknown"
	repgenFollowerPhaseNewEpoch    repgenFollowerPhase = "newEpoch"
	repgenFollowerPhaseReady       repgenFollowerPhase = "ready"
	repgenFollowerPhaseSentObserve repgenFollowerPhase = "sentObserve"
	repgenFollowerPhaseSentPrepare repgenFollowerPhase = "sentPrepare"
	repgenFollowerPhaseSentCommit  repgenFollowerPhase = "sentCommit"
)

///////////////////////////////////////////////////////////
// Report Generation Follower
///////////////////////////////////////////////////////////

func (repgen *reportGenerationState[RI]) messageStartEpoch(msg MessageStartEpoch[RI], sender commontypes.OracleID) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessageStartEpoch for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if sender != repgen.l {
		repgen.logger.Warn("Non-leader sent MessageStartEpoch", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	if repgen.followerState.phase != repgenFollowerPhaseNewEpoch {
		repgen.logger.Warn("Got MessageStartEpoch for wrong phase", commontypes.LogFields{
			"sender": sender,
			"phase":  repgen.followerState.phase,
		})
		return
	}

	{
		err := msg.StartEpochProof.Verify(
			repgen.Timestamp(),
			repgen.config.OracleIdentities,
			repgen.config.N(),
			repgen.config.F,
		)
		if err != nil {
			repgen.logger.Warn("MessageStartEpoch contains invalid StartRoundQuorumCertificate", commontypes.LogFields{
				"sender": repgen.l,
				"error":  err,
			})
			return
		}
	}

	if msg.StartEpochProof.HighestCertified.IsGenesis() {
		repgen.followerState.firstSeqNrOfEpoch = repgen.followerState.deliveredSeqNr + 1
		repgen.startSubsequentFollowerRound()
	} else if commitQC, ok := msg.StartEpochProof.HighestCertified.(*CertifiedPrepareOrCommitCommit); ok {
		repgen.deliver(*commitQC)
		repgen.followerState.firstSeqNrOfEpoch = repgen.followerState.deliveredSeqNr + 1
		repgen.startSubsequentFollowerRound()
	} else {
		// We're dealing with a re-proposal from a failed epoch

		prepareQc := msg.StartEpochProof.HighestCertified.(*CertifiedPrepareOrCommitPrepare)

		outcomeDigest := MakeOutcomeDigest(prepareQc.Outcome)

		prepareSignature, err := MakePrepareSignature(
			repgen.Timestamp(),
			prepareQc.SeqNr,
			OutcomeInputsDigest{},
			outcomeDigest,
			repgen.offchainKeyring.OffchainSign,
		)
		if err != nil {
			repgen.logger.Critical("Failed to sign Prepare", commontypes.LogFields{
				"error": err,
			})
			return
		}

		repgen.followerState.phase = repgenFollowerPhaseSentPrepare
		repgen.followerState.firstSeqNrOfEpoch = prepareQc.SeqNr + 1
		repgen.followerState.seqNr = prepareQc.SeqNr
		repgen.followerState.currentOutcome = prepareQc.Outcome
		repgen.followerState.currentOutcomeDigest = outcomeDigest
		repgen.logger.Debug("Broadcasting MessagePrepare (reproposal)", commontypes.LogFields{
			"seqNr": prepareQc.SeqNr,
		})
		repgen.netSender.Broadcast(MessagePrepare[RI]{
			repgen.e,
			prepareQc.SeqNr,
			prepareSignature,
		})
	}
}

func (repgen *reportGenerationState[RI]) startSubsequentFollowerRound() {
	repgen.followerState.phase = repgenFollowerPhaseReady
	repgen.followerState.seqNr = repgen.followerState.deliveredSeqNr + 1
	repgen.followerState.query = nil
	repgen.followerState.currentOutcome = nil
	repgen.followerState.currentOutcomeDigest = OutcomeDigest{}

	repgen.tryProcessStartRoundPool()
}

func (repgen *reportGenerationState[RI]) messageStartRound(msg MessageStartRound[RI], sender commontypes.OracleID) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessageStartRound for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if sender != repgen.l {
		repgen.logger.Warn("Non-leader sent MessageStartRound", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	if putResult := repgen.followerState.startRoundPool.Put(msg.SeqNr, sender, msg); putResult != pool.PutResultOK {
		repgen.logger.Warn("Dropping MessageStartRound", commontypes.LogFields{
			"sender": sender,
			"seqNr":  msg.SeqNr,
			"reason": putResult,
		})
		return
	}

	repgen.logger.Debug("Pooled MessageStartRound", commontypes.LogFields{
		"sender": sender,
		"seqNr":  msg.SeqNr,
	})

	repgen.tryProcessStartRoundPool()
}

func (repgen *reportGenerationState[RI]) tryProcessStartRoundPool() {
	if repgen.followerState.phase != repgenFollowerPhaseReady {
		repgen.logger.Debug("cannot process StartRoundPool, wrong phase", commontypes.LogFields{
			"phase": repgen.followerState.phase,
		})
		return
	}

	poolEntries := repgen.followerState.startRoundPool.Entries(repgen.followerState.seqNr)

	if poolEntries == nil || poolEntries[repgen.l] == nil {

		repgen.logger.Debug("cannot process StartRoundPool, it's empty", commontypes.LogFields{
			"followerStateSeqNr": repgen.followerState.seqNr,
		})
		return
	}

	if repgen.followerState.query != nil {
		repgen.logger.Warn("cannot process StartRoundPool, query already set", commontypes.LogFields{
			"seqNr": repgen.followerState.seqNr,
		})
		return
	}

	msg := poolEntries[repgen.l].Item

	repgen.followerState.query = &msg.Query

	o, ok := callPlugin[types.Observation](
		repgen,
		"Observation",
		repgen.config.MaxDurationObservation,
		repgen.OutcomeCtx(repgen.followerState.seqNr),
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Observation, error) {
			return repgen.reportingPlugin.Observation(ctx, outctx, *repgen.followerState.query)
		},
	)
	if !ok {
		return
	}

	so, err := MakeSignedObservation(repgen.Timestamp(), msg.Query, o, repgen.offchainKeyring.OffchainSign)
	if err != nil {
		repgen.logger.Error("messageStartRound: could not make SignedObservation observation", commontypes.LogFields{
			"seqNr": repgen.followerState.seqNr,
			"error": err,
		})
		return
	}

	if err := so.Verify(repgen.Timestamp(), msg.Query, repgen.offchainKeyring.OffchainPublicKey()); err != nil {
		repgen.logger.Error("MakeSignedObservation produced invalid signature:", commontypes.LogFields{
			"seqNr": repgen.followerState.seqNr,
			"error": err,
		})
		return
	}

	repgen.followerState.phase = repgenFollowerPhaseSentObserve
	repgen.logger.Debug("sent observation to leader", commontypes.LogFields{
		"seqNr": repgen.followerState.seqNr,
	})
	repgen.netSender.SendTo(MessageObserve[RI]{
		repgen.e,
		repgen.followerState.seqNr,
		so,
	}, repgen.l)

	repgen.tryProcessProposePool()
}

func (repgen *reportGenerationState[RI]) messagePropose(msg MessagePropose[RI], sender commontypes.OracleID) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessagePropose for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if sender != repgen.l {
		repgen.logger.Warn("Non-leader sent MessagePropose", commontypes.LogFields{
			"msgSeqNr": msg.SeqNr,
			"sender":   sender,
		})
		return
	}

	if putResult := repgen.followerState.proposePool.Put(msg.SeqNr, sender, msg); putResult != pool.PutResultOK {
		repgen.logger.Warn("Dropping MessagePropose", commontypes.LogFields{
			"sender": sender,
			"seqNr":  msg.SeqNr,
			"reason": putResult,
		})
		return
	}

	repgen.logger.Debug("Pooled MessagePropose", commontypes.LogFields{
		"sender": sender,
		"seqNr":  msg.SeqNr,
	})

	repgen.tryProcessProposePool()
}

func (repgen *reportGenerationState[RI]) tryProcessProposePool() {
	if repgen.followerState.phase != repgenFollowerPhaseSentObserve {
		repgen.logger.Debug("cannot process ProposePool, wrong phase", commontypes.LogFields{
			"phase": repgen.followerState.phase,
		})
		return
	}

	poolEntries := repgen.followerState.proposePool.Entries(repgen.followerState.seqNr)

	if poolEntries == nil || poolEntries[repgen.l] == nil {

		return
	}

	msg := poolEntries[repgen.l].Item

	if msg.SeqNr <= repgen.followerState.deliveredSeqNr {
		repgen.logger.Critical("MessagePropose contains invalid SeqNr", commontypes.LogFields{
			"sender":         repgen.l,
			"msgSeqNr":       msg.SeqNr,
			"deliveredSeqNr": repgen.followerState.deliveredSeqNr,
		})
		return
	}

	attributedObservations := []types.AttributedObservation{}
	{
		if len(msg.AttributedSignedObservations) <= 2*repgen.config.F {
			repgen.logger.Debug("MessagePropose contains too few signed observations", nil)
			return
		}
		seen := map[commontypes.OracleID]bool{}
		for _, aso := range msg.AttributedSignedObservations {
			if !(0 <= int(aso.Observer) && int(aso.Observer) <= repgen.config.N()) {
				repgen.logger.Debug("MessagePropose contains signed observation with invalid observer", commontypes.LogFields{
					"invalidObserver": aso.Observer,
				})
				return
			}

			if seen[aso.Observer] {
				repgen.logger.Debug("MessagePropose contains duplicate signed observation", nil)
				return
			}

			seen[aso.Observer] = true

			if err := aso.SignedObservation.Verify(repgen.Timestamp(), *repgen.followerState.query, repgen.config.OracleIdentities[aso.Observer].OffchainPublicKey); err != nil {
				repgen.logger.Debug("MessagePropose contains signed observation with invalid signature", nil)
				return
			}

			attributedObservations = append(attributedObservations, types.AttributedObservation{
				aso.SignedObservation.Observation,
				aso.Observer,
			})
		}
	}

	outcomeInputsDigest := MakeOutcomeInputsDigest(
		repgen.Timestamp(),
		repgen.followerState.deliveredOutcome,
		repgen.followerState.seqNr,
		*repgen.followerState.query,
		attributedObservations,
	)

	outcome, ok := callPlugin[ocr3types.Outcome](
		repgen,
		"Outcome",
		0,
		repgen.OutcomeCtx(repgen.followerState.seqNr),
		func(_ context.Context, outctx ocr3types.OutcomeContext) (ocr3types.Outcome, error) {
			return repgen.reportingPlugin.Outcome(outctx, *repgen.followerState.query, attributedObservations)
		},
	)
	if !ok {
		return
	}

	outcomeDigest := MakeOutcomeDigest(outcome)

	prepareSignature, err := MakePrepareSignature(
		repgen.Timestamp(),
		msg.SeqNr,
		outcomeInputsDigest,
		outcomeDigest,
		repgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		repgen.logger.Critical("Failed to sign Prepare", commontypes.LogFields{
			"error": err,
		})
		return
	}

	repgen.followerState.phase = repgenFollowerPhaseSentPrepare
	repgen.followerState.currentOutcomeInputsDigest = outcomeInputsDigest
	repgen.followerState.currentOutcome = outcome
	repgen.followerState.currentOutcomeDigest = outcomeDigest

	repgen.logger.Debug("Broadcasting MessagePrepare", commontypes.LogFields{
		"seqNr": msg.SeqNr,
	})
	repgen.netSender.Broadcast(MessagePrepare[RI]{
		repgen.e,
		msg.SeqNr,
		prepareSignature,
	})
}

func (repgen *reportGenerationState[RI]) messagePrepare(msg MessagePrepare[RI], sender commontypes.OracleID) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessagePrepare for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if putResult := repgen.followerState.preparePool.Put(msg.SeqNr, sender, msg.Signature); putResult != pool.PutResultOK {
		repgen.logger.Debug("Dropping MessagePrepare", commontypes.LogFields{
			"sender": sender,
			"seqNr":  msg.SeqNr,
			"reason": putResult,
		})
		return
	}

	repgen.logger.Debug("Pooled MessagePrepare", commontypes.LogFields{
		"sender": sender,
		"seqNr":  msg.SeqNr,
	})

	repgen.tryProcessPreparePool()
}

func (repgen *reportGenerationState[RI]) tryProcessPreparePool() {
	if repgen.followerState.phase != repgenFollowerPhaseSentPrepare {
		repgen.logger.Debug("cannot process PreparePool, wrong phase", commontypes.LogFields{
			"phase": repgen.followerState.phase,
		})
		return
	}

	byzQuorumSize := ByzQuorumSize(repgen.config.N(), repgen.config.F)

	poolEntries := repgen.followerState.preparePool.Entries(repgen.followerState.seqNr)
	if len(poolEntries) < byzQuorumSize {

		return
	}

	for sender, preparePoolEntry := range poolEntries {
		if preparePoolEntry == nil {
			continue
		}
		if preparePoolEntry.Verified != nil {
			continue
		}
		err := preparePoolEntry.Item.Verify(
			repgen.Timestamp(),
			repgen.followerState.seqNr,
			repgen.followerState.currentOutcomeInputsDigest,
			repgen.followerState.currentOutcomeDigest,
			repgen.config.OracleIdentities[sender].OffchainPublicKey,
		)
		ok := err == nil
		repgen.followerState.preparePool.StoreVerified(repgen.followerState.seqNr, sender, ok)
		if !ok {
			repgen.logger.Warn("Got invalid MessagePrepare", commontypes.LogFields{
				"sender": sender,
				"seqNr":  repgen.followerState.seqNr,
				"error":  err,
			})
		}
	}

	var prepareQuorumCertificate []AttributedPrepareSignature
	for sender, preparePoolEntry := range poolEntries {
		if preparePoolEntry.Verified != nil && *preparePoolEntry.Verified {
			prepareQuorumCertificate = append(prepareQuorumCertificate, AttributedPrepareSignature{
				preparePoolEntry.Item,
				sender,
			})
			if len(prepareQuorumCertificate) == byzQuorumSize {
				break
			}
		}
	}

	if len(prepareQuorumCertificate) < byzQuorumSize {
		return
	}

	commitSignature, err := MakeCommitSignature(
		repgen.Timestamp(),
		repgen.followerState.seqNr,
		repgen.followerState.currentOutcomeDigest,
		repgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		repgen.logger.Critical("Failed to sign Commit", commontypes.LogFields{
			"error": err,
		})
		return
	}

	repgen.followerState.cert = &CertifiedPrepareOrCommitPrepare{
		repgen.e,
		repgen.followerState.seqNr,
		repgen.followerState.currentOutcomeInputsDigest,
		repgen.followerState.currentOutcome,
		prepareQuorumCertificate,
	}
	if !repgen.persistCert() {
		return
	}

	repgen.followerState.phase = repgenFollowerPhaseSentCommit

	repgen.logger.Debug("Broadcasting MessageCommit", commontypes.LogFields{})
	repgen.netSender.Broadcast(MessageCommit[RI]{
		repgen.e,
		repgen.followerState.seqNr,
		commitSignature,
	})
}

func (repgen *reportGenerationState[RI]) messageCommit(msg MessageCommit[RI], sender commontypes.OracleID) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessageCommit for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if putResult := repgen.followerState.commitPool.Put(msg.SeqNr, sender, msg.Signature); putResult != pool.PutResultOK {
		repgen.logger.Debug("Dropping MessageCommit", commontypes.LogFields{
			"sender": sender,
			"seqNr":  msg.SeqNr,
			"reason": putResult,
		})
		return
	}

	repgen.logger.Debug("Pooled MessageCommit", commontypes.LogFields{
		"sender": sender,
		"seqNr":  msg.SeqNr,
	})

	repgen.tryProcessCommitPool()
}

func (repgen *reportGenerationState[RI]) tryProcessCommitPool() {
	if repgen.followerState.phase != repgenFollowerPhaseSentCommit {
		repgen.logger.Debug("cannot process CommitPool, wrong phase", commontypes.LogFields{
			"phase": repgen.followerState.phase,
		})
		return
	}

	byzQuorumSize := ByzQuorumSize(repgen.config.N(), repgen.config.F)

	poolEntries := repgen.followerState.commitPool.Entries(repgen.followerState.seqNr)
	if len(poolEntries) < byzQuorumSize {

		return
	}

	for sender, commitPoolEntry := range poolEntries {
		if commitPoolEntry == nil {
			continue
		}
		if commitPoolEntry.Verified != nil {
			continue
		}
		err := commitPoolEntry.Item.Verify(
			repgen.Timestamp(),
			repgen.followerState.seqNr,
			repgen.followerState.currentOutcomeDigest,
			repgen.config.OracleIdentities[sender].OffchainPublicKey,
		)
		ok := err == nil
		commitPoolEntry.Verified = &ok
		if !ok {
			repgen.logger.Warn("Got invalid MessageCommit", commontypes.LogFields{
				"sender": sender,
			})
		}
	}

	var commitQuorumCertificate []AttributedCommitSignature
	for sender, commitPoolEntry := range poolEntries {
		if commitPoolEntry.Verified != nil && *commitPoolEntry.Verified {
			commitQuorumCertificate = append(commitQuorumCertificate, AttributedCommitSignature{
				commitPoolEntry.Item,
				sender,
			})
			if len(commitQuorumCertificate) == byzQuorumSize {
				break
			}
		}
	}

	if len(commitQuorumCertificate) < byzQuorumSize {
		return
	}

	repgen.deliver(CertifiedPrepareOrCommitCommit{
		repgen.e,
		repgen.followerState.seqNr,
		repgen.followerState.currentOutcome,
		commitQuorumCertificate,
	})

	if uint64(repgen.config.RMax) <= repgen.followerState.seqNr-repgen.followerState.firstSeqNrOfEpoch+1 {
		repgen.logger.Debug("epoch has been going on for too long, sending EventChangeLeader to Pacemaker", commontypes.LogFields{
			"seqNr": repgen.followerState.seqNr,
		})
		select {
		case repgen.chReportGenerationToPacemaker <- EventChangeLeader[RI]{}:
		case <-repgen.ctx.Done():
			return
		}
		return
	} else {
		repgen.logger.Debug("sending EventProgress to Pacemaker", commontypes.LogFields{
			"seqNr": repgen.followerState.seqNr,
		})
		select {
		case repgen.chReportGenerationToPacemaker <- EventProgress[RI]{}:
		case <-repgen.ctx.Done():
			return
		}
	}

	repgen.startSubsequentFollowerRound()
	if repgen.id == repgen.l {
		repgen.startSubsequentLeaderRound()
	}

	repgen.tryProcessStartRoundPool()
}

func (repgen *reportGenerationState[RI]) deliver(commit CertifiedPrepareOrCommitCommit) {
	if commit.SeqNr < repgen.followerState.deliveredSeqNr {
		repgen.logger.Critical("Assumption violation, commitSeqNr is less than deliveredSeqNr", commontypes.LogFields{
			"commitSeqNr":    commit.SeqNr,
			"deliveredSeqNr": repgen.followerState.deliveredSeqNr,
		})
		return
	}

	if commit.SeqNr <= repgen.followerState.deliveredSeqNr {

		repgen.logger.Debug("Skipping delivery of already delivered outcome", commontypes.LogFields{
			"seqNr":          commit.SeqNr,
			"deliveredSeqNr": repgen.followerState.deliveredSeqNr,
		})
	} else {
		repgen.followerState.cert = &commit
		if !repgen.persistCert() {
			return
		}

		repgen.followerState.deliveredSeqNr = commit.SeqNr
		repgen.followerState.deliveredOutcome = commit.Outcome

		repgen.logger.Debug("✅ Delivered outcome", commontypes.LogFields{
			"seqNr": commit.SeqNr,
		})

		select {
		case repgen.chReportGenerationToReportFinalization <- EventDeliver[RI]{commit}:
		case <-repgen.ctx.Done():
			return
		}
	}

	repgen.followerState.startRoundPool.ReapDelivered(repgen.followerState.deliveredSeqNr)
	repgen.followerState.proposePool.ReapDelivered(repgen.followerState.deliveredSeqNr)
	repgen.followerState.preparePool.ReapDelivered(repgen.followerState.deliveredSeqNr)
	repgen.followerState.commitPool.ReapDelivered(repgen.followerState.deliveredSeqNr)
}

func (repgen *reportGenerationState[RI]) persistCert() (ok bool) {
	ctx, cancel := context.WithTimeout(repgen.ctx, repgen.localConfig.DatabaseTimeout)
	defer cancel()
	if err := repgen.database.WriteCert(ctx, repgen.config.ConfigDigest, repgen.followerState.cert); err != nil {
		repgen.logger.Error("Error persisting cert to database. Cannot safely continue current round.", commontypes.LogFields{
			"error": err,
		})
		return false
	}
	return true
}
