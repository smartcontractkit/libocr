package ocrintegrationtesthelpers

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type MockContractConfigTracker struct {
	offchainConfigDigester types.OffchainConfigDigester

	mu      sync.RWMutex
	configs []types.ContractConfig
}

var _ types.ContractConfigTracker = &MockContractConfigTracker{}

// NewMockContractConfigTracker returns a tracker that simulates on-chain config
// history without a chain: LatestBlockHeight is the number of configs set via
// SetConfig, and LatestConfigDetails returns changedInBlock as the index of the
// latest config (0 when there are no configs set, 1 when one config is set,
// etc.) LatestBlockHeight is always +1 to the index/changedInBlock of the
// latest config.
//
// Requires ContractConfigConfirmations == 1 or SkipContractConfigConfirmations
// in [types.LocalConfig].
//
// Developers must invoke SetConfig at least once to set a meaningful config,
// otherwise oracles using this tracker will not be able to start.
func NewMockContractConfigTracker(offchainConfigDigester types.OffchainConfigDigester) *MockContractConfigTracker {
	configs := []types.ContractConfig{{}}
	return &MockContractConfigTracker{offchainConfigDigester, sync.RWMutex{}, configs}
}

// SetConfig does not perform *any* validity checks on the supplied config. You
// are responsible for ensuring that the config is valid for your use case.
func (c *MockContractConfigTracker) SetConfig(signers []types.OnchainPublicKey, transmitters []types.Account, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	config := types.ContractConfig{types.ConfigDigest{}, uint64(len(c.configs)), signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig}
	configDigest, err := c.offchainConfigDigester.ConfigDigest(context.TODO(), config)
	if err != nil {
		return fmt.Errorf("error computing config digest: %w", err)
	}
	config.ConfigDigest = configDigest
	c.configs = append(c.configs, config)
	return nil
}

func (c *MockContractConfigTracker) Notify() <-chan struct{} {
	return nil
}

func (c *MockContractConfigTracker) LatestConfigDetails(ctx context.Context) (uint64, types.ConfigDigest, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return uint64(len(c.configs) - 1), c.configs[len(c.configs)-1].ConfigDigest, nil
}

func (c *MockContractConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !(0 < changedInBlock && changedInBlock < uint64(len(c.configs))) {
		return types.ContractConfig{}, fmt.Errorf("changedInBlock %d is out of range (0, %d)", changedInBlock, len(c.configs))
	}
	return c.configs[changedInBlock], nil
}

func (c *MockContractConfigTracker) LatestBlockHeight(ctx context.Context) (uint64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return uint64(len(c.configs)), nil
}
