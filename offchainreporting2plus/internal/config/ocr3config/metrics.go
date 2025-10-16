package ocr3config

import (
	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/metricshelper"
	"github.com/prometheus/client_golang/prometheus"
)

type PublicConfigMetrics struct {
	registerer                  prometheus.Registerer
	deltaProgress               prometheus.Gauge
	deltaResend                 prometheus.Gauge
	deltaInitial                prometheus.Gauge
	deltaRound                  prometheus.Gauge
	deltaGrace                  prometheus.Gauge
	deltaCertifiedCommitRequest prometheus.Gauge
	deltaStage                  prometheus.Gauge
	rMax                        prometheus.Gauge
	// skip S

	maxDurationInitialization               prometheus.Gauge
	maxDurationQuery                        prometheus.Gauge
	maxDurationObservation                  prometheus.Gauge
	maxDurationShouldAcceptAttestedReport   prometheus.Gauge
	maxDurationShouldTransmitAcceptedReport prometheus.Gauge

	n prometheus.Gauge
	f prometheus.Gauge

	minRoundInterval prometheus.Gauge
}

func NewPublicConfigMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	publicConfig PublicConfig,
) *PublicConfigMetrics {

	deltaProgress := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_delta_progress_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	deltaProgress.Set(publicConfig.DeltaProgress.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaProgress, "ocr3_config_delta_progress_seconds")

	deltaResend := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_delta_resend_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	deltaResend.Set(publicConfig.DeltaResend.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaResend, "ocr3_config_delta_resend_seconds")

	deltaInitial := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_delta_initial_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	deltaInitial.Set(publicConfig.DeltaInitial.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaInitial, "ocr3_config_delta_initial_seconds")

	deltaRound := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_delta_round_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	deltaRound.Set(publicConfig.DeltaRound.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaRound, "ocr3_config_delta_round_seconds")

	deltaGrace := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_delta_grace_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	deltaGrace.Set(publicConfig.DeltaGrace.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaGrace, "ocr3_config_delta_grace_seconds")

	deltaCertifiedCommitRequest := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_delta_certified_commit_request_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	deltaCertifiedCommitRequest.Set(publicConfig.DeltaCertifiedCommitRequest.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaCertifiedCommitRequest, "ocr3_config_delta_certified_commit_request_seconds")

	deltaStage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_delta_stage_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	deltaStage.Set(publicConfig.DeltaStage.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaStage, "ocr3_config_delta_stage_seconds")

	rMax := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_r_max",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	rMax.Set(float64(publicConfig.RMax))
	metricshelper.RegisterOrLogError(logger, registerer, rMax, "ocr3_config_r_max")

	maxDurationInitialization := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_max_duration_initialization_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	maxDurationInitialization.Set(0)
	if publicConfig.MaxDurationInitialization != nil {
		maxDurationInitialization.Set(publicConfig.MaxDurationInitialization.Seconds())
	}
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationInitialization, "ocr3_config_max_duration_initialization_seconds")

	maxDurationQuery := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_max_duration_query_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	maxDurationQuery.Set(publicConfig.MaxDurationQuery.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationQuery, "ocr3_config_max_duration_query_seconds")

	maxDurationObservation := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_max_duration_observation_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	maxDurationObservation.Set(publicConfig.MaxDurationObservation.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationObservation, "ocr3_config_max_duration_observation_seconds")

	maxDurationShouldAcceptAttestedReport := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_max_duration_should_accept_attested_report_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	maxDurationShouldAcceptAttestedReport.Set(publicConfig.MaxDurationShouldAcceptAttestedReport.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationShouldAcceptAttestedReport, "ocr3_config_max_duration_should_accept_attested_report_seconds")

	maxDurationShouldTransmitAcceptedReport := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_max_duration_should_transmit_accepted_report_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	maxDurationShouldTransmitAcceptedReport.Set(publicConfig.MaxDurationShouldTransmitAcceptedReport.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationShouldTransmitAcceptedReport, "ocr3_config_max_duration_should_transmit_accepted_report_seconds")

	n := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_n",
		Help: "The number of oracles participating in this protocol instance",
	})
	n.Set(float64(publicConfig.N()))
	metricshelper.RegisterOrLogError(logger, registerer, n, "ocr3_config_n")

	f := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_f",
		Help: "The maximum number of oracles that are assumed to be faulty while the protocol can retain liveness and safety",
	})
	f.Set(float64(publicConfig.F))
	metricshelper.RegisterOrLogError(logger, registerer, f, "ocr3_config_f")

	minRoundInterval := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_config_min_round_interval_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config#PublicConfig for details",
	})
	minRoundInterval.Set(publicConfig.MinRoundInterval().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, minRoundInterval, "ocr3_config_min_round_interval_seconds")

	return &PublicConfigMetrics{
		registerer,

		deltaProgress,
		deltaResend,
		deltaInitial,
		deltaRound,
		deltaGrace,
		deltaCertifiedCommitRequest,
		deltaStage,
		rMax,

		maxDurationInitialization,
		maxDurationQuery,
		maxDurationObservation,
		maxDurationShouldAcceptAttestedReport,
		maxDurationShouldTransmitAcceptedReport,

		n,
		f,
		minRoundInterval,
	}
}

func (pm *PublicConfigMetrics) Close() {
	pm.registerer.Unregister(pm.deltaProgress)
	pm.registerer.Unregister(pm.deltaResend)
	pm.registerer.Unregister(pm.deltaInitial)
	pm.registerer.Unregister(pm.deltaRound)
	pm.registerer.Unregister(pm.deltaGrace)
	pm.registerer.Unregister(pm.deltaCertifiedCommitRequest)
	pm.registerer.Unregister(pm.deltaStage)
	pm.registerer.Unregister(pm.rMax)

	pm.registerer.Unregister(pm.maxDurationInitialization)
	pm.registerer.Unregister(pm.maxDurationQuery)
	pm.registerer.Unregister(pm.maxDurationObservation)
	pm.registerer.Unregister(pm.maxDurationShouldAcceptAttestedReport)
	pm.registerer.Unregister(pm.maxDurationShouldTransmitAcceptedReport)

	pm.registerer.Unregister(pm.n)
	pm.registerer.Unregister(pm.f)
	pm.registerer.Unregister(pm.minRoundInterval)
}
