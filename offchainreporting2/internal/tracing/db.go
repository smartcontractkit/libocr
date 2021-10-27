package tracing

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// Database is a wrapper of an instance of type.Database which reports frames to the tracer.
type Database struct {
	tracer   *Tracer
	oracleID OracleID
	backend  types.Database
}

var _ types.Database = (*Database)(nil)

func MakeDatabase(tracer *Tracer, oracleID OracleID, backend types.Database) *Database {
	return &Database{tracer, oracleID, backend}
}

func (db *Database) ReadState(ctx context.Context, digest types.ConfigDigest) (*types.PersistentState, error) {
	state, err := db.backend.ReadState(ctx, digest)
	db.tracer.Append(NewReadState(db.oracleID, digest, *state, err))
	return state, err
}

func (db *Database) WriteState(ctx context.Context, digest types.ConfigDigest, state types.PersistentState) error {
	err := db.backend.WriteState(ctx, digest, state)
	db.tracer.Append(NewWriteState(db.oracleID, digest, state, err))
	return err
}

func (db *Database) ReadConfig(ctx context.Context) (*types.ContractConfig, error) {
	config, err := db.backend.ReadConfig(ctx)
	db.tracer.Append(NewReadConfig(db.oracleID, *config, err))
	return config, err
}

func (db *Database) WriteConfig(ctx context.Context, config types.ContractConfig) error {
	err := db.backend.WriteConfig(ctx, config)
	db.tracer.Append(NewWriteConfig(db.oracleID, config, err))
	return err
}

func (db *Database) StorePendingTransmission(ctx context.Context, ts types.ReportTimestamp, transmission types.PendingTransmission) error {
	err := db.backend.StorePendingTransmission(ctx, ts, transmission)
	db.tracer.Append(NewStorePendingTransmission(db.oracleID, ts, transmission, err))
	return err
}

func (db *Database) PendingTransmissionsWithConfigDigest(ctx context.Context, digest types.ConfigDigest) (map[types.ReportTimestamp]types.PendingTransmission, error) {
	transmissions, err := db.backend.PendingTransmissionsWithConfigDigest(ctx, digest)
	db.tracer.Append(NewPendingTransmissionsWithConfigDigest(db.oracleID, digest, err))
	return transmissions, err
}

func (db *Database) DeletePendingTransmission(ctx context.Context, ts types.ReportTimestamp) error {
	err := db.backend.DeletePendingTransmission(ctx, ts)
	db.tracer.Append(NewDeletePendingTransmission(db.oracleID, ts, err))
	return err
}

func (db *Database) DeletePendingTransmissionsOlderThan(ctx context.Context, cutoff time.Time) error {
	err := db.backend.DeletePendingTransmissionsOlderThan(ctx, cutoff)
	db.tracer.Append(NewDeletePendingTransmissionsOlderThan(db.oracleID, cutoff, err))
	return err
}
