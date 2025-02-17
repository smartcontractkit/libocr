package types

import (
	"context"
	"time"
)

// ConfigDatabase persistently stores configuration-related information on
// disk.
//
// All its functions should be thread-safe.
type ConfigDatabase interface {
	ReadConfig(ctx context.Context) (*ContractConfig, error)
	WriteConfig(ctx context.Context, config ContractConfig) error
}

// Database persistently stores information on-disk.
// All its functions should be thread-safe.
type Database interface {
	ConfigDatabase

	ReadState(ctx context.Context, configDigest ConfigDigest) (*PersistentState, error)
	WriteState(ctx context.Context, configDigest ConfigDigest, state PersistentState) error

	StorePendingTransmission(context.Context, ReportTimestamp, PendingTransmission) error
	PendingTransmissionsWithConfigDigest(context.Context, ConfigDigest) (map[ReportTimestamp]PendingTransmission, error)
	DeletePendingTransmission(context.Context, ReportTimestamp) error
	DeletePendingTransmissionsOlderThan(context.Context, time.Time) error
}

type PendingTransmission struct {
	Time                 time.Time
	ExtraHash            [32]byte
	Report               Report
	AttributedSignatures []AttributedOnchainSignature
}

type PersistentState struct {
	Epoch                uint32
	HighestSentEpoch     uint32
	HighestReceivedEpoch []uint32 // length: at most MaxOracles
}

func (ps PersistentState) Equal(ps2 PersistentState) bool {
	if ps.Epoch != ps2.Epoch {
		return false
	}
	if ps.HighestSentEpoch != ps2.HighestSentEpoch {
		return false
	}
	if len(ps.HighestReceivedEpoch) != len(ps2.HighestReceivedEpoch) {
		return false
	}
	for i := 0; i < len(ps.HighestReceivedEpoch); i++ {
		if ps.HighestReceivedEpoch[i] != ps2.HighestReceivedEpoch[i] {
			return false
		}
	}
	return true
}
