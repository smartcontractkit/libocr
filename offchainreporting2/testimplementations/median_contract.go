package testimplementations

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// var _ types.ContractConfigTracker = (*MedianContract)(nil)

type MedianContract struct {
	backend  ContractBackendWithBlockheaders
	contract *ocr2aggregator.OCR2Aggregator
}

func NewMedianContract(address common.Address, backend ContractBackendWithBlockheaders) (*MedianContract, error) {
	contract, err := ocr2aggregator.NewOCR2Aggregator(address, backend)
	if err != nil {
		return nil, err
	}

	return &MedianContract{
		backend,
		contract,
	}, nil
}

func (t *MedianContract) LatestTransmissionDetails(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	latestAnswer *big.Int,
	latestTimestamp time.Time,
	err error,
) {
	result, err := t.contract.LatestTransmissionDetails(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return configDigest, 0, 0, nil, time.Time{}, fmt.Errorf("error getting LatestTransmissionDetails: %w", err)
	}
	return result.ConfigDigest, result.Epoch, result.Round, result.LatestAnswer, time.Unix(int64(result.LatestTimestamp), 0), nil
}

func (t *MedianContract) LatestRoundRequested(
	ctx context.Context,
	lookback time.Duration,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	err error,
) {
	const k = 100
	tip, err := t.backend.HeaderByNumber(ctx, nil)
	if err != nil {
		return types.ConfigDigest{}, 0, 0, errors.Wrap(err, "failed to get tip in LatestRoundRequested")
	}

	// Estimate start block for log filter:
	// go back k blocks.
	// based on that, estimate range of blocks to query
	// then finetune in intervals of k blocks
	start := uint64(0)
	if tip.Number.Uint64() >= k {
		n := big.NewInt(int64(tip.Number.Uint64()) - k)
		tipMinusK, err := t.backend.HeaderByNumber(ctx, n)
		if err != nil {
			return types.ConfigDigest{}, 0, 0, errors.Wrap(err, "failed to get tipMinusK in LatestRoundRequested")
		}

		blockInterval := float64(tip.Time-tipMinusK.Time) / k
		lookbackBlocks := uint64(math.Ceil(lookback.Seconds() / blockInterval))

		for {
			if tip.Number.Uint64() <= lookbackBlocks {
				break
			}
			n := big.NewInt(int64(tip.Number.Uint64() - lookbackBlocks))
			tipMinusLookback, err := t.backend.HeaderByNumber(ctx, n)
			if err != nil {
				return types.ConfigDigest{}, 0, 0, errors.Wrap(err, "failed to get tipMinusLookback in LatestRoundRequested")
			}

			if uint64(lookback.Seconds()) < tip.Time-tipMinusLookback.Time {
				start = tip.Number.Uint64() - lookbackBlocks
				break
			}

			lookbackBlocks += k
		}
	}

	it, err := t.contract.FilterRoundRequested(&bind.FilterOpts{start, nil, ctx}, nil)
	if err != nil {
		return types.ConfigDigest{}, 0, 0, errors.Wrapf(err, "could not search for RoundRequested logs")
	}
	defer it.Close()

	for it.Next() {
		if it.Event != nil {
			configDigest, epoch, round = it.Event.ConfigDigest, it.Event.Epoch, it.Event.Round
		}
	}
	return configDigest, epoch, round, nil

}

func (t *MedianContract) LatestBlockHeight(ctx context.Context) (blockheight uint64, err error) {
	header, err := t.backend.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}
