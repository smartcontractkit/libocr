package protocol

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common/scheduler"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/permutation"
	"github.com/RoSpaceDev/libocr/subprocesses"
)

const ContractTransmitterTimeoutWarningGracePeriod = 50 * time.Millisecond

func RunTransmission[RI any](
	ctx context.Context,

	chReportAttestationToTransmission <-chan EventToTransmission[RI],
	config ocr3config.SharedConfig,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	reportingPlugin ocr3types.ReportingPlugin[RI],
	telemetrySender TelemetrySender,
) {
	sched := scheduler.NewScheduler[EventAttestedReport[RI]]()
	defer sched.Close()

	t := transmissionState[RI]{
		ctx,
		subprocesses.Subprocesses{},

		chReportAttestationToTransmission,
		config,
		contractTransmitter,
		id,
		localConfig,
		logger.MakeUpdated(commontypes.LogFields{"proto": "transmission"}),
		reportingPlugin,
		telemetrySender,

		sched,
	}
	t.run()
}

type transmissionState[RI any] struct {
	ctx  context.Context
	subs subprocesses.Subprocesses

	chReportAttestationToTransmission <-chan EventToTransmission[RI]
	config                            ocr3config.SharedConfig
	contractTransmitter               ocr3types.ContractTransmitter[RI]
	id                                commontypes.OracleID
	localConfig                       types.LocalConfig
	logger                            loghelper.LoggerWithContext
	reportingPlugin                   ocr3types.ReportingPlugin[RI]
	telemetrySender                   TelemetrySender

	scheduler *scheduler.Scheduler[EventAttestedReport[RI]]
}

// run runs the event loop for the local transmission protocol
func (t *transmissionState[RI]) run() {
	t.logger.Info("Transmission: running", nil)

	chDone := t.ctx.Done()
	for {
		select {
		case ev := <-t.chReportAttestationToTransmission:
			ev.processTransmission(t)
		case ev := <-t.scheduler.Scheduled():
			t.scheduled(ev)
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			t.logger.Info("Transmission: winding down", nil)
			t.subs.Wait()
			t.logger.Info("Transmission: exiting", nil)
			return
		default:
		}
	}
}

func (t *transmissionState[RI]) eventAttestedReport(ev EventAttestedReport[RI]) {
	now := time.Now()

	t.subs.Go(func() {
		t.backgroundEventAttestedReport(t.ctx, now, ev)
	})
}

func (t *transmissionState[RI]) backgroundEventAttestedReport(ctx context.Context, start time.Time, ev EventAttestedReport[RI]) {
	delays, ok := t.transmitDelays(ev.SeqNr, ev.Index, ev.TransmissionScheduleOverride)
	t.telemetrySender.TransmissionScheduleComputed(
		t.config.ConfigDigest,
		ev.SeqNr,
		ev.Index,
		start,
		ev.TransmissionScheduleOverride != nil,
		delays,
		ok,
	)
	if !ok {
		return
	}
	delay, ok := delays[t.id]
	if !ok {
		t.logger.Debug("dropping EventAttestedReport because we're not included in transmission schedule", commontypes.LogFields{
			"seqNr":                        ev.SeqNr,
			"index":                        ev.Index,
			"transmissionScheduleOverride": ev.TransmissionScheduleOverride != nil,
		})
		return
	}

	shouldAccept, ok := common.CallPlugin[bool](
		ctx,
		t.logger,
		commontypes.LogFields{
			"seqNr": ev.SeqNr,
			"index": ev.Index,
		},
		"ShouldAcceptAttestedReport",
		t.config.MaxDurationShouldAcceptAttestedReport,
		func(ctx context.Context) (bool, error) {
			return t.reportingPlugin.ShouldAcceptAttestedReport(
				ctx,
				ev.SeqNr,
				ev.AttestedReport.ReportWithInfo,
			)
		},
	)
	t.telemetrySender.TransmissionShouldAcceptAttestedReportComputed(t.config.ConfigDigest, ev.SeqNr, ev.Index, shouldAccept, ok)
	if !ok {
		return
	}

	if !shouldAccept {
		t.logger.Debug("ReportingPlugin.ShouldAcceptAttestedReport returned false", commontypes.LogFields{
			"seqNr": ev.SeqNr,
			"index": ev.Index,
		})
		return
	}

	t.logger.Debug("accepted AttestedReport for transmission", commontypes.LogFields{
		"seqNr":                        ev.SeqNr,
		"index":                        ev.Index,
		"delay":                        delay.String(),
		"transmissionScheduleOverride": ev.TransmissionScheduleOverride != nil,
	})
	t.scheduler.ScheduleDeadline(ev, start.Add(delay))
}

func (t *transmissionState[RI]) scheduled(ev EventAttestedReport[RI]) {
	t.subs.Go(func() {
		t.backgroundScheduled(t.ctx, ev)
	})
}

func (t *transmissionState[RI]) backgroundScheduled(ctx context.Context, ev EventAttestedReport[RI]) {
	shouldTransmit, ok := common.CallPlugin[bool](
		ctx,
		t.logger,
		commontypes.LogFields{
			"seqNr": ev.SeqNr,
			"index": ev.Index,
		},
		"ShouldTransmitAcceptedReport",
		t.config.MaxDurationShouldTransmitAcceptedReport,
		func(ctx context.Context) (bool, error) {
			return t.reportingPlugin.ShouldTransmitAcceptedReport(
				ctx,
				ev.SeqNr,
				ev.AttestedReport.ReportWithInfo,
			)
		},
	)
	t.telemetrySender.TransmissionShouldTransmitAcceptedReportComputed(t.config.ConfigDigest, ev.SeqNr, ev.Index, shouldTransmit, ok)
	if !ok {
		return
	}

	if !shouldTransmit {
		t.logger.Info("ReportingPlugin.ShouldTransmitAcceptedReport returned false", commontypes.LogFields{
			"seqNr": ev.SeqNr,
			"index": ev.Index,
		})
		return
	}

	t.logger.Debug("transmitting report", commontypes.LogFields{
		"seqNr": ev.SeqNr,
		"index": ev.Index,
	})

	{
		transmitCtx, transmitCancel := context.WithTimeout(
			ctx,
			t.localConfig.ContractTransmitterTransmitTimeout,
		)
		defer transmitCancel()

		ins := loghelper.NewIfNotStopped(
			t.localConfig.ContractTransmitterTransmitTimeout+ContractTransmitterTimeoutWarningGracePeriod,
			func() {
				t.logger.Error("ContractTransmitter.Transmit is taking too long", commontypes.LogFields{
					"maxDuration": t.localConfig.ContractTransmitterTransmitTimeout.String(),
					"seqNr":       ev.SeqNr,
					"index":       ev.Index,
				})
			},
		)

		err := t.contractTransmitter.Transmit(
			transmitCtx,
			t.config.ConfigDigest,
			ev.SeqNr,
			ev.AttestedReport.ReportWithInfo,
			ev.AttestedReport.AttributedSignatures,
		)

		ins.Stop()

		if err != nil {
			t.logger.Error("ContractTransmitter.Transmit error", commontypes.LogFields{"error": err})
			return
		}

	}

	t.logger.Info("🚀 successfully invoked ContractTransmitter.Transmit", commontypes.LogFields{
		"seqNr": ev.SeqNr,
		"index": ev.Index,
	})
}

func (t *transmissionState[RI]) transmitPermutationKey(seqNr uint64, index int) [16]byte {
	transmissionOrderKey := t.config.TransmissionOrderKey()
	mac := hmac.New(sha256.New, transmissionOrderKey[:])
	_ = binary.Write(mac, binary.BigEndian, seqNr)
	_ = binary.Write(mac, binary.BigEndian, uint64(index))

	var key [16]byte
	_ = copy(key[:], mac.Sum(nil))
	return key
}

func (t *transmissionState[RI]) transmitDelaysFromOverride(seqNr uint64, index int, transmissionScheduleOverride ocr3types.TransmissionSchedule) (delays map[commontypes.OracleID]time.Duration, ok bool) {
	if len(transmissionScheduleOverride.TransmissionDelays) != len(transmissionScheduleOverride.Transmitters) {
		t.logger.Error("invalid TransmissionScheduleOverride, cannot compute delay, lengths do not match", commontypes.LogFields{
			"seqNr":                        seqNr,
			"index":                        index,
			"transmissionScheduleOverride": transmissionScheduleOverride,
		})
		return nil, false
	}

	for _, oid := range transmissionScheduleOverride.Transmitters {
		if !(0 <= int(oid) && int(oid) < t.config.N()) {
			t.logger.Error("invalid TransmissionScheduleOverride, cannot compute delay, oracle id out of bounds", commontypes.LogFields{
				"seqNr":                        seqNr,
				"index":                        index,
				"transmissionScheduleOverride": transmissionScheduleOverride,
			})
			return nil, false
		}
	}

	// Permutation from index of oracle in transmissionScheduleOverride.Transmitters to transmission order index.
	pi := permutation.Permutation(len(transmissionScheduleOverride.Transmitters), t.transmitPermutationKey(seqNr, index))

	result := make(map[commontypes.OracleID]time.Duration, len(transmissionScheduleOverride.Transmitters))

	for transmitterIndex, oid := range transmissionScheduleOverride.Transmitters {
		if _, ok := result[oid]; ok {
			t.logger.Error("invalid TransmissionScheduleOverride, cannot compute delay, duplicate oracle id", commontypes.LogFields{
				"seqNr":                        seqNr,
				"index":                        index,
				"transmissionScheduleOverride": transmissionScheduleOverride,
			})
			return nil, false
		}
		result[oid] = transmissionScheduleOverride.TransmissionDelays[pi[transmitterIndex]]
	}

	return result, true
}

func (t *transmissionState[RI]) transmitDelaysDefault(seqNr uint64, index int) map[commontypes.OracleID]time.Duration {
	// Permutation from transmission order index to oracle id
	piInv := make([]int, t.config.N())
	{
		// Permutation from oracle id to transmission order index. The
		// permutations are structured in this "inverted" way for historical
		// compatibility
		pi := permutation.Permutation(t.config.N(), t.transmitPermutationKey(seqNr, index))
		for i := range pi {
			piInv[pi[i]] = i
		}
	}

	result := make(map[commontypes.OracleID]time.Duration, t.config.N())

	accumulatedStageSize := 0
	for stageIdx, stageSize := range t.config.S {
		// i is the index of the oracle sorted by transmission order
		for i := accumulatedStageSize; i < accumulatedStageSize+stageSize; i++ {
			if i >= len(piInv) {
				// Index is larger than index of the last oracle. This happens
				// when sum(S) > N.
				break
			}
			oracleId := commontypes.OracleID(piInv[i])
			result[oracleId] = time.Duration(stageIdx) * t.config.DeltaStage
		}

		accumulatedStageSize += stageSize
	}

	return result
}

func (t *transmissionState[RI]) transmitDelays(seqNr uint64, index int, transmissionScheduleOverride *ocr3types.TransmissionSchedule) (delays map[commontypes.OracleID]time.Duration, ok bool) {
	if transmissionScheduleOverride != nil {
		return t.transmitDelaysFromOverride(seqNr, index, *transmissionScheduleOverride)
	} else {
		return t.transmitDelaysDefault(seqNr, index), true
	}
}
