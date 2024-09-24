package networking

import (
	"fmt"
	"strings"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

func correspondingStreamNamePrefix(configDigestPrefix ocr2types.ConfigDigestPrefix) (streamNamePrefix string, ok bool) {
	switch configDigestPrefix { // nolint:exhaustive
	case ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo:
		return "ccip-rmn/", true
	default:
		return "", false
	}
}

type genericNetworkEndpointFactory struct {
	peer *concretePeerV2
}

type GenericNetworkEndpointFactory interface {
	// This call is necessary for peer discovery to work among the group
	// configured in the parameters. Once the generic endpoint is closed, peer
	// discovery will cease, and all streams created under the generic endpoint
	// will be automatically closed. For pure bootstrapping, it is expected to
	// invoke NewStreamFactory and never create any streams using NewStream.
	NewGenericEndpoint(
		configDigest ocr2types.ConfigDigest,
		v2peerIDs []string,
		v2bootstrappers []commontypes.BootstrapperLocator,
	) (GenericNetworkEndpoint, error)
}

var _ GenericNetworkEndpointFactory = &genericNetworkEndpointFactory{}

type GenericNetworkEndpoint interface {
	// See ragep2p.Host NewStream for details on the parameters. The stream will
	// be automatically closed upon closure of the generic endpoint. The
	// streamName must be prefixed as dictated by the config digest prefix used
	// in NewStreamFactory.
	NewStream(
		other ragetypes.PeerID,
		streamName string,
		outgoingBufferSize int, // number of messages that fit in the outgoing buffer
		incomingBufferSize int, // number of messages that fit in the incoming buffer
		maxMessageLength int,
		messagesLimit ragep2p.TokenBucketParams, // rate limit for incoming messages
		bytesLimit ragep2p.TokenBucketParams, // rate limit for incoming messages
	) (*ragep2p.Stream, error)

	// Close closes all opened streams, and stops peer discovery for the group.
	// Future calls to NewStream will error.
	Close() error
}

var _ GenericNetworkEndpoint = &genericNetworkEndpoint{}

func (f *genericNetworkEndpointFactory) NewGenericEndpoint(
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
) (GenericNetworkEndpoint, error) {
	configDigestPrefix := ocr2types.ConfigDigestPrefixFromConfigDigest(configDigest)
	streamNamePrefix, ok := correspondingStreamNamePrefix(configDigestPrefix)
	if !ok {
		return nil, fmt.Errorf("config digest prefix %s is not allowed", configDigestPrefix)
	}

	decodedv2PeerIDs, err := decodev2PeerIDs(v2peerIDs)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 peer IDs: %w", err)
	}

	peerIDSet := make(map[ragetypes.PeerID]struct{}, len(decodedv2PeerIDs))
	for _, id := range decodedv2PeerIDs {
		if _, ok := peerIDSet[id]; ok {
			return nil, fmt.Errorf("duplicate v2 peer ID: %s", id)
		}
		peerIDSet[id] = struct{}{}
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 bootstrappers: %w", err)
	}

	registration, err := f.peer.register(configDigest, decodedv2PeerIDs, decodedv2Bootstrappers)
	if err != nil {
		return nil, err
	}

	return &genericNetworkEndpoint{
		registration,
		f.peer.host,
		streamNamePrefix,
		peerIDSet,
		sync.Mutex{},
		[]*ragep2p.Stream{},
		genericNetworkEndpointOpen,
	}, nil
}

type genericNetworkEndpointState int

const (
	_ genericNetworkEndpointState = iota
	genericNetworkEndpointOpen
	genericNetworkEndpointClosed
)

type genericNetworkEndpoint struct {
	reg  *endpointRegistration
	host *ragep2p.Host

	streamNamePrefix string
	peerIDSet        map[ragetypes.PeerID]struct{}

	mu            sync.Mutex
	openedStreams []*ragep2p.Stream
	state         genericNetworkEndpointState
}

func (f *genericNetworkEndpoint) NewStream(
	other ragetypes.PeerID,
	streamName string,
	outgoingBufferSize int,
	incomingBufferSize int,
	maxMessageLength int,
	messagesLimit ragep2p.TokenBucketParams,
	bytesLimit ragep2p.TokenBucketParams,
) (*ragep2p.Stream, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.state != genericNetworkEndpointOpen {
		return nil, fmt.Errorf("generic endpoint has been closed")
	}

	if !strings.HasPrefix(streamName, f.streamNamePrefix) {
		return nil, fmt.Errorf("stream name does not have expected prefix %q", f.streamNamePrefix)
	}
	if _, ok := f.peerIDSet[other]; !ok {
		return nil, fmt.Errorf("peer ID %s is not in the set of peer IDs registered with this generic endpoint", other)
	}

	stream, err := f.host.NewStream(
		other,
		streamName,
		outgoingBufferSize,
		incomingBufferSize,
		maxMessageLength,
		messagesLimit,
		bytesLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}
	f.openedStreams = append(f.openedStreams, stream)
	return stream, nil
}

func (f *genericNetworkEndpoint) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = genericNetworkEndpointClosed
	for len(f.openedStreams) > 0 {
		head, tail := f.openedStreams[0], f.openedStreams[1:]
		f.openedStreams = tail

		head.Close()
	}
	return f.reg.Close()
}
