package protocol

import (
	"context"
	"fmt"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
)

const (
	maxStateKeysToDestroyInSingleTransaction = 1_000_000
)

type destroyStateIfNeededResult int

const (
	_ destroyStateIfNeededResult = iota
	destroyStateIfNeededResultDone
	destroyStateIfNeededResultDoneButNeedMore
	destroyStateIfNeededResultNotNeeded
	destroyStateIfNeededResultError
)

func destroyStateIfNeeded(kvDb KeyValueDatabase, logger commontypes.Logger) (destroyStateIfNeededResult, error) {
	tx, err := kvDb.NewSerializedReadWriteTransactionUnchecked()
	if err != nil {
		return destroyStateIfNeededResultError, fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	treeSyncStatus, err := tx.ReadTreeSyncStatus()
	if err != nil {
		return destroyStateIfNeededResultError, fmt.Errorf("failed to read tree sync status: %w", err)
	}

	logger.Info("StateDestroyIfNeeded: read tree sync status", commontypes.LogFields{
		"treeSyncStatus": treeSyncStatus,
	})

	if treeSyncStatus.Phase != TreeSyncPhaseWaiting {
		return destroyStateIfNeededResultNotNeeded, nil
	}

	done, err := tx.DestructiveDestroyForTreeSync(maxStateKeysToDestroyInSingleTransaction)
	if err != nil {
		return destroyStateIfNeededResultError, fmt.Errorf("failed to delete everything but tree sync status: %w", err)
	}

	if done {
		err := tx.WriteTreeSyncStatus(TreeSyncStatus{
			TreeSyncPhaseActive,
			treeSyncStatus.TargetSeqNr,
			treeSyncStatus.TargetStateRootDigest,
			treeSyncStatus.PendingKeyDigestRanges,
		})
		if err != nil {
			return destroyStateIfNeededResultError, fmt.Errorf("failed to write tree sync status after being done destroying state: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return destroyStateIfNeededResultError, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if done {
		return destroyStateIfNeededResultDone, nil
	} else {
		return destroyStateIfNeededResultDoneButNeedMore, nil
	}
}

func RunStateSyncDestroyIfNeeded(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	kvDb KeyValueDatabase,
	chNotificationFromStateSync <-chan struct{},
) {
	logger = logger.MakeChild(commontypes.LogFields{"proto": "stateSyncDestroyIfNeeded"})

	chDone := ctx.Done()

	for {
		select {
		case <-chNotificationFromStateSync:
		case <-chDone:
			return
		}

		for {
			logger.Trace("RunStateSyncDestroyIfNeeded: destroying state if needed...", nil)
			destroyStateIfNeededResult, err := destroyStateIfNeeded(kvDb, logger)

			followupImmediately := false
			switch destroyStateIfNeededResult {
			case destroyStateIfNeededResultDone:
				logger.Info("RunStateSyncDestroyIfNeeded: destroyed state 💣", nil)
			case destroyStateIfNeededResultDoneButNeedMore:
				logger.Debug("RunStateSyncDestroyIfNeeded: destroyed state, but need more", nil)
				followupImmediately = true
			case destroyStateIfNeededResultNotNeeded:
				logger.Trace("RunStateSyncDestroyIfNeeded: not needed to destroy state", nil)
			case destroyStateIfNeededResultError:
				logger.Warn("RunStateSyncDestroyIfNeeded: failed to destroy state", commontypes.LogFields{
					"error": err,
				})
			}

			if !followupImmediately {
				break
			}

			select {
			case <-chDone:
				return
			default:
			}
		}
	}
}
