package shim

import (
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting/internal/serialization/protobuf"
	"github.com/RoSpaceDev/libocr/offchainreporting/types"
)

type TelemetrySender struct {
	chTelemetry chan<- *protobuf.TelemetryWrapper
	logger      commontypes.Logger
	taper       loghelper.LogarithmicTaper
}

func NewTelemetrySender(chTelemetry chan<- *protobuf.TelemetryWrapper, logger commontypes.Logger) *TelemetrySender {
	return &TelemetrySender{chTelemetry, logger, loghelper.LogarithmicTaper{}}
}

func (ts *TelemetrySender) send(t *protobuf.TelemetryWrapper) {
	select {
	case ts.chTelemetry <- t:
		ts.taper.Reset(func(oldCount uint64) {
			ts.logger.Info("TelemetrySender: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		ts.taper.Trigger(func(newCount uint64) {
			ts.logger.Warn("TelemetrySender: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (ts *TelemetrySender) RoundStarted(
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	leader commontypes.OracleID,
) {
	t := time.Now().UnixNano()
	ts.send(&protobuf.TelemetryWrapper{
		Wrapped: &protobuf.TelemetryWrapper_RoundStarted{&protobuf.TelemetryRoundStarted{
			ConfigDigest: configDigest[:],
			Epoch:        uint64(epoch),
			Round:        uint64(round),
			Leader:       uint64(leader),
			Time:         uint64(t),
		}},
		UnixTimeNanoseconds: t,
	})
}
