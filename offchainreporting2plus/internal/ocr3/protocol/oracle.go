package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
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

	config ocr3config.SharedConfig,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	database Database,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	netEndpoint NetworkEndpoint[RI],
	offchainKeyring types.OffchainKeyring,
	onchainKeyring ocr3types.OnchainKeyring[RI],
	reportingPlugin ocr3types.ReportingPlugin[RI],
	telemetrySender TelemetrySender,
) {
	o := oracleState[RI]{
		ctx: ctx,

		config:              config,
		contractTransmitter: contractTransmitter,
		database:            database,
		id:                  id,
		localConfig:         localConfig,
		logger:              logger,
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

	config              ocr3config.SharedConfig
	contractTransmitter ocr3types.ContractTransmitter[RI]
	database            Database
	id                  commontypes.OracleID
	localConfig         types.LocalConfig
	logger              loghelper.LoggerWithContext
	netEndpoint         NetworkEndpoint[RI]
	offchainKeyring     types.OffchainKeyring
	onchainKeyring      ocr3types.OnchainKeyring[RI]
	reportingPlugin     ocr3types.ReportingPlugin[RI]
	telemetrySender     TelemetrySender

	chNetToPacemaker         chan<- MessageToPacemakerWithSender[RI]
	chNetToOutcomeGeneration chan<- MessageToOutcomeGenerationWithSender[RI]
	chNetToReportAttestation chan<- MessageToReportAttestationWithSender[RI]
	childCancel              context.CancelFunc
	childCtx                 context.Context
	epoch                    uint64
	subprocesses             subprocesses.Subprocesses
}

// TODO: This comment is outdated
// run ensures safe shutdown of the Oracle's "child routines",
// (Pacemaker, OutcomeGeneration and Transmission) upon o.ctx.Done()
// being closed.
//
// Here is a graph of the various channels involved and what they
// transport.
//
//	    ┌────────────epoch changes──────────────┐
//	    ▼                                       │
//	┌──────┐                               ┌────┴────┐
//	│Oracle├─────pacemaker messages───────►│Pacemaker│
//	└──┬─┬─┘                               └─────────┘
//	   │ │                                       ▲
//	   │ └───────rep. gen. messages───────────┐  │
//	   │rep. fin. messages                    │  │
//	   ▼                                      ▼  │progress events
//	┌──────────────────┐                   ┌─────┴──────────┐
//	│ReportAttestation│◄───final events───┤OutcomeGeneration│
//	└────────┬─────────┘                   └────────────────┘
//	         │
//	         │transmit events
//	         ▼
//	    ┌────────────┐
//	    │Transmission│
//	    └────────────┘
//
// All channels are unbuffered.
//
// Once o.ctx.Done() is closed, the Oracle runloop will enter the
// corresponding select case and no longer forward network messages
// to Pacemaker and OutcomeGeneration. It will then cancel o.childCtx,
// making all children exit. To prevent deadlocks, all channel sends and
// receives in Oracle, Pacemaker, OutcomeGeneration, Transmission, etc...
// are contained in select{} statements that also contain a case for context
// cancellation.
//
// Finally, all sub-goroutines spawned in the protocol are attached to o.subprocesses
// (with the exception of OutcomeGeneration which is explicitly managed by Pacemaker).
// This enables us to wait for their completion before exiting.
func (o *oracleState[RI]) run() {
	o.logger.Info("Running", nil)

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

	o.childCtx, o.childCancel = context.WithCancel(context.Background())
	defer o.childCancel()

	cert, err := o.restoreCertFromDatabase()
	if err != nil {
		o.logger.Info("restoreCertFromDatabase returned an error, exiting oracle", commontypes.LogFields{
			"error": err,
		})
		return
	}

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
			o.netEndpoint,
			o.offchainKeyring,
			o.telemetrySender,

			cert.Epoch(),
		)
	})
	o.subprocesses.Go(func() {
		RunOutcomeGeneration[RI](
			o.childCtx,
			&o.subprocesses,

			chNetToOutcomeGeneration,
			chPacemakerToOutcomeGeneration,
			chOutcomeGenerationToPacemaker,
			chOutcomeGenerationToReportAttestation,
			o.config,
			o.database,
			o.id,
			o.localConfig,
			o.logger,
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
			chReportAttestationToTransmission,
			chOutcomeGenerationToReportAttestation,
			o.config,
			o.contractTransmitter,
			o.logger,
			o.netEndpoint,
			o.onchainKeyring,
			o.reportingPlugin,
		)
	})
	o.subprocesses.Go(func() {
		RunTransmission(
			o.childCtx,
			&o.subprocesses,

			chReportAttestationToTransmission,
			o.config,
			o.contractTransmitter,
			o.id,
			o.localConfig,
			o.logger,
			o.reportingPlugin,
		)
	})

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

func (o *oracleState[RI]) restoreCertFromDatabase() (CertifiedPrepareOrCommit, error) {
	const retryWait = 5 * time.Second

	for {
		ctx, cancel := context.WithTimeout(o.ctx, o.localConfig.DatabaseTimeout)
		defer cancel()
		cert, err := o.database.ReadCert(ctx, o.config.ConfigDigest)
		if err == nil {
			if cert != nil {
				o.logger.Info("restoreCertFromDatabase: successfully restored cert", commontypes.LogFields{
					"cert": cert,
				})
				return cert, nil
			} else {
				o.logger.Info("restoreCertFromDatabase: did not find cert, starting at genesis", nil)
				return &CertifiedCommit{}, nil
			}
		}

		o.logger.Error("restoreCertFromDatabase: database read failed, retrying", commontypes.LogFields{
			"error":     err,
			"retryWait": retryWait,
		})

		select {
		case <-time.After(retryWait):
		case <-o.ctx.Done():
			return nil, o.ctx.Err()
		}
	}
}
