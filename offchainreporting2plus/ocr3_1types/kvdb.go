package ocr3_1types

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Deprecated: Use KeyValueDatabaseReadWriteTransaction instead.
type KeyValueReadWriteTransaction = KeyValueDatabaseReadWriteTransaction

// A KeyValueDatabaseReadWriteTransaction is a read-write transaction over the
// database. Writes and deletes are visible to subsequent reads within the
// same transaction. See [KeyValueDatabaseReadTransaction] for the lifetime
// rules.
//
// # Concurrency
//
// Unlike a read transaction, a read-write transaction is not safe for
// concurrent use.
type KeyValueDatabaseReadWriteTransaction interface {
	KeyValueDatabaseReadTransaction
	// A value of nil is interpreted as an empty slice, and does *not* delete
	// the key. For deletions you must use the Delete method.
	Write(key []byte, value []byte) error
	Delete(key []byte) error

	// Commit commits the transaction. Whether it succeeds or fails, no
	// further calls other than Discard may be made on the transaction once
	// Commit returns.
	Commit() error
}

// Deprecated: Use KeyValueDatabaseReadTransaction instead.
type KeyValueReadTransaction = KeyValueDatabaseReadTransaction

// A KeyValueDatabaseReadTransaction is a read transaction over the database.
//
// # Lifetime
//
// Every transaction must be discarded — or, for a
// [KeyValueDatabaseReadWriteTransaction], committed (or discarded to abort
// it). All iterators obtained from Range must be closed before the
// transaction is discarded or committed. After Discard or Commit return, no
// further calls may be made on the transaction, except that Discard may be
// called multiple times.
//
// # Concurrency
//
// Read and Range may be called concurrently with each other and themselves.
// Iterators from different Range calls may be used concurrently with each
// other and with Read and Range calls, but an individual iterator is not
// safe for concurrent use. Discard must not be called concurrently with any
// other call on the transaction.
type KeyValueDatabaseReadTransaction interface {
	// If the key exists, the returned value must not be nil!
	Read(key []byte) ([]byte, error)
	// Range iterates over the key-value pairs with keys in the range [loKey,
	// hiKeyExcl), in ascending order of key. Key-value stores typically store
	// keys in a sorted order, making this a fast operation. loKey can be set to
	// 0 length or nil for iteration without a lower bound. hiKeyExcl can be set
	// to 0 length or nil for iteration without an upper bound.
	//
	// WARNING: DO NOT perform any writes/deletes to the key-value store while
	// the iterator is opened.
	Range(loKey []byte, hiKeyExcl []byte) KeyValueDatabaseIterator
	// Discard discards the transaction and releases its resources. See
	// [KeyValueDatabaseReadTransaction] for the lifetime rules.
	Discard()
}

// Deprecated: Use KeyValueDatabaseIterator instead.
type KeyValueIterator = KeyValueDatabaseIterator

// KeyValueDatabaseIterator is an iterator over key-value pairs, in ascending
// order of keys. See [KeyValueDatabaseReadTransaction] for the lifetime and
// concurrency rules that apply to iterators.
//

// Example usage:
//
//	it := kvReader.Range(loKey, hiKeyExcl)
//	defer it.Close()
//	for it.Next() {
//	    key := it.Key()
//	    value, err := it.Value()
//	    if err != nil {
//	        // handle error
//	    }
//	    // process key and value
//	}
//	if err := it.Err(); err != nil {
//	    // handle error
//	}
type KeyValueDatabaseIterator interface {
	// Next prepares the next key-value pair for reading. It returns true on
	// success, or false if there is no next key-value pair or an error occurred
	// while preparing it.
	Next() bool
	// Key returns the key of the current key-value pair.
	Key() []byte
	// Value returns the value of the current key-value pair. An error value
	// indicates a failure to retrieve the value, and the caller is responsible
	// for handling it. Even if all errors are nil, [KeyValueDatabaseIterator.Err] must
	// be checked after iteration is completed.
	Value() ([]byte, error)
	// Err returns any error encountered during iteration. Must be checked after
	// the end of the iteration, to ensure that no key-value pairs were missed
	// due to iteration errors. Errors in [KeyValueDatabaseIterator.Value] are distinct
	// and will not cause a non-nil error.
	Err() error
	// Close closes the iterator and releases any resources associated with it.
	// Must be called in any case, even if the iteration encountered any
	// error through [KeyValueDatabaseIterator.Value] or [KeyValueDatabaseIterator.Err].
	// After Close returns, no further calls may be made on the iterator,
	// including another Close.
	Close() error
}

// KeyValueDatabase is a key-value database.
//
// NewReadTransaction and NewReadWriteTransaction may be called concurrently
// with each other and themselves. At most one read-write transaction is
// active at a time: NewReadWriteTransaction blocks until the previous
// read-write transaction has been committed or discarded.
//
// Every transaction must be committed or discarded (see
// [KeyValueDatabaseReadTransaction]) before Close is called. Close must not
// be called concurrently with any other method, and no method may be called
// after Close.
type KeyValueDatabase interface {
	NewReadWriteTransaction() (KeyValueDatabaseReadWriteTransaction, error)
	NewReadTransaction() (KeyValueDatabaseReadTransaction, error)

	Close() error
}

var ErrKeyValueDatabaseDoesNotExist = fmt.Errorf("key value database does not exist")

type KeyValueDatabaseFactory interface {
	NewKeyValueDatabase(configDigest types.ConfigDigest) (KeyValueDatabase, error)
	NewKeyValueDatabaseIfExists(configDigest types.ConfigDigest) (KeyValueDatabase, error)
}
