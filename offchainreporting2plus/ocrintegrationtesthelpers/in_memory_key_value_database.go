package ocrintegrationtesthelpers

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/google/btree"
	"github.com/smartcontractkit/libocr/internal/util"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

var (
	errKeyValueDatabaseClosed          = fmt.Errorf("key value database closed")
	errTransactionCommittedOrDiscarded = fmt.Errorf("transaction committed or discarded")
	errKeyValueDatabaseAlreadyOpen     = fmt.Errorf("key value database already open")
)

type item struct {
	key   []byte
	value []byte
}

func (i item) Less(than item) bool {
	return bytes.Compare(i.key, than.key) < 0
}

// StatelessInMemoryKeyValueDatabaseFactory is a stateless factory for testing.
// Each call to NewKeyValueDatabase returns a fresh, empty database.
// NewKeyValueDatabaseIfExists always returns ErrKeyValueDatabaseDoesNotExist.
// No exclusive access is enforced: multiple concurrent calls to
// NewKeyValueDatabase with the same configDigest will each succeed and return
// independent databases. Use [StatefulInMemoryKeyValueDatabaseFactory] if you
// need a NewKeyValueDatabase* call to return a database with contents written
// before a previous Close call.
type StatelessInMemoryKeyValueDatabaseFactory struct{}

var _ ocr3_1types.KeyValueDatabaseFactory = StatelessInMemoryKeyValueDatabaseFactory{}

func NewStatelessInMemoryKeyValueDatabaseFactory() StatelessInMemoryKeyValueDatabaseFactory {
	return StatelessInMemoryKeyValueDatabaseFactory{}
}

func (StatelessInMemoryKeyValueDatabaseFactory) NewKeyValueDatabase(types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	return NewInMemoryKeyValueDatabase(), nil
}

func (StatelessInMemoryKeyValueDatabaseFactory) NewKeyValueDatabaseIfExists(types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	return nil, ocr3_1types.ErrKeyValueDatabaseDoesNotExist
}

// StatefulInMemoryKeyValueDatabaseFactory is a factory for testing that retains
// database contents in memory (not on disk) for the lifetime of the factory.
// After calling db.Close(), a subsequent NewKeyValueDatabase or
// NewKeyValueDatabaseIfExists call for the same configDigest returns a database
// with the contents written before Close. Enforces exclusive access:
// NewKeyValueDatabase* returns an error if a database for that configDigest has
// been opened but not yet closed. Use [StatefulInMemoryKeyValueDatabaseFactory.ForgetKeyValueDatabaseForTests] to make
// the factory forget about a configDigest (akin to deleting the database from
// the filesystem for a disk-based implementation).
type StatefulInMemoryKeyValueDatabaseFactory struct {
	mu  sync.Mutex
	dbs map[types.ConfigDigest]*InMemoryKeyValueDatabase
}

var _ ocr3_1types.KeyValueDatabaseFactory = &StatefulInMemoryKeyValueDatabaseFactory{}

func NewStatefulInMemoryKeyValueDatabaseFactory() *StatefulInMemoryKeyValueDatabaseFactory {
	return &StatefulInMemoryKeyValueDatabaseFactory{
		sync.Mutex{},
		make(map[types.ConfigDigest]*InMemoryKeyValueDatabase),
	}
}

func (f *StatefulInMemoryKeyValueDatabaseFactory) NewKeyValueDatabase(configDigest types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	return f.newKeyValueDatabase(configDigest, true)
}

func (f *StatefulInMemoryKeyValueDatabaseFactory) NewKeyValueDatabaseIfExists(configDigest types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	return f.newKeyValueDatabase(configDigest, false)
}

// ForgetKeyValueDatabaseForTests makes the factory forget about a configDigest.
// The next NewKeyValueDatabase call for this configDigest will create a fresh
// database; NewKeyValueDatabaseIfExists will return ErrKeyValueDatabaseDoesNotExist.
// Any currently open database for this configDigest is unaffected and continues
// to work, but the exclusivity guarantee is lost: the database opened before
// [StatefulInMemoryKeyValueDatabaseFactory.ForgetKeyValueDatabaseForTests] and the database opened after can be open
// simultaneously.
func (f *StatefulInMemoryKeyValueDatabaseFactory) ForgetKeyValueDatabaseForTests(configDigest types.ConfigDigest) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.dbs, configDigest)
}

func (f *StatefulInMemoryKeyValueDatabaseFactory) newKeyValueDatabase(configDigest types.ConfigDigest, createIfNotExists bool) (ocr3_1types.KeyValueDatabase, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	existingDb, exists := f.dbs[configDigest]
	if exists {
		existingDb.mu.Lock()
		if !existingDb.closed {
			existingDb.mu.Unlock()
			return nil, errKeyValueDatabaseAlreadyOpen
		}
		// Closed - transplant tree to new database
		tree := existingDb.tree
		existingDb.mu.Unlock()

		newDb := newInMemoryKeyValueDatabaseWithTree(tree)
		f.dbs[configDigest] = newDb
		return newDb, nil
	}

	if !createIfNotExists {
		return nil, ocr3_1types.ErrKeyValueDatabaseDoesNotExist
	}

	newDb := NewInMemoryKeyValueDatabase()
	f.dbs[configDigest] = newDb
	return newDb, nil
}

type InMemoryKeyValueDatabase struct {
	ctx    context.Context
	cancel context.CancelFunc
	subs   subprocesses.Subprocesses

	rwSerializationLock sync.Mutex

	mu     sync.Mutex
	tree   *btree.BTreeG[item]
	closed bool
}

var _ ocr3_1types.KeyValueDatabase = &InMemoryKeyValueDatabase{}

// NewInMemoryKeyValueDatabase creates a standalone in-memory database without
// factory tracking. Use [NewStatefulInMemoryKeyValueDatabaseFactory] for
// exclusive access and persistence semantics.
func NewInMemoryKeyValueDatabase() *InMemoryKeyValueDatabase {
	return newInMemoryKeyValueDatabaseWithTree(btree.NewG(32, func(a, b item) bool {
		return a.Less(b)
	}))
}

func newInMemoryKeyValueDatabaseWithTree(tree *btree.BTreeG[item]) *InMemoryKeyValueDatabase {
	ctx, cancel := context.WithCancel(context.Background())
	return &InMemoryKeyValueDatabase{
		ctx,
		cancel,
		subprocesses.Subprocesses{},
		sync.Mutex{},
		sync.Mutex{},
		tree,
		false,
	}
}

func (db *InMemoryKeyValueDatabase) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		return errKeyValueDatabaseClosed
	}

	db.closed = true
	db.cancel()
	db.subs.Wait()
	return nil
}

func (db *InMemoryKeyValueDatabase) NewReadTransaction() (ocr3_1types.KeyValueDatabaseReadTransaction, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		return nil, errKeyValueDatabaseClosed
	}

	ctx, cancel := context.WithCancel(db.ctx)
	return &inMemoryReadWriteTransaction{
		ctx,
		cancel,
		db,
		sync.Mutex{},
		db.tree.Clone(),
		true,
	}, nil
}

func (db *InMemoryKeyValueDatabase) NewReadWriteTransaction() (ocr3_1types.KeyValueDatabaseReadWriteTransaction, error) {
	db.rwSerializationLock.Lock()

	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		db.rwSerializationLock.Unlock()
		return nil, errKeyValueDatabaseClosed
	}

	ctx, cancel := context.WithCancel(db.ctx)
	updatedTree := db.tree.Clone()
	return &inMemoryReadWriteTransaction{
		ctx,
		cancel,
		db,
		sync.Mutex{},
		updatedTree,
		false,
	}, nil
}

func (db *InMemoryKeyValueDatabase) commit(updatedTree *btree.BTreeG[item]) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		return errKeyValueDatabaseClosed
	}
	db.tree = updatedTree
	return nil
}

type inMemoryReadWriteTransaction struct {
	ctx    context.Context
	cancel context.CancelFunc

	db *InMemoryKeyValueDatabase

	mu          sync.Mutex
	updatedTree *btree.BTreeG[item]
	readOnly    bool
}

func (tx *inMemoryReadWriteTransaction) Range(loKey []byte, hiKeyExcl []byte) ocr3_1types.KeyValueDatabaseIterator {
	loKey = bytes.Clone(loKey)
	hiKeyExcl = bytes.Clone(hiKeyExcl)

	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.updatedTree == nil {
		return &inMemoryRangeIterator{
			errTransactionCommittedOrDiscarded,
			nil,
			nil,
			nil,
			nil,
		}
	}

	ctx, cancel := context.WithCancel(tx.ctx)
	ch := make(chan item)
	// Btree is not safe to share across goroutines with concurrent writes. We
	// intentionally do not clone because as per the
	// ocr3_1types.KeyValueDatabaseIterator interface, users must not write or
	// delete keys while an iterator is opened. We want this to cause a nil
	// dereference, panic, or other problematic behavior, so that we'll know we
	// misuse the interface in protocol code.
	tree := tx.updatedTree

	iter := &inMemoryRangeIterator{
		nil,
		cancel,
		ch,
		tx,
		nil,
	}

	tx.db.subs.Go(func() {
		defer close(ch)
		tree.AscendGreaterOrEqual(item{loKey, nil}, func(i item) bool {
			if len(hiKeyExcl) > 0 && bytes.Compare(i.key, hiKeyExcl) >= 0 {
				return false
			}
			select {
			case ch <- i:
				return true
			case <-ctx.Done():
				return false
			}
		})
	})

	return iter
}

type inMemoryRangeIterator struct {
	err    error
	cancel context.CancelFunc
	ch     chan item
	tx     *inMemoryReadWriteTransaction

	currentKey []byte
}

var _ ocr3_1types.KeyValueDatabaseIterator = &inMemoryRangeIterator{}

func (i *inMemoryRangeIterator) Close() error {
	if i.cancel != nil {
		i.cancel()
	}
	return nil
}

func (i *inMemoryRangeIterator) Next() bool {
	if i.err != nil {
		return false
	}

	nextItem, ok := <-i.ch
	if !ok {
		return false
	}

	i.currentKey = nextItem.key
	return true
}

func (i *inMemoryRangeIterator) Key() []byte {
	return bytes.Clone(i.currentKey)
}

func (i *inMemoryRangeIterator) Value() ([]byte, error) {
	return i.tx.Read(i.currentKey)
}

func (i *inMemoryRangeIterator) Err() error {
	return i.err
}

var _ ocr3_1types.KeyValueDatabaseReadWriteTransaction = &inMemoryReadWriteTransaction{}
var _ ocr3_1types.KeyValueDatabaseReadTransaction = &inMemoryReadWriteTransaction{}

func (tx *inMemoryReadWriteTransaction) Commit() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.readOnly {
		return fmt.Errorf("read only transaction cannot be committed")
	}
	if tx.updatedTree == nil {
		return errTransactionCommittedOrDiscarded
	}
	defer tx.db.rwSerializationLock.Unlock()

	if err := tx.db.commit(tx.updatedTree); err != nil {
		return err
	}

	tx.updatedTree = nil
	tx.cancel()
	return nil
}

func (tx *inMemoryReadWriteTransaction) Delete(key []byte) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.readOnly {
		return fmt.Errorf("read only transaction cannot delete")
	}
	if tx.updatedTree == nil {
		return errTransactionCommittedOrDiscarded
	}

	tx.updatedTree.Delete(item{key, nil})
	return nil
}

func (tx *inMemoryReadWriteTransaction) Write(key []byte, value []byte) error {
	// protect against later modification
	key = bytes.Clone(key)
	value = bytes.Clone(value)
	value = util.NilCoalesceSlice(value)

	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.readOnly {
		return fmt.Errorf("read only transaction cannot write")
	}
	if tx.updatedTree == nil {
		return errTransactionCommittedOrDiscarded
	}

	tx.updatedTree.ReplaceOrInsert(item{key, value})
	return nil
}

func (tx *inMemoryReadWriteTransaction) Discard() {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.updatedTree == nil {
		return
	}

	tx.updatedTree = nil
	tx.cancel()
	if !tx.readOnly {
		tx.db.rwSerializationLock.Unlock()
	}
}

func (tx *inMemoryReadWriteTransaction) Read(key []byte) ([]byte, error) {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.updatedTree == nil {
		return nil, errTransactionCommittedOrDiscarded
	}

	if v, ok := tx.updatedTree.Get(item{key, nil}); ok {
		// protect against later modification
		return bytes.Clone(v.value), nil
	}
	return nil, nil
}
