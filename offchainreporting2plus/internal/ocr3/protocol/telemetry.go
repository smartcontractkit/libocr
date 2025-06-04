package protocol

import (
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type TelemetrySender interface {
	EpochStarted(
		configDigest types.ConfigDigest,
		epoch uint32,
		leader commontypes.OracleID,
	)

	RoundStarted(
		configDigest types.ConfigDigest,
		epoch uint64,
		seqNr uint64,
		round uint64,
		leader commontypes.OracleID,
	)

	TransmissionScheduleComputed(
		configDigest types.ConfigDigest,
		seqNr uint64,
		index int,
		now time.Time,
		isOverride bool,
		schedule map[commontypes.OracleID]time.Duration,
		ok bool,
	)

	TransmissionShouldAcceptAttestedReportComputed(
		configDigest types.ConfigDigest,
		seqNr uint64,
		index int,
		result bool,
		ok bool,
	)

	TransmissionShouldTransmitAcceptedReportComputed(
		configDigest types.ConfigDigest,
		seqNr uint64,
		index int,
		result bool,
		ok bool,
	)
}
