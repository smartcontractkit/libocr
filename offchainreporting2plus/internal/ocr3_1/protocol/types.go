package protocol

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type AttestedReportMany[RI any] struct {
	ReportWithInfo       ocr3types.ReportWithInfo[RI]
	AttributedSignatures []types.AttributedOnchainSignature
}

type StateTransitionBlock struct {
	StateTransitionInputs       StateTransitionInputs
	StateTransitionOutputDigest StateTransitionOutputDigest
	ReportsPrecursorDigest      ReportsPlusPrecursorDigest
}

func (stb *StateTransitionBlock) SeqNr() uint64 {
	return stb.StateTransitionInputs.SeqNr
}

type StateTransitionInputs struct {
	SeqNr                  uint64
	Epoch                  uint64
	Round                  uint64
	Query                  types.Query
	AttributedObservations []types.AttributedObservation
}

func (stis *StateTransitionInputs) isGenesis() bool {
	return stis.SeqNr == 0 && stis.Epoch == 0 && stis.Round == 0 &&
		len(stis.Query) == 0 && len(stis.AttributedObservations) == 0
}
