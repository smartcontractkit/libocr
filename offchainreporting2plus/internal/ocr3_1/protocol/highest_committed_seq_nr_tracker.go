package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
)

func RunHighestCommittedSeqNrTracker(
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	kvDb KeyValueDatabase,
	interval time.Duration,
	chOut chan<- uint64,
) {
	chDone := ctx.Done()
	chTick := time.After(0)

	highestCommittedSeqNr := uint64(0)
	notifyOut := false
	for {
		var nilOrChOut chan<- uint64
		if notifyOut {
			nilOrChOut = chOut
		} else {
			nilOrChOut = nil
		}

		select {
		case nilOrChOut <- highestCommittedSeqNr:
			notifyOut = false
		case <-chTick:
		case <-chDone:
			return
		}

		tx, err := kvDb.NewReadTransactionUnchecked()
		if err != nil {
			logger.Warn("RunHighestCommittedSeqNrTracker: failed to create read transaction", commontypes.LogFields{
				"error": err,
			})
			chTick = time.After(interval)
			continue
		}

		highestCommittedSeqNrNew, err := tx.ReadHighestCommittedSeqNr()
		tx.Discard()
		if err != nil {
			logger.Warn("RunHighestCommittedSeqNrTracker: failed to read highest committed seq nr", commontypes.LogFields{
				"error": err,
			})
			chTick = time.After(interval)
			continue
		}

		chTick = time.After(interval)

		if highestCommittedSeqNrNew > highestCommittedSeqNr {
			logger.Trace("RunHighestCommittedSeqNrTracker: new highestCommittedSeqNr", commontypes.LogFields{
				"highestCommittedSeqNr":     highestCommittedSeqNrNew,
				"prevHighestCommittedSeqNr": highestCommittedSeqNr,
			})
			notifyOut = true
			highestCommittedSeqNr = highestCommittedSeqNrNew
		}
	}
}
