package limits

import (
	"crypto/ed25519"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type ocr3_1serializedLengthLimits struct {
	maxLenMsgNewEpoch               int
	maxLenMsgEpochStartRequest      int
	maxLenMsgEpochStart             int
	maxLenMsgRoundStart             int
	maxLenMsgObservation            int
	maxLenMsgProposal               int
	maxLenMsgPrepare                int
	maxLenMsgCommit                 int
	maxLenMsgReportSignatures       int
	maxLenMsgCertifiedCommitRequest int
	maxLenMsgCertifiedCommit        int
	maxLenMsgBlockSyncRequest       int
	maxLenMsgBlockSync              int
	maxLenMsgBlockSyncSummary       int
}

func ocr3_1limits(cfg ocr3config.PublicConfig, pluginLimits ocr3_1types.ReportingPluginLimits, maxSigLen int) (types.BinaryNetworkEndpointLimits, types.BinaryNetworkEndpointLimits, ocr3_1serializedLengthLimits, error) {
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

	const sigOverhead = 10
	const overhead = 256

	maxStateTransitionInputs := add(pluginLimits.MaxQueryLength, mul(pluginLimits.MaxObservationLength, cfg.N()), overhead)
	maxLenCertifiedPrepareOrCommit := add(mul(ed25519.SignatureSize+sigOverhead, cfg.ByzQuorumSize()), maxStateTransitionInputs, overhead)
	maxLenCertifiedCommittedReports := add(mul(ed25519.SignatureSize+sigOverhead, cfg.ByzQuorumSize()), pluginLimits.MaxReportsPlusPrecursorLength, overhead)
	maxStateTransitionBlock := add(maxStateTransitionInputs, overhead)
	maxAttestedStateTransitionBlock := add(mul(ed25519.SignatureSize+sigOverhead, cfg.ByzQuorumSize()), maxStateTransitionBlock, overhead)

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
	maxLenMsgBlockSyncRequest := overhead
	maxLenMsgBlockSync := add(protocol.MaxBlocksSent*maxAttestedStateTransitionBlock, overhead)
	maxLenMsgBlockSyncSummary := overhead

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
		maxLenMsgBlockSyncSummary,
	)

	maxLowPriorityMessageSize := max(
		maxLenMsgBlockSyncRequest,
		maxLenMsgBlockSync,
	)

	minEpochInterval := math.Min(float64(cfg.DeltaProgress), math.Min(float64(cfg.DeltaInitial), float64(cfg.RMax)*float64(cfg.DeltaRound)))

	defaultPriorityMessagesRate := (1.0*float64(time.Second)/float64(cfg.DeltaResend) +
		3.0*float64(time.Second)/minEpochInterval +
		8.0*float64(time.Second)/float64(cfg.DeltaRound) +
		1.0*float64(time.Second)/float64(protocol.DeltaBlockSyncHeartbeat)) * 1.2

	lowPriorityMessagesRate :=
		1.0 * float64(time.Second) / float64(protocol.DeltaMinBlockSyncRequest) * 1.2

	defaultPriorityMessagesCapacity := mul(13, 3)
	lowPriorityMessagesCapacity := mul(1, 3)

	// we don't multiply bytesRate by a safetyMargin since we already have a generous overhead on each message

	defaultPriorityBytesRate := float64(time.Second)/float64(cfg.DeltaResend)*float64(maxLenMsgNewEpoch) +
		float64(time.Second)/float64(minEpochInterval)*float64(maxLenMsgNewEpoch) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenMsgPrepare) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenMsgCommit) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenMsgReportSignatures) +
		float64(time.Second)/float64(minEpochInterval)*float64(maxLenMsgEpochStart) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenMsgRoundStart) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenMsgProposal) +
		float64(time.Second)/float64(minEpochInterval)*float64(maxLenMsgEpochStartRequest) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenMsgObservation) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenMsgCertifiedCommitRequest) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenMsgCertifiedCommit) +
		float64(time.Second)/float64(protocol.DeltaBlockSyncHeartbeat)*float64(maxLenMsgBlockSyncSummary)

	lowPriorityBytesRate := float64(time.Second) / float64(protocol.DeltaMinBlockSyncRequest) * float64(maxLenMsgBlockSyncRequest)

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
		maxLenMsgBlockSyncSummary,
	), 3)

	lowPriorityBytesCapacity := mul(add(
		maxLenMsgBlockSyncRequest,
	), 3)

	if overflow {
		// this should not happen due to us checking the limits in types.go
		return types.BinaryNetworkEndpointLimits{}, types.BinaryNetworkEndpointLimits{}, ocr3_1serializedLengthLimits{}, fmt.Errorf("int32 overflow while computing bandwidth limits")
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
		ocr3_1serializedLengthLimits{
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
			maxLenMsgBlockSyncRequest,
			maxLenMsgBlockSync,
			maxLenMsgBlockSyncSummary,
		},
		nil
}

func OCR3_1Limits(
	cfg ocr3config.PublicConfig,
	pluginLimits ocr3_1types.ReportingPluginLimits,
	maxSigLen int,
) (
	defaultLimits types.BinaryNetworkEndpointLimits,
	lowPriorityLimits types.BinaryNetworkEndpointLimits,
	err error,
) {
	defaultLimits, lowPriorityLimits, _, err = ocr3_1limits(cfg, pluginLimits, maxSigLen)
	return defaultLimits, lowPriorityLimits, err
}
