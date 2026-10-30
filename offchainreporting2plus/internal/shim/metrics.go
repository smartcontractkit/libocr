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
	registerer                   prometheus.Registerer
	oldestOpenTransactionAge     prometheus.GaugeFunc
	openedTransactionsTotal      prometheus.Counter
	discardedTransactionsTotal   prometheus.Counter
	refusedInsidePluginCallTotal prometheus.Counter
	refusedUnavailableTotal      prometheus.Counter
	transactionOpeningDuration   prometheus.Histogram
	transactionAgeAtDiscard      prometheus.Histogram
	transactionReads             prometheus.Histogram
	transactionReadBytes         prometheus.Histogram
}

func newReadOnlyKeyValueStateMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	oldestOpenTransactionAgeSeconds func() float64,
) *readOnlyKeyValueStateMetrics {
	// The number that predicts trouble. A read transaction keeps the database
	// from reclaiming the versions of the data it can still see, and that cost
	// is driven by how long the oldest one has been open rather than by how many
	// are open.
	oldestOpenTransactionAge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "ocr3_1_read_only_key_value_state_oldest_open_transaction_age_seconds",
		Help: "Seconds since creation of the oldest KeyValueStateReadTransaction returned by " +
			"ReadOnlyKeyValueState.NewReadTransaction that has not been discarded. Zero if none.",
	}, oldestOpenTransactionAgeSeconds)
	metricshelper.RegisterOrLogError(logger, registerer, oldestOpenTransactionAge, "ocr3_1_read_only_key_value_state_oldest_open_transaction_age_seconds")

	openedTransactionsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_read_only_key_value_state_opened_transactions_total",
		Help: "The number of successful ReadOnlyKeyValueState.NewReadTransaction calls. Subtract " +
			"ocr3_1_read_only_key_value_state_discarded_transactions_total for matching label sets to obtain the number currently open.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, openedTransactionsTotal, "ocr3_1_read_only_key_value_state_opened_transactions_total")

	discardedTransactionsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_1_read_only_key_value_state_discarded_transactions_total",
		Help: "The number of transactions returned by ReadOnlyKeyValueState.NewReadTransaction that have been " +
			"discarded through KeyValueStateReadTransaction.Discard. " +
			"Repeated Discard calls count once.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, discardedTransactionsTotal, "ocr3_1_read_only_key_value_state_discarded_transactions_total")

	// The vec itself is deliberately not kept: the per-reason counters are all
	// the callers need, and unregistering any one of them unregisters the vec,
	// see Close.
	refusedTransactionsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ocr3_1_read_only_key_value_state_refused_transactions_total",
		Help: "The number of ReadOnlyKeyValueState.NewReadTransaction calls returning an error, labeled by reason: " +
			"inside_plugin_call for use from a prohibited ReportingPlugin method, " +
			"unavailable for unavailable state or failure to create a transaction.",
	}, []string{"reason"})
	metricshelper.RegisterOrLogError(logger, registerer, refusedTransactionsTotal, "ocr3_1_read_only_key_value_state_refused_transactions_total")

	transactionOpeningDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "ocr3_1_experimental_read_only_key_value_state_transaction_opening_duration_seconds",
		Help: "The duration in seconds of successful ReadOnlyKeyValueState.NewReadTransaction calls. " +
			"Calls returning an error are excluded.",
		// Explicit boundaries keep integer-second le labels exact.
		Buckets: []float64{1e-6, 1e-5, 1e-4, 1e-3, 1e-2, 1e-1, 1, 2, 5, 10},
	})
	metricshelper.RegisterOrLogError(logger, registerer, transactionOpeningDuration, "ocr3_1_experimental_read_only_key_value_state_transaction_opening_duration_seconds")

	transactionAgeAtDiscard := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "ocr3_1_experimental_read_only_key_value_state_transaction_age_at_discard_seconds",
		Help: "The age in seconds of each KeyValueStateReadTransaction returned by ReadOnlyKeyValueState.NewReadTransaction " +
			"when discarded. Observed once per transaction on KeyValueStateReadTransaction.Discard.",
		Buckets: []float64{1e-4, 1e-3, 1e-2, 1e-1, 1, 2, 5, 10, 100, 1000},
	})
	metricshelper.RegisterOrLogError(logger, registerer, transactionAgeAtDiscard, "ocr3_1_experimental_read_only_key_value_state_transaction_age_at_discard_seconds")

	transactionReads := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "ocr3_1_experimental_read_only_key_value_state_transaction_reads",
		Help: "The number of Read calls on each KeyValueStateReadTransaction returned by ReadOnlyKeyValueState.NewReadTransaction " +
			"before discard, including errors and missing keys. Observed once on KeyValueStateReadTransaction.Discard. " +
			"Calls after discard are excluded.",
		Buckets: append([]float64{0}, prometheus.ExponentialBuckets(1, 10, 6)...),
	})
	metricshelper.RegisterOrLogError(logger, registerer, transactionReads, "ocr3_1_experimental_read_only_key_value_state_transaction_reads")

	transactionReadBytes := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "ocr3_1_experimental_read_only_key_value_state_transaction_read_bytes",
		Help: "The sum of value lengths in bytes returned without error by Read on each KeyValueStateReadTransaction " +
			"returned by ReadOnlyKeyValueState.NewReadTransaction. Repeated reads count separately. Observed once on " +
			"KeyValueStateReadTransaction.Discard.",
		Buckets: []float64{0, 100, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9},
	})
	metricshelper.RegisterOrLogError(logger, registerer, transactionReadBytes, "ocr3_1_experimental_read_only_key_value_state_transaction_read_bytes")

	return &readOnlyKeyValueStateMetrics{
		registerer,
		oldestOpenTransactionAge,
		openedTransactionsTotal,
		discardedTransactionsTotal,
		refusedTransactionsTotal.WithLabelValues("inside_plugin_call"),
		refusedTransactionsTotal.WithLabelValues("unavailable"),
		transactionOpeningDuration,
		transactionAgeAtDiscard,
		transactionReads,
		transactionReadBytes,
	}
}

func (m *readOnlyKeyValueStateMetrics) Close() {
	m.registerer.Unregister(m.oldestOpenTransactionAge)
	m.registerer.Unregister(m.openedTransactionsTotal)
	m.registerer.Unregister(m.discardedTransactionsTotal)
	m.registerer.Unregister(m.transactionOpeningDuration)
	m.registerer.Unregister(m.transactionAgeAtDiscard)
	m.registerer.Unregister(m.transactionReads)
	m.registerer.Unregister(m.transactionReadBytes)
	// A child of a vec describes the vec's descriptor, and prometheus.Registerer
	// documents that Unregister matches collectors by the descriptors they
	// describe. So unregistering the first child unregisters the whole vec, and
	// the remaining ones are no-ops. Unregistering all of them keeps this
	// method oblivious to which is which.
	m.registerer.Unregister(m.refusedInsidePluginCallTotal)
	m.registerer.Unregister(m.refusedUnavailableTotal)
}
