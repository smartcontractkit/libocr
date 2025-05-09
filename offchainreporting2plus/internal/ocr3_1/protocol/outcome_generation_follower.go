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
	outgenFollowerPhaseUnknown                           outgenFollowerPhase = "unknown"
	outgenFollowerPhaseNewEpoch                          outgenFollowerPhase = "newEpoch"
	outgenFollowerPhaseNewRound                          outgenFollowerPhase = "newRound"
	outgenFollowerPhaseBackgroundObservation             outgenFollowerPhase = "backgroundObservation"
	outgenFollowerPhaseSentObservation                   outgenFollowerPhase = "sentObservation"
	outgenFollowerPhaseBackgroundProposalStateTransition outgenFollowerPhase = "backgroundProposalStateTransition"
	outgenFollowerPhaseSentPrepare                       outgenFollowerPhase = "sentPrepare"
	outgenFollowerPhaseSentCommit                        outgenFollowerPhase = "sentCommit"
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

	outgen.refreshCommittedSeqNrAndCert()
	if !outgen.ensureHighestCertifiedIsCompatible(msg.EpochStartProof.HighestCertified, "MessageEpochStart") {
		return
	}

	if msg.EpochStartProof.HighestCertified.IsGenesis() {
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentFollowerRound()
	} else if commitQC, ok := msg.EpochStartProof.HighestCertified.(*CertifiedCommit); ok {
		outgen.commit(*commitQC)

		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
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

		stateTransitionInputs := prepareQc.StateTransitionInputs
		prepareQcSeqNr := stateTransitionInputs.SeqNr

		outgen.sharedState.firstSeqNrOfEpoch = prepareQcSeqNr + 1
		outgen.sharedState.seqNr = prepareQcSeqNr
		outgen.sharedState.observationQuorum = nil

		outgen.followerState.query = nil
		outgen.ensureOpenKVTransactionDiscarded()

		outgen.followerState.phase = outgenFollowerPhaseBackgroundProposalStateTransition
		outgen.followerState.stateTransitionInfo = stateTransitionInfo{}
		outgen.ensureOpenKVTransactionDiscarded()

		outgen.logger.Debug("re-executing StateTransition from MessagePrepare (reproposal)", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})

		kvReadWriteTxn, err := outgen.kvStore.NewReadWriteTransaction(prepareQcSeqNr)
		if err != nil {
			outgen.logger.Warn("could not create kv read/write transaction", commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
				"err":   err,
			})
			return
		}

		{
			ctx := outgen.ctx
			logger := outgen.logger
			ogid := outgen.ID()
			roundCtx := outgen.RoundCtx(prepareQcSeqNr)
			query := stateTransitionInputs.Query
			aos := stateTransitionInputs.AttributedObservations
			outgen.subs.Go(func() {
				outgen.backgroundProposalStateTransition(
					ctx,
					logger,
					ogid,
					roundCtx,
					query,
					nil,
					aos,
					kvReadWriteTxn,
				)
			})
		}
	}
}

func (outgen *outcomeGenerationState[RI]) startSubsequentFollowerRound() {
	outgen.sharedState.seqNr = outgen.sharedState.committedSeqNr + 1
	outgen.sharedState.observationQuorum = nil

	outgen.followerState.phase = outgenFollowerPhaseNewRound
	outgen.followerState.query = nil
	outgen.followerState.stateTransitionInfo = stateTransitionInfo{}
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
		query := *outgen.followerState.query
		kvReadTxn, err := outgen.kvStore.NewReadTransaction(roundCtx.SeqNr)
		if err != nil {
			outgen.logger.Warn("failed to create new transaction, aborting tryProcessRoundStartPool", commontypes.LogFields{
				"seqNr": roundCtx.SeqNr,
				"error": err,
			})
			return
		}
		outgen.subs.Go(func() {
			outgen.backgroundObservation(ctx, logger, roundCtx, query, kvReadTxn)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundObservation(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	roundCtx ocr3_1types.RoundContext,
	query types.Query,
	kvReadTxn KeyValueStoreReadTransaction,
) {
	observation, ok := callPluginFromOutcomeGenerationBackground[types.Observation](
		ctx,
		logger,
		"Observation",
		outgen.config.MaxDurationObservation,
		roundCtx,
		func(ctx context.Context, rondCtx ocr3_1types.RoundContext) (types.Observation, error) {
			return outgen.reportingPlugin.Observation(ctx, roundCtx, query, kvReadTxn)
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

	so, err := MakeSignedObservation(outgen.ID(), outgen.sharedState.seqNr, ev.Query, ev.Observation, outgen.offchainKeyring.OffchainSign)
	if err != nil {
		outgen.logger.Error("MakeSignedObservation returned error", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	if err := so.Verify(outgen.ID(), outgen.sharedState.seqNr, ev.Query, outgen.offchainKeyring.OffchainPublicKey()); err != nil {
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

	kvReadWriteTxn, err := outgen.kvStore.NewReadWriteTransaction(outgen.sharedState.seqNr)
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
		query := *outgen.followerState.query

		asos := msg.AttributedSignedObservations

		outgen.subs.Go(func() {
			outgen.backgroundProposalStateTransition(
				ctx,
				logger,
				ogid,
				roundCtx,
				query,
				asos,
				nil,
				kvReadWriteTxn,
			)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundProposalStateTransition(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	ogid OutcomeGenerationID,
	roundCtx ocr3_1types.RoundContext,
	query types.Query,
	asos []AttributedSignedObservation,
	aos []types.AttributedObservation,
	kvReadWriteTxn KeyValueStoreReadWriteTransaction,
) {
	shouldDiscardKVTxn := true
	defer func() {
		if shouldDiscardKVTxn {
			kvReadWriteTxn.Discard()
		}
	}()

	if aos == nil {
		var ok bool
		aos, ok = outgen.checkAttributedSignedObservations(ctx, logger, ogid, roundCtx, query, asos, kvReadWriteTxn)
		if !ok {
			return
		}
	} else {
		// We do not enter this branch for re-proposals. aos is coming from the prepareQc.
		// If we have previously prepared this sequence number more >= f correct oracles have
		// checked that the observation quorum is satisfied, the observations are valid.
		// Moreover, they have checked that the attributed observations signatures are valid,
		// and we have not included signatures in asos.

	}

	reportsPlusPrecursor, ok := callPluginFromOutcomeGenerationBackground[ocr3_1types.ReportsPlusPrecursor](
		ctx,
		logger,
		"StateTransition",
		0, // StateTransition is a pure function and should finish "instantly"
		roundCtx,
		func(ctx context.Context, roundCtx ocr3_1types.RoundContext) (ocr3_1types.ReportsPlusPrecursor, error) {
			return outgen.reportingPlugin.StateTransition(ctx, roundCtx, query, aos, kvReadWriteTxn)
		},
	)
	if !ok {
		return
	}

	stateTransitionInputsDigest := MakeStateTransitionInputsDigest(
		ogid,
		roundCtx.SeqNr,
		query,
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

	stateTransitionOutputDigest := MakeStateTransitionOutputDigest(ogid, roundCtx.SeqNr, writeSet)
	reportsPlusPrecursorDigest := MakeReportsPlusPrecursorDigest(ogid, roundCtx.SeqNr, reportsPlusPrecursor)

	stateTransitionInputs := StateTransitionInputs{
		roundCtx.SeqNr,
		roundCtx.Epoch,
		roundCtx.Round,
		query,
		aos,
	}

	select {
	case outgen.chLocalEvent <- EventComputedProposalStateTransition[RI]{
		roundCtx.Epoch,
		roundCtx.SeqNr,
		kvReadWriteTxn,
		stateTransitionInfo{
			stateTransitionInputs,
			reportsPlusPrecursor,
			stateTransitionInputsDigest,
			stateTransitionOutputDigest,
			reportsPlusPrecursorDigest,
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

	outgen.followerState.openKVTxn = ev.KeyValueStoreReadWriteTransaction

	prepareSignature, err := MakePrepareSignature(
		outgen.ID(),
		outgen.sharedState.seqNr,
		ev.stateTransitionInfo.InputsDigest,
		ev.stateTransitionInfo.OutputDigest,
		ev.stateTransitionInfo.ReportsPlusPrecursorDigest,
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
	outgen.followerState.stateTransitionInfo = ev.stateTransitionInfo

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
			outgen.ID(),
			outgen.sharedState.seqNr,
			outgen.followerState.stateTransitionInfo.InputsDigest,
			outgen.followerState.stateTransitionInfo.OutputDigest,
			outgen.followerState.stateTransitionInfo.ReportsPlusPrecursorDigest,
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
		outgen.followerState.stateTransitionInfo.InputsDigest,
		outgen.followerState.stateTransitionInfo.OutputDigest,
		outgen.followerState.stateTransitionInfo.ReportsPlusPrecursorDigest,
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
		outgen.followerState.stateTransitionInfo.Inputs,
		outgen.followerState.stateTransitionInfo.OutputDigest,
		outgen.followerState.stateTransitionInfo.ReportsPlusPrecursorDigest,
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
			outgen.ID(),
			outgen.sharedState.seqNr,
			outgen.followerState.stateTransitionInfo.InputsDigest,
			outgen.followerState.stateTransitionInfo.OutputDigest,
			outgen.followerState.stateTransitionInfo.ReportsPlusPrecursorDigest,
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

	persistedBlockAndCert := outgen.commit(CertifiedCommit{
		outgen.followerState.stateTransitionInfo.Inputs,
		outgen.followerState.stateTransitionInfo.OutputDigest,
		outgen.followerState.stateTransitionInfo.ReportsPlusPrecursorDigest,
		commitQuorumCertificate,
	})

	if !persistedBlockAndCert {
		outgen.logger.Error("failed to persist block/cert, stopping to not advance kv further than persisted blocks", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})

		return
	}

	if outgen.followerState.openKVTxn == nil {
		outgen.logger.Critical("assumption violation, open kv transaction must exist in this phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		panic("")
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
			kvSeqNr, err := outgen.kvStore.HighestCommittedSeqNr()
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

	if uint64(outgen.config.RMax) <= outgen.sharedState.seqNr-outgen.sharedState.firstSeqNrOfEpoch+1 {
		outgen.logger.Debug("epoch has been going on for too long, sending EventChangeLeader to Pacemaker", commontypes.LogFields{
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

func (outgen *outcomeGenerationState[RI]) commit(commit CertifiedCommit) (persistedBlockAndCert bool) {
	if commit.SeqNr() < outgen.sharedState.committedSeqNr {
		outgen.logger.Critical("assumption violation, commitSeqNr is less than committedSeqNr", commontypes.LogFields{
			"commitSeqNr":    commit.SeqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
		panic("")
	}

	if commit.SeqNr() <= outgen.sharedState.committedSeqNr {

		outgen.logger.Debug("skipping commit of already committed seqNr", commontypes.LogFields{
			"commitSeqNr ":   commit.SeqNr(),
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
	} else { // commit.SeqNr >= outgen.sharedState.committedSeqNr + 1

		if !outgen.persistCommitAsBlock(&commit) {
			return
		}

		if !outgen.persistAndUpdateCertIfGreater(&commit) {
			return
		}

		persistedBlockAndCert = true

		outgen.sharedState.committedSeqNr = commit.SeqNr()
		outgen.metrics.committedSeqNr.Set(float64(commit.SeqNr()))

		outgen.logger.Debug("✅ committed", commontypes.LogFields{
			"seqNr": commit.SeqNr(),
		})

		if outgen.followerState.phase != outgenFollowerPhaseSentCommit {
			outgen.logger.Debug("skipping notification of report attestation, we don't have the reports plus precursor", commontypes.LogFields{
				"committedSeqNr": outgen.sharedState.committedSeqNr,
				"phase":          outgen.followerState.phase,
			})
			goto ReapCompleted
		}

		reportsPlusPrecursor := outgen.followerState.stateTransitionInfo.ReportsPlusPrecursor

		select {
		case outgen.chOutcomeGenerationToReportAttestation <- EventCertifiedCommit[RI]{
			CertifiedCommittedReports[RI]{
				commit.Epoch(),
				commit.SeqNr(),
				MakeStateTransitionInputsDigest(
					outgen.ID(),
					commit.SeqNr(),
					commit.StateTransitionInputs.Query,
					commit.StateTransitionInputs.AttributedObservations,
				),
				commit.StateTransitionOutputDigest,
				reportsPlusPrecursor,
				commit.CommitQuorumCertificate,
			},
		}:
		case <-outgen.ctx.Done():
			return
		}
	}

ReapCompleted:
	outgen.followerState.roundStartPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.proposalPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.preparePool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.commitPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	return
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

func (outgen *outcomeGenerationState[RI]) persistCommitAsBlock(commit *CertifiedCommit) bool {
	ctx := outgen.ctx
	configDigest := outgen.config.ConfigDigest
	seqNr := commit.SeqNr()
	astb := AttestedStateTransitionBlock{
		StateTransitionBlock{
			commit.StateTransitionInputs,
			commit.StateTransitionOutputDigest,
			commit.ReportsPlusPrecursorDigest,
		},
		commit.CommitQuorumCertificate,
	}

	werr := outgen.database.WriteAttestedStateTransitionBlock(
		ctx,
		configDigest,
		seqNr,
		astb,
	)

	if werr != nil {

		astb, rerr := outgen.database.ReadAttestedStateTransitionBlock(
			ctx,
			configDigest,
			seqNr,
		)
		if astb.StateTransitionBlock.SeqNr() == seqNr && rerr == nil {
			// already persisted by someone else
			return true
		} else {
			outgen.logger.Error("error persisting commit as attested state transition block", commontypes.LogFields{
				"seqNr": seqNr,
				"error": werr,
			})
			return false
		}
	} else {
		// persited now by us
		outgen.logger.Trace("persisted block", commontypes.LogFields{
			"seqNr": seqNr,
		})
		return true
	}
}

func (outgen *outcomeGenerationState[RI]) refreshCommittedSeqNrAndCert() {
	preRefreshCommittedSeqNr := outgen.sharedState.committedSeqNr

	postRefreshCommittedSeqNr, err := outgen.kvStore.HighestCommittedSeqNr()
	if err != nil {
		outgen.logger.Error("kvStore.HighestCommittedSeqNr() failed during refresh", commontypes.LogFields{
			"preRefreshCommittedSeqNr": preRefreshCommittedSeqNr,
			"error":                    err,
		})
		return
	}

	logger := outgen.logger.MakeChild(commontypes.LogFields{
		"preRefreshCommittedSeqNr":  preRefreshCommittedSeqNr,
		"postRefreshCommittedSeqNr": postRefreshCommittedSeqNr,
	})

	if postRefreshCommittedSeqNr == preRefreshCommittedSeqNr {
		return
	} else if postRefreshCommittedSeqNr < preRefreshCommittedSeqNr {
		logger.Critical("assumption violation, kv is behind what outgen knows as committed", nil)
		panic("")
	}

	ctx := outgen.ctx
	configDigest := outgen.config.ConfigDigest
	astb, err := outgen.database.ReadAttestedStateTransitionBlock(
		ctx,
		configDigest,
		postRefreshCommittedSeqNr,
	)
	if err != nil {
		logger.Error("error reading attested state transition block during refresh", commontypes.LogFields{
			"error": err,
		})
		return
	}
	if astb.StateTransitionBlock.SeqNr() == 0 {
		logger.Critical("assumption violation, attested state transition block for kv committed seq nr does not exist", nil)
		panic("")
	}
	if astb.StateTransitionBlock.SeqNr() != postRefreshCommittedSeqNr {
		logger.Critical("assumption violation, attested state transition block has unexpected seq nr", commontypes.LogFields{
			"expectedSeqNr": postRefreshCommittedSeqNr,
			"actualSeqNr":   astb.StateTransitionBlock.SeqNr(),
		})
		panic("")
	}
	stb := astb.StateTransitionBlock

	persistedBlockAndCert := outgen.commit(CertifiedCommit{
		stb.StateTransitionInputs,
		stb.StateTransitionOutputDigest,
		stb.ReportsPrecursorDigest,
		astb.AttributedSignatures,
	})

	if !persistedBlockAndCert {
		logger.Warn("outgen.commit() failed, aborting refresh", nil)
		return
	}

	if outgen.sharedState.committedSeqNr != postRefreshCommittedSeqNr {
		logger.Critical("assumption violation, outgen local committed seq nr did not progress even though we persisted block and cert", commontypes.LogFields{
			"expectedCommittedSeqNr": postRefreshCommittedSeqNr,
			"actualCommittedSeqNr":   outgen.sharedState.committedSeqNr,
		})
		panic("")
	}

	logger.Debug("refreshed cert", nil)
}

func (outgen *outcomeGenerationState[RI]) ensureHighestCertifiedIsCompatible(highestCertified CertifiedPrepareOrCommit, name string) bool {
	logger := outgen.logger
	if highestCertified.IsGenesis() {
		return true
	} else if commitQC, ok := highestCertified.(*CertifiedCommit); ok {
		commitQcSeqNr := commitQC.SeqNr()
		if commitQcSeqNr != outgen.sharedState.committedSeqNr {

			logger.Warn("dropping "+name+" because we are behind (commitQc)", commontypes.LogFields{
				"commitQcSeqNr":         commitQcSeqNr,
				"expectedCommitQcSeqNr": outgen.sharedState.committedSeqNr,
			})

			return false
		}
	} else {
		// We're dealing with a re-proposal from a failed epoch

		prepareQc, ok := highestCertified.(*CertifiedPrepare)
		if !ok {
			logger.Critical("cast to CertifiedPrepare failed while processing "+name, commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
			})
			return false
		}

		expectedPrepareQcSeqNr := outgen.sharedState.committedSeqNr + 1

		stateTransitionInputs := prepareQc.StateTransitionInputs
		prepareQcSeqNr := stateTransitionInputs.SeqNr

		if expectedPrepareQcSeqNr == prepareQcSeqNr {

		} else if expectedPrepareQcSeqNr < prepareQcSeqNr {
			logger.Warn("dropping "+name+" because we are behind (prepareQc)", commontypes.LogFields{
				"prepareQcSeqNr":         prepareQcSeqNr,
				"expectedPrepareQcSeqNr": expectedPrepareQcSeqNr,
			})
			return false
		} else {
			logger.Critical("dropping "+name+" because we are ahead (prepareQc)", commontypes.LogFields{
				"prepareQcSeqNr":         prepareQcSeqNr,
				"expectedPrepareQcSeqNr": expectedPrepareQcSeqNr,
			})

			return false
		}
	}
	return true
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
