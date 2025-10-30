package limits

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"time"

	"github.com/smartcontractkit/libocr/internal/jmt"
	"github.com/smartcontractkit/libocr/internal/mt"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type OCR3_1SerializedLengthLimits struct {
	MaxLenMsgNewEpoch                    int
	MaxLenMsgEpochStartRequest           int
	MaxLenMsgEpochStart                  int
	MaxLenMsgRoundStart                  int
	MaxLenMsgObservation                 int
	MaxLenMsgProposal                    int
	MaxLenMsgPrepare                     int
	MaxLenMsgCommit                      int
	MaxLenMsgReportSignatures            int
	MaxLenMsgReportsPlusPrecursorRequest int
	MaxLenMsgReportsPlusPrecursor        int
	MaxLenMsgStateSyncSummary            int
	MaxLenMsgBlockSyncRequest            int
	MaxLenMsgBlockSyncResponse           int
	MaxLenMsgTreeSyncChunkRequest        int
	MaxLenMsgTreeSyncChunkResponse       int
	MaxLenMsgBlobOffer                   int
	MaxLenMsgBlobOfferResponse           int
	MaxLenMsgBlobChunkRequest            int
	MaxLenMsgBlobChunkResponse           int
}

func OCR3_1Limits(cfg ocr3_1config.PublicConfig, pluginLimits ocr3_1types.ReportingPluginLimits, maxSigLen int) (types.BinaryNetworkEndpointLimits, types.BinaryNetworkEndpointLimits, OCR3_1SerializedLengthLimits, error) {
	overflow := false

	// These two helper functions add/multiply together a bunch of numbers and set overflow to true if the result
	// lies outside the range [0; math.MaxInt32]. We compare with int32 rather than int to be independent of
	// the underlying architecture.
	add := func(xs ...int) int {
		sum := big.NewInt(0)
		for _, x := range xs {
			sum.Add(sum, big.NewInt(int64(x)))
		}
		if !(big.NewInt(0).Cmp(sum) <= 0 && sum.Cmp(big.NewInt(int64(math.MaxInt32))) <= 0) {
			overflow = true
		}
		return int(sum.Int64())
	}
	mul := func(xs ...int) int {
		prod := big.NewInt(1)
		for _, x := range xs {
			prod.Mul(prod, big.NewInt(int64(x)))
		}
		if !(big.NewInt(0).Cmp(prod) <= 0 && prod.Cmp(big.NewInt(int64(math.MaxInt32))) <= 0) {
			overflow = true
		}
		return int(prod.Int64())
	}

	const repeatedOverhead = 10
	const sigOverhead = 10
	const overhead = 256

	maxLenStateWriteSet := add(
		pluginLimits.MaxKeyValueModifiedKeysPlusValuesBytes,
		mul(
			pluginLimits.MaxKeyValueModifiedKeys,
			repeatedOverhead,
		),
	)
	maxLenCertifiedPrepareOrCommit := add(mul(ed25519.SignatureSize+sigOverhead, cfg.ByzQuorumSize()),
		sha256.Size*5,
		overhead)
	maxLenMsgNewEpoch := overhead
	maxLenMsgEpochStartRequest := add(maxLenCertifiedPrepareOrCommit, overhead)
	maxLenMsgEpochStart := add(maxLenCertifiedPrepareOrCommit, mul(ed25519.SignatureSize+sigOverhead, cfg.ByzQuorumSize()), overhead)
	maxLenMsgRoundStart := add(pluginLimits.MaxQueryBytes, overhead)
	maxLenMsgObservation := add(pluginLimits.MaxObservationBytes, overhead)
	maxLenMsgProposal := add(mul(add(pluginLimits.MaxObservationBytes, ed25519.SignatureSize+sigOverhead), cfg.N()), overhead)
	maxLenMsgPrepare := overhead
	maxLenMsgCommit := overhead
	maxLenMsgReportSignatures := add(mul(add(maxSigLen, sigOverhead), pluginLimits.MaxReportCount), overhead)
	maxLenMsgReportsPlusPrecursorRequest := overhead
	maxLenMsgReportsPlusPrecursor := add(pluginLimits.MaxReportsPlusPrecursorBytes, overhead)
	maxLenMsgStateSyncSummary := overhead

	// tree sync messages
	maxLenMsgTreeSyncChunkRequest := overhead
	maxLenMsgTreeSyncChunkResponseBoundingLeaf := add(
		len(jmt.Digest{}), // leaf key
		len(jmt.Digest{}), // leaf value
		mul( // siblings
			jmt.MaxProofLength,
			add(
				repeatedOverhead,
				len(jmt.Digest{}),
			),
		),
	)
	maxLenMsgTreeSyncChunkResponseKeyValues := add(
		cfg.GetMaxTreeSyncChunkKeysPlusValuesBytes(),
		mul( // repeated overheads
			cfg.GetMaxTreeSyncChunkKeys(),
			repeatedOverhead, // key-value
			add(
				repeatedOverhead, // key
				repeatedOverhead, // value
			),
		),
	)
	maxLenMsgTreeSyncChunkResponse := add(
		overhead,
		mul(
			jmt.MaxBoundingLeaves,
			add(
				repeatedOverhead,
				maxLenMsgTreeSyncChunkResponseBoundingLeaf,
			),
		),
		maxLenMsgTreeSyncChunkResponseKeyValues,
	)

	// block sync messages
	maxLenMsgBlockSyncRequest := overhead
	maxLenAttestedStateTransitionBlock := add(maxLenCertifiedPrepareOrCommit, maxLenStateWriteSet)
	maxLenMsgBlockSyncResponse := add(mul(cfg.GetMaxBlocksPerBlockSyncResponse(), maxLenAttestedStateTransitionBlock), overhead)

	// blob exchange messages
	const blobDigestSize = len(protocol.BlobDigest{})
	cfgBlobChunkSize := cfg.GetBlobChunkBytes()
	maxNumBlobChunks := (pluginLimits.MaxBlobPayloadBytes + cfgBlobChunkSize - 1) / cfgBlobChunkSize
	maxBlobChunksDigestProofElements := bits.Len(uint(maxNumBlobChunks)) + 1
	maxLenMsgBlobOffer := add(blobDigestSize, overhead)
	maxLenMsgBlobChunkRequest := add(blobDigestSize, overhead)
	maxLenMsgBlobChunkResponse := add(blobDigestSize, cfgBlobChunkSize,
		mul(maxBlobChunksDigestProofElements, add(repeatedOverhead, len(mt.Digest{}))), overhead)
	maxLenMsgBlobOfferResponse := add(blobDigestSize, ed25519.SignatureSize+sigOverhead, overhead)

	maxDefaultPriorityMessageSize := max(
		maxLenMsgNewEpoch,
		maxLenMsgEpochStartRequest,
		maxLenMsgEpochStart,
		maxLenMsgRoundStart,
		maxLenMsgProposal,
		maxLenMsgPrepare,
		maxLenMsgCommit,
		maxLenMsgReportSignatures,
		maxLenMsgReportsPlusPrecursorRequest,
		maxLenMsgBlobOffer,
		maxLenMsgBlobChunkRequest,
	)

	maxLowPriorityMessageSize := max(
		maxLenMsgStateSyncSummary,
		maxLenMsgBlockSyncRequest,
		maxLenMsgTreeSyncChunkRequest,
	)

	minRoundInterval := math.Max(float64(cfg.DeltaRound), float64(cfg.DeltaGrace))

	minEpochInterval := math.Min(float64(cfg.DeltaProgress), math.Min(float64(cfg.GetDeltaInitial()), float64(cfg.RMax)*float64(minRoundInterval)))

	defaultPriorityMessagesRate := (1.0*float64(time.Second)/float64(cfg.GetDeltaResend()) +
		3.0*float64(time.Second)/minEpochInterval +
		6.0*float64(time.Second)/float64(minRoundInterval) +
		2.0*float64(time.Second)/float64(cfg.GetDeltaBlobOfferMinRequestToSameOracleInterval()) +
		1.0*float64(time.Second)/float64(cfg.GetDeltaBlobChunkMinRequestToSameOracleInterval())) * 1.2

	lowPriorityMessagesRate := (1.0*float64(time.Second)/float64(cfg.GetDeltaBlockSyncMinRequestToSameOracleInterval()) +
		1.0*float64(time.Second)/float64(cfg.GetDeltaTreeSyncMinRequestToSameOracleInterval()) +
		1.0*float64(time.Second)/float64(cfg.GetDeltaStateSyncSummaryInterval())) * 1.2

	defaultPriorityMessagesCapacity := mul(13, 3)
	lowPriorityMessagesCapacity := mul(3, 3)

	// we don't multiply bytesRate by a safetyMargin since we already have a generous overhead on each message

	defaultPriorityBytesRate := float64(time.Second)/float64(cfg.GetDeltaResend())*float64(maxLenMsgNewEpoch) +
		float64(time.Second)/float64(minEpochInterval)*float64(maxLenMsgNewEpoch) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgPrepare) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgCommit) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgReportSignatures) +
		float64(time.Second)/float64(minEpochInterval)*float64(maxLenMsgEpochStart) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgRoundStart) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgProposal) +
		float64(time.Second)/float64(minEpochInterval)*float64(maxLenMsgEpochStartRequest) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgReportsPlusPrecursorRequest) +
		float64(time.Second)/float64(cfg.GetDeltaBlobOfferMinRequestToSameOracleInterval())*float64(maxLenMsgBlobOffer) + // blob-related messages
		float64(time.Second)/float64(cfg.GetDeltaBlobChunkMinRequestToSameOracleInterval())*float64(maxLenMsgBlobChunkRequest)

	lowPriorityBytesRate := float64(time.Second)/float64(cfg.GetDeltaStateSyncSummaryInterval())*float64(maxLenMsgStateSyncSummary) +
		float64(time.Second)/float64(cfg.GetDeltaBlockSyncMinRequestToSameOracleInterval())*float64(maxLenMsgBlockSyncRequest) +
		float64(time.Second)/float64(cfg.GetDeltaTreeSyncMinRequestToSameOracleInterval())*float64(maxLenMsgTreeSyncChunkRequest)

	defaultPriorityBytesCapacity := mul(add(
		maxLenMsgNewEpoch,
		maxLenMsgNewEpoch,
		maxLenMsgEpochStartRequest,
		maxLenMsgEpochStart,
		maxLenMsgRoundStart,
		maxLenMsgProposal,
		maxLenMsgPrepare,
		maxLenMsgCommit,
		maxLenMsgReportSignatures,
		maxLenMsgReportsPlusPrecursorRequest,
		maxLenMsgBlobOffer,
		maxLenMsgBlobChunkRequest,
	), 3)

	lowPriorityBytesCapacity := mul(add(
		maxLenMsgStateSyncSummary,
		maxLenMsgBlockSyncRequest,
		maxLenMsgTreeSyncChunkRequest,
	), 3)

	if overflow {
		// this should not happen due to us checking the limits in types.go
		return types.BinaryNetworkEndpointLimits{}, types.BinaryNetworkEndpointLimits{}, OCR3_1SerializedLengthLimits{}, fmt.Errorf("int32 overflow while computing bandwidth limits")
	}

	return types.BinaryNetworkEndpointLimits{
			maxDefaultPriorityMessageSize,
			defaultPriorityMessagesRate,
			defaultPriorityMessagesCapacity,
			defaultPriorityBytesRate,
			defaultPriorityBytesCapacity,
		},
		types.BinaryNetworkEndpointLimits{
			maxLowPriorityMessageSize,
			lowPriorityMessagesRate,
			lowPriorityMessagesCapacity,
			lowPriorityBytesRate,
			lowPriorityBytesCapacity,
		},
		OCR3_1SerializedLengthLimits{
			maxLenMsgNewEpoch,
			maxLenMsgEpochStartRequest,
			maxLenMsgEpochStart,
			maxLenMsgRoundStart,
			maxLenMsgObservation,
			maxLenMsgProposal,
			maxLenMsgPrepare,
			maxLenMsgCommit,
			maxLenMsgReportSignatures,
			maxLenMsgReportsPlusPrecursorRequest,
			maxLenMsgReportsPlusPrecursor,
			maxLenMsgStateSyncSummary,
			maxLenMsgBlockSyncRequest,
			maxLenMsgBlockSyncResponse,
			maxLenMsgTreeSyncChunkRequest,
			maxLenMsgTreeSyncChunkResponse,
			maxLenMsgBlobOffer,
			maxLenMsgBlobOfferResponse,
			maxLenMsgBlobChunkRequest,
			maxLenMsgBlobChunkResponse,
		},
		nil
}
