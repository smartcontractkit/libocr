package types

import (
	"context"

	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type DiscovererDatabase interface {
	// StoreAnnouncement has key-value-store semantics and stores a peerID (key) and an associated serialized
	//announcement (value).
	StoreAnnouncement(ctx context.Context, peerID string, ann []byte) error

	// ReadAnnouncements returns one serialized announcement (if available) for each of the peerIDs in the form of a map
	// keyed by each announcement's corresponding peer ID.
	ReadAnnouncements(ctx context.Context, peerIDs []string) (map[string][]byte, error)
}

type Peer interface {
	PeerID() string
	OCR1BinaryNetworkEndpointFactory() *ocr1types.BinaryNetworkEndpointFactory
	OCR1BootstrapperFactory() *ocr1types.BootstrapperFactory
	OCR2BinaryNetworkEndpointFactory() *ocr2types.BinaryNetworkEndpointFactory
	OCR2BootstrapperFactory() *ocr2types.BootstrapperFactory
	Close() error
}
