package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
)

const (
	blobReapInterval                  = 10 * time.Second
	maxBlobsToReapInSingleTransaction = 100
)

func reapBlobs(ctx context.Context, kvDb KeyValueDatabase) (done bool, err error) {
	chDone := ctx.Done()

	tx, err := kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return false, fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	committedSeqNr, err := tx.ReadHighestCommittedSeqNr()
	if err != nil {
		return false, fmt.Errorf("failed to read highest committed seq nr: %w", err)
	}

	staleBlobs, err := tx.ReadStaleBlobIndex(committedSeqNr, maxBlobsToReapInSingleTransaction+1)
	if err != nil {
		return false, fmt.Errorf("failed to read stale blob index: %w", err)
	}

	if len(staleBlobs) == 0 {

		return true, nil
	}

	for i, staleBlob := range staleBlobs {
		if i >= maxBlobsToReapInSingleTransaction {
			break
		}

		select {
		case <-chDone:
			return true, ctx.Err()
		default:
		}

		if err := reapSingleBlob(tx, staleBlob); err != nil {
			return false, fmt.Errorf("failed to reap single blob: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return len(staleBlobs) <= maxBlobsToReapInSingleTransaction, nil
}

func reapSingleBlob(tx KeyValueDatabaseReadWriteTransaction, staleBlob StaleBlob) error {
	meta, err := tx.ReadBlobMeta(staleBlob.BlobDigest)
	if err != nil {
		return fmt.Errorf("failed to read blob meta: %w", err)
	}

	if meta == nil {
		return fmt.Errorf("blob meta is nil")
	}

	for chunkIndex, chunkHave := range meta.ChunksHave {
		if !chunkHave {
			continue
		}

		if err := tx.DeleteBlobChunk(staleBlob.BlobDigest, uint64(chunkIndex)); err != nil {
			return fmt.Errorf("failed to delete blob chunk: %w", err)
		}
	}

	if err := tx.DeleteBlobMeta(staleBlob.BlobDigest); err != nil {
		return fmt.Errorf("failed to delete blob meta: %w", err)
	}
	if err := tx.DeleteStaleBlobIndex(staleBlob); err != nil {
		return fmt.Errorf("failed to delete stale blob index: %w", err)
	}

	return nil
}

func RunBlobReap(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
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

		done, err := reapBlobs(ctx, kvDb)
		if err != nil {
			logger.Warn("BlobReap: failed to reap blobs", commontypes.LogFields{
				"error": err,
			})
		}
		if done {
			chTick = time.After(blobReapInterval)
		} else {
			chTick = time.After(0)
		}
	}
}
