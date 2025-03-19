package shim

import (
	"context"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"
)

type SerializingOCR3_1Database struct {
	BinaryDb ocr3_1types.Database
}

var _ protocol.Database = (*SerializingOCR3_1Database)(nil)

const statePersistenceKey = "state"

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

func (db *SerializingOCR3_1Database) ReadStatePersistenceState(ctx context.Context, configDigest types.ConfigDigest) (protocol.StatePersistenceState, error) {
	raw, err := db.BinaryDb.ReadProtocolState(ctx, configDigest, statePersistenceKey)
	if err != nil {
		return protocol.StatePersistenceState{}, err
	}

	if len(raw) == 0 {
		return protocol.StatePersistenceState{}, nil
	}

	return serialization.DeserializeStatePersistenceState(raw)
}

// Writing with an empty value is the same as deleting.
func (db *SerializingOCR3_1Database) WriteStatePersistenceState(ctx context.Context, configDigest types.ConfigDigest, state protocol.StatePersistenceState) error {
	raw, err := serialization.SerializeStatePersistenceState(state)
	if err != nil {
		return err
	}

	return db.BinaryDb.WriteProtocolState(ctx, configDigest, statePersistenceKey, raw)
}

func (db *SerializingOCR3_1Database) ReadAttestedStateTransitionBlock(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64) (protocol.AttestedStateTransitionBlock, error) {
	raw, err := db.BinaryDb.ReadBlock(ctx, configDigest, seqNr)
	if err != nil {
		return protocol.AttestedStateTransitionBlock{}, err
	}

	if len(raw) == 0 {
		return protocol.AttestedStateTransitionBlock{}, nil
	}

	astb := serialization.AttestedStateTransitionBlock{}
	if err := proto.Unmarshal(raw, &astb); err != nil {
		return protocol.AttestedStateTransitionBlock{}, err
	}

	return serialization.DeserializeAttestedStateTransitionBlock(raw)
}

// Writing with an empty value is the same as deleting.
func (db *SerializingOCR3_1Database) WriteAttestedStateTransitionBlock(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64, astb protocol.AttestedStateTransitionBlock) error {
	raw, err := serialization.SerializeAttestedStateTransitionBlock(astb)
	if err != nil {
		return err
	}

	return db.BinaryDb.WriteBlock(ctx, configDigest, seqNr, raw)
}
