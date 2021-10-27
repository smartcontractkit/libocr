package tracing

import (
	"context"
	"sync"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// lockedContractTransmitter will lock access to the backend contract transmitter
// interface so it can be shared by multiple oracles.
type lockedContractTransmitter struct {
	types.ContractTransmitter
	mu sync.Mutex
}

// Transmit needs to be protected by a mutex when twins are sharing the same transmitter instance.
// The twins may transmit transactions to the simulated backend at the same time,
// leading to both transactions having the same nonce and a panic from the backend.
func (c *lockedContractTransmitter) Transmit(ctx context.Context, repCtx types.ReportContext, report types.Report, sigs []types.AttributedOnchainSignature) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ContractTransmitter.Transmit(ctx, repCtx, report, sigs)
}

func MakeLockedContractTransmitter(backend types.ContractTransmitter) types.ContractTransmitter {
	return &lockedContractTransmitter{backend, sync.Mutex{}}
}

// ContractTransmitter is a wrapper over types.ContractTransmitter to enable tracing instrumentation.
type ContractTransmitter struct {
	tracer   *Tracer
	oracleID OracleID
	backend  types.ContractTransmitter
}

var _ types.ContractTransmitter = (*ContractTransmitter)(nil)

func MakeContractTransmitter(tracer *Tracer, oracleID OracleID, backend types.ContractTransmitter) *ContractTransmitter {
	return &ContractTransmitter{tracer, oracleID, backend}
}

func (c *ContractTransmitter) Transmit(ctx context.Context, repCtx types.ReportContext, report types.Report, sigs []types.AttributedOnchainSignature) error {
	err := c.backend.Transmit(ctx, repCtx, report, sigs)
	c.tracer.Append(NewTransmit(c.oracleID, repCtx, report, sigs, err))
	return err
}

func (c *ContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (
	configDigest types.ConfigDigest,
	epoch uint32,
	err error,
) {
	configDigest, epoch, err = c.backend.LatestConfigDigestAndEpoch(ctx)
	c.tracer.Append(NewLatestConfigDigestAndEpoch(c.oracleID, configDigest, epoch, err))
	return configDigest, epoch, err
}

func (c *ContractTransmitter) FromAccount() types.Account {
	account := c.backend.FromAccount()
	c.tracer.Append(NewFromAccount(c.oracleID, account))
	return account
}
