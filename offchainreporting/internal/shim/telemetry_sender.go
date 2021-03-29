package shim

import (
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting/internal/serialization/protobuf"
	"github.com/smartcontractkit/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

type TelemetrySender struct {
	chTelemetry chan<- *protobuf.TelemetryWrapper
	logger      types.Logger
	taper       loghelper.LogarithmicTaper
}

func MakeTelemetrySender(chTelemetry chan<- *protobuf.TelemetryWrapper, logger types.Logger) TelemetrySender {
	return TelemetrySender{chTelemetry, logger, loghelper.LogarithmicTaper{}}
}

func (ts TelemetrySender) send(t *protobuf.TelemetryWrapper) {
	select {
	case ts.chTelemetry <- t:
		ts.taper.Reset(func(oldCount uint64) {
			ts.logger.Info("TelemetrySender: stopped dropping telemetry", types.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		ts.taper.Trigger(func(newCount uint64) {
			ts.logger.Warn("TelemetrySender: dropping telemetry", types.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (ts TelemetrySender) RoundStarted(
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	leader types.OracleID,
) {
	ts.send(&protobuf.TelemetryWrapper{
		Wrapped: &protobuf.TelemetryWrapper_RoundStarted{&protobuf.TelemetryRoundStarted{
			ConfigDigest: configDigest[:],
			Epoch:        uint64(epoch),
			Round:        uint64(round),
			Leader:       uint64(leader),
			Time:         uint64(time.Now().UnixNano()),
		}},
	})
}
