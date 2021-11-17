package testimplementations

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ types.Database = (*TestDatabase)(nil)

// Noop Database
type TestDatabase struct{}

func (db *TestDatabase) ReadState(context.Context, types.ConfigDigest) (*types.PersistentState, error) {
	return nil, nil
}

func (db *TestDatabase) WriteState(context.Context, types.ConfigDigest, types.PersistentState) error {
	return nil
}

func (db *TestDatabase) ReadConfig(context.Context) (*types.ContractConfig, error) {
	return nil, nil
}

func (db *TestDatabase) WriteConfig(context.Context, types.ContractConfig) error {
	return nil
}

func (db *TestDatabase) StorePendingTransmission(context.Context, types.ReportTimestamp, types.PendingTransmission) error {
	return nil
}

func (db *TestDatabase) PendingTransmissionsWithConfigDigest(context.Context, types.ConfigDigest) (map[types.ReportTimestamp]types.PendingTransmission, error) {
	return nil, nil
}

func (db *TestDatabase) DeletePendingTransmission(context.Context, types.ReportTimestamp) error {
	return nil
}

func (db *TestDatabase) DeletePendingTransmissionsOlderThan(context.Context, time.Time) error {
	return nil
}

// In-memory database

type InMemoryDatabase struct {
	state         types.PersistentState
	config        types.ContractConfig
	transmissions map[types.ReportTimestamp]types.PendingTransmission
}

var _ types.Database = (*InMemoryDatabase)(nil)

func (db *InMemoryDatabase) ReadState(_ context.Context, digest types.ConfigDigest) (*types.PersistentState, error) {
	return &db.state, nil
}

func (db *InMemoryDatabase) WriteState(_ context.Context, digest types.ConfigDigest, state types.PersistentState) error {
	db.state = state
	return nil
}

func (db *InMemoryDatabase) ReadConfig(_ context.Context) (*types.ContractConfig, error) {
	return &db.config, nil
}

func (db *InMemoryDatabase) WriteConfig(_ context.Context, config types.ContractConfig) error {
	db.config = config
	return nil
}

func (db *InMemoryDatabase) StorePendingTransmission(_ context.Context, ts types.ReportTimestamp, transmission types.PendingTransmission) error {
	db.transmissions[ts] = transmission
	return nil
}

func (db *InMemoryDatabase) PendingTransmissionsWithConfigDigest(_ context.Context, digest types.ConfigDigest) (map[types.ReportTimestamp]types.PendingTransmission, error) {
	out := make(map[types.ReportTimestamp]types.PendingTransmission)
	for key, transmission := range db.transmissions {
		if key.ConfigDigest == digest {
			out[key] = transmission
		}
	}
	return out, nil
}

func (db *InMemoryDatabase) DeletePendingTransmission(_ context.Context, ts types.ReportTimestamp) error {
	delete(db.transmissions, ts)
	return nil
}

func (db *InMemoryDatabase) DeletePendingTransmissionsOlderThan(_ context.Context, cutoff time.Time) error {
	clean := make(map[types.ReportTimestamp]types.PendingTransmission)
	removed := make(map[types.ReportTimestamp]types.PendingTransmission)
	for key, transmission := range db.transmissions {
		if transmission.Time.After(cutoff) {
			clean[key] = transmission
		} else {
			removed[key] = transmission
		}
	}
	db.transmissions = clean
	return nil
}

// Factory for in-memory databases

type InMemoryDatabaseFactory struct {
	dbs map[int]*InMemoryDatabase
}

func NewInMemoryDatabaseFactory() *InMemoryDatabaseFactory {
	return &InMemoryDatabaseFactory{
		make(map[int]*InMemoryDatabase),
	}
}

func (d *InMemoryDatabaseFactory) GetDatabase(oracleID int) *InMemoryDatabase {
	db := d.dbs[oracleID]
	return db
}

func (d *InMemoryDatabaseFactory) MakeDatabase(oracleID int) *InMemoryDatabase {
	d.dbs[oracleID] = &InMemoryDatabase{
		state:         types.PersistentState{},
		config:        types.ContractConfig{},
		transmissions: make(map[types.ReportTimestamp]types.PendingTransmission),
	}
	return d.dbs[oracleID]
}
