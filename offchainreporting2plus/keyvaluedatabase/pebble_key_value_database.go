package keyvaluedatabase

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"sync"

	"github.com/cockroachdb/pebble"

	"github.com/RoSpaceDev/libocr/internal/util"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

var ErrPebbleReadWriteTransactionAlreadyOpen = errors.New("a read-write transaction is already open")

// NewPebbleKeyValueDatabaseFactory produces a
// [ocr3_1types.KeyValueDatabaseFactory] that creates Pebble databases under
// the directory indicated by baseDir. The directory must exist and be
// writeable. NewKeyValueDatabase may fail if not. The factory requires
// exclusive control of the directory: external changes are forbidden.
func NewPebbleKeyValueDatabaseFactory(baseDir string) ocr3_1types.KeyValueDatabaseFactory {
	return &pebbleKeyValueDatabaseFactory{baseDir}
}

type pebbleKeyValueDatabaseFactory struct{ baseDir string }

var _ ocr3_1types.KeyValueDatabaseFactory = &pebbleKeyValueDatabaseFactory{}

func (p *pebbleKeyValueDatabaseFactory) NewKeyValueDatabase(configDigest types.ConfigDigest) (ocr3_1types.KeyValueDatabase, error) {
	dbPath := path.Join(p.baseDir, fmt.Sprintf("%s.db", configDigest.String()))

	db, err := pebble.Open(dbPath, nil)
	if err != nil {
		return nil, err
	}

	return &pebbleKeyValueDatabase{
		db,
		sync.Mutex{},
		sync.Once{},
	}, nil
}

type pebbleKeyValueDatabase struct {
	db *pebble.DB

	// This lock enforces that we can have at most one active committable
	// read-write transaction open at any point in time.
	rwSerializationLock sync.Mutex
	closeOnce           sync.Once
}

var _ ocr3_1types.KeyValueDatabase = &pebbleKeyValueDatabase{}

func (p *pebbleKeyValueDatabase) Close() error {
	err := fmt.Errorf("database already closed")
	p.closeOnce.Do(func() {
		err = p.db.Close()
	})
	return err
}

// The resulting transaction is NOT thread-safe.

func (p *pebbleKeyValueDatabase) NewReadTransaction() (ocr3_1types.KeyValueDatabaseReadTransaction, error) {
	snapshot := p.db.NewSnapshot()
	return &pebbleReadTransaction{
		snapshot,
		false,
	}, nil
}

// The resulting transaction is NOT thread-safe.

func (p *pebbleKeyValueDatabase) NewReadWriteTransaction() (ocr3_1types.KeyValueDatabaseReadWriteTransaction, error) {
	p.rwSerializationLock.Lock()
	batch := p.db.NewIndexedBatch()
	return &pebbleReadWriteTransaction{
		batch,
		false,
		func() {
			p.rwSerializationLock.Unlock()
		},
	}, nil
}

type pebbleReadTransaction struct {
	snapshot  *pebble.Snapshot
	discarded bool
}

var _ ocr3_1types.KeyValueDatabaseReadTransaction = &pebbleReadTransaction{}

func (p *pebbleReadTransaction) Discard() {
	if p.discarded {
		return
	}
	p.discarded = true
	_ = p.snapshot.Close()
}

func (p *pebbleReadTransaction) Read(key []byte) ([]byte, error) {
	return readFromPebbleReader(p.snapshot, key)
}

func (p *pebbleReadTransaction) Range(loKey []byte, hiKeyExcl []byte) ocr3_1types.KeyValueDatabaseIterator {
	return newPebbleIterator(p.snapshot, loKey, hiKeyExcl)
}

type pebbleReadWriteTransaction struct {
	batch *pebble.Batch

	committedOrDiscarded     bool
	afterCommitOrDiscardFunc func()
}

var _ ocr3_1types.KeyValueDatabaseReadWriteTransaction = &pebbleReadWriteTransaction{}
var _ ocr3_1types.KeyValueDatabaseReadTransaction = &pebbleReadWriteTransaction{}

func (p *pebbleReadWriteTransaction) Discard() {
	if p.committedOrDiscarded {
		return
	}
	defer p.afterCommitOrDiscardFunc()
	p.committedOrDiscarded = true
	_ = p.batch.Close()
}

func (p *pebbleReadWriteTransaction) Read(key []byte) ([]byte, error) {
	return readFromPebbleReader(p.batch, key)
}

func (p *pebbleReadWriteTransaction) Range(loKey []byte, hiKeyExcl []byte) ocr3_1types.KeyValueDatabaseIterator {
	return newPebbleIterator(p.batch, loKey, hiKeyExcl)
}

func (p *pebbleReadWriteTransaction) Write(key, value []byte) error {
	return p.batch.Set(key, util.NilCoalesceSlice(value), nil)
}

func (p *pebbleReadWriteTransaction) Delete(key []byte) error {
	return p.batch.Delete(key, nil)
}

func (p *pebbleReadWriteTransaction) Commit() error {
	if p.committedOrDiscarded {
		return fmt.Errorf("transaction has been committed or discarded")
	}
	defer p.afterCommitOrDiscardFunc()
	p.committedOrDiscarded = true
	return p.batch.Commit(nil)
}

type pebbleIterator struct {
	iter      *pebble.Iterator
	loKey     []byte
	hiKeyExcl []byte
	err       error
	firstCall bool
}

var _ ocr3_1types.KeyValueDatabaseIterator = &pebbleIterator{}

func newPebbleIterator(reader pebble.Reader, loKey []byte, hiKeyExcl []byte) *pebbleIterator {
	loKey = util.NilCoalesceSlice(bytes.Clone(loKey))
	hiKeyExcl = bytes.Clone(hiKeyExcl)

	opts := &pebble.IterOptions{
		LowerBound: loKey,
		UpperBound: hiKeyExcl,
	}

	errorneousPebbleIterator := func(err error) *pebbleIterator {
		return &pebbleIterator{
			nil,
			nil,
			nil,
			err,
			true,
		}
	}

	iter, err := reader.NewIter(opts)
	if err != nil {
		return errorneousPebbleIterator(err)
	}

	return &pebbleIterator{
		iter,
		loKey,
		hiKeyExcl,
		nil,
		true,
	}
}

func (p *pebbleIterator) Close() error {
	if p.iter == nil {
		return p.err
	}
	return p.iter.Close()
}

func (p *pebbleIterator) Next() bool {
	if p.iter == nil || p.err != nil {
		return false
	}

	if p.firstCall {
		p.firstCall = false
		if !p.iter.First() {
			p.err = p.iter.Error()
			return false
		}
	} else {
		if !p.iter.Next() {
			p.err = p.iter.Error()
			return false
		}
	}

	return p.iter.Valid()
}

func (p *pebbleIterator) Key() []byte {
	if p.iter == nil {
		return nil
	}
	return bytes.Clone(p.iter.Key())
}

func (p *pebbleIterator) Value() ([]byte, error) {
	if p.iter == nil {
		return nil, nil
	}
	return util.NilCoalesceSlice(bytes.Clone(p.iter.Value())), nil
}

func (p *pebbleIterator) Err() error {
	if p.iter == nil || p.err != nil {
		return p.err
	}
	return p.iter.Error()
}

func readFromPebbleReader(reader pebble.Reader, key []byte) ([]byte, error) {
	value, closer, err := reader.Get(key)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	defer closer.Close()
	return util.NilCoalesceSlice(bytes.Clone(value)), nil
}
