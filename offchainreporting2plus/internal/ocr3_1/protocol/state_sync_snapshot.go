package protocol

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
)

func snapshotIndexFromSeqNr(seqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	if seqNr == 0 {
		return 0
	}
	return (seqNr + config.GetSnapshotInterval() - 1) / config.GetSnapshotInterval()
}

func maxSeqNrWithSnapshotIndex(snapshotIndex uint64, config ocr3_1config.PublicConfig) uint64 {
	if snapshotIndex == 0 {
		return 0
	}
	return snapshotIndex * config.GetSnapshotInterval()
}

func desiredLowestPersistedSeqNr(highestCommittedSeqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	highestSnapshotIndex := snapshotIndexFromSeqNr(highestCommittedSeqNr, config)
	var lowestDesiredSnapshotIndex uint64
	cfgMaxHistoricalSnapshotsRetained := config.GetMaxHistoricalSnapshotsRetained()
	if highestSnapshotIndex > cfgMaxHistoricalSnapshotsRetained {
		lowestDesiredSnapshotIndex = highestSnapshotIndex - cfgMaxHistoricalSnapshotsRetained
	} else {
		lowestDesiredSnapshotIndex = 0
	}
	return maxSeqNrWithSnapshotIndex(lowestDesiredSnapshotIndex, config)
}

func snapshotSeqNr(seqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	return maxSeqNrWithSnapshotIndex(snapshotIndexFromSeqNr(seqNr, config), config)
}

// prevRootVersion returns the version number of the JMT root referring to the
// state as of seqNr - 1. This is used as the "old version" for writing the
// modifications of seqNr. We only maintain trees with versions that are
// multiples of SnapshotInterval.
func PrevRootVersion(seqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	return snapshotSeqNr(seqNr-1, config)
}

func RootVersion(seqNr uint64, config ocr3_1config.PublicConfig) uint64 {
	return snapshotSeqNr(seqNr, config)
}
