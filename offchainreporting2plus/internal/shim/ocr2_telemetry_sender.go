package shim

import (
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr2/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type OCR2TelemetrySender struct {
	chTelemetry                 chan<- *serialization.TelemetryWrapper
	logger                      commontypes.Logger
	taper                       loghelper.LogarithmicTaper
	enableTransmissionTelemetry bool
}

func NewOCR2TelemetrySender(chTelemetry chan<- *serialization.TelemetryWrapper, logger commontypes.Logger, enableTransmissionTelemetry bool) *OCR2TelemetrySender {
	return &OCR2TelemetrySender{chTelemetry, logger, loghelper.LogarithmicTaper{}, enableTransmissionTelemetry}
}

func (ts *OCR2TelemetrySender) send(t *serialization.TelemetryWrapper) {
	select {
	case ts.chTelemetry <- t:
		ts.taper.Reset(func(oldCount uint64) {
			ts.logger.Info("OCR2TelemetrySender: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		ts.taper.Trigger(func(newCount uint64) {
			ts.logger.Warn("OCR2TelemetrySender: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (ts *OCR2TelemetrySender) EpochStarted(
	configDigest types.ConfigDigest,
	epoch uint32,
	leader commontypes.OracleID,
) {
	t := time.Now().UnixNano()
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_EpochStarted{&serialization.TelemetryEpochStarted{
			ConfigDigest: configDigest[:],
			Epoch:        uint64(epoch),
			Leader:       uint64(leader),
		}},
		UnixTimeNanoseconds: t,
	})
}

func (ts *OCR2TelemetrySender) RoundStarted(
	reportTimestamp types.ReportTimestamp,
	leader commontypes.OracleID,
) {
	t := time.Now().UnixNano()
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_RoundStarted{&serialization.TelemetryRoundStarted{
			ConfigDigest: reportTimestamp.ConfigDigest[:],
			Epoch:        uint64(reportTimestamp.Epoch),
			Round:        uint64(reportTimestamp.Round),
			Leader:       uint64(leader),
			Time:         uint64(t),
		}},
		UnixTimeNanoseconds: t,
	})
}

func (ts *OCR2TelemetrySender) TransmissionScheduleComputed(
	reportTimestamp types.ReportTimestamp,
	now time.Time,
	schedule map[commontypes.OracleID]time.Duration,
) {
	if !ts.enableTransmissionTelemetry {
		return
	}

	t := time.Now().UnixNano()

	scheduleDelayNanosecondsPerNode := make(map[uint64]uint64, len(schedule))
	for oracle, delay := range schedule {
		scheduleDelayNanosecondsPerNode[uint64(oracle)] = uint64(delay.Nanoseconds())
	}

	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_TransmissionScheduleComputed{&serialization.TelemetryTransmissionScheduleComputed{
			ConfigDigest:                    reportTimestamp.ConfigDigest[:],
			Epoch:                           uint64(reportTimestamp.Epoch),
			Round:                           uint64(reportTimestamp.Round),
			UnixTimeNanoseconds:             uint64(now.UnixNano()),
			ScheduleDelayNanosecondsPerNode: scheduleDelayNanosecondsPerNode,
		}},
		UnixTimeNanoseconds: t,
	})
}

func (ts *OCR2TelemetrySender) TransmissionShouldAcceptFinalizedReportComputed(
	reportTimestamp types.ReportTimestamp,
	result bool,
	ok bool,
) {
	if !ts.enableTransmissionTelemetry {
		return
	}

	t := time.Now().UnixNano()
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_TransmissionShouldAcceptFinalizedReportComputed{&serialization.TelemetryTransmissionShouldAcceptFinalizedReportComputed{
			ConfigDigest: reportTimestamp.ConfigDigest[:],
			Epoch:        uint64(reportTimestamp.Epoch),
			Round:        uint64(reportTimestamp.Round),
			Result:       result,
			Ok:           ok,
		}},
		UnixTimeNanoseconds: t,
	})
}

func (ts *OCR2TelemetrySender) TransmissionShouldTransmitAcceptedReportComputed(
	reportTimestamp types.ReportTimestamp,
	result bool,
	ok bool,
) {
	if !ts.enableTransmissionTelemetry {
		return
	}

	t := time.Now().UnixNano()
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_TransmissionShouldTransmitAcceptedReportComputed{&serialization.TelemetryTransmissionShouldTransmitAcceptedReportComputed{
			ConfigDigest: reportTimestamp.ConfigDigest[:],
			Epoch:        uint64(reportTimestamp.Epoch),
			Round:        uint64(reportTimestamp.Round),
			Result:       result,
			Ok:           ok,
		}},
		UnixTimeNanoseconds: t,
	})
}
