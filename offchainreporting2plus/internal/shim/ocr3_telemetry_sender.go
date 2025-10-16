package shim

import (
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type OCR3TelemetrySender struct {
	chTelemetry                 chan<- *serialization.TelemetryWrapper
	logger                      commontypes.Logger
	taper                       loghelper.LogarithmicTaper
	enableTransmissionTelemetry bool
}

func NewOCR3TelemetrySender(chTelemetry chan<- *serialization.TelemetryWrapper, logger commontypes.Logger, enableTransmissionTelemetry bool) *OCR3TelemetrySender {
	return &OCR3TelemetrySender{chTelemetry, logger, loghelper.LogarithmicTaper{}, enableTransmissionTelemetry}
}

func (ts *OCR3TelemetrySender) send(t *serialization.TelemetryWrapper) {
	select {
	case ts.chTelemetry <- t:
		ts.taper.Reset(func(oldCount uint64) {
			ts.logger.Info("OCR3TelemetrySender: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		ts.taper.Trigger(func(newCount uint64) {
			ts.logger.Warn("OCR3TelemetrySender: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (ts *OCR3TelemetrySender) EpochStarted(
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

func (ts *OCR3TelemetrySender) RoundStarted(
	configDigest types.ConfigDigest,
	epoch uint64,
	seqNr uint64,
	round uint64,
	leader commontypes.OracleID,
) {
	t := time.Now().UnixNano()
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_RoundStarted{&serialization.TelemetryRoundStarted{
			ConfigDigest: configDigest[:],
			Epoch:        epoch,
			Round:        round,
			Leader:       uint64(leader),
			Time:         uint64(t),
			SeqNr:        seqNr,
		}},
		UnixTimeNanoseconds: t,
	})
}

func (ts *OCR3TelemetrySender) TransmissionScheduleComputed(
	configDigest types.ConfigDigest,
	seqNr uint64,
	index int,
	now time.Time,
	isOverride bool,
	schedule map[commontypes.OracleID]time.Duration,
	ok bool,
) {
	if !ts.enableTransmissionTelemetry {
		return
	}
	scheduleDelayNanosecondsPerNode := make(map[uint64]uint64, len(schedule))
	for oracle, delay := range schedule {
		scheduleDelayNanosecondsPerNode[uint64(oracle)] = uint64(delay.Nanoseconds())
	}
	t := time.Now().UnixNano()
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_TransmissionScheduleComputed{&serialization.TelemetryTransmissionScheduleComputed{
			ConfigDigest:                    configDigest[:],
			SeqNr:                           seqNr,
			Index:                           uint64(index),
			UnixTimeNanoseconds:             uint64(now.UnixNano()),
			IsOverride:                      isOverride,
			ScheduleDelayNanosecondsPerNode: scheduleDelayNanosecondsPerNode,
			Ok:                              ok,
		}},
		UnixTimeNanoseconds: t,
	})
}

func (ts *OCR3TelemetrySender) TransmissionShouldAcceptAttestedReportComputed(
	configDigest types.ConfigDigest,
	seqNr uint64,
	index int,
	result bool,
	ok bool,
) {
	if !ts.enableTransmissionTelemetry {
		return
	}
	t := time.Now().UnixNano()
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_TransmissionShouldAcceptAttestedReportComputed{&serialization.TelemetryTransmissionShouldAcceptAttestedReportComputed{
			ConfigDigest: configDigest[:],
			SeqNr:        seqNr,
			Index:        uint64(index),
			Result:       result,
			Ok:           ok,
		}},
		UnixTimeNanoseconds: t,
	})
}

func (ts *OCR3TelemetrySender) TransmissionShouldTransmitAcceptedReportComputed(
	configDigest types.ConfigDigest,
	seqNr uint64,
	index int,
	result bool,
	ok bool,
) {
	if !ts.enableTransmissionTelemetry {
		return
	}
	t := time.Now().UnixNano()
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_TransmissionShouldTransmitAcceptedReportComputed{&serialization.TelemetryTransmissionShouldTransmitAcceptedReportComputed{
			ConfigDigest: configDigest[:],
			SeqNr:        seqNr,
			Index:        uint64(index),
			Result:       result,
			Ok:           ok,
		}},
		UnixTimeNanoseconds: t,
	})
}
