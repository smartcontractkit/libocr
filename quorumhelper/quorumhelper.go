package quorumhelper

import (
	"github.com/RoSpaceDev/libocr/internal/byzquorum"
	"github.com/RoSpaceDev/libocr/offchainreporting2/types"
)

type Quorum int

const (
	_ Quorum = iota
	// Guarantees at least one honest observation
	QuorumFPlusOne
	// Guarantees an honest majority of observations
	QuorumTwoFPlusOne
	// Guarantees that all sets of observations overlap in at least one honest oracle
	QuorumByzQuorum
	// Maximal number of observations we can rely on being available
	QuorumNMinusF
)

func ObservationCountReachesObservationQuorum(quorum Quorum, n, f int, aos []types.AttributedObservation) bool {
	switch quorum {
	case QuorumFPlusOne:
		return len(aos) >= f+1
	case QuorumTwoFPlusOne:
		return len(aos) >= 2*f+1
	case QuorumByzQuorum:
		return len(aos) >= byzquorum.Size(n, f)
	case QuorumNMinusF:
		return len(aos) >= n-f
	default:
		panic("Unknown quorum")
	}
}
