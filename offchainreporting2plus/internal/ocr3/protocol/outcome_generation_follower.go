package protocol

import (
	"context"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common/pool"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type outgenFollowerPhase string

const (
	outgenFollowerPhaseUnknown                   outgenFollowerPhase = "unknown"
	outgenFollowerPhaseNewEpoch                  outgenFollowerPhase = "newEpoch"
	outgenFollowerPhaseNewRound                  outgenFollowerPhase = "newRound"
	outgenFollowerPhaseBackgroundObservation     outgenFollowerPhase = "backgroundObservation"
	outgenFollowerPhaseSentObservation           outgenFollowerPhase = "sentObservation"
	outgenFollowerPhaseBackgroundProposalOutcome outgenFollowerPhase = "backgroundProposalOutcome"
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

	if msg.EpochStartProof.HighestCertified.IsGenesis() {
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentFollowerRound()
	} else if commitQC, ok := msg.EpochStartProof.HighestCertified.(*CertifiedCommit); ok {
		outgen.commit(*commitQC)
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentFollowerRound()
	} else {
		// We're dealing with a re-proposal from a failed epoch

		prepareQC, ok := msg.EpochStartProof.HighestCertified.(*CertifiedPrepare)
		if !ok {
			outgen.logger.Critical("cast to CertifiedPrepare failed while processing MessageEpochStart", commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
			})
			return
		}

		// We don't know the actual inputs, so we always use the empty OutcomeInputsDigest
		// in case of a re-proposal.
		outcomeInputsDigest := OutcomeInputsDigest{}

		outcomeDigest := MakeOutcomeDigest(prepareQC.Outcome)

		prepareSignature, err := MakePrepareSignature(
			outgen.ID(),
			prepareQC.SeqNr,
			outcomeInputsDigest,
			outcomeDigest,
			outgen.offchainKeyring.OffchainSign,
		)
		if err != nil {
			outgen.logger.Critical("failed to sign Prepare", commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
				"error": err,
			})
			return
		}

		outgen.sharedState.firstSeqNrOfEpoch = prepareQC.SeqNr + 1
		outgen.sharedState.seqNr = prepareQC.SeqNr

		outgen.followerState.phase = outgenFollowerPhaseSentPrepare
		outgen.followerState.outcome = outcomeAndDigests{
			prepareQC.Outcome,
			outcomeInputsDigest,
			outcomeDigest,
		}
		outgen.logger.Debug("broadcasting MessagePrepare (reproposal)", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		outgen.netSender.Broadcast(MessagePrepare[RI]{
			outgen.sharedState.e,
			prepareQC.SeqNr,
			prepareSignature,
		})
	}
}

func (outgen *outcomeGenerationState[RI]) startSubsequentFollowerRound() {
	outgen.sharedState.seqNr = outgen.sharedState.committedSeqNr + 1

	outgen.followerState.phase = outgenFollowerPhaseNewRound
	outgen.followerState.query = nil
	outgen.followerState.outcome = outcomeAndDigests{}

	outgen.tryProcessRoundStartPool()
}

func (outgen *outcomeGenerationState[RI]) messageRoundStart(msg MessageRoundStart[RI], sender commontypes.OracleID) {
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
	outctx := outgen.OutcomeCtx(outgen.sharedState.seqNr)

	outgen.followerState.phase = outgenFollowerPhaseBackgroundObservation
	outgen.followerState.query = &msg.Query

	outgen.telemetrySender.RoundStarted(
		outgen.config.ConfigDigest,
		outctx.Epoch,
		outctx.SeqNr,
		outctx.Round,
		outgen.sharedState.l,
	)

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		query := *outgen.followerState.query
		outgen.subs.Go(func() {
			outgen.backgroundObservation(ctx, logger, outctx, query)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundObservation(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	outctx ocr3types.OutcomeContext,
	query types.Query,
) {
	observation, ok := callPluginFromOutcomeGenerationBackground[types.Observation, RI](
		ctx,
		logger,
		"Observation",
		outgen.config.MaxDurationObservation,
		outctx,
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Observation, error) {
			return outgen.reportingPlugin.Observation(ctx, outctx, query)
		},
	)
	if !ok {
		return
	}

	select {
	case outgen.chLocalEvent <- EventComputedObservation[RI]{
		outctx.Epoch,
		outctx.SeqNr,
		query,
		observation,
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedObservation(ev EventComputedObservation[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("dropping EventComputedObservation from old round", commontypes.LogFields{
			"seqNr":   outgen.sharedState.seqNr,
			"evEpoch": ev.Epoch,
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.followerState.phase != outgenFollowerPhaseBackgroundObservation {
		outgen.logger.Debug("dropping EventComputedObservation, wrong phase", commontypes.LogFields{
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

	outgen.followerState.phase = outgenFollowerPhaseBackgroundProposalOutcome

	{
		ctx := outgen.epochCtx
		logger := outgen.logger
		outctx := outgen.OutcomeCtx(outgen.sharedState.seqNr)
		ogid := outgen.ID()
		query := *outgen.followerState.query
		outgen.subs.Go(func() {
			outgen.backgroundProposalOutcome(
				ctx,
				logger,
				ogid,
				outctx,
				query,
				msg.AttributedSignedObservations,
			)
		})
	}
}

func (outgen *outcomeGenerationState[RI]) backgroundProposalOutcome(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	ogid OutcomeGenerationID,
	outctx ocr3types.OutcomeContext,
	query types.Query,
	asos []AttributedSignedObservation,
) {

	attributedObservations := []types.AttributedObservation{}
	{
		seen := map[commontypes.OracleID]bool{}
		for _, aso := range asos {
			if !(0 <= int(aso.Observer) && int(aso.Observer) < outgen.config.N()) {
				logger.Warn("dropping MessageProposal that contains signed observation with invalid observer", commontypes.LogFields{
					"seqNr":           outctx.SeqNr,
					"invalidObserver": aso.Observer,
				})
				return
			}

			if seen[aso.Observer] {
				logger.Warn("dropping MessageProposal that contains duplicate signed observation", commontypes.LogFields{
					"seqNr": outctx.SeqNr,
				})
				return
			}

			seen[aso.Observer] = true

			if err := aso.SignedObservation.Verify(ogid, outctx.SeqNr, query, outgen.config.OracleIdentities[aso.Observer].OffchainPublicKey); err != nil {
				logger.Warn("dropping MessageProposal that contains signed observation with invalid signature", commontypes.LogFields{
					"seqNr": outctx.SeqNr,
					"error": err,
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
					return outgen.reportingPlugin.ValidateObservation(
						ctx,
						outctx,
						query,
						types.AttributedObservation{aso.SignedObservation.Observation, aso.Observer},
					), nil
				},
			)
			if !ok {
				logger.Error("dropping MessageProposal containing observation that could not be validated", commontypes.LogFields{
					"seqNr":    outctx.SeqNr,
					"observer": aso.Observer,
				})
				return
			}
			if err != nil {
				logger.Warn("dropping MessageProposal that contains an invalid observation", commontypes.LogFields{
					"seqNr":    outctx.SeqNr,
					"error":    err,
					"observer": aso.Observer,
				})
				return
			}

			attributedObservations = append(attributedObservations, types.AttributedObservation{
				aso.SignedObservation.Observation,
				aso.Observer,
			})
		}

		observationQuorum, ok := callPluginFromOutcomeGenerationBackground[bool, RI](
			ctx,
			logger,
			"ObservationQuorum",
			0, // ObservationQuorum is a pure function and should finish "instantly"
			outctx,
			func(ctx context.Context, outctx ocr3types.OutcomeContext) (bool, error) {
				return outgen.reportingPlugin.ObservationQuorum(ctx, outctx, query, attributedObservations)
			},
		)

		if !ok {
			return
		}

		if !observationQuorum {
			logger.Warn("dropping MessageProposal that doesn't achieve observation quorum", commontypes.LogFields{
				"seqNr": outctx.SeqNr,
			})
			return
		}

		if seen[outgen.id] {
			outgen.metrics.includedObservationsTotal.Inc()
		}
	}

	outcomeInputsDigest := MakeOutcomeInputsDigest(
		ogid,
		outgen.sharedState.committedOutcome,
		outctx.SeqNr,
		query,
		attributedObservations,
	)

	outcome, ok := callPluginFromOutcomeGenerationBackground[ocr3types.Outcome, RI](
		ctx,
		logger,
		"Outcome",
		0, // Outcome is a pure function and should finish "instantly"
		outctx,
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (ocr3types.Outcome, error) {
			return outgen.reportingPlugin.Outcome(ctx, outctx, query, attributedObservations)
		},
	)
	if !ok {
		return
	}

	outcomeDigest := MakeOutcomeDigest(outcome)

	select {
	case outgen.chLocalEvent <- EventComputedProposalOutcome[RI]{
		outctx.Epoch,
		outctx.SeqNr,
		outcomeAndDigests{
			outcome,
			outcomeInputsDigest,
			outcomeDigest,
		},
	}:
	case <-ctx.Done():
	}
}

func (outgen *outcomeGenerationState[RI]) eventComputedProposalOutcome(ev EventComputedProposalOutcome[RI]) {
	if ev.Epoch != outgen.sharedState.e || ev.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("dropping EventComputedProposalOutcome from old round", commontypes.LogFields{
			"seqNr":   outgen.sharedState.seqNr,
			"evEpoch": ev.Epoch,
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	if outgen.followerState.phase != outgenFollowerPhaseBackgroundProposalOutcome {
		outgen.logger.Debug("dropping EventComputedProposalOutcome, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	prepareSignature, err := MakePrepareSignature(
		outgen.ID(),
		outgen.sharedState.seqNr,
		ev.outcomeAndDigests.InputsDigest,
		ev.outcomeAndDigests.Digest,
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
	outgen.followerState.outcome = ev.outcomeAndDigests

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
			outgen.followerState.outcome.InputsDigest,
			outgen.followerState.outcome.Digest,
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
		outgen.followerState.outcome.Digest,
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
		outgen.followerState.outcome.InputsDigest,
		outgen.followerState.outcome.Outcome,
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
			outgen.followerState.outcome.Digest,
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
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		outgen.followerState.outcome.Outcome,
		commitQuorumCertificate,
	})
	if outgen.id == outgen.sharedState.l {
		outgen.metrics.ledCommittedRoundsTotal.Inc()
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

func (outgen *outcomeGenerationState[RI]) commit(commit CertifiedCommit) {
	if commit.SeqNr < outgen.sharedState.committedSeqNr {
		outgen.logger.Critical("assumption violation, commitSeqNr is less than committedSeqNr", commontypes.LogFields{
			"commitSeqNr":    commit.SeqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
		return
	}

	if commit.SeqNr <= outgen.sharedState.committedSeqNr {

		outgen.logger.Debug("skipping commit of already committed outcome", commontypes.LogFields{
			"commitSeqNr":    commit.SeqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
	} else {

		if !outgen.persistAndUpdateCertIfGreater(&commit) {
			return
		}
		outgen.sharedState.committedSeqNr = commit.SeqNr
		outgen.sharedState.committedOutcome = commit.Outcome
		outgen.metrics.committedSeqNr.Set(float64(commit.SeqNr))

		outgen.logger.Debug("✅ committed outcome", commontypes.LogFields{
			"seqNr": commit.SeqNr,
		})

		select {
		case outgen.chOutcomeGenerationToReportAttestation <- EventCommittedOutcome[RI]{commit}:
		case <-outgen.ctx.Done():
			return
		}
	}

	outgen.followerState.roundStartPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.proposalPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.preparePool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.commitPool.ReapCompleted(outgen.sharedState.committedSeqNr)
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
