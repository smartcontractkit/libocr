package protocol

import (
	"context"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/scheduler"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func RunReportFinalization[RI any](
	ctx context.Context,

	chNetToReportFinalization <-chan MessageToReportFinalizationWithSender[RI],
	chReportFinalizationToTransmission chan<- EventToTransmission[RI],
	chReportGenerationToReportFinalization <-chan EventToReportFinalization[RI],
	config ocr3config.SharedConfig,
	contractSigner ocr3types.OnchainKeyring[RI],
	contractTransmitter ocr3types.ContractTransmitter[RI],
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	reportingPlugin ocr3types.OCR3Plugin[RI],
) {
	newReportFinalizationState(ctx, chNetToReportFinalization,
		chReportFinalizationToTransmission, chReportGenerationToReportFinalization,
		config, contractSigner, contractTransmitter, logger, netSender, reportingPlugin).run()
}

const minExpirationAgeRounds int = 10
const expirationAgeDuration = 10 * time.Minute
const maxExpirationAgeRounds int = 1_000

const deltaRequestCertifiedCommit = 200 * time.Millisecond

type reportFinalizationState[RI any] struct {
	ctx context.Context

	chNetToReportFinalization              <-chan MessageToReportFinalizationWithSender[RI]
	chReportFinalizationToTransmission     chan<- EventToTransmission[RI]
	chReportGenerationToReportFinalization <-chan EventToReportFinalization[RI]
	config                                 ocr3config.SharedConfig
	contractSigner                         ocr3types.OnchainKeyring[RI]
	contractTransmitter                    ocr3types.ContractTransmitter[RI]
	logger                                 loghelper.LoggerWithContext
	netSender                              NetworkSender[RI]
	reportingPlugin                        ocr3types.OCR3Plugin[RI]

	scheduler *scheduler.Scheduler[EventMissingOutcome[RI]]
	// reap() is used to prevent unbounded state growth of finalized
	finalized             map[uint64]*finalizationRound[RI]
	finalizedHighestSeqNr uint64
}

type finalizationRound[RI any] struct {
	certifiedCommit *CertifiedPrepareOrCommitCommit
	reportsWithInfo []ocr3types.ReportWithInfo[RI]
	signatures      map[commontypes.OracleID]*reportSignatures
	startedFetch    bool
	complete        bool
}

type reportSignatures struct {
	signatures       [][]byte
	validSignatures  *bool
	requestedOutcome bool
	suppliedOutcome  bool
}

// func (fr finalizationRound) finalized(f int) bool {
// 	return len(fr.reportSignatures) > f
// }

func (repfin *reportFinalizationState[RI]) run() {
	for {
		select {
		case msg := <-repfin.chNetToReportFinalization:
			msg.msg.processReportFinalization(repfin, msg.sender)
		case ev := <-repfin.chReportGenerationToReportFinalization:
			ev.processReportFinalization(repfin)
		case ev := <-repfin.scheduler.Scheduled():
			ev.processReportFinalization(repfin)
		case <-repfin.ctx.Done():
		}

		// ensure prompt exit
		select {
		case <-repfin.ctx.Done():
			repfin.logger.Info("ReportFinalization: exiting", nil)
			repfin.scheduler.Close()
			return
		default:
		}
	}
}

func (repfin *reportFinalizationState[RI]) messageFinal(
	msg MessageFinal[RI],
	sender commontypes.OracleID,
) {

	if repfin.isExpired(msg.SeqNr) {
		repfin.logger.Debug("ignoring MessageFinal for expired SeqNr", commontypes.LogFields{
			"seqNr":  msg.SeqNr,
			"sender": sender,
		})
		return
	}

	if _, ok := repfin.finalized[msg.SeqNr]; !ok {
		repfin.finalized[msg.SeqNr] = &finalizationRound[RI]{
			nil,
			nil,
			map[commontypes.OracleID]*reportSignatures{},
			false,
			false,
		}
	}

	if _, ok := repfin.finalized[msg.SeqNr].signatures[sender]; ok {
		repfin.logger.Debug("ignoring MessageFinal with duplicate signature", commontypes.LogFields{
			"seqNr":  msg.SeqNr,
			"sender": sender,
		})
		return
	}

	repfin.finalized[msg.SeqNr].signatures[sender] = &reportSignatures{
		msg.ReportSignatures,
		nil,
		false,
		false,
	}

	repfin.tryComplete(msg.SeqNr)
}

func (repfin *reportFinalizationState[RI]) eventMissingOutcome(ev EventMissingOutcome[RI]) {
	if len(repfin.finalized[ev.SeqNr].reportsWithInfo) != 0 {
		repfin.logger.Debug("dropping EventMissingOutcome, already have Outcome", commontypes.LogFields{
			"seqNr": ev.SeqNr,
		})
		return
	}

	repfin.tryRequestCertifiedCommit(ev.SeqNr)
}

func (repfin *reportFinalizationState[RI]) messageRequestCertifiedCommit(msg MessageRequestCertifiedCommit[RI], sender commontypes.OracleID) {
	if repfin.finalized[msg.SeqNr] == nil || repfin.finalized[msg.SeqNr].certifiedCommit == nil {
		repfin.logger.Warn("dropping MessageRequestCertifiedCommit for outcome with unknown certified commit", commontypes.LogFields{
			"seqNr":  msg.SeqNr,
			"sender": sender,
		})
		return
	}

	repfin.logger.Debug("sending MessageSupplyCertifiedCommit", commontypes.LogFields{
		"seqNr": msg.SeqNr,
		"to":    sender,
	})
	repfin.netSender.SendTo(MessageSupplyCertifiedCommit[RI]{*repfin.finalized[msg.SeqNr].certifiedCommit}, sender)

}

func (repfin *reportFinalizationState[RI]) messageSupplyCertifiedCommit(msg MessageSupplyCertifiedCommit[RI], sender commontypes.OracleID) {
	if repfin.finalized[msg.CertifiedCommit.SeqNr] == nil {
		repfin.logger.Warn("dropping MessageSupplyCertifiedCommit for unknown seqNr", commontypes.LogFields{
			"seqNr":  msg.CertifiedCommit.SeqNr,
			"sender": sender,
		})
		return
	}

	senderSigs := repfin.finalized[msg.CertifiedCommit.SeqNr].signatures[sender]
	requestedOutcome := senderSigs != nil && senderSigs.requestedOutcome
	suppliedOutcome := senderSigs != nil && senderSigs.suppliedOutcome
	if !(requestedOutcome && !suppliedOutcome) {
		repfin.logger.Warn("dropping MessageSupplyCertifiedCommit for sender who doesn't have pending request", commontypes.LogFields{
			"seqNr":            msg.CertifiedCommit.SeqNr,
			"sender":           sender,
			"requestedOutcome": requestedOutcome,
			"suppliedOutcome":  suppliedOutcome,
		})
		return
	}

	senderSigs.suppliedOutcome = true

	if repfin.finalized[msg.CertifiedCommit.SeqNr].certifiedCommit != nil {
		repfin.logger.Debug("dropping redundant MessageSupplyCertifiedCommit", commontypes.LogFields{
			"seqNr":  msg.CertifiedCommit.SeqNr,
			"sender": sender,
		})
		return
	}

	if err := msg.CertifiedCommit.Verify(repfin.config.ConfigDigest, repfin.config.OracleIdentities, repfin.config.N(), repfin.config.F); err != nil {
		repfin.logger.Warn("dropping MessageSupplyCertifiedCommit with invalid certified commit", commontypes.LogFields{
			"seqNr":  msg.CertifiedCommit.SeqNr,
			"sender": sender,
		})
		return
	}

	repfin.logger.Debug("triggering eventDeliver based on valid MessageSupplyCertifiedCommit", commontypes.LogFields{
		"seqNr":  msg.CertifiedCommit.SeqNr,
		"sender": sender,
	})

	repfin.eventDeliver(EventDeliver[RI]{msg.CertifiedCommit})
}

func (repfin *reportFinalizationState[RI]) tryRequestCertifiedCommit(seqNr uint64) {
	candidates := make([]commontypes.OracleID, 0, repfin.config.N())
	for signer, sig := range repfin.finalized[seqNr].signatures {
		if sig.requestedOutcome {
			continue
		}
		candidates = append(candidates, signer)
	}

	if len(candidates) == 0 {

		return
	}

	randomOracle := candidates[rand.Intn(len(candidates))]
	repfin.finalized[seqNr].signatures[randomOracle].requestedOutcome = true
	repfin.logger.Debug("sending MessageRequestCertifiedCommit", commontypes.LogFields{
		"seqNr": seqNr,
		"to":    randomOracle,
	})
	repfin.netSender.SendTo(MessageRequestCertifiedCommit[RI]{seqNr}, randomOracle)
	repfin.scheduler.ScheduleDelay(EventMissingOutcome[RI]{seqNr}, deltaRequestCertifiedCommit)
}

func (repfin *reportFinalizationState[RI]) tryComplete(seqNr uint64) {
	if repfin.finalized[seqNr].complete {
		repfin.logger.Debug("cannot complete, already completed", commontypes.LogFields{
			"seqNr": seqNr,
		})
		return
	}

	if len(repfin.finalized[seqNr].reportsWithInfo) == 0 {
		if len(repfin.finalized[seqNr].signatures) <= repfin.config.F {
			repfin.logger.Debug("cannot complete, missing reports and signatures", commontypes.LogFields{
				"seqNr": seqNr,
			})
		} else if !repfin.finalized[seqNr].startedFetch {
			repfin.finalized[seqNr].startedFetch = true
			repfin.tryRequestCertifiedCommit(seqNr)
		}
		return
	}

	reportsWithInfo := repfin.finalized[seqNr].reportsWithInfo
	goodSigs := 0
	var aossPerReport [][]types.AttributedOnchainSignature = make([][]types.AttributedOnchainSignature, len(reportsWithInfo))
	for signer, sig := range repfin.finalized[seqNr].signatures {
		if sig.validSignatures == nil {
			validSignatures := repfin.verifySignatures(repfin.config.OracleIdentities[signer].OnchainPublicKey, seqNr, reportsWithInfo, sig.signatures)
			sig.validSignatures = &validSignatures
		}
		if sig.validSignatures != nil && *sig.validSignatures {
			goodSigs++

			for i := range reportsWithInfo {
				aossPerReport[i] = append(aossPerReport[i], types.AttributedOnchainSignature{sig.signatures[i], signer})
			}
		}
		if goodSigs > repfin.config.F {
			break
		}
	}

	if goodSigs <= repfin.config.F {
		repfin.logger.Debug("cannot complete, insufficient number of signatures", commontypes.LogFields{
			"seqNr":    seqNr,
			"goodSigs": goodSigs,
		})
		return
	}

	repfin.finalized[seqNr].complete = true

	repfin.logger.Info("🚀 Ready to broadcast", commontypes.LogFields{
		"seqNr":   seqNr,
		"reports": len(reportsWithInfo),
	})

	for i := range reportsWithInfo {
		select {
		case repfin.chReportFinalizationToTransmission <- EventTransmit[RI]{
			seqNr,
			i,
			AttestedReportMany[RI]{
				reportsWithInfo[i],
				aossPerReport[i],
			},
		}:
		case <-repfin.ctx.Done():
		}
	}

	repfin.reap()
}

func (repfin *reportFinalizationState[RI]) verifySignatures(publicKey types.OnchainPublicKey, seqNr uint64, reportsWithInfo []ocr3types.ReportWithInfo[RI], signatures [][]byte) bool {
	if len(reportsWithInfo) != len(signatures) {
		return false
	}

	n := runtime.GOMAXPROCS(0)
	if (len(reportsWithInfo)+3)/4 < n {
		n = (len(reportsWithInfo) + 3) / 4
	}

	var wg sync.WaitGroup
	wg.Add(n)

	var mutex sync.Mutex
	allValid := true

	for k := 0; k < n; k++ {
		k := k

		go func() {
			defer wg.Done()
			for i := k; i < len(reportsWithInfo); i += n {
				if i%n != k {
					panic("bug")
				}

				mutex.Lock()
				allValidCopy := allValid
				mutex.Unlock()

				if !allValidCopy {
					return
				}

				if !repfin.contractSigner.Verify(publicKey, repfin.config.ConfigDigest, seqNr, reportsWithInfo[i], signatures[i]) {
					mutex.Lock()
					allValid = false
					mutex.Unlock()
					return
				}
			}
		}()
	}

	wg.Wait()

	return allValid
}

func (repfin *reportFinalizationState[RI]) eventDeliver(ev EventDeliver[RI]) {
	if repfin.finalized[ev.CertifiedCommit.SeqNr] != nil && repfin.finalized[ev.CertifiedCommit.SeqNr].reportsWithInfo != nil {
		repfin.logger.Debug("Skipping delivery of already delivered outcome", commontypes.LogFields{
			"seqNr": ev.CertifiedCommit.SeqNr,
		})
		return
	}

	reportsWithInfo, err := repfin.reportingPlugin.Reports(ev.CertifiedCommit.SeqNr, ev.CertifiedCommit.Outcome)
	if err != nil {
		repfin.logger.Error("ReportingPlugin.Reports failed", commontypes.LogFields{
			"seqNr": ev.CertifiedCommit.SeqNr,
			"error": err,
		})
		return
	}

	if reportsWithInfo == nil {
		repfin.logger.Info("ReportingPlugin.Reports returned no reports, skipping", commontypes.LogFields{
			"seqNr": ev.CertifiedCommit.SeqNr,
		})
		return
	}

	var sigs [][]byte
	for i, reportWithInfo := range reportsWithInfo {
		sig, err := repfin.contractSigner.Sign(repfin.config.ConfigDigest, ev.CertifiedCommit.SeqNr, reportWithInfo)
		if err != nil {
			repfin.logger.Error("Error while signing report", commontypes.LogFields{
				"seqNr": ev.CertifiedCommit.SeqNr,
				"index": i,
				"error": err,
			})
			return
		}
		sigs = append(sigs, sig)
	}

	if _, ok := repfin.finalized[ev.CertifiedCommit.SeqNr]; !ok {
		repfin.finalized[ev.CertifiedCommit.SeqNr] = &finalizationRound[RI]{
			&ev.CertifiedCommit,
			reportsWithInfo,
			map[commontypes.OracleID]*reportSignatures{},
			false,
			false,
		}
	} else {
		repfin.finalized[ev.CertifiedCommit.SeqNr].certifiedCommit = &ev.CertifiedCommit
		repfin.finalized[ev.CertifiedCommit.SeqNr].reportsWithInfo = reportsWithInfo
	}

	if repfin.finalizedHighestSeqNr < ev.CertifiedCommit.SeqNr {
		repfin.finalizedHighestSeqNr = ev.CertifiedCommit.SeqNr
	}

	repfin.logger.Debug("Broadcasting MessageFinal", commontypes.LogFields{
		"seqNr": ev.CertifiedCommit.SeqNr,
	})

	repfin.netSender.Broadcast(MessageFinal[RI]{
		ev.CertifiedCommit.SeqNr,
		sigs,
	})

	// no need to call tryComplete since receipt of our own MessageFinal will do so
}

// func (repfin *reportFinalizationState[RI]) finalize(msg MessageFinal) {
// 	repfin.logger.Debug("finalizing report", commontypes.LogFields{
// 		"epoch": msg.Epoch,
// 		"round": msg.Round,
// 	})

// 	epochRound := EpochRound{msg.Epoch, msg.Round}

// 	repfin.finalized[epochRound] = struct{}{}
// 	if repfin.finalizedLatest.Less(epochRound) {
// 		repfin.finalizedLatest = epochRound
// 	}

// 	repfin.netSender.Broadcast(MessageFinalEcho{msg}) // send [ FINALECHO, e, r, O] to all p_j ∈ P

// 	select {
// 	case repfin.chReportFinalizationToTransmission <- EventTransmit(msg):
// 	case <-repfin.ctx.Done():
// 	}

// 	repfin.reap()
// }

func (repfin *reportFinalizationState[RI]) isExpired(seqNr uint64) bool {
	highest := repfin.finalizedHighestSeqNr
	expired := uint64(0)
	expirationAgeRounds := uint64(repfin.expirationAgeRounds())
	if highest > expirationAgeRounds {
		expired = highest - expirationAgeRounds
	}
	return seqNr <= expired
}

// reap expired entries from repfin.finalized to prevent unbounded state growth
func (repfin *reportFinalizationState[RI]) reap() {
	if len(repfin.finalized) <= 2*repfin.expirationAgeRounds() {
		return
	}
	// A long time ago in a galaxy far, far away, Go used to leak memory when
	// repeatedly adding and deleting from the same map without ever exceeding
	// some maximum length. Fortunately, this is no longer the case
	// https://go-review.googlesource.com/c/go/+/25049/
	for seqNr := range repfin.finalized {
		if repfin.isExpired(seqNr) {
			delete(repfin.finalized, seqNr)
		}
	}
}

// The age (denoted in rounds) after which a report is considered expired and
// will automatically be dropped
func (repfin *reportFinalizationState[RI]) expirationAgeRounds() int {
	// number of rounds in a window of duration expirationAgeDuration
	age := math.Ceil(expirationAgeDuration.Seconds() / repfin.config.DeltaRound.Seconds())

	if age < float64(minExpirationAgeRounds) {
		age = float64(minExpirationAgeRounds)
	}
	if math.IsNaN(age) || age > float64(maxExpirationAgeRounds) {
		age = float64(maxExpirationAgeRounds)
	}

	return int(age)
}

func newReportFinalizationState[RI any](
	ctx context.Context,

	chNetToReportFinalization <-chan MessageToReportFinalizationWithSender[RI],
	chReportFinalizationToTransmission chan<- EventToTransmission[RI],
	chReportGenerationToReportFinalization <-chan EventToReportFinalization[RI],
	config ocr3config.SharedConfig,
	contractSigner ocr3types.OnchainKeyring[RI],
	contractTransmitter ocr3types.ContractTransmitter[RI],
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	reportingPlugin ocr3types.OCR3Plugin[RI],
) *reportFinalizationState[RI] {
	return &reportFinalizationState[RI]{
		ctx,

		chNetToReportFinalization,
		chReportFinalizationToTransmission,
		chReportGenerationToReportFinalization,
		config,
		contractSigner,
		contractTransmitter,
		logger.MakeUpdated(commontypes.LogFields{"proto": "repfin"}),
		netSender,
		reportingPlugin,

		scheduler.NewScheduler[EventMissingOutcome[RI]](),
		map[uint64]*finalizationRound[RI]{},
		0,
	}
}
