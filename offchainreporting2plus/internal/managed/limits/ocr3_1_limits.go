package limits

import (
	"crypto/ed25519"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/RoSpaceDev/libocr/internal/jmt"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type OCR3_1SerializedLengthLimits struct {
	MaxLenMsgNewEpoch               int
	MaxLenMsgEpochStartRequest      int
	MaxLenMsgEpochStart             int
	MaxLenMsgRoundStart             int
	MaxLenMsgObservation            int
	MaxLenMsgProposal               int
	MaxLenMsgPrepare                int
	MaxLenMsgCommit                 int
	MaxLenMsgReportSignatures       int
	MaxLenMsgCertifiedCommitRequest int
	MaxLenMsgCertifiedCommit        int
	MaxLenMsgStateSyncSummary       int
	MaxLenMsgBlockSyncRequest       int
	MaxLenMsgBlockSyncResponse      int
	MaxLenMsgTreeSyncChunkRequest   int
	MaxLenMsgTreeSyncChunkResponse  int
	MaxLenMsgBlobOffer              int
	MaxLenMsgBlobOfferResponse      int
	MaxLenMsgBlobChunkRequest       int
	MaxLenMsgBlobChunkResponse      int
}

func OCR3_1Limits(cfg ocr3config.PublicConfig, pluginLimits ocr3_1types.ReportingPluginLimits, maxSigLen int) (types.BinaryNetworkEndpointLimits, types.BinaryNetworkEndpointLimits, OCR3_1SerializedLengthLimits, error) {
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

	maxLenStateTransitionOutputs := add(mul(2, pluginLimits.MaxKeyValueModifiedKeysPlusValuesLength), overhead)
	maxLenCertifiedPrepareOrCommit := add(mul(ed25519.SignatureSize+sigOverhead, cfg.ByzQuorumSize()),
		len(protocol.StateTransitionInputsDigest{}),
		maxLenStateTransitionOutputs,
		len(protocol.StateRootDigest{}),
		pluginLimits.MaxReportsPlusPrecursorLength,
		overhead)
	maxLenCertifiedCommittedReports := add(mul(ed25519.SignatureSize+sigOverhead, cfg.ByzQuorumSize()),
		len(protocol.StateTransitionInputsDigest{}),
		len(protocol.StateTransitionOutputDigest{}),
		len(protocol.StateRootDigest{}),
		pluginLimits.MaxReportsPlusPrecursorLength,
		overhead)
	maxLenMsgNewEpoch := overhead
	maxLenMsgEpochStartRequest := add(maxLenCertifiedPrepareOrCommit, overhead)
	maxLenMsgEpochStart := add(maxLenCertifiedPrepareOrCommit, mul(ed25519.SignatureSize+sigOverhead, cfg.ByzQuorumSize()), overhead)
	maxLenMsgRoundStart := add(pluginLimits.MaxQueryLength, overhead)
	maxLenMsgObservation := add(pluginLimits.MaxObservationLength, overhead)
	maxLenMsgProposal := add(mul(add(pluginLimits.MaxObservationLength, ed25519.SignatureSize+sigOverhead), cfg.N()), overhead)
	maxLenMsgPrepare := overhead
	maxLenMsgCommit := overhead
	maxLenMsgReportSignatures := add(mul(add(maxSigLen, sigOverhead), pluginLimits.MaxReportCount), overhead)
	maxLenMsgCertifiedCommitRequest := overhead
	maxLenMsgCertifiedCommit := add(maxLenCertifiedCommittedReports, overhead)
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
		protocol.MaxTreeSyncChunkKeysPlusValuesLength,
		mul( // repeated overheads
			protocol.MaxTreeSyncChunkKeys,
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
	maxLenAttestedStateTransitionBlock := maxLenCertifiedPrepareOrCommit
	maxLenMsgBlockSyncResponse := add(mul(protocol.MaxBlocksPerBlockSyncResponse, maxLenAttestedStateTransitionBlock), overhead)

	// blob exchange messages
	const blobChunkDigestSize = len(protocol.BlobChunkDigest{})
	maxNumBlobChunks := (pluginLimits.MaxBlobPayloadLength + protocol.BlobChunkSize - 1) / protocol.BlobChunkSize
	maxLenMsgBlobOffer := add(mul(blobChunkDigestSize, maxNumBlobChunks), overhead)
	maxLenMsgBlobChunkRequest := add(blobChunkDigestSize, overhead)
	maxLenMsgBlobChunkResponse := add(blobChunkDigestSize, protocol.BlobChunkSize, overhead)
	maxLenMsgBlobOfferResponse := add(blobChunkDigestSize, ed25519.SignatureSize+sigOverhead, overhead)

	maxDefaultPriorityMessageSize := max(
		maxLenMsgNewEpoch,
		maxLenMsgEpochStartRequest,
		maxLenMsgEpochStart,
		maxLenMsgRoundStart,
		maxLenMsgObservation,
		maxLenMsgProposal,
		maxLenMsgPrepare,
		maxLenMsgCommit,
		maxLenMsgReportSignatures,
		maxLenMsgCertifiedCommitRequest,
		maxLenMsgCertifiedCommit,
		maxLenMsgBlobOffer,
		maxLenMsgBlobChunkRequest,
	)

	maxLowPriorityMessageSize := max(
		maxLenMsgStateSyncSummary,
		maxLenMsgBlockSyncRequest,
		maxLenMsgTreeSyncChunkRequest,
	)

	minRoundInterval := math.Max(float64(cfg.DeltaRound), float64(cfg.DeltaGrace))

	minEpochInterval := math.Min(float64(cfg.DeltaProgress), math.Min(float64(cfg.DeltaInitial), float64(cfg.RMax)*float64(minRoundInterval)))

	defaultPriorityMessagesRate := (1.0*float64(time.Second)/float64(cfg.DeltaResend) +
		3.0*float64(time.Second)/minEpochInterval +
		8.0*float64(time.Second)/float64(minRoundInterval) +
		2.0*float64(time.Second)/float64(protocol.DeltaBlobOfferBroadcast) +
		1.0*float64(time.Second)/float64(protocol.DeltaBlobChunkRequest)) * 1.2

	lowPriorityMessagesRate := (1.0*float64(time.Second)/float64(protocol.DeltaMinBlockSyncRequest) +
		1.0*float64(time.Second)/float64(protocol.DeltaMinTreeSyncRequest) +
		1.0*float64(time.Second)/float64(protocol.DeltaStateSyncHeartbeat)) * 1.2

	defaultPriorityMessagesCapacity := mul(15, 3)
	lowPriorityMessagesCapacity := mul(3, 3)

	// we don't multiply bytesRate by a safetyMargin since we already have a generous overhead on each message

	defaultPriorityBytesRate := float64(time.Second)/float64(cfg.DeltaResend)*float64(maxLenMsgNewEpoch) +
		float64(time.Second)/float64(minEpochInterval)*float64(maxLenMsgNewEpoch) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgPrepare) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgCommit) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgReportSignatures) +
		float64(time.Second)/float64(minEpochInterval)*float64(maxLenMsgEpochStart) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgRoundStart) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgProposal) +
		float64(time.Second)/float64(minEpochInterval)*float64(maxLenMsgEpochStartRequest) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgObservation) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgCertifiedCommitRequest) +
		float64(time.Second)/float64(minRoundInterval)*float64(maxLenMsgCertifiedCommit) +
		float64(time.Second)/float64(protocol.DeltaBlobOfferBroadcast)*float64(maxLenMsgBlobOffer) + // blob-related messages
		float64(time.Second)/float64(protocol.DeltaBlobChunkRequest)*float64(maxLenMsgBlobChunkRequest)

	lowPriorityBytesRate := float64(time.Second)/float64(protocol.DeltaStateSyncHeartbeat)*float64(maxLenMsgStateSyncSummary) +
		float64(time.Second)/float64(protocol.DeltaMinBlockSyncRequest)*float64(maxLenMsgBlockSyncRequest) +
		float64(time.Second)/float64(protocol.DeltaMinTreeSyncRequest)*float64(maxLenMsgTreeSyncChunkRequest)

	defaultPriorityBytesCapacity := mul(add(
		maxLenMsgNewEpoch,
		maxLenMsgNewEpoch,
		maxLenMsgEpochStartRequest,
		maxLenMsgEpochStart,
		maxLenMsgRoundStart,
		maxLenMsgObservation,
		maxLenMsgProposal,
		maxLenMsgPrepare,
		maxLenMsgCommit,
		maxLenMsgReportSignatures,
		maxLenMsgCertifiedCommitRequest,
		maxLenMsgCertifiedCommit,
		maxLenMsgBlobOffer,
		maxLenMsgBlobChunkRequest,
		maxLenMsgBlobOfferResponse,
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
			maxLenMsgCertifiedCommitRequest,
			maxLenMsgCertifiedCommit,
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
