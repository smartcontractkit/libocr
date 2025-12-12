package protocol

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
)

func genesisSeqNr(config ocr3_1config.PublicConfig) uint64 {
	if prev, ok := config.GetPrevFields(); ok {
		return prev.PrevSeqNr
	}
	return 0
}

func snapshotIndexFromSeqNrAssumingZeroGenesis(seqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	if seqNr == 0 {
		return 0
	}
	return (seqNr + config.GetSnapshotInterval() - 1) / config.GetSnapshotInterval()
}

func maxSeqNrWithSnapshotIndexAssumingZeroGenesis(snapshotIndex uint64, config ocr3_1config.PublicConfig) uint64 {
	return snapshotIndex * config.GetSnapshotInterval()
}

func highestCompleteSnapshotSeqNrNotAbove(seqNr uint64, config ocr3_1config.PublicConfig) (uint64, bool) {
	if isCompleteSnapshotSeqNr(seqNr, config) {
		return seqNr, true
	}
	snapshotIndex := snapshotIndexFromSeqNrAssumingZeroGenesis(seqNr, config)
	if snapshotIndex > 0 {
		prevSnapshotMaxSeqNr := maxSeqNrWithSnapshotIndexAssumingZeroGenesis(snapshotIndex-1, config)
		prevSnapshotMaxSeqNrConsideringGenesis := max(genesisSeqNr(config), prevSnapshotMaxSeqNr)
		if prevSnapshotMaxSeqNrConsideringGenesis == 0 || prevSnapshotMaxSeqNrConsideringGenesis > seqNr {
			return 0, false
		} else {
			return prevSnapshotMaxSeqNrConsideringGenesis, true
		}
	}
	return 0, false
}

func isCompleteSnapshotSeqNr(seqNr uint64, config ocr3_1config.PublicConfig) bool {
	if seqNr == 0 {
		return false
	}
	gen := genesisSeqNr(config)
	if seqNr < gen {
		return false
	}
	return RootVersion(seqNr, config) == seqNr
}

func desiredLowestPersistedSeqNr(highestCommittedSeqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	gen := genesisSeqNr(config)
	if highestCommittedSeqNr < gen {
		return 0 // ensures that lowest <= highest
	}

	highestSnapshotIndex := snapshotIndexFromSeqNrAssumingZeroGenesis(highestCommittedSeqNr, config)
	var lowestDesiredSnapshotIndex uint64
	cfgMaxHistoricalSnapshotsRetained := config.GetMaxHistoricalSnapshotsRetained()
	if highestSnapshotIndex > cfgMaxHistoricalSnapshotsRetained {
		lowestDesiredSnapshotIndex = highestSnapshotIndex - cfgMaxHistoricalSnapshotsRetained
	} else {
		lowestDesiredSnapshotIndex = 0
	}

	lowestDesiredSeqNrAssumingZeroGenesis := maxSeqNrWithSnapshotIndexAssumingZeroGenesis(lowestDesiredSnapshotIndex, config)
	return max(gen, lowestDesiredSeqNrAssumingZeroGenesis)
}

func snapshotSeqNr(seqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	gen := genesisSeqNr(config)
	if seqNr < gen {
		return 0
	}
	if seqNr <= gen {
		return gen
	}
	return maxSeqNrWithSnapshotIndexAssumingZeroGenesis(snapshotIndexFromSeqNrAssumingZeroGenesis(seqNr, config), config)
}

// prevRootVersion returns the version number of the JMT root referring to the
// state as of seqNr - 1. This is used as the "old version" for writing the
// modifications of seqNr. We only maintain trees with versions that are
// multiples of SnapshotInterval.
func PrevRootVersion(seqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	if seqNr <= 0 {
		return 0
	}
	return snapshotSeqNr(seqNr-1, config)
}

func RootVersion(seqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	return snapshotSeqNr(seqNr, config)
}
