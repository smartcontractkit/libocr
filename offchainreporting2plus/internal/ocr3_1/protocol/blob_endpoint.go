package protocol

import (
	"context"
	"fmt"
	"sync"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3_1types"
)

type BlobEndpointWrapper struct {
	mu      sync.Mutex
	wrapped *BlobEndpoint
}

func (bew *BlobEndpointWrapper) locked() *BlobEndpoint {
	bew.mu.Lock()
	wrapped := bew.wrapped
	bew.mu.Unlock()
	return wrapped
}

var _ ocr3_1types.BlobBroadcaster = &BlobEndpointWrapper{}

func (bew *BlobEndpointWrapper) BroadcastBlob(ctx context.Context, payload []byte, expirationHint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	wrapped := bew.locked()
	if wrapped == nil {
		return ocr3_1types.BlobHandle{}, errBlobEndpointUnavailable
	}
	return wrapped.BroadcastBlob(ctx, payload, expirationHint)
}

var _ ocr3_1types.BlobFetcher = &BlobEndpointWrapper{}

func (bew *BlobEndpointWrapper) FetchBlob(ctx context.Context, handle ocr3_1types.BlobHandle) ([]byte, error) {
	wrapped := bew.locked()
	if wrapped == nil {
		return nil, errBlobEndpointUnavailable
	}
	return wrapped.FetchBlob(ctx, handle)
}

func (bew *BlobEndpointWrapper) setBlobEndpoint(wrapped *BlobEndpoint) {
	bew.mu.Lock()
	bew.wrapped = wrapped
	bew.mu.Unlock()
}

type BlobEndpoint struct {
	ctx context.Context

	chBlobBroadcastRequest chan<- blobBroadcastRequest
	chBlobFetchRequest     chan<- blobFetchRequest
}

var errBlobEndpointUnavailable = fmt.Errorf("blob endpoint unavailable")

func expirySeqNr(expirationHint ocr3_1types.BlobExpirationHint) uint64 {
	switch beh := expirationHint.(type) {
	case ocr3_1types.BlobExpirationHintSequenceNumber:
		return beh.SeqNr
	default:
		panic(fmt.Sprintf("unexpected blob expiration hint type %T", beh))
	}
}

func (be *BlobEndpoint) BroadcastBlob(ctx context.Context, payload []byte, expirationHint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	chRequestDone := ctx.Done()
	chEndpointDone := be.ctx.Done()

	chResponse := make(chan blobBroadcastResponse)
	chDone := make(chan struct{})
	defer close(chDone)

	request := blobBroadcastRequest{
		payload,
		expirySeqNr(expirationHint),
		chResponse,
		chDone,
	}

	select {
	case be.chBlobBroadcastRequest <- request:
		select {
		case response := <-chResponse:
			if response.err != nil {
				return ocr3_1types.BlobHandle{}, response.err
			}
			return blobtypes.MakeBlobHandle(&response.cert), nil
		case <-chEndpointDone:
			return ocr3_1types.BlobHandle{}, be.ctx.Err()
		case <-chRequestDone:
			return ocr3_1types.BlobHandle{}, ctx.Err()
		}
	case <-chEndpointDone:
		return ocr3_1types.BlobHandle{}, be.ctx.Err()
	case <-chRequestDone:
		return ocr3_1types.BlobHandle{}, ctx.Err()
	}
}

var _ ocr3_1types.BlobBroadcaster = &BlobEndpoint{}

func (be *BlobEndpoint) FetchBlob(ctx context.Context, handle ocr3_1types.BlobHandle) ([]byte, error) {
	chRequestDone := ctx.Done()
	chEndpointDone := be.ctx.Done()

	chResponse := make(chan blobFetchResponse)
	chDone := make(chan struct{})
	defer close(chDone)

	blobHandleSumType := blobtypes.ExtractBlobHandleSumType(handle)
	if blobHandleSumType == nil {
		return nil, fmt.Errorf("zero value blob handle provided")
	}

	switch handle := blobHandleSumType.(type) {
	case *LightCertifiedBlob:
		if handle == nil {
			return nil, fmt.Errorf("zero value blob handle provided")
		}

		request := blobFetchRequest{
			*handle,
			chResponse,
			chDone,
		}

		select {
		case be.chBlobFetchRequest <- request:
			select {
			case response := <-chResponse:
				if response.err != nil {
					return nil, response.err
				}
				return response.payload, nil
			case <-chEndpointDone:
				return nil, be.ctx.Err()
			case <-chRequestDone:
				return nil, ctx.Err()
			}
		case <-chEndpointDone:
			return nil, be.ctx.Err()
		case <-chRequestDone:
			return nil, ctx.Err()
		}
	default:
		panic(fmt.Sprintf("unexpected blob handle type %T", handle))
	}
}

var _ ocr3_1types.BlobFetcher = &BlobEndpoint{}

// RoundBlobBroadcastFetcher is a thin wrapper around a blob broadcast fetcher
// which enforces that no expired blobs as of the current round at seqNr are
// fetched.
type RoundBlobBroadcastFetcher struct {
	seqNr                uint64
	blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher
}

func NewRoundBlobBroadcastFetcher(seqNr uint64, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) *RoundBlobBroadcastFetcher {
	return &RoundBlobBroadcastFetcher{seqNr, blobBroadcastFetcher}
}

var _ ocr3_1types.BlobBroadcastFetcher = &RoundBlobBroadcastFetcher{}

func (r *RoundBlobBroadcastFetcher) BroadcastBlob(ctx context.Context, payload []byte, expirationHint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	return r.blobBroadcastFetcher.BroadcastBlob(ctx, payload, expirationHint)
}

func (r *RoundBlobBroadcastFetcher) FetchBlob(ctx context.Context, handle ocr3_1types.BlobHandle) ([]byte, error) {
	blobHandleSumType := blobtypes.ExtractBlobHandleSumType(handle)
	switch cert := blobHandleSumType.(type) {
	case *blobtypes.LightCertifiedBlob:
		if cert != nil && cert.ExpirySeqNr < r.seqNr {
			return nil, fmt.Errorf("blob expired")
		}
		return r.blobBroadcastFetcher.FetchBlob(ctx, handle)
	default:
		panic(fmt.Sprintf("unexpected blob handle type %T", handle))
	}
}
