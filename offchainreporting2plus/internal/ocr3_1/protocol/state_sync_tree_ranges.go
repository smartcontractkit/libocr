package protocol

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/RoSpaceDev/libocr/internal/jmt"
)

// KeyDigestRange represents a contiguous range [StartIndex, EndInclIndex] in
// the key digest space that needs to be fetched during tree synchronization.
type KeyDigestRange struct {
	StartIndex   jmt.Digest
	EndInclIndex jmt.Digest
}

// splitKeyDigestSpaceEvenly splits the key digest space [0, max] into n even ranges.
// Returns an error if n >= 256 (split granularity would exceed byte precision).
func splitKeyDigestSpaceEvenly(n int) ([]KeyDigestRange, error) {
	if n < 1 {
		return nil, fmt.Errorf("n must be at least 1, got %d", n)
	}
	if n >= 256 {
		return nil, fmt.Errorf("n must be less than 256, got %d", n)
	}

	ranges := make([]KeyDigestRange, 0, n)
	startIndices := make([]jmt.Digest, 0, n)

	for i := range n {
		var startIndex jmt.Digest
		startIndex[0] = byte(i * 256 / n)
		startIndices = append(startIndices, startIndex)
	}

	for i, startIndex := range startIndices {
		var endInclIndex jmt.Digest
		if i+1 == len(startIndices) {
			endInclIndex = jmt.MaxDigest
		} else {
			var ok bool
			endInclIndex, ok = jmt.DecrementDigest(startIndices[i+1])
			if !ok {
				return nil, fmt.Errorf("unexpected: could not decrement nonzero digest")
			}
		}
		ranges = append(ranges, KeyDigestRange{startIndex, endInclIndex})
	}
	return ranges, nil
}

// PendingKeyDigestRanges tracks which key digest ranges still need to be
// fetched during tree synchronization. As chunks are received, the
// corresponding ranges are removed or updated.
type PendingKeyDigestRanges struct {
	ranges []KeyDigestRange
}

func NewPendingKeyDigestRanges(ranges []KeyDigestRange) PendingKeyDigestRanges {
	return PendingKeyDigestRanges{ranges}
}

// WithReceivedRange returns a new PendingKeyDigestRanges with the given range
// marked as received. Does not mutate the receiver.
func (pkdr PendingKeyDigestRanges) WithReceivedRange(receivedRange KeyDigestRange) PendingKeyDigestRanges {
	// Find the range with the startIndex of the received chunk
	i := slices.IndexFunc(pkdr.ranges, func(r KeyDigestRange) bool {
		return r.StartIndex == receivedRange.StartIndex
	})
	if i == -1 {
		// Range not found - return unchanged
		return pkdr
	}

	// Make a copy of the ranges slice to avoid mutating the original
	newRanges := slices.Clone(pkdr.ranges)

	nextStartIndex, ok := jmt.IncrementDigest(receivedRange.EndInclIndex)
	if !ok || bytes.Compare(receivedRange.EndInclIndex[:], newRanges[i].EndInclIndex[:]) >= 0 {
		// The received range covers the entire pending range - remove it
		newRanges = slices.Delete(newRanges, i, i+1)
	} else {
		// The received range covers only part of the pending range - update it
		newRanges[i].StartIndex = nextStartIndex
	}

	return PendingKeyDigestRanges{newRanges}
}

// All returns all pending key digest ranges that still need to be fetched.
func (pkdr PendingKeyDigestRanges) All() []KeyDigestRange {
	return pkdr.ranges
}
