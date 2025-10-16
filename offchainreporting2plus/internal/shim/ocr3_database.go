package shim

import (
	"context"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3/protocol"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3/serialization"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type SerializingOCR3Database struct {
	BinaryDb ocr3types.Database
}

var _ protocol.Database = (*SerializingOCR3Database)(nil)

const pacemakerKey = "pacemaker"

const certKey = "cert"

func (db *SerializingOCR3Database) ReadConfig(ctx context.Context) (*types.ContractConfig, error) {
	return db.BinaryDb.ReadConfig(ctx)
}

func (db *SerializingOCR3Database) WriteConfig(ctx context.Context, config types.ContractConfig) error {
	return db.BinaryDb.WriteConfig(ctx, config)
}

func (db *SerializingOCR3Database) ReadPacemakerState(ctx context.Context, configDigest types.ConfigDigest) (protocol.PacemakerState, error) {
	raw, err := db.BinaryDb.ReadProtocolState(ctx, configDigest, pacemakerKey)
	if err != nil {
		return protocol.PacemakerState{}, err
	}

	if len(raw) == 0 {
		return protocol.PacemakerState{}, nil
	}

	return serialization.DeserializePacemakerState(raw)
}

func (db *SerializingOCR3Database) WritePacemakerState(ctx context.Context, configDigest types.ConfigDigest, state protocol.PacemakerState) error {
	raw, err := serialization.SerializePacemakerState(state)
	if err != nil {
		return err
	}

	return db.BinaryDb.WriteProtocolState(ctx, configDigest, pacemakerKey, raw)
}

func (db *SerializingOCR3Database) ReadCert(ctx context.Context, configDigest types.ConfigDigest) (protocol.CertifiedPrepareOrCommit, error) {
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
func (db *SerializingOCR3Database) WriteCert(ctx context.Context, configDigest types.ConfigDigest, cert protocol.CertifiedPrepareOrCommit) error {
	if cert == nil {
		return db.BinaryDb.WriteProtocolState(ctx, configDigest, certKey, nil)
	}

	raw, err := serialization.SerializeCertifiedPrepareOrCommit(cert)
	if err != nil {
		return err
	}

	return db.BinaryDb.WriteProtocolState(ctx, configDigest, certKey, raw)
}
