package protocol

import (
	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting/types"
)

type TelemetrySender interface {
	RoundStarted(
		configDigest types.ConfigDigest,
		epoch uint32,
		round uint8,
		leader commontypes.OracleID,
	)
}
