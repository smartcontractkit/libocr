package protocol

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/pool"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type outgenFollowerPhase string

const (
	outgenFollowerPhaseUnknown                   outgenFollowerPhase = "unknown"
	outgenFollowerPhaseNewEpoch                  outgenFollowerPhase = "newEpoch"
	outgenFollowerPhaseNewRound                  outgenFollowerPhase = "newRound"
	outgenFollowerPhaseBackgroundObservation     outgenFollowerPhase = "backgroundObservation"
	outgenFollowerPhaseSentObservation           outgenFollowerPhase = "sentObservation"
	outgenFollowerPhaseBackgroundStateTransition outgenFollowerPhase = "backgroundProposalOutcome"
	outgenFollowerPhaseSentPrepare               outgenFollowerPhase = "sentPrepare"
	outgenFollowerPhaseSentCommit                outgenFollowerPhase = "sentCommit"
)

func (outgen *outcomeGenerationState[RI]) eventTInitialTimeout() {
	outgen.logger.Debug("TInitial fired", commontypes.LogFields{
		"seqNr":        outgen.sharedState.seqNr,
		"deltaInitial": outgen.config.DeltaInitial.String(),
	})
	select {
	case outgen.chOutcomeGenerationToPacemaker <- EventNewEpochRequest[RI]{}:
	case <-outgen.ctx.Done():
		return
	}
}

func (outgen *outcomeGenerationState[RI]) messageEpochStart(msg MessageEpochStart[RI], sender commontypes.OracleID) {
	outgen.logger.Debug("received MessageEpochStart", commontypes.LogFields{
		"sender":   sender,
		"msgEpoch": msg.Epoch,
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

	{
		err := msg.EpochStartProof.Verify(
			outgen.ID(outgen.sharedState.e),
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

	if msg.EpochStartProof.HighestCertified.IsGenesis() {
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedToKVStoreSeqNr + 1
		outgen.startSubsequentFollowerRound()
	} else if commitQC, ok := msg.EpochStartProof.HighestCertified.(*CertifiedCommit); ok {
		outgen.sharedState.firstSeqNrOfEpoch = commitQC.SeqNr() + 1
		if commitQC.SeqNr() == outgen.sharedState.committedToKVStoreSeqNr {
			outgen.logger.Debug("starting new epoch gracefully",
				commontypes.LogFields{
					"committedToKVStoreSeqNr": outgen.sharedState.committedToKVStoreSeqNr,
					"committedSeqNr":          outgen.sharedState.committedSeqNr,
					"commitQCSeqNr":           commitQC.SeqNr(),
				})
			outgen.startSubsequentFollowerRound()
			outgen.tryProcessRoundStartPool()
		} else {
			outgen.logger.Debug("trying to start a new epoch after replaying the highest committed state transition",
				commontypes.LogFields{
					"committedToKVStoreSeqNr": outgen.sharedState.committedToKVStoreSeqNr,
					"committedSeqNr":          outgen.sharedState.committedSeqNr,
					"commitQCSeqNr":           commitQC.SeqNr(),
				})
			outgen.initRound()
			outgen.replayStateTransition(
				commitQC.StateTransitionInputs,
				commitQC.StateTransitionOutputDigest,
				commitQC.ReportsPlusPrecursorDigest,
				commitQC.CommitQuorumCertificate,
			)
		}
	} else {
		// We're dealing with a re-proposal from a failed epoch
		prepareQc, ok := msg.EpochStartProof.HighestCertified.(*CertifiedPrepare)
		if !ok {
			outgen.logger.Critical("cast to CertifiedPrepare failed while processing MessageEpochStart", commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
			})
			return
		}
		outgen.sharedState.firstSeqNrOfEpoch = prepareQc.SeqNr() + 1
		if outgen.sharedState.committedToKVStoreSeqNr+1 != prepareQc.SeqNr() {
			outgen.logger.Debug("cannot start new epoch, we need to state sync first",
				commontypes.LogFields{
					"committedToKVStoreSeqNr": outgen.sharedState.committedToKVStoreSeqNr,
					"committedSeqNr":          outgen.sharedState.committedSeqNr,
					"prepareQCSeqNr":          prepareQc.SeqNr(),
				})
			select {
			case outgen.chOutcomeGenerationToStatePersistence <- EventStateSyncRequest[RI]{prepareQc.SeqNr()}:
			case <-outgen.ctx.Done():
			}
			return
		}
		outgen.logger.Debug("starting new epoch with a re-proposal",
			commontypes.LogFields{
				"committedToKVStoreSeqNr": outgen.sharedState.committedToKVStoreSeqNr,
				"committedSeqNr":          outgen.sharedState.committedSeqNr,
				"prepareQCSeqNr":          prepareQc.SeqNr(),
			})
		outgen.initRound()
		select {
		case outgen.chOutcomeGenerationToStatePersistence <- EventKVTransactionRequest[RI]{
			ocr3_1types.RoundContext{
				prepareQc.StateTransitionInputs.SeqNr,
				prepareQc.StateTransitionInputs.Epoch,
				prepareQc.StateTransitionInputs.Round,
			},
			prepareQc.StateTransitionInputs.Query,
			attributedSignedObservationsFromAttributedObservations(prepareQc.StateTransitionInputs.AttributedObservations),
			true,
			prepareQc.StateTransitionOutputDigest,
			prepareQc.ReportsPlusPrecursorDigest,
			nil,
		}:
		case <-outgen.ctx.Done():
		}
	}
}

func (outgen *outcomeGenerationState[RI]) initRound() {
	outgen.sharedState.seqNr = outgen.sharedState.committedToKVStoreSeqNr + 1
	outgen.sharedState.observationQuorum = nil
	outgen.followerState.query = nil
	outgen.followerState.roundInfo = roundInfo[RI]{}
}

func (outgen *outcomeGenerationState[RI]) startSubsequentFollowerRound() {
	outgen.initRound()
	outgen.followerState.phase = outgenFollowerPhaseNewRound
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
		query := *outgen.followerState.query
		kvReader := outgen.kvStore.GetReader()
		outgen.subs.Go(func() {
			outgen.backgroundObservation(ctx, logger, roundCtx, query, kvReader)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundObservation(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	roundCtx ocr3_1types.RoundContext,
	query types.Query,
	kvReader ocr3_1types.KeyValueReaderDiscardable,
) {
	observation, ok := callPluginFromOutcomeGenerationBackground[types.Observation](
		ctx,
		logger,
		"Observation",
		outgen.config.MaxDurationObservation,
		roundCtx,
		func(ctx context.Context, rondCtx ocr3_1types.RoundContext) (types.Observation, error) {
			return outgen.reportingPlugin.Observation(ctx, roundCtx, query, kvReader)
		},
	)
	kvReader.Discard()
	if !ok {
		return
	}

	select {
	case outgen.chLocalEvent <- EventComputedObservation[RI]{
		roundCtx.Epoch,
		roundCtx.SeqNr,
		query,
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

	so, err := MakeSignedObservation(outgen.ID(outgen.sharedState.e), outgen.sharedState.seqNr, ev.Query, ev.Observation, outgen.offchainKeyring.OffchainSign)
	if err != nil {
		outgen.logger.Error("MakeSignedObservation returned error", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	if err := so.Verify(outgen.ID(outgen.sharedState.e), outgen.sharedState.seqNr, ev.Query, outgen.offchainKeyring.OffchainPublicKey()); err != nil {
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

	roundCtx := outgen.RoundCtx(msg.SeqNr)

	select {
	case outgen.chOutcomeGenerationToStatePersistence <- EventKVTransactionRequest[RI]{
		roundCtx,
		*outgen.followerState.query,
		msg.AttributedSignedObservations,
		false,
		StateTransitionOutputDigest{},
		ReportsPlusPrecursorDigest{},
		nil,
	}:
	case <-outgen.ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) produceStateTransition(
	roundCtx ocr3_1types.RoundContext,
	kvStoreTxn ocr3_1types.KeyValueStoreTransaction,
	query types.Query,
	asos []AttributedSignedObservation,
	prepared bool,
	commitOutputDigest StateTransitionOutputDigest,
	commitReportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	commitQC []AttributedCommitSignature,
) {
	outgen.followerState.phase = outgenFollowerPhaseBackgroundStateTransition

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		ogid := outgen.ID(roundCtx.Epoch)
		query := query
		kvTxn := kvStoreTxn
		outgen.subs.Go(func() {
			outgen.backgroundStateTransition(
				ctx,
				logger,
				ogid,
				roundCtx,
				query,
				asos,
				kvTxn,
				prepared,
				commitOutputDigest,
				commitReportsPlusPrecursorDigest,
				commitQC,
			)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundStateTransition(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	ogid OutcomeGenerationID,
	roundCtx ocr3_1types.RoundContext,
	query types.Query,
	asos []AttributedSignedObservation,
	kvTxn ocr3_1types.KeyValueStoreTransaction,
	prepared bool,
	commitOutputDigest StateTransitionOutputDigest,
	commitReportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	commitQC []AttributedCommitSignature,
) {
	shouldDiscardKVTxn := true
	defer func() {
		if shouldDiscardKVTxn {
			select {
			case outgen.chOutcomeGenerationToStatePersistence <- EventDiscardKVTransaction[RI]{kvTxn.SeqNr()}:
			case <-outgen.ctx.Done():
				return
			}
		}
	}()

	kvReadWriter, err := kvTxn.GetReadWriter()
	if err != nil {
		outgen.logger.Error("could not get kv transaction read writer", commontypes.LogFields{
			"seqNr": roundCtx.SeqNr,
			"err":   err,
		})
		return
	}

	var attributedObservations []types.AttributedObservation
	// If we have previously prepared this sequence number more >= f correct oracles have
	// checked that the observation quorum is satisfied, the observations are valid.
	// Moreover, they have checked that the attributed observations signatures are valid,
	// and we have not included signatures in asos.
	if !prepared {
		aos, ok := outgen.checkAttributedSignedObservations(ctx, logger, ogid, roundCtx, query, asos, kvReadWriter)
		if !ok {
			return
		}
		attributedObservations = aos
	} else {
		attributedObservations = attributedObservationsFromAttributedSignedObservations(asos)
	}

	reportsPlusPrecursor, ok := callPluginFromOutcomeGenerationBackground[ocr3_1types.ReportsPlusPrecursor](
		ctx,
		logger,
		"StateTransition",
		0, // StateTransition is a pure function and should finish "instantly"
		roundCtx,
		func(ctx context.Context, roundCtx ocr3_1types.RoundContext) (ocr3_1types.ReportsPlusPrecursor, error) {
			return outgen.reportingPlugin.StateTransition(ctx, roundCtx, query, attributedObservations, kvReadWriter)
		},
	)
	if !ok {
		return
	}

	stateTransitionInputs := StateTransitionInputs{
		roundCtx.SeqNr,
		roundCtx.Epoch,
		roundCtx.Round,
		query,
		attributedObservations,
	}
	inputsDigest := MakeStateTransitionInputsDigest(
		ogid,
		roundCtx.SeqNr,
		query,
		attributedObservations,
	)

	outputDigest := MakeStateTransitionOutputDigest(ogid, roundCtx.SeqNr, kvTxn.GetWriteSet())
	reportsPlusPrecursorDigest, err := MakeReportsPlusPrecursorDigest(ogid, roundCtx.SeqNr, reportsPlusPrecursor)
	if err != nil {
		outgen.logger.Error("failed to compute reports digest", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	// If we are dealing with a replayed state transition we need to verify that the return values
	// of the StateTransition are consistent.
	if commitQC != nil {
		if commitOutputDigest != outputDigest {
			outgen.logger.Error("StateTransitionOutput digests do not match. "+
				"This is commonly caused by non-determinism in the ReportingPlugin.", commontypes.LogFields{
				"replayedSeqNr": roundCtx.SeqNr,
				"seqNr":         outgen.sharedState.seqNr,
			})
			return
		}
		if commitReportsPlusPrecursorDigest != reportsPlusPrecursorDigest {
			outgen.logger.Error("ReportsPlusPrecursor digests do not match. "+
				"This is commonly caused by non-determinism in the ReportingPlugin.", commontypes.LogFields{
				"replayedSeqNr": roundCtx.SeqNr,
				"seqNr":         outgen.sharedState.seqNr,
			})
			return
		}
	}

	outgen.followerState.roundInfo = roundInfo[RI]{
		stateTransitionInputs,
		reportsPlusPrecursor,
		inputsDigest,
		outputDigest,
		reportsPlusPrecursorDigest,
		commitQC,
	}
	select {
	case outgen.chOutcomeGenerationToStatePersistence <- EventComputedStateTransition[RI]{
		kvTxn.SeqNr(),
		outgen.sharedState.e,
	}:
		shouldDiscardKVTxn = false
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) processAcknowledgedComputedStateTransition(ev EventAcknowledgedComputedStateTransition[RI]) {
	if ev.SeqNr != outgen.followerState.roundInfo.inputs.SeqNr {
		outgen.logger.Error("event AcknowledgedComputedStateTransition seqNr does not match round info", commontypes.LogFields{
			"seqNr":      outgen.followerState.roundInfo.inputs.SeqNr,
			"eventSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.followerState.phase != outgenFollowerPhaseBackgroundStateTransition {
		outgen.logger.Debug("discarding event AcknowledgedComputedStateTransition, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	// We are in state-sync
	if outgen.followerState.roundInfo.commitQuorumCertificate != nil {
		outgen.commit(CertifiedCommit{
			outgen.followerState.roundInfo.inputs,
			outgen.followerState.roundInfo.outputDigest,
			outgen.followerState.roundInfo.reportsPlusPrecursorDigest,
			outgen.followerState.roundInfo.commitQuorumCertificate,
		})
		return
	}

	prepareSignature, err := MakePrepareSignature(
		outgen.ID(outgen.sharedState.e),
		outgen.followerState.roundInfo.inputs.SeqNr,
		outgen.followerState.roundInfo.inputsDigest,
		outgen.followerState.roundInfo.outputDigest,
		outgen.followerState.roundInfo.reportsPlusPrecursorDigest,
		outgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		outgen.logger.Critical("failed to sign Prepare", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseSentPrepare

	outgen.logger.Debug("broadcasting MessagePrepare", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})
	outgen.netSender.Broadcast(MessagePrepare[RI]{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
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

	for sender, preparePoolEntry := range poolEntries {
		if preparePoolEntry.Verified != nil {
			continue
		}
		err := preparePoolEntry.Item.Verify(
			outgen.ID(outgen.sharedState.e),
			outgen.sharedState.seqNr,
			outgen.followerState.roundInfo.inputsDigest,
			outgen.followerState.roundInfo.outputDigest,
			outgen.followerState.roundInfo.reportsPlusPrecursorDigest,
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
		outgen.ID(outgen.sharedState.e),
		outgen.sharedState.seqNr,
		outgen.followerState.roundInfo.inputsDigest,
		outgen.followerState.roundInfo.outputDigest,
		outgen.followerState.roundInfo.reportsPlusPrecursorDigest,
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
		outgen.followerState.roundInfo.inputs,
		outgen.followerState.roundInfo.outputDigest,
		outgen.followerState.roundInfo.reportsPlusPrecursorDigest,
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

	for sender, commitPoolEntry := range poolEntries {
		if commitPoolEntry.Verified != nil {
			continue
		}
		err := commitPoolEntry.Item.Verify(
			outgen.ID(outgen.sharedState.e),
			outgen.sharedState.seqNr,
			outgen.followerState.roundInfo.inputsDigest,
			outgen.followerState.roundInfo.outputDigest,
			outgen.followerState.roundInfo.reportsPlusPrecursorDigest,
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

	outgen.commit(CertifiedCommit{
		outgen.followerState.roundInfo.inputs,
		outgen.followerState.roundInfo.outputDigest,
		outgen.followerState.roundInfo.reportsPlusPrecursorDigest,
		commitQuorumCertificate,
	})
	if outgen.id == outgen.sharedState.l {
		outgen.metrics.ledCommittedRoundsTotal.Inc()
	}
}

func (outgen *outcomeGenerationState[RI]) commit(commit CertifiedCommit) {
	if commit.SeqNr() < outgen.sharedState.committedToKVStoreSeqNr {
		outgen.logger.Critical("assumption violation, commitSeqNr is less than committedToKVStoreSeqNr", commontypes.LogFields{
			"commitSeqNr":             commit.SeqNr(),
			"committedToKVStoreSeqNr": outgen.sharedState.committedToKVStoreSeqNr,
		})
		return
	}

	if commit.SeqNr() == outgen.sharedState.committedSeqNr {

		outgen.logger.Debug("skipping commit of already committed seqNr", commontypes.LogFields{
			"commitSeqNr ":   commit.SeqNr(),
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
	} else { // commit.SeqNr == outgen.sharedState.committedSeqNr + 1

		if commit.SeqNr() != outgen.sharedState.committedSeqNr+1 {
			outgen.logger.Critical("assumption violation, committing out of order", commontypes.LogFields{
				"commitSeqNr":    commit.SeqNr(),
				"committedSeqNr": outgen.sharedState.committedSeqNr,
			})
			return
		}

		if !outgen.persistAndUpdateCertIfGreater(&commit) {
			return
		}

		outgen.sharedState.committedSeqNr = commit.SeqNr()
		outgen.metrics.committedSeqNr.Set(float64(commit.SeqNr()))

		outgen.logger.Debug("✅ committed", commontypes.LogFields{
			"seqNr": commit.SeqNr(),
		})

		select {
		case outgen.chOutcomeGenerationToReportAttestation <- EventCertifiedCommit[RI]{
			CertifiedCommittedReports[RI]{
				commit.Epoch(),
				commit.SeqNr(),
				MakeStateTransitionInputsDigest(
					outgen.ID(commit.Epoch()),
					commit.SeqNr(),
					commit.StateTransitionInputs.Query,
					commit.StateTransitionInputs.AttributedObservations,
				),
				commit.StateTransitionOutputDigest,
				outgen.followerState.roundInfo.reportsPlusPrecursor,
				commit.CommitQuorumCertificate,
			},
		}:
		case <-outgen.ctx.Done():
			return
		}
	}

	// We might have already processed the commit before a crash, but not persisted to KV
	if outgen.sharedState.committedToKVStoreSeqNr+1 == outgen.sharedState.committedSeqNr {
		select {
		case outgen.chOutcomeGenerationToStatePersistence <- EventAttestedStateTransitionBlock[RI]{
			AttestedStateTransitionBlock{
				StateTransitionBlock{
					outgen.followerState.roundInfo.inputs,
					outgen.followerState.roundInfo.outputDigest,
					outgen.followerState.roundInfo.reportsPlusPrecursorDigest,
				},
				commit.CommitQuorumCertificate,
			},
		}:
		case <-outgen.ctx.Done():
			return
		}
	}

	outgen.followerState.roundStartPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.proposalPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.preparePool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.commitPool.ReapCompleted(outgen.sharedState.committedSeqNr)
}

func (outgen *outcomeGenerationState[RI]) replayStateTransition(
	inputs StateTransitionInputs,
	outputDigest StateTransitionOutputDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	commitQC []AttributedCommitSignature,
) {
	if outgen.sharedState.committedToKVStoreSeqNr+1 != inputs.SeqNr {
		outgen.logger.Warn("cannot replay out of order state transitions",
			commontypes.LogFields{
				"committedToKVStoreSeqNr": outgen.sharedState.committedToKVStoreSeqNr,
				"replayedSeqNr":           inputs.SeqNr,
			})
		select {
		case outgen.chOutcomeGenerationToStatePersistence <- EventStateSyncRequest[RI]{inputs.SeqNr}:
		case <-outgen.ctx.Done():
		}
		return
	}

	select {
	case outgen.chOutcomeGenerationToStatePersistence <- EventKVTransactionRequest[RI]{
		ocr3_1types.RoundContext{
			inputs.SeqNr,
			inputs.Epoch,
			inputs.Round,
		},
		inputs.Query,
		attributedSignedObservationsFromAttributedObservations(inputs.AttributedObservations),
		true,
		outputDigest,
		reportsPlusPrecursorDigest,
		commitQC,
	}:
	case <-outgen.ctx.Done():
	}
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

// If the attributed signed observations have valid signature, and they satisfy ValidateObservation
// and ObservationQuorum plugin methods, this function returns the vector of corresponding
// AttributedObservations and true.
func (outgen *outcomeGenerationState[RI]) checkAttributedSignedObservations(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	ogid OutcomeGenerationID,
	roundCtx ocr3_1types.RoundContext,
	query types.Query,
	asos []AttributedSignedObservation,
	kvReader ocr3_1types.KeyValueReader,
) ([]types.AttributedObservation, bool) {

	attributedObservations := []types.AttributedObservation{}

	seen := map[commontypes.OracleID]bool{}
	for _, aso := range asos {
		if !(0 <= int(aso.Observer) && int(aso.Observer) <= outgen.config.N()) {
			logger.Warn("dropping MessageProposal that contains signed observation with invalid observer", commontypes.LogFields{
				"seqNr":           roundCtx.SeqNr,
				"invalidObserver": aso.Observer,
			})
			return nil, false
		}

		if seen[aso.Observer] {
			logger.Warn("dropping MessageProposal that contains duplicate signed observation", commontypes.LogFields{
				"seqNr": roundCtx.SeqNr,
			})
			return nil, false
		}

		seen[aso.Observer] = true

		if err := aso.SignedObservation.Verify(ogid, roundCtx.SeqNr, query, outgen.config.OracleIdentities[aso.Observer].OffchainPublicKey); err != nil {
			logger.Warn("dropping MessageProposal that contains signed observation with invalid signature", commontypes.LogFields{
				"seqNr": roundCtx.SeqNr,
				"error": err,
			})
			return nil, false
		}

		err, ok := callPluginFromOutcomeGenerationBackground[error](
			ctx,
			logger,
			"ValidateObservation",
			0, // ValidateObservation is a pure function and should finish "instantly"
			roundCtx,
			func(ctx context.Context, roundCtx ocr3_1types.RoundContext) (error, error) {
				return outgen.reportingPlugin.ValidateObservation(
					ctx,
					roundCtx,
					query,
					types.AttributedObservation{aso.SignedObservation.Observation, aso.Observer},
					kvReader,
				), nil
			},
		)
		// kvReader.Discard() must not happen here, because
		// backgroundStateTransition (our caller) manages the lifecycle of the
		// underlying transaction.
		if !ok {
			logger.Error("dropping MessageProposal containing observation that could not be validated", commontypes.LogFields{
				"seqNr":    roundCtx.SeqNr,
				"observer": aso.Observer,
			})
			return nil, false
		}
		if err != nil {
			logger.Warn("dropping MessageProposal that contains an invalid observation", commontypes.LogFields{
				"seqNr":    roundCtx.SeqNr,
				"error":    err,
				"observer": aso.Observer,
			})
			return nil, false
		}

		attributedObservations = append(attributedObservations, types.AttributedObservation{
			aso.SignedObservation.Observation,
			aso.Observer,
		})
	}

	observationQuorum, ok := callPluginFromOutcomeGenerationBackground[bool](
		ctx,
		logger,
		"ObservationQuorum",
		0, // ObservationQuorum is a pure function and should finish "instantly"
		roundCtx,
		func(ctx context.Context, roundCtx ocr3_1types.RoundContext) (bool, error) {
			return outgen.reportingPlugin.ObservationQuorum(ctx, roundCtx, query, attributedObservations, kvReader)
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

	if seen[outgen.id] {
		outgen.metrics.includedObservationsTotal.Inc()
	}

	return attributedObservations, true
}

func attributedObservationsFromAttributedSignedObservations(asos []AttributedSignedObservation) []types.AttributedObservation {
	aos := make([]types.AttributedObservation, len(asos))
	for i, aso := range asos {
		aos[i] = types.AttributedObservation{
			aso.SignedObservation.Observation,
			aso.Observer,
		}
	}
	return aos
}

func attributedSignedObservationsFromAttributedObservations(aos []types.AttributedObservation) []AttributedSignedObservation {
	asos := make([]AttributedSignedObservation, len(aos))
	for i, ao := range aos {
		asos[i] = AttributedSignedObservation{
			SignedObservation{
				ao.Observation,
				nil,
			},
			ao.Observer,
		}
	}
	return asos
}
