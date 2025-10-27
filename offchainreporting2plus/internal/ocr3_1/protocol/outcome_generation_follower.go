package protocol

import (
	"context"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/pool"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type outgenFollowerPhase string

const (
	outgenFollowerPhaseUnknown                           outgenFollowerPhase = "unknown"
	outgenFollowerPhaseNewEpoch                          outgenFollowerPhase = "newEpoch"
	outgenFollowerPhaseEpochStartReceived                outgenFollowerPhase = "epochStartReceived"
	outgenFollowerPhaseNewRound                          outgenFollowerPhase = "newRound"
	outgenFollowerPhaseBackgroundObservation             outgenFollowerPhase = "backgroundObservation"
	outgenFollowerPhaseSentObservation                   outgenFollowerPhase = "sentObservation"
	outgenFollowerPhaseBackgroundProposalStateTransition outgenFollowerPhase = "backgroundProposalStateTransition"
	outgenFollowerPhaseSentPrepare                       outgenFollowerPhase = "sentPrepare"
	outgenFollowerPhaseSentCommit                        outgenFollowerPhase = "sentCommit"
	outgenFollowerPhaseBackgroundCommitted               outgenFollowerPhase = "backgroundCommitted"
)

func (outgen *outcomeGenerationState[RI]) eventTInitialTimeout() {
	outgen.logger.Debug("TInitial fired", commontypes.LogFields{
		"seqNr":        outgen.sharedState.seqNr,
		"deltaInitial": outgen.config.GetDeltaInitial().String(),
	})
	select {
	case outgen.chOutcomeGenerationToPacemaker <- EventNewEpochRequest[RI]{}:
	case <-outgen.ctx.Done():
		return
	}
}

func (outgen *outcomeGenerationState[RI]) messageEpochStart(msg MessageEpochStart[RI], sender commontypes.OracleID) {
	outgen.logger.Debug("received MessageEpochStart", commontypes.LogFields{
		"sender":                       sender,
		"msgEpoch":                     msg.Epoch,
		"msgHighestCertifiedTimestamp": msg.EpochStartProof.HighestCertified.Timestamp(),
		"msgAbdicate":                  msg.Abdicate,
	})

	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageEpochStart for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if sender != outgen.sharedState.l {
		outgen.logger.Warn("dropping MessageEpochStart from non-leader", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	if outgen.followerState.phase != outgenFollowerPhaseNewEpoch {
		outgen.logger.Warn("dropping MessageEpochStart for wrong phase", commontypes.LogFields{
			"phase": outgen.followerState.phase,
		})
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseEpochStartReceived

	{
		err := msg.EpochStartProof.Verify(
			outgen.ID(),
			outgen.config.OracleIdentities,
			outgen.config.ByzQuorumSize(),
		)
		if err != nil {
			outgen.logger.Warn("dropping MessageEpochStart containing invalid StartRoundQuorumCertificate", commontypes.LogFields{
				"error": err,
			})
			return
		}
	}

	outgen.followerState.tInitial = nil
	outgen.followerState.leaderAbdicated = msg.Abdicate

	outgen.refreshCommittedSeqNrAndCert() // call in case stasy has made progress in the meanwhile
	outgen.sendStateSyncRequestFromCertifiedPrepareOrCommit(msg.EpochStartProof.HighestCertified)

	if msg.EpochStartProof.HighestCertified.IsGenesis() {
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentFollowerRound()
	} else if commitQC, ok := msg.EpochStartProof.HighestCertified.(*CertifiedCommit); ok {
		outgen.tryToMoveCertAndKVStateToCommitQC(commitQC)
		if outgen.sharedState.committedSeqNr != commitQC.CommitSeqNr {
			outgen.logger.Warn("cannot process MessageEpochStart, mismatching committedSeqNr, will not be able to participate in epoch", commontypes.LogFields{
				"commitSeqNr":    commitQC.CommitSeqNr,
				"committedSeqNr": outgen.sharedState.committedSeqNr,
			})
			return
		}
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		if msg.Abdicate {
			outgen.sendNewEpochRequestToPacemakerDueToLeaderAbdication()
			return
		}
		outgen.startSubsequentFollowerRound()
	} else {
		// We're dealing with a re-proposal from a failed epoch
		prepareQc, ok := msg.EpochStartProof.HighestCertified.(*CertifiedPrepare)
		if !ok {
			outgen.logger.Critical("cast to CertifiedPrepare failed while processing MessageEpochStart", commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
			})
			return
		}

		if prepareQc.SeqNr() < outgen.sharedState.committedSeqNr {
			outgen.logger.Warn("cannot process MessageEpochStart, prepareQC seqNr is less than committedSeqNr, will not be able to participate in epoch", commontypes.LogFields{
				"prepareQCSeqNr": prepareQc.SeqNr(),
				"committedSeqNr": outgen.sharedState.committedSeqNr,
			})
			return
		}

		outgen.followerState.stateTransitionInfo = stateTransitionInfoDigests{
			prepareQc.StateTransitionInputsDigest,
			prepareQc.StateTransitionOutputsDigest,
			prepareQc.StateRootDigest,
			prepareQc.ReportsPlusPrecursorDigest,
		}

		outgen.sharedState.firstSeqNrOfEpoch = prepareQc.SeqNr() + 1
		outgen.sharedState.seqNr = prepareQc.SeqNr()
		outgen.sharedState.observationQuorum = nil

		outgen.followerState.query = nil
		outgen.ensureOpenKVTransactionDiscarded()

		outgen.broadcastMessagePrepare()
	}
}

func (outgen *outcomeGenerationState[RI]) startSubsequentFollowerRound() {
	outgen.sharedState.seqNr = outgen.sharedState.committedSeqNr + 1
	outgen.sharedState.observationQuorum = nil

	outgen.followerState.phase = outgenFollowerPhaseNewRound
	outgen.followerState.query = nil
	outgen.followerState.stateTransitionInfo = stateTransitionInfoDigests{}
	outgen.ensureOpenKVTransactionDiscarded()
	outgen.logger.Debug("starting new follower round", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})

	outgen.tryProcessRoundStartPool()
}

func (outgen *outcomeGenerationState[RI]) messageRoundStart(msg MessageRoundStart[RI], sender commontypes.OracleID) {
	outgen.logger.Debug("received MessageRoundStart", commontypes.LogFields{
		"sender":   sender,
		"msgSeqNr": msg.SeqNr,
		"msgEpoch": msg.Epoch,
	})

	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageRoundStart for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if sender != outgen.sharedState.l {
		outgen.logger.Warn("dropping MessageRoundStart from non-leader", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if putResult := outgen.followerState.roundStartPool.Put(msg.SeqNr, sender, msg); putResult != pool.PutResultOK {
		outgen.logger.Debug("dropping MessageRoundStart", commontypes.LogFields{
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
			"reason":   putResult,
		})
		return
	}

	outgen.logger.Debug("pooled MessageRoundStart", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})

	outgen.tryProcessRoundStartPool()
}

func (outgen *outcomeGenerationState[RI]) tryProcessRoundStartPool() {
	if outgen.followerState.phase != outgenFollowerPhaseNewRound {
		outgen.logger.Debug("cannot process RoundStartPool, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	poolEntries := outgen.followerState.roundStartPool.Entries(outgen.sharedState.seqNr)

	if poolEntries == nil || poolEntries[outgen.sharedState.l] == nil {

		outgen.logger.Debug("cannot process RoundStartPool, it's empty", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		return
	}

	if outgen.followerState.query != nil {
		outgen.logger.Warn("cannot process RoundStartPool, query already set", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		return
	}

	msg := poolEntries[outgen.sharedState.l].Item
	roundCtx := outgen.RoundCtx(outgen.sharedState.seqNr)

	outgen.followerState.phase = outgenFollowerPhaseBackgroundObservation
	outgen.followerState.query = &msg.Query

	outgen.telemetrySender.RoundStarted(
		outgen.config.ConfigDigest,
		roundCtx.Epoch,
		roundCtx.SeqNr,
		roundCtx.Round,
		outgen.sharedState.l,
	)

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		aq := types.AttributedQuery{
			*outgen.followerState.query,
			outgen.sharedState.l,
		}
		kvReadTxn, err := outgen.kvDb.NewReadTransaction(roundCtx.SeqNr)

		if err != nil {
			outgen.logger.Warn("failed to create new transaction, aborting tryProcessRoundStartPool", commontypes.LogFields{
				"seqNr": roundCtx.SeqNr,
				"error": err,
			})
			return
		}
		outgen.subs.Go(func() {
			outgen.backgroundObservation(ctx, logger, roundCtx, aq, kvReadTxn)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundObservation(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	roundCtx RoundContext,
	aq types.AttributedQuery,
	kvReadTxn KeyValueDatabaseReadTransaction,
) {
	observation, ok := callPluginFromOutcomeGenerationBackground[types.Observation](
		ctx,
		logger,
		"Observation",
		outgen.config.WarnDurationObservation,
		roundCtx,
		func(ctx context.Context, roundCtx RoundContext) (types.Observation, error) {
			return outgen.reportingPlugin.Observation(ctx,
				roundCtx.SeqNr,
				aq,
				kvReadTxn,
				NewRoundBlobBroadcastFetcher(
					roundCtx.SeqNr,
					outgen.blobBroadcastFetcher,
				),
			)
		},
	)
	kvReadTxn.Discard()
	if !ok {
		return
	}

	select {
	case outgen.chLocalEvent <- EventComputedObservation[RI]{
		roundCtx.Epoch,
		roundCtx.SeqNr,
		aq,
		observation,
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedObservation(ev EventComputedObservation[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("discarding EventComputedObservation from old round", commontypes.LogFields{
			"seqNr":   outgen.sharedState.seqNr,
			"evEpoch": ev.Epoch,
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.followerState.phase != outgenFollowerPhaseBackgroundObservation {
		outgen.logger.Debug("discarding EventComputedObservation, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	so, err := MakeSignedObservation(outgen.ID(), outgen.sharedState.seqNr, ev.AttributedQuery, ev.Observation, outgen.offchainKeyring.OffchainSign)
	if err != nil {
		outgen.logger.Error("MakeSignedObservation returned error", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	if err := so.Verify(outgen.ID(), outgen.sharedState.seqNr, ev.AttributedQuery, outgen.offchainKeyring.OffchainPublicKey()); err != nil {
		outgen.logger.Error("MakeSignedObservation produced invalid signature", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseSentObservation
	outgen.metrics.sentObservationsTotal.Inc()
	outgen.logger.Debug("sent MessageObservation to leader", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})
	outgen.netSender.SendTo(MessageObservation[RI]{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		so,
	}, outgen.sharedState.l)

	outgen.tryProcessProposalPool()
}

func (outgen *outcomeGenerationState[RI]) messageProposal(msg MessageProposal[RI], sender commontypes.OracleID) {
	outgen.logger.Debug("received MessageProposal", commontypes.LogFields{
		"sender":   sender,
		"msgSeqNr": msg.SeqNr,
		"msgEpoch": msg.Epoch,
	})

	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageProposal for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if sender != outgen.sharedState.l {
		outgen.logger.Warn("dropping MessageProposal from non-leader", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if putResult := outgen.followerState.proposalPool.Put(msg.SeqNr, sender, msg); putResult != pool.PutResultOK {
		outgen.logger.Debug("dropping MessageProposal", commontypes.LogFields{
			"seqNr":        outgen.sharedState.seqNr,
			"messageSeqNr": msg.SeqNr,
			"reason":       putResult,
		})
		return
	}

	outgen.logger.Debug("pooled MessageProposal", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})

	outgen.tryProcessProposalPool()
}

func (outgen *outcomeGenerationState[RI]) tryProcessProposalPool() {
	if outgen.followerState.phase != outgenFollowerPhaseSentObservation {
		outgen.logger.Debug("cannot process ProposalPool, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	poolEntries := outgen.followerState.proposalPool.Entries(outgen.sharedState.seqNr)

	if poolEntries == nil || poolEntries[outgen.sharedState.l] == nil {

		return
	}

	msg := poolEntries[outgen.sharedState.l].Item

	if msg.SeqNr <= outgen.sharedState.committedSeqNr {
		outgen.logger.Critical("MessageProposal contains invalid SeqNr", commontypes.LogFields{
			"msgSeqNr":       msg.SeqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseBackgroundProposalStateTransition

	kvReadWriteTxn, err := outgen.kvDb.NewSerializedReadWriteTransaction(outgen.sharedState.seqNr)
	if err != nil {
		outgen.logger.Warn("could not create kv read/write transaction", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"err":   err,
		})
		return

	}

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		roundCtx := outgen.RoundCtx(outgen.sharedState.seqNr)
		ogid := outgen.ID()
		aq := types.AttributedQuery{
			*outgen.followerState.query,
			outgen.sharedState.l,
		}

		asos := msg.AttributedSignedObservations

		outgen.subs.Go(func() {
			outgen.backgroundProposalStateTransition(
				ctx,
				logger,
				ogid,
				roundCtx,
				aq,
				asos,
				kvReadWriteTxn,
			)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundProposalStateTransition(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	ogid OutcomeGenerationID,
	roundCtx RoundContext,
	aq types.AttributedQuery,
	asos []AttributedSignedObservation,
	kvReadWriteTxn KeyValueDatabaseReadWriteTransaction,
) {
	shouldDiscardKVTxn := true
	defer func() {
		if shouldDiscardKVTxn {
			kvReadWriteTxn.Discard()
		}
	}()

	aos, ok := outgen.backgroundCheckAttributedSignedObservations(ctx, logger, ogid, roundCtx, aq, asos, kvReadWriteTxn)
	if !ok {
		return
	}
	reportsPlusPrecursor, ok := callPluginFromOutcomeGenerationBackground[ocr3_1types.ReportsPlusPrecursor](
		ctx,
		logger,
		"StateTransition",
		outgen.config.WarnDurationStateTransition,
		roundCtx,
		func(ctx context.Context, roundCtx RoundContext) (ocr3_1types.ReportsPlusPrecursor, error) {
			return outgen.reportingPlugin.StateTransition(
				ctx,
				roundCtx.SeqNr,
				aq,
				aos,
				kvReadWriteTxn,
				NewRoundBlobBroadcastFetcher(
					roundCtx.SeqNr,
					outgen.blobBroadcastFetcher,
				),
			)
		},
	)
	if !ok {
		return
	}

	stateTransitionInputsDigest := MakeStateTransitionInputsDigest(
		outgen.config.ConfigDigest,
		roundCtx.SeqNr,
		aq,
		aos,
	)

	writeSet, err := kvReadWriteTxn.GetWriteSet()
	if err != nil {
		outgen.logger.Warn("failed to get write set from kv read/write transaction", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}
	stateRootDigest, err := kvReadWriteTxn.CloseWriteSet()
	if err != nil {
		outgen.logger.Warn("failed to close the transaction WriteSet", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	stateTransitionOutputDigest := MakeStateTransitionOutputDigest(outgen.config.ConfigDigest, roundCtx.SeqNr, writeSet)
	reportsPlusPrecursorDigest := MakeReportsPlusPrecursorDigest(outgen.config.ConfigDigest, roundCtx.SeqNr, reportsPlusPrecursor)

	stateTransitionOutputs := StateTransitionOutputs{writeSet}

	select {
	case outgen.chLocalEvent <- EventComputedProposalStateTransition[RI]{
		roundCtx.Epoch,
		roundCtx.SeqNr,
		kvReadWriteTxn,
		stateTransitionInfoDigestsAndPreimages{
			stateTransitionInfoDigests{
				stateTransitionInputsDigest,
				stateTransitionOutputDigest,
				stateRootDigest,
				reportsPlusPrecursorDigest,
			},
			stateTransitionOutputs,
			reportsPlusPrecursor,
		},
	}:
		shouldDiscardKVTxn = false
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedProposalStateTransition(ev EventComputedProposalStateTransition[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("discarding EventComputedProposalStateTransition from old round", commontypes.LogFields{
			"seqNr":   outgen.sharedState.seqNr,
			"evEpoch": ev.Epoch,
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.followerState.phase != outgenFollowerPhaseBackgroundProposalStateTransition {
		outgen.logger.Debug("discarding EventComputedProposalStateTransition, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	outgen.followerState.openKVTxn = ev.KeyValueDatabaseReadWriteTransaction

	var stidap stateTransitionInfoDigestsAndPreimages
	switch sti := ev.stateTransitionInfo.(type) {
	case stateTransitionInfoDigestsAndPreimages:
		stidap = sti
	case stateTransitionInfoDigests:
		outgen.logger.Critical("assumption violation, EventComputedProposalStateTransition state transition info contains only digests and no preimages", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		return
	}

	outgen.followerState.stateTransitionInfo = stidap

	err := outgen.persistUnattestedStateTransitionBlockAndReportsPlusPrecursor(
		StateTransitionBlock{
			ev.Epoch,
			ev.SeqNr,
			stidap.InputsDigest,
			stidap.Outputs,
			stidap.StateRootDigest,
			stidap.ReportsPlusPrecursorDigest,
		},
		stidap.InputsDigest,
		stidap.ReportsPlusPrecursor,
	)
	if err != nil {
		outgen.logger.Error("failed to persist unattested state transition block and reports plus precursor", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	outgen.broadcastMessagePrepare()
}

// broadcasts MessagePrepare for outgen.sharedState.e, outgen.sharedState.seqNr
// and outgen.followerState.stateTransitionInfo
func (outgen *outcomeGenerationState[RI]) broadcastMessagePrepare() {
	ogid := outgen.ID()
	seqNr := outgen.sharedState.seqNr
	stid := outgen.followerState.stateTransitionInfo.digests()
	prepareSignature, err := MakePrepareSignature(
		ogid,
		seqNr,
		stid.InputsDigest,
		stid.OutputDigest,
		stid.StateRootDigest,
		stid.ReportsPlusPrecursorDigest,
		outgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		outgen.logger.Critical("failed to sign Prepare", commontypes.LogFields{
			"seqNr": seqNr,
			"error": err,
		})
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseSentPrepare

	outgen.logger.Debug("broadcasting MessagePrepare", commontypes.LogFields{
		"seqNr": seqNr,
		"phase": outgen.followerState.phase,
	})
	outgen.netSender.Broadcast(MessagePrepare[RI]{
		outgen.sharedState.e,
		seqNr,
		prepareSignature,
	})
}

func (outgen *outcomeGenerationState[RI]) messagePrepare(msg MessagePrepare[RI], sender commontypes.OracleID) {
	outgen.logger.Debug("received MessagePrepare", commontypes.LogFields{
		"sender":   sender,
		"msgSeqNr": msg.SeqNr,
		"msgEpoch": msg.Epoch,
	})

	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessagePrepare for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if putResult := outgen.followerState.preparePool.Put(msg.SeqNr, sender, msg.Signature); putResult != pool.PutResultOK {
		outgen.logger.Debug("dropping MessagePrepare", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
			"reason":   putResult,
		})
		return
	}

	outgen.logger.Debug("pooled MessagePrepare", commontypes.LogFields{
		"sender":   sender,
		"seqNr":    outgen.sharedState.seqNr,
		"msgSeqNr": msg.SeqNr,
	})

	outgen.tryProcessPreparePool()
}

func (outgen *outcomeGenerationState[RI]) tryProcessPreparePool() {
	if outgen.followerState.phase != outgenFollowerPhaseSentPrepare {
		outgen.logger.Debug("cannot process PreparePool, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	poolEntries := outgen.followerState.preparePool.Entries(outgen.sharedState.seqNr)
	if len(poolEntries) < outgen.config.ByzQuorumSize() {

		return
	}

	stid := outgen.followerState.stateTransitionInfo.digests()

	for sender, preparePoolEntry := range poolEntries {
		if preparePoolEntry.Verified != nil {
			continue
		}
		err := preparePoolEntry.Item.Verify(
			outgen.ID(),
			outgen.sharedState.seqNr,
			stid.InputsDigest,
			stid.OutputDigest,
			stid.StateRootDigest,
			stid.ReportsPlusPrecursorDigest,
			outgen.config.OracleIdentities[sender].OffchainPublicKey,
		)
		ok := err == nil
		outgen.followerState.preparePool.StoreVerified(outgen.sharedState.seqNr, sender, ok)
		if !ok {
			outgen.logger.Warn("dropping invalid MessagePrepare", commontypes.LogFields{
				"sender": sender,
				"seqNr":  outgen.sharedState.seqNr,
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
			if len(prepareQuorumCertificate) == outgen.config.ByzQuorumSize() {
				break
			}
		}
	}

	if len(prepareQuorumCertificate) < outgen.config.ByzQuorumSize() {
		return
	}

	commitSignature, err := MakeCommitSignature(
		outgen.ID(),
		outgen.sharedState.seqNr,
		stid.InputsDigest,
		stid.OutputDigest,
		stid.StateRootDigest,
		stid.ReportsPlusPrecursorDigest,
		outgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		outgen.logger.Critical("failed to sign Commit", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	if !outgen.persistAndUpdateCertIfGreater(&CertifiedPrepare{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		stid.InputsDigest,
		stid.OutputDigest,
		stid.StateRootDigest,
		stid.ReportsPlusPrecursorDigest,
		prepareQuorumCertificate,
	}) {
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseSentCommit

	outgen.logger.Debug("broadcasting MessageCommit", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})
	outgen.netSender.Broadcast(MessageCommit[RI]{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		commitSignature,
	})
}

func (outgen *outcomeGenerationState[RI]) messageCommit(msg MessageCommit[RI], sender commontypes.OracleID) {
	outgen.logger.Debug("received MessageCommit", commontypes.LogFields{
		"sender":   sender,
		"msgSeqNr": msg.SeqNr,
		"msgEpoch": msg.Epoch,
	})

	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageCommit for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if putResult := outgen.followerState.commitPool.Put(msg.SeqNr, sender, msg.Signature); putResult != pool.PutResultOK {
		outgen.logger.Debug("dropping MessageCommit", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
			"reason":   putResult,
		})
		return
	}

	outgen.logger.Debug("pooled MessageCommit", commontypes.LogFields{
		"sender":   sender,
		"seqNr":    outgen.sharedState.seqNr,
		"msgSeqNr": msg.SeqNr,
	})

	outgen.tryProcessCommitPool()
}

func (outgen *outcomeGenerationState[RI]) tryProcessCommitPool() {
	if outgen.followerState.phase != outgenFollowerPhaseSentCommit {
		outgen.logger.Debug("cannot process CommitPool, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	poolEntries := outgen.followerState.commitPool.Entries(outgen.sharedState.seqNr)
	if len(poolEntries) < outgen.config.ByzQuorumSize() {

		return
	}

	stid := outgen.followerState.stateTransitionInfo.digests()

	for sender, commitPoolEntry := range poolEntries {
		if commitPoolEntry.Verified != nil {
			continue
		}
		err := commitPoolEntry.Item.Verify(
			outgen.ID(),
			outgen.sharedState.seqNr,
			stid.InputsDigest,
			stid.OutputDigest,
			stid.StateRootDigest,
			stid.ReportsPlusPrecursorDigest,
			outgen.config.OracleIdentities[sender].OffchainPublicKey,
		)
		ok := err == nil
		commitPoolEntry.Verified = &ok
		if !ok {
			outgen.logger.Warn("dropping invalid MessageCommit", commontypes.LogFields{
				"sender": sender,
				"seqNr":  outgen.sharedState.seqNr,
				"error":  err,
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
			if len(commitQuorumCertificate) == outgen.config.ByzQuorumSize() {
				break
			}
		}
	}

	if len(commitQuorumCertificate) < outgen.config.ByzQuorumSize() {
		return
	}

	if outgen.id == outgen.sharedState.l {
		outgen.metrics.ledCommittedRoundsTotal.Inc()
	}

	commitQC := &CertifiedCommit{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		stid.InputsDigest,
		stid.OutputDigest,
		stid.StateRootDigest,
		stid.ReportsPlusPrecursorDigest,
		commitQuorumCertificate,
	}

	switch sti := outgen.followerState.stateTransitionInfo.(type) {
	case stateTransitionInfoDigests:
		// We re-prepared
		outgen.tryToMoveCertAndKVStateToCommitQC(commitQC)
	case stateTransitionInfoDigestsAndPreimages:
		// Regular round progression, we already should have an open transaction
		persistedCert := outgen.commit(*commitQC)
		if !persistedCert {
			outgen.logger.Error("commit() failed to persist cert", commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
			})
			return
		}

		if outgen.followerState.openKVTxn == nil {
			outgen.logger.Error("no open kv transaction, unexpected", commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
				"phase": outgen.followerState.phase,
			})
			return
		}

		// Write attested state transition block
		{
			stb := StateTransitionBlock{
				commitQC.CommitEpoch,
				commitQC.CommitSeqNr,
				sti.InputsDigest,
				sti.Outputs,
				sti.StateRootDigest,
				sti.ReportsPlusPrecursorDigest,
			}
			astb := AttestedStateTransitionBlock{
				stb,
				commitQC.CommitQuorumCertificate,
			}

			err := outgen.followerState.openKVTxn.WriteAttestedStateTransitionBlock(commitQC.CommitSeqNr, astb)
			if err != nil {
				outgen.logger.Error("error writing attested state transition block", commontypes.LogFields{
					"seqNr": commitQC.CommitSeqNr,
					"error": err,
				})
				return
			}
		}

		err := outgen.followerState.openKVTxn.Commit()
		outgen.followerState.openKVTxn.Discard()
		outgen.followerState.openKVTxn = nil
		if err != nil {
			outgen.logger.Warn("failed to commit kv transaction", commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
				"error": err,
			})

			{
				kvSeqNr, err := outgen.kvDb.HighestCommittedSeqNr()
				if err != nil {
					outgen.logger.Error("failed to validate kv commit post-condition, upon kv commit failure", commontypes.LogFields{
						"seqNr": outgen.sharedState.seqNr,
						"error": err,
					})
					return
				}

				if kvSeqNr < outgen.sharedState.seqNr {
					outgen.logger.Error("kv commit failed and post-condition (seqNr <= kvSeqNr) is not satisfied", commontypes.LogFields{
						"seqNr":   outgen.sharedState.seqNr,
						"kvSeqNr": kvSeqNr,
					})
					return
				}
			}
		}
	}

	if outgen.sharedState.committedSeqNr != commitQC.CommitSeqNr {
		outgen.logger.Warn("could not move committed seq nr to commit qc, abandoning epoch", commontypes.LogFields{
			"commitSeqNr":    commitQC.CommitSeqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
		return
	}

	kvReadTxn, err := outgen.kvDb.NewReadTransaction(outgen.sharedState.seqNr + 1)
	if err != nil {
		outgen.logger.Warn("skipping call to ReportingPlugin.Committed", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		outgen.completeRound()
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseBackgroundCommitted

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		roundCtx := RoundContext{
			outgen.sharedState.seqNr,
			outgen.sharedState.e,
			outgen.sharedState.seqNr - outgen.sharedState.firstSeqNrOfEpoch + 1,
		}
		kvReadTxn := kvReadTxn
		outgen.subs.Go(func() {
			outgen.backgroundCommitted(
				ctx,
				logger,
				roundCtx,
				kvReadTxn,
			)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundCommitted(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	roundCtx RoundContext,
	kvReadTxn KeyValueDatabaseReadTransaction,
) {
	_, ok := callPluginFromOutcomeGenerationBackground[error](
		ctx,
		logger,
		"Committed",
		outgen.config.WarnDurationCommitted,
		roundCtx,
		func(ctx context.Context, roundCtx RoundContext) (error, error) {
			return outgen.reportingPlugin.Committed(ctx, roundCtx.SeqNr, kvReadTxn), nil
		},
	)
	kvReadTxn.Discard()

	if !ok {
		outgen.logger.Info("continuing after ReportingPlugin.Committed returned an error", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})

	}

	select {
	case outgen.chLocalEvent <- EventComputedCommitted[RI]{
		roundCtx.Epoch,
		roundCtx.SeqNr,
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedCommitted(ev EventComputedCommitted[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("discarding EventComputedCommitted from old round", commontypes.LogFields{
			"seqNr":   outgen.sharedState.seqNr,
			"evEpoch": ev.Epoch,
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.followerState.phase != outgenFollowerPhaseBackgroundCommitted {
		outgen.logger.Debug("discarding EventComputedCommitted, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	outgen.completeRound()
}

func (outgen *outcomeGenerationState[RI]) completeRound() {
	if uint64(outgen.config.RMax) <= outgen.sharedState.seqNr-outgen.sharedState.firstSeqNrOfEpoch+1 {
		outgen.logger.Debug("epoch has been going on for too long, sending EventNewEpochRequest to Pacemaker", commontypes.LogFields{
			"firstSeqNrOfEpoch": outgen.sharedState.firstSeqNrOfEpoch,
			"seqNr":             outgen.sharedState.seqNr,
			"rMax":              outgen.config.RMax,
		})
		select {
		case outgen.chOutcomeGenerationToPacemaker <- EventNewEpochRequest[RI]{}:
		case <-outgen.ctx.Done():
			return
		}
		return
	} else if outgen.followerState.leaderAbdicated {
		outgen.sendNewEpochRequestToPacemakerDueToLeaderAbdication()
		return
	} else {
		outgen.logger.Debug("sending EventProgress to Pacemaker", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		select {
		case outgen.chOutcomeGenerationToPacemaker <- EventProgress[RI]{}:
		case <-outgen.ctx.Done():
			return
		}
	}

	outgen.startSubsequentFollowerRound()
	if outgen.id == outgen.sharedState.l {
		outgen.startSubsequentLeaderRound()
	}

	outgen.tryProcessRoundStartPool()
}

func (outgen *outcomeGenerationState[RI]) sendNewEpochRequestToPacemakerDueToLeaderAbdication() {
	outgen.logger.Debug("leader abdicated in MessageEpochStart, sending EventNewEpochRequest to Pacemaker", commontypes.LogFields{
		"firstSeqNrOfEpoch": outgen.sharedState.firstSeqNrOfEpoch,
		"seqNr":             outgen.sharedState.seqNr,
	})
	select {
	case outgen.chOutcomeGenerationToPacemaker <- EventNewEpochRequest[RI]{}:
	case <-outgen.ctx.Done():
		return
	}
}

func (outgen *outcomeGenerationState[RI]) commit(commit CertifiedCommit) bool {
	if commit.SeqNr() <= outgen.sharedState.committedSeqNr {

		outgen.logger.Debug("skipping commit of already committed seqNr", commontypes.LogFields{
			"commitSeqNr":    commit.SeqNr(),
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
	} else { // commit.SeqNr >= outgen.sharedState.committedSeqNr + 1

		if !outgen.persistAndUpdateCertIfGreater(&commit) {
			return false
		}

		outgen.sharedState.committedSeqNr = commit.SeqNr()
		outgen.metrics.committedSeqNr.Set(float64(commit.SeqNr()))

		outgen.logger.Debug("✅ committed", commontypes.LogFields{
			"seqNr": commit.SeqNr(),
		})

		select {
		case outgen.chOutcomeGenerationToReportAttestation <- EventNewCertifiedCommit[RI]{
			commit.SeqNr(),
			commit.ReportsPlusPrecursorDigest,
		}:
		case <-outgen.ctx.Done():
		}
	}

	outgen.followerState.roundStartPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.proposalPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.preparePool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.commitPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	return true
}

// Updates and persists cert if it is greater than the current cert.
// Returns false if the cert could not be persisted, in which case the round should be aborted.
func (outgen *outcomeGenerationState[RI]) persistAndUpdateCertIfGreater(cert CertifiedPrepareOrCommit) (ok bool) {
	if outgen.followerState.cert.Timestamp().Less(cert.Timestamp()) {
		ctx, cancel := context.WithTimeout(outgen.ctx, outgen.localConfig.DatabaseTimeout)
		defer cancel()
		if err := outgen.database.WriteCert(ctx, outgen.config.ConfigDigest, cert); err != nil {
			outgen.logger.Error("error persisting cert to database, cannot safely continue current round", commontypes.LogFields{
				"seqNr":                    outgen.sharedState.seqNr,
				"lastWrittenCertTimestamp": outgen.followerState.cert.Timestamp(),
				"certTimestamp":            cert.Timestamp(),
				"error":                    err,
			})
			return false
		}

		outgen.followerState.cert = cert
	}
	return true
}

func (outgen *outcomeGenerationState[RI]) backgroundCheckAttributedSignedObservation(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	ogid OutcomeGenerationID,
	roundCtx RoundContext,
	aq types.AttributedQuery,
	aso AttributedSignedObservation,
	kvReader ocr3_1types.KeyValueStateReader, // we don't discard the kvReader in this function because it is managed further up the call stack
) bool {
	if err := aso.SignedObservation.Verify(ogid, roundCtx.SeqNr, aq, outgen.config.OracleIdentities[aso.Observer].OffchainPublicKey); err != nil {
		logger.Warn("dropping MessageProposal that contains signed observation with invalid signature", commontypes.LogFields{
			"seqNr": roundCtx.SeqNr,
			"error": err,
		})
		return false
	}

	err, ok := callPluginFromOutcomeGenerationBackground[error](
		ctx,
		logger,
		"ValidateObservation",
		outgen.config.WarnDurationValidateObservation,
		roundCtx,
		func(ctx context.Context, roundCtx RoundContext) (error, error) {
			return outgen.reportingPlugin.ValidateObservation(
				ctx,
				roundCtx.SeqNr,
				aq,
				types.AttributedObservation{
					aso.SignedObservation.Observation,
					aso.Observer,
				},
				kvReader,
				NewRoundBlobBroadcastFetcher(
					roundCtx.SeqNr,
					outgen.blobBroadcastFetcher,
				),
			), nil
		},
	)

	if !ok {
		logger.Error("dropping MessageProposal containing observation that could not be validated", commontypes.LogFields{
			"seqNr":    roundCtx.SeqNr,
			"observer": aso.Observer,
		})
		return false
	}

	if err != nil {
		logger.Warn("dropping MessageProposal that contains an invalid observation", commontypes.LogFields{
			"seqNr":    roundCtx.SeqNr,
			"error":    err,
			"observer": aso.Observer,
		})
		return false
	}

	return true
}

// If the attributed signed observations have valid signature, and they satisfy ValidateObservation
// and ObservationQuorum plugin methods, this function returns the vector of corresponding
// AttributedObservations and true.
func (outgen *outcomeGenerationState[RI]) backgroundCheckAttributedSignedObservations(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	ogid OutcomeGenerationID,
	roundCtx RoundContext,
	aq types.AttributedQuery,
	asos []AttributedSignedObservation,
	kvReader ocr3_1types.KeyValueStateReader, // we don't discard the kvReader in this function because it is managed further up the call stack
) ([]types.AttributedObservation, bool) {

	attributedObservations := make([]types.AttributedObservation, 0, len(asos))

	subs, allValidMutex, allValid := subprocesses.Subprocesses{}, sync.Mutex{}, true

	myObservationIncluded := false

	for i, aso := range asos {
		if !(0 <= int(aso.Observer) && int(aso.Observer) < outgen.config.N()) {
			logger.Warn("dropping MessageProposal that contains signed observation with invalid observer", commontypes.LogFields{
				"seqNr":           roundCtx.SeqNr,
				"invalidObserver": aso.Observer,
			})
			return nil, false
		}

		if i > 0 && !(asos[i-1].Observer < aso.Observer) {
			logger.Warn("dropping MessageProposal that contains duplicate signed observation", commontypes.LogFields{
				"seqNr": roundCtx.SeqNr,
			})
			return nil, false
		}

		if aso.Observer == outgen.id {
			myObservationIncluded = true
		}

		attributedObservations = append(attributedObservations, types.AttributedObservation{
			aso.SignedObservation.Observation,
			aso.Observer,
		})

		subs.Go(func() {
			if !outgen.backgroundCheckAttributedSignedObservation(ctx, logger, ogid, roundCtx, aq, aso, kvReader) {
				allValidMutex.Lock()
				allValid = false
				allValidMutex.Unlock()
			}
		})
	}

	subs.Wait()
	if !allValid {
		// no need to log, since backgroundCheckAttributedSignedObservation will already have done so
		return nil, false
	}

	observationQuorum, ok := callPluginFromOutcomeGenerationBackground[bool](
		ctx,
		logger,
		"ObservationQuorum",
		outgen.config.WarnDurationObservationQuorum,
		roundCtx,
		func(ctx context.Context, roundCtx RoundContext) (bool, error) {
			return outgen.reportingPlugin.ObservationQuorum(
				ctx,
				roundCtx.SeqNr,
				aq,
				attributedObservations,
				kvReader,
				NewRoundBlobBroadcastFetcher(
					roundCtx.SeqNr,
					outgen.blobBroadcastFetcher,
				),
			)
		},
	)

	if !ok {
		return nil, false
	}

	if !observationQuorum {
		logger.Warn("dropping MessageProposal that doesn't achieve observation quorum", commontypes.LogFields{
			"seqNr": roundCtx.SeqNr,
		})
		return nil, false
	}

	if myObservationIncluded {
		outgen.metrics.includedObservationsTotal.Inc()
	}

	return attributedObservations, true
}

func (outgen *outcomeGenerationState[RI]) refreshCommittedSeqNrAndCert() {
	tx, err := outgen.kvDb.NewReadTransactionUnchecked()
	if err != nil {
		outgen.logger.Error("error creating read transaction", commontypes.LogFields{
			"error": err,
		})
		return
	}
	defer tx.Discard()

	committedKVSeqNr, err := tx.ReadHighestCommittedSeqNr()
	if err != nil {
		outgen.logger.Error("tx.ReadHighestCommittedSeqNr() failed during refresh", commontypes.LogFields{
			"committedSeqNr": outgen.sharedState.committedSeqNr,
			"error":          err,
		})
		return
	}

	logger := outgen.logger.MakeChild(commontypes.LogFields{
		"committedKVSeqNr": committedKVSeqNr,
		"committedSeqNr":   outgen.sharedState.committedSeqNr,
	})

	if committedKVSeqNr <= outgen.sharedState.committedSeqNr {
		logger.Info("refresh: committedKVSeqNr <= outgen.sharedState.committedSeqNr, nothing to do", nil)

		return
	}

	astb, err := tx.ReadAttestedStateTransitionBlock(committedKVSeqNr)
	if err != nil {
		logger.Error("error reading attested state transition block during refresh", commontypes.LogFields{
			"error": err,
		})
		return
	}
	if astb.StateTransitionBlock.SeqNr() == 0 { // The block does not exist in the database
		logger.Critical("assumption violation, attested state transition block for kv committed seq nr does not exist", nil)
		panic("")
	}
	if astb.StateTransitionBlock.SeqNr() != committedKVSeqNr {
		logger.Critical("assumption violation, attested state transition block has unexpected seq nr", commontypes.LogFields{
			"expectedSeqNr": committedKVSeqNr,
			"actualSeqNr":   astb.StateTransitionBlock.SeqNr(),
		})
		panic("")
	}

	persistedCert := outgen.commit(astb.ToCertifiedCommit(outgen.config.ConfigDigest))

	if !persistedCert {
		logger.Warn("outgen.commit() failed, aborting refresh", nil)
		return
	}

	if outgen.sharedState.committedSeqNr != committedKVSeqNr {
		logger.Critical("assumption violation, outgen local committed seq nr did not progress even though commit() succeeded", commontypes.LogFields{
			"expectedCommittedSeqNr": committedKVSeqNr,
			"actualCommittedSeqNr":   outgen.sharedState.committedSeqNr,
		})
		panic("")
	}

	logger.Debug("refreshed committed seq nr and cert", nil)
}

func (outgen *outcomeGenerationState[RI]) ensureOpenKVTransactionDiscarded() {
	if outgen.followerState.openKVTxn != nil {
		outgen.logger.Warn("discarding open transaction from probably previously failed round", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"spicy": "🫑",
		})
		outgen.followerState.openKVTxn.Discard()
		outgen.followerState.openKVTxn = nil
	}
}
