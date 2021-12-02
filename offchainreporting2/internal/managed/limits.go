package managed

import (
	"crypto/ed25519"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

func limits(cfg config.PublicConfig, reportingPluginInfo types.ReportingPluginInfo, maxSigLen int) types.BinaryNetworkEndpointLimits {
	const overhead = 256

	maxLenNewEpoch := overhead
	maxLenObserveReq := reportingPluginInfo.MaxQueryLen + overhead
	maxLenObserve := reportingPluginInfo.MaxObservationLen + overhead
	maxLenReportReq := (reportingPluginInfo.MaxObservationLen+ed25519.SignatureSize)*cfg.N() + overhead
	maxLenReport := reportingPluginInfo.MaxReportLen + ed25519.SignatureSize + overhead
	maxLenFinal := reportingPluginInfo.MaxReportLen + maxSigLen*cfg.N() + overhead
	maxLenFinalEcho := maxLenFinal

	maxMessageSize := max(maxLenObserveReq, maxLenObserve, maxLenReportReq, maxLenReport, maxLenFinal, maxLenFinalEcho)

	messagesRate := (1.0*float64(time.Second)/float64(cfg.DeltaResend) +
		1.0*float64(time.Second)/float64(cfg.DeltaProgress) +
		1.0*float64(time.Second)/float64(cfg.DeltaRound) +
		3.0*float64(time.Second)/float64(cfg.DeltaRound) +
		2.0*float64(time.Second)/float64(cfg.DeltaRound)) * 2.0

	messagesCapacity := (2 + 6) * 2

	bytesRate := float64(time.Second)/float64(cfg.DeltaResend)*float64(maxLenNewEpoch) +
		float64(time.Second)/float64(cfg.DeltaProgress)*float64(maxLenNewEpoch) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenObserveReq) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenObserve) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenReportReq) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenReport) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenFinal) +
		float64(time.Second)/float64(cfg.DeltaRound)*float64(maxLenFinalEcho)

	bytesCapacity := (maxLenNewEpoch + maxLenObserveReq + maxLenObserve + maxLenReportReq + maxLenReport + maxLenFinal + maxLenFinalEcho) * 2

	return types.BinaryNetworkEndpointLimits{
		maxMessageSize,
		messagesRate,
		messagesCapacity,
		bytesRate,
		bytesCapacity,
	}
}
