package protocol

import (
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type TelemetrySender interface {
	RoundStarted(
		configDigest types.ConfigDigest,
		epoch uint64,
		round uint8,
		leader commontypes.OracleID,
	)
}
