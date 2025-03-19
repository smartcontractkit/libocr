package protocol

import (
	"context"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/pool"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// Identifies an instance of the outcome generation protocol
type OutcomeGenerationID struct {
	ConfigDigest types.ConfigDigest
	Epoch        uint64
}

const futureMessageBufferSize = 10 // big enough for a couple of full rounds of outgen protocol
const poolSize = 3

func RunOutcomeGeneration[RI any](
	ctx context.Context,

	chNetToOutcomeGeneration <-chan MessageToOutcomeGenerationWithSender[RI],
	chPacemakerToOutcomeGeneration <-chan EventToOutcomeGeneration[RI],
	chOutcomeGenerationToPacemaker chan<- EventToPacemaker[RI],
	chOutcomeGenerationToReportAttestation chan<- EventToReportAttestation[RI],
	chOutcomeGenerationToStatePersistence chan<- EventToStatePersistence[RI],
	chStatePersistenceToOutcomeGeneration <-chan EventToOutcomeGeneration[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	kvStore ocr3_1types.KeyValueStore,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	reportingPlugin ocr3_1types.ReportingPlugin[RI],
	telemetrySender TelemetrySender,

	restoredCert CertifiedPrepareOrCommit,
	restoredHighestCommittedToKVSeqNr uint64,
) {

	outgen := outcomeGenerationState[RI]{
		ctx:  ctx,
		subs: subprocesses.Subprocesses{},

		chLocalEvent:                           make(chan EventToOutcomeGeneration[RI]),
		chNetToOutcomeGeneration:               chNetToOutcomeGeneration,
		chPacemakerToOutcomeGeneration:         chPacemakerToOutcomeGeneration,
		chOutcomeGenerationToPacemaker:         chOutcomeGenerationToPacemaker,
		chOutcomeGenerationToReportAttestation: chOutcomeGenerationToReportAttestation,
		chOutcomeGenerationToStatePersistence:  chOutcomeGenerationToStatePersistence,
		chStatePersistenceToOutcomeGeneration:  chStatePersistenceToOutcomeGeneration,
		config:                                 config,
		database:                               database,
		id:                                     id,
		kvStore:                                kvStore,
		localConfig:                            localConfig,
		logger:                                 logger.MakeUpdated(commontypes.LogFields{"proto": "outgen"}),
		metrics:                                newOutcomeGenerationMetrics(metricsRegisterer, logger),
		netSender:                              netSender,
		offchainKeyring:                        offchainKeyring,
		reportingPlugin:                        reportingPlugin,
		telemetrySender:                        telemetrySender,
	}
	outgen.run(restoredCert, restoredHighestCommittedToKVSeqNr)
}

type outcomeGenerationState[RI any] struct {
	ctx  context.Context
	subs subprocesses.Subprocesses

	chLocalEvent                           chan EventToOutcomeGeneration[RI]
	chNetToOutcomeGeneration               <-chan MessageToOutcomeGenerationWithSender[RI]
	chPacemakerToOutcomeGeneration         <-chan EventToOutcomeGeneration[RI]
	chReportAttestationToOutcomeGeneration <-chan EventToOutcomeGeneration[RI]
	chOutcomeGenerationToPacemaker         chan<- EventToPacemaker[RI]
	chOutcomeGenerationToReportAttestation chan<- EventToReportAttestation[RI]
	chOutcomeGenerationToStatePersistence  chan<- EventToStatePersistence[RI]
	chStatePersistenceToOutcomeGeneration  <-chan EventToOutcomeGeneration[RI]
	config                                 ocr3config.SharedConfig
	database                               Database
	id                                     commontypes.OracleID
	kvStore                                ocr3_1types.KeyValueStore
	localConfig                            types.LocalConfig
	logger                                 loghelper.LoggerWithContext
	metrics                                *outcomeGenerationMetrics
	netSender                              NetworkSender[RI]
	offchainKeyring                        types.OffchainKeyring
	reportingPlugin                        ocr3_1types.ReportingPlugin[RI]
	telemetrySender                        TelemetrySender

	epochCtx         context.Context
	epochCtxCancel   context.CancelFunc
	bufferedMessages []*MessageBuffer[RI]
	leaderState      leaderState[RI]
	followerState    followerState[RI]
	sharedState      sharedState
}

type leaderState[RI any] struct {
	phase outgenLeaderPhase

	epochStartRequests map[commontypes.OracleID]*epochStartRequest[RI]

	readyToStartRound bool
	tRound            <-chan time.Time

	query           types.Query
	observationPool *pool.Pool[SignedObservation]
	tGrace          <-chan time.Time
}

type epochStartRequest[RI any] struct {
	message MessageEpochStartRequest[RI]
	bad     bool
}

type followerState[RI any] struct {
	phase outgenFollowerPhase

	tInitial <-chan time.Time

	roundStartPool *pool.Pool[MessageRoundStart[RI]]

	query *types.Query

	proposalPool *pool.Pool[MessageProposal[RI]]

	roundInfo roundInfo[RI]

	// lock
	cert CertifiedPrepareOrCommit

	preparePool *pool.Pool[PrepareSignature]
	commitPool  *pool.Pool[CommitSignature]
}

type roundInfo[RI any] struct {
	inputs                     StateTransitionInputs
	reportsPlusPrecursor       ocr3_1types.ReportsPlusPrecursor
	inputsDigest               StateTransitionInputsDigest
	outputDigest               StateTransitionOutputDigest
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest
	commitQuorumCertificate    []AttributedCommitSignature
}

type sharedState struct {
	e uint64               // Current epoch number
	l commontypes.OracleID // Current leader number

	firstSeqNrOfEpoch uint64
	seqNr             uint64
	observationQuorum *int
	committedSeqNr    uint64

	kvStoreTxn              *kvStoreTxn
	committedToKVStoreSeqNr uint64 // The sequence number of the key-value store
}

func (outgen *outcomeGenerationState[RI]) run(restoredCert CertifiedPrepareOrCommit, restoredCommittedToKVStoreSeqNr uint64) {
	var restoredCommitedSeqNr uint64
	if restoredCert != nil {
		if commitQC, ok := restoredCert.(*CertifiedCommit); ok {
			restoredCommitedSeqNr = commitQC.SeqNr()
		} else if prepareQc, ok := restoredCert.(*CertifiedPrepare); ok {
			if prepareQc.SeqNr() > 1 {
				restoredCommitedSeqNr = prepareQc.SeqNr() - 1
			}
		}
	}

	outgen.logger.Info("OutcomeGeneration: running", commontypes.LogFields{
		"restoredCommittedSeqNr":          restoredCommitedSeqNr,
		"restoredCommittedToKVStoreSeqNr": restoredCommittedToKVStoreSeqNr,
	})

	// Initialization
	outgen.epochCtx, outgen.epochCtxCancel = context.WithCancel(outgen.ctx)

	for i := 0; i < outgen.config.N(); i++ {
		outgen.bufferedMessages = append(outgen.bufferedMessages, NewMessageBuffer[RI](futureMessageBufferSize))
	}

	outgen.leaderState = leaderState[RI]{
		outgenLeaderPhaseUnknown,
		map[commontypes.OracleID]*epochStartRequest[RI]{},
		false,
		nil,
		nil,
		nil,
		nil,
	}

	outgen.followerState = followerState[RI]{
		outgenFollowerPhaseUnknown,
		nil,
		nil,
		nil,
		nil,
		roundInfo[RI]{},
		restoredCert,
		nil,
		nil,
	}

	outgen.sharedState = sharedState{
		0,
		0,

		0,
		restoredCommitedSeqNr,
		nil,
		restoredCommitedSeqNr,
		nil,
		restoredCommittedToKVStoreSeqNr,
	}

	// Event Loop
	chDone := outgen.ctx.Done()
	for {
		select {
		case ev := <-outgen.chLocalEvent:
			ev.processOutcomeGeneration(outgen)
		case msg := <-outgen.chNetToOutcomeGeneration:
			outgen.messageToOutcomeGeneration(msg)
		case ev := <-outgen.chPacemakerToOutcomeGeneration:
			ev.processOutcomeGeneration(outgen)
		case ev := <-outgen.chReportAttestationToOutcomeGeneration:
			ev.processOutcomeGeneration(outgen)
		case ev := <-outgen.chStatePersistenceToOutcomeGeneration:
			ev.processOutcomeGeneration(outgen)
		case <-outgen.followerState.tInitial:
			outgen.eventTInitialTimeout()
		case <-outgen.leaderState.tGrace:
			outgen.eventTGraceTimeout()
		case <-outgen.leaderState.tRound:
			outgen.eventTRoundTimeout()
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			outgen.logger.Info("OutcomeGeneration: winding down", commontypes.LogFields{
				"e": outgen.sharedState.e,
				"l": outgen.sharedState.l,
			})
			outgen.subs.Wait()
			outgen.metrics.Close()
			outgen.logger.Info("OutcomeGeneration: exiting", commontypes.LogFields{
				"e": outgen.sharedState.e,
				"l": outgen.sharedState.l,
			})
			return
		default:
		}
	}
}

func (outgen *outcomeGenerationState[RI]) messageToOutcomeGeneration(msg MessageToOutcomeGenerationWithSender[RI]) {
	msgEpoch := msg.msg.epoch()
	if msgEpoch < outgen.sharedState.e {
		// drop
		outgen.logger.Debug("dropping message for past epoch", commontypes.LogFields{
			"epoch":    outgen.sharedState.e,
			"msgEpoch": msgEpoch,
			"sender":   msg.sender,
		})
	} else if msgEpoch == outgen.sharedState.e {
		msg.msg.processOutcomeGeneration(outgen, msg.sender)
	} else {
		outgen.bufferedMessages[msg.sender].Push(msg.msg)
		outgen.logger.Trace("buffering message for future epoch", commontypes.LogFields{
			"msgEpoch": msgEpoch,
			"sender":   msg.sender,
		})
	}
}

func (outgen *outcomeGenerationState[RI]) unbufferMessages() {
	outgen.logger.Trace("getting messages for new epoch", nil)
	for i, buffer := range outgen.bufferedMessages {
		sender := commontypes.OracleID(i)
		for buffer.Length() > 0 {
			msg := buffer.Peek()
			msgEpoch := msg.epoch()
			if msgEpoch < outgen.sharedState.e {
				buffer.Pop()
				outgen.logger.Debug("unbuffered and dropped message", commontypes.LogFields{
					"msgEpoch": msgEpoch,
					"sender":   sender,
				})
			} else if msgEpoch == outgen.sharedState.e {
				buffer.Pop()
				outgen.logger.Trace("unbuffered message for new epoch", commontypes.LogFields{
					"msgEpoch": msgEpoch,
					"sender":   sender,
				})
				msg.processOutcomeGeneration(outgen, sender)
			} else { // msgEpoch > e
				// this and all subsequent messages are for future epochs
				// leave them in the buffer
				break
			}
		}
	}
	outgen.logger.Trace("done unbuffering messages for new epoch", nil)
}

func (outgen *outcomeGenerationState[RI]) eventReplayVerifiedStateTransition(ev EventReplayVerifiedStateTransition[RI]) {
	outgen.logger.Debug("received EventReplayVerifiedStateTransition", commontypes.LogFields{
		"stateTransitionBlockSeqNr": ev.AttestedStateTransitionBlock.StateTransitionBlock.SeqNr(),
	})

	outgen.replayStateTransition(
		ev.AttestedStateTransitionBlock.StateTransitionBlock.StateTransitionInputs,
		ev.AttestedStateTransitionBlock.StateTransitionBlock.StateTransitionOutputDigest,
		ev.AttestedStateTransitionBlock.StateTransitionBlock.ReportsPrecursorDigest,
		ev.AttestedStateTransitionBlock.AttributedSignatures)
}

func (outgen *outcomeGenerationState[RI]) eventProduceStateTransition(ev EventProduceStateTransition[RI]) {
	outgen.logger.Debug("received EventProduceStateTransition", commontypes.LogFields{
		"seqNr": ev.RoundCtx.SeqNr,
	})
	outgen.sharedState.kvStoreTxn = &ev.Txn
	outgen.produceStateTransition(
		ev.RoundCtx,
		ev.Txn.transaction,
		ev.Query,
		ev.Asos,
		ev.Prepared,
		ev.StateTransitionOutputDigest,
		ev.ReportsPlusPrecursorDigest,
		ev.CommitQC)
}

func (outgen *outcomeGenerationState[RI]) eventAcknowledgedComputedStateTransition(ev EventAcknowledgedComputedStateTransition[RI]) {
	outgen.logger.Debug("received EventAcknowledgedComputedStateTransition", commontypes.LogFields{
		"stateTransitionSeqNr": ev.SeqNr,
	})
	outgen.processAcknowledgedComputedStateTransition(ev)
}

// state as of ev.SeqNr has been written to the key-value store
func (outgen *outcomeGenerationState[RI]) eventCommittedKVTransaction(ev EventCommittedKVTransaction[RI]) {
	outgen.logger.Debug("received EventCommittedKVTransaction", commontypes.LogFields{
		"evSeqNr": ev.SeqNr,
	})
	if ev.SeqNr != outgen.sharedState.committedSeqNr {
		outgen.logger.Critical("we received a CommittedKVTransaction event out of order", commontypes.LogFields{
			"eventSeqNr":     ev.SeqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
	}
	outgen.sharedState.committedToKVStoreSeqNr = ev.SeqNr
	outgen.sharedState.kvStoreTxn = nil
	select {
	case outgen.chOutcomeGenerationToStatePersistence <- EventAcknowledgedCommittedKVTransaction[RI]{outgen.sharedState.committedToKVStoreSeqNr}:
	case <-outgen.ctx.Done():
		return
	}

	if uint64(outgen.config.RMax) <= outgen.sharedState.committedToKVStoreSeqNr-outgen.sharedState.firstSeqNrOfEpoch+1 {
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

func (outgen *outcomeGenerationState[RI]) eventNewEpochStart(ev EventNewEpochStart[RI]) {
	// Initialization
	outgen.logger.Info("starting new epoch", commontypes.LogFields{
		"epoch": ev.Epoch,
	})

	outgen.epochCtxCancel()
	outgen.epochCtx, outgen.epochCtxCancel = context.WithCancel(outgen.ctx)

	outgen.sharedState.e = ev.Epoch
	outgen.sharedState.l = Leader(outgen.sharedState.e, outgen.config.N(), outgen.config.LeaderSelectionKey())

	outgen.logger = outgen.logger.MakeUpdated(commontypes.LogFields{
		"e": outgen.sharedState.e,
		"l": outgen.sharedState.l,
	})

	outgen.sharedState.firstSeqNrOfEpoch = 0
	outgen.sharedState.seqNr = 0

	// In case we have an open transaction that is not being replayed we should discard it.
	// Note that relying on cancelling the epoch context in order to discard the transaction
	// through backgroundStateTransition is not enough. The state transition might have been
	// already computed but the transaction might have not been committed yet.
	// If the transaction is being replayed it will commit anyway so we don't have to cancel it.
	if outgen.sharedState.kvStoreTxn != nil && !outgen.sharedState.kvStoreTxn.replayed {
		select {
		case outgen.chOutcomeGenerationToStatePersistence <- EventDiscardKVTransaction[RI]{outgen.sharedState.kvStoreTxn.transaction.SeqNr()}:
		case <-outgen.ctx.Done():
			return
		}
		outgen.sharedState.kvStoreTxn = nil
	}

	outgen.followerState.phase = outgenFollowerPhaseNewEpoch
	outgen.followerState.tInitial = time.After(outgen.config.DeltaInitial)
	outgen.followerState.roundInfo = roundInfo[RI]{}

	outgen.followerState.roundStartPool = pool.NewPool[MessageRoundStart[RI]](poolSize)
	outgen.followerState.proposalPool = pool.NewPool[MessageProposal[RI]](poolSize)
	outgen.followerState.preparePool = pool.NewPool[PrepareSignature](poolSize)
	outgen.followerState.commitPool = pool.NewPool[CommitSignature](poolSize)

	outgen.leaderState.phase = outgenLeaderPhaseNewEpoch
	outgen.leaderState.epochStartRequests = map[commontypes.OracleID]*epochStartRequest[RI]{}
	outgen.leaderState.readyToStartRound = false
	outgen.leaderState.observationPool = pool.NewPool[SignedObservation](1) // only one observation per sender & round, and we do not need to worry about observations from the future
	outgen.leaderState.tGrace = nil

	var highestCertified CertifiedPrepareOrCommit
	var highestCertifiedTimestamp HighestCertifiedTimestamp
	highestCertified = outgen.followerState.cert
	highestCertifiedTimestamp = outgen.followerState.cert.Timestamp()

	signedHighestCertifiedTimestamp, err := MakeSignedHighestCertifiedTimestamp(
		outgen.ID(outgen.sharedState.e),
		highestCertifiedTimestamp,
		outgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		outgen.logger.Error("error signing timestamp", commontypes.LogFields{
			"error": err,
		})
		return
	}

	outgen.logger.Info("sending MessageEpochStartRequest to leader", commontypes.LogFields{
		"highestCertifiedTimestamp": highestCertifiedTimestamp,
	})
	outgen.netSender.SendTo(MessageEpochStartRequest[RI]{
		outgen.sharedState.e,
		highestCertified,
		signedHighestCertifiedTimestamp,
	}, outgen.sharedState.l)

	if outgen.id == outgen.sharedState.l {
		outgen.leaderState.tRound = time.After(outgen.config.DeltaRound)
	}

	outgen.unbufferMessages()
}

func (outgen *outcomeGenerationState[RI]) ID(epoch uint64) OutcomeGenerationID {
	return OutcomeGenerationID{outgen.config.ConfigDigest, epoch}
}

func (outgen *outcomeGenerationState[RI]) RoundCtx(seqNr uint64) ocr3_1types.RoundContext {
	if seqNr != outgen.sharedState.committedToKVStoreSeqNr+1 {
		outgen.logger.Critical("assumption violation, seqNr isn't successor to committedSeqToKVSeqNr", commontypes.LogFields{
			"seqNr":                   seqNr,
			"committedToKVStoreSeqNr": outgen.sharedState.committedToKVStoreSeqNr,
		})
		panic("")
	}
	return ocr3_1types.RoundContext{
		seqNr,
		outgen.sharedState.e,
		seqNr - outgen.sharedState.firstSeqNrOfEpoch + 1,
	}
}

func callPluginFromOutcomeGenerationBackground[T any](
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	name string,
	recommendedMaxDuration time.Duration,
	roundCtx ocr3_1types.RoundContext,
	f func(context.Context, ocr3_1types.RoundContext) (T, error),
) (T, bool) {
	return common.CallPluginFromBackground[T](
		ctx,
		logger,
		commontypes.LogFields{
			"seqNr": roundCtx.SeqNr,
			"round": roundCtx.Round, // nolint: staticcheck
		},
		name,
		recommendedMaxDuration,
		func(ctx context.Context) (T, error) {
			return f(ctx, roundCtx)
		},
	)
}
