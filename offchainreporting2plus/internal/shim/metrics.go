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

type readOnlyKeyValueStateMetrics struct {
	registerer                        prometheus.Registerer
	openReadTransactions              prometheus.GaugeFunc
	oldestOpenReadTransactionDuration prometheus.GaugeFunc
	openedReadTransactionsTotal       prometheus.Counter
	refusedInsidePluginCallTotal      prometheus.Counter
	refusedTooManyTotal               prometheus.Counter
	refusedUnavailableTotal           prometheus.Counter
}

func newReadOnlyKeyValueStateMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	openReadTransactionCount func() float64,
	oldestOpenReadTransactionDurationSeconds func() float64,
) *readOnlyKeyValueStateMetrics {
	openReadTransactions := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "ocr3_1_experimental_key_value_state_open_read_transactions",
		Help: "The number of read transactions the reporting plugin currently holds open.",
	}, openReadTransactionCount)
	metricshelper.RegisterOrLogError(logger, registerer, openReadTransactions, "ocr3_1_experimental_key_value_state_open_read_transactions")

	// The number that predicts trouble. A read transaction keeps the database
	// from reclaiming the versions of the data it can still see, and that cost
	// is driven by how long the oldest one has been open rather than by how many
	// are open.
	oldestOpenReadTransactionDuration := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "ocr3_1_experimental_key_value_state_oldest_open_read_transaction_duration_seconds",
		Help: "How long the oldest read transaction the reporting plugin holds open has been open, zero if none.",
	}, oldestOpenReadTransactionDurationSeconds)
	metricshelper.RegisterOrLogError(logger, registerer, oldestOpenReadTransactionDuration, "ocr3_1_experimental_key_value_state_oldest_open_read_transaction_duration_seconds")

	openedReadTransactionsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_key_value_state_opened_read_transactions_total",
		Help: "The number of read transactions opened by the reporting plugin.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, openedReadTransactionsTotal, "ocr3_1_experimental_key_value_state_opened_read_transactions_total")

	// The vec itself is deliberately not kept: the per-reason counters are all
	// the callers need, and unregistering any one of them unregisters the vec,
	// see Close.
	refusedReadTransactionsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ocr3_1_experimental_key_value_state_refused_read_transactions_total",
		Help: "The number of read transactions refused to the reporting plugin, by reason.",
	}, []string{"reason"})
	metricshelper.RegisterOrLogError(logger, registerer, refusedReadTransactionsTotal, "ocr3_1_experimental_key_value_state_refused_read_transactions_total")

	return &readOnlyKeyValueStateMetrics{
		registerer,
		openReadTransactions,
		oldestOpenReadTransactionDuration,
		openedReadTransactionsTotal,
		refusedReadTransactionsTotal.WithLabelValues("inside_plugin_call"),
		refusedReadTransactionsTotal.WithLabelValues("too_many"),
		refusedReadTransactionsTotal.WithLabelValues("unavailable"),
	}
}

func (m *readOnlyKeyValueStateMetrics) Close() {
	m.registerer.Unregister(m.openReadTransactions)
	m.registerer.Unregister(m.oldestOpenReadTransactionDuration)
	m.registerer.Unregister(m.openedReadTransactionsTotal)
	// A child of a vec describes the vec's descriptor, and prometheus.Registerer
	// documents that Unregister matches collectors by the descriptors they
	// describe. So unregistering the first child unregisters the whole vec, and
	// the remaining ones are no-ops. Unregistering all of them keeps this
	// method oblivious to which is which.
	m.registerer.Unregister(m.refusedInsidePluginCallTotal)
	m.registerer.Unregister(m.refusedTooManyTotal)
	m.registerer.Unregister(m.refusedUnavailableTotal)
}
