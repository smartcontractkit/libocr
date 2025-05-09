package protocol

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type KeyValueStoreReadTransaction interface {
	// The only read part of the interface that the plugin might see. The rest
	// of the methods might only be called by protocol code.
	ocr3_1types.KeyValueReader
	KeyValueStoreSemanticRead
	Discard()
}

type KeyValueStoreSemanticRead interface {
	// Returns the sequence number of which the state the transaction
	// represents. Really read from the database here, no cached values allowed.
	ReadHighestCommittedSeqNr() (uint64, error)
}

type KeyValueStoreReadWriteTransaction interface {
	KeyValueStoreReadTransaction
	// The only write part of the interface that the plugin might see. The rest
	// of the methods might only be called by protocol code.
	ocr3_1types.KeyValueReadWriter
	KeyValueStoreSemanticWrite
	// Commit writes the new highest committed sequence number to the magic key
	// (if the transaction is _not_ unchecked) and commits the transaction to
	// the key value store, then discards the transaction.
	Commit() error
}

type KeyValueStoreSemanticWrite interface {
	// GetWriteSet returns a map from keys in string encoding to values that
	// have been written in this transaction. If the value of a key has been
	// deleted, it is mapped to nil.

	GetWriteSet() ([]KeyValuePair, error)

	// WriteHighestCommittedSeqNr writes the given sequence number to the magic
	// key. It is called before Commit on checked transactions.
	WriteHighestCommittedSeqNr(seqNr uint64) error
}

type KeyValuePair struct {
	Key   []byte
	Value []byte
}

type KeyValueStore interface {
	// Must error if the key value store is not ready to apply state transition
	// for the given sequence number. Must update the highest committed sequence
	// number magic key upon commit. Convenience method for synchronization
	// between outcome generation & state persistence.
	NewReadWriteTransaction(postSeqNr uint64) (KeyValueStoreReadWriteTransaction, error)
	// Must error if the key value store is not ready to apply state transition
	// for the given sequence number. Convenience method for synchronization
	// between outcome generation & state persistence.
	NewReadTransaction(postSeqNr uint64) (KeyValueStoreReadTransaction, error)

	// Unchecked transactions are useful when you don't care that the
	// transaction state represents the kv state as of some particular sequence
	// number, mostly when writing auxiliary data to the kv store. Unchecked
	// transactions do not update the highest committed sequence number magic
	// key upon commit, as would checked transactions.
	NewReadWriteTransactionUnchecked() (KeyValueStoreReadWriteTransaction, error)
	// Unchecked transactions are useful when you don't care that the
	// transaction state represents the kv state as of some particular sequence
	// number, mostly when reading auxiliary data from the kv store.
	NewReadTransactionUnchecked() (KeyValueStoreReadTransaction, error)

	// Deprecated: Kept for convenience/small diff, consider using
	// [KeyValueStoreSemanticRead.ReadHighestCommittedSeqNr] instead.
	HighestCommittedSeqNr() (uint64, error)
	Close() error
}

type KeyValueStoreFactory interface {
	NewKeyValueStore(configDigest types.ConfigDigest) (KeyValueStore, error)
}
