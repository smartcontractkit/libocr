package protocol

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
)

const (
	stateBlockReplayInterval          = 10 * time.Second
	stateBlockReplayFastFollowOnError = stateBlockReplayInterval / 10

	maxBlocksToReplayInOneGo = 100
)

type errReplayVerifiedBlock struct {
	SeqNr uint64
	Cause error
}

func (e *errReplayVerifiedBlock) Error() string {
	return fmt.Sprintf("failed to replay verified block %d: %v", e.SeqNr, e.Cause)
}

func (e *errReplayVerifiedBlock) Unwrap() error {
	return e.Cause
}

func tryReplay(ctx context.Context, kvDb KeyValueDatabase, logger loghelper.LoggerWithContext, metrics *stateSyncMetrics) error {
	kvReadTxn, err := kvDb.NewReadTransactionUnchecked()
	if err != nil {
		return fmt.Errorf("failed to create read transaction")
	}
	defer kvReadTxn.Discard()

	committedSeqNr, err := kvReadTxn.ReadHighestCommittedSeqNr()
	if err != nil {
		return fmt.Errorf("failed to read highest committed seq nr: %w", err)
	}

	for {
		astbsToReplay, more, err := getReplayableBlocks(kvReadTxn, committedSeqNr)
		if err != nil {
			return fmt.Errorf("failed to get blocks to replay: %w", err)
		}

		for _, astb := range astbsToReplay {
			block := astb.StateTransitionBlock
			seqNr := block.SeqNr()

			logger.Trace("StateBlockReplay: trying to replay block", commontypes.LogFields{
				"seqNr": seqNr,
			})

			err := func() error {
				tx, err := kvDb.NewSerializedReadWriteTransaction(seqNr)
				if err != nil {
					return fmt.Errorf("failed to create kv read/write transaction: %w", err)
				}
				defer tx.Discard()

				// next block found, has been verified before being persisted so we don't check again
				err = replayVerifiedBlock(logger, tx, &block)
				if err != nil {
					return &errReplayVerifiedBlock{seqNr, err}
				}
				err = tx.Commit()
				if err != nil {
					return fmt.Errorf("failed to commit transaction: %w", err)
				}
				return nil
			}()
			if err != nil {
				return fmt.Errorf("failed to replay block %d: %w", seqNr, err)
			}
			metrics.attestedBlocksReplayedTotal.Inc()
			logger.Debug("StateBlockReplay: 🐌✅ committed", commontypes.LogFields{
				"seqNr": seqNr,
			})
			committedSeqNr = seqNr
		}

		if !more {
			break
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}

func getReplayableBlocks(kvReadTxn KeyValueDatabaseReadTransaction, committedSeqNr uint64) ([]AttestedStateTransitionBlock, bool, error) {
	blocks, more, err := kvReadTxn.ReadAttestedStateTransitionBlocks(committedSeqNr+1, maxBlocksToReplayInOneGo)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read attested state transition blocks: %w", err)
	}
	return blocks, more, nil
}

func replayVerifiedBlock(logger loghelper.LoggerWithContext, kvReadWriteTxn KeyValueDatabaseReadWriteTransaction, stb *StateTransitionBlock) error {
	seqNr := stb.SeqNr()
	logger = logger.MakeChild(commontypes.LogFields{
		"replay": "YES",
		"seqNr":  seqNr,
	})

	logger.Trace("replaying state transition block", nil)

	stateRootDigest, err := kvReadWriteTxn.ApplyWriteSet(stb.StateWriteSet.Entries)
	if err != nil {
		return fmt.Errorf("failed to apply write set for seq nr %d: %w", seqNr, err)
	}

	if stateRootDigest != stb.StateRootDigest {
		return fmt.Errorf("state root digest mismatch from block replay for seq nr %d: expected %v, actual %v", seqNr, stb.StateRootDigest, stateRootDigest)
	}

	return nil
}

func RunStateSyncBlockReplay(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	kvDb KeyValueDatabase,
	metrics *stateSyncMetrics,
	chStateSyncToStateSyncBlockReplay <-chan EventStateSyncBlockReplayWake,
	chStateSyncBlockReplayToStateSync chan<- EventStateSyncBlockReplayFailure,
) {
	chDone := ctx.Done()
	chTick := time.After(0)
	notifyStateSyncOfFailedReplay := false
	var failedReplaySeqNr uint64

	for {
		if notifyStateSyncOfFailedReplay {

			select {
			case chStateSyncBlockReplayToStateSync <- EventStateSyncBlockReplayFailure{failedReplaySeqNr}:
				notifyStateSyncOfFailedReplay = false
				logger.Info("StateBlockReplay: notified state sync of failed replay", commontypes.LogFields{
					"failedReplaySeqNr": failedReplaySeqNr,
				})

				chTick = time.After(stateBlockReplayInterval)
			case <-chDone:
				return
			}
			continue
		}

		select {
		case <-chTick:
		case <-chStateSyncToStateSyncBlockReplay:
		case <-chDone:
			return
		}

		logger.Trace("StateBlockReplay: calling tryReplay", nil)
		err := tryReplay(ctx, kvDb, logger, metrics)
		if err != nil {
			logger.Warn("StateBlockReplay: failed while trying to replay blocks", commontypes.LogFields{
				"error": err,
			})
			var errReplayVerifiedBlock *errReplayVerifiedBlock
			if errors.As(err, &errReplayVerifiedBlock) {
				failedReplaySeqNr = errReplayVerifiedBlock.SeqNr
				notifyStateSyncOfFailedReplay = true
				logger.Info("StateBlockReplay: will notify state sync of failed replay", commontypes.LogFields{
					"failedReplaySeqNr": errReplayVerifiedBlock.SeqNr,
				})
			} else {

				chTick = time.After(stateBlockReplayFastFollowOnError)
			}

		} else {
			chTick = time.After(stateBlockReplayInterval)
		}
	}
}
