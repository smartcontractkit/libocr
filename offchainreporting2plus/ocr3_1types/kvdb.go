package ocr3_1types

import (
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

// Deprecated: Use KeyValueDatabaseReadWriteTransaction instead.
type KeyValueReadWriteTransaction = KeyValueDatabaseReadWriteTransaction

type KeyValueDatabaseReadWriteTransaction interface {
	KeyValueDatabaseReadTransaction
	// A value of nil is interpreted as an empty slice, and does *not* delete
	// the key. For deletions you must use the Delete method.
	Write(key []byte, value []byte) error
	Delete(key []byte) error

	Commit() error
}

// Deprecated: Use KeyValueDatabaseReadTransaction instead.
type KeyValueReadTransaction = KeyValueDatabaseReadTransaction

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
	Discard()
}

// Deprecated: Use KeyValueDatabaseIterator instead.
type KeyValueIterator = KeyValueDatabaseIterator

// KeyValueDatabaseIterator is a iterator over key-value pairs, in ascending order of
// keys.
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
	// Further iteration is prevented, i.e., [KeyValueDatabaseIterator.Next] will return
	// false. Must be called in any case, even if the iteration encountered any
	// error through [KeyValueDatabaseIterator.Value] or [KeyValueDatabaseIterator.Err].
	Close() error
}

type KeyValueDatabase interface {
	NewReadWriteTransaction() (KeyValueDatabaseReadWriteTransaction, error)
	NewReadTransaction() (KeyValueDatabaseReadTransaction, error)

	Close() error
}

type KeyValueDatabaseFactory interface {
	NewKeyValueDatabase(configDigest types.ConfigDigest) (KeyValueDatabase, error)
}
