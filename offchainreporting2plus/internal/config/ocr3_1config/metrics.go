package ocr3_1config

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
)

type PublicConfigMetrics struct {
	registerer                       prometheus.Registerer
	deltaProgress                    prometheus.Gauge
	deltaResend                      prometheus.Gauge
	deltaInitial                     prometheus.Gauge
	deltaRound                       prometheus.Gauge
	deltaGrace                       prometheus.Gauge
	deltaReportsPlusPrecursorRequest prometheus.Gauge
	deltaStage                       prometheus.Gauge

	// state sync
	deltaStateSyncSummaryInterval prometheus.Gauge

	// block sync
	deltaBlockSyncMinRequestToSameOracleInterval prometheus.Gauge
	deltaBlockSyncResponseTimeout                prometheus.Gauge
	maxBlocksPerBlockSyncResponse                prometheus.Gauge
	maxParallelRequestedBlocks                   prometheus.Gauge

	// tree sync
	deltaTreeSyncMinRequestToSameOracleInterval prometheus.Gauge
	deltaTreeSyncResponseTimeout                prometheus.Gauge
	maxTreeSyncChunkKeys                        prometheus.Gauge
	maxTreeSyncChunkKeysPlusValuesBytes         prometheus.Gauge
	maxParallelTreeSyncChunkFetches             prometheus.Gauge

	// snapshotting
	snapshotInterval               prometheus.Gauge
	maxHistoricalSnapshotsRetained prometheus.Gauge

	// blobs
	deltaBlobOfferMinRequestToSameOracleInterval prometheus.Gauge
	deltaBlobOfferResponseTimeout                prometheus.Gauge
	deltaBlobBroadcastGrace                      prometheus.Gauge
	deltaBlobChunkMinRequestToSameOracleInterval prometheus.Gauge
	deltaBlobChunkResponseTimeout                prometheus.Gauge
	blobChunkBytes                               prometheus.Gauge

	rMax prometheus.Gauge
	// skip S

	maxDurationInitialization               prometheus.Gauge
	warnDurationQuery                       prometheus.Gauge
	warnDurationObservation                 prometheus.Gauge
	warnDurationValidateObservation         prometheus.Gauge
	warnDurationObservationQuorum           prometheus.Gauge
	warnDurationStateTransition             prometheus.Gauge
	warnDurationCommitted                   prometheus.Gauge
	maxDurationShouldAcceptAttestedReport   prometheus.Gauge
	maxDurationShouldTransmitAcceptedReport prometheus.Gauge

	prevSeqNr prometheus.Gauge

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
		Name: "ocr3_1_config_delta_progress_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaProgress.Set(publicConfig.DeltaProgress.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaProgress, "ocr3_1_config_delta_progress_seconds")

	deltaResend := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_resend_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaResend.Set(publicConfig.GetDeltaResend().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaResend, "ocr3_1_config_delta_resend_seconds")

	deltaInitial := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_initial_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaInitial.Set(publicConfig.GetDeltaInitial().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaInitial, "ocr3_1_config_delta_initial_seconds")

	deltaRound := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_round_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaRound.Set(publicConfig.DeltaRound.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaRound, "ocr3_1_config_delta_round_seconds")

	deltaGrace := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_grace_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaGrace.Set(publicConfig.DeltaGrace.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaGrace, "ocr3_1_config_delta_grace_seconds")

	deltaReportsPlusPrecursorRequest := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_reports_plus_precursor_request_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaReportsPlusPrecursorRequest.Set(publicConfig.GetDeltaReportsPlusPrecursorRequest().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaReportsPlusPrecursorRequest, "ocr3_1_config_delta_reports_plus_precursor_request_seconds")

	deltaStage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_stage_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaStage.Set(publicConfig.DeltaStage.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaStage, "ocr3_1_config_delta_stage_seconds")

	// state sync
	deltaStateSyncSummaryInterval := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_state_sync_summary_interval_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaStateSyncSummaryInterval.Set(publicConfig.GetDeltaStateSyncSummaryInterval().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaStateSyncSummaryInterval, "ocr3_1_config_delta_state_sync_summary_interval_seconds")

	// block sync
	deltaBlockSyncMinRequestToSameOracleInterval := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_block_sync_min_request_to_same_oracle_interval_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaBlockSyncMinRequestToSameOracleInterval.Set(publicConfig.GetDeltaBlockSyncMinRequestToSameOracleInterval().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaBlockSyncMinRequestToSameOracleInterval, "ocr3_1_config_delta_block_sync_min_request_to_same_oracle_interval_seconds")

	deltaBlockSyncResponseTimeout := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_block_sync_response_timeout_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaBlockSyncResponseTimeout.Set(publicConfig.GetDeltaBlockSyncResponseTimeout().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaBlockSyncResponseTimeout, "ocr3_1_config_delta_block_sync_response_timeout_seconds")

	maxBlocksPerBlockSyncResponse := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_max_blocks_per_block_sync_response",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	maxBlocksPerBlockSyncResponse.Set(float64(publicConfig.GetMaxBlocksPerBlockSyncResponse()))
	metricshelper.RegisterOrLogError(logger, registerer, maxBlocksPerBlockSyncResponse, "ocr3_1_config_max_blocks_per_block_sync_response")

	maxParallelRequestedBlocks := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_max_parallel_requested_blocks",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	maxParallelRequestedBlocks.Set(float64(publicConfig.GetMaxParallelRequestedBlocks()))
	metricshelper.RegisterOrLogError(logger, registerer, maxParallelRequestedBlocks, "ocr3_1_config_max_parallel_requested_blocks")

	// tree sync
	deltaTreeSyncMinRequestToSameOracleInterval := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_tree_sync_min_request_to_same_oracle_interval_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaTreeSyncMinRequestToSameOracleInterval.Set(publicConfig.GetDeltaTreeSyncMinRequestToSameOracleInterval().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaTreeSyncMinRequestToSameOracleInterval, "ocr3_1_config_delta_tree_sync_min_request_to_same_oracle_interval_seconds")

	deltaTreeSyncResponseTimeout := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_tree_sync_response_timeout_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaTreeSyncResponseTimeout.Set(publicConfig.GetDeltaTreeSyncResponseTimeout().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaTreeSyncResponseTimeout, "ocr3_1_config_delta_tree_sync_response_timeout_seconds")

	maxTreeSyncChunkKeys := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_max_tree_sync_chunk_keys",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	maxTreeSyncChunkKeys.Set(float64(publicConfig.GetMaxTreeSyncChunkKeys()))
	metricshelper.RegisterOrLogError(logger, registerer, maxTreeSyncChunkKeys, "ocr3_1_config_max_tree_sync_chunk_keys")

	maxTreeSyncChunkKeysPlusValuesBytes := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_max_tree_sync_chunk_keys_plus_values_bytes",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	maxTreeSyncChunkKeysPlusValuesBytes.Set(float64(publicConfig.GetMaxTreeSyncChunkKeysPlusValuesBytes()))
	metricshelper.RegisterOrLogError(logger, registerer, maxTreeSyncChunkKeysPlusValuesBytes, "ocr3_1_config_max_tree_sync_chunk_keys_plus_values_bytes")

	maxParallelTreeSyncChunkFetches := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_max_parallel_tree_sync_chunk_fetches",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	maxParallelTreeSyncChunkFetches.Set(float64(publicConfig.GetMaxParallelTreeSyncChunkFetches()))
	metricshelper.RegisterOrLogError(logger, registerer, maxParallelTreeSyncChunkFetches, "ocr3_1_config_max_parallel_tree_sync_chunk_fetches")

	// snapshotting
	snapshotInterval := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_snapshot_interval",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	snapshotInterval.Set(float64(publicConfig.GetSnapshotInterval()))
	metricshelper.RegisterOrLogError(logger, registerer, snapshotInterval, "ocr3_1_config_snapshot_interval")

	maxHistoricalSnapshotsRetained := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_max_historical_snapshots_retained",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	maxHistoricalSnapshotsRetained.Set(float64(publicConfig.GetMaxHistoricalSnapshotsRetained()))
	metricshelper.RegisterOrLogError(logger, registerer, maxHistoricalSnapshotsRetained, "ocr3_1_config_max_historical_snapshots_retained")

	// blobs
	deltaBlobOfferMinRequestToSameOracleInterval := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_blob_offer_min_request_to_same_oracle_interval_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaBlobOfferMinRequestToSameOracleInterval.Set(publicConfig.GetDeltaBlobOfferMinRequestToSameOracleInterval().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaBlobOfferMinRequestToSameOracleInterval, "ocr3_1_config_delta_blob_offer_min_request_to_same_oracle_interval_seconds")

	deltaBlobOfferResponseTimeout := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_blob_offer_response_timeout_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaBlobOfferResponseTimeout.Set(publicConfig.GetDeltaBlobOfferResponseTimeout().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaBlobOfferResponseTimeout, "ocr3_1_config_delta_blob_offer_response_timeout_seconds")

	deltaBlobBroadcastGrace := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_blob_broadcast_grace_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaBlobBroadcastGrace.Set(publicConfig.GetDeltaBlobBroadcastGrace().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaBlobBroadcastGrace, "ocr3_1_config_delta_blob_broadcast_grace_seconds")

	deltaBlobChunkMinRequestToSameOracleInterval := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_blob_chunk_min_request_to_same_oracle_interval_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaBlobChunkMinRequestToSameOracleInterval.Set(publicConfig.GetDeltaBlobChunkMinRequestToSameOracleInterval().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaBlobChunkMinRequestToSameOracleInterval, "ocr3_1_config_delta_blob_chunk_min_request_to_same_oracle_interval_seconds")

	deltaBlobChunkResponseTimeout := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_delta_blob_chunk_response_timeout_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	deltaBlobChunkResponseTimeout.Set(publicConfig.GetDeltaBlobChunkResponseTimeout().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, deltaBlobChunkResponseTimeout, "ocr3_1_config_delta_blob_chunk_response_timeout_seconds")

	blobChunkBytes := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_blob_chunk_bytes",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	blobChunkBytes.Set(float64(publicConfig.GetBlobChunkBytes()))
	metricshelper.RegisterOrLogError(logger, registerer, blobChunkBytes, "ocr3_1_config_blob_chunk_bytes")

	rMax := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_r_max",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	rMax.Set(float64(publicConfig.RMax))
	metricshelper.RegisterOrLogError(logger, registerer, rMax, "ocr3_1_config_r_max")

	maxDurationInitialization := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_max_duration_initialization_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	maxDurationInitialization.Set(publicConfig.MaxDurationInitialization.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationInitialization, "ocr3_1_config_max_duration_initialization_seconds")

	warnDurationQuery := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_warn_duration_query_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	warnDurationQuery.Set(publicConfig.WarnDurationQuery.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, warnDurationQuery, "ocr3_1_config_warn_duration_query_seconds")

	warnDurationObservation := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_warn_duration_observation_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	warnDurationObservation.Set(publicConfig.WarnDurationObservation.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, warnDurationObservation, "ocr3_1_config_warn_duration_observation_seconds")

	warnDurationValidateObservation := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_warn_duration_validate_observation_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	warnDurationValidateObservation.Set(publicConfig.WarnDurationValidateObservation.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, warnDurationValidateObservation, "ocr3_1_config_warn_duration_validate_observation_seconds")

	warnDurationObservationQuorum := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_warn_duration_observation_quorum_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	warnDurationObservationQuorum.Set(publicConfig.WarnDurationObservationQuorum.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, warnDurationObservationQuorum, "ocr3_1_config_warn_duration_observation_quorum_seconds")

	warnDurationStateTransition := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_warn_duration_state_transition_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	warnDurationStateTransition.Set(publicConfig.WarnDurationStateTransition.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, warnDurationStateTransition, "ocr3_1_config_warn_duration_state_transition_seconds")

	warnDurationCommitted := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_warn_duration_committed_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	warnDurationCommitted.Set(publicConfig.WarnDurationCommitted.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, warnDurationCommitted, "ocr3_1_config_warn_duration_committed_seconds")

	maxDurationShouldAcceptAttestedReport := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_max_duration_should_accept_attested_report_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	maxDurationShouldAcceptAttestedReport.Set(publicConfig.MaxDurationShouldAcceptAttestedReport.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationShouldAcceptAttestedReport, "ocr3_1_config_max_duration_should_accept_attested_report_seconds")

	maxDurationShouldTransmitAcceptedReport := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_max_duration_should_transmit_accepted_report_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	maxDurationShouldTransmitAcceptedReport.Set(publicConfig.MaxDurationShouldTransmitAcceptedReport.Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, maxDurationShouldTransmitAcceptedReport, "ocr3_1_config_max_duration_should_transmit_accepted_report_seconds")

	prevSeqNr := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_prev_seq_nr",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	if publicConfig.PrevSeqNr != nil {
		prevSeqNr.Set(float64(*publicConfig.PrevSeqNr))
	} else {
		prevSeqNr.Set(0)
	}
	metricshelper.RegisterOrLogError(logger, registerer, prevSeqNr, "ocr3_1_config_prev_seq_nr")

	n := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_n",
		Help: "The number of oracles participating in this protocol instance",
	})
	n.Set(float64(publicConfig.N()))
	metricshelper.RegisterOrLogError(logger, registerer, n, "ocr3_1_config_n")

	f := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_f",
		Help: "The maximum number of oracles that are assumed to be faulty while the protocol can retain liveness and safety",
	})
	f.Set(float64(publicConfig.F))
	metricshelper.RegisterOrLogError(logger, registerer, f, "ocr3_1_config_f")

	minRoundInterval := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_1_config_min_round_interval_seconds",
		Help: "See https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config#PublicConfig for details",
	})
	minRoundInterval.Set(publicConfig.MinRoundInterval().Seconds())
	metricshelper.RegisterOrLogError(logger, registerer, minRoundInterval, "ocr3_1_config_min_round_interval_seconds")

	return &PublicConfigMetrics{
		registerer,

		deltaProgress,
		deltaResend,
		deltaInitial,
		deltaRound,
		deltaGrace,
		deltaReportsPlusPrecursorRequest,
		deltaStage,

		// state sync
		deltaStateSyncSummaryInterval,

		// block sync
		deltaBlockSyncMinRequestToSameOracleInterval,
		deltaBlockSyncResponseTimeout,
		maxBlocksPerBlockSyncResponse,
		maxParallelRequestedBlocks,

		// tree sync
		deltaTreeSyncMinRequestToSameOracleInterval,
		deltaTreeSyncResponseTimeout,
		maxTreeSyncChunkKeys,
		maxTreeSyncChunkKeysPlusValuesBytes,
		maxParallelTreeSyncChunkFetches,

		// snapshotting
		snapshotInterval,
		maxHistoricalSnapshotsRetained,

		// blobs
		deltaBlobOfferMinRequestToSameOracleInterval,
		deltaBlobOfferResponseTimeout,
		deltaBlobBroadcastGrace,
		deltaBlobChunkMinRequestToSameOracleInterval,
		deltaBlobChunkResponseTimeout,
		blobChunkBytes,

		rMax,

		maxDurationInitialization,
		warnDurationQuery,
		warnDurationObservation,
		warnDurationValidateObservation,
		warnDurationObservationQuorum,
		warnDurationStateTransition,
		warnDurationCommitted,
		maxDurationShouldAcceptAttestedReport,
		maxDurationShouldTransmitAcceptedReport,

		prevSeqNr,

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
	pm.registerer.Unregister(pm.deltaReportsPlusPrecursorRequest)
	pm.registerer.Unregister(pm.deltaStage)

	// state sync
	pm.registerer.Unregister(pm.deltaStateSyncSummaryInterval)

	// block sync
	pm.registerer.Unregister(pm.deltaBlockSyncMinRequestToSameOracleInterval)
	pm.registerer.Unregister(pm.deltaBlockSyncResponseTimeout)
	pm.registerer.Unregister(pm.maxBlocksPerBlockSyncResponse)
	pm.registerer.Unregister(pm.maxParallelRequestedBlocks)

	// tree sync
	pm.registerer.Unregister(pm.deltaTreeSyncMinRequestToSameOracleInterval)
	pm.registerer.Unregister(pm.deltaTreeSyncResponseTimeout)
	pm.registerer.Unregister(pm.maxTreeSyncChunkKeys)
	pm.registerer.Unregister(pm.maxTreeSyncChunkKeysPlusValuesBytes)
	pm.registerer.Unregister(pm.maxParallelTreeSyncChunkFetches)

	// snapshotting
	pm.registerer.Unregister(pm.snapshotInterval)
	pm.registerer.Unregister(pm.maxHistoricalSnapshotsRetained)

	// blobs
	pm.registerer.Unregister(pm.deltaBlobOfferMinRequestToSameOracleInterval)
	pm.registerer.Unregister(pm.deltaBlobOfferResponseTimeout)
	pm.registerer.Unregister(pm.deltaBlobBroadcastGrace)
	pm.registerer.Unregister(pm.deltaBlobChunkMinRequestToSameOracleInterval)
	pm.registerer.Unregister(pm.deltaBlobChunkResponseTimeout)
	pm.registerer.Unregister(pm.blobChunkBytes)

	pm.registerer.Unregister(pm.rMax)

	pm.registerer.Unregister(pm.maxDurationInitialization)
	pm.registerer.Unregister(pm.warnDurationQuery)
	pm.registerer.Unregister(pm.warnDurationObservation)
	pm.registerer.Unregister(pm.warnDurationValidateObservation)
	pm.registerer.Unregister(pm.warnDurationObservationQuorum)
	pm.registerer.Unregister(pm.warnDurationStateTransition)
	pm.registerer.Unregister(pm.warnDurationCommitted)
	pm.registerer.Unregister(pm.maxDurationShouldAcceptAttestedReport)
	pm.registerer.Unregister(pm.maxDurationShouldTransmitAcceptedReport)

	pm.registerer.Unregister(pm.prevSeqNr)

	pm.registerer.Unregister(pm.n)
	pm.registerer.Unregister(pm.f)
	pm.registerer.Unregister(pm.minRoundInterval)
}
