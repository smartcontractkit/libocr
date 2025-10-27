package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
)

const (
	unattestedStateTransitionBlockFetchingReapInterval = 5 * time.Second

	maxUnattestedStateTransitionBlocksToReapInOneGo = 100_000
)

func reapUnattestedBlocks(ctx context.Context, kvDb KeyValueDatabase, logger commontypes.Logger) error {
	committedKVSeqNr, err := committedKVSeqNr(kvDb)
	if err != nil {
		return fmt.Errorf("failed to read highest committed kv seq nr: %w", err)
	}

	logger.Info("RunOutcomeGenerationReap: reaping unattested state transition blocks", commontypes.LogFields{
		"committedKVSeqNr": committedKVSeqNr,
	})

	for {
		done, err := reapSomeUnattestedBlocks(kvDb, committedKVSeqNr)
		if err != nil {
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if done {
			break
		}
	}
	return nil
}

func reapSomeUnattestedBlocks(kvDb KeyValueDatabase, committedKVSeqNr uint64) (done bool, err error) {
	kvTxn, err := kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return false, fmt.Errorf("failed to open unserialized transaction: %w", err)
	}
	defer kvTxn.Discard()

	done, err = kvTxn.DeleteUnattestedStateTransitionBlocks(committedKVSeqNr, maxUnattestedStateTransitionBlocksToReapInOneGo)
	if err != nil {
		return false, fmt.Errorf("failed to delete unattested state transition blocks: %w", err)
	}
	err = kvTxn.Commit()
	if err != nil {
		return false, fmt.Errorf("failed to commit: %w", err)
	}
	return done, err
}

func RunOutcomeGenerationReap(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	kvDb KeyValueDatabase,
) {
	chDone := ctx.Done()
	ticker := time.NewTicker(unattestedStateTransitionBlockFetchingReapInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		case <-chDone:
			return
		}

		logger.Info("RunOutcomeGenerationReap: calling reapUnattestedBlocks", nil)
		err := reapUnattestedBlocks(ctx, kvDb, logger)
		if err != nil {
			logger.Warn("RunOutcomeGenerationReap: failed to reap unattested state transition blocks. Will retry soon.", commontypes.LogFields{
				"error":           err,
				"waitBeforeRetry": unattestedStateTransitionBlockFetchingReapInterval.String(),
			})
		}
	}
}
