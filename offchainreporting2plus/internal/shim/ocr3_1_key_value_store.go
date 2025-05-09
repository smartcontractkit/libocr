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
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type SemanticOCR3_1KeyValueStoreFactory struct {
	KeyValueDatabaseFactory ocr3_1types.KeyValueDatabaseFactory
}

var _ protocol.KeyValueStoreFactory = &SemanticOCR3_1KeyValueStoreFactory{}

func (s *SemanticOCR3_1KeyValueStoreFactory) NewKeyValueStore(configDigest types.ConfigDigest) (protocol.KeyValueStore, error) {
	keyValueDatabase, err := s.KeyValueDatabaseFactory.NewKeyValueDatabase(configDigest)
	if err != nil {
		return nil, fmt.Errorf("error while creating key value database: %w", err)
	}
	return &SemanticOCR3_1KeyValueStore{keyValueDatabase}, nil
}

type SemanticOCR3_1KeyValueStore struct {
	KeyValueDatabase ocr3_1types.KeyValueDatabase
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
		&SemanticOCR3_1KeyValueStoreReadTransaction{tx},
		tx,
		sync.Mutex{},
		make(map[string][]byte),
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
	return &SemanticOCR3_1KeyValueStoreReadTransaction{tx}, nil
}

type SemanticOCR3_1KeyValueStoreReadWriteTransaction struct {
	protocol.KeyValueStoreReadTransaction // inherit all read implementations

	rawTransaction ocr3_1types.KeyValueReadWriteTransaction

	mu            sync.Mutex
	nilOrWriteSet map[string][]byte
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
	s.mu.Lock()
	if s.nilOrWriteSet == nil {
		s.mu.Unlock()
		return fmt.Errorf("transaction has been discarded")
	}
	s.nilOrWriteSet[string(key)] = nil // tombstone
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
	var writeSet []protocol.KeyValuePair
	s.mu.Lock()
	if s.nilOrWriteSet == nil {
		s.mu.Unlock()
		return nil, fmt.Errorf("transaction has been discarded")
	}
	for k, v := range s.nilOrWriteSet {
		writeSet = append(writeSet, protocol.KeyValuePair{[]byte(k), v})
	}
	s.mu.Unlock()

	sort.Slice(writeSet, func(i, j int) bool {
		return bytes.Compare(writeSet[i].Key, writeSet[j].Key) < 0
	})
	return writeSet, nil
}

func (s *SemanticOCR3_1KeyValueStoreReadWriteTransaction) Write(key []byte, value []byte) error {
	value = bytes.Clone(value) // protect write set against later modification
	value = util.NilCoalesceSlice(value)

	s.mu.Lock()
	if s.nilOrWriteSet == nil {
		s.mu.Unlock()
		return fmt.Errorf("transaction has been discarded")
	}
	s.nilOrWriteSet[string(key)] = value
	s.mu.Unlock()

	return s.rawTransaction.Write(pluginPrefixedKey(key), value)
}

type SemanticOCR3_1KeyValueStoreReadTransaction struct {
	rawTransaction ocr3_1types.KeyValueReadTransaction
}

var _ protocol.KeyValueStoreReadTransaction = &SemanticOCR3_1KeyValueStoreReadTransaction{}

func (s *SemanticOCR3_1KeyValueStoreReadTransaction) Discard() {
	s.rawTransaction.Discard()
}

func (s *SemanticOCR3_1KeyValueStoreReadTransaction) Read(key []byte) ([]byte, error) {
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

const (
	protocolPrefix                 = byte(0)
	pluginPrefix                   = byte(1)
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

func prefixKey(prefix byte, key []byte) []byte {
	return append([]byte{prefix}, key...)
}
