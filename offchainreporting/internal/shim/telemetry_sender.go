package shim

import (
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting/internal/serialization/protobuf"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

type TelemetrySender struct {
	chTelemetry chan<- *protobuf.TelemetryWrapper
}

func MakeTelemetrySender(chTelemetry chan<- *protobuf.TelemetryWrapper) TelemetrySender {
	return TelemetrySender{chTelemetry}
}

func (ts TelemetrySender) send(t *protobuf.TelemetryWrapper) {
	select {
	case ts.chTelemetry <- t:
	default:
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
