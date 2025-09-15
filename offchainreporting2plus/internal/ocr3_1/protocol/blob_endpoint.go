package protocol

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

type BlobEndpointWrapper struct {
	mu      sync.Mutex
	wrapped *BlobEndpoint
	limits  ocr3_1types.ReportingPluginLimits
}

func (bew *BlobEndpointWrapper) locked() (*BlobEndpoint, ocr3_1types.ReportingPluginLimits) {
	bew.mu.Lock()
	wrapped := bew.wrapped
	limits := bew.limits
	bew.mu.Unlock()
	return wrapped, limits
}

var _ ocr3_1types.BlobBroadcaster = &BlobEndpointWrapper{}

func (bew *BlobEndpointWrapper) BroadcastBlob(ctx context.Context, payload []byte, expirationHint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	wrapped, limits := bew.locked()
	if wrapped == nil {
		return ocr3_1types.BlobHandle{}, errBlobEndpointClosed
	}
	if len(payload) > limits.MaxBlobPayloadLength {
		return ocr3_1types.BlobHandle{}, fmt.Errorf("blob payload length %d exceeds maximum allowed length %d",
			len(payload), limits.MaxBlobPayloadLength)
	}
	return wrapped.BroadcastBlob(ctx, payload, expirationHint)
}

var _ ocr3_1types.BlobFetcher = &BlobEndpointWrapper{}

func (bew *BlobEndpointWrapper) FetchBlob(ctx context.Context, handle ocr3_1types.BlobHandle) ([]byte, error) {
	wrapped, _ := bew.locked()
	if wrapped == nil {
		return nil, errBlobEndpointClosed
	}
	return wrapped.FetchBlob(ctx, handle)
}

func (bew *BlobEndpointWrapper) setBlobEndpoint(wrapped *BlobEndpoint) {
	bew.mu.Lock()
	bew.wrapped = wrapped
	bew.mu.Unlock()
}

func (bew *BlobEndpointWrapper) SetLimits(limits ocr3_1types.ReportingPluginLimits) {
	bew.mu.Lock()
	bew.limits = limits
	bew.mu.Unlock()
}

type BlobEndpoint struct {
	ctx context.Context

	chBlobBroadcastRequest  chan<- blobBroadcastRequest
	chBlobBroadcastResponse <-chan blobBroadcastResponse

	chBlobFetchRequest  chan<- blobFetchRequest
	chBlobFetchResponse <-chan blobFetchResponse
}

var (
	errBlobEndpointClosed     = fmt.Errorf("blob endpoint closed")
	errReceivingChannelClosed = fmt.Errorf("receiving channel closed")
)

func expirySeqNr(expirationHint ocr3_1types.BlobExpirationHint) uint64 {
	switch beh := expirationHint.(type) {
	case ocr3_1types.BlobExpirationHintSequenceNumber:
		return beh.SeqNr
	default:
		panic(fmt.Sprintf("unexpected blob expiration hint type %T", beh))
	}
}

func (be *BlobEndpoint) BroadcastBlob(_ context.Context, payload []byte, expirationHint ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	chDone := be.ctx.Done()

	select {
	case be.chBlobBroadcastRequest <- blobBroadcastRequest{
		payload,
		expirySeqNr(expirationHint),
	}:
		response := <-be.chBlobBroadcastResponse
		if response.err != nil {
			return ocr3_1types.BlobHandle{}, response.err
		}

		select {
		case cert, ok := <-response.chCert:
			if !ok {
				return ocr3_1types.BlobHandle{}, errReceivingChannelClosed
			}
			return blobtypes.MakeBlobHandle(&cert), nil
		case <-chDone:
			return ocr3_1types.BlobHandle{}, errBlobEndpointClosed
		}
	case <-chDone:
		return ocr3_1types.BlobHandle{}, errBlobEndpointClosed
	}
}

var _ ocr3_1types.BlobBroadcaster = &BlobEndpoint{}

func (be *BlobEndpoint) FetchBlob(_ context.Context, handle ocr3_1types.BlobHandle) ([]byte, error) {
	chDone := be.ctx.Done()

	blobHandleSumType := blobtypes.ExtractBlobHandleSumType(handle)
	if blobHandleSumType == nil {
		return nil, fmt.Errorf("zero value blob handle provided")
	}
	switch handle := blobHandleSumType.(type) {
	case *LightCertifiedBlob:
		select {
		case be.chBlobFetchRequest <- blobFetchRequest{*handle}:
			response := <-be.chBlobFetchResponse
			if response.err != nil {
				return nil, response.err
			}
			select {
			case payload, ok := <-response.chPayload:
				if !ok {
					return nil, errReceivingChannelClosed
				}
				return payload, nil
			case <-chDone:
				return nil, errBlobEndpointClosed
			}
		case <-chDone:
			return nil, errBlobEndpointClosed
		}
	default:
		panic(fmt.Sprintf("unexpected blob handle type %T", handle))
	}
}

var _ ocr3_1types.BlobFetcher = &BlobEndpoint{}
