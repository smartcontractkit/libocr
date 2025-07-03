package protocol

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type AttestedReportMany[RI any] struct {
	ReportWithInfo       ocr3types.ReportWithInfo[RI]
	AttributedSignatures []types.AttributedOnchainSignature
}

type StateTransitionBlock struct {
	Epoch                       uint64
	BlockSeqNr                  uint64
	StateTransitionInputsDigest StateTransitionInputsDigest
	StateTransitionOutputs      StateTransitionOutputs
	ReportsPlusPrecursor        ocr3_1types.ReportsPlusPrecursor
}

func (stb *StateTransitionBlock) SeqNr() uint64 {
	return stb.BlockSeqNr
}

type StateTransitionOutputs struct {
	WriteSet []KeyValuePair
}
