package protocol

import (
	"sort"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting/types"
)





type phase int

const (
	phaseObserve phase = iota
	phaseGrace
	phaseReport
	phaseFinal
)

var englishPhase = map[phase]string{
	phaseObserve: "observe",
	phaseGrace:   "grace",
	phaseReport:  "report",
	phaseFinal:   "final",
}

func (repgen *reportGenerationState) leaderReportContext() ReportContext {
	return ReportContext{repgen.config.ConfigDigest, repgen.e, repgen.leaderState.r}
}





func (repgen *reportGenerationState) eventTRoundTimeout() {
	repgen.startRound()
}







func (repgen *reportGenerationState) startRound() {
	rPlusOne := repgen.leaderState.r + 1
	if rPlusOne <= repgen.leaderState.r {
		repgen.logger.Error("ReportGeneration: round overflows, cannot start new round", nil)
		return
	}
	repgen.leaderState.r = rPlusOne
	repgen.leaderState.observe = make([]*SignedObservation, repgen.config.N())
	repgen.leaderState.phase = phaseObserve
	repgen.netSender.Broadcast(MessageObserveReq{Epoch: repgen.e, Round: repgen.leaderState.r})
	repgen.leaderState.tRound = time.After(repgen.config.DeltaRound)
}






func (repgen *reportGenerationState) messageObserve(msg MessageObserve, sender types.OracleID) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessageObserve for wrong epoch", types.LogFields{
			"epoch":    repgen.e,
			"round":    repgen.leaderState.r,
			"sender":   sender,
			"msgEpoch": msg.Epoch,
			"msgRound": msg.Round,
		})
		return
	}

	if repgen.l != repgen.id {
		repgen.logger.Warn("Non-leader received MessageObserve", types.LogFields{
			"sender": sender,
			"msg":    msg,
		})
		return
	}

	if msg.Round != repgen.leaderState.r {
		repgen.logger.Debug("Got MessageObserve for wrong round", types.LogFields{
			"epoch":    repgen.e,
			"round":    repgen.leaderState.r,
			"sender":   sender,
			"msgEpoch": msg.Epoch,
			"msgRound": msg.Round,
		})
		return
	}

	if repgen.leaderState.phase != phaseObserve && repgen.leaderState.phase != phaseGrace {
		repgen.logger.Debug("received MessageObserve after grace phase", nil)
		return
	}

	if repgen.leaderState.observe[sender] != nil {
		
		
		repgen.logger.Debug("already sent an observation", types.LogFields{
			"sender": sender})
		return
	}

	if err := msg.SignedObservation.Verify(repgen.leaderReportContext(), repgen.config.OracleIdentities[sender].OffchainPublicKey); err != nil {
		repgen.logger.Warn("MessageObserve carries invalid SignedObservation", types.LogFields{
			"round":  repgen.leaderState.r,
			"sender": sender,
			"msg":    msg,
			"error":  err,
		})
		return
	}

	repgen.logger.Debug("MessageObserve has valid SignedObservation", types.LogFields{
		"round":    repgen.leaderState.r,
		"sender":   sender,
		"msgEpoch": msg.Epoch,
		"msgRound": msg.Round,
	})

	repgen.leaderState.observe[sender] = &msg.SignedObservation

	
	switch repgen.leaderState.phase {
	case phaseObserve:
		observationCount := 0 
		for _, so := range repgen.leaderState.observe {
			if so != nil {
				observationCount++
			}
		}
		repgen.logger.Debug("One more observation", types.LogFields{
			"observationCount":         observationCount,
			"requiredObservationCount": (2 * repgen.config.F) + 1,
		})
		if observationCount > 2*repgen.config.F {
			
			repgen.logger.Debug("starting observation grace period", nil)
			repgen.leaderState.tGrace = time.After(repgen.config.DeltaGrace)
			repgen.leaderState.phase = phaseGrace
		}
	case phaseGrace:
		repgen.logger.Debug("accepted extra observation during grace period", nil)
	}
}




func (repgen *reportGenerationState) eventTGraceTimeout() {
	if repgen.leaderState.phase != phaseGrace {
		repgen.logger.Error("leader's phase conflicts tGrace timeout", types.LogFields{
			"phase": englishPhase[repgen.leaderState.phase],
		})
		return
	}
	asos := []AttributedSignedObservation{}
	for oid, so := range repgen.leaderState.observe {
		if so != nil {
			asos = append(asos, AttributedSignedObservation{
				*so,
				types.OracleID(oid),
			})
		}
	}
	sort.Slice(asos, func(i, j int) bool {
		return asos[i].SignedObservation.Observation.Less(asos[j].SignedObservation.Observation)
	})
	repgen.netSender.Broadcast(MessageReportReq{
		repgen.e,
		repgen.leaderState.r,
		asos,
	})
	repgen.leaderState.phase = phaseReport
}

func (repgen *reportGenerationState) messageReport(msg MessageReport, sender types.OracleID) {
	dropPrefix := "messageReport: dropping MessageReport due to "
	if msg.Epoch != repgen.e {
		repgen.logger.Debug(dropPrefix+"wrong epoch",
			types.LogFields{"epoch": repgen.e, "msgEpoch": msg.Epoch})
		return
	}
	if repgen.l != repgen.id {
		repgen.logger.Warn(dropPrefix+"not being leader of the current epoch",
			types.LogFields{"leader": repgen.l})
		return
	}
	if msg.Round != repgen.leaderState.r {
		repgen.logger.Debug(dropPrefix+"wrong round",
			types.LogFields{"round": repgen.leaderState.r, "msgRound": msg.Round})
		return
	}
	if repgen.leaderState.phase != phaseReport {
		repgen.logger.Debug(dropPrefix+"not being in report phase",
			types.LogFields{"currentPhase": englishPhase[repgen.leaderState.phase]})
		return
	}
	if repgen.leaderState.report[sender] != nil {
		repgen.logger.Warn(dropPrefix+"having already received sender's report",
			types.LogFields{"sender": sender, "msg": msg})
		return
	}

	a := types.OnChainSigningAddress(repgen.config.OracleIdentities[sender].OnChainSigningAddress)
	err := msg.Report.Verify(repgen.leaderReportContext(), a)
	if err != nil {
		repgen.logger.Error("could not validate signature", types.LogFields{
			"error": err,
			"msg":   msg,
		})
		return
	}

	repgen.leaderState.report[sender] = &msg.Report

	
	{ 
		sigs := [][]byte{}
		for _, report := range repgen.leaderState.report {
			if report == nil {
				continue
			}
			if report.AttributedObservations.Equal(msg.Report.AttributedObservations) {
				sigs = append(sigs, report.Signature)
			} else {
				repgen.logger.Warn("received disparate reports messages", types.LogFields{
					"previousReport": report,
					"msgReport":      msg,
				})
			}
		}

		if repgen.config.F < len(sigs) {
			repgen.netSender.Broadcast(MessageFinal{
				repgen.e,
				repgen.leaderState.r,
				AttestedReportMany{
					msg.Report.AttributedObservations,
					sigs,
				},
			})
			repgen.leaderState.phase = phaseFinal
		}
	}
}
