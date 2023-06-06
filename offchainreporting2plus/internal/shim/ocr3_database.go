package shim

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"
)

type SerializingOCR3Database struct {
	BinaryDb ocr3types.Database
}

var _ protocol.Database = (*SerializingOCR3Database)(nil)

const certKey = "cert"

func (db *SerializingOCR3Database) ReadConfig(ctx context.Context) (*types.ContractConfig, error) {
	return db.BinaryDb.ReadConfig(ctx)
}

func (db *SerializingOCR3Database) WriteConfig(ctx context.Context, config types.ContractConfig) error {
	return db.BinaryDb.WriteConfig(ctx, config)
}

func (db *SerializingOCR3Database) ReadCert(ctx context.Context, configDigest types.ConfigDigest) (protocol.CertifiedPrepareOrCommit, error) {
	raw, err := db.BinaryDb.ReadProtocolState(ctx, configDigest, certKey)
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return nil, nil
	}

	p := serialization.CertifiedPrepareOrCommit{}
	if err := proto.Unmarshal(raw, &p); err != nil {
		return nil, err
	}

	return serialization.CertifiedPrepareOrCommitFromProtoMessage(&p)
}

// Writing with an empty value is the same as deleting.
func (db *SerializingOCR3Database) WriteCert(ctx context.Context, configDigest types.ConfigDigest, cert protocol.CertifiedPrepareOrCommit) error {
	p := serialization.CertifiedPrepareOrCommitToProtoMessage(cert)

	raw, err := proto.Marshal(p)
	if err != nil {
		return err
	}

	return db.BinaryDb.WriteProtocolState(ctx, configDigest, certKey, raw)
}
