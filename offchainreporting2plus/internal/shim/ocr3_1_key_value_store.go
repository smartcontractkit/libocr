package shim

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
	"sync"

	"github.com/smartcontractkit/libocr/internal/util"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

type SemanticOCR3_1KeyValueStore struct {
	KeyValueDatabase ocr3_1types.KeyValueDatabase
	Limits           ocr3_1types.ReportingPluginLimits
}

var _ protocol.KeyValueStore = &SemanticOCR3_1KeyValueStore{}

func (s *SemanticOCR3_1KeyValueStore) Close() error {
	return s.KeyValueDatabase.Close()
}

func (s *SemanticOCR3_1KeyValueStore) HighestCommittedSeqNr() (uint64, error) {
	tx, err := s.NewReadTransactionUnchecked()
	if err != nil {
		return 0, fmt.Errorf("failed to create read transaction: %w", err)
	}
	defer tx.Discard()
	return tx.ReadHighestCommittedSeqNr()
}

func (s *SemanticOCR3_1KeyValueStore) NewReadWriteTransaction(postSeqNr uint64) (protocol.KeyValueStoreReadWriteTransaction, error) {
	tx, err := s.NewReadWriteTransactionUnchecked()
	if err != nil {
		return nil, fmt.Errorf("failed to create read write transaction: %w", err)
	}
	highestCommittedSeqNr, err := tx.ReadHighestCommittedSeqNr()
	if err != nil {
		tx.Discard()
		return nil, fmt.Errorf("failed to get highest committed seq nr: %w", err)
	}
	if highestCommittedSeqNr+1 != postSeqNr {
		tx.Discard()
		return nil, fmt.Errorf("post seq nr %d must be equal to highest committed seq nr + 1 (%d)", postSeqNr, highestCommittedSeqNr+1)
	}
	return &SemanticOCR3_1KeyValueStoreReadWriteTransactionWithPreCommitHook{
		tx,
		func() error {
			if err := tx.WriteHighestCommittedSeqNr(postSeqNr); err != nil {
				return fmt.Errorf("WriteHighestCommittedSeqNr: %w", err)
			}
			return nil
		},
	}, nil
}

func (s *SemanticOCR3_1KeyValueStore) NewReadWriteTransactionUnchecked() (protocol.KeyValueStoreReadWriteTransaction, error) {
	tx, err := s.KeyValueDatabase.NewReadWriteTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to create read write transaction: %w", err)
	}
	return &SemanticOCR3_1KeyValueStoreReadWriteTransaction{
		&SemanticOCR3_1KeyValueStoreReadTransaction{tx, s.Limits},
		tx,
		sync.Mutex{},
		newLimitCheckWriteSet(s.Limits.MaxKeyValueModifiedKeysPlusValuesLength),
	}, nil
}

func (s *SemanticOCR3_1KeyValueStore) NewReadTransaction(postSeqNr uint64) (protocol.KeyValueStoreReadTransaction, error) {
	tx, err := s.NewReadTransactionUnchecked()
	if err != nil {
		return nil, fmt.Errorf("failed to create read transaction: %w", err)
	}
	highestCommittedSeqNr, err := tx.ReadHighestCommittedSeqNr()
	if err != nil {
		tx.Discard()
		return nil, fmt.Errorf("failed to get highest committed seq nr: %w", err)
	}
	if highestCommittedSeqNr+1 != postSeqNr {
		tx.Discard()
		return nil, fmt.Errorf("post seq nr %d must be equal to highest committed seq nr + 1 (%d)", postSeqNr, highestCommittedSeqNr+1)
	}
	return tx, nil
}

func (s *SemanticOCR3_1KeyValueStore) NewReadTransactionUnchecked() (protocol.KeyValueStoreReadTransaction, error) {
	tx, err := s.KeyValueDatabase.NewReadTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to create read transaction: %w", err)
	}
	return &SemanticOCR3_1KeyValueStoreReadTransaction{tx, s.Limits}, nil
}

type SemanticOCR3_1KeyValueStoreReadWriteTransaction struct {
	protocol.KeyValueStoreReadTransaction // inherit all read implementations

	rawTransaction ocr3_1types.KeyValueReadWriteTransaction

	mu            sync.Mutex
	nilOrWriteSet *limitCheckWriteSet
}

var _ protocol.KeyValueStoreReadWriteTransaction = &SemanticOCR3_1KeyValueStoreReadWriteTransaction{}

type SemanticOCR3_1KeyValueStoreReadWriteTransactionWithPreCommitHook struct {
	protocol.KeyValueStoreReadWriteTransaction
	preCommitHook func() error // must be idempotent
}

var _ protocol.KeyValueStoreReadWriteTransaction = &SemanticOCR3_1KeyValueStoreReadWriteTransactionWithPreCommitHook{}

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransactionWithPreCommitHook) Commit() error {
	if err := s.preCommitHook(); err != nil {
		return fmt.Errorf("failed while executing preCommit: %w", err)
	}
	return s.KeyValueStoreReadWriteTransaction.Commit()
}

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransaction) Commit() error {
	err := s.rawTransaction.Commit()
	// Transactions might persistently fail to commit, due to another txn having
	// gone in before that causes a conflict, so we need to discard in any case
	// to avoid memory leaks.
	s.Discard()
	return err
}

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransaction) Delete(key []byte) error {
	if !(len(key) <= ocr3_1types.MaxMaxKeyValueKeyLength) {
		return fmt.Errorf("key length %d exceeds maximum %d", len(key), ocr3_1types.MaxMaxKeyValueKeyLength)
	}

	s.mu.Lock()
	if s.nilOrWriteSet == nil {
		s.mu.Unlock()
		return fmt.Errorf("transaction has been discarded")
	}
	if err := s.nilOrWriteSet.Delete(key); err != nil {

		s.mu.Unlock()
		return fmt.Errorf("failed to delete key %s from write set: %w", key, err)
	}
	s.mu.Unlock()

	return s.rawTransaction.Delete(pluginPrefixedKey(key))
}

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransaction) Discard() {
	s.mu.Lock()
	s.nilOrWriteSet = nil // tombstone
	s.mu.Unlock()

	s.rawTransaction.Discard()
}

// GetWriteSet returns a map from keys in string encoding to values that have been written in
// this transaction. If the value of a key has been deleted, it is mapped to nil. The write set
// must fit in memory.

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransaction) GetWriteSet() ([]protocol.KeyValuePair, error) {
	s.mu.Lock()
	if s.nilOrWriteSet == nil {
		s.mu.Unlock()
		return nil, fmt.Errorf("transaction has been discarded")
	}
	writeSet := s.nilOrWriteSet.Pairs()
	s.mu.Unlock()

	sort.Slice(writeSet, func(i, j int) bool {
		return bytes.Compare(writeSet[i].Key, writeSet[j].Key) < 0
	})
	return writeSet, nil
}

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransaction) Write(key []byte, value []byte) error {
	if !(len(key) <= ocr3_1types.MaxMaxKeyValueKeyLength) {
		return fmt.Errorf("key length %d exceeds maximum %d", len(key), ocr3_1types.MaxMaxKeyValueKeyLength)
	}
	if !(len(value) <= ocr3_1types.MaxMaxKeyValueValueLength) {
		return fmt.Errorf("value length %d exceeds maximum %d", len(value), ocr3_1types.MaxMaxKeyValueValueLength)
	}

	value = util.NilCoalesceSlice(value)

	s.mu.Lock()
	if s.nilOrWriteSet == nil {
		s.mu.Unlock()
		return fmt.Errorf("transaction has been discarded")
	}
	if err := s.nilOrWriteSet.Write(key, value); err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to write key %s to write set: %w", key, err)
	}
	s.mu.Unlock()

	return s.rawTransaction.Write(pluginPrefixedKey(key), value)
}

type SemanticOCR3_1KeyValueStoreReadTransaction struct {
	rawTransaction ocr3_1types.KeyValueReadTransaction
	limits         ocr3_1types.ReportingPluginLimits
}

var _ protocol.KeyValueStoreReadTransaction = &SemanticOCR3_1KeyValueStoreReadTransaction{}

func (s *SemanticOCR3_1KeyValueStoreReadTransaction) Discard() {
	s.rawTransaction.Discard()
}

func (s *SemanticOCR3_1KeyValueStoreReadTransaction) Read(key []byte) ([]byte, error) {
	if !(len(key) <= ocr3_1types.MaxMaxKeyValueKeyLength) {
		return nil, fmt.Errorf("key length %d exceeds maximum %d", len(key), ocr3_1types.MaxMaxKeyValueKeyLength)
	}
	return s.rawTransaction.Read(pluginPrefixedKey(key))
}

func (s *SemanticOCR3_1KeyValueStoreReadTransaction) ReadHighestCommittedSeqNr() (uint64, error) {
	seqNrRaw, err := s.rawTransaction.Read(highestCommittedSeqNrKey())
	if err != nil {
		return 0, err
	}
	if seqNrRaw == nil { // indicates that we are starting from scratch
		return 0, nil
	}
	if len(seqNrRaw) != 8 {
		return 0, fmt.Errorf("expected 8 bytes for seqNr, got %d", len(seqNrRaw))
	}
	return binary.BigEndian.Uint64(seqNrRaw), nil
}

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransaction) WriteHighestCommittedSeqNr(seqNr uint64) error {
	return s.rawTransaction.Write(highestCommittedSeqNrKey(), binary.BigEndian.AppendUint64(nil, seqNr))
}

func (s *SemanticOCR3_1KeyValueStoreReadTransaction) ReadBlob(blobDigest protocol.BlobDigest) ([]byte, error) {
	var blob []byte

	length, err := s.ReadBlobMeta(blobDigest)
	if err != nil {
		return nil, fmt.Errorf("error reading blob meta for %s: %w", blobDigest, err)
	}
	if length == 0 {

		return nil, nil
	}

	it := s.rawTransaction.Range(blobChunkPrefixedKey(blobDigest), nil)
	defer it.Close()

	residualLength := length

	for i := uint64(0); residualLength > 0 && it.Next(); i++ {
		key := it.Key()
		if !bytes.Equal(key, blobChunkKey(blobDigest, i)) {
			// gap in keys, we're missing a chunk
			return nil, nil
		}

		value, err := it.Value()
		if err != nil {
			return nil, fmt.Errorf("error reading value for key %s: %w", key, err)
		}

		expectedChunkSize := min(protocol.BlobChunkSize, residualLength)
		actualChunkSize := uint64(len(value))
		if actualChunkSize != expectedChunkSize {
			// we don't have the full blob yet
			return nil, nil
		}

		residualLength -= actualChunkSize
		blob = append(blob, value...)
	}

	err = it.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating over blob chunks: %w", err)
	}

	if residualLength != 0 {
		// we somehow don't have the full blob yet, defensive
		return nil, nil
	}

	return blob, nil
}

func (s *SemanticOCR3_1KeyValueStoreReadTransaction) ReadBlobChunk(blobDigest protocol.BlobDigest, chunkIndex uint64) ([]byte, error) {
	return s.rawTransaction.Read(blobChunkKey(blobDigest, chunkIndex))
}

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransaction) WriteBlobChunk(blobDigest protocol.BlobDigest, chunkIndex uint64, chunk []byte) error {
	return s.rawTransaction.Write(blobChunkKey(blobDigest, chunkIndex), chunk)
}

func (s *SemanticOCR3_1KeyValueStoreReadTransaction) ReadBlobMeta(blobDigest protocol.BlobDigest) (uint64, error) {
	lengthBytes, err := s.rawTransaction.Read(blobMetaPrefixKey(blobDigest))
	if err != nil {
		return 0, fmt.Errorf("error reading blob meta for %s: %w", blobDigest, err)
	}
	if lengthBytes == nil {
		// no record of the blob at all
		return 0, nil
	}
	if len(lengthBytes) != 8 {
		return 0, fmt.Errorf("expected 8 bytes for blob meta length, got %d", len(lengthBytes))
	}
	return binary.BigEndian.Uint64(lengthBytes), nil
}

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransaction) WriteBlobMeta(blobDigest protocol.BlobDigest, length uint64) error {
	if length == 0 {
		return fmt.Errorf("cannot write blob meta with length 0 for blob %s", blobDigest)
	}
	lengthBytes := binary.BigEndian.AppendUint64(nil, length)
	return s.rawTransaction.Write(blobMetaPrefixKey(blobDigest), lengthBytes)
}

const (
	protocolPrefix = byte(0)
	pluginPrefix   = byte(1)

	blobChunkSuffix = "blob chunk"
	blobMetaSuffix  = "blob meta"

	highestCommittedSeqNrKeySuffix = "highestCommittedSeqNo"
)

func highestCommittedSeqNrKey() []byte {
	return protocolPrefixedKey([]byte(highestCommittedSeqNrKeySuffix))
}

func pluginPrefixedKey(key []byte) []byte {
	return prefixKey(pluginPrefix, key)
}

func protocolPrefixedKey(key []byte) []byte {
	return prefixKey(protocolPrefix, key)
}

func blobChunkPrefixedKey(blobDigest protocol.BlobDigest) []byte {
	return append(protocolPrefixedKey([]byte(blobChunkSuffix)), blobDigest[:]...)
}

func blobChunkKey(blobDigest protocol.BlobDigest, chunkIndex uint64) []byte {
	chunkIndexBytes := binary.BigEndian.AppendUint64(nil, chunkIndex)
	return append(blobChunkPrefixedKey(blobDigest), chunkIndexBytes...)
}

func blobMetaPrefixKey(blobDigest protocol.BlobDigest) []byte {
	return append(protocolPrefixedKey([]byte(blobMetaSuffix)), blobDigest[:]...)
}

func prefixKey(prefix byte, key []byte) []byte {
	return append([]byte{prefix}, key...)
}

type limitCheckWriteSet struct {
	m                         map[string][]byte
	keysPlusValuesLength      int
	keysPlusValuesLengthLimit int
}

func newLimitCheckWriteSet(keysPlusValuesLengthLimit int) *limitCheckWriteSet {
	return &limitCheckWriteSet{
		make(map[string][]byte),
		0,
		keysPlusValuesLengthLimit,
	}
}

func (l *limitCheckWriteSet) modify(key []byte, value []byte) error {

	add, sub := 0, 0
	if prevValue, ok := l.m[string(key)]; ok {
		add = len(value)
		sub = len(prevValue)
	} else {
		if len(key)+len(value) < len(key) {
			return fmt.Errorf("key + value length overflow")
		}
		add = len(key) + len(value)
	}

	keysPlusValuesLengthMinusExistingValue := l.keysPlusValuesLength - sub
	if keysPlusValuesLengthMinusExistingValue+add < keysPlusValuesLengthMinusExistingValue {
		return fmt.Errorf("keys + values length overflow")
	}
	if keysPlusValuesLengthMinusExistingValue+add > l.keysPlusValuesLengthLimit {
		return fmt.Errorf("keys + values length %d exceeds limit %d", keysPlusValuesLengthMinusExistingValue+add, l.keysPlusValuesLengthLimit)
	}
	l.m[string(key)] = bytes.Clone(value)
	l.keysPlusValuesLength = keysPlusValuesLengthMinusExistingValue + add
	return nil
}

func (l *limitCheckWriteSet) Write(key []byte, value []byte) error {
	return l.modify(key, value)
}

func (l *limitCheckWriteSet) Delete(key []byte) error {
	return l.modify(key, nil)
}

func (l *limitCheckWriteSet) Pairs() []protocol.KeyValuePair {
	pairs := make([]protocol.KeyValuePair, 0, len(l.m))
	for k, v := range l.m {
		pairs = append(pairs, protocol.KeyValuePair{
			[]byte(k),
			v,
			v == nil,
		})
	}
	return pairs
}
