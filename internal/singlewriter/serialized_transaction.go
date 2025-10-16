package singlewriter

import (
	"fmt"
	"sync"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
)

type SerializedTransaction struct {
	*overlayTransaction
	timestamp                                  uint64
	maxCommittedTransactionTimestampAtCreation uint64

	tracker *ConflictTracker
}

var _ ocr3_1types.KeyValueDatabaseReadWriteTransaction = &SerializedTransaction{}

func NewSerializedTransaction(
	keyValueDatabase ocr3_1types.KeyValueDatabase,
	conflictTracker *ConflictTracker,
) (ocr3_1types.KeyValueDatabaseReadWriteTransaction, error) {
	timestamp, maxCommittedTransactionTimestamp := conflictTracker.beginTransaction()
	rawReaderTx, err := keyValueDatabase.NewReadTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to create read transaction to create a SerializedTransaction: %w", err)
	}
	return &SerializedTransaction{
		&overlayTransaction{
			make(map[string]operation),
			keyValueDatabase,
			rawReaderTx,
			sync.Mutex{},
			sync.WaitGroup{},
			txStatusOpen,
		},
		timestamp,
		maxCommittedTransactionTimestamp,
		conflictTracker,
	}, nil
}

// Commit the overlay to the underlying DB as long as no other SerializedTransaction is committed since this
// SerializedTransaction was created.
func (st *SerializedTransaction) Commit() error {
	st.mu.Lock()
	defer st.mu.Unlock()

	if st.status != txStatusOpen {
		return ErrClosed
	}
	overlay := st.overlay
	st.lockedDiscard()

	if err := st.tracker.lockAndPrepareToCommit(st.maxCommittedTransactionTimestampAtCreation); err != nil {
		return err
	}
	committed := false
	defer func() {
		st.tracker.finalizeCommitAndUnlock(committed, st.timestamp)
	}()

	rawReadWriteTx, err := st.keyValueDatabase.NewReadWriteTransaction()
	if err != nil {
		return fmt.Errorf("failed to create read write transaction to commit SerializedTransaction: %w", err)
	}

	defer rawReadWriteTx.Discard()

	err = lockedCommit(overlay, rawReadWriteTx)
	if err != nil {
		return err
	}

	st.status = txStatusCommitted
	committed = true
	return nil
}
