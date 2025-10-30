package shim

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/jmt"
	"github.com/smartcontractkit/libocr/internal/singlewriter"
	"github.com/smartcontractkit/libocr/internal/util"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

type SemanticOCR3_1KeyValueDatabase struct {
	conflictTracker  *singlewriter.ConflictTracker
	KeyValueDatabase ocr3_1types.KeyValueDatabase
	Limits           ocr3_1types.ReportingPluginLimits
	config           ocr3_1config.PublicConfig
	logger           commontypes.Logger
	metrics          *semanticKeyValueDatabaseMetrics
}

var _ protocol.KeyValueDatabase = &SemanticOCR3_1KeyValueDatabase{}

func NewSemanticOCR3_1KeyValueDatabase(
	keyValueDatabase ocr3_1types.KeyValueDatabase,
	limits ocr3_1types.ReportingPluginLimits,
	config ocr3_1config.PublicConfig,
	logger commontypes.Logger,
	metricsRegisterer prometheus.Registerer,
) (*SemanticOCR3_1KeyValueDatabase, error) {
	if err := initializeSchema(keyValueDatabase); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}
	return &SemanticOCR3_1KeyValueDatabase{
		singlewriter.NewConflictTracker(),
		keyValueDatabase,
		limits,
		config,
		logger,
		newSemanticKeyValueDatabaseMetrics(metricsRegisterer, logger),
	}, nil
}

func (s *SemanticOCR3_1KeyValueDatabase) newReadWriteTransaction(tx ocr3_1types.KeyValueDatabaseReadWriteTransaction, nilOrSeqNr *uint64) *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction {
	return &SemanticOCR3_1KeyValueDatabaseReadWriteTransaction{
		SemanticOCR3_1KeyValueDatabaseReadTransaction{
			tx,
			s.config,
		},
		tx,
		s.metrics,
		sync.Mutex{},
		newLimitCheckWriteSet(s.Limits),
		nilOrSeqNr,
		false,
	}
}

// SemanticOCR3_1KeyValueDatabase does not manage the lifetime of the underlying
// database, and is expected to be closed first, before the database is closed.
func (s *SemanticOCR3_1KeyValueDatabase) Close() error {
	s.metrics.Close()
	return nil
}

func (s *SemanticOCR3_1KeyValueDatabase) HighestCommittedSeqNr() (uint64, error) {
	tx, err := s.NewReadTransactionUnchecked()
	if err != nil {
		return 0, fmt.Errorf("failed to create read transaction: %w", err)
	}
	defer tx.Discard()
	return tx.ReadHighestCommittedSeqNr()
}

func (s *SemanticOCR3_1KeyValueDatabase) NewSerializedReadWriteTransaction(postSeqNr uint64) (protocol.KeyValueDatabaseReadWriteTransaction, error) {
	rawTx, err := singlewriter.NewSerializedTransaction(s.KeyValueDatabase, s.conflictTracker)
	if err != nil {
		return nil, fmt.Errorf("failed to create read write transaction: %w", err)
	}
	tx := s.newReadWriteTransaction(rawTx, &postSeqNr)
	highestCommittedSeqNr, err := tx.ReadHighestCommittedSeqNr()
	if err != nil {
		tx.Discard()
		return nil, fmt.Errorf("failed to get highest committed seq nr: %w", err)
	}
	if highestCommittedSeqNr+1 != postSeqNr {
		tx.Discard()
		return nil, fmt.Errorf("post seq nr %d must be equal to highest committed seq nr + 1 (%d)", postSeqNr, highestCommittedSeqNr+1)
	}
	if err := checkNotClobbered(tx); err != nil {
		tx.Discard()
		return nil, err
	}
	return &SemanticOCR3_1KeyValueDatabaseReadWriteTransactionWithPreCommitHook{
		tx,
		func() error {
			if err := tx.WriteHighestCommittedSeqNr(postSeqNr); err != nil {
				return fmt.Errorf("WriteHighestCommittedSeqNr: %w", err)
			}
			return nil
		},
	}, nil
}

func (s *SemanticOCR3_1KeyValueDatabase) NewSerializedReadWriteTransactionUnchecked() (protocol.KeyValueDatabaseReadWriteTransaction, error) {
	rawTx, err := singlewriter.NewSerializedTransaction(s.KeyValueDatabase, s.conflictTracker)
	if err != nil {
		return nil, fmt.Errorf("failed to create read write transaction: %w", err)
	}
	return s.newReadWriteTransaction(rawTx, nil), nil
}

func (s *SemanticOCR3_1KeyValueDatabase) NewUnserializedReadWriteTransactionUnchecked() (protocol.KeyValueDatabaseReadWriteTransaction, error) {
	rawTx, err := singlewriter.NewUnserializedTransaction(s.KeyValueDatabase)
	if err != nil {
		return nil, fmt.Errorf("failed to create read write transaction: %w", err)
	}
	return s.newReadWriteTransaction(rawTx, nil), nil
}

func (s *SemanticOCR3_1KeyValueDatabase) NewReadTransaction(postSeqNr uint64) (protocol.KeyValueDatabaseReadTransaction, error) {
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
	if err := checkNotClobbered(tx); err != nil {
		tx.Discard()
		return nil, err
	}
	return tx, nil
}

func checkNotClobbered(tx protocol.KeyValueDatabaseReadTransaction) error {
	treeSyncStatus, err := tx.ReadTreeSyncStatus()
	if err != nil {
		return fmt.Errorf("failed to read tree sync status: %w", err)
	}
	if treeSyncStatus.Phase != protocol.TreeSyncPhaseInactive {
		return fmt.Errorf("tree sync might be in progress")
	}
	return nil
}

func (s *SemanticOCR3_1KeyValueDatabase) NewReadTransactionUnchecked() (protocol.KeyValueDatabaseReadTransaction, error) {
	tx, err := s.KeyValueDatabase.NewReadTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to create read transaction: %w", err)
	}
	return &SemanticOCR3_1KeyValueDatabaseReadTransaction{tx, s.config}, nil
}

type SemanticOCR3_1KeyValueDatabaseReadWriteTransaction struct {
	SemanticOCR3_1KeyValueDatabaseReadTransaction // inherit all read implementations
	rawTransaction                                ocr3_1types.KeyValueDatabaseReadWriteTransaction
	metrics                                       *semanticKeyValueDatabaseMetrics
	mu                                            sync.Mutex
	nilOrWriteSet                                 *limitCheckWriteSet
	nilOrSeqNr                                    *uint64
	closedForWriting                              bool
}

var _ protocol.KeyValueDatabaseReadWriteTransaction = &SemanticOCR3_1KeyValueDatabaseReadWriteTransaction{}

type SemanticOCR3_1KeyValueDatabaseReadWriteTransactionWithPreCommitHook struct {
	protocol.KeyValueDatabaseReadWriteTransaction
	preCommitHook func() error // must be idempotent
}

var _ protocol.KeyValueDatabaseReadWriteTransaction = &SemanticOCR3_1KeyValueDatabaseReadWriteTransactionWithPreCommitHook{}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransactionWithPreCommitHook) Commit() error {
	if err := s.preCommitHook(); err != nil {
		return fmt.Errorf("failed while executing preCommit: %w", err)
	}
	return s.KeyValueDatabaseReadWriteTransaction.Commit()
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) Commit() error {
	start := time.Now()
	defer func() {
		s.metrics.txCommitDurationSeconds.Observe(float64(time.Since(start).Seconds()))
	}()

	err := s.rawTransaction.Commit()
	// Transactions might persistently fail to commit, due to another txn having
	// gone in before that causes a conflict, so we need to discard in any case
	// to avoid memory leaks.
	s.Discard()
	return err
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) Delete(key []byte) error {
	if !(len(key) <= ocr3_1types.MaxMaxKeyValueKeyBytes) {
		return fmt.Errorf("key length %d exceeds maximum %d", len(key), ocr3_1types.MaxMaxKeyValueKeyBytes)
	}

	s.mu.Lock()
	if s.closedForWriting {
		s.mu.Unlock()
		return fmt.Errorf("transaction has been closed for writing")
	}
	if s.nilOrWriteSet == nil {
		s.mu.Unlock()
		return fmt.Errorf("transaction has been discarded")
	}
	if err := s.nilOrWriteSet.Delete(key); err != nil {

		s.mu.Unlock()
		return fmt.Errorf("failed to delete key %s from write set: %w", key, err)
	}
	s.mu.Unlock()
	return s.rawTransaction.Delete(pluginPrefixedUnhashedKey(key))
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) deletePrefixedKeys(prefix []byte, n int) (done bool, numDeleted int, err error) {
	// We cannot delete the keys while iterating them, if we want to be agnostic
	// to kvdb implementation semantics.
	var keysToDelete [][]byte

	it := s.rawTransaction.Range(prefix, nil)
	for it.Next() && len(keysToDelete) < n+1 {
		key := it.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		keysToDelete = append(keysToDelete, key)
	}
	if err := it.Err(); err != nil {
		it.Close()
		return false, 0, fmt.Errorf("failed to range: %w", err)
	}
	it.Close()

	for _, key := range keysToDelete {
		if !(numDeleted+1 <= n) {
			break
		}
		if err := s.rawTransaction.Delete(key); err != nil {
			return false, numDeleted, fmt.Errorf("failed to delete key %s: %w", key, err)
		}
		numDeleted++
	}

	return len(keysToDelete) <= n, numDeleted, nil
}

// Caller must ensure to make committed state inaccessible to other transactions
// until completed. Must be reinvoked until done=true.
func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) DestructiveDestroyForTreeSync(n int) (done bool, err error) {
	prefixesToDelete := destructiveDestroyForTreeSyncPrefixesToDelete

	prefixesDone := make([]bool, len(prefixesToDelete))
	budget := n
	for i, prefix := range prefixesToDelete {
		done, numDeleted, err := s.deletePrefixedKeys(prefix, budget)
		if err != nil {
			return false, fmt.Errorf("failed to delete prefixed keys for prefix %q: %w", prefix, err)
		}
		budget -= numDeleted
		if done {
			prefixesDone[i] = true
		} else if budget == 0 {
			break
		}
	}

	allPrefixesDone := true
	for _, done := range prefixesDone {
		if !done {
			allPrefixesDone = false
			break
		}
	}
	return allPrefixesDone, nil
}

// Helper for reaping methods that require large ranges over multiple transactions
func partialExclusiveRangeKeys(readTransaction ocr3_1types.KeyValueDatabaseReadTransaction, loKey []byte, hiKeyExcl []byte, maxItems int) (keys [][]byte, more bool, err error) {
	it := readTransaction.Range(loKey, hiKeyExcl)
	defer it.Close()

	for it.Next() {
		if len(keys) == maxItems {
			more = true
			break
		}
		keys = append(keys, it.Key())
	}
	if err := it.Err(); err != nil {
		return nil, false, fmt.Errorf("failed to range: %w", err)
	}
	return keys, more, nil
}

func partialInclusiveRangeKeys(readTransaction ocr3_1types.KeyValueDatabaseReadTransaction, loKey []byte, hiKeyIncl []byte, maxItems int) (keys [][]byte, more bool, err error) {
	hiKeyExcl := append(bytes.Clone(hiKeyIncl), 0)
	return partialExclusiveRangeKeys(readTransaction, loKey, hiKeyExcl, maxItems)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) partialExclusiveRangeKeys(loKey []byte, hiKeyExcl []byte, maxItems int) (keys [][]byte, more bool, err error) {
	return partialExclusiveRangeKeys(s.rawTransaction, loKey, hiKeyExcl, maxItems)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) partialInclusiveRangeKeys(loKey []byte, hiKeyIncl []byte, maxItems int) (keys [][]byte, more bool, err error) {
	return partialInclusiveRangeKeys(s.rawTransaction, loKey, hiKeyIncl, maxItems)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) partialInclusiveRangeKeys(loKey []byte, hiKeyIncl []byte, maxItems int) (keys [][]byte, more bool, err error) {
	return partialInclusiveRangeKeys(s.rawTransaction, loKey, hiKeyIncl, maxItems)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) Discard() {
	s.mu.Lock()
	s.nilOrWriteSet = nil // tombstone
	s.mu.Unlock()

	s.rawTransaction.Discard()
}

// GetWriteSet returns sorted list of key-value pairs that have been modified as
// part of this transaction.
func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) GetWriteSet() ([]protocol.KeyValuePairWithDeletions, error) {
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

// CloseWriteSet updates the state tree according to the write set and returns
// the root. After this function is invoked the transaction for writing: any
// future attempts for Writes or Deletes on this transaction will fail.
func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) CloseWriteSet() (protocol.StateRootDigest, error) {
	start := time.Now()
	defer func() {
		s.metrics.closeWriteSetDurationSeconds.Observe(float64(time.Since(start).Seconds()))
	}()

	s.mu.Lock()
	if s.nilOrWriteSet == nil {
		s.mu.Unlock()
		return protocol.StateRootDigest{}, fmt.Errorf("transaction has been discarded")
	}
	writeSet := s.nilOrWriteSet.Pairs()
	s.nilOrWriteSet = nil
	s.closedForWriting = true
	s.mu.Unlock()

	if s.nilOrSeqNr == nil {
		return protocol.StateRootDigest{}, fmt.Errorf("transaction seqNr should not be nil")
	}

	keyValueUpdates := make([]jmt.KeyValue, 0, len(writeSet))
	for _, pair := range writeSet {
		var value []byte
		if !pair.Deleted {
			value = pair.Value
		}
		keyValueUpdates = append(keyValueUpdates, jmt.KeyValue{
			pair.Key,
			value,
		})
	}

	_, err := jmt.BatchUpdate(
		s,
		s,
		s,
		protocol.PrevRootVersion(*s.nilOrSeqNr, s.config),
		protocol.RootVersion(*s.nilOrSeqNr, s.config),
		keyValueUpdates,
	)
	if err != nil {
		return protocol.StateRootDigest{}, fmt.Errorf("failed to batch update: %w", err)
	}

	stateRootDigest, err := jmt.ReadRootDigest(s, s, protocol.RootVersion(*s.nilOrSeqNr, s.config))
	if err != nil {
		return protocol.StateRootDigest{}, fmt.Errorf("failed to read root digest: %w", err)
	}
	return stateRootDigest, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) ApplyWriteSet(writeSet []protocol.KeyValuePairWithDeletions) (protocol.StateRootDigest, error) {
	if s.nilOrSeqNr == nil {
		return protocol.StateRootDigest{}, fmt.Errorf("transaction seqNr should not be nil")
	}
	seqNr := *s.nilOrSeqNr
	for i, m := range writeSet {
		var err error
		switch m.Deleted {
		case false:
			err = s.Write(m.Key, m.Value)
		case true:
			err = s.Delete(m.Key)
		}
		if err != nil {
			return protocol.StateRootDigest{}, fmt.Errorf("failed to write %d-th write-set modification for seq nr %d: %w", i, seqNr, err)
		}
	}
	stateRootDigest, err := s.CloseWriteSet()
	if err != nil {
		return protocol.StateRootDigest{}, fmt.Errorf("failed to close write set: %w", err)
	}
	return stateRootDigest, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) Write(key []byte, value []byte) error {
	start := time.Now()
	defer func() {
		s.metrics.txWriteDurationSeconds.Observe(float64(time.Since(start).Seconds()))
	}()

	if !(len(key) <= ocr3_1types.MaxMaxKeyValueKeyBytes) {
		return fmt.Errorf("key length %d exceeds maximum %d", len(key), ocr3_1types.MaxMaxKeyValueKeyBytes)
	}
	if !(len(value) <= ocr3_1types.MaxMaxKeyValueValueBytes) {
		return fmt.Errorf("value length %d exceeds maximum %d", len(value), ocr3_1types.MaxMaxKeyValueValueBytes)
	}

	value = util.NilCoalesceSlice(value)

	s.mu.Lock()
	if s.closedForWriting {
		s.mu.Unlock()
		return fmt.Errorf("transaction has been closed for writing")
	}
	if s.nilOrWriteSet == nil {
		s.mu.Unlock()
		return fmt.Errorf("transaction has been discarded")
	}
	if err := s.nilOrWriteSet.Write(key, value); err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to write key %s to write set: %w", key, err)
	}
	s.mu.Unlock()

	err := s.rawTransaction.Write(pluginPrefixedUnhashedKey(key), value)
	if err != nil {
		return fmt.Errorf("failed to write key %s to write set: %w", key, err)
	}
	return nil
}

type SemanticOCR3_1KeyValueDatabaseReadTransaction struct {
	rawTransaction ocr3_1types.KeyValueDatabaseReadTransaction
	config         ocr3_1config.PublicConfig
}

var _ protocol.KeyValueDatabaseReadTransaction = &SemanticOCR3_1KeyValueDatabaseReadTransaction{}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) Discard() {
	s.rawTransaction.Discard()
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) Read(key []byte) ([]byte, error) {
	if !(len(key) <= ocr3_1types.MaxMaxKeyValueKeyBytes) {
		return nil, fmt.Errorf("key length %d exceeds maximum %d", len(key), ocr3_1types.MaxMaxKeyValueKeyBytes)
	}
	return s.rawTransaction.Read(pluginPrefixedUnhashedKey(key))
}

func readUint64ValueOrZero(raw []byte) (uint64, error) {
	if raw == nil {
		return 0, nil
	}
	if len(raw) != 8 {
		return 0, fmt.Errorf("expected 8 bytes for seqNr, got %d", len(raw))
	}
	return binary.BigEndian.Uint64(raw), nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadHighestCommittedSeqNr() (uint64, error) {
	seqNrRaw, err := s.rawTransaction.Read([]byte(highestCommittedSeqNrKey))
	if err != nil {
		return 0, err
	}
	return readUint64ValueOrZero(seqNrRaw)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadLowestPersistedSeqNr() (uint64, error) {
	seqNrRaw, err := s.rawTransaction.Read([]byte(lowestPersistedSeqNrKey))
	if err != nil {
		return 0, err
	}
	return readUint64ValueOrZero(seqNrRaw)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadAttestedStateTransitionBlock(seqNr uint64) (protocol.AttestedStateTransitionBlock, error) {
	blockRaw, err := s.rawTransaction.Read(blockKey(seqNr))
	if err != nil {
		return protocol.AttestedStateTransitionBlock{}, err
	}
	if blockRaw == nil {
		return protocol.AttestedStateTransitionBlock{}, nil
	}
	block, err := serialization.DeserializeAttestedStateTransitionBlock(blockRaw)
	if err != nil {
		return protocol.AttestedStateTransitionBlock{}, fmt.Errorf("failed to deserialize attested state transition block %d: %w", seqNr, err)
	}
	return block, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadAttestedStateTransitionBlocks(minSeqNr uint64, maxItems int) (blocks []protocol.AttestedStateTransitionBlock, more bool, err error) {
	blockKeys, more, err := s.partialInclusiveRangeKeys(blockKey(minSeqNr), blockKey(math.MaxUint64), maxItems)
	if err != nil {
		return nil, false, fmt.Errorf("failed to range: %w", err)
	}

	for _, blockKey := range blockKeys {
		seqNr, err := deserializeBlockKey(blockKey)
		if err != nil {
			return nil, false, fmt.Errorf("failed to deserialize block key: %w", err)
		}
		block, err := s.ReadAttestedStateTransitionBlock(seqNr)
		if err != nil {
			return nil, false, fmt.Errorf("failed to read attested state transition block %d: %w", seqNr, err)
		}
		blocks = append(blocks, block)
	}
	return blocks, more, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteAttestedStateTransitionBlock(seqNr uint64, block protocol.AttestedStateTransitionBlock) error {
	blockBytes, err := serialization.SerializeAttestedStateTransitionBlock(block)
	if err != nil {
		return fmt.Errorf("failed to serialize attested state transition block %d: %w", seqNr, err)
	}
	return s.rawTransaction.Write(blockKey(seqNr), blockBytes)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) DeleteAttestedStateTransitionBlocks(maxSeqNrToDelete uint64, maxItems int) (done bool, err error) {
	keys, more, err := s.partialInclusiveRangeKeys(blockKey(0), blockKey(maxSeqNrToDelete), maxItems)
	if err != nil {
		return false, fmt.Errorf("failed to range: %w", err)
	}
	for _, key := range keys {
		if err := s.rawTransaction.Delete(key); err != nil {
			return false, fmt.Errorf("failed to delete key %s: %w", key, err)
		}
	}
	return !more, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadTreeSyncStatus() (protocol.TreeSyncStatus, error) {
	statusRaw, err := s.rawTransaction.Read([]byte(treeSyncStatusKey))
	if err != nil {
		return protocol.TreeSyncStatus{}, err
	}
	if statusRaw == nil {
		return protocol.TreeSyncStatus{}, nil
	}
	return serialization.DeserializeTreeSyncStatus(statusRaw)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadTreeSyncChunk(
	toSeqNr uint64,
	startIndex jmt.Digest,
	requestEndInclIndex jmt.Digest,
	maxCumulativeKeysPlusValuesBytes int,
) (
	endInclIndex jmt.Digest,
	boundingLeaves []jmt.BoundingLeaf,
	keyValues []protocol.KeyValuePair,
	err error,
) {
	if !(0 < toSeqNr) {
		return jmt.Digest{}, nil, nil, fmt.Errorf("toSeqNr (%d) must be > 0", toSeqNr)
	}

	highestCommittedSeqNr, err := s.ReadHighestCommittedSeqNr()
	if err != nil {
		return jmt.Digest{}, nil, nil, fmt.Errorf("failed to read highest committed seq nr")
	}

	lowestPersistedSeqNr, err := s.ReadLowestPersistedSeqNr()
	if err != nil {
		return jmt.Digest{}, nil, nil, fmt.Errorf("failed to read lowest persisted seq nr")
	}

	if !(lowestPersistedSeqNr <= toSeqNr && toSeqNr <= highestCommittedSeqNr) {
		return jmt.Digest{}, nil, nil, fmt.Errorf("toSeqNr (%d) must be >= lowest persisted seq nr (%d) and <= highest committed seq nr (%d)", toSeqNr, lowestPersistedSeqNr, highestCommittedSeqNr)
	}

	keyValues, truncated, err := jmt.ReadRange(
		s,
		s,
		protocol.RootVersion(toSeqNr, s.config),
		startIndex,
		requestEndInclIndex,
		maxCumulativeKeysPlusValuesBytes,
		s.config.GetMaxTreeSyncChunkKeys(),
	)
	if err != nil {
		return jmt.Digest{}, nil, nil, fmt.Errorf("failed to read range: %w", err)
	}

	if truncated {
		if len(keyValues) == 0 {
			return jmt.Digest{}, nil, nil, fmt.Errorf("read range could not even fit a single kv pair in required limits, the limits are probably wrong")
		}
		endInclIndex = jmt.DigestKey(keyValues[len(keyValues)-1].Key)
	} else {
		endInclIndex = requestEndInclIndex
	}

	boundingLeaves, err = jmt.ProveSubrange(
		s,
		s,
		protocol.RootVersion(toSeqNr, s.config),
		startIndex,
		endInclIndex,
	)
	if err != nil {
		return jmt.Digest{}, nil, nil, fmt.Errorf("failed to prove range: %w", err)
	}

	return
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteHighestCommittedSeqNr(seqNr uint64) error {
	preHighestCommittedSeqNr, err := s.ReadHighestCommittedSeqNr()
	if err != nil {
		return fmt.Errorf("failed to read highest committed seq nr: %w", err)
	}
	if preHighestCommittedSeqNr > seqNr {
		return fmt.Errorf("pre highest committed seq nr %d must be <= highest committed seq nr %d", preHighestCommittedSeqNr, seqNr)
	}
	return s.rawTransaction.Write([]byte(highestCommittedSeqNrKey), encodeBigEndianUint64(seqNr))
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteLowestPersistedSeqNr(seqNr uint64) error {
	preLowestPersistedSeqNr, err := s.ReadLowestPersistedSeqNr()
	if err != nil {
		return fmt.Errorf("failed to read lowest persisted seq nr: %w", err)
	}
	if seqNr < preLowestPersistedSeqNr {
		return fmt.Errorf("pre lowest persisted seq nr %d must be <= lowest persisted seq nr %d", preLowestPersistedSeqNr, seqNr)
	}
	return s.rawTransaction.Write([]byte(lowestPersistedSeqNrKey), encodeBigEndianUint64(seqNr))
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteTreeSyncStatus(status protocol.TreeSyncStatus) error {
	rawStatus, err := serialization.SerializeTreeSyncStatus(status)
	if err != nil {
		return err
	}
	return s.rawTransaction.Write([]byte(treeSyncStatusKey), rawStatus)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) VerifyAndWriteTreeSyncChunk(
	targetRootDigest protocol.StateRootDigest,
	targetSeqNr uint64,
	startIndex jmt.Digest,
	endInclIndex jmt.Digest,
	boundingLeaves []jmt.BoundingLeaf,
	keyValues []protocol.KeyValuePair,
) (protocol.VerifyAndWriteTreeSyncChunkResult, error) {
	if len(keyValues) > s.config.GetMaxTreeSyncChunkKeys() {
		return protocol.VerifyAndWriteTreeSyncChunkResultByzantine, fmt.Errorf("too many leaves: %d > %d",
			len(keyValues), s.config.GetMaxTreeSyncChunkKeys())
	}
	var byteBudget int
	for _, kv := range keyValues {
		byteBudget += len(kv.Key) + len(kv.Value)
	}
	if byteBudget > s.config.GetMaxTreeSyncChunkKeysPlusValuesBytes() {
		return protocol.VerifyAndWriteTreeSyncChunkResultByzantine, fmt.Errorf("chunk exceeds byte limit: %d > %d",
			byteBudget, s.config.GetMaxTreeSyncChunkKeysPlusValuesBytes())
	}

	prevIdx := startIndex
	for i, kv := range keyValues {
		if kv.Value == nil {
			return protocol.VerifyAndWriteTreeSyncChunkResultByzantine, fmt.Errorf("leaf %v has nil value", kv)
		}
		idx := hashPluginKey(kv.Key)
		if bytes.Compare(idx[:], startIndex[:]) < 0 {
			return protocol.VerifyAndWriteTreeSyncChunkResultByzantine, fmt.Errorf("index of leaf %d out of chunk range, want index >= startIndex:%x got index:%x", i, startIndex, idx)
		}
		if bytes.Compare(idx[:], endInclIndex[:]) > 0 {
			return protocol.VerifyAndWriteTreeSyncChunkResultByzantine, fmt.Errorf("index of leaf %d out of chunk range, want index <= endInclIndex:%x got index:%x", i, endInclIndex, idx)
		}
		if i > 0 && bytes.Compare(idx[:], prevIdx[:]) <= 0 {
			return protocol.VerifyAndWriteTreeSyncChunkResultByzantine, fmt.Errorf("leaves not strictly ascending")
		}
		prevIdx = idx
	}

	// verify subrange proof
	{
		err := jmt.VerifySubrange(
			targetRootDigest,
			startIndex,
			endInclIndex,
			keyValues,
			boundingLeaves,
		)
		if err != nil {
			return protocol.VerifyAndWriteTreeSyncChunkResultByzantine, fmt.Errorf("invalid subrange proof: %w", err)
		}
	}

	// apply the updates as indicated by the leaves
	{
		_, err := jmt.BatchUpdate(
			s,
			s,
			s,
			protocol.RootVersion(targetSeqNr, s.config),
			protocol.RootVersion(targetSeqNr, s.config),
			keyValues,
		)
		if err != nil {
			return protocol.VerifyAndWriteTreeSyncChunkResultUnrelatedError, fmt.Errorf("failed to batch update: %w", err)
		}
	}

	// write flat representation
	for _, kv := range keyValues {
		err := s.rawTransaction.Write(pluginPrefixedUnhashedKey(kv.Key), kv.Value)
		if err != nil {
			return protocol.VerifyAndWriteTreeSyncChunkResultUnrelatedError, fmt.Errorf("could not write the key-value pair to store: %w", err)
		}
	}

	rootDigest, err := jmt.ReadRootDigest(
		s,
		s,
		protocol.RootVersion(targetSeqNr, s.config),
	)
	if err != nil {
		return protocol.VerifyAndWriteTreeSyncChunkResultUnrelatedError, fmt.Errorf("failed to read root digest: %w", err)
	}

	if rootDigest == targetRootDigest {
		return protocol.VerifyAndWriteTreeSyncChunkResultOkComplete, nil
	}
	return protocol.VerifyAndWriteTreeSyncChunkResultOkNeedMore, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadBlobPayload(blobDigest protocol.BlobDigest) ([]byte, error) {
	blobMeta, err := s.ReadBlobMeta(blobDigest)
	if err != nil {
		return nil, fmt.Errorf("error reading blob meta for %s: %w", blobDigest, err)
	}
	if blobMeta == nil {
		return nil, nil
	}
	if slices.Contains(blobMeta.ChunkHaves, false) {
		return nil, fmt.Errorf("blob has missing chunks")
	}

	highestCommittedSeqNr, err := s.ReadHighestCommittedSeqNr()
	if err != nil {
		return nil, fmt.Errorf("error reading highest committed seq nr: %w", err)
	}
	if blobMeta.ExpirySeqNr < highestCommittedSeqNr {
		return nil, fmt.Errorf("blob has expired")
	}

	it := s.rawTransaction.Range(blobChunkPrefixedKey(blobDigest), nil)
	defer it.Close()

	residualLength := blobMeta.PayloadLength
	payload := make([]byte, 0, residualLength)

	for i := uint64(0); residualLength > 0 && it.Next(); i++ {
		key := it.Key()
		if !bytes.Equal(key, blobChunkKey(blobDigest, i)) {
			return nil, fmt.Errorf("unexpected key for %v-th chunk: %x", i, key)
		}

		value, err := it.Value()
		if err != nil {
			return nil, fmt.Errorf("error reading value for key %s: %w", key, err)
		}

		expectedChunkSize := min(uint64(s.config.GetBlobChunkBytes()), residualLength)
		actualChunkSize := uint64(len(value))
		if actualChunkSize != expectedChunkSize {
			return nil, fmt.Errorf("actual chunk size %v != expected chunk size %v", actualChunkSize, expectedChunkSize)
		}

		residualLength -= actualChunkSize
		payload = append(payload, value...)
	}

	err = it.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating over blob chunks: %w", err)
	}

	if residualLength != 0 {
		return nil, fmt.Errorf("residual length %v != 0 even though we have all chunks", residualLength)
	}

	return payload, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadBlobChunk(blobDigest protocol.BlobDigest, chunkIndex uint64) ([]byte, error) {
	return s.rawTransaction.Read(blobChunkKey(blobDigest, chunkIndex))
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadNode(nodeKey jmt.NodeKey) (jmt.Node, error) {
	rawNode, err := s.rawTransaction.Read(treePrefixedKey(nodeKey))
	if err != nil {
		return nil, fmt.Errorf("failed to read jmt node: %w", err)
	}
	if rawNode == nil {

		return nil, nil
	}
	return serialization.DeserializeJmtNode(rawNode)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadRoot(version jmt.Version) (jmt.NodeKey, error) {
	rawNodeKey, err := s.rawTransaction.Read(rootKey(version))
	if err != nil {
		return jmt.NodeKey{}, fmt.Errorf("failed to read jmt root: %w", err)
	}
	if rawNodeKey == nil {
		return jmt.NodeKey{}, nil
	}
	return serialization.DeserializeNodeKey(rawNodeKey)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteNode(nodeKey jmt.NodeKey, nodeOrNil jmt.Node) error {
	if nodeOrNil == nil {
		return s.rawTransaction.Delete(treePrefixedKey(nodeKey))
	}

	rawNode, err := serialization.SerializeJmtNode(nodeOrNil)
	if err != nil {
		return fmt.Errorf("failed to serialize jmt node: %w", err)
	}
	return s.rawTransaction.Write(treePrefixedKey(nodeKey), rawNode)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteRoot(version jmt.Version, nodeKey jmt.NodeKey) error {
	return s.rawTransaction.Write(rootKey(version), serialization.AppendSerializeNodeKey(nil, nodeKey))
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) DeleteRoots(minVersionToKeep jmt.Version, maxItems int) (done bool, err error) {
	keys, more, err := s.partialExclusiveRangeKeys(rootKey(0), rootKey(minVersionToKeep), maxItems)
	if err != nil {
		return false, fmt.Errorf("failed to range: %w", err)
	}
	for _, key := range keys {
		if err := s.rawTransaction.Delete(key); err != nil {
			return false, fmt.Errorf("failed to delete key %s: %w", key, err)
		}
	}
	return !more, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteStaleNode(staleNode jmt.StaleNode) error {
	return s.rawTransaction.Write(stalePrefixedKey(staleNode), nil)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) DeleteStaleNodes(maxStaleSinceVersion jmt.Version, maxItems int) (done bool, err error) {
	if maxStaleSinceVersion+1 < maxStaleSinceVersion {
		return false, fmt.Errorf("maxStaleSinceVersion overflow")
	}

	staleIndexNodeKeys, more, err := s.partialExclusiveRangeKeys(
		staleKeyWithStaleSinceVersionBase(0),
		staleKeyWithStaleSinceVersionBase(maxStaleSinceVersion+1),
		maxItems,
	)
	if err != nil {
		return false, fmt.Errorf("failed to range: %w", err)
	}

	for _, staleIndexNodeKey := range staleIndexNodeKeys {
		staleNode, err := deserializeStalePrefixedKey(staleIndexNodeKey)
		if err != nil {
			return false, fmt.Errorf("failed to deserialize stale node: %w", err)
		}

		err = s.WriteNode(staleNode.NodeKey, nil)
		if err != nil {
			return false, fmt.Errorf("error writing node %v: %w", staleNode.NodeKey, err)
		}
		err = s.deleteStaleNode(staleNode)
		if err != nil {
			return false, fmt.Errorf("error deleting stale node %v: %w", staleNode.NodeKey, err)
		}
	}
	return !more, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) deleteStaleNode(staleNode jmt.StaleNode) error {
	return s.rawTransaction.Delete(stalePrefixedKey(staleNode))
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteBlobChunk(blobDigest protocol.BlobDigest, chunkIndex uint64, chunk []byte) error {
	return s.rawTransaction.Write(blobChunkKey(blobDigest, chunkIndex), chunk)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) DeleteBlobChunk(blobDigest protocol.BlobDigest, chunkIndex uint64) error {
	return s.rawTransaction.Delete(blobChunkKey(blobDigest, chunkIndex))
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadBlobMeta(blobDigest protocol.BlobDigest) (*protocol.BlobMeta, error) {
	metaBytes, err := s.rawTransaction.Read(blobMetaPrefixKey(blobDigest))
	if err != nil {
		return nil, fmt.Errorf("error reading blob meta for %s: %w", blobDigest, err)
	}
	if metaBytes == nil {
		// no record of the blob at all
		return nil, nil
	}

	blobMeta, err := serialization.DeserializeBlobMeta(metaBytes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling blob meta for %s: %w", blobDigest, err)
	}
	return &blobMeta, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteBlobMeta(blobDigest protocol.BlobDigest, blobMeta protocol.BlobMeta) error {
	metaBytes, err := serialization.SerializeBlobMeta(blobMeta)
	if err != nil {
		return fmt.Errorf("error marshaling blob meta for %s: %w", blobDigest, err)
	}
	return s.rawTransaction.Write(blobMetaPrefixKey(blobDigest), metaBytes)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) DeleteBlobMeta(blobDigest protocol.BlobDigest) error {
	return s.rawTransaction.Delete(blobMetaPrefixKey(blobDigest))
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadBlobQuotaStats(blobQuotaStatsType protocol.BlobQuotaStatsType, submitter commontypes.OracleID) (protocol.BlobQuotaStats, error) {
	statsBytes, err := s.rawTransaction.Read(blobQuotaStatsPrefixKey(blobQuotaStatsType, submitter))
	if err != nil {
		return protocol.BlobQuotaStats{}, fmt.Errorf("error reading blob quota stats for %s/%d: %w", blobQuotaStatsType, submitter, err)
	}
	if statsBytes == nil {
		return protocol.BlobQuotaStats{}, nil
	}
	blobQuotaStats, err := serialization.DeserializeBlobQuotaStats(statsBytes)
	if err != nil {
		return protocol.BlobQuotaStats{}, fmt.Errorf("error unmarshaling blob quota stats for %s/%d: %w", blobQuotaStatsType, submitter, err)
	}
	return blobQuotaStats, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteBlobQuotaStats(blobQuotaStatsType protocol.BlobQuotaStatsType, submitter commontypes.OracleID, blobQuotaStats protocol.BlobQuotaStats) error {
	statsBytes, err := serialization.SerializeBlobQuotaStats(blobQuotaStats)
	if err != nil {
		return fmt.Errorf("error marshaling blob quota stats for %s/%d: %w", blobQuotaStatsType, submitter, err)
	}
	return s.rawTransaction.Write(blobQuotaStatsPrefixKey(blobQuotaStatsType, submitter), statsBytes)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadStaleBlobIndex(maxStaleSinceSeqNr uint64, limit int) ([]protocol.StaleBlob, error) {
	if maxStaleSinceSeqNr+1 < maxStaleSinceSeqNr {
		return nil, fmt.Errorf("maxStaleSinceSeqNr overflow")
	}

	it := s.rawTransaction.Range(staleBlobIndexPrefixKey(protocol.StaleBlob{0, blobtypes.BlobDigest{}}), staleBlobIndexPrefixKey(protocol.StaleBlob{maxStaleSinceSeqNr + 1, blobtypes.BlobDigest{}}))
	defer it.Close()

	var staleBlobs []protocol.StaleBlob

	for i := 0; i < limit && it.Next(); i++ {
		key := it.Key()
		staleBlob, err := deserializeStaleBlobIndexKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize stale blob index key: %w", err)
		}
		staleBlobs = append(staleBlobs, staleBlob)
	}

	if err := it.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over stale blob index: %w", err)
	}

	return staleBlobs, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteStaleBlobIndex(staleBlob protocol.StaleBlob) error {
	return s.rawTransaction.Write(staleBlobIndexPrefixKey(staleBlob), []byte{})
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) DeleteStaleBlobIndex(staleBlob protocol.StaleBlob) error {
	return s.rawTransaction.Delete(staleBlobIndexPrefixKey(staleBlob))
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadReportsPlusPrecursor(seqNr uint64, reportsPlusPrecursorDigest protocol.ReportsPlusPrecursorDigest) (*ocr3_1types.ReportsPlusPrecursor, error) {
	reportsPlusPrecursor, err := s.rawTransaction.Read(reportsPlusPrecursorKey(seqNr, reportsPlusPrecursorDigest))
	if err != nil {
		return nil, fmt.Errorf("error reading reports plus precursor: %w", err)
	}
	if reportsPlusPrecursor == nil {
		return nil, nil
	}
	return util.PointerTo(ocr3_1types.ReportsPlusPrecursor(reportsPlusPrecursor)), nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteReportsPlusPrecursor(seqNr uint64, reportsPlusPrecursorDigest protocol.ReportsPlusPrecursorDigest, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) error {
	return s.rawTransaction.Write(reportsPlusPrecursorKey(seqNr, reportsPlusPrecursorDigest), reportsPlusPrecursor)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) DeleteReportsPlusPrecursors(minSeqNrToKeep uint64, maxItems int) (done bool, err error) {
	keys, more, err := s.partialExclusiveRangeKeys(reportsPlusPrecursorKey(0, protocol.ReportsPlusPrecursorDigest{}), reportsPlusPrecursorKey(minSeqNrToKeep, protocol.ReportsPlusPrecursorDigest{}), maxItems)
	if err != nil {
		return false, fmt.Errorf("failed to range: %w", err)
	}
	for _, key := range keys {
		if err := s.rawTransaction.Delete(key); err != nil {
			return false, fmt.Errorf("failed to delete key %s: %w", key, err)
		}
	}
	return !more, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ReadUnattestedStateTransitionBlock(seqNr uint64, stateTransitionInputsDigest protocol.StateTransitionInputsDigest) (*protocol.StateTransitionBlock, error) {
	unattestedStateTransitionBlockRaw, err := s.rawTransaction.Read(unattestedBlockKey(seqNr, stateTransitionInputsDigest))
	if err != nil {
		return nil, fmt.Errorf("error reading unattested state transition block: %w", err)
	}
	if unattestedStateTransitionBlockRaw == nil {
		return nil, nil
	}
	ret, err := serialization.DeserializeStateTransitionBlock(unattestedStateTransitionBlockRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize unattested state transition block: %w", err)
	}
	return &ret, nil
}

func (s *SemanticOCR3_1KeyValueDatabaseReadTransaction) ExistsUnattestedStateTransitionBlock(seqNr uint64, stateTransitionInputsDigest protocol.StateTransitionInputsDigest) (bool, error) {
	keys, _, err := s.partialInclusiveRangeKeys(
		unattestedBlockKey(seqNr, stateTransitionInputsDigest),
		unattestedBlockKey(seqNr, stateTransitionInputsDigest),
		1,
	)
	return len(keys) == 1, err
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) WriteUnattestedStateTransitionBlock(seqNr uint64, stateTransitionInputsDigest protocol.StateTransitionInputsDigest, stb protocol.StateTransitionBlock) error {
	unattestedStateTransitionBlockRaw, err := serialization.SerializeStateTransitionBlock(stb)
	if err != nil {
		return fmt.Errorf("failed to serialize unattested state transition block: %w", err)
	}
	return s.rawTransaction.Write(unattestedBlockKey(seqNr, stateTransitionInputsDigest), unattestedStateTransitionBlockRaw)
}

func (s *SemanticOCR3_1KeyValueDatabaseReadWriteTransaction) DeleteUnattestedStateTransitionBlocks(maxSeqNrToDelete uint64, maxItems int) (done bool, err error) {
	if maxSeqNrToDelete+1 < maxSeqNrToDelete {
		return false, fmt.Errorf("maxSeqNrToDelete overflow")
	}

	keys, more, err := s.partialExclusiveRangeKeys(
		unattestedBlockKey(0, protocol.StateTransitionInputsDigest{}),
		unattestedBlockKey(maxSeqNrToDelete+1, protocol.StateTransitionInputsDigest{}),
		maxItems,
	)
	for _, key := range keys {
		if err := s.rawTransaction.Delete(key); err != nil {
			return false, fmt.Errorf("failed to delete key %s: %w", key, err)
		}
	}
	return !more, err
}

var destructiveDestroyForTreeSyncPrefixesToDelete = [][]byte{
	[]byte(pluginPrefix),
	[]byte(treeNodePrefix),
	[]byte(treeRootPrefix),
	[]byte(treeStaleNodePrefix),
	[]byte(lowestPersistedSeqNrKey),
}

const (
	schemaVersionKey       = "OCR3_1_SCHEMA_VERSION"
	supportedSchemaVersion = "1"

	blockPrefix                = "B|"
	pluginPrefix               = "P|"
	blobChunkPrefix            = "BC|"
	blobMetaPrefix             = "BM|"
	blobQuotaStatsPrefix       = "BQS|"
	staleBlobIndexPrefix       = "BI|"
	treeNodePrefix             = "TN|"
	treeRootPrefix             = "TR|"
	treeStaleNodePrefix        = "TSN|"
	reportsPlusPrecursorPrefix = "RP|"
	unattestedBlockPrefix      = "UB|"

	treeSyncStatusKey        = "TSS"
	highestCommittedSeqNrKey = "HCS"
	lowestPersistedSeqNrKey  = "LPS"
)

func initializeSchema(keyValueDatabase ocr3_1types.KeyValueDatabase) error {
	rawTx, err := keyValueDatabase.NewReadWriteTransaction()
	if err != nil {
		return fmt.Errorf("failed to create read write transaction: %w", err)
	}
	defer rawTx.Discard()

	schemaVersion, err := readSchemaVersion(rawTx)
	if err != nil {
		return fmt.Errorf("failed to read schema version: %w", err)
	}
	if schemaVersion != nil {
		if *schemaVersion != supportedSchemaVersion {
			return fmt.Errorf("unsupported schema version: %q, we support: %q", *schemaVersion, supportedSchemaVersion)
		}
		// already initialized
		return nil
	}

	// ensure database is empty
	err = func() error {
		it := rawTx.Range(nil, nil)
		defer it.Close()
		for it.Next() {
			return fmt.Errorf("database is not empty")
		}
		if err := it.Err(); err != nil {
			return fmt.Errorf("failed to ensure database is not empty, iteration error: %w", err)
		}
		return nil
	}()

	if err != nil {
		return err
	}

	if err := writeSchemaVersion(rawTx); err != nil {
		return fmt.Errorf("failed to write schema version: %w", err)
	}

	if err := rawTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func writeSchemaVersion(rawTx ocr3_1types.KeyValueDatabaseReadWriteTransaction) error {
	return rawTx.Write([]byte(schemaVersionKey), []byte(supportedSchemaVersion))
}

func readSchemaVersion(rawTx ocr3_1types.KeyValueDatabaseReadTransaction) (*string, error) {
	schemaVersion, err := rawTx.Read([]byte(schemaVersionKey))
	if err != nil {
		return nil, fmt.Errorf("error reading schema version: %w", err)
	}
	if schemaVersion == nil {
		return nil, nil
	}
	return util.PointerTo(string(schemaVersion)), nil
}

func encodeBigEndianUint64(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

func hashPluginKey(key []byte) jmt.Digest {
	return jmt.DigestKey(key)
}

func pluginPrefixedUnhashedKey(unhashedKey []byte) []byte {
	return append([]byte(pluginPrefix), unhashedKey[:]...)
}

// ────────────────────────── blocks ───────────────────────────

func blockKey(seqNr uint64) []byte {
	return append([]byte(blockPrefix), encodeBigEndianUint64(seqNr)...)
}

func deserializeBlockKey(enc []byte) (uint64, error) {
	if len(enc) < len(blockPrefix) {
		return 0, fmt.Errorf("encoding too short")
	}
	enc = enc[len(blockPrefix):]
	if len(enc) < 8 {
		return 0, fmt.Errorf("encoding too short to contain seqnr")
	}
	return binary.BigEndian.Uint64(enc), nil
}

// ────────────────────────── blobs ───────────────────────────

func blobChunkPrefixedKey(blobDigest protocol.BlobDigest) []byte {
	return append([]byte(blobChunkPrefix), blobDigest[:]...)
}

func blobChunkKey(blobDigest protocol.BlobDigest, chunkIndex uint64) []byte {
	return append(blobChunkPrefixedKey(blobDigest), encodeBigEndianUint64(chunkIndex)...)
}

func blobMetaPrefixKey(blobDigest protocol.BlobDigest) []byte {
	return append([]byte(blobMetaPrefix), blobDigest[:]...)
}

func blobQuotaStatsPrefixKey(blobQuotaStatsType protocol.BlobQuotaStatsType, submitter commontypes.OracleID) []byte {
	base := append([]byte(blobQuotaStatsPrefix), []byte(blobQuotaStatsType)...)
	return append(base, []byte(fmt.Sprintf("%d", submitter))...)
}

// ───────────────────────── meta ────────────────────────────

func rootKey(version uint64) []byte {
	return append([]byte(treeRootPrefix), encodeBigEndianUint64(version)...)
}

// ────────────────────────── tree ───────────────────────────

func treePrefixedKey(nodeKey jmt.NodeKey) []byte {
	base := []byte(treeNodePrefix)
	return serialization.AppendSerializeNodeKey(base, nodeKey)
}

// ────────────────────────── stale tree nodes ───────────────────────────

func staleKeyWithStaleSinceVersionBase(staleSinceVersion jmt.Version) []byte {
	return append([]byte(treeStaleNodePrefix), encodeBigEndianUint64(staleSinceVersion)...)
}

func stalePrefixedKey(staleNode jmt.StaleNode) []byte {
	base := staleKeyWithStaleSinceVersionBase(staleNode.StaleSinceVersion)
	return serialization.AppendSerializeNodeKey(base, staleNode.NodeKey)
}

func deserializeStalePrefixedKey(enc []byte) (jmt.StaleNode, error) {
	base := []byte(treeStaleNodePrefix)
	if len(enc) < len(base) {
		return jmt.StaleNode{}, fmt.Errorf("encoding too short")
	}
	enc = enc[len(base):]
	return serialization.DeserializeStaleNode(enc)
}

// ────────────────────────── stale blobs ───────────────────────────

func staleBlobIndexPrefixKey(staleBlob protocol.StaleBlob) []byte {
	base := []byte(staleBlobIndexPrefix)
	base = binary.BigEndian.AppendUint64(base, staleBlob.StaleSinceSeqNr)
	base = append(base, staleBlob.BlobDigest[:]...)
	return base
}

func deserializeStaleBlobIndexKey(enc []byte) (protocol.StaleBlob, error) {
	base := []byte(staleBlobIndexPrefix)
	if len(enc) < len(base) {
		return protocol.StaleBlob{}, fmt.Errorf("encoding too short")
	}
	enc = enc[len(base):]
	if len(enc) < 8 {
		return protocol.StaleBlob{}, fmt.Errorf("encoding too short to contain seqnr")
	}
	staleSinceSeqNr := binary.BigEndian.Uint64(enc[:8])
	enc = enc[8:]
	if len(enc) < len(protocol.BlobDigest{}) {
		return protocol.StaleBlob{}, fmt.Errorf("encoding too short to contain blob digest")
	}
	blobDigest := protocol.BlobDigest(enc[:len(protocol.BlobDigest{})])
	enc = enc[len(protocol.BlobDigest{}):]
	if len(enc) != 0 {
		return protocol.StaleBlob{}, fmt.Errorf("encoding too long")
	}
	return protocol.StaleBlob{staleSinceSeqNr, blobDigest}, nil
}

// ────────────────────────── reports plus precursor ───────────────────────────

func reportsPlusPrecursorKey(seqNr uint64, reportsPlusPrecursorDigest protocol.ReportsPlusPrecursorDigest) []byte {
	base := []byte(reportsPlusPrecursorPrefix)
	base = append(base, encodeBigEndianUint64(seqNr)...)
	base = append(base, reportsPlusPrecursorDigest[:]...)
	return base
}

// ────────────────────────── unattested blocks ───────────────────────────

func unattestedBlockKey(seqNr uint64, stateTransitionInputsDigest protocol.StateTransitionInputsDigest) []byte {
	base := []byte(unattestedBlockPrefix)
	base = append(base, encodeBigEndianUint64(seqNr)...)
	base = append(base, stateTransitionInputsDigest[:]...)
	return base
}

type limitCheckWriteSet struct {
	m                         map[string][]byte
	keys                      int
	keysLimit                 int
	keysPlusValuesLength      int
	keysPlusValuesLengthLimit int
}

func newLimitCheckWriteSet(limits ocr3_1types.ReportingPluginLimits) *limitCheckWriteSet {
	return &limitCheckWriteSet{
		make(map[string][]byte),
		0,
		limits.MaxKeyValueModifiedKeys,
		0,
		limits.MaxKeyValueModifiedKeysPlusValuesBytes,
	}
}

func (l *limitCheckWriteSet) modify(key []byte, value []byte) error {

	add, sub := 0, 0
	if prevValue, ok := l.m[string(key)]; ok {
		add = len(value)
		sub = len(prevValue)

		// key already in write set, no change to l.keys
	} else {
		if len(key)+len(value) < len(key) {
			return fmt.Errorf("key + value length overflow")
		}
		add = len(key) + len(value)

		// new key being added
		if l.keys+1 < l.keys {
			return fmt.Errorf("keys overflow")
		}
		if l.keys+1 > l.keysLimit {
			return fmt.Errorf("keys %d exceed limit %d", l.keys+1, l.keysLimit)
		}
		l.keys++
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

func (l *limitCheckWriteSet) Pairs() []protocol.KeyValuePairWithDeletions {
	pairs := make([]protocol.KeyValuePairWithDeletions, 0, len(l.m))
	for k, v := range l.m {
		pairs = append(pairs, protocol.KeyValuePairWithDeletions{
			[]byte(k),
			v,
			v == nil,
		})
	}
	return pairs
}
