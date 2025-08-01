package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// RunOracle runs one oracle instance of the offchain reporting protocol and manages
// the lifecycle of all underlying goroutines.
//
// RunOracle runs forever until ctx is cancelled. It will only shut down
// after all its sub-goroutines have exited.
func RunOracle[RI any](
	ctx context.Context,

	blobEndpointWrapper *BlobEndpointWrapper,
	config ocr3config.SharedConfig,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	database Database,
	id commontypes.OracleID,
	kvStore KeyValueStore,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netEndpoint NetworkEndpoint[RI],
	offchainKeyring types.OffchainKeyring,
	onchainKeyring ocr3types.OnchainKeyring[RI],
	reportingPlugin ocr3_1types.ReportingPlugin[RI],
	telemetrySender TelemetrySender,
) {
	o := oracleState[RI]{
		ctx: ctx,

		blobEndpointWrapper: blobEndpointWrapper,
		config:              config,
		contractTransmitter: contractTransmitter,
		database:            database,
		id:                  id,
		kvStore:             kvStore,
		localConfig:         localConfig,
		logger:              logger,
		metricsRegisterer:   metricsRegisterer,
		netEndpoint:         netEndpoint,
		offchainKeyring:     offchainKeyring,
		onchainKeyring:      onchainKeyring,
		reportingPlugin:     reportingPlugin,
		telemetrySender:     telemetrySender,
	}
	o.run()
}

type oracleState[RI any] struct {
	ctx context.Context

	blobEndpointWrapper *BlobEndpointWrapper
	config              ocr3config.SharedConfig
	contractTransmitter ocr3types.ContractTransmitter[RI]
	database            Database
	id                  commontypes.OracleID
	kvStore             KeyValueStore
	localConfig         types.LocalConfig
	logger              loghelper.LoggerWithContext
	metricsRegisterer   prometheus.Registerer
	netEndpoint         NetworkEndpoint[RI]
	offchainKeyring     types.OffchainKeyring
	onchainKeyring      ocr3types.OnchainKeyring[RI]
	reportingPlugin     ocr3_1types.ReportingPlugin[RI]
	telemetrySender     TelemetrySender

	chNetToPacemaker         chan<- MessageToPacemakerWithSender[RI]
	chNetToOutcomeGeneration chan<- MessageToOutcomeGenerationWithSender[RI]
	chNetToReportAttestation chan<- MessageToReportAttestationWithSender[RI]
	chNetToStatePersistence  chan<- MessageToStatePersistenceWithSender[RI]
	chNetToBlobExchange      chan<- MessageToBlobExchangeWithSender[RI]
	childCancel              context.CancelFunc
	childCtx                 context.Context
	epoch                    uint64
	subprocesses             subprocesses.Subprocesses
}

// run ensures safe shutdown of the Oracle's "child routines",
// (Pacemaker, OutcomeGeneration, Attestation, State, Transmission, and BlobExchange) upon
// o.ctx.Done()
//
// Here is a graph of the various channels involved and what they
// transport.
//
//	                        state
//	                        message
//	    ┌───────────────────────────────────────────────────────┐
//		   │                  ┌─────────┐                          │
//		   ├─────────────────►│Pacemaker│                          │
//		   │     pacemaker    └──────┬──┘                          │
//		   │     message         ▲   │                             │
//		   │                     │   │                             │
//		   │             progress│   │epoch                        │
//		   │              /change│   │start                        │
//		   │                epoch│   │notification                 │
//		   │              request│   │                             │
//		   ▼                     │   ▼                             ▼
//		┌──────┐              ┌──┴───────────────┐              ┌─────────────────┐
//		│Oracle│◄────────────►│Outcome Generation│◄────────────►│State Persistence│
//		└──────┘  out.gen.    └──────┬───────────┘              └─────────────────┘
//		   ▲      message            │                             ▲
//		   │                         │certified                    │
//		   │                         │outcome                      │
//		   │                         │                             │
//		   │                         ▼                             │
//		   │                  ┌────────────┐                       │
//		   └─────────────────►│ Attestation│◄──────────────────────┘
//		          rep.att.    └──────┬─────┘
//		          message            │
//		                             │attested
//		                             │report
//		                             │
//		                             ▼
//		                      ┌────────────┐
//		                      │Transmission│
//		                      └────────────┘
//
// All channels are unbuffered.
//
// Once o.ctx.Done() is closed, the Oracle runloop will enter the corresponding
// select case and no longer forward network messages to Pacemaker,
// OutcomeGeneration, etc... It will then cancel o.childCtx, making all children
// exit. To prevent deadlocks, all channel sends and receives in Oracle,
// Pacemaker, OutcomeGeneration, etc... are (1) contained in select{} statements
// that also contain a case for context cancellation or (2) guaranteed to occur
// before o.childCtx is cancelled.
//
// Finally, all sub-goroutines spawned in the protocol are attached to o.subprocesses
// This enables us to wait for their completion before exiting.
func (o *oracleState[RI]) run() {
	o.logger.Info("Oracle: running", commontypes.LogFields{
		"localConfig":  fmt.Sprintf("%+v", o.localConfig),
		"publicConfig": fmt.Sprintf("%+v", o.config.PublicConfig),
	})

	chNetToPacemaker := make(chan MessageToPacemakerWithSender[RI])
	o.chNetToPacemaker = chNetToPacemaker

	chNetToOutcomeGeneration := make(chan MessageToOutcomeGenerationWithSender[RI])
	o.chNetToOutcomeGeneration = chNetToOutcomeGeneration

	chPacemakerToOutcomeGeneration := make(chan EventToOutcomeGeneration[RI])

	chOutcomeGenerationToPacemaker := make(chan EventToPacemaker[RI])

	chNetToReportAttestation := make(chan MessageToReportAttestationWithSender[RI])
	o.chNetToReportAttestation = chNetToReportAttestation

	chOutcomeGenerationToReportAttestation := make(chan EventToReportAttestation[RI])

	chReportAttestationToTransmission := make(chan EventToTransmission[RI])

	chNetToStatePersistence := make(chan MessageToStatePersistenceWithSender[RI])
	o.chNetToStatePersistence = chNetToStatePersistence

	chReportAttestationToStatePersistence := make(chan EventToStatePersistence[RI])

	chNetToBlobExchange := make(chan MessageToBlobExchangeWithSender[RI])
	o.chNetToBlobExchange = chNetToBlobExchange

	chOutcomeGenerationToBlobExchange := make(chan EventToBlobExchange[RI])

	// communication between blob exchange and blob endpoint
	chBlobBroadcastRequest := make(chan blobBroadcastRequest)
	chBlobBroadcastResponse := make(chan blobBroadcastResponse)

	chBlobFetchRequest := make(chan blobFetchRequest)
	chBlobFetchResponse := make(chan blobFetchResponse)

	// be careful if you want to change anything here.
	// chNetTo* sends in message.go assume that their recipients are running.
	o.childCtx, o.childCancel = context.WithCancel(context.Background())
	defer o.childCancel()

	defer o.kvStore.Close()

	paceState, cert, statePersistenceState, err := o.restoreFromDatabase()
	if err != nil {
		o.logger.Error("restoreFromDatabase returned an error, exiting oracle", commontypes.LogFields{
			"error": err,
		})
		return
	}
	highestCommittedToKVdSeqNr, err := o.kvStore.HighestCommittedSeqNr()
	if err != nil {
		o.logger.Error("cannot read highest committed seqNr from key value store, exiting oracle",
			commontypes.LogFields{
				"error": err,
			})
		return
	}

	blobEndpoint := BlobEndpoint{
		o.childCtx,

		chBlobBroadcastRequest,
		chBlobBroadcastResponse,

		chBlobFetchRequest,
		chBlobFetchResponse,
	}
	o.blobEndpointWrapper.setBlobEndpoint(&blobEndpoint) // pass through to plugin

	o.subprocesses.Go(func() {
		RunPacemaker[RI](
			o.childCtx,

			chNetToPacemaker,
			chPacemakerToOutcomeGeneration,
			chOutcomeGenerationToPacemaker,
			o.config,
			o.database,
			o.id,
			o.localConfig,
			o.logger,
			o.metricsRegisterer,
			o.netEndpoint,
			o.offchainKeyring,
			o.telemetrySender,

			paceState,
		)
	})
	o.subprocesses.Go(func() {
		RunOutcomeGeneration[RI](
			o.childCtx,

			chNetToOutcomeGeneration,
			chPacemakerToOutcomeGeneration,
			chOutcomeGenerationToPacemaker,
			chOutcomeGenerationToReportAttestation,
			&blobEndpoint,
			o.config,
			o.database,
			o.id,
			o.kvStore,
			o.localConfig,
			o.logger,
			o.metricsRegisterer,
			o.netEndpoint,
			o.offchainKeyring,
			o.reportingPlugin,
			o.telemetrySender,

			cert,
		)
	})

	o.subprocesses.Go(func() {
		RunReportAttestation[RI](
			o.childCtx,

			chNetToReportAttestation,
			chOutcomeGenerationToReportAttestation,
			chReportAttestationToStatePersistence,
			chReportAttestationToTransmission,
			o.config,
			o.contractTransmitter,
			o.logger,
			o.netEndpoint,
			o.onchainKeyring,
			o.reportingPlugin,
		)
	})

	o.subprocesses.Go(func() {
		RunStatePersistence[RI](
			o.childCtx,

			chNetToStatePersistence,
			chReportAttestationToStatePersistence,
			o.config,
			o.database,
			o.id,
			o.kvStore,
			o.logger,
			o.netEndpoint,
			o.reportingPlugin,
			statePersistenceState,
			highestCommittedToKVdSeqNr,
		)
	})

	o.subprocesses.Go(func() {
		RunTransmission(
			o.childCtx,

			chReportAttestationToTransmission,
			o.config,
			o.contractTransmitter,
			o.id,
			o.localConfig,
			o.logger,
			o.reportingPlugin,
		)
	})

	o.subprocesses.Go(func() {
		RunBlobExchange[RI](
			o.childCtx,

			chNetToBlobExchange,
			chOutcomeGenerationToBlobExchange,

			chBlobBroadcastRequest,
			chBlobBroadcastResponse,

			chBlobFetchRequest,
			chBlobFetchResponse,

			o.config,
			o.kvStore,
			o.id,
			o.localConfig,
			o.logger,
			o.metricsRegisterer,
			o.netEndpoint,
			o.offchainKeyring,
			o.telemetrySender,
		)
	})

	publicConfigMetrics := ocr3config.NewPublicConfigMetrics(o.metricsRegisterer, o.logger, o.config.PublicConfig)
	defer publicConfigMetrics.Close()

	chNet := o.netEndpoint.Receive()

	chDone := o.ctx.Done()
	for {
		select {
		case msg := <-chNet:
			// This bounds check should never trigger since it's the netEndpoint's
			// responsibility to only provide valid senders. We perform it for
			// defense-in-depth.
			if 0 <= int(msg.Sender) && int(msg.Sender) < o.config.N() {
				msg.Msg.process(o, msg.Sender)
			} else {
				o.logger.Critical("msg.Sender out of bounds. This should *never* happen.", commontypes.LogFields{
					"sender": msg.Sender,
					"n":      o.config.N(),
				})
			}
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			o.logger.Debug("Oracle: winding down", nil)
			o.childCancel()
			o.subprocesses.Wait()
			o.logger.Debug("Oracle: exiting", nil)
			return
		default:
		}
	}
}

func tryUntilSuccess[T any](ctx context.Context, logger commontypes.Logger, retryPeriod time.Duration, fnTimeout time.Duration, fnName string, fn func(context.Context) (T, error)) (T, error) {
	for {
		var result T
		var err error
		func() {
			fnCtx, cancel := context.WithTimeout(ctx, fnTimeout)
			defer cancel()
			result, err = fn(fnCtx)
		}()
		if err == nil {
			return result, nil
		}
		logger.Error(fmt.Sprintf("error during %s, retrying", fnName), commontypes.LogFields{
			"error":       err,
			"retryPeriod": retryPeriod.String(),
		})

		select {
		case <-time.After(retryPeriod):
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		}
	}
}

func (o *oracleState[RI]) restoreFromDatabase() (PacemakerState, CertifiedPrepareOrCommit, StatePersistenceState, error) {
	const retryPeriod = 5 * time.Second

	paceState, err := tryUntilSuccess[PacemakerState](
		o.ctx,
		o.logger,
		retryPeriod,
		o.localConfig.DatabaseTimeout,
		"Database.ReadPacemakerState",
		func(ctx context.Context) (PacemakerState, error) {
			return o.database.ReadPacemakerState(ctx, o.config.ConfigDigest)
		},
	)
	if err != nil {
		return PacemakerState{}, nil, StatePersistenceState{}, err
	}

	o.logger.Info("restoreFromDatabase: successfully restored pacemaker state", commontypes.LogFields{
		"state": paceState,
	})

	cert, err := tryUntilSuccess[CertifiedPrepareOrCommit](
		o.ctx,
		o.logger,
		retryPeriod,
		o.localConfig.DatabaseTimeout,
		"Database.ReadCert",
		func(ctx context.Context) (CertifiedPrepareOrCommit, error) {
			return o.database.ReadCert(ctx, o.config.ConfigDigest)
		},
	)
	if err != nil {
		return PacemakerState{}, nil, StatePersistenceState{}, err
	}

	if cert != nil {
		o.logger.Info("restoreFromDatabase: successfully restored cert", commontypes.LogFields{
			"certTimestamp": cert.Timestamp(),
		})
	} else {
		o.logger.Info("restoreFromDatabase: did not find cert, starting at genesis", nil)
		cert = &CertifiedCommit{}
	}

	statePersistenceState, err := tryUntilSuccess[StatePersistenceState](
		o.ctx,
		o.logger,
		retryPeriod,
		o.localConfig.DatabaseTimeout,
		"Database.ReadStatePersistenceState",
		func(ctx context.Context) (StatePersistenceState, error) {
			return o.database.ReadStatePersistenceState(ctx, o.config.ConfigDigest)
		},
	)
	if err != nil {
		return PacemakerState{}, nil, StatePersistenceState{}, err
	}

	o.logger.Info("restoreFromDatabase: successfully restored state persistence state", commontypes.LogFields{
		"state": statePersistenceState,
	})

	return paceState, cert, statePersistenceState, nil
}
