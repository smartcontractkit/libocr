package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
)

const (
	blobPopulateStaleIndexShortInterval = 10 * time.Second
	blobPopulateStaleIndexLongInterval  = 60 * time.Second

	maxBlobDigestExpirySeqNrsToReadInSingleTransaction = 100_000
	maxBlobsToPopulateStaleIndexForInSingleTransaction = 1_000
)

type populateStaleBlobIndexStats struct {
	numScannedBlobs          int
	numPopulatedIndexEntries int
}

func populateStaleBlobIndex(
	ctx context.Context,
	logger commontypes.Logger,
	kvDb KeyValueDatabase,
	minBlobDigest blobtypes.BlobDigest,
) (nextMinBlobDigest BlobDigest, _ populateStaleBlobIndexStats, done bool, err error) {
	chDone := ctx.Done()

	tx, err := kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return BlobDigest{}, populateStaleBlobIndexStats{}, false, fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	items, readMore, err := tx.ReadBlobDigestExpirySeqNrs(minBlobDigest, maxBlobDigestExpirySeqNrsToReadInSingleTransaction)
	if err != nil {
		return BlobDigest{}, populateStaleBlobIndexStats{}, false, fmt.Errorf("failed to read blob digest expiry seq nrs: %w", err)
	}

	var highestProcessed *BlobDigest
	stats := populateStaleBlobIndexStats{
		0,
		0,
	}
	populateMore := false

	for _, item := range items {
		select {
		case <-chDone:
			return BlobDigest{}, populateStaleBlobIndexStats{}, false, ctx.Err()
		default:
		}

		staleBlob := staleBlob(item.ExpirySeqNr, item.BlobDigest)
		exists, err := tx.ExistsStaleBlobIndex(staleBlob)
		if err != nil {
			return BlobDigest{}, populateStaleBlobIndexStats{}, false, fmt.Errorf("failed to check existence of stale blob index for %v: %w", staleBlob, err)
		}

		stats.numScannedBlobs++

		if exists {
			highestProcessed = &item.BlobDigest
			continue
		}

		if stats.numPopulatedIndexEntries >= maxBlobsToPopulateStaleIndexForInSingleTransaction {
			populateMore = true
			break
		}

		logger.Trace("populateStaleBlobIndex: blob is missing from stale blob index, populating index entry", commontypes.LogFields{
			"expirySeqNr": item.ExpirySeqNr,
			"blobDigest":  item.BlobDigest,
		})

		err = tx.WriteStaleBlobIndex(staleBlob)
		if err != nil {
			return BlobDigest{}, populateStaleBlobIndexStats{}, false, fmt.Errorf("failed to write to stale blob index for %v: %w", staleBlob, err)
		}

		highestProcessed = &item.BlobDigest
		stats.numPopulatedIndexEntries++
	}

	if err := tx.Commit(); err != nil {
		return BlobDigest{}, populateStaleBlobIndexStats{}, false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if readMore || populateMore {
		var nextMinBlobDigest BlobDigest
		if highestProcessed != nil {
			nextMinBlobDigest = blobtypes.WrappingIncrementBlobDigest(*highestProcessed)
		} else {
			// If we never even tried to process anything, while knowing
			// that there are more things in the range to process or populate,
			// something's very wrong. Try skipping over the minBlobDigest,
			// though it's not clear that this will help.
			logger.Error("populateStaleBlobIndex: need to read more blob digest expiry seq nrs or populate more stale blob indices, "+
				"but did not read or populate anything in this iteration, so skipping over the minBlobDigest", commontypes.LogFields{
				"minBlobDigest": minBlobDigest,
			})
			nextMinBlobDigest = blobtypes.WrappingIncrementBlobDigest(minBlobDigest)
		}
		return nextMinBlobDigest, stats, false, nil
	}
	return blobtypes.MinBlobDigest(), stats, true, nil
}

func RunBlobPopulateStaleIndex(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	kvDb KeyValueDatabase,
) {
	chDone := ctx.Done()
	chTick := time.After(0)

	minBlobDigest := blobtypes.MinBlobDigest()
	interval := blobPopulateStaleIndexShortInterval

	for {
		select {
		case <-chTick:
		case <-chDone:
			return
		}

		nextMinBlobDigest, stats, done, err := populateStaleBlobIndex(ctx, logger, kvDb, minBlobDigest)
		if err != nil {
			logger.Warn("BlobPopulateStaleIndex: failed to populate stale index. Will retry soon.", commontypes.LogFields{
				"minBlobDigest":   minBlobDigest,
				"error":           err,
				"waitBeforeRetry": interval.String(),
			})
		} else {
			if done {
				interval = blobPopulateStaleIndexLongInterval
			}
			var logFn func(msg string, fields commontypes.LogFields)
			if stats.numPopulatedIndexEntries > 0 {
				logFn = logger.Info
			} else {
				logFn = logger.Debug
			}
			logFn("BlobPopulateStaleIndex: populated stale index", commontypes.LogFields{
				"minBlobDigest":           minBlobDigest,
				"nextMinBlobDigest":       nextMinBlobDigest,
				"done":                    done,
				"stats":                   fmt.Sprintf("%+v", stats),
				"waitBeforeNextIteration": interval.String(),
			})
			minBlobDigest = nextMinBlobDigest
		}

		chTick = time.After(interval)
	}
}
