package shim

import (
	"context"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type SerializingOCR3_1Database struct {
	BinaryDb ocr3_1types.Database
}

var _ protocol.Database = (*SerializingOCR3_1Database)(nil)

func (db *SerializingOCR3_1Database) ReadConfig(ctx context.Context) (*types.ContractConfig, error) {
	return db.BinaryDb.ReadConfig(ctx)
}

func (db *SerializingOCR3_1Database) WriteConfig(ctx context.Context, config types.ContractConfig) error {
	return db.BinaryDb.WriteConfig(ctx, config)
}

func (db *SerializingOCR3_1Database) ReadPacemakerState(ctx context.Context, configDigest types.ConfigDigest) (protocol.PacemakerState, error) {
	raw, err := db.BinaryDb.ReadProtocolState(ctx, configDigest, pacemakerKey)
	if err != nil {
		return protocol.PacemakerState{}, err
	}

	if len(raw) == 0 {
		return protocol.PacemakerState{}, nil
	}

	return serialization.DeserializePacemakerState(raw)
}

func (db *SerializingOCR3_1Database) WritePacemakerState(ctx context.Context, configDigest types.ConfigDigest, state protocol.PacemakerState) error {
	raw, err := serialization.SerializePacemakerState(state)
	if err != nil {
		return err
	}

	return db.BinaryDb.WriteProtocolState(ctx, configDigest, pacemakerKey, raw)
}

func (db *SerializingOCR3_1Database) ReadCert(ctx context.Context, configDigest types.ConfigDigest) (protocol.CertifiedPrepareOrCommit, error) {
	raw, err := db.BinaryDb.ReadProtocolState(ctx, configDigest, certKey)
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return nil, nil
	}

	// This oracle wrote the PrepareOrCommit, so it's fine to trust the value.
	return serialization.DeserializeTrustedPrepareOrCommit(raw)
}

// Writing with an empty value is the same as deleting.
func (db *SerializingOCR3_1Database) WriteCert(ctx context.Context, configDigest types.ConfigDigest, cert protocol.CertifiedPrepareOrCommit) error {
	if cert == nil {
		return db.BinaryDb.WriteProtocolState(ctx, configDigest, certKey, nil)
	}

	raw, err := serialization.SerializeCertifiedPrepareOrCommit(cert)
	if err != nil {
		return err
	}

	return db.BinaryDb.WriteProtocolState(ctx, configDigest, certKey, raw)
}
