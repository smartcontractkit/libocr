package protocol

import (
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type TelemetrySender interface {
	EpochStarted(
		configDigest types.ConfigDigest,
		epoch uint32,
		leader commontypes.OracleID,
	)

	RoundStarted(
		reportTimestamp types.ReportTimestamp,
		leader commontypes.OracleID,
	)

	TransmissionScheduleComputed(
		reportTimestamp types.ReportTimestamp,
		now time.Time,
		schedule map[commontypes.OracleID]time.Duration,
	)

	TransmissionShouldAcceptFinalizedReportComputed(
		reportTimestamp types.ReportTimestamp,
		result bool,
		ok bool,
	)

	TransmissionShouldTransmitAcceptedReportComputed(
		reportTimestamp types.ReportTimestamp,
		result bool,
		ok bool,
	)
}
