package shim

import (
	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/metricshelper"
	"github.com/prometheus/client_golang/prometheus"
)

type serializingEndpointMetrics struct {
	registerer           prometheus.Registerer
	sentMessagesTotal    prometheus.Counter
	droppedMessagesTotal prometheus.Counter
}

func newSerializingEndpointMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
) *serializingEndpointMetrics {
	sentMessagesTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr_telemetry_sent_messages_total",
		Help: "The number of telemetry messages sent.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, sentMessagesTotal, "ocr_telemetry_sent_messages_total")

	droppedMessagesTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr_telemetry_dropped_messages_total",
		Help: "The number of telemetry messages dropped.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, droppedMessagesTotal, "ocr_telemetry_dropped_messages_total")

	return &serializingEndpointMetrics{
		registerer,
		sentMessagesTotal,
		droppedMessagesTotal,
	}
}

func (m *serializingEndpointMetrics) Close() {
	m.registerer.Unregister(m.sentMessagesTotal)
	m.registerer.Unregister(m.droppedMessagesTotal)
}

type keyValueMetrics struct {
	registerer                       prometheus.Registerer
	closeWriteSetDurationNanoseconds prometheus.Histogram
	txWriteDurationNanoseconds       prometheus.Histogram
	txCommitDurationNanoseconds      prometheus.Histogram
}

const (
	hist_bucket_start  = 2
	hist_bucket_factor = 2
	hist_bucket_count  = 35
)

func newKeyValueMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
) *keyValueMetrics {
	closeWriteSetDurationNanoseconds := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ocr3_1_experimental_key_value_close_write_set_duration_ns",
		Help:    "How long it takes to close the write set.",
		Buckets: prometheus.ExponentialBuckets(hist_bucket_start, hist_bucket_factor, hist_bucket_count),
	})
	metricshelper.RegisterOrLogError(logger, registerer, closeWriteSetDurationNanoseconds, "ocr3_1_experimental_key_value_close_write_set_duration_ns")

	txWriteDurationNanoseconds := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ocr3_1_experimental_key_value_tx_write_duration_ns",
		Help:    "How long it takes to write to the transaction.",
		Buckets: prometheus.ExponentialBuckets(hist_bucket_start, hist_bucket_factor, hist_bucket_count),
	})
	metricshelper.RegisterOrLogError(logger, registerer, txWriteDurationNanoseconds, "ocr3_1_experimental_key_value_tx_write_duration_ns")

	txCommitDurationNanoseconds := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ocr3_1_experimental_key_value_tx_commit_duration_nanoseconds",
		Help:    "How long it takes to commit a transaction.",
		Buckets: prometheus.ExponentialBuckets(hist_bucket_start, hist_bucket_factor, hist_bucket_count),
	})
	metricshelper.RegisterOrLogError(logger, registerer, txCommitDurationNanoseconds, "ocr3_1_experimental_key_value_tx_commit_duration_ns")
	return &keyValueMetrics{
		registerer,
		closeWriteSetDurationNanoseconds,
		txWriteDurationNanoseconds,
		txCommitDurationNanoseconds,
	}
}

func (m *keyValueMetrics) Close() {
	m.registerer.Unregister(m.closeWriteSetDurationNanoseconds)
	m.registerer.Unregister(m.txWriteDurationNanoseconds)
	m.registerer.Unregister(m.txCommitDurationNanoseconds)
}
