package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
)

const (
	stateBlockReplayInterval          = 10 * time.Second
	stateBlockReplayFastFollowOnError = stateBlockReplayInterval / 10

	maxBlocksToReplayInOneGo = 100
)

func tryReplay(ctx context.Context, kvDb KeyValueDatabase, logger loghelper.LoggerWithContext) error {
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
					return fmt.Errorf("failed to replay verified block %d: %w", seqNr, err)
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

	stateRootDigest, err := kvReadWriteTxn.ApplyWriteSet(stb.StateTransitionOutputs.WriteSet)
	if err != nil {
		return fmt.Errorf("failed to apply write set for seq nr %d: %w", seqNr, err)
	}

	if stateRootDigest != stb.StateRootDigest {
		return fmt.Errorf("state root digest mismatch from block replay for seq nr %d: expected %s, actual %s", seqNr, stb.StateRootDigest, stateRootDigest)
	}

	return nil
}

func RunStateSyncBlockReplay(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	kvDb KeyValueDatabase,
	chNotificationFromStateSync <-chan struct{},
) {
	chDone := ctx.Done()
	chTick := time.After(0)

	for {
		select {
		case <-chTick:
		case <-chNotificationFromStateSync:
		case <-chDone:
			return
		}

		logger.Trace("StateBlockReplay: calling tryReplay", nil)
		err := tryReplay(ctx, kvDb, logger)
		if err != nil {
			logger.Warn("StateBlockReplay: failed while trying to replay blocks", commontypes.LogFields{
				"error": err,
			})
			chTick = time.After(stateBlockReplayFastFollowOnError)
		} else {
			chTick = time.After(stateBlockReplayInterval)
		}
	}
}
