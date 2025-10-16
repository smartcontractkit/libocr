package protocol

import (
	"context"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common/pool"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
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
		outgen.ID(),
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

	notBadCount := 0 // Note: just because a request is not marked bad does not mean it's good. Tertium datur!
	for _, epochStartRequest := range outgen.leaderState.epochStartRequests {
		if epochStartRequest.bad {
			continue
		}
		notBadCount++
	}

	if notBadCount < outgen.config.ByzQuorumSize() {
		return
	}

	// The not-bad entries in epochStartRequests here are guaranteed to be
	// nonempty due to definition of ByzQuorumSize.
	var maxSender *commontypes.OracleID
	for sender, epochStartRequest := range outgen.leaderState.epochStartRequests {
		if epochStartRequest.bad {
			continue
		}
		if maxSender != nil {
			maxTimestamp := outgen.leaderState.epochStartRequests[*maxSender].message.HighestCertified.Timestamp()
			if !maxTimestamp.Less(epochStartRequest.message.HighestCertified.Timestamp()) {
				continue
			}
		}
		maxSender = &sender
	}

	if maxSender == nil {
		return
	}

	maxRequest := outgen.leaderState.epochStartRequests[*maxSender]

	if !maxRequest.message.HighestCertified.Timestamp().Equal(maxRequest.message.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp) {
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
	if err := epochStartProof.Verify(outgen.ID(), outgen.config.OracleIdentities, outgen.config.ByzQuorumSize()); err != nil {
		outgen.logger.Critical("EpochStartProof is invalid, very surprising!", commontypes.LogFields{
			"proof": epochStartProof,
		})
		return
	}

	outgen.leaderState.phase = outgenLeaderPhaseSentEpochStart

	outgen.logger.Info("broadcasting MessageEpochStart", commontypes.LogFields{
		"contributors": contributors,
	})

	epochStartSignature31, err := MakeEpochStartSignature31(
		outgen.ID(),
		epochStartProof,
		outgen.offchainKeyring.OffchainSign,
	)

	if err != nil {
		outgen.logger.Error("MakeEpochStartSignature31 returned error", commontypes.LogFields{
			"error": err,
		})
		return
	}

	outgen.netSender.Broadcast(MessageEpochStart[RI]{
		outgen.sharedState.e,
		epochStartProof,
		epochStartSignature31,
	})

	if epochStartProof.HighestCertified.IsGenesis() {
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentLeaderRound()
	} else if commitQC, ok := epochStartProof.HighestCertified.(*CertifiedCommit); ok {
		outgen.commit(*commitQC)
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentLeaderRound()
	} else {
		prepareQC, ok := epochStartProof.HighestCertified.(*CertifiedPrepare)
		if !ok {
			outgen.logger.Critical("cast to CertifiedPrepare failed while processing MessageEpochStartRequest", nil)
			return
		}
		outgen.sharedState.firstSeqNrOfEpoch = prepareQC.SeqNr + 1
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
	if !outgen.leaderState.readyToStartRound {
		outgen.leaderState.readyToStartRound = true
		return
	}
	outgen.leaderState.readyToStartRound = false

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		outctx := outgen.OutcomeCtx(outgen.sharedState.committedSeqNr + 1)
		outgen.subs.Go(func() {
			outgen.backgroundQuery(ctx, logger, outctx)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundQuery(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	outctx ocr3types.OutcomeContext,
) {
	query, ok := callPluginFromOutcomeGenerationBackground[types.Query, RI](
		ctx,
		logger,
		"Query",
		outgen.config.MaxDurationQuery,
		outctx,
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
			return outgen.reportingPlugin.Query(ctx, outctx)
		},
	)
	if !ok {
		return
	}

	select {
	case outgen.chLocalEvent <- EventComputedQuery[RI]{
		outctx.Epoch,
		outctx.SeqNr,
		query,
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedQuery(ev EventComputedQuery[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.committedSeqNr+1 {
		outgen.logger.Debug("dropping EventComputedQuery from old round", commontypes.LogFields{
			"seqNr":          outgen.sharedState.seqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
			"evEpoch":        ev.Epoch,
			"evSeqNr":        ev.SeqNr,
		})
		return
	}

	roundStartSignature31, err := MakeRoundStartSignature31(
		outgen.ID(),
		outgen.sharedState.committedSeqNr+1,
		ev.Query,
		outgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		outgen.logger.Error("MakeRoundStartSignature31 returned error", commontypes.LogFields{
			"error":          err,
			"seqNr":          outgen.sharedState.seqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
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
		roundStartSignature31,
	})
}

func (outgen *outcomeGenerationState[RI]) messageObservation(msg MessageObservation[RI], sender commontypes.OracleID) {

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
		ogid := outgen.ID()
		outctx := outgen.OutcomeCtx(outgen.sharedState.seqNr)
		query := outgen.leaderState.query
		outgen.subs.Go(func() {
			outgen.backgroundVerifyValidateObservation(ctx, logger, ogid, outctx, sender, msg.SignedObservation, query)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundVerifyValidateObservation(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	ogid OutcomeGenerationID,
	outctx ocr3types.OutcomeContext,
	sender commontypes.OracleID,
	signedObservation SignedObservation,
	query types.Query,
) {

	if err := signedObservation.Verify(
		ogid,
		outctx.SeqNr,
		query,
		outgen.config.OracleIdentities[sender].OffchainPublicKey,
	); err != nil {
		logger.Warn("dropping MessageObservation carrying invalid SignedObservation", commontypes.LogFields{
			"sender": sender,
			"seqNr":  outctx.SeqNr,
			"error":  err,
		})
		return
	}

	err, ok := callPluginFromOutcomeGenerationBackground[error, RI](
		ctx,
		logger,
		"ValidateObservation",
		0, // ValidateObservation is a pure function and should finish "instantly"
		outctx,
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (error, error) {
			return outgen.reportingPlugin.ValidateObservation(ctx, outctx, query, types.AttributedObservation{signedObservation.Observation, sender}), nil
		},
	)
	if !ok {
		logger.Error("dropping MessageObservation carrying Observation that could not be validated", commontypes.LogFields{
			"sender": sender,
			"seqNr":  outctx.SeqNr,
		})
		return
	}

	if err != nil {
		logger.Warn("dropping MessageObservation carrying invalid Observation", commontypes.LogFields{
			"sender": sender,
			"seqNr":  outctx.SeqNr,
			"error":  err,
		})
		return
	}

	select {
	case outgen.chLocalEvent <- EventComputedValidateVerifyObservation[RI]{
		outctx.Epoch,
		outctx.SeqNr,
		sender,
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedValidateVerifyObservation(ev EventComputedValidateVerifyObservation[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("dropping EventComputedValidateVerifyObservation from old round", commontypes.LogFields{
			"seqNr":   outgen.sharedState.seqNr,
			"evEpoch": ev.Epoch,
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.leaderState.phase != outgenLeaderPhaseSentRoundStart && outgen.leaderState.phase != outgenLeaderPhaseGrace {
		outgen.logger.Debug("dropping EventComputedValidateVerifyObservation, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.leaderState.phase,
		})
		return
	}

	outgen.logger.Debug("got valid MessageObservation", commontypes.LogFields{
		"sender": ev.Sender,
		"seqNr":  outgen.sharedState.seqNr,
	})

	outgen.leaderState.observationPool.StoreVerified(outgen.sharedState.seqNr, ev.Sender, true)

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		outctx := outgen.OutcomeCtx(outgen.sharedState.seqNr)
		query := outgen.leaderState.query
		aos := []types.AttributedObservation{}
		for sender, observationPoolEntry := range outgen.leaderState.observationPool.Entries(outgen.sharedState.seqNr) {
			if observationPoolEntry.Verified == nil || !*observationPoolEntry.Verified {
				continue
			}
			aos = append(aos, types.AttributedObservation{observationPoolEntry.Item.Observation, sender})
		}

		outgen.subs.Go(func() {
			outgen.backgroundObservationQuorum(
				ctx,
				logger,
				outctx,
				query,
				aos,
			)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundObservationQuorum(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	outctx ocr3types.OutcomeContext,
	query types.Query,
	aos []types.AttributedObservation,
) {
	observationQuorum, ok := callPluginFromOutcomeGenerationBackground[bool, RI](
		ctx,
		logger,
		"ObservationQuorum",
		0, // ObservationQuorum is a pure function and should finish "instantly"
		outctx,
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (bool, error) {
			return outgen.reportingPlugin.ObservationQuorum(ctx, outctx, query, aos)
		},
	)

	if !ok {
		return
	}

	if !observationQuorum {
		if len(aos) >= outgen.config.N()-outgen.config.F {
			logger.Warn("ObservationQuorum returned false despite there being at least n-f valid observations. This is the maximum number of valid observations we are guaranteed to receive. Maybe there is a bug in the ReportingPlugin.", commontypes.LogFields{
				"attributedObservationCount": len(aos),
				"nMinusF":                    outgen.config.N() - outgen.config.F,
				"seqNr":                      outctx.SeqNr,
			})
		}
		return
	}

	select {
	case outgen.chLocalEvent <- EventComputedObservationQuorumSuccess[RI]{
		outctx.Epoch,
		outctx.SeqNr,
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedObservationQuorumSuccess(ev EventComputedObservationQuorumSuccess[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("dropping EventComputedObservationQuorumSuccess from old round", commontypes.LogFields{
			"seqNr":   outgen.sharedState.seqNr,
			"evEpoch": ev.Epoch,
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.leaderState.phase != outgenLeaderPhaseSentRoundStart {
		outgen.logger.Debug("dropping EventComputedObservationQuorumSuccess, wrong phase", commontypes.LogFields{
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
		nil,
	})

	outgen.leaderState.observationPool.ReapCompleted(outgen.sharedState.seqNr)
}
