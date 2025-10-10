package protocol

const (
	SnapshotInterval = 128
	// MaxHistoricalSnapshotsRetained must be a non-zero value, denoting the
	// number of complete snapshots prior to the current (potentially
	// incomplete) one that will be retained to help other oracles with state
	// sync. All blocks starting from the highest block of the earliest retained
	// snapshot will be retained.
	MaxHistoricalSnapshotsRetained = 64
)

func snapshotIndexFromSeqNr(seqNr uint64) uint64 {
	if seqNr == 0 {
		return 0
	}
	return (seqNr + SnapshotInterval - 1) / SnapshotInterval
}

func maxSeqNrWithSnapshotIndex(snapshotIndex uint64) uint64 {
	if snapshotIndex == 0 {
		return 0
	}
	return snapshotIndex * SnapshotInterval
}

func desiredLowestPersistedSeqNr(highestCommittedSeqNr uint64) uint64 {
	highestSnapshotIndex := snapshotIndexFromSeqNr(highestCommittedSeqNr)
	var lowestDesiredSnapshotIndex uint64
	if highestSnapshotIndex > MaxHistoricalSnapshotsRetained {
		lowestDesiredSnapshotIndex = highestSnapshotIndex - MaxHistoricalSnapshotsRetained
	} else {
		lowestDesiredSnapshotIndex = 0
	}
	return maxSeqNrWithSnapshotIndex(lowestDesiredSnapshotIndex)
}

func snapshotSeqNr(seqNr uint64) uint64 {
	return maxSeqNrWithSnapshotIndex(snapshotIndexFromSeqNr(seqNr))
}

// prevRootVersion returns the version number of the JMT root referring to the
// state as of seqNr - 1. This is used as the "old version" for writing the
// modifications of seqNr. We only maintain trees with versions that are
// multiples of SnapshotInterval.
func PrevRootVersion(seqNr uint64) uint64 {
	return snapshotSeqNr(seqNr - 1)
}

func RootVersion(seqNr uint64) uint64 {
	return snapshotSeqNr(seqNr)
}
