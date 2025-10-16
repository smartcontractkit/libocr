package protocol

import (
	"context"
	"crypto/rand"
	"math"
	"math/big"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common/scheduler"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
)

func RunReportAttestation[RI any](
	ctx context.Context,

	chNetToReportAttestation <-chan MessageToReportAttestationWithSender[RI],
	chOutcomeGenerationToReportAttestation <-chan EventToReportAttestation[RI],
	chReportAttestationToTransmission chan<- EventToTransmission[RI],
	config ocr3config.SharedConfig,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	onchainKeyring ocr3types.OnchainKeyring[RI],
	reportingPlugin ocr3types.ReportingPlugin[RI],
) {
	sched := scheduler.NewScheduler[EventMissingOutcome[RI]]()
	defer sched.Close()

	newReportAttestationState(ctx, chNetToReportAttestation,
		chOutcomeGenerationToReportAttestation, chReportAttestationToTransmission,
		config, contractTransmitter, logger, netSender, onchainKeyring, reportingPlugin, sched).run()
}

const expiryMinRounds int = 10
const expiryDuration = 1 * time.Minute
const expiryMaxRounds int = 50

const lookaheadMinRounds int = 4
const lookaheadDuration = 30 * time.Second
const lookaheadMaxRounds int = 10

type reportAttestationState[RI any] struct {
	ctx  context.Context
	subs subprocesses.Subprocesses

	chNetToReportAttestation               <-chan MessageToReportAttestationWithSender[RI]
	chOutcomeGenerationToReportAttestation <-chan EventToReportAttestation[RI]
	chReportAttestationToTransmission      chan<- EventToTransmission[RI]
	config                                 ocr3config.SharedConfig
	contractTransmitter                    ocr3types.ContractTransmitter[RI]
	logger                                 loghelper.LoggerWithContext
	netSender                              NetworkSender[RI]
	onchainKeyring                         ocr3types.OnchainKeyring[RI]
	reportingPlugin                        ocr3types.ReportingPlugin[RI]

	scheduler    *scheduler.Scheduler[EventMissingOutcome[RI]]
	chLocalEvent chan EventComputedReports[RI]
	// reap() is used to prevent unbounded state growth of rounds.

	rounds map[uint64]*round[RI]

	// Highest sequence number for which we know a certified commit exists.
	// This is used for determining the window of rounds we keep in memory.
	// Computed as select_largest(f+1, highestReportSignaturesSeqNr).
	highWaterMark                uint64
	highestReportSignaturesSeqNr []uint64
}

type round[RI any] struct {
	ctx                     context.Context             // should always be initialized when a round[RI] is initiated
	ctxCancel               context.CancelFunc          // should always be initialized when a round[RI] is initiated
	verifiedCertifiedCommit *CertifiedCommit            // only stores certifiedCommit whose qc has been verified
	reportsPlus             *[]ocr3types.ReportPlus[RI] // cache result of ReportingPlugin.Reports(certifiedCommit.SeqNr, certifiedCommit.Outcome)
	oracles                 []oracle                    // always initialized to be of length n
	startedFetch            bool
	complete                bool
}

// oracle contains information about interactions with oracles (self & others)
type oracle struct {
	signatures      [][]byte
	validSignatures *bool
	weRequested     bool
	theyServiced    bool
	weServiced      bool
}

func (repatt *reportAttestationState[RI]) run() {
	repatt.logger.Info("ReportAttestation: running", nil)

	for {
		select {
		case ev := <-repatt.chLocalEvent:
			ev.processReportAttestation(repatt)
		case msg := <-repatt.chNetToReportAttestation:
			msg.msg.processReportAttestation(repatt, msg.sender)
		case ev := <-repatt.chOutcomeGenerationToReportAttestation:
			ev.processReportAttestation(repatt)
		case ev := <-repatt.scheduler.Scheduled():
			ev.processReportAttestation(repatt)
		case <-repatt.ctx.Done():
		}

		// ensure prompt exit
		select {
		case <-repatt.ctx.Done():
			repatt.logger.Info("ReportAttestation: winding down", nil)
			repatt.subs.Wait()
			repatt.scheduler.Close()
			repatt.logger.Info("ReportAttestation: exiting", nil)
			return
		default:
		}
	}
}

func (repatt *reportAttestationState[RI]) messageReportSignatures(
	msg MessageReportSignatures[RI],
	sender commontypes.OracleID,
) {
	repatt.tryReap(msg.SeqNr, sender)

	if repatt.isBeyondExpiry(msg.SeqNr) {
		repatt.logger.Debug("dropping MessageReportSignatures for expired seqNr", commontypes.LogFields{
			"msgSeqNr": msg.SeqNr,
			"sender":   sender,
		})
		return
	}

	if repatt.isBeyondLookahead(msg.SeqNr) {
		repatt.logger.Debug("dropping MessageReportSignatures for seqNr beyond lookahead", commontypes.LogFields{
			"msgSeqNr": msg.SeqNr,
			"sender":   sender,
		})
		return
	}

	if _, ok := repatt.rounds[msg.SeqNr]; !ok {
		ctx, cancel := context.WithCancel(repatt.ctx)
		repatt.rounds[msg.SeqNr] = &round[RI]{
			ctx,
			cancel,
			nil,
			nil,
			make([]oracle, repatt.config.N()),
			false,
			false,
		}
	}

	if len(repatt.rounds[msg.SeqNr].oracles[sender].signatures) != 0 {
		repatt.logger.Debug("dropping MessageReportSignatures with duplicate signature", commontypes.LogFields{
			"msgSeqNr": msg.SeqNr,
			"sender":   sender,
		})
		return
	}

	repatt.rounds[msg.SeqNr].oracles[sender].signatures = msg.ReportSignatures

	repatt.tryComplete(msg.SeqNr)
}

func (repatt *reportAttestationState[RI]) eventMissingOutcome(ev EventMissingOutcome[RI]) {
	if repatt.rounds[ev.SeqNr] == nil {
		repatt.logger.Debug("dropping EventMissingOutcome for unknown seqNr", commontypes.LogFields{
			"evSeqNr":       ev.SeqNr,
			"highWaterMark": repatt.highWaterMark,
			"expiryRounds":  repatt.expiryRounds(),
		})
		return
	}

	if repatt.rounds[ev.SeqNr].verifiedCertifiedCommit != nil {
		repatt.logger.Debug("dropping EventMissingOutcome, already have Outcome", commontypes.LogFields{
			"evSeqNr": ev.SeqNr,
		})
		return
	}

	repatt.tryRequestCertifiedCommit(ev.SeqNr)
}

func (repatt *reportAttestationState[RI]) messageCertifiedCommitRequest(msg MessageCertifiedCommitRequest[RI], sender commontypes.OracleID) {
	if repatt.rounds[msg.SeqNr] == nil {
		repatt.logger.Debug("dropping MessageCertifiedCommitRequest for unknown seqNr", commontypes.LogFields{
			"msgSeqNr":      msg.SeqNr,
			"sender":        sender,
			"highWaterMark": repatt.highWaterMark,
			"expiryRounds":  repatt.expiryRounds(),
		})
		return
	}

	if repatt.rounds[msg.SeqNr].verifiedCertifiedCommit == nil {
		repatt.logger.Debug("dropping MessageCertifiedCommitRequest for outcome with unknown certified commit", commontypes.LogFields{
			"msgSeqNr": msg.SeqNr,
			"sender":   sender,
		})
		return
	}

	if repatt.rounds[msg.SeqNr].oracles[sender].weServiced {
		repatt.logger.Warn("dropping duplicate MessageCertifiedCommitRequest", commontypes.LogFields{
			"msgSeqNr": msg.SeqNr,
			"sender":   sender,
		})
		return
	}

	repatt.rounds[msg.SeqNr].oracles[sender].weServiced = true

	repatt.logger.Debug("sending MessageCertifiedCommit", commontypes.LogFields{
		"msgSeqNr": msg.SeqNr,
		"to":       sender,
	})
	repatt.netSender.SendTo(MessageCertifiedCommit[RI]{*repatt.rounds[msg.SeqNr].verifiedCertifiedCommit}, sender)
}

func (repatt *reportAttestationState[RI]) messageCertifiedCommit(msg MessageCertifiedCommit[RI], sender commontypes.OracleID) {
	if repatt.rounds[msg.CertifiedCommit.SeqNr] == nil {
		repatt.logger.Warn("dropping MessageCertifiedCommit for unknown seqNr", commontypes.LogFields{
			"msgSeqNr":      msg.CertifiedCommit.SeqNr,
			"sender":        sender,
			"highWaterMark": repatt.highWaterMark,
			"expiryRounds":  repatt.expiryRounds(),
		})
		return
	}

	oracle := &repatt.rounds[msg.CertifiedCommit.SeqNr].oracles[sender]
	if !(oracle.weRequested && !oracle.theyServiced) {
		repatt.logger.Warn("dropping unexpected MessageCertifiedCommit", commontypes.LogFields{
			"msgSeqNr":     msg.CertifiedCommit.SeqNr,
			"sender":       sender,
			"weRequested":  oracle.weRequested,
			"theyServiced": oracle.theyServiced,
		})
		return
	}

	oracle.theyServiced = true

	if repatt.rounds[msg.CertifiedCommit.SeqNr].verifiedCertifiedCommit != nil {
		repatt.logger.Debug("dropping redundant MessageCertifiedCommit", commontypes.LogFields{
			"msgSeqNr": msg.CertifiedCommit.SeqNr,
			"sender":   sender,
		})
		return
	}

	if err := msg.CertifiedCommit.Verify(repatt.config.ConfigDigest, repatt.config.OracleIdentities, repatt.config.ByzQuorumSize()); err != nil {
		repatt.logger.Warn("dropping MessageCertifiedCommit with invalid certified commit", commontypes.LogFields{
			"msgSeqNr": msg.CertifiedCommit.SeqNr,
			"sender":   sender,
		})
		return
	}

	repatt.logger.Debug("received valid MessageCertifiedCommit", commontypes.LogFields{
		"msgSeqNr": msg.CertifiedCommit.SeqNr,
		"sender":   sender,
	})

	repatt.receivedVerifiedCertifiedCommit(msg.CertifiedCommit)
}

func (repatt *reportAttestationState[RI]) tryRequestCertifiedCommit(seqNr uint64) {
	candidates := make([]commontypes.OracleID, 0, repatt.config.N())
	for oracleID, oracle := range repatt.rounds[seqNr].oracles {
		// avoid duplicate requests
		if oracle.weRequested {
			continue
		}
		// avoid requesting from oracles that haven't sent MessageReportSignatures
		if len(oracle.signatures) == 0 {
			continue
		}
		candidates = append(candidates, commontypes.OracleID(oracleID))
	}

	if len(candidates) == 0 {

		return
	}

	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates))))
	if err != nil {
		repatt.logger.Critical("unexpected error returned by rand.Int", commontypes.LogFields{
			"error": err,
		})
		return
	}
	randomCandidate := candidates[int(randomIndex.Int64())]
	repatt.rounds[seqNr].oracles[randomCandidate].weRequested = true
	repatt.logger.Debug("sending MessageCertifiedCommitRequest", commontypes.LogFields{
		"seqNr": seqNr,
		"to":    randomCandidate,
	})
	repatt.netSender.SendTo(MessageCertifiedCommitRequest[RI]{seqNr}, randomCandidate)
	repatt.scheduler.ScheduleDelay(EventMissingOutcome[RI]{seqNr}, repatt.config.DeltaCertifiedCommitRequest)
}

func (repatt *reportAttestationState[RI]) tryComplete(seqNr uint64) {
	if repatt.rounds[seqNr].complete {
		repatt.logger.Debug("cannot complete, already completed", commontypes.LogFields{
			"seqNr": seqNr,
		})
		return
	}

	if repatt.rounds[seqNr].verifiedCertifiedCommit == nil {
		oraclesThatSentNonemptySignatures := 0
		for _, oracle := range repatt.rounds[seqNr].oracles {
			if len(oracle.signatures) == 0 {
				continue
			}
			oraclesThatSentNonemptySignatures++
		}

		if oraclesThatSentNonemptySignatures <= repatt.config.F {
			repatt.logger.Debug("cannot complete, missing CertifiedCommit and signatures", commontypes.LogFields{
				"oraclesThatSentNonemptySignatures": oraclesThatSentNonemptySignatures,
				"seqNr":                             seqNr,
				"threshold":                         repatt.config.F + 1,
			})
		} else if !repatt.rounds[seqNr].startedFetch {
			repatt.rounds[seqNr].startedFetch = true
			repatt.scheduler.ScheduleDelay(EventMissingOutcome[RI]{seqNr}, repatt.config.DeltaCertifiedCommitRequest)
		}
		return
	}

	if repatt.rounds[seqNr].reportsPlus == nil {
		repatt.logger.Debug("cannot complete, reportsPlus not computed yet", commontypes.LogFields{
			"seqNr": seqNr,
		})
		return
	}
	reportsPlus := *repatt.rounds[seqNr].reportsPlus

	goodSigs := 0
	var aossPerReport [][]types.AttributedOnchainSignature = make([][]types.AttributedOnchainSignature, len(reportsPlus))
	for oracleID := range repatt.rounds[seqNr].oracles {
		oracle := &repatt.rounds[seqNr].oracles[oracleID]
		if len(oracle.signatures) == 0 {
			continue
		}
		if oracle.validSignatures == nil {
			validSignatures := repatt.verifySignatures(
				repatt.config.OracleIdentities[oracleID].OnchainPublicKey,
				seqNr,
				reportsPlus,
				oracle.signatures,
			)
			oracle.validSignatures = &validSignatures
			if !validSignatures {
				// Other less common causes include actually invalid signatures.
				repatt.logger.Warn("report signatures failed to verify. This is commonly caused by non-determinism in the ReportingPlugin", commontypes.LogFields{
					"sender":        oracleID,
					"seqNr":         seqNr,
					"signaturesLen": len(oracle.signatures),
					"reportsLen":    len(reportsPlus),
				})
			}
		}
		if oracle.validSignatures != nil && *oracle.validSignatures {
			goodSigs++

			for i := range reportsPlus {
				aossPerReport[i] = append(aossPerReport[i], types.AttributedOnchainSignature{
					oracle.signatures[i],
					commontypes.OracleID(oracleID),
				})
			}
		}
		if goodSigs > repatt.config.F {
			break
		}
	}

	if goodSigs <= repatt.config.F {
		repatt.logger.Debug("cannot complete, insufficient number of signatures", commontypes.LogFields{
			"seqNr":     seqNr,
			"goodSigs":  goodSigs,
			"threshold": repatt.config.F + 1,
		})
		return
	}

	repatt.rounds[seqNr].complete = true

	repatt.logger.Debug("sending attested reports to transmission protocol", commontypes.LogFields{
		"seqNr":   seqNr,
		"reports": len(reportsPlus),
	})

	for i := range reportsPlus {
		select {
		case repatt.chReportAttestationToTransmission <- EventAttestedReport[RI]{
			seqNr,
			i,
			AttestedReportMany[RI]{
				reportsPlus[i].ReportWithInfo,
				aossPerReport[i],
			},
			reportsPlus[i].TransmissionScheduleOverride,
		}:
		case <-repatt.ctx.Done():
		}
	}
}

func (repatt *reportAttestationState[RI]) verifySignatures(publicKey types.OnchainPublicKey, seqNr uint64, reportsPlus []ocr3types.ReportPlus[RI], signatures [][]byte) bool {
	if len(reportsPlus) != len(signatures) {
		return false
	}

	n := runtime.GOMAXPROCS(0)
	if (len(reportsPlus)+3)/4 < n {
		n = (len(reportsPlus) + 3) / 4
	}

	var wg sync.WaitGroup
	wg.Add(n)

	var mutex sync.Mutex
	allValid := true

	for k := 0; k < n; k++ {
		go func() {
			defer wg.Done()
			for i := k; i < len(reportsPlus); i += n {
				mutex.Lock()
				allValidCopy := allValid
				mutex.Unlock()

				if !allValidCopy {
					return
				}

				if !repatt.onchainKeyring.Verify(publicKey, repatt.config.ConfigDigest, seqNr, reportsPlus[i].ReportWithInfo, signatures[i]) {
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

func (repatt *reportAttestationState[RI]) eventCommittedOutcome(ev EventCommittedOutcome[RI]) {
	repatt.receivedVerifiedCertifiedCommit(ev.CertifiedCommit)
}

func (repatt *reportAttestationState[RI]) receivedVerifiedCertifiedCommit(certifiedCommit CertifiedCommit) {
	if repatt.rounds[certifiedCommit.SeqNr] != nil && repatt.rounds[certifiedCommit.SeqNr].verifiedCertifiedCommit != nil {
		repatt.logger.Debug("dropping redundant CertifiedCommit", commontypes.LogFields{
			"seqNr": certifiedCommit.SeqNr,
		})
		return
	}

	if _, ok := repatt.rounds[certifiedCommit.SeqNr]; !ok {
		ctx, cancel := context.WithCancel(repatt.ctx)
		repatt.rounds[certifiedCommit.SeqNr] = &round[RI]{
			ctx,
			cancel,
			nil,
			nil,
			make([]oracle, repatt.config.N()),
			false,
			false,
		}
	}

	repatt.rounds[certifiedCommit.SeqNr].verifiedCertifiedCommit = &certifiedCommit

	{
		ctx := repatt.rounds[certifiedCommit.SeqNr].ctx
		repatt.subs.Go(func() {
			repatt.backgroundComputeReports(ctx, certifiedCommit)
		})
	}
}

func (repatt *reportAttestationState[RI]) backgroundComputeReports(ctx context.Context, verifiedCertifiedCommit CertifiedCommit) {
	reportsPlus, ok := common.CallPluginFromBackground(
		ctx,
		repatt.logger,
		commontypes.LogFields{"seqNr": verifiedCertifiedCommit.SeqNr},
		"Reports",
		0, // Reports is a pure function and should finish "instantly"
		func(ctx context.Context) ([]ocr3types.ReportPlus[RI], error) {
			return repatt.reportingPlugin.Reports(ctx, verifiedCertifiedCommit.SeqNr, verifiedCertifiedCommit.Outcome)
		},
	)
	if !ok {
		return
	}

	repatt.logger.Debug("successfully invoked ReportingPlugin.Reports", commontypes.LogFields{
		"seqNr":   verifiedCertifiedCommit.SeqNr,
		"reports": len(reportsPlus),
	})

	select {
	case repatt.chLocalEvent <- EventComputedReports[RI]{verifiedCertifiedCommit.SeqNr, reportsPlus}:
	case <-ctx.Done():
		return
	}
}

func (repatt *reportAttestationState[RI]) eventComputedReports(ev EventComputedReports[RI]) {
	if repatt.rounds[ev.SeqNr] == nil {
		repatt.logger.Debug("dropping EventComputedReports for unknown seqNr", commontypes.LogFields{
			"evSeqNr":       ev.SeqNr,
			"highWaterMark": repatt.highWaterMark,
			"expiryRounds":  repatt.expiryRounds(),
		})
		return
	}

	repatt.rounds[ev.SeqNr].reportsPlus = &ev.ReportsPlus

	var sigs [][]byte
	for i, reportPlus := range ev.ReportsPlus {
		sig, err := repatt.onchainKeyring.Sign(repatt.config.ConfigDigest, ev.SeqNr, reportPlus.ReportWithInfo)
		if err != nil {
			repatt.logger.Error("error while signing report", commontypes.LogFields{
				"evSeqNr": ev.SeqNr,
				"index":   i,
				"error":   err,
			})
			return
		}
		sigs = append(sigs, sig)
	}

	repatt.logger.Debug("broadcasting MessageReportSignatures", commontypes.LogFields{
		"evSeqNr": ev.SeqNr,
	})

	repatt.netSender.Broadcast(MessageReportSignatures[RI]{
		ev.SeqNr,
		sigs,
	})

	// no need to call tryComplete since receipt of our own MessageReportSignatures will do so
}

// reap expired rounds if there is a new high water mark
func (repatt *reportAttestationState[RI]) tryReap(seqNr uint64, sender commontypes.OracleID) {
	if repatt.highestReportSignaturesSeqNr[sender] >= seqNr {
		return
	}

	repatt.highestReportSignaturesSeqNr[sender] = seqNr

	var newHighWaterMark uint64
	{
		highestReportSignaturesSeqNr := append([]uint64{}, repatt.highestReportSignaturesSeqNr...)
		sort.Slice(highestReportSignaturesSeqNr, func(i, j int) bool {
			return highestReportSignaturesSeqNr[i] > highestReportSignaturesSeqNr[j]
		})
		newHighWaterMark = highestReportSignaturesSeqNr[repatt.config.F] // (f+1)th largest seqNr
	}

	if repatt.highWaterMark >= newHighWaterMark {
		return
	}

	repatt.highWaterMark = newHighWaterMark // (f+1)th largest seqNr
	repatt.reap()
}

func (repatt *reportAttestationState[RI]) isBeyondExpiry(seqNr uint64) bool {
	expiry := uint64(repatt.expiryRounds())
	if repatt.highWaterMark <= expiry {
		return false
	}
	return seqNr < repatt.highWaterMark-expiry
}

func (repatt *reportAttestationState[RI]) isBeyondLookahead(seqNr uint64) bool {
	lookahead := uint64(repatt.lookaheadRounds())
	if seqNr <= lookahead {
		return false
	}
	return repatt.highWaterMark < seqNr-lookahead
}

// reap expired entries from repatt.rounds to prevent unbounded state growth
func (repatt *reportAttestationState[RI]) reap() {
	maxActiveRoundCount := repatt.expiryRounds() + repatt.lookaheadRounds()
	// only reap if more than ~ a third of the rounds can potentially be discarded
	if 3*len(repatt.rounds) <= 4*maxActiveRoundCount {
		return
	}

	beforeRounds := len(repatt.rounds)

	// A long time ago in a galaxy far, far away, Go used to leak memory when
	// repeatedly adding and deleting from the same map without ever exceeding
	// some maximum length. Fortunately, this is no longer the case
	// https://go-review.googlesource.com/c/go/+/25049/
	for seqNr := range repatt.rounds {
		if repatt.isBeyondExpiry(seqNr) {
			repatt.rounds[seqNr].ctxCancel()
			delete(repatt.rounds, seqNr)
		}
	}

	repatt.logger.Debug("reaped expired rounds", commontypes.LogFields{
		"before":        beforeRounds,
		"after":         len(repatt.rounds),
		"highWaterMark": repatt.highWaterMark,
	})
}

// The age (denoted in rounds) after which a report is considered expired and
// will automatically be dropped
func (repatt *reportAttestationState[RI]) expiryRounds() int {
	return repatt.roundWindowSize(expiryMinRounds, expiryMaxRounds, expiryDuration)
}

// The lookahead (denoted in rounds) after which a report is considered too far in the future and
// will automatically be dropped
func (repatt *reportAttestationState[RI]) lookaheadRounds() int {
	return repatt.roundWindowSize(lookaheadMinRounds, lookaheadMaxRounds, lookaheadDuration)
}

func (repatt *reportAttestationState[RI]) roundWindowSize(minWindowSize int, maxWindowSize int, windowDuration time.Duration) int {
	// number of rounds in a window of duration expirationAgeDuration
	size := math.Ceil(windowDuration.Seconds() / repatt.config.MinRoundInterval().Seconds())

	if size < float64(minWindowSize) {
		size = float64(minWindowSize)
	}
	if math.IsNaN(size) || size > float64(maxWindowSize) {
		size = float64(maxWindowSize)
	}

	return int(math.Ceil(size))
}

func newReportAttestationState[RI any](
	ctx context.Context,

	chNetToReportAttestation <-chan MessageToReportAttestationWithSender[RI],
	chOutcomeGenerationToReportAttestation <-chan EventToReportAttestation[RI],
	chReportAttestationToTransmission chan<- EventToTransmission[RI],
	config ocr3config.SharedConfig,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	logger loghelper.LoggerWithContext,
	netSender NetworkSender[RI],
	onchainKeyring ocr3types.OnchainKeyring[RI],
	reportingPlugin ocr3types.ReportingPlugin[RI],
	sched *scheduler.Scheduler[EventMissingOutcome[RI]],
) *reportAttestationState[RI] {
	return &reportAttestationState[RI]{
		ctx,
		subprocesses.Subprocesses{},

		chNetToReportAttestation,
		chOutcomeGenerationToReportAttestation,
		chReportAttestationToTransmission,
		config,
		contractTransmitter,
		logger.MakeUpdated(commontypes.LogFields{"proto": "repatt"}),
		netSender,
		onchainKeyring,
		reportingPlugin,

		sched,
		make(chan EventComputedReports[RI]),
		map[uint64]*round[RI]{},
		0,
		make([]uint64, config.N()),
	}
}
