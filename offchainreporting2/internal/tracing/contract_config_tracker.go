package tracing

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type ContractConfigTracker struct {
	tracer   *Tracer
	oracleID OracleID
	backend  types.ContractConfigTracker
}

var _ types.ContractConfigTracker = (*ContractConfigTracker)(nil)

func MakeContractConfigTracker(tracer *Tracer, oracleID OracleID, backend types.ContractConfigTracker) *ContractConfigTracker {
	return &ContractConfigTracker{
		tracer,
		oracleID,
		backend,
	}
}

func (c *ContractConfigTracker) Notify() <-chan struct{} {
	c.tracer.Append(NewNotify(c.oracleID))
	return c.backend.Notify()
}

func (c *ContractConfigTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	changedInBlock, configDigest, err = c.backend.LatestConfigDetails(ctx)
	c.tracer.Append(NewLatestConfigDetails(c.oracleID, changedInBlock, configDigest, err))
	return changedInBlock, configDigest, err
}

func (c *ContractConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	cfg, err := c.backend.LatestConfig(ctx, changedInBlock)
	c.tracer.Append(NewLatestConfig(c.oracleID, changedInBlock, cfg, err))
	return cfg, err
}

func (c *ContractConfigTracker) LatestBlockHeight(ctx context.Context) (blockheight uint64, err error) {
	blockheight, err = c.backend.LatestBlockHeight(ctx)
	c.tracer.Append(NewLatestBlockHeight(c.oracleID, blockheight, err))
	return blockheight, err
}
