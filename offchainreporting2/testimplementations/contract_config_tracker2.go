package testimplementations

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type ContractBackendWithBlockheaders interface {
	bind.ContractBackend
	HeaderByNumber(ctx context.Context, block *big.Int) (*ethtypes.Header, error)
}

var _ types.ContractConfigTracker = (*ContractConfigTracker2)(nil)

type ContractConfigTracker2 struct {
	backend  ContractBackendWithBlockheaders
	contract *ocr2aggregator.OCR2Abstract
}

func NewContractConfigTracker2(address common.Address, backend ContractBackendWithBlockheaders) (*ContractConfigTracker2, error) {
	contract, err := ocr2aggregator.NewOCR2Abstract(address, backend)
	if err != nil {
		return nil, err
	}

	return &ContractConfigTracker2{
		backend,
		contract,
	}, nil
}

func contractConfigFromConfigSetEvent(
	cc *ocr2aggregator.OCR2AbstractConfigSet,
) *types.ContractConfig {
	if cc == nil {
		return nil
	}
	transmitAccounts := []types.Account{}
	for _, addr := range cc.Transmitters {
		transmitAccounts = append(transmitAccounts, types.Account(addr.Hex()))
	}
	signers := []types.OnchainPublicKey{}
	for _, addr := range cc.Signers {
		addr := addr
		signers = append(signers, types.OnchainPublicKey(addr[:]))
	}
	return &types.ContractConfig{
		cc.ConfigDigest,
		cc.ConfigCount,
		signers,
		transmitAccounts,
		cc.F,
		cc.OnchainConfig,
		cc.OffchainConfigVersion,
		cc.OffchainConfig,
	}
}

func (t *ContractConfigTracker2) Notify() <-chan struct{} {
	return nil
}

func (t *ContractConfigTracker2) LatestConfigDetails(
	ctx context.Context,
) (
	changedInBlock uint64, configDigest types.ConfigDigest, err error,
) {
	details, err := t.contract.LatestConfigDetails(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, types.ConfigDigest{}, err
	}

	configDigest, err = types.BytesToConfigDigest(details.ConfigDigest[:])
	if err != nil {
		return 0, types.ConfigDigest{}, err
	}

	return uint64(details.BlockNumber), configDigest, nil
}

func (t *ContractConfigTracker2) LatestConfig(
	ctx context.Context,
	changedInBlock uint64,
) (
	types.ContractConfig,
	error,
) {
	configChanges, err := t.contract.FilterConfigSet(&bind.FilterOpts{
		Start:   changedInBlock,
		End:     &changedInBlock,
		Context: ctx,
	})
	if err != nil {
		return types.ContractConfig{}, errors.Wrapf(err, "could not search for ContractConfigs")
	}
	defer configChanges.Close()
	latestChangePtr := (*types.ContractConfig)(nil)
	for {
		latestChangePtr = contractConfigFromConfigSetEvent(configChanges.Event)
		if !configChanges.Next() {
			break
		}
	}
	if latestChangePtr == nil {
		return types.ContractConfig{}, errors.Errorf("found no config in block %d", changedInBlock)
	}
	return *latestChangePtr, nil
}

func (t *ContractConfigTracker2) LatestBlockHeight(ctx context.Context) (blockheight uint64, err error) {
	header, err := t.backend.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}
