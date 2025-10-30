package shim

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
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

type semanticKeyValueDatabaseMetrics struct {
	registerer                   prometheus.Registerer
	closeWriteSetDurationSeconds prometheus.Histogram
	txWriteDurationSeconds       prometheus.Histogram
	txCommitDurationSeconds      prometheus.Histogram
}

func newSemanticKeyValueDatabaseMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
) *semanticKeyValueDatabaseMetrics {
	closeWriteSetDurationSeconds := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "ocr3_1_experimental_semantic_kvdb_close_write_set_duration_seconds",
		Help: "How long it takes to close the write set.",
		Buckets: prometheus.ExponentialBucketsRange(
			(50 * time.Microsecond).Seconds(),
			(10 * time.Second).Seconds(),
			20,
		),
	})
	metricshelper.RegisterOrLogError(logger, registerer, closeWriteSetDurationSeconds, "ocr3_1_experimental_semantic_kvdb_close_write_set_duration_seconds")

	txWriteDurationSeconds := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "ocr3_1_experimental_semantic_kvdb_tx_write_duration_seconds",
		Help: "How long it takes to write to the transaction.",
		Buckets: prometheus.ExponentialBucketsRange(
			(1 * time.Microsecond).Seconds(),
			(10 * time.Millisecond).Seconds(),
			10,
		),
	})
	metricshelper.RegisterOrLogError(logger, registerer, txWriteDurationSeconds, "ocr3_1_experimental_semantic_kvdb_tx_write_duration_seconds")

	txCommitDurationSeconds := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "ocr3_1_experimental_semantic_kvdb_tx_commit_duration_seconds",
		Help: "How long it takes to commit a transaction.",
		Buckets: prometheus.ExponentialBucketsRange(
			(500 * time.Microsecond).Seconds(),
			(10 * time.Second).Seconds(),
			20,
		),
	})
	metricshelper.RegisterOrLogError(logger, registerer, txCommitDurationSeconds, "ocr3_1_experimental_semantic_kvdb_tx_commit_duration_seconds")
	return &semanticKeyValueDatabaseMetrics{
		registerer,
		closeWriteSetDurationSeconds,
		txWriteDurationSeconds,
		txCommitDurationSeconds,
	}
}

func (m *semanticKeyValueDatabaseMetrics) Close() {
	m.registerer.Unregister(m.closeWriteSetDurationSeconds)
	m.registerer.Unregister(m.txWriteDurationSeconds)
	m.registerer.Unregister(m.txCommitDurationSeconds)
}

type keyValueDatabaseMetrics struct {
	registerer                          prometheus.Registerer
	openedReadTransactionsTotal         prometheus.Counter
	discardedReadTransactionsTotal      prometheus.Counter
	openedReadWriteTransactionsTotal    prometheus.Counter
	committedReadWriteTransactionsTotal prometheus.Counter
	discardedReadWriteTransactionsTotal prometheus.Counter
}

func newKeyValueDatabaseMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
) *keyValueDatabaseMetrics {
	openedReadTransactionsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_kvdb_opened_read_transactions_total",
		Help: "The number of opened read transactions.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, openedReadTransactionsTotal, "ocr3_1_experimental_kvdb_opened_read_transactions_total")

	discardedReadTransactionsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_kvdb_discarded_read_transactions_total",
		Help: "The number of discarded read transactions.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, discardedReadTransactionsTotal, "ocr3_1_experimental_kvdb_discarded_read_transactions_total")

	openedReadWriteTransactionsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_kvdb_opened_read_write_transactions_total",
		Help: "The number of opened read-write transactions.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, openedReadWriteTransactionsTotal, "ocr3_1_experimental_kvdb_opened_read_write_transactions_total")

	committedReadWriteTransactionsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_kvdb_committed_read_write_transactions_total",
		Help: "The number of committed read-write transactions.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, committedReadWriteTransactionsTotal, "ocr3_1_experimental_kvdb_committed_read_write_transactions_total")

	discardedReadWriteTransactionsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_kvdb_discarded_read_write_transactions_total",
		Help: "The number of discarded read-write transactions.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, discardedReadWriteTransactionsTotal, "ocr3_1_experimental_kvdb_discarded_read_write_transactions_total")

	return &keyValueDatabaseMetrics{
		registerer,
		openedReadTransactionsTotal,
		discardedReadTransactionsTotal,
		openedReadWriteTransactionsTotal,
		committedReadWriteTransactionsTotal,
		discardedReadWriteTransactionsTotal,
	}
}

func (m *keyValueDatabaseMetrics) Close() {
	m.registerer.Unregister(m.openedReadTransactionsTotal)
	m.registerer.Unregister(m.discardedReadTransactionsTotal)
	m.registerer.Unregister(m.openedReadWriteTransactionsTotal)
	m.registerer.Unregister(m.committedReadWriteTransactionsTotal)
	m.registerer.Unregister(m.discardedReadWriteTransactionsTotal)
}
