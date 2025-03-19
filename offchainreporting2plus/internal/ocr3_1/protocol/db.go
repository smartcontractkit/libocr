package protocol

import (
	"context"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type PacemakerState struct {
	Epoch                   uint64
	HighestSentNewEpochWish uint64
}

type StatePersistenceState struct {
	HighestPersistedStateTransitionBlockSeqNr uint64
}

type Database interface {
	types.ConfigDatabase

	ReadPacemakerState(ctx context.Context, configDigest types.ConfigDigest) (PacemakerState, error)
	WritePacemakerState(ctx context.Context, configDigest types.ConfigDigest, state PacemakerState) error

	ReadCert(ctx context.Context, configDigest types.ConfigDigest) (CertifiedPrepareOrCommit, error)
	WriteCert(ctx context.Context, configDigest types.ConfigDigest, cert CertifiedPrepareOrCommit) error

	ReadStatePersistenceState(ctx context.Context, configDigest types.ConfigDigest) (StatePersistenceState, error)
	WriteStatePersistenceState(ctx context.Context, configDigest types.ConfigDigest, state StatePersistenceState) error

	ReadAttestedStateTransitionBlock(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64) (AttestedStateTransitionBlock, error)
	WriteAttestedStateTransitionBlock(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64, ast AttestedStateTransitionBlock) error
}
