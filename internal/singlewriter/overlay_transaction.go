package singlewriter

import (
	"bytes"
	"fmt"
	"sort"
	"sync"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
)

type opKind int

const (
	_ opKind = iota
	opWrite
	opDelete
)

type operation struct {
	kind  opKind
	value []byte // only for writes
}

var ErrClosed = fmt.Errorf("transaction already committed/discarded")

type txStatus int

const (
	_ txStatus = iota
	txStatusOpen
	txStatusCommitted
	txStatusDiscarded
)

type overlayTransaction struct {
	overlay map[string]operation

	keyValueDatabase ocr3_1types.KeyValueDatabase
	rawReaderTx      ocr3_1types.KeyValueDatabaseReadTransaction

	// but using the mu mutex

	mu sync.Mutex

	iterators sync.WaitGroup
	status    txStatus
}

func (ot *overlayTransaction) Read(key []byte) ([]byte, error) {
	k := string(key)
	ot.mu.Lock()
	defer ot.mu.Unlock()
	if ot.status != txStatusOpen {
		return nil, ErrClosed
	}
	if op, ok := ot.overlay[k]; ok {
		switch op.kind {
		case opWrite:
			v := bytes.Clone(op.value)
			return v, nil
		case opDelete:
			return nil, nil
		default:
			return nil, fmt.Errorf("unknown op kind: %v", op.kind)
		}
	}
	return ot.rawReaderTx.Read(key)
}

type keyOperationPair struct {
	key []byte
	op  operation
}

type overlayIterator struct {
	hiKeyExcl []byte

	currentKey   []byte
	currentValue []byte
	err          error
	done         bool

	sortedTouchedLocalSnapshot []keyOperationPair
	idx                        int // index on sortedTouchedLocalSnapshot

	base     ocr3_1types.KeyValueDatabaseIterator
	baseNext bool
	baseKey  []byte

	onClose func()
}

var _ ocr3_1types.KeyValueDatabaseIterator = &overlayIterator{}

func (oi *overlayIterator) advanceBase() {
	if !oi.base.Next() {
		oi.baseNext = false
		oi.baseKey = nil
		return
	}

	k := oi.base.Key()
	if oi.hiKeyExcl != nil && bytes.Compare(k, oi.hiKeyExcl) >= 0 {
		oi.baseNext = false
		oi.baseKey = nil
		return
	}

	oi.baseNext = true
	oi.baseKey = bytes.Clone(k)
}

func (oi *overlayIterator) Next() bool {
	if oi.done {
		return false
	}

	if oi.err != nil {
		return false
	}

	for {
		var pickedTouchedLocal bool
		if oi.baseNext && oi.idx < len(oi.sortedTouchedLocalSnapshot) {
			// both ranges have more keys, we must compare
			pickedTouchedLocal = bytes.Compare(oi.sortedTouchedLocalSnapshot[oi.idx].key, oi.baseKey) <= 0
		} else if oi.idx < len(oi.sortedTouchedLocalSnapshot) {
			pickedTouchedLocal = true
		} else if oi.baseNext {
			pickedTouchedLocal = false
		} else {
			// both ranges are exhausted
			oi.done = true
			return false
		}

		if pickedTouchedLocal {
			kop := oi.sortedTouchedLocalSnapshot[oi.idx]
			oi.idx++
			// If it was a tie we must advance base to avoid duplicates
			if oi.baseNext && bytes.Equal(kop.key, oi.baseKey) {
				oi.advanceBase()
			}

			if kop.op.kind == opDelete {
				continue
			}
			oi.currentKey = kop.key
			oi.currentValue = bytes.Clone(kop.op.value)
			return true
		}
		// else we picked from rawReaderTx.Range()
		oi.currentKey = bytes.Clone(oi.baseKey)
		v, err := oi.base.Value()
		if err != nil {
			oi.err = err
			oi.done = true
			return false
		}
		oi.currentValue = v
		oi.advanceBase()
		return true
	}
}

func (oi *overlayIterator) Key() []byte {
	return oi.currentKey
}

func (oi *overlayIterator) Value() ([]byte, error) {
	if oi.err != nil {
		return nil, oi.err
	}
	return oi.currentValue, nil
}

func (oi *overlayIterator) Err() error {
	return oi.err
}

func (oi *overlayIterator) Close() error {
	err := oi.base.Close()
	if oi.onClose != nil {
		oi.onClose()
		oi.onClose = nil
	}
	if err != nil {
		return err
	}
	return nil
}

type closedIterator struct {
	err error
}

var _ ocr3_1types.KeyValueDatabaseIterator = &closedIterator{}

func (ti *closedIterator) Next() bool {
	return false
}

func (ti *closedIterator) Key() []byte {
	return nil
}

func (ti *closedIterator) Value() ([]byte, error) {
	return nil, ti.err
}

func (ti *closedIterator) Err() error {
	return ti.err
}

func (ti *closedIterator) Close() error {
	return ti.err
}

// Range iterates over the merged overlay and rawReaderTx keys.
// If a key exists both in the rawReaderTx and the overlay, it returns the overlay value.
// If a key is deleted in the overlay it skips it.
// The caller is expected to call Close() on the returned iterator before Committing or Discarding the transaction.
func (ot *overlayTransaction) Range(loKey []byte, hiKeyExcl []byte) ocr3_1types.KeyValueDatabaseIterator {
	loKey = bytes.Clone(loKey)
	hiKeyExcl = bytes.Clone(hiKeyExcl)

	ot.mu.Lock()
	defer ot.mu.Unlock()
	if ot.status != txStatusOpen {
		return &closedIterator{err: ErrClosed}
	}

	kops := make([]keyOperationPair, 0, len(ot.overlay))
	for k, op := range ot.overlay {
		kb := []byte(k)
		if (loKey == nil || bytes.Compare(kb, loKey) >= 0) &&
			(hiKeyExcl == nil || bytes.Compare(kb, hiKeyExcl) < 0) {
			kops = append(kops, keyOperationPair{bytes.Clone(kb), op})
		}
	}

	sort.Slice(kops,
		func(i, j int) bool {
			return bytes.Compare(kops[i].key, kops[j].key) < 0
		})
	ot.iterators.Add(1)
	oi := &overlayIterator{
		hiKeyExcl,
		nil,
		nil,
		nil,
		false,
		kops,
		0,
		ot.rawReaderTx.Range(loKey, hiKeyExcl),
		false,
		nil,
		ot.iterators.Done,
	}
	oi.advanceBase()
	return oi
}

func (ot *overlayTransaction) lockedDiscard() {
	ot.overlay = nil
	ot.status = txStatusDiscarded
	// wait for any open iterators to close before dropping the read tx
	ot.iterators.Wait()
	ot.rawReaderTx.Discard()
}

func (ot *overlayTransaction) Discard() {
	ot.mu.Lock()
	defer ot.mu.Unlock()
	if ot.status != txStatusOpen {
		return
	}
	ot.lockedDiscard()
}

func (ot *overlayTransaction) Write(key []byte, value []byte) error {
	ot.mu.Lock()
	defer ot.mu.Unlock()
	if ot.status != txStatusOpen {
		return ErrClosed
	}
	ot.overlay[string(key)] = operation{opWrite, bytes.Clone(value)}
	return nil
}

func (ot *overlayTransaction) Delete(key []byte) error {
	ot.mu.Lock()
	defer ot.mu.Unlock()
	if ot.status != txStatusOpen {
		return ErrClosed
	}
	ot.overlay[string(key)] = operation{opDelete, nil}
	return nil
}

func lockedCommit(overlay map[string]operation, rawReadWriteTx ocr3_1types.KeyValueDatabaseReadWriteTransaction) error {
	for k := range overlay {
		op := overlay[k]
		switch op.kind {
		case opWrite:
			if err := rawReadWriteTx.Write([]byte(k), op.value); err != nil {
				return err
			}
		case opDelete:
			if err := rawReadWriteTx.Delete([]byte(k)); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown op kind: %v", op.kind)
		}
	}
	if err := rawReadWriteTx.Commit(); err != nil {
		return err
	}
	return nil
}
