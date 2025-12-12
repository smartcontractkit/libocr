package protocol

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func AttestedToGenesisStateTransitionBlock(prevConfigDigest types.ConfigDigest, astb AttestedStateTransitionBlock) GenesisStateTransitionBlock {
	stb := astb.StateTransitionBlock
	stateWriteSetDigest := MakeStateWriteSetDigest(
		prevConfigDigest,
		stb.BlockSeqNr,
		stb.StateWriteSet.Entries,
	)
	return GenesisStateTransitionBlock{
		stb.PrevHistoryDigest,
		stb.BlockSeqNr,
		stb.StateTransitionInputsDigest,
		stateWriteSetDigest,
		stb.StateRootDigest,
		stb.ReportsPlusPrecursorDigest,
	}
}

func VerifyGenesisStateTransitionBlockFromPrevInstance(cfg ocr3_1config.PublicConfig, gstb GenesisStateTransitionBlock) error {
	prev, ok := cfg.GetPrevFields()
	if !ok {
		return fmt.Errorf("previous instance is not specified in PublicConfig")
	}
	if gstb.SeqNr != prev.PrevSeqNr {
		return fmt.Errorf("genesis state transition block seqNr mismatch, expected PrevSeqNr %d but got %d", prev.PrevSeqNr, gstb.SeqNr)
	}

	actualPrevHistoryDigest := MakeHistoryDigest(
		prev.PrevConfigDigest,
		gstb.PrevHistoryDigest,
		gstb.SeqNr,
		gstb.StateTransitionInputsDigest,
		gstb.StateWriteSetDigest,
		gstb.StateRootDigest,
		gstb.ReportsPlusPrecursorDigest,
	)
	if actualPrevHistoryDigest != prev.PrevHistoryDigest {
		return fmt.Errorf("history digest mismatch, expected PrevHistoryDigest %s but got %s", prev.PrevHistoryDigest, actualPrevHistoryDigest)
	}
	return nil
}
