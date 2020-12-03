package protocol

import (
	"context"
	"math/big"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

func (repgen *reportGenerationState) followerReportContext() ReportContext {
	return ReportContext{repgen.config.ConfigDigest, repgen.e, repgen.followerState.r}
}











func (repgen *reportGenerationState) messageObserveReq(msg MessageObserveReq, sender types.OracleID) {
	dropPrefix := "messageObserveReq: dropping MessageObserveReq from "
	
	
	if msg.Epoch != repgen.e {
		repgen.logger.Debug(dropPrefix+"wrong epoch",
			types.LogFields{"round": repgen.followerState.r, "msgEpoch": msg.Epoch},
		)
		return
	}
	if sender != repgen.l {
		
		repgen.logger.Warn(dropPrefix+"non-leader",
			types.LogFields{"round": repgen.followerState.r, "sender": sender})
		return
	}
	if msg.Round <= repgen.followerState.r {
		
		repgen.logger.Debug(dropPrefix+"earlier round",
			types.LogFields{"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if int64(repgen.config.RMax)+1 < int64(msg.Round) {
		
		
		
		
		
		
		
		repgen.logger.Warn(dropPrefix+"out of bounds round",
			types.LogFields{"round": repgen.followerState.r, "rMax": repgen.config.RMax, "msgRound": msg.Round})
		return
	}

	
	
	
	repgen.followerState.r = msg.Round

	if repgen.followerState.r > repgen.config.RMax {
		repgen.logger.Debug(
			"messageReportReq: leader sent MessageObserveReq past its expiration "+
				"round. Time to change leader",
			types.LogFields{
				"round":        repgen.followerState.r,
				"messageRound": msg.Round,
				"roundMax":     repgen.config.RMax,
			})
		select {
		case repgen.chReportGenerationToPacemaker <- EventChangeLeader{}:
		case <-repgen.ctx.Done():
		}

		
		
		
		
		return
	}
	
	
	
	
	
	
	
	repgen.followerState.sentEcho = nil
	repgen.followerState.sentReport = false
	repgen.followerState.completedRound = false
	repgen.followerState.receivedEcho = make([]bool, repgen.config.N())

	value := repgen.observeValue()
	if value.IsMissingValue() {
		
		
		return
	}

	so, err := MakeSignedObservation(value, repgen.followerReportContext(), repgen.privateKeys.SignOffChain)
	if err != nil {
		repgen.logger.Error("messageObserveReq: could not make SignedObservation observation", types.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return
	}

	if err := so.Verify(repgen.followerReportContext(), repgen.privateKeys.PublicKeyOffChain()); err != nil {
		repgen.logger.Error("MakeSignedObservation produced invalid signature:", types.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return
	}

	repgen.logger.Debug("sent observation to leader", types.LogFields{
		"round":       repgen.followerState.r,
		"observation": value,
	})
	repgen.netSender.SendTo(MessageObserve{
		repgen.e,
		repgen.followerState.r,
		so,
	}, repgen.l)
}




func (repgen *reportGenerationState) messageReportReq(msg MessageReportReq, sender types.OracleID) {
	
	
	if repgen.e != msg.Epoch {
		repgen.logger.Debug("messageReportReq from wrong epoch", types.LogFields{
			"round":    repgen.followerState.r,
			"msgEpoch": msg.Epoch})
		return
	}
	if sender != repgen.l {
		
		repgen.logger.Warn("messageReportReq from non-leader", types.LogFields{
			"round": repgen.followerState.r, "sender": sender})
		return
	}
	if repgen.followerState.r != msg.Round {
		
		
		repgen.logger.Debug("messageReportReq from wrong round", types.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if repgen.followerState.sentReport {
		repgen.logger.Warn("messageReportReq after report sent", types.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if repgen.followerState.completedRound {
		repgen.logger.Warn("messageReportReq after round completed", types.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	err := repgen.verifyReportReq(msg)
	if err != nil {
		repgen.logger.Error("messageReportReq: could not validate report sent by leader", types.LogFields{
			"round": repgen.followerState.r,
			"error": err,
			"msg":   msg,
		})
		return
	}

	if repgen.shouldReport(msg.AttributedSignedObservations) {
		attributedValues := make([]AttributedObservation, len(msg.AttributedSignedObservations))
		for i, aso := range msg.AttributedSignedObservations {
			
			attributedValues[i] = AttributedObservation{
				aso.SignedObservation.Observation,
				aso.Observer,
			}
		}

		report, err := MakeAttestedReportOne(
			attributedValues,
			repgen.followerReportContext(),
			repgen.privateKeys.SignOnChain,
		)
		if err != nil {
			
			
			repgen.logger.Error("messageReportReq: failed to sign report", types.LogFields{
				"round":  repgen.followerState.r,
				"error":  err,
				"id":     repgen.id,
				"report": report,
				"pubkey": repgen.privateKeys.PublicKeyAddressOnChain(),
			})
			return
		}

		{
			err := report.Verify(repgen.followerReportContext(), repgen.privateKeys.PublicKeyAddressOnChain())
			if err != nil {
				repgen.logger.Error("could not verify my own signature", types.LogFields{
					"round":  repgen.followerState.r,
					"error":  err,
					"id":     repgen.id,
					"report": report, 
					"pubkey": repgen.privateKeys.PublicKeyAddressOnChain()})
				return
			}
		}

		repgen.followerState.sentReport = true
		repgen.netSender.SendTo(
			MessageReport{
				repgen.e,
				repgen.followerState.r,
				report,
			},
			repgen.l,
		)
	} else {
		repgen.completeRound()
	}
}




func (repgen *reportGenerationState) messageFinal(
	msg MessageFinal, sender types.OracleID,
) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("wrong epoch from MessageFinal", types.LogFields{
			"round": repgen.followerState.r, "msgEpoch": msg.Epoch, "sender": sender})
		return
	}
	if msg.Round != repgen.followerState.r {
		repgen.logger.Debug("wrong round from MessageFinal", types.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if sender != repgen.l {
		repgen.logger.Warn("MessageFinal from non-leader", types.LogFields{
			"msgEpoch": msg.Epoch, "sender": sender,
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if repgen.followerState.sentEcho != nil {
		repgen.logger.Debug("MessageFinal after already sent MessageFinalEcho", nil)
		return
	}
	if !repgen.verifyAttestedReport(msg.Report, sender) {
		return
	}
	repgen.followerState.sentEcho = &msg.Report
	repgen.netSender.Broadcast(MessageFinalEcho{MessageFinal: msg})
}







func (repgen *reportGenerationState) messageFinalEcho(msg MessageFinalEcho,
	sender types.OracleID,
) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("wrong epoch from MessageFinalEcho", types.LogFields{
			"round": repgen.followerState.r, "msgEpoch": msg.Epoch, "sender": sender})
		return
	}
	if msg.Round != repgen.followerState.r {
		repgen.logger.Debug("wrong round from MessageFinalEcho", types.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round, "sender": sender})
		return
	}
	if repgen.followerState.receivedEcho[sender] {
		repgen.logger.Warn("extra MessageFinalEcho received", types.LogFields{
			"round": repgen.followerState.r, "sender": sender})
		return
	}
	if repgen.followerState.completedRound {
		repgen.logger.Debug("received final echo after round completion", nil)
		return
	}
	if !repgen.verifyAttestedReport(msg.Report, sender) { 
		
		return
	}
	repgen.followerState.receivedEcho[sender] = true 

	if repgen.followerState.sentEcho == nil { 
		repgen.followerState.sentEcho = &msg.Report 
		repgen.netSender.Broadcast(msg)             
	}

	
	{
		count := 0 
		for _, receivedEcho := range repgen.followerState.receivedEcho {
			if receivedEcho {
				count++
			}
		}
		if repgen.config.F < count {
			select {
			case repgen.chReportGenerationToTransmission <- EventTransmit{
				repgen.e,
				repgen.followerState.r,
				*repgen.followerState.sentEcho,
			}:
			case <-repgen.ctx.Done():
			}
			repgen.completeRound()
		}
	}

}



func (repgen *reportGenerationState) observeValue() observation.Observation {
	var value observation.Observation
	var err error
	
	
	
	
	ok := repgen.subprocesses.BlockForAtMost(
		repgen.ctx,
		repgen.localConfig.DataSourceTimeout,
		func(ctx context.Context) {
			var rawValue types.Observation
			rawValue, err = repgen.datasource.Observe(ctx)
			if err != nil {
				return
			}
			value, err = observation.MakeObservation((*big.Int)(rawValue))
		},
	)

	if !ok {
		repgen.logger.Error("DataSource timed out", types.LogFields{
			"round":   repgen.followerState.r,
			"timeout": repgen.localConfig.DataSourceTimeout,
		})
		return observation.Observation{}
	}

	if err != nil {
		repgen.logger.Error("DataSource errored", types.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return observation.Observation{}
	}

	return value
}

func (repgen *reportGenerationState) shouldReport(observations []AttributedSignedObservation) bool {
	ctx, cancel := context.WithTimeout(repgen.ctx, repgen.localConfig.BlockchainTimeout)
	defer cancel()
	contractConfigDigest, contractEpoch, contractRound, rawAnswer, timestamp,
		err := repgen.contractTransmitter.LatestTransmissionDetails(ctx)
	if err != nil {
		repgen.logger.Error("shouldReport: Error during LatestTransmissionDetails", types.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		
		
		
		return true
	}

	answer, err := observation.MakeObservation(rawAnswer)
	if err != nil {
		repgen.logger.Error("shouldReport: Error during observation.NewObservation", types.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return false
	}

	initialRound := contractConfigDigest == repgen.config.ConfigDigest && contractEpoch == 0 && contractRound == 0
	deviation := observations[len(observations)/2].SignedObservation.Observation.Deviates(answer, repgen.config.AlphaPPB)
	deltaCTimeout := timestamp.Add(repgen.config.DeltaC).Before(time.Now())
	result := initialRound || deviation || deltaCTimeout

	repgen.logger.Info("shouldReport: returning result", types.LogFields{
		"round":         repgen.followerState.r,
		"result":        result,
		"initialRound":  initialRound,
		"deviation":     deviation,
		"deltaCTimeout": deltaCTimeout,
	})

	return result
}





func (repgen *reportGenerationState) completeRound() {
	repgen.logger.Debug("ReportGeneration: completed round", types.LogFields{
		"round": repgen.followerState.r,
	})
	repgen.followerState.completedRound = true

	select {
	case repgen.chReportGenerationToPacemaker <- EventProgress{}:
	case <-repgen.ctx.Done():
	}
}




func (repgen *reportGenerationState) verifyReportReq(msg MessageReportReq) error {
	
	if !sort.SliceIsSorted(msg.AttributedSignedObservations,
		func(i, j int) bool {
			return msg.AttributedSignedObservations[i].SignedObservation.Observation.Less(msg.AttributedSignedObservations[j].SignedObservation.Observation)
		}) {
		return errors.Errorf("messages not sorted by value")
	}

	
	{
		counted := map[types.OracleID]bool{}
		for _, obs := range msg.AttributedSignedObservations {
			
			numOracles := len(repgen.config.OracleIdentities)
			if int(obs.Observer) < 0 || numOracles <= int(obs.Observer) {
				return errors.Errorf("given oracle ID of %v is out of bounds (only "+
					"have %v public keys)", obs.Observer, numOracles)
			}
			if counted[obs.Observer] {
				return errors.Errorf("duplicate observation by oracle id %v", obs.Observer)
			} else {
				counted[obs.Observer] = true
			}
			observerOffchainPublicKey := repgen.config.OracleIdentities[obs.Observer].OffchainPublicKey
			if err := obs.SignedObservation.Verify(repgen.followerReportContext(), observerOffchainPublicKey); err != nil {
				return errors.Errorf("invalid signed observation: %s", err)
			}
		}
		bound := 2 * repgen.config.F
		if len(counted) <= bound {
			return errors.Errorf("not enough observations in report; got %d, "+
				"need more than %d", len(counted), bound)
		}
	}
	return nil
}



func (repgen *reportGenerationState) verifyAttestedReport(
	report AttestedReportMany, sender types.OracleID,
) bool {
	if len(report.Signatures) <= repgen.config.F {
		repgen.logger.Warn("verifyAttestedReport: dropping final report because "+
			"it has too few signatures", types.LogFields{"sender": sender,
			"numSignatures": len(report.Signatures), "F": repgen.config.F})
		return false
	}

	keys := make(signature.EthAddresses)
	for oid, id := range repgen.config.OracleIdentities {
		keys[types.OnChainSigningAddress(id.OnChainSigningAddress)] =
			types.OracleID(oid)
	}

	err := report.VerifySignatures(repgen.followerReportContext(), keys)
	if err != nil {
		repgen.logger.Error("could not validate signatures on final report",
			types.LogFields{
				"round":  repgen.followerState.r,
				"error":  err,
				"report": report,
				"sender": sender,
			})
		return false
	}
	return true
}
