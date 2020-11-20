

package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)


func RunReportGeneration(
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	chNetToReportGeneration <-chan MessageToReportGenerationWithSender,
	chReportGenerationToPacemaker chan<- EventToPacemaker,
	chReportGenerationToTransmission chan<- EventToTransmission,
	config config.SharedConfig,
	contractTransmitter types.ContractTransmitter,
	datasource types.DataSource,
	e uint32,
	id types.OracleID,
	l types.OracleID,
	localConfig types.LocalConfig,
	logger types.Logger,
	netSender NetworkSender,
	privateKeys types.PrivateKeys,
) {
	repgen := reportGenerationState{
		ctx:          ctx,
		subprocesses: subprocesses,

		chNetToReportGeneration:          chNetToReportGeneration,
		chReportGenerationToPacemaker:    chReportGenerationToPacemaker,
		chReportGenerationToTransmission: chReportGenerationToTransmission,
		config:                           config,
		contractTransmitter:              contractTransmitter,
		datasource:                       datasource,
		e:                                e,
		id:                               id,
		l:                                l,
		localConfig:                      localConfig,
		logger:                           loghelper.MakeLoggerWithContext(logger, types.LogFields{"epoch": e, "leader": l}),
		netSender:                        netSender,
		privateKeys:                      privateKeys,
	}
	repgen.run()
}

type reportGenerationState struct {
	ctx          context.Context
	subprocesses *subprocesses.Subprocesses

	chNetToReportGeneration          <-chan MessageToReportGenerationWithSender
	chReportGenerationToPacemaker    chan<- EventToPacemaker
	chReportGenerationToTransmission chan<- EventToTransmission
	config                           config.SharedConfig
	contractTransmitter              types.ContractTransmitter
	datasource                       types.DataSource
	e                                uint32 
	id                               types.OracleID
	l                                types.OracleID 
	localConfig                      types.LocalConfig
	logger                           types.Logger
	netSender                        NetworkSender
	privateKeys                      types.PrivateKeys

	leaderState   leaderState
	followerState followerState
}

type leaderState struct {
	
	r uint8

	
	observe []*SignedObservation

	
	report []*AttestedReportOne

	
	
	tRound <-chan time.Time

	
	
	
	tGrace <-chan time.Time

	phase phase
}

type followerState struct {
	
	r uint8

	
	
	receivedEcho []bool

	
	
	sentEcho *AttestedReportMany

	
	
	sentReport bool

	
	
	completedRound bool
}


func (repgen *reportGenerationState) run() {
	repgen.logger.Info("Running ReportGeneration", nil)

	
	repgen.leaderState.r = 0
	repgen.leaderState.report = make([]*AttestedReportOne, repgen.config.N())
	repgen.followerState.r = 0
	repgen.followerState.receivedEcho = make([]bool, repgen.config.N())
	repgen.followerState.sentEcho = nil
	repgen.followerState.completedRound = false

	
	if repgen.id == repgen.l {
		repgen.startRound()
	}

	
	chDone := repgen.ctx.Done()
	for {
		select {
		case msg := <-repgen.chNetToReportGeneration:
			msg.msg.processReportGeneration(repgen, msg.sender)
		case <-repgen.leaderState.tGrace:
			repgen.eventTGraceTimeout()
		case <-repgen.leaderState.tRound:
			repgen.eventTRoundTimeout()
		case <-chDone:
		}

		
		select {
		case <-chDone:
			repgen.logger.Info("ReportGeneration: exiting", types.LogFields{
				"e": repgen.e,
				"l": repgen.l,
			})
			return
		default:
		}
	}
}
