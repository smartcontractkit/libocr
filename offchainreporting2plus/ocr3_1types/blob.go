package ocr3_1types

import (
	"context"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
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
