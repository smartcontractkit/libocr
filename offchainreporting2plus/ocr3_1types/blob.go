package ocr3_1types

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
)

type BlobHandle = blobtypes.BlobHandle

//go-sumtype:decl BlobExpirationHint

type BlobExpirationHint interface {
	isBlobExpirationHint()
}

var _ BlobExpirationHint = BlobExpirationHintSequenceNumber{}

type BlobExpirationHintSequenceNumber struct{ SeqNr uint64 }

func (BlobExpirationHintSequenceNumber) isBlobExpirationHint() {}

type BlobBroadcaster interface {
	BroadcastBlob(ctx context.Context, payload []byte, expirationHint BlobExpirationHint) (BlobHandle, error)
}

type BlobFetcher interface {
	FetchBlob(ctx context.Context, handle BlobHandle) ([]byte, error)
}

type BlobBroadcastFetcher interface {
	BlobBroadcaster
	BlobFetcher
}

// BlobHandleMarshalledBytesUpperBound returns a conservative upper bound on the
// marshalled length of a BlobHandle. The n & f parameters must match those in
// [ocr3types.ReportingPluginConfig].
//
// This is useful for computing limits like MaxObservationBytes when blob
// handles will be included in observations.
//
// Example usage in a ReportingPluginFactory:
//
//	func (f *Factory) NewReportingPlugin(
//	    ctx context.Context,
//	    config ocr3types.ReportingPluginConfig,
//	    bbf ocr3_1types.BlobBroadcastFetcher,
//	) (ocr3_1types.ReportingPlugin[RI], ocr3_1types.ReportingPluginInfo, error) {
//	    maxBlobHandleBytes := ocr3_1types.BlobHandleMarshalledBytesUpperBound(config.N, config.F)
//	    maxObservationBytes := baseObservationSize + numBlobHandles*maxBlobHandleBytes
//	    // ...
//	}
func BlobHandleMarshalledBytesUpperBound(n int, f int) int {
	return blobtypes.BlobHandleMarshalledBytesUpperBound(n, f)
}
