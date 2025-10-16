package singlewriter

import (
	"fmt"
	"sync"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
)

type UnserializedTransaction struct {
	*overlayTransaction
}

var _ ocr3_1types.KeyValueDatabaseReadWriteTransaction = &UnserializedTransaction{}

func NewUnserializedTransaction(
	keyValueDatabase ocr3_1types.KeyValueDatabase,
) (ocr3_1types.KeyValueDatabaseReadWriteTransaction, error) {
	rawReaderTx, err := keyValueDatabase.NewReadTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to create read transaction to create an UnserializedTransactionImpl: %w", err)
	}
	return &UnserializedTransaction{
		&overlayTransaction{
			make(map[string]operation),
			keyValueDatabase,
			rawReaderTx,
			sync.Mutex{},
			sync.WaitGroup{},
			txStatusOpen,
		},
	}, nil
}

func (ut *UnserializedTransaction) Commit() error {
	ut.mu.Lock()
	defer ut.mu.Unlock()

	if ut.status != txStatusOpen {
		return ErrClosed
	}
	overlay := ut.overlay
	ut.lockedDiscard()

	rawReadWriteTx, err := ut.keyValueDatabase.NewReadWriteTransaction()
	if err != nil {
		return fmt.Errorf("failed to create read write transaction to commit UnserializedTransaction: %w", err)
	}

	defer rawReadWriteTx.Discard()

	err = lockedCommit(overlay, rawReadWriteTx)
	if err != nil {
		return err
	}

	ut.status = txStatusCommitted
	ut.overlay = nil
	return nil
}
