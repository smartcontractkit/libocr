package ocr3_1config

import (
	"math"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/maxmaxserializationlimits"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const (
	defaultSmallRequestSizeMinRequestToSameOracleInterval = 10 * time.Millisecond

	assumedRTT                     = 500 * time.Millisecond
	assumedBandwidthBitsPerSecond  = 100e6 // 100Mbit
	assumedBandwidthBytesPerSecond = assumedBandwidthBitsPerSecond / 8
)

// transferDuration calculates the duration required to transfer the given
// number of bytes at the assumed bandwidth
func transferDuration(bytes int) time.Duration {
	seconds := float64(bytes) / float64(assumedBandwidthBytesPerSecond)
	return time.Duration(seconds * float64(time.Second))
}

func roundUpToTenthOfSecond(duration time.Duration) time.Duration {
	tenthsOfSecond := float64(duration.Milliseconds()) / 100
	return time.Duration(math.Ceil(tenthsOfSecond)) * 100 * time.Millisecond
}

func DefaultDeltaInitial() time.Duration {
	return roundUpToTenthOfSecond(
		3*assumedRTT/2 +
			transferDuration(maxmaxserializationlimits.MaxMaxEpochStartRequestBytes*types.MaxOracles+maxmaxserializationlimits.MaxMaxEpochStartBytes))
}

func DefaultDeltaReportsPlusPrecursorRequest() time.Duration {
	return roundUpToTenthOfSecond(
		assumedRTT +
			transferDuration(maxmaxserializationlimits.MaxMaxReportsPlusPrecursorRequestBytes+maxmaxserializationlimits.MaxMaxReportsPlusPrecursorBytes))
}

func DefaultDeltaBlockSyncResponseTimeout() time.Duration {
	return roundUpToTenthOfSecond(
		assumedRTT +
			transferDuration(maxmaxserializationlimits.MaxMaxBlockSyncRequestBytes+maxmaxserializationlimits.MaxMaxBlockSyncResponseBytes))
}

func DefaultDeltaTreeSyncResponseTimeout() time.Duration {
	return roundUpToTenthOfSecond(
		assumedRTT +
			transferDuration(maxmaxserializationlimits.MaxMaxTreeSyncChunkRequestBytes+maxmaxserializationlimits.MaxMaxTreeSyncChunkResponseBytes))
}

func DefaultDeltaBlobChunkResponseTimeout() time.Duration {
	return roundUpToTenthOfSecond(
		assumedRTT +
			transferDuration(maxmaxserializationlimits.MaxMaxBlobChunkRequestBytes+maxmaxserializationlimits.MaxMaxBlobChunkResponseBytes))
}

const (
	DefaultDeltaResend = 5 * time.Second

	DefaultDeltaStateSyncSummaryInterval                = 5 * time.Second
	DefaultDeltaBlockSyncMinRequestToSameOracleInterval = defaultSmallRequestSizeMinRequestToSameOracleInterval

	DefaultMaxBlocksPerBlockSyncResponse = 2
	DefaultMaxParallelRequestedBlocks    = 100

	DefaultDeltaTreeSyncMinRequestToSameOracleInterval = defaultSmallRequestSizeMinRequestToSameOracleInterval

	DefaultMaxTreeSyncChunkKeys = 1024

	// A tree sync chunk must always fit at least 1 maximally sized (using maxmax) key-value pair
	DefaultMaxTreeSyncChunkKeysPlusValuesBytes = ocr3_1types.MaxMaxKeyValueKeyBytes + ocr3_1types.MaxMaxKeyValueValueBytes

	DefaultMaxParallelTreeSyncChunkFetches = 8

	DefaultSnapshotInterval               = 10_000
	DefaultMaxHistoricalSnapshotsRetained = 10

	DefaultDeltaBlobOfferMinRequestToSameOracleInterval = defaultSmallRequestSizeMinRequestToSameOracleInterval
	DefaultDeltaBlobOfferResponseTimeout                = 10 * time.Second

	DefaultDeltaBlobBroadcastGrace = 10 * time.Millisecond

	DefaultDeltaBlobChunkMinRequestToSameOracleInterval = defaultSmallRequestSizeMinRequestToSameOracleInterval

	DefaultBlobChunkBytes = 1_000_000 // 1MB
)
