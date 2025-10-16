package singlewriter

import (
	"fmt"
	"sync"
)

type ConflictTracker struct {
	mu                      sync.Mutex
	maxCreatedTxTimestamp   uint64
	maxCommittedTxTimestamp uint64
}

func NewConflictTracker() *ConflictTracker {
	return &ConflictTracker{
		sync.Mutex{},
		uint64(0),
		uint64(0),
	}
}

func (ct *ConflictTracker) beginTransaction() (uint64, uint64) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.maxCreatedTxTimestamp++
	return ct.maxCreatedTxTimestamp, ct.maxCommittedTxTimestamp
}

func (ct *ConflictTracker) lockAndPrepareToCommit(maxCommittedTxTimestampAtCreation uint64) error {
	ct.mu.Lock()
	if maxCommittedTxTimestampAtCreation != ct.maxCommittedTxTimestamp {
		ct.mu.Unlock()
		return fmt.Errorf("concurrent conflict detected: expected maxCommittedTxTimestamp: %d, got: %d", maxCommittedTxTimestampAtCreation, ct.maxCommittedTxTimestamp)
	}
	return nil
}

func (ct *ConflictTracker) finalizeCommitAndUnlock(success bool, timestamp uint64) {
	if success {
		ct.maxCommittedTxTimestamp = timestamp
	}
	ct.mu.Unlock()
}
