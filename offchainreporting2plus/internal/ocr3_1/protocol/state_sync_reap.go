package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
)

const (
	stateReapInterval          = 10 * time.Second
	stateReapFastFollowOnError = 120 * time.Millisecond

	maxBlocksToReapInOneGo    = 100_000
	maxTreeNodesToReapInOneGo = 10_000
	maxTreeRootsToReapInOneGo = 100_000
)

func reapState(ctx context.Context, kvDb KeyValueDatabase, logger commontypes.Logger) (done bool, err error) {

	tx, err := kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return false, fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	treeSyncStatus, err := tx.ReadTreeSyncStatus()
	if err != nil {
		return false, fmt.Errorf("failed to read tree sync status: %w", err)
	}
	if treeSyncStatus.Phase != TreeSyncPhaseInactive {
		return false, fmt.Errorf("tree sync is not inactive")
	}
	highestCommittedSeqNr, err := tx.ReadHighestCommittedSeqNr()
	if err != nil {
		return false, fmt.Errorf("failed to read highest committed seq nr: %w", err)
	}

	lowestPersistedSeqNr, err := tx.ReadLowestPersistedSeqNr()
	if err != nil {
		return false, fmt.Errorf("failed to read lowest persisted seq nr: %w", err)
	}

	desiredLowestPersistedSeqNr := desiredLowestPersistedSeqNr(highestCommittedSeqNr)
	if desiredLowestPersistedSeqNr > lowestPersistedSeqNr {
		logger.Info("RunStateSyncReap: new lowest persisted seq nr", commontypes.LogFields{
			"desiredLowestPersistedSeqNr": desiredLowestPersistedSeqNr,
			"lowestPersistedSeqNr":        lowestPersistedSeqNr,
		})

		// write new lowest persisted seq nr first
		if err := tx.WriteLowestPersistedSeqNr(desiredLowestPersistedSeqNr); err != nil {
			return false, fmt.Errorf("failed to write lowest persisted seq nr: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("failed to commit transaction: %w", err)
		}
	} else {
		tx.Discard()
	}

	// Reap unneeded blocks

	logger.Info("RunStateSyncReap: reaping blocks", commontypes.LogFields{
		"desiredLowestPersistedSeqNr": desiredLowestPersistedSeqNr,
		"lowestPersistedSeqNr":        lowestPersistedSeqNr,
	})

	for {
		done, err := reapBlocks(kvDb, desiredLowestPersistedSeqNr)
		if err != nil {
			return false, fmt.Errorf("failed to reap blocks: %w", err)
		}
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
		if done {
			break
		}
	}

	// Reap unneeded tree nodes

	logger.Info("RunStateSyncReap: reaping stale nodes from tree", commontypes.LogFields{
		"desiredLowestPersistedSeqNr": desiredLowestPersistedSeqNr,
	})

	for {
		done, err := reapTreeNodes(kvDb, desiredLowestPersistedSeqNr)
		if err != nil {
			return false, fmt.Errorf("failed to reap tree nodes: %w", err)
		}
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
		if done {
			break
		}
	}

	logger.Info("RunStateSyncReap: reaping stale roots", commontypes.LogFields{
		"desiredLowestPersistedSeqNr": desiredLowestPersistedSeqNr,
	})

	for {
		done, err := reapTreeRoots(kvDb, desiredLowestPersistedSeqNr)
		if err != nil {
			return false, fmt.Errorf("failed to reap tree roots: %w", err)
		}
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
		if done {
			break
		}
	}

	return true, nil
}

func reapBlocks(kvDb KeyValueDatabase, desiredLowestPersistedSeqNr uint64) (done bool, err error) {
	tx, err := kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return false, fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	if desiredLowestPersistedSeqNr == 0 {
		return true, nil
	}

	done, err = tx.DeleteAttestedStateTransitionBlocks(desiredLowestPersistedSeqNr-1, maxBlocksToReapInOneGo)
	if err != nil {
		return false, fmt.Errorf("failed to delete stale blocks: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return done, nil
}

func reapTreeNodes(kvDb KeyValueDatabase, desiredLowestPersistedSeqNr uint64) (done bool, err error) {
	tx, err := kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return false, fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	done, err = tx.DeleteStaleNodes(RootVersion(desiredLowestPersistedSeqNr), maxTreeNodesToReapInOneGo)
	if err != nil {
		return false, fmt.Errorf("failed to delete stale nodes: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return done, nil
}

func reapTreeRoots(kvDb KeyValueDatabase, desiredLowestPersistedSeqNr uint64) (done bool, err error) {
	tx, err := kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return false, fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	done, err = tx.DeleteRoots(RootVersion(desiredLowestPersistedSeqNr), maxTreeRootsToReapInOneGo)
	if err != nil {
		return false, fmt.Errorf("failed to delete roots: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return done, nil
}

func RunStateSyncReap(
	ctx context.Context,
	config ocr3config.SharedConfig,
	logger loghelper.LoggerWithContext,
	database Database,
	kvDb KeyValueDatabase,
) {
	chDone := ctx.Done()
	chTick := time.After(0)

	for {
		select {
		case <-chTick:
		case <-chDone:
			return
		}

		logger.Info("RunStateSyncReap: calling reapState", nil)
		done, err := reapState(ctx, kvDb, logger)
		if err != nil {
			logger.Warn("RunStateSyncReap: failed to reap state. Will retry soon.", commontypes.LogFields{
				"error":           err,
				"waitBeforeRetry": stateReapFastFollowOnError.String(),
			})
			chTick = time.After(stateReapFastFollowOnError)
		} else if !done {
			chTick = time.After(0)
		} else {
			chTick = time.After(stateReapInterval)
		}
	}
}
