package shim

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

type keyValueDatabaseWithMetrics struct {
	ocr3_1types.KeyValueDatabase
	metrics *keyValueDatabaseMetrics
}

type keyValueDatabaseReadTransactionWithMetrics struct {
	ocr3_1types.KeyValueDatabaseReadTransaction
	metrics *keyValueDatabaseMetrics
	once    sync.Once
}

// Discard implements ocr3_1types.KeyValueDatabaseReadTransaction.
func (k *keyValueDatabaseReadTransactionWithMetrics) Discard() {
	k.once.Do(func() {
		k.metrics.discardedReadTransactionsTotal.Inc()
	})
	k.KeyValueDatabaseReadTransaction.Discard()
}

type keyValueDatabaseReadWriteTransactionWithMetrics struct {
	ocr3_1types.KeyValueDatabaseReadWriteTransaction
	metrics *keyValueDatabaseMetrics
	once    sync.Once
}

// Discard implements ocr3_1types.KeyValueDatabaseReadWriteTransaction.
func (k *keyValueDatabaseReadWriteTransactionWithMetrics) Discard() {
	k.once.Do(func() {
		k.metrics.discardedReadWriteTransactionsTotal.Inc()
	})
	k.KeyValueDatabaseReadWriteTransaction.Discard()
}

// Commit implements ocr3_1types.KeyValueDatabaseReadWriteTransaction.
func (k *keyValueDatabaseReadWriteTransactionWithMetrics) Commit() error {
	err := k.KeyValueDatabaseReadWriteTransaction.Commit()
	k.once.Do(func() {
		if err == nil {
			k.metrics.committedReadWriteTransactionsTotal.Inc()
		} else {

			k.metrics.discardedReadWriteTransactionsTotal.Inc()
		}
	})
	return err
}

// Close implements ocr3_1types.KeyValueDatabase. It does not manage the lifetime of the underlying database,
// and is expected to be closed first, before the underlying database is closed.
func (k *keyValueDatabaseWithMetrics) Close() error {
	k.metrics.Close()
	return nil
}

// NewReadTransaction implements ocr3_1types.KeyValueDatabase.
func (k *keyValueDatabaseWithMetrics) NewReadTransaction() (ocr3_1types.KeyValueDatabaseReadTransaction, error) {
	tx, err := k.KeyValueDatabase.NewReadTransaction()
	if err != nil {
		return nil, err
	}
	k.metrics.openedReadTransactionsTotal.Inc()
	return &keyValueDatabaseReadTransactionWithMetrics{
		tx,
		k.metrics,
		sync.Once{},
	}, nil
}

// NewReadWriteTransaction implements ocr3_1types.KeyValueDatabase.
func (k *keyValueDatabaseWithMetrics) NewReadWriteTransaction() (ocr3_1types.KeyValueDatabaseReadWriteTransaction, error) {
	tx, err := k.KeyValueDatabase.NewReadWriteTransaction()
	if err != nil {
		return nil, err
	}
	k.metrics.openedReadWriteTransactionsTotal.Inc()
	return &keyValueDatabaseReadWriteTransactionWithMetrics{
		tx,
		k.metrics,
		sync.Once{},
	}, nil
}

func NewKeyValueDatabaseWithMetrics(
	kvDb ocr3_1types.KeyValueDatabase,
	metricsRegisterer prometheus.Registerer,
	logger commontypes.Logger,
) ocr3_1types.KeyValueDatabase {
	return &keyValueDatabaseWithMetrics{
		kvDb,
		newKeyValueDatabaseMetrics(metricsRegisterer, logger),
	}
}

var _ ocr3_1types.KeyValueDatabase = &keyValueDatabaseWithMetrics{}
