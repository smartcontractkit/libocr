package protocol

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// RunOracle runs one oracle instance of the offchain reporting protocol and manages
// the lifecycle of all underlying goroutines.
//
// RunOracle runs forever until ctx is cancelled. It will only shut down
// after all its sub-goroutines have exited.
func RunOracle(
	ctx context.Context,

	config config.SharedConfig,
	contractTransmitter types.ContractTransmitter,
	database types.Database,
	datasource types.DataSource,
	id types.OracleID,
	keys types.PrivateKeys,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	netEndpoint NetworkEndpoint,
	telemetrySender TelemetrySender,
) {
	o := oracleState{
		ctx: ctx,

		Config:              config,
		contractTransmitter: contractTransmitter,
		database:            database,
		datasource:          datasource,
		id:                  id,
		localConfig:         localConfig,
		logger:              logger,
		netEndpoint:         netEndpoint,
		PrivateKeys:         keys,
		telemetrySender:     telemetrySender,
	}
	o.run()
}

type oracleState struct {
	ctx context.Context

	Config              config.SharedConfig
	contractTransmitter types.ContractTransmitter
	database            types.Database
	datasource          types.DataSource
	id                  types.OracleID
	localConfig         types.LocalConfig
	logger              loghelper.LoggerWithContext
	netEndpoint         NetworkEndpoint
	PrivateKeys         types.PrivateKeys
	telemetrySender     TelemetrySender

	chNetToPacemaker        chan<- MessageToPacemakerWithSender
	chNetToReportGeneration chan<- MessageToReportGenerationWithSender
	childCancel             context.CancelFunc
	childCtx                context.Context
	subprocesses            subprocesses.Subprocesses
}

// run ensures safe shutdown of the Oracle's "child routines",
// (Pacemaker, ReportGeneration and Transmission) upon o.ctx.Done()
// being closed.
//
// Here is a graph of the various channels involved:
//
//   Oracle +-----------------------------> Pacemaker
//        +                                     ^
//        |                                     |
//        +---------------------------------+   |
//                                          |   |
//                                          v   +
//   Transmission <-----------------------+ ReportGeneration
//
// All channels are unbuffered.
//
// Once o.ctx.Done() is closed, the Oracle runloop will enter the
// corresponding select case and no longer forward network messages
// to Pacemaker and ReportGeneration. It will then cancel o.childCtx,
// making all children exit. To prevent a deadlock on send in ReportGeneration
// all channel sends in ReportGeneration are contained in a select statements
// that also checks for cancellation. Similarly, all channel receives in
// Pacemaker and Transmission are contained in select statements that also
// check for cancellation.
//
// Finally, all sub-goroutines spawned in the protocol are attached to o.subprocesses
// (with the exception of ReportGeneration which is explicitly managed by Pacemaker).
// This enables us to wait for their completion before exiting.
func (o *oracleState) run() {
	o.logger.Info("Running", nil)

	chNetToPacemaker := make(chan MessageToPacemakerWithSender)
	o.chNetToPacemaker = chNetToPacemaker

	chNetToReportGeneration := make(chan MessageToReportGenerationWithSender)
	o.chNetToReportGeneration = chNetToReportGeneration

	chReportGenerationToTransmission := make(chan EventToTransmission)

	o.childCtx, o.childCancel = context.WithCancel(context.Background())
	defer o.childCancel()

	o.subprocesses.Go(func() {
		RunPacemaker(
			o.childCtx,
			&o.subprocesses,

			chNetToPacemaker,
			chNetToReportGeneration,
			chReportGenerationToTransmission,
			o.Config,
			o.contractTransmitter,
			o.database,
			o.datasource,
			o.id,
			o.localConfig,
			o.logger,
			o.netEndpoint,
			o.PrivateKeys,
			o.telemetrySender,
		)
	})
	o.subprocesses.Go(func() {
		RunTransmission(
			o.childCtx,
			&o.subprocesses,

			o.Config,
			chReportGenerationToTransmission,
			o.database,
			o.id,
			o.localConfig,
			o.logger,
			o.contractTransmitter,
		)
	})

	chNet := o.netEndpoint.Receive()

	chDone := o.ctx.Done()
	for {
		select {
		case msg := <-chNet:
			msg.Msg.process(o, msg.Sender)
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
