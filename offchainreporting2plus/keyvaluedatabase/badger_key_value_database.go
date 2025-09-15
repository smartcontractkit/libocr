package keyvaluedatabase

import (
	"bytes"
	"errors"
	"path"

	badger "github.com/dgraph-io/badger/v4"

	"github.com/smartcontractkit/libocr/internal/util"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// NewBadgerKeyValueDatabaseFactory produces a
// [ocr3_1types.KeyValueDatabaseFactory] that creates [Badger] databases under
// the directory indicated by baseDir. The directory must exist and be
// writeable. NewKeyValueDatabase may fail if not. The factory requires
// exclusive control of the directory: external changes are forbidden.
//
// [Badger]: https://pkg.go.dev/github.com/dgraph-io/badger/v4
func NewBadgerKeyValueDatabaseFactory(baseDir string) ocr3_1types.KeyValueDatabaseFactory {
	return &badgerKeyValueDatabaseFactory{baseDir}
}

type badgerKeyValueDatabaseFactory struct{ baseDir string }

var _ ocr3_1types.KeyValueDatabaseFactory = &badgerKeyValueDatabaseFactory{}

func (b *badgerKeyValueDatabaseFactory) NewKeyValueDatabase(configDigest types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	path := path.Join(b.baseDir, configDigest.String())
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return &badgerKeyValueDatabase{db}, nil
}

type badgerKeyValueDatabase struct{ raw *badger.DB }

var _ ocr3_1types.KeyValueDatabase = &badgerKeyValueDatabase{}

func (b *badgerKeyValueDatabase) Close() error {
	return b.raw.Close()
}

func (b *badgerKeyValueDatabase) NewReadTransaction() (ocr3_1types.KeyValueReadTransaction, error) {
	txn := b.raw.NewTransaction(false)
	return &badgerReadWriteTransaction{txn}, nil // write funcs are type erased
}

func (b *badgerKeyValueDatabase) NewReadWriteTransaction() (ocr3_1types.KeyValueReadWriteTransaction, error) {
	txn := b.raw.NewTransaction(true)
	return &badgerReadWriteTransaction{txn}, nil
}

type badgerReadWriteTransaction struct{ view *badger.Txn }

var _ ocr3_1types.KeyValueReadWriteTransaction = &badgerReadWriteTransaction{}
var _ ocr3_1types.KeyValueReadTransaction = &badgerReadWriteTransaction{}

func (b *badgerReadWriteTransaction) Discard() {
	b.view.Discard()
}

func (b *badgerReadWriteTransaction) Read(key []byte) ([]byte, error) {
	item, err := b.view.Get(key)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}
	val, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}
	return util.NilCoalesceSlice(val), nil
}

type badgerRangeIterator struct {
	hiKeyExcl []byte

	it        *badger.Iterator
	firstTime bool

	currentItem    *badger.Item
	currentItemKey []byte
}

var _ ocr3_1types.KeyValueIterator = &badgerRangeIterator{}

func (b *badgerRangeIterator) Close() error {
	b.it.Close()
	return nil
}

func (b *badgerRangeIterator) Next() bool {
	if b.firstTime {
		b.firstTime = false
	} else {
		b.it.Next()
	}

	if !b.it.Valid() {
		return false
	}

	item := b.it.Item()
	key := item.KeyCopy(nil)

	if len(b.hiKeyExcl) > 0 && bytes.Compare(key, b.hiKeyExcl) >= 0 {
		return false
	}

	b.currentItem = item
	b.currentItemKey = key

	return true
}

func (b *badgerRangeIterator) Key() []byte {
	return bytes.Clone(b.currentItemKey)
}

func (b *badgerRangeIterator) Value() ([]byte, error) {
	val, err := b.currentItem.ValueCopy(nil)
	if err != nil {
		return nil, err
	}
	return util.NilCoalesceSlice(val), nil
}

func (b *badgerRangeIterator) Err() error {
	return nil
}

func (b *badgerReadWriteTransaction) Range(loKey []byte, hiKeyExcl []byte) ocr3_1types.KeyValueIterator {
	loKey = bytes.Clone(loKey)
	hiKeyExcl = bytes.Clone(hiKeyExcl)

	loKey = util.NilCoalesceSlice(loKey)

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false // iterator offers values only on demand
	opts.AllVersions = false    // so that we don't have to check [badger.Item.IsDeletedOrExpired]
	it := b.view.NewIterator(opts)

	it.Seek(loKey)

	return &badgerRangeIterator{
		hiKeyExcl,
		it,
		true,
		nil,
		nil,
	}
}

func (b *badgerReadWriteTransaction) Write(key, value []byte) error {
	// Badger: The current transaction keeps a reference to the key and val byte
	// slice arguments. Users must not modify key and val until the end of the
	// transaction.
	key = bytes.Clone(key)
	value = bytes.Clone(value)

	return b.view.Set(key, util.NilCoalesceSlice(value))
}

func (b *badgerReadWriteTransaction) Delete(key []byte) error {
	// Badger: The current transaction keeps a reference to the key byte slice
	// argument. Users must not modify the key until the end of the transaction.
	key = bytes.Clone(key)

	return b.view.Delete(key)
}

func (b *badgerReadWriteTransaction) Commit() error {
	return b.view.Commit()
}
