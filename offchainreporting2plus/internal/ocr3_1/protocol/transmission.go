package protocol

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"slices"
	"time"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common/scheduler"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
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
	reportingPlugin ocr3_1types.ReportingPlugin[RI],
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
	reportingPlugin                   ocr3_1types.ReportingPlugin[RI]

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
	var delay time.Duration
	{
		delayMaybe := t.transmitDelay(ev.SeqNr, ev.Index, ev.TransmissionScheduleOverride)
		if delayMaybe == nil {
			t.logger.Debug("dropping EventAttestedReport because we're not included in transmission schedule", commontypes.LogFields{
				"seqNr":                        ev.SeqNr,
				"index":                        ev.Index,
				"transmissionScheduleOverride": ev.TransmissionScheduleOverride != nil,
			})
			return
		}
		delay = *delayMaybe
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

func (t *transmissionState[RI]) transmitDelayFromOverride(seqNr uint64, index int, transmissionScheduleOverride ocr3types.TransmissionSchedule) *time.Duration {
	if len(transmissionScheduleOverride.TransmissionDelays) != len(transmissionScheduleOverride.Transmitters) {
		t.logger.Error("invalid TransmissionScheduleOverride, cannot compute delay", commontypes.LogFields{
			"seqNr":                        seqNr,
			"index":                        index,
			"transmissionScheduleOverride": transmissionScheduleOverride,
		})
		return nil
	}

	oracleIndex := slices.Index(transmissionScheduleOverride.Transmitters, t.id)
	if oracleIndex < 0 {
		return nil
	}
	pi := permutation.Permutation(len(transmissionScheduleOverride.TransmissionDelays), t.transmitPermutationKey(seqNr, index))
	delay := transmissionScheduleOverride.TransmissionDelays[pi[oracleIndex]]
	return &delay
}

func (t *transmissionState[RI]) transmitDelayDefault(seqNr uint64, index int) *time.Duration {
	pi := permutation.Permutation(t.config.N(), t.transmitPermutationKey(seqNr, index))
	sum := 0
	for i, s := range t.config.S {
		sum += s
		if pi[t.id] < sum {
			result := time.Duration(i) * t.config.DeltaStage
			return &result
		}
	}
	return nil
}

func (t *transmissionState[RI]) transmitDelay(seqNr uint64, index int, transmissionScheduleOverride *ocr3types.TransmissionSchedule) *time.Duration {
	if transmissionScheduleOverride != nil {
		return t.transmitDelayFromOverride(seqNr, index, *transmissionScheduleOverride)
	} else {
		return t.transmitDelayDefault(seqNr, index)
	}
}
