package protocol

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type Database interface {
	types.ConfigDatabase

	ReadCert(ctx context.Context, configDigest types.ConfigDigest) (CertifiedPrepareOrCommit, error)
	WriteCert(ctx context.Context, configDigest types.ConfigDigest, cert CertifiedPrepareOrCommit) error
}
