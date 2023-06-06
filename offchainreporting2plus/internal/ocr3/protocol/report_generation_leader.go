package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type repgenLeaderPhase string

const (
	repgenLeaderPhaseUnknown        repgenLeaderPhase = "unknown"
	repgenLeaderPhaseNewEpoch       repgenLeaderPhase = "newEpoch"
	repgenLeaderPhaseSentStartEpoch repgenLeaderPhase = "sentStartEpoch"
	repgenLeaderPhaseSentStartRound repgenLeaderPhase = "sentStartRound"
	repgenLeaderPhaseGrace          repgenLeaderPhase = "grace"
	repgenLeaderPhaseSentPropose    repgenLeaderPhase = "sentPropose"
)

func (repgen *reportGenerationState[RI]) messageReconcile(msg MessageReconcile[RI], sender commontypes.OracleID) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessageReconcile for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if repgen.l != repgen.id {
		repgen.logger.Warn("Non-leader received MessageReconcile", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	if repgen.leaderState.phase != repgenLeaderPhaseNewEpoch {
		repgen.logger.Debug("Got MessageReconcile for wrong phase", commontypes.LogFields{
			"sender": sender,
			"phase":  repgen.leaderState.phase,
		})
		return
	}

	{
		err := msg.HighestCertified.Verify(
			repgen.config.ConfigDigest,
			repgen.config.OracleIdentities,
			repgen.config.N(),
			repgen.config.F,
		)
		if err != nil {
			repgen.logger.Warn("MessageReconcile.HighestCertified is invalid", commontypes.LogFields{
				"sender": sender,
				"error":  err,
			})
			return
		}
	}

	{
		err := msg.SignedHighestCertifiedTimestamp.Verify(
			repgen.Timestamp(),
			repgen.config.OracleIdentities[sender].OffchainPublicKey,
		)
		if err != nil {
			repgen.logger.Warn("MessageReconcile.SignedHighestCertifiedTimestamp is invalid", commontypes.LogFields{
				"sender": sender,
				"error":  err,
			})
			return
		}
	}

	if msg.HighestCertified.Timestamp() != msg.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp {
		repgen.logger.Warn("Timestamp mismatch in MessageReconcile", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	for _, ashct := range repgen.leaderState.startRoundQuorumCertificate.HighestCertifiedProof {
		if ashct.Signer == sender {
			repgen.logger.Warn("MessageReconcile.HighestCertified is duplicate", commontypes.LogFields{
				"sender": sender,
			})
			return
		}
	}

	repgen.logger.Debug("Received valid MessageReconcile", commontypes.LogFields{
		"sender":                    sender,
		"highestCertifiedTimestamp": msg.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp,
	})

	repgen.leaderState.startRoundQuorumCertificate.HighestCertifiedProof = append(repgen.leaderState.startRoundQuorumCertificate.HighestCertifiedProof, AttributedSignedHighestCertifiedTimestamp{
		msg.SignedHighestCertifiedTimestamp,
		sender,
	})

	if repgen.leaderState.startRoundQuorumCertificate.HighestCertified == nil || repgen.leaderState.startRoundQuorumCertificate.HighestCertified.Timestamp().Less(msg.HighestCertified.Timestamp()) {
		repgen.leaderState.startRoundQuorumCertificate.HighestCertified = msg.HighestCertified
	}

	if len(repgen.leaderState.startRoundQuorumCertificate.HighestCertifiedProof) == ByzQuorumSize(repgen.config.N(), repgen.config.F) {
		if err := repgen.leaderState.startRoundQuorumCertificate.Verify(repgen.Timestamp(), repgen.config.OracleIdentities, repgen.config.N(), repgen.config.F); err != nil {
			repgen.logger.Critical("StartRoundQuorumCertificate is invalid, very surprising!", commontypes.LogFields{
				"qc": repgen.leaderState.startRoundQuorumCertificate,
			})
			return
		}

		repgen.leaderState.phase = repgenLeaderPhaseSentStartEpoch

		repgen.logger.Info("Broadcasting MessageStartEpoch", nil)

		repgen.netSender.Broadcast(MessageStartEpoch[RI]{
			repgen.e,
			repgen.leaderState.startRoundQuorumCertificate,
		})

		if repgen.leaderState.startRoundQuorumCertificate.HighestCertified.IsGenesis() {
			repgen.followerState.firstSeqNrOfEpoch = repgen.followerState.deliveredSeqNr + 1
			repgen.startSubsequentLeaderRound()
		} else if commitQC, ok := repgen.leaderState.startRoundQuorumCertificate.HighestCertified.(*CertifiedPrepareOrCommitCommit); ok {
			repgen.deliver(*commitQC)
			repgen.followerState.firstSeqNrOfEpoch = repgen.followerState.deliveredSeqNr + 1
			repgen.startSubsequentLeaderRound()
		} else {
			prepareQc := repgen.leaderState.startRoundQuorumCertificate.HighestCertified.(*CertifiedPrepareOrCommitPrepare)
			repgen.followerState.firstSeqNrOfEpoch = prepareQc.SeqNr + 1
			// We're dealing with a re-proposal from a failed epoch based on a
			// prepare qc.
			// We don't want to send OBSERVER-REQ.
		}
	}
}

func (repgen *reportGenerationState[RI]) eventTRoundTimeout() {
	repgen.logger.Debug("TRound fired", commontypes.LogFields{
		"deltaRoundMilliseconds": repgen.config.DeltaRound.Milliseconds(),
	})
	repgen.startSubsequentLeaderRound()
}

func (repgen *reportGenerationState[RI]) startSubsequentLeaderRound() {
	if !repgen.leaderState.readyToStartRound {
		repgen.leaderState.readyToStartRound = true
		return
	}

	query, ok := callPlugin[types.Query](
		repgen,
		"Query",
		repgen.config.MaxDurationQuery,
		repgen.OutcomeCtx(repgen.followerState.deliveredSeqNr+1),
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
			return repgen.reportingPlugin.Query(ctx, outctx)
		},
	)
	if !ok {
		return
	}

	repgen.leaderState.query = query

	repgen.leaderState.observations = map[commontypes.OracleID]*SignedObservation{}

	repgen.leaderState.tRound = time.After(repgen.config.DeltaRound)
	repgen.leaderState.readyToStartRound = false

	repgen.leaderState.phase = repgenLeaderPhaseSentStartRound
	repgen.logger.Debug("Broadcasting MessageStartRound", commontypes.LogFields{
		"seqNr": repgen.followerState.deliveredSeqNr + 1,
	})
	repgen.netSender.Broadcast(MessageStartRound[RI]{
		repgen.e,
		repgen.followerState.deliveredSeqNr + 1,
		query, // query
	})
}

func (repgen *reportGenerationState[RI]) messageObserve(msg MessageObserve[RI], sender commontypes.OracleID) {

	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessageObserve for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if repgen.l != repgen.id {
		repgen.logger.Warn("Non-leader received MessageObserve", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	if repgen.leaderState.phase != repgenLeaderPhaseSentStartRound && repgen.leaderState.phase != repgenLeaderPhaseGrace {
		repgen.logger.Debug("Got MessageObserve for wrong phase", commontypes.LogFields{
			"sender": sender,
			"phase":  repgen.leaderState.phase,
		})
		return
	}

	if msg.SeqNr != repgen.followerState.seqNr {
		repgen.logger.Debug("Got MessageObserve with invalid SeqNr", commontypes.LogFields{
			"sender":   sender,
			"msgSeqNr": msg.SeqNr,
			"seqNr":    repgen.followerState.seqNr,
		})
		return
	}

	if repgen.leaderState.observations[sender] != nil {
		repgen.logger.Warn("Got duplicate MessageObserve", commontypes.LogFields{
			"sender": sender,
			"seqNr":  repgen.followerState.seqNr,
		})
		return
	}

	if err := msg.SignedObservation.Verify(repgen.Timestamp(), repgen.leaderState.query, repgen.config.OracleIdentities[sender].OffchainPublicKey); err != nil {
		repgen.logger.Warn("MessageObserve carries invalid SignedObservation", commontypes.LogFields{
			"sender": sender,
			"error":  err,
		})
		return
	}

	repgen.logger.Debug("Got valid MessageObserve", commontypes.LogFields{
		"seqNr": repgen.followerState.seqNr,
	})

	repgen.leaderState.observations[sender] = &msg.SignedObservation

	observationCount := 0
	for _, so := range repgen.leaderState.observations {
		if so != nil {
			observationCount++
		}
	}
	if observationCount == 2*repgen.config.F+1 {
		repgen.logger.Debug("starting observation grace period", commontypes.LogFields{})
		repgen.leaderState.phase = repgenLeaderPhaseGrace
		repgen.leaderState.tGrace = time.After(repgen.config.DeltaGrace)
	}
}

func (repgen *reportGenerationState[RI]) eventTGraceTimeout() {
	if repgen.leaderState.phase != repgenLeaderPhaseGrace {
		repgen.logger.Error("leader's phase conflicts tGrace timeout", commontypes.LogFields{
			"phase": repgen.leaderState.phase,
		})
		return
	}
	asos := []AttributedSignedObservation{}
	for oid, so := range repgen.leaderState.observations {
		if so != nil {
			asos = append(asos, AttributedSignedObservation{
				*so,
				commontypes.OracleID(oid),
			})
		}
	}

	repgen.leaderState.phase = repgenLeaderPhaseSentPropose

	repgen.logger.Debug("Broadcasting MessagePropose", commontypes.LogFields{})
	repgen.netSender.Broadcast(MessagePropose[RI]{
		repgen.e,
		repgen.followerState.seqNr,
		asos,
	})
}
