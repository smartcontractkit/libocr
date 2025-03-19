package ocr3_1types

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type KeyNotFoundError string

func (e KeyNotFoundError) Error() string {
	return string(e)
}

const ErrKeyNotFound KeyNotFoundError = "key not found"

type DuplicateTransactionError string

func (e DuplicateTransactionError) Error() string {
	return string(e)
}

const ErrDuplicateTransaction DuplicateTransactionError = "a kv stored transaction is already open"

type KeyValueReader interface {
	// Read returns the value of the given key. If the key is not found, ErrKeyNotFound should be returned.
	Read(key []byte) ([]byte, error)

	// The addition of the following methods is still being evaluated

	// ReadRange returns values for keys in the range [key1, key2]
	//ReadRange(key1 []byte, key2 []byte) ([][]byte, error)
	// KeyRange returns the keys in the range [key1, key2]
	//KeyRange(key1 []byte, key2 []byte) ([][]byte, error)
	// ReadBatch returns the values for the specified keys. If some key is not found ErrKeyNotFound should be returned.
	//ReadBatch(keys [][]byte) ([][]byte, error)

	// Returns the maximum allowed key size
	//MaxKeySize() int
	// Returns the maximum allowed value size
	//MaxValueSize() int
}

type KeyValueReaderDiscardable interface {
	KeyValueReader
	Discard()
}

type KeyValueWriter interface {
	// Write sets a value to the given key. Creates a key if the key does not already exist.
	Write(key []byte, value []byte) error
	// Delete deletes the key.
	Delete(key []byte) error
}

type KeyValueReadWriter interface {
	KeyValueReader
	KeyValueWriter
}

type KeyValuePair struct {
	Key   []byte
	Value []byte
}

type KeyValueStoreTransaction interface {
	// SeqNr returns the sequence number of the OCR3 round for which this transaction was created
	SeqNr() uint64
	// GetReadWriter returns a KeyValueReadWriter implementation for this transaction
	GetReadWriter() (KeyValueReadWriter, error)
	// GetWriteSet returns a map from keys in string encoding to values that have been written in
	// this transaction. If the value of a key has been deleted, it is mapped to nil.
	GetWriteSet() []KeyValuePair
	// Commit commits the transaction to the key value store and then discards the transaction.
	Commit() error
	// Discard discards the transaction.
	Discard()
}

// We assume that all reads and writes in a transaction commit atomically
type KeyValueStore interface {
	NewTransaction(seqNr uint64) (KeyValueStoreTransaction, error)
	GetReader() KeyValueReaderDiscardable
	HighestCommittedSeqNr() (uint64, error)
	Close() error
}

type KeyValueStoreFactory interface {
	NewKeyValueStore(configDigest types.ConfigDigest) (KeyValueStore, error)
}
