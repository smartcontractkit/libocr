package ocr2config

import (
	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/metricshelper"
	"github.com/prometheus/client_golang/prometheus"
)

type PublicConfigMetrics struct {
	registerer prometheus.Registerer

	deltaProgress prometheus.Gauge
	deltaResend   prometheus.Gauge
	deltaRound    prometheus.Gauge
	deltaGrace    prometheus.Gauge
	deltaStage    prometheus.Gauge
	rMax          prometheus.Gauge

	maxDurationInitialization               prometheus.Gauge
	maxDurationQuery                        prometheus.Gauge
	maxDurationObservation                  prometheus.Gauge
	maxDurationReport                       prometheus.Gauge
	maxDurationShouldAcceptFinalizedReport  prometheus.Gauge
	maxDurationShouldTransmitAcceptedReport prometheus.Gauge

	n prometheus.Gauge
	f prometheus.Gauge
}

func NewPublicConfigMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	publicConfig PublicConfig,
) *PublicConfigMetrics {
	deltaProgress := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_delta_progress_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	deltaProgress.Set(publicConfig.DeltaProgress.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaProgress, "ocr2_config_delta_progress_seconds")

	deltaResend := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_delta_resend_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	deltaResend.Set(publicConfig.DeltaResend.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaResend, "ocr2_config_delta_resend_seconds")

	deltaRound := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_delta_round_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	deltaRound.Set(publicConfig.DeltaRound.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaRound, "ocr2_config_delta_round_seconds")

	deltaGrace := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_delta_grace_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	deltaGrace.Set(publicConfig.DeltaGrace.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaGrace, "ocr2_config_delta_grace_seconds")

	deltaStage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_delta_stage_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	deltaStage.Set(publicConfig.DeltaStage.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaStage, "ocr2_config_delta_stage_seconds")

	rMax := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_r_max",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	rMax.Set(float64(publicConfig.RMax))
	metricshelper.RegisterOrLogError(logger, registerer, rMax, "ocr2_config_r_max")

	maxDurationInitialization := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_max_duration_initialization_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	maxDurationInitialization.Set(0)
	if publicConfig.MaxDurationInitialization != nil {
		maxDurationInitialization.Set(publicConfig.MaxDurationInitialization.Seconds())
	}
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationInitialization, "ocr2_config_max_duration_initialization_seconds")

	maxDurationQuery := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_max_duration_query_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	maxDurationQuery.Set(publicConfig.MaxDurationQuery.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationQuery, "ocr2_config_max_duration_query_seconds")

	maxDurationObservation := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_max_duration_observation_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	maxDurationObservation.Set(publicConfig.MaxDurationObservation.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationObservation, "ocr2_config_max_duration_observation_seconds")

	maxDurationReport := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_max_duration_report_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	maxDurationReport.Set(publicConfig.MaxDurationReport.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationReport, "ocr2_config_max_duration_report_seconds")

	maxDurationShouldAcceptFinalizedReport := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_max_duration_should_accept_finalize_report_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	maxDurationShouldAcceptFinalizedReport.Set(publicConfig.MaxDurationShouldAcceptFinalizedReport.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationShouldAcceptFinalizedReport, "ocr2_config_max_duration_should_accept_finalize_report_seconds")

	maxDurationShouldTransmitAcceptedReport := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_max_duration_should_transmit_accepted_report_seconds",
		Help: "See https://pkg.go.dev/github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr2config#PublicConfig for details",
	})
	maxDurationShouldTransmitAcceptedReport.Set(publicConfig.MaxDurationShouldTransmitAcceptedReport.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationShouldTransmitAcceptedReport, "ocr2_config_max_duration_should_transmit_accepted_report_seconds")

	n := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_n",
		Help: "The number of oracles participating in this protocol instance",
	})
	n.Set(float64(publicConfig.N()))
	metricshelper.RegisterOrLogError(logger, registerer, n, "ocr2_config_n")

	f := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_config_f",
		Help: "The maximum number of oracles that are assumed to be faulty while the protocol can retain liveness and safety",
	})
	f.Set(float64(publicConfig.F))
	metricshelper.RegisterOrLogError(logger, registerer, f, "ocr2_config_f")

	return &PublicConfigMetrics{
		registerer,

		deltaProgress,
		deltaResend,
		deltaRound,
		deltaGrace,
		deltaStage,
		rMax,

		maxDurationInitialization,
		maxDurationQuery,
		maxDurationObservation,
		maxDurationReport,
		maxDurationShouldAcceptFinalizedReport,
		maxDurationShouldTransmitAcceptedReport,

		n,
		f,
	}
}

func (p *PublicConfigMetrics) Close() {
	p.registerer.Unregister(p.deltaProgress)
	p.registerer.Unregister(p.deltaResend)
	p.registerer.Unregister(p.deltaRound)
	p.registerer.Unregister(p.deltaGrace)
	p.registerer.Unregister(p.deltaStage)
	p.registerer.Unregister(p.rMax)

	p.registerer.Unregister(p.maxDurationInitialization)
	p.registerer.Unregister(p.maxDurationQuery)
	p.registerer.Unregister(p.maxDurationObservation)
	p.registerer.Unregister(p.maxDurationReport)
	p.registerer.Unregister(p.maxDurationShouldAcceptFinalizedReport)
	p.registerer.Unregister(p.maxDurationShouldTransmitAcceptedReport)

	p.registerer.Unregister(p.n)
	p.registerer.Unregister(p.f)
}
