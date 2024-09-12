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

type genericStreamFactoryFactory struct {
	peer *concretePeerV2

	configDigestPrefix ocr2types.ConfigDigestPrefix
	streamNamePrefix   string
}

func newGenericStreamFactoryFactory(peer *concretePeerV2, configDigestPrefix ocr2types.ConfigDigestPrefix) (*genericStreamFactoryFactory, error) {
	streamNamePrefix, ok := correspondingStreamNamePrefix(configDigestPrefix)
	if !ok {
		return nil, fmt.Errorf("config digest prefix %s is not allowed", configDigestPrefix)
	}
	return &genericStreamFactoryFactory{peer, configDigestPrefix, streamNamePrefix}, nil
}

// This call is necessary for peer discovery to work among the group configured
// in the parameters. Once the stream factory is closed, peer discovery will
// cease, and all streams created under the stream factory will be automatically
// closed. For pure bootstrapping, it is expected to invoke NewStreamFactory and
// never create any streams using NewStream.
func (f *genericStreamFactoryFactory) NewStreamFactory(
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
) (*streamFactory, error) {
	if !f.configDigestPrefix.IsPrefixOf(configDigest) {
		return nil, fmt.Errorf("config digest %s does not have supported prefix %q", configDigest, f.configDigestPrefix)
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

	return &streamFactory{
		registration,
		f.peer.host,
		f.streamNamePrefix,
		peerIDSet,
		sync.Mutex{},
		[]*ragep2p.Stream{},
		streamFactoryOpen,
	}, nil
}

type streamFactoryState int

const (
	_ streamFactoryState = iota
	streamFactoryOpen
	streamFactoryClosed
)

type streamFactory struct {
	reg  *endpointRegistration
	host *ragep2p.Host

	streamNamePrefix string
	peerIDSet        map[ragetypes.PeerID]struct{}

	mu            sync.Mutex
	openedStreams []*ragep2p.Stream
	state         streamFactoryState
}

// See ragep2p.Host NewStream for details on the parameters. The stream will be
// automatically closed upon closure of the stream factory. The streamName must
// be prefixed as dictated by the config digest prefix used in NewStreamFactory.
func (f *streamFactory) NewStream(
	other ragetypes.PeerID,
	streamName string,
	outgoingBufferSize int, // number of messages that fit in the outgoing buffer
	incomingBufferSize int, // number of messages that fit in the incoming buffer
	maxMessageLength int,
	messagesLimit ragep2p.TokenBucketParams, // rate limit for incoming messages
	bytesLimit ragep2p.TokenBucketParams, // rate limit for incoming messages
) (*ragep2p.Stream, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.state != streamFactoryOpen {
		return nil, fmt.Errorf("stream factory has been closed")
	}

	if !strings.HasPrefix(streamName, f.streamNamePrefix) {
		return nil, fmt.Errorf("stream name does not have expected prefix %q", f.streamNamePrefix)
	}
	if _, ok := f.peerIDSet[other]; !ok {
		return nil, fmt.Errorf("peer ID %s is not in the set of peer IDs registered with this stream factory", other)
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

// Close closes all opened streams, and stops peer discovery for the group.
// Future calls to NewStream will error.
func (f *streamFactory) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = streamFactoryClosed
	for len(f.openedStreams) > 0 {
		head, tail := f.openedStreams[0], f.openedStreams[1:]
		f.openedStreams = tail

		head.Close()
	}
	return f.reg.Close()
}
