package shim

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/common/list"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

const (
	// A read transaction held for longer than this is almost certainly a bug in
	// the plugin: reads take microseconds, so anything this slow is blocking on
	// something else while keeping the transaction open. We only complain, we
	// never discard it out from under a running plugin.
	recommendedMaxKeyValueStateReadTransactionDuration = 10 * time.Second
	longHeldKeyValueStateReadTransactionCheckInterval  = 10 * time.Second
)

var (
	errKeyValueStateUnavailable              = errors.New("key value state is unavailable")
	errKeyValueStateReadInsidePluginCall     = errors.New("key value state must not be read via ReadOnlyKeyValueState inside this ReportingPlugin method; use the method's KeyValueStateReader or KeyValueStateReadWriter argument if available")
	errKeyValueStateReadTransactionDiscarded = errors.New("key value state read transaction has been discarded")
)

// ReadOnlyKeyValueStateWrapper implements [ocr3_1types.ReadOnlyKeyValueState],
// the plugin's read access to the KeyValueState from outside its own methods.
//
// It exists to enable deferred initialization: the reporting plugin is
// constructed before the key value database is opened (the database needs the
// limits that the plugin's factory reports), so the plugin is handed the
// wrapper up front and the database is supplied later through Enable.
type ReadOnlyKeyValueStateWrapper struct {
	logger  commontypes.Logger
	metrics *readOnlyKeyValueStateMetrics

	subs     subprocesses.Subprocesses
	chStop   chan struct{}
	stopOnce sync.Once

	mu                    sync.Mutex
	nilOrKeyValueDatabase protocol.KeyValueDatabase
	// openReadTransactions is ordered from oldest to newest: a transaction's
	// openedAt is assigned and the transaction registered while holding mu, so
	// each registration is newer than the previous one. The front is therefore
	// the oldest open transaction, in O(1).
	openReadTransactions *list.List[*keyValueStateReadTransaction]
}

var _ ocr3_1types.ReadOnlyKeyValueState = &ReadOnlyKeyValueStateWrapper{}

func NewReadOnlyKeyValueStateWrapper(logger commontypes.Logger, registerer prometheus.Registerer) *ReadOnlyKeyValueStateWrapper {
	w := &ReadOnlyKeyValueStateWrapper{
		logger,
		nil,
		subprocesses.Subprocesses{},
		make(chan struct{}),
		sync.Once{},
		sync.Mutex{},
		nil,
		list.NewList[*keyValueStateReadTransaction](),
	}
	w.metrics = newReadOnlyKeyValueStateMetrics(
		registerer,
		logger,
		func() float64 { return w.openReadTransactionsSummary().oldestDuration.Seconds() },
	)
	w.subs.Go(w.complainAboutLongHeldReadTransactions)
	return w
}

type openReadTransactionsSummary struct {
	count int
	// oldestSeqNr and oldestDuration are zero if count is zero.
	oldestSeqNr    uint64
	oldestDuration time.Duration
}

func (w *ReadOnlyKeyValueStateWrapper) openReadTransactionsSummary() openReadTransactionsSummary {
	w.mu.Lock()
	defer w.mu.Unlock()

	summary := openReadTransactionsSummary{w.openReadTransactions.Len(), 0, 0}
	if oldest, ok := w.openReadTransactions.Front(); ok {
		summary.oldestSeqNr = oldest.seqNr
		summary.oldestDuration = time.Since(oldest.openedAt)
	}
	return summary
}

// complainAboutLongHeldReadTransactions polls rather than arming a timer per
// read transaction, both to keep the plugin's read path free of per-transaction
// overhead and so that a plugin holding one forever keeps being complained
// about, not just once.
func (w *ReadOnlyKeyValueStateWrapper) complainAboutLongHeldReadTransactions() {
	ticker := time.NewTicker(longHeldKeyValueStateReadTransactionCheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-w.chStop:
			return
		case <-ticker.C:
			summary := w.openReadTransactionsSummary()
			if summary.oldestDuration <= recommendedMaxKeyValueStateReadTransactionDuration {
				continue
			}
			w.logger.Warn("KeyValueStateReadTransaction opened by ReportingPlugin is being held open for longer than recommended", commontypes.LogFields{
				"oldestOpenReadTransactionSeqNr":    summary.oldestSeqNr,
				"oldestOpenReadTransactionDuration": summary.oldestDuration.String(),
				"recommendedMaxDuration":            recommendedMaxKeyValueStateReadTransactionDuration.String(),
				"currentlyOpenReadTransactions":     summary.count,
			})
		}
	}
}

func (w *ReadOnlyKeyValueStateWrapper) Enable(keyValueDatabase protocol.KeyValueDatabase) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.nilOrKeyValueDatabase = keyValueDatabase
}

// Disable stops handing out read transactions and discards those still open.
// It undoes Enable; it does not release the resources acquired by
// NewReadOnlyKeyValueStateWrapper, that is what Close is for.
//
// It is called after the reporting plugin has been closed and before the
// underlying key value database is closed. A well-behaved plugin has stopped
// its goroutines and discarded its transactions by the time Close returns, so
// anything still open here is a leak: we log it and discard it on the plugin's
// behalf, which releases the database's resources and makes the plugin's next
// Read fail with a clear error rather than reaching a closed database.
//
// Read holds the transaction's shared lock across the call into the database.
// Discard takes the exclusive lock, waiting for in-flight reads to finish before
// releasing the transaction. Subsequent reads fail with
// errKeyValueStateReadTransactionDiscarded. We do not wait for the plugin to
// discard the transaction itself.
func (w *ReadOnlyKeyValueStateWrapper) Disable() {
	w.mu.Lock()
	w.nilOrKeyValueDatabase = nil
	leakedReadTransactions := slices.Collect(w.openReadTransactions.All())
	w.mu.Unlock()

	if len(leakedReadTransactions) > 0 {
		oldest := leakedReadTransactions[0]
		w.logger.Error("ReportingPlugin has been closed but is still holding read transactions open. This is a bug in the ReportingPlugin.", commontypes.LogFields{
			"oldestOpenReadTransactionSeqNr":    oldest.seqNr,
			"oldestOpenReadTransactionDuration": time.Since(oldest.openedAt).String(),
			"currentlyOpenReadTransactions":     len(leakedReadTransactions),
		})
	}

	// Discard takes w.mu to unregister, so we must not be holding it here.
	for _, readTransaction := range leakedReadTransactions {
		readTransaction.Discard()
	}
}

// Close releases everything NewReadOnlyKeyValueStateWrapper acquired: it stops
// the goroutine that complains about long-held read transactions and
// unregisters the metrics. It Disables first, so that on its own it is a
// complete cleanup and can be deferred right after construction as a safety
// net, regardless of whether Enable was ever called. Closing more than once is
// a no-op.
//
// Close inherits Disable's synchronization when discarding transactions.
func (w *ReadOnlyKeyValueStateWrapper) Close() {
	// Idempotent, so this is a no-op if the caller has already Disabled, as the
	// managed oracle does on its happy path.
	w.Disable()
	w.stopOnce.Do(func() { close(w.chStop) })
	w.subs.Wait()
	w.metrics.Close()
}

func (w *ReadOnlyKeyValueStateWrapper) NewReadTransaction(ctx context.Context) (ocr3_1types.KeyValueStateReadTransaction, error) {
	if outsideKeyValueStateReadsDisallowed(ctx) {
		w.metrics.refusedInsidePluginCallTotal.Inc()
		return nil, errKeyValueStateReadInsidePluginCall
	}

	startedAt := time.Now()
	readTransaction, err := w.newReadTransaction()
	if err != nil {
		w.metrics.refusedUnavailableTotal.Inc()
		return nil, err
	}
	w.metrics.transactionOpeningDuration.Observe(time.Since(startedAt).Seconds())
	return readTransaction, nil
}

func (w *ReadOnlyKeyValueStateWrapper) newReadTransaction() (*keyValueStateReadTransaction, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.nilOrKeyValueDatabase == nil {
		return nil, errKeyValueStateUnavailable
	}

	rawReadTransaction, err := w.nilOrKeyValueDatabase.NewReadTransactionUnchecked()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create read transaction: %w", errKeyValueStateUnavailable, err)
	}

	// Both of the following read from the very transaction we are about to hand
	// to the plugin, so the state it sees and the sequence number we report for
	// it are guaranteed to agree.
	if err := checkNotClobbered(rawReadTransaction); err != nil {
		rawReadTransaction.Discard()
		return nil, fmt.Errorf("%w: %w", errKeyValueStateUnavailable, err)
	}
	seqNr, err := rawReadTransaction.ReadHighestCommittedSeqNr()
	if err != nil {
		rawReadTransaction.Discard()
		return nil, fmt.Errorf("%w: failed to read highest committed seq nr: %w", errKeyValueStateUnavailable, err)
	}

	readTransaction := &keyValueStateReadTransaction{
		seqNr,
		time.Now(),
		list.Element[*keyValueStateReadTransaction]{},
		nil,
		sync.RWMutex{},
		rawReadTransaction,
		sync.Mutex{},
		0,
		0,
	}
	readTransaction.element = w.openReadTransactions.PushBack(readTransaction)
	readTransaction.onDiscard = func() {
		w.unregisterReadTransaction(readTransaction)
	}
	w.metrics.openedTransactionsTotal.Inc()
	return readTransaction, nil
}

func (w *ReadOnlyKeyValueStateWrapper) unregisterReadTransaction(readTransaction *keyValueStateReadTransaction) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.openReadTransactions.Remove(readTransaction.element)
	w.metrics.discardedTransactionsTotal.Inc()
	w.metrics.transactionAgeAtDiscard.Observe(time.Since(readTransaction.openedAt).Seconds())
	w.metrics.transactionReads.Observe(float64(readTransaction.reads))
	w.metrics.transactionReadBytes.Observe(float64(readTransaction.readBytes))
}

type keyValueStateReadTransaction struct {
	seqNr uint64
	// openedAt and element are only accessed by the wrapper, under its mu.
	openedAt time.Time
	element  list.Element[*keyValueStateReadTransaction]
	// onDiscard is invoked exactly once, by Discard.
	onDiscard func()

	// mu makes the transaction safe for concurrent use, as promised by
	// [ocr3_1types.KeyValueStateReadTransaction]. Reads take it for reading and
	// run in parallel; the underlying database read is itself safe for that.
	// Discard takes it for writing so that it cannot close the transaction
	// underneath an in-flight read. Disable calls Discard from another
	// goroutine, so this is not merely defensive.
	mu sync.RWMutex
	// nilOrRawReadTransaction is nil once discarded.
	nilOrRawReadTransaction protocol.KeyValueDatabaseReadTransaction

	// readStatsMu serializes counter updates without serializing database
	// reads. Read also holds mu for reading, so Discard's exclusive mu lock
	// ensures all updates finish before onDiscard observes the totals.
	readStatsMu sync.Mutex
	reads       uint64
	readBytes   uint64
}

var _ ocr3_1types.KeyValueStateReadTransaction = &keyValueStateReadTransaction{}

func (r *keyValueStateReadTransaction) Read(key []byte) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.nilOrRawReadTransaction == nil {
		return nil, errKeyValueStateReadTransactionDiscarded
	}
	value, err := r.nilOrRawReadTransaction.Read(key)
	r.readStatsMu.Lock()
	defer r.readStatsMu.Unlock()
	r.reads++
	if err != nil {
		return value, err
	}
	r.readBytes += uint64(len(value))
	return value, nil
}

func (r *keyValueStateReadTransaction) SeqNr() uint64 {
	return r.seqNr
}

func (r *keyValueStateReadTransaction) Discard() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.nilOrRawReadTransaction == nil {
		return
	}
	r.nilOrRawReadTransaction.Discard()
	r.nilOrRawReadTransaction = nil
	r.onDiscard()
}
