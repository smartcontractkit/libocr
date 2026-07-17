package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
)

const (
	blobReapMinInterval               = 3 * trackHighestCommittedSeqNrInterval
	maxBlobsToReapInSingleTransaction = 1_000
)

type blobReapStats struct {
	numReaped int
}

func reapBlobs(ctx context.Context, logger commontypes.Logger, kvDb KeyValueDatabase, maxStaleSinceSeqNr uint64, perOracleMetrics []*blobOracleMetrics) (done bool, _ blobReapStats, err error) {
	chDone := ctx.Done()

	tx, err := kvDb.NewUnserializedReadWriteTransactionUnchecked()
	if err != nil {
		return false, blobReapStats{}, fmt.Errorf("failed to create read/write transaction: %w", err)
	}
	defer tx.Discard()

	staleBlobs, err := tx.ReadStaleBlobIndex(maxStaleSinceSeqNr, maxBlobsToReapInSingleTransaction+1)
	if err != nil {
		return false, blobReapStats{}, fmt.Errorf("failed to read stale blob index: %w", err)
	}

	if len(staleBlobs) == 0 {

		return true, blobReapStats{}, nil
	}

	stats := blobReapStats{}
	updatedReapedStats := make(map[commontypes.OracleID]BlobQuotaStats)
	for i, staleBlob := range staleBlobs {
		if i >= maxBlobsToReapInSingleTransaction {
			break
		}

		select {
		case <-chDone:
			return true, blobReapStats{}, ctx.Err()
		default:
		}

		if err := reapSingleBlob(tx, logger, staleBlob, updatedReapedStats); err != nil {
			return false, blobReapStats{}, fmt.Errorf("failed to reap single blob: %w", err)
		}
		stats.numReaped++
	}

	if err := tx.Commit(); err != nil {
		return false, blobReapStats{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	for submitter, stats := range updatedReapedStats {
		if int(submitter) < len(perOracleMetrics) {
			perOracleMetrics[submitter].SetTheirReapedStats(stats)
		}
	}

	return len(staleBlobs) <= maxBlobsToReapInSingleTransaction, stats, nil
}

func reapSingleBlob(tx KeyValueDatabaseReadWriteTransaction, logger commontypes.Logger, staleBlob StaleBlob, updatedReapedStats map[commontypes.OracleID]BlobQuotaStats) error {
	if err := tx.DeleteStaleBlobIndex(staleBlob); err != nil {
		return fmt.Errorf("failed to delete stale blob index: %w", err)
	}

	meta, err := tx.ReadBlobMeta(staleBlob.BlobDigest)
	if err != nil {
		return fmt.Errorf("failed to read blob meta: %w", err)
	}
	if meta == nil {
		logger.Warn("reapSingleBlob: orphan stale blob index entry (no blob meta), dropped index entry", commontypes.LogFields{
			"staleSinceSeqNr": staleBlob.StaleSinceSeqNr,
			"blobDigest":      staleBlob.BlobDigest,
		})
		return nil
	}

	for chunkIndex, chunkHave := range meta.ChunkHaves {
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

	// increase reaped quota stats

	existingQuotaStats, err := tx.ReadBlobQuotaStats(BlobQuotaStatsTypeReaped, meta.Submitter)
	if err != nil {
		return fmt.Errorf("failed to read blob quota stats: %w", err)
	}
	thisBlob := BlobQuotaStats{1, meta.PayloadLength}
	updatedQuotaStats, ok := existingQuotaStats.Add(thisBlob)
	if !ok {
		return fmt.Errorf("quotaStats overflow")
	}
	err = tx.WriteBlobQuotaStats(BlobQuotaStatsTypeReaped, meta.Submitter, updatedQuotaStats)
	if err != nil {
		return fmt.Errorf("failed to write blob quota stats: %w", err)
	}

	updatedReapedStats[meta.Submitter] = updatedQuotaStats

	return nil
}

func RunBlobReap(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	kvDb KeyValueDatabase,
	chBlobExchangeToBlobReap <-chan uint64,
	perOracleMetrics []*blobOracleMetrics,
) {
	chDone := ctx.Done()
	var chTick <-chan time.Time
	var maxStaleSinceSeqNr uint64
	haveMaxStaleSinceSeqNr := false

	for {
		select {
		case highestCommittedSeqNr := <-chBlobExchangeToBlobReap:
			if !haveMaxStaleSinceSeqNr || maxStaleSinceSeqNr < highestCommittedSeqNr {
				logger.Debug("BlobReap: received higher highest committed seq nr from blob exchange", commontypes.LogFields{
					"highestCommittedSeqNr": highestCommittedSeqNr,
				})

				maxStaleSinceSeqNr = highestCommittedSeqNr
				haveMaxStaleSinceSeqNr = true
				chTick = time.After(0)
			}
		case <-chTick:
			if !haveMaxStaleSinceSeqNr {
				chTick = nil
				continue
			}

			done, stats, err := reapBlobs(ctx, logger, kvDb, maxStaleSinceSeqNr, perOracleMetrics)
			if err != nil {
				logger.Warn("BlobReap: failed to reap blobs", commontypes.LogFields{
					"maxStaleSinceSeqNr": maxStaleSinceSeqNr,
					"error":              err,
				})
			} else {
				logger.Info("BlobReap: finished reaping blobs", commontypes.LogFields{
					"maxStaleSinceSeqNr": maxStaleSinceSeqNr,
					"done":               done,
					"stats":              fmt.Sprintf("%+v", stats),
				})
			}
			if done {
				chTick = time.After(blobReapMinInterval)
			} else {
				chTick = time.After(0)
			}
		case <-chDone:
			return
		}
	}
}
