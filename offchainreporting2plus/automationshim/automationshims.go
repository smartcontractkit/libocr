package automationshim

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type AutomationReportInfo struct {
	Epoch uint32
	Round uint8
}

type AutomationOCR3OnchainKeyring struct {
	ocr2OnchainKeyring types.OnchainKeyring
}

var _ ocr3types.OnchainKeyring[AutomationReportInfo] = &AutomationOCR3OnchainKeyring{}

func NewAutomationOCR3OnchainKeyring(ocr2OnchainKeyring types.OnchainKeyring) *AutomationOCR3OnchainKeyring {
	return &AutomationOCR3OnchainKeyring{ocr2OnchainKeyring}
}

func (ok *AutomationOCR3OnchainKeyring) MaxSignatureLength() int {
	return ok.ocr2OnchainKeyring.MaxSignatureLength()
}

func (ok *AutomationOCR3OnchainKeyring) Sign(configDigest types.ConfigDigest, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[AutomationReportInfo]) (signature []byte, err error) {
	return ok.ocr2OnchainKeyring.Sign(
		types.ReportContext{
			types.ReportTimestamp{
				configDigest,
				reportWithInfo.Info.Epoch,
				reportWithInfo.Info.Round,
			},
			[32]byte{},
		},
		reportWithInfo.Report,
	)
}

func (ok *AutomationOCR3OnchainKeyring) PublicKey() types.OnchainPublicKey {
	return ok.ocr2OnchainKeyring.PublicKey()
}

func (ok *AutomationOCR3OnchainKeyring) Verify(pubkey types.OnchainPublicKey, configDigest types.ConfigDigest, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[AutomationReportInfo], sig []byte) bool {
	return ok.ocr2OnchainKeyring.Verify(
		pubkey,
		types.ReportContext{
			types.ReportTimestamp{
				configDigest,
				reportWithInfo.Info.Epoch,
				reportWithInfo.Info.Round,
			},
			[32]byte{},
		},
		reportWithInfo.Report,
		sig,
	)
}

type AutomationOCR3ContractTransmitter struct {
	ocr2ContractTransmitter types.ContractTransmitter
}

var _ ocr3types.ContractTransmitter[AutomationReportInfo] = &AutomationOCR3ContractTransmitter{}

func NewAutomationOCR3ContractTransmitter(ocr2ContractTransmitter types.ContractTransmitter) *AutomationOCR3ContractTransmitter {
	return &AutomationOCR3ContractTransmitter{ocr2ContractTransmitter}
}

func (t *AutomationOCR3ContractTransmitter) Transmit(
	ctx context.Context,
	configDigest types.ConfigDigest,
	reportWithInfo ocr3types.ReportWithInfo[AutomationReportInfo],
	aoss []types.AttributedOnchainSignature,
) error {
	return t.ocr2ContractTransmitter.Transmit(
		ctx,
		types.ReportContext{
			types.ReportTimestamp{
				configDigest,
				reportWithInfo.Info.Epoch,
				reportWithInfo.Info.Round,
			},
			[32]byte{},
		},
		reportWithInfo.Report,
		aoss,
	)
}

func (t *AutomationOCR3ContractTransmitter) FromAccount() (types.Account, error) {
	return t.ocr2ContractTransmitter.FromAccount()
}

type MercuryOCR3Plugin struct {
	Config ocr3types.OCR3PluginConfig
	Plugin ocr3types.MercuryPlugin
}
