package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol/pool"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

const futureMessageBufferSize = 10 // big enough for a couple of full rounds of repgen protocol

func RunReportGeneration[RI any](
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	chNetToReportGeneration <-chan MessageToReportGenerationWithSender[RI],
	chPacemakerToReportGeneration <-chan EventToReportGeneration[RI],
	chReportGenerationToPacemaker chan<- EventToPacemaker[RI],
	chReportGenerationToReportFinalization chan<- EventToReportFinalization[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	reportingPlugin ocr3types.OCR3Plugin[RI],
	telemetrySender TelemetrySender,

	restoredCert CertifiedPrepareOrCommit,
) {
	repgen := reportGenerationState[RI]{
		ctx:          ctx,
		subprocesses: subprocesses,

		chNetToReportGeneration:                chNetToReportGeneration,
		chPacemakerToReportGeneration:          chPacemakerToReportGeneration,
		chReportGenerationToPacemaker:          chReportGenerationToPacemaker,
		chReportGenerationToReportFinalization: chReportGenerationToReportFinalization,
		config:                                 config,
		database:                               database,
		id:                                     id,
		localConfig:                            localConfig,
		logger:                                 logger.MakeUpdated(commontypes.LogFields{"proto": "repgen"}),
		netSender:                              netSender,
		offchainKeyring:                        offchainKeyring,
		reportingPlugin:                        reportingPlugin,
		telemetrySender:                        telemetrySender,
	}
	repgen.run(restoredCert)
}

type reportGenerationState[RI any] struct {
	ctx          context.Context
	subprocesses *subprocesses.Subprocesses

	chNetToReportGeneration                <-chan MessageToReportGenerationWithSender[RI]
	chPacemakerToReportGeneration          <-chan EventToReportGeneration[RI]
	chReportGenerationToPacemaker          chan<- EventToPacemaker[RI]
	chReportGenerationToReportFinalization chan<- EventToReportFinalization[RI]
	config                                 ocr3config.SharedConfig
	database                               Database
	e                                      uint64 // Current epoch number
	id                                     commontypes.OracleID
	l                                      commontypes.OracleID // Current leader number
	localConfig                            types.LocalConfig
	logger                                 loghelper.LoggerWithContext
	netSender                              NetworkSender[RI]
	offchainKeyring                        types.OffchainKeyring
	reportingPlugin                        ocr3types.OCR3Plugin[RI]
	telemetrySender                        TelemetrySender

	bufferedMessages []*MessageBuffer[RI]
	leaderState      leaderState
	followerState    followerState[RI]
}

type leaderState struct {
	phase repgenLeaderPhase

	startRoundQuorumCertificate StartEpochProof

	readyToStartRound bool
	tRound            <-chan time.Time

	query        types.Query
	observations map[commontypes.OracleID]*SignedObservation
	tGrace       <-chan time.Time
}

type followerState[RI any] struct {
	phase repgenFollowerPhase

	firstSeqNrOfEpoch uint64

	seqNr uint64

	observeReqPool *pool.Pool[MessageStartRound[RI]]

	query *types.Query

	proposePool *pool.Pool[MessagePropose[RI]]

	currentOutcomeInputsDigest OutcomeInputsDigest
	currentOutcome             ocr3types.Outcome
	currentOutcomeDigest       OutcomeDigest

	// lock
	cert CertifiedPrepareOrCommit

	preparePool *pool.Pool[PrepareSignature]
	commitPool  *pool.Pool[CommitSignature]

	deliveredSeqNr   uint64
	deliveredOutcome ocr3types.Outcome
}

// Run starts the event loop for the report-generation protocol
func (repgen *reportGenerationState[RI]) run(restoredCert CertifiedPrepareOrCommit) {
	repgen.logger.Info("Running ReportGeneration", nil)

	for i := 0; i < repgen.config.N(); i++ {
		repgen.bufferedMessages = append(repgen.bufferedMessages, NewMessageBuffer[RI](futureMessageBufferSize))
	}

	// Initialization
	repgen.leaderState = leaderState{
		repgenLeaderPhaseUnknown,
		StartEpochProof{
			nil,
			nil,
		},
		false,
		nil,
		nil,
		nil,
		nil,
	}

	repgen.followerState = followerState[RI]{
		repgenFollowerPhaseUnknown,
		0,
		0,
		nil,
		nil,
		nil,
		OutcomeInputsDigest{},
		nil,
		OutcomeDigest{},
		restoredCert,
		nil,
		nil,

		0,
		nil,
	}

	// AXE
	// if repgen.id == repgen.l {
	// 	repgen.startRound()
	// 	repgen.startRound()
	// }

	// Event Loop
	chDone := repgen.ctx.Done()
	for {
		select {
		case msg := <-repgen.chNetToReportGeneration:
			msg.msg.processReportGeneration(repgen, msg.sender)
		case ev := <-repgen.chPacemakerToReportGeneration:
			ev.processReportGeneration(repgen)
		case <-repgen.leaderState.tGrace:
			repgen.eventTGraceTimeout()
		case <-repgen.leaderState.tRound:
			repgen.eventTRoundTimeout()
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			repgen.logger.Info("ReportGeneration: exiting", commontypes.LogFields{
				"e": repgen.e,
				"l": repgen.l,
			})
			return
		default:
		}
	}
}

func (repgen *reportGenerationState[RI]) messageToReportGeneration(msg MessageToReportGeneration[RI], sender commontypes.OracleID) {
	msgEpoch := msg.epoch()
	if msgEpoch < repgen.e {
		// drop
		repgen.logger.Debug("dropping message for past epoch", commontypes.LogFields{
			"epoch":    repgen.e,
			"msgEpoch": msgEpoch,
			"sender":   sender,
		})
	} else if msgEpoch == repgen.e {
		msg.processReportGeneration(repgen, sender)
	} else {
		repgen.bufferedMessages[sender].Push(msg)
		repgen.logger.Trace("buffering message for future epoch", commontypes.LogFields{
			"epoch":    repgen.e,
			"msgEpoch": msgEpoch,
			"sender":   sender,
		})
	}
}

func (repgen *reportGenerationState[RI]) unbufferMessages() {
	repgen.logger.Trace("getting messages for new epoch", commontypes.LogFields{
		"epoch": repgen.e,
	})
	for i, buffer := range repgen.bufferedMessages {
		sender := commontypes.OracleID(i)
		for {
			msg := buffer.Peek()
			if msg == nil {
				// no messages left in buffer
				break
			}
			msgEpoch := (*msg).epoch()
			if msgEpoch < repgen.e {
				buffer.Pop()
				repgen.logger.Debug("unbuffered and dropped message", commontypes.LogFields{
					"epoch":    repgen.e,
					"msgEpoch": msgEpoch,
					"sender":   sender,
				})
			} else if msgEpoch == repgen.e {
				buffer.Pop()
				repgen.logger.Trace("unbuffered messages for new epoch", commontypes.LogFields{
					"epoch":    repgen.e,
					"msgEpoch": msgEpoch,
					"sender":   sender,
				})
				(*msg).processReportGeneration(repgen, sender)
			} else { // msgEpoch > e
				// this and all subsequent messages are for future epochs
				// leave them in the buffer
				break
			}
		}
	}
	repgen.logger.Trace("done unbuffering messages for new epoch", commontypes.LogFields{
		"epoch": repgen.e,
	})
}

func (repgen *reportGenerationState[RI]) eventStartNewEpoch(ev EventStartNewEpoch[RI]) {
	// Initialization
	repgen.logger.Info("Starting new epoch", commontypes.LogFields{
		"epoch": ev.Epoch,
	})

	repgen.e = ev.Epoch
	repgen.l = Leader(repgen.e, repgen.config.N(), repgen.config.LeaderSelectionKey())

	repgen.logger = repgen.logger.MakeUpdated(commontypes.LogFields{
		"e": repgen.e,
		"l": repgen.l,
	})

	repgen.followerState.phase = repgenFollowerPhaseNewEpoch
	repgen.followerState.firstSeqNrOfEpoch = 0
	repgen.followerState.seqNr = 0
	repgen.followerState.currentOutcomeInputsDigest = OutcomeInputsDigest{}
	repgen.followerState.currentOutcome = nil
	repgen.followerState.currentOutcomeDigest = OutcomeDigest{}

	repgen.followerState.observeReqPool = pool.NewPool[MessageStartRound[RI]](10)
	repgen.followerState.proposePool = pool.NewPool[MessagePropose[RI]](10)
	repgen.followerState.preparePool = pool.NewPool[PrepareSignature](10)
	repgen.followerState.commitPool = pool.NewPool[CommitSignature](10)

	repgen.leaderState.phase = repgenLeaderPhaseNewEpoch

	repgen.leaderState.startRoundQuorumCertificate = StartEpochProof{
		nil,
		nil,
	}
	repgen.leaderState.readyToStartRound = false
	repgen.leaderState.tGrace = nil

	var highestCertified CertifiedPrepareOrCommit
	var highestCertifiedTimestamp HighestCertifiedTimestamp
	highestCertified = repgen.followerState.cert
	highestCertifiedTimestamp = repgen.followerState.cert.Timestamp()

	signedHighestCertifiedTimestamp, err := MakeSignedHighestCertifiedTimestamp(
		repgen.Timestamp(),
		highestCertifiedTimestamp,
		repgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		repgen.logger.Error("error signing timestamp", commontypes.LogFields{
			"error": err,
		})
		return
	}

	repgen.logger.Info("Sending MessageReconcile to leader", commontypes.LogFields{
		"epoch":                     ev.Epoch,
		"leader":                    repgen.l,
		"highestCertifiedTimestamp": highestCertifiedTimestamp,
	})
	repgen.netSender.SendTo(MessageReconcile[RI]{
		repgen.e,
		highestCertified,
		signedHighestCertifiedTimestamp,
	}, repgen.l)

	if repgen.id == repgen.l {
		repgen.leaderState.tRound = time.After(repgen.config.DeltaRound)
	}

	repgen.unbufferMessages()
}

func (repgen *reportGenerationState[RI]) Timestamp() Timestamp {
	return Timestamp{repgen.config.ConfigDigest, repgen.e}
}

func (repgen *reportGenerationState[RI]) OutcomeCtx(seqNr uint64) ocr3types.OutcomeContext {
	if seqNr != repgen.followerState.deliveredSeqNr+1 {
		repgen.logger.Critical("Assumption violation, seqNr isn't successor to deliveredSeqNr", commontypes.LogFields{
			"seqNr":          seqNr,
			"deliveredSeqNr": repgen.followerState.deliveredSeqNr,
		})
		panic("")
	}
	return ocr3types.OutcomeContext{
		seqNr,
		repgen.followerState.deliveredOutcome,
		uint64(repgen.e),
		seqNr - repgen.followerState.firstSeqNrOfEpoch + 1,
	}
}

func callPlugin[T any, RI any](
	repgen *reportGenerationState[RI],
	name string,
	maxDuration time.Duration,
	outctx ocr3types.OutcomeContext,
	f func(context.Context, ocr3types.OutcomeContext) (T, error),
) (T, bool) {
	ctx, cancel := context.WithTimeout(repgen.ctx, maxDuration)
	defer cancel()

	repgen.logger.Debug(fmt.Sprintf("calling ReportingPlugin.%s", name), commontypes.LogFields{
		"seqNr":       outctx.SeqNr,
		"round":       outctx.Round,
		"maxDuration": maxDuration,
	})

	// copy to avoid races when used inside the following closure
	logger := repgen.logger

	ins := loghelper.NewIfNotStopped(
		maxDuration+ReportingPluginTimeoutWarningGracePeriod,
		func() {
			logger.Error(fmt.Sprintf("call to ReportingPlugin.%s is taking too long", name), commontypes.LogFields{
				"seqNr":       outctx.SeqNr,
				"maxDuration": maxDuration,
			})
		},
	)

	result, err := f(ctx, outctx)

	ins.Stop()

	if err != nil {
		repgen.logger.ErrorIfNotCanceled(fmt.Sprintf("call to ReportingPlugin.%s errored", name), repgen.ctx, commontypes.LogFields{
			"seqNr": outctx.SeqNr,
			"error": err,
		})
		// failed to get data, nothing to be done
		var zero T
		return zero, false
	}

	return result, true
}
