package testimplementations

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

var _ types.ContractTransmitter = (*ContractTransmitter2)(nil)

type ContractTransmitter2 struct {
	contract *ocr2aggregator.OCR2Abstract
	opts     *bind.TransactOpts
}

func NewContractTransmitter2(address common.Address, backend bind.ContractBackend, txopts *bind.TransactOpts) (*ContractTransmitter2, error) {
	contract, err := ocr2aggregator.NewOCR2Abstract(address, backend)
	if err != nil {
		return nil, err
	}

	return &ContractTransmitter2{
		contract,
		txopts,
	}, nil
}

func (ct *ContractTransmitter2) Transmit(
	ctx context.Context,
	repctx types.ReportContext,
	report types.Report,
	ass []types.AttributedOnchainSignature,
) error {
	_, err := ct.TransmitAndReturnTx(ctx, repctx, report, ass)
	return err
}

func (ct *ContractTransmitter2) TransmitAndReturnTx(
	ctx context.Context,
	repctx types.ReportContext,
	report types.Report,
	sigs []types.AttributedOnchainSignature,
) (
	*gethtypes.Transaction, error,
) {
	txopts := *ct.opts
	txopts.Context = ctx
	txopts.GasLimit = 500_000

	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	for i, as := range sigs {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			panic("eventTransmit(ev): error in SplitSignature")
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}
	rawReportContext := evmutil.RawReportContext(repctx)
	return ct.contract.Transmit(&txopts, rawReportContext, report, rs, ss, vs)
}

func (ct *ContractTransmitter2) LatestConfigDigestAndEpoch(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	err error,
) {
	latestConfigDigestAndEpoch, err := ct.contract.LatestConfigDigestAndEpoch(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return types.ConfigDigest{}, 0, err
	}

	if !latestConfigDigestAndEpoch.ScanLogs {
		return types.ConfigDigest(latestConfigDigestAndEpoch.ConfigDigest), latestConfigDigestAndEpoch.Epoch, nil
	}

	it, err := ct.contract.FilterTransmitted(&bind.FilterOpts{
		0,
		nil,
		ctx,
	})
	if err != nil {
		return types.ConfigDigest{}, 0, err
	}
	defer it.Close()
	for it.Next() {
		fmt.Println("LatestConfigDigestAndEpoch:", it.Event)
		configDigest = it.Event.ConfigDigest
		epoch = it.Event.Epoch
	}

	if it.Error() != nil {
		return types.ConfigDigest{}, 0, it.Error()
	}
	return configDigest, epoch, nil
}

func (ct *ContractTransmitter2) FromAccount() types.Account {
	return types.Account(ct.opts.From.Hex())
}
