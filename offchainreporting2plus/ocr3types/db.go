package ocr3types

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type Database interface {
	types.ConfigDatabase
	ProtocolStateDatabase
}

// ProtocolStateDatabase persistently stores protocol state to survive process restarts.
// Expect Write to be called far more frequently than Read.
//
// All its functions should be thread-safe.
type ProtocolStateDatabase interface {
	ReadProtocolState(ctx context.Context, configDigest types.ConfigDigest, key string) ([]byte, error)
	// Writing with an empty value is the same as deleting.
	WriteProtocolState(ctx context.Context, configDigest types.ConfigDigest, key string, value []byte) error
}
