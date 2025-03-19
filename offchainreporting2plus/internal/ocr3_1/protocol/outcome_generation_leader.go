package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/pool"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type outgenLeaderPhase string

const (
	outgenLeaderPhaseUnknown        outgenLeaderPhase = "unknown"
	outgenLeaderPhaseNewEpoch       outgenLeaderPhase = "newEpoch"
	outgenLeaderPhaseSentEpochStart outgenLeaderPhase = "sentEpochStart"
	outgenLeaderPhaseSentRoundStart outgenLeaderPhase = "sentRoundStart"
	outgenLeaderPhaseGrace          outgenLeaderPhase = "grace"
	outgenLeaderPhaseSentProposal   outgenLeaderPhase = "sentProposal"
)

func (outgen *outcomeGenerationState[RI]) messageEpochStartRequest(msg MessageEpochStartRequest[RI], sender commontypes.OracleID) {
	outgen.logger.Debug("received MessageEpochStartRequest", commontypes.LogFields{
		"sender":                       sender,
		"msgHighestCertifiedTimestamp": msg.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp,
		"msgEpoch":                     msg.Epoch,
	})

	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageEpochStartRequest for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if outgen.sharedState.l != outgen.id {
		outgen.logger.Warn("dropping MessageEpochStartRequest to non-leader", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	if outgen.leaderState.phase != outgenLeaderPhaseNewEpoch {
		outgen.logger.Debug("dropping MessageEpochStartRequest for wrong phase", commontypes.LogFields{
			"sender": sender,
			"phase":  outgen.leaderState.phase,
		})
		return
	}

	if outgen.leaderState.epochStartRequests[sender] != nil {
		outgen.logger.Warn("dropping duplicate MessageEpochStartRequest", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	outgen.leaderState.epochStartRequests[sender] = &epochStartRequest[RI]{}

	if err := msg.SignedHighestCertifiedTimestamp.Verify(
		outgen.ID(outgen.sharedState.e),
		outgen.config.OracleIdentities[sender].OffchainPublicKey,
	); err != nil {
		outgen.leaderState.epochStartRequests[sender].bad = true
		outgen.logger.Warn("MessageEpochStartRequest.SignedHighestCertifiedTimestamp is invalid", commontypes.LogFields{
			"sender": sender,
			"error":  err,
		})
		return
	}

	// Note that the MessageEpochStartRequest might still be invalid, e.g. if its HighestCertified is invalid.
	outgen.logger.Debug("got MessageEpochStartRequest with valid SignedHighestCertifiedTimestamp", commontypes.LogFields{
		"sender":                       sender,
		"msgHighestCertifiedTimestamp": msg.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp,
	})

	outgen.leaderState.epochStartRequests[sender].message = msg

	if len(outgen.leaderState.epochStartRequests) < outgen.config.ByzQuorumSize() {
		return
	}

	goodCount := 0
	var maxSender *commontypes.OracleID
	for sender, epochStartRequest := range outgen.leaderState.epochStartRequests {
		if epochStartRequest.bad {
			continue
		}
		goodCount++

		if maxSender == nil || outgen.leaderState.epochStartRequests[*maxSender].message.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp.Less(epochStartRequest.message.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp) {
			sender := sender
			maxSender = &sender
		}
	}

	if maxSender == nil || goodCount < outgen.config.ByzQuorumSize() {
		return
	}

	maxRequest := outgen.leaderState.epochStartRequests[*maxSender]

	if maxRequest.message.HighestCertified.Timestamp() != maxRequest.message.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp {
		maxRequest.bad = true
		outgen.logger.Warn("timestamp mismatch in MessageEpochStartRequest", commontypes.LogFields{
			"sender":                          *maxSender,
			"highestCertified.Timestamp":      maxRequest.message.HighestCertified.Timestamp(),
			"signedHighestCertifiedTimestamp": maxRequest.message.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp,
		})
		return
	}

	if err := maxRequest.message.HighestCertified.Verify(
		outgen.config.ConfigDigest,
		outgen.config.OracleIdentities,
		outgen.config.ByzQuorumSize(),
	); err != nil {
		maxRequest.bad = true
		outgen.logger.Warn("MessageEpochStartRequest.HighestCertified is invalid", commontypes.LogFields{
			"sender": *maxSender,
			"error":  err,
		})
		return
	}

	highestCertifiedProof := make([]AttributedSignedHighestCertifiedTimestamp, 0, outgen.config.ByzQuorumSize())
	contributors := make([]commontypes.OracleID, 0, outgen.config.ByzQuorumSize())
	for sender, epochStartRequest := range outgen.leaderState.epochStartRequests {
		if epochStartRequest.bad {
			continue
		}
		highestCertifiedProof = append(highestCertifiedProof, AttributedSignedHighestCertifiedTimestamp{
			epochStartRequest.message.SignedHighestCertifiedTimestamp,
			sender,
		})
		contributors = append(contributors, sender)
		// not necessary, but hopefully helps with readability
		if len(highestCertifiedProof) == outgen.config.ByzQuorumSize() {
			break
		}
	}

	epochStartProof := EpochStartProof{
		maxRequest.message.HighestCertified,
		highestCertifiedProof,
	}

	// This is a sanity check to ensure that we only construct epochStartProofs that are actually valid.
	// This should never fail.
	if err := epochStartProof.Verify(outgen.ID(outgen.sharedState.e), outgen.config.OracleIdentities, outgen.config.ByzQuorumSize()); err != nil {
		outgen.logger.Critical("EpochStartProof is invalid, very surprising!", commontypes.LogFields{
			"proof": epochStartProof,
		})
		return
	}

	outgen.leaderState.phase = outgenLeaderPhaseSentEpochStart

	outgen.logger.Info("broadcasting MessageEpochStart", commontypes.LogFields{
		"contributors":              contributors,
		"highestCertifiedTimestamp": epochStartProof.HighestCertified.Timestamp(),
		"highestCertifiedQCSeqNr":   epochStartProof.HighestCertified.SeqNr(),
	})

	outgen.netSender.Broadcast(MessageEpochStart[RI]{
		outgen.sharedState.e,
		epochStartProof,
	})

	if epochStartProof.HighestCertified.IsGenesis() {
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedToKVStoreSeqNr + 1
		outgen.startSubsequentLeaderRound()
	} else if commitQC, ok := epochStartProof.HighestCertified.(*CertifiedCommit); ok {
		outgen.sharedState.firstSeqNrOfEpoch = commitQC.SeqNr() + 1
		if commitQC.SeqNr() == outgen.sharedState.committedToKVStoreSeqNr {
			outgen.startSubsequentLeaderRound()
		} else {
			outgen.logger.Debug("leader cannot start subsequent round yet, will try to state sync first", commontypes.LogFields{
				"commitQCSeqNr":           outgen.sharedState.seqNr,
				"committedToKVStoreSeqNr": commitQC.SeqNr(),
			})
			// The leader will try to state sync in the follower code path and
		}
	} else {
		prepareQc, ok := epochStartProof.HighestCertified.(*CertifiedPrepare)
		if !ok {
			outgen.logger.Critical("cast to CertifiedPrepare failed while processing MessageEpochStartRequest", nil)
			return
		}
		outgen.sharedState.firstSeqNrOfEpoch = prepareQc.SeqNr() + 1
		// We're dealing with a re-proposal from a failed epoch based on a
		// prepare qc.
		// We don't want to send MessageRoundStart.
	}
}

func (outgen *outcomeGenerationState[RI]) eventTRoundTimeout() {
	outgen.logger.Debug("TRound fired", commontypes.LogFields{
		"seqNr":          outgen.sharedState.seqNr,
		"committedSeqNr": outgen.sharedState.committedSeqNr,
		"deltaRound":     outgen.config.DeltaRound.String(),
	})
	outgen.startSubsequentLeaderRound()
}

func (outgen *outcomeGenerationState[RI]) startSubsequentLeaderRound() {
	outgen.logger.Debug("trying to start new leader round", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})

	if !outgen.leaderState.readyToStartRound {
		outgen.logger.Debug("not ready to start new leader round yet", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		outgen.leaderState.readyToStartRound = true
		return
	}
	outgen.leaderState.readyToStartRound = false
	outgen.logger.Debug("starting new leader round", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		roundCtx := outgen.RoundCtx(outgen.sharedState.committedToKVStoreSeqNr + 1)
		kvReader := outgen.kvStore.GetReader()
		outgen.subs.Go(func() {
			outgen.backgroundQuery(ctx, logger, roundCtx, kvReader)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundQuery(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	roundCtx ocr3_1types.RoundContext,
	kvReader ocr3_1types.KeyValueReaderDiscardable,
) {
	query, ok := callPluginFromOutcomeGenerationBackground[types.Query](
		ctx,
		logger,
		"Query",
		outgen.config.MaxDurationQuery,
		roundCtx,
		func(ctx context.Context, outctx ocr3_1types.RoundContext) (types.Query, error) {
			return outgen.reportingPlugin.Query(ctx, roundCtx, kvReader)
		},
	)
	kvReader.Discard()
	if !ok {
		return
	}

	select {
	case outgen.chLocalEvent <- EventComputedQuery[RI]{
		roundCtx.Epoch,
		roundCtx.SeqNr,
		query,
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedQuery(ev EventComputedQuery[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.committedSeqNr+1 {
		outgen.logger.Debug("discarding EventComputedQuery from old round", commontypes.LogFields{
			"seqNr":          outgen.sharedState.seqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
			"evEpoch":        ev.Epoch,
			"evSeqNr":        ev.SeqNr,
		})
		return
	}

	outgen.leaderState.query = ev.Query

	outgen.leaderState.observationPool.ReapCompleted(outgen.sharedState.committedSeqNr)

	outgen.leaderState.tRound = time.After(outgen.config.DeltaRound)

	outgen.leaderState.phase = outgenLeaderPhaseSentRoundStart
	outgen.logger.Debug("broadcasting MessageRoundStart", commontypes.LogFields{
		"seqNr": outgen.sharedState.committedSeqNr + 1,
	})
	outgen.netSender.Broadcast(MessageRoundStart[RI]{
		outgen.sharedState.e,
		outgen.sharedState.committedSeqNr + 1,
		ev.Query,
	})
}

func (outgen *outcomeGenerationState[RI]) messageObservation(msg MessageObservation[RI], sender commontypes.OracleID) {

	outgen.logger.Debug("received MessageObservation", commontypes.LogFields{
		"sender":   sender,
		"msgSeqNr": msg.SeqNr,
		"msgEpoch": msg.Epoch,
	})

	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageObservation for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if outgen.sharedState.l != outgen.id {
		outgen.logger.Warn("dropping MessageObservation to non-leader", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if outgen.leaderState.phase != outgenLeaderPhaseSentRoundStart && outgen.leaderState.phase != outgenLeaderPhaseGrace {
		outgen.logger.Debug("dropping MessageObservation for wrong phase", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
			"phase":    outgen.leaderState.phase,
		})
		return
	}

	if msg.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("dropping MessageObservation with invalid SeqNr", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if putResult := outgen.leaderState.observationPool.Put(msg.SeqNr, sender, msg.SignedObservation); putResult != pool.PutResultOK {
		outgen.logger.Warn("dropping MessageObservation", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
			"reason":   putResult,
		})
		return
	}

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		ogid := outgen.ID(outgen.sharedState.e)
		roundCtx := outgen.RoundCtx(outgen.sharedState.seqNr)
		query := outgen.leaderState.query
		kvReader := outgen.kvStore.GetReader()
		outgen.subs.Go(func() {
			outgen.backgroundVerifyValidateObservation(ctx, logger, ogid, roundCtx, sender, msg.SignedObservation, query, kvReader)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundVerifyValidateObservation(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	ogid OutcomeGenerationID,
	roundCtx ocr3_1types.RoundContext,
	sender commontypes.OracleID,
	signedObservation SignedObservation,
	query types.Query,
	kvReader ocr3_1types.KeyValueReaderDiscardable,
) {

	if err := signedObservation.Verify(
		ogid,
		roundCtx.SeqNr,
		query,
		outgen.config.OracleIdentities[sender].OffchainPublicKey,
	); err != nil {
		logger.Warn("dropping MessageObservation carrying invalid SignedObservation", commontypes.LogFields{
			"sender": sender,
			"seqNr":  roundCtx.SeqNr,
			"error":  err,
		})
		return
	}

	err, ok := callPluginFromOutcomeGenerationBackground[error](
		ctx,
		logger,
		"ValidateObservation",
		0, // ValidateObservation is a pure function and should finish "instantly"
		roundCtx,
		func(ctx context.Context, roundCtx ocr3_1types.RoundContext) (error, error) {
			return outgen.reportingPlugin.ValidateObservation(ctx,
				roundCtx,
				query,
				types.AttributedObservation{signedObservation.Observation, sender},
				kvReader,
			), nil
		},
	)
	kvReader.Discard()
	if !ok {
		logger.Error("dropping MessageObservation carrying Observation that could not be validated", commontypes.LogFields{
			"sender": sender,
			"seqNr":  roundCtx.SeqNr,
		})
		return
	}

	if err != nil {
		logger.Warn("dropping MessageObservation carrying invalid Observation", commontypes.LogFields{
			"sender": sender,
			"seqNr":  roundCtx.SeqNr,
			"error":  err,
		})
		return
	}

	select {
	case outgen.chLocalEvent <- EventComputedValidateVerifyObservation[RI]{
		roundCtx.Epoch,
		roundCtx.SeqNr,
		sender,
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedValidateVerifyObservation(ev EventComputedValidateVerifyObservation[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("discarding EventComputedValidateVerifyObservation from old round", commontypes.LogFields{
			"seqNr":   outgen.sharedState.seqNr,
			"evEpoch": ev.Epoch,
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.leaderState.phase != outgenLeaderPhaseSentRoundStart && outgen.leaderState.phase != outgenLeaderPhaseGrace {
		outgen.logger.Debug("discarding EventComputedValidateVerifyObservation, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.leaderState.phase,
		})
		return
	}

	outgen.logger.Debug("got valid MessageObservation", commontypes.LogFields{
		"sender": ev.Sender,
		"seqNr":  ev.SeqNr,
	})

	outgen.leaderState.observationPool.StoreVerified(outgen.sharedState.seqNr, ev.Sender, true)

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		outctx := outgen.RoundCtx(outgen.sharedState.seqNr)
		query := outgen.leaderState.query
		aos := []types.AttributedObservation{}
		for sender, observationPoolEntry := range outgen.leaderState.observationPool.Entries(outgen.sharedState.seqNr) {
			if observationPoolEntry.Verified == nil || !*observationPoolEntry.Verified {
				continue
			}
			aos = append(aos, types.AttributedObservation{observationPoolEntry.Item.Observation, sender})
		}
		kvReader := outgen.kvStore.GetReader()

		outgen.subs.Go(func() {
			outgen.backgroundObservationQuorum(
				ctx,
				logger,
				outctx,
				query,
				aos,
				kvReader,
			)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundObservationQuorum(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	roundCtx ocr3_1types.RoundContext,
	query types.Query,
	aos []types.AttributedObservation,
	kvReader ocr3_1types.KeyValueReaderDiscardable,
) {
	observationQuorum, ok := callPluginFromOutcomeGenerationBackground[bool](
		ctx,
		logger,
		"ObservationQuorum",
		0, // ObservationQuorum is a pure function and should finish "instantly"
		roundCtx,
		func(ctx context.Context, roundCtx ocr3_1types.RoundContext) (bool, error) {
			return outgen.reportingPlugin.ObservationQuorum(ctx, roundCtx, query, aos, kvReader)
		},
	)
	kvReader.Discard()

	if !ok {
		return
	}

	if !observationQuorum {
		if len(aos) >= outgen.config.N()-outgen.config.F {
			logger.Warn("ObservationQuorum returned false despite there being at least n-f valid observations. This is the maximum number of valid observations we are guaranteed to receive. Maybe there is a bug in the ReportingPlugin.", commontypes.LogFields{
				"attributedObservationCount": len(aos),
				"nMinusF":                    outgen.config.N() - outgen.config.F,
				"seqNr":                      roundCtx.SeqNr,
			})
		}
		return
	}

	select {
	case outgen.chLocalEvent <- EventComputedObservationQuorumSuccess[RI]{
		roundCtx.Epoch,
		roundCtx.SeqNr,
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedObservationQuorumSuccess(ev EventComputedObservationQuorumSuccess[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("discarding EventComputedObservationQuorumSuccess from old round", commontypes.LogFields{
			"seqNr":   outgen.sharedState.seqNr,
			"evEpoch": ev.Epoch,
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.leaderState.phase != outgenLeaderPhaseSentRoundStart {
		outgen.logger.Debug("discarding EventComputedObservationQuorumSuccess, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.leaderState.phase,
		})
		return
	}

	outgen.logger.Debug("reached observation quorum, starting observation grace period", commontypes.LogFields{
		"seqNr":      outgen.sharedState.seqNr,
		"deltaGrace": outgen.config.DeltaGrace.String(),
	})
	outgen.leaderState.phase = outgenLeaderPhaseGrace
	outgen.leaderState.tGrace = time.After(outgen.config.DeltaGrace)
}

func (outgen *outcomeGenerationState[RI]) eventTGraceTimeout() {
	if outgen.leaderState.phase != outgenLeaderPhaseGrace {
		outgen.logger.Error("leader's phase conflicts TGrace timeout", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.leaderState.phase,
		})
		return
	}
	asos := make([]AttributedSignedObservation, 0, outgen.config.N())
	contributors := make([]commontypes.OracleID, 0, outgen.config.N())
	for sender, observationPoolEntry := range outgen.leaderState.observationPool.Entries(outgen.sharedState.seqNr) {
		if observationPoolEntry.Verified == nil || !*observationPoolEntry.Verified {
			continue
		}
		asos = append(asos, AttributedSignedObservation{SignedObservation: observationPoolEntry.Item, Observer: sender})
		contributors = append(contributors, sender)
	}

	outgen.leaderState.phase = outgenLeaderPhaseSentProposal

	outgen.logger.Debug("broadcasting MessageProposal after TGrace fired", commontypes.LogFields{
		"seqNr":        outgen.sharedState.seqNr,
		"contributors": contributors,
		"deltaGrace":   outgen.config.DeltaGrace.String(),
	})
	outgen.netSender.Broadcast(MessageProposal[RI]{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		asos,
	})

	outgen.leaderState.observationPool.ReapCompleted(outgen.sharedState.seqNr)
}
