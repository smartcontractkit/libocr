package persist

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type TransmissionDBUpdate struct {
	Timestamp           types.ReportTimestamp
	PendingTransmission *types.PendingTransmission
}

// Persists state from the transmission protocol to the database to allow for recovery
// after restarts
func PersistTransmission(
	ctx context.Context,
	chPersist <-chan TransmissionDBUpdate,
	db ocr3types.Database,
	dbTimeout time.Duration,
	logger loghelper.LoggerWithContext,
) {
	for {
		select {
		case update, ok := <-chPersist:
			if !ok {
				logger.Error("PersistTransmission: chPersist closed unexpectedly, exiting", nil)
				return
			}

			// persisting them. We should obviously fix this, though it does not matter
			// for mercury
			_ = update
			// func() {
			// 	dbCtx, dbCancel := context.WithTimeout(ctx, dbTimeout)
			// 	defer dbCancel()

			// 	store := update.PendingTransmission != nil
			// 	var err error
			// 	if store {
			// 		err = db.StorePendingTransmission(dbCtx, update.Timestamp, *update.PendingTransmission)
			// 	} else {
			// 		err = db.DeletePendingTransmission(dbCtx, update.Timestamp)
			// 	}
			// 	if err != nil {
			// 		logger.ErrorIfNotCanceled(
			// 			"PersistTransmission: error updating database",
			// 			dbCtx,
			// 			commontypes.LogFields{"error": err, "store": store},
			// 		)
			// 		return
			// 	}
			// }()

		case <-ctx.Done():
			logger.Info("PersistTransmission: exiting", nil)
			return
		}
	}
}
