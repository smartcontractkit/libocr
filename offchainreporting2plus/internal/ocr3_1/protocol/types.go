package protocol

import (
	"github.com/RoSpaceDev/libocr/internal/jmt"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

type KeyValuePair = jmt.KeyValue

type StateRootDigest = jmt.Digest

type AttestedReportMany[RI any] struct {
	ReportWithInfo       ocr3types.ReportWithInfo[RI]
	AttributedSignatures []types.AttributedOnchainSignature
}

type KeyValuePairWithDeletions struct {
	Key     []byte
	Value   []byte
	Deleted bool
}

type StateTransitionOutputs struct {
	WriteSet []KeyValuePairWithDeletions
}

type TreeSyncPhase int

const (
	// Tree sync was never started, or was completed. Regardless, it's not
	// happening right now.
	TreeSyncPhaseInactive TreeSyncPhase = iota
	// Tree sync is waiting for the necessary parts of the key-value store to be
	// cleaned up before it can start.
	TreeSyncPhaseWaiting
	// Tree sync is actively progressing now.
	TreeSyncPhaseActive
)

type TreeSyncStatus struct {
	Phase                  TreeSyncPhase
	TargetSeqNr            uint64
	TargetStateRootDigest  StateRootDigest
	PendingKeyDigestRanges PendingKeyDigestRanges
}
