package ocr3_1types

import (
	"context"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type Database interface {
	ocr3types.Database
	BlockDatabase
}

type BlockNotFoundError string

func (e BlockNotFoundError) Error() string {
	return string(e)
}

const ErrBlockNotFound BlockNotFoundError = "block not found"

// BlockDatabase persistently stores state transition blocks to support state transfer requests
// Expect Write to be called far more frequently than Read.
//
// All its functions should be thread-safe.

type BlockDatabase interface {
	// ReadBlock retrieves a block from the database.
	// If the block is not found, ErrBlockNotFound should be returned.
	ReadBlock(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64) ([]byte, error)
	// WriteBlock writes a block to the database.
	// Writing with a nil value is the same as deleting.
	WriteBlock(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64, block []byte) error
}
