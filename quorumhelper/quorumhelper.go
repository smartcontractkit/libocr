package quorumhelper

import (
	"github.com/smartcontractkit/libocr/internal/byzquorum"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
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
	//
	// We discourage use of this quorum for OCR ReportingPlugins. Unlike the quorums
	// used in OCR protocol logic, this quorum is not monotone in f (i.e. decreasing
	// f increases the threshold n−f) and thus the common practice of
	// setting a non-maximum f (i.e. choosing f s.t. 3f+1 < n) to improve
	// availability in case of crash-faults while reducing byzantine-fault
	// tolerance would actually *reduce* availability when this quorum is used.
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
