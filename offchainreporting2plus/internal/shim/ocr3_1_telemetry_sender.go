package shim

import (
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type OCR3_1TelemetrySender struct {
	chTelemetry chan<- *serialization.TelemetryWrapper
	logger      commontypes.Logger
	taper       loghelper.LogarithmicTaper
}

func NewOCR3_1TelemetrySender(chTelemetry chan<- *serialization.TelemetryWrapper, logger commontypes.Logger) *OCR3_1TelemetrySender {
	return &OCR3_1TelemetrySender{chTelemetry, logger, loghelper.LogarithmicTaper{}}
}

func (ts *OCR3_1TelemetrySender) send(t *serialization.TelemetryWrapper) {
	select {
	case ts.chTelemetry <- t:
		ts.taper.Reset(func(oldCount uint64) {
			ts.logger.Info("NewOCR3_1TelemetrySender: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		ts.taper.Trigger(func(newCount uint64) {
			ts.logger.Warn("NewOCR3_1TelemetrySender: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (ts *OCR3_1TelemetrySender) RoundStarted(
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
