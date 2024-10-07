package networking

import (
	"container/list"
	"fmt"
	"strings"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"go.uber.org/multierr"
)

func peerGroupStreamNamePrefix(configDigestPrefix ocr2types.ConfigDigestPrefix) (streamNamePrefix string, ok bool) {
	switch configDigestPrefix { // nolint:exhaustive
	case ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo:
		return "ccip-rmn/", true
	default:
		return "", false
	}
}

type peerGroupFactory struct {
	peer *concretePeerV2
}

type PeerGroupFactory interface {
	// This call is necessary for peer discovery to work among the group
	// configured in the parameters. Once the peer group is closed, peer
	// discovery will cease, and all streams created under the peer group
	// will be automatically closed. For pure bootstrapping, it is expected to
	// invoke NewPeerGroup and never create any streams using NewStream.
	NewPeerGroup(
		configDigest ocr2types.ConfigDigest,
		peerIDs []string,
		bootstrappers []commontypes.BootstrapperLocator,
	) (PeerGroup, error)
}

var _ PeerGroupFactory = &peerGroupFactory{}

type Stream interface {
	SendMessage(data []byte)
	ReceiveMessages() <-chan []byte
	Close() error
}

var _ Stream = &ragep2p.Stream{}

//sumtype:decl
type NewStreamArgs interface {
	isNewStreamArgs()
}

type NewStreamArgs1 struct {
	StreamName         string
	OutgoingBufferSize int // number of messages that fit in the outgoing buffer
	IncomingBufferSize int // number of messages that fit in the incoming buffer
	MaxMessageLength   int
	MessagesLimit      ragep2p.TokenBucketParams // rate limit for incoming messages
	BytesLimit         ragep2p.TokenBucketParams // rate limit for incoming messages
}

func (NewStreamArgs1) isNewStreamArgs() {}

type PeerGroup interface {
	// See ragep2p.Host NewStream for details on the parameters. The stream will
	// be automatically closed upon closure of the peer group. The
	// streamName must be prefixed as dictated by the config digest prefix used
	// in NewPeerGroup.
	NewStream(remotePeerID string, newStreamArgs NewStreamArgs) (Stream, error)

	// Close closes all opened streams, and stops peer discovery for the group.
	// Future calls to NewStream will error.
	Close() error
}

var _ PeerGroup = &peerGroup{}

func (f *peerGroupFactory) NewPeerGroup(
	configDigest ocr2types.ConfigDigest,
	peerIDs []string,
	bootstrappers []commontypes.BootstrapperLocator,
) (PeerGroup, error) {
	configDigestPrefix := ocr2types.ConfigDigestPrefixFromConfigDigest(configDigest)
	streamNamePrefix, ok := peerGroupStreamNamePrefix(configDigestPrefix)
	if !ok {
		return nil, fmt.Errorf("config digest prefix %s is not allowed", configDigestPrefix)
	}

	decodedv2PeerIDs, err := decodev2PeerIDs(peerIDs)
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

	decodedv2Bootstrappers, err := decodev2Bootstrappers(bootstrappers)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 bootstrappers: %w", err)
	}

	registration, err := f.peer.register(configDigest, decodedv2PeerIDs, decodedv2Bootstrappers)
	if err != nil {
		return nil, err
	}

	return &peerGroup{
		registration,
		f.peer.host,
		streamNamePrefix,
		peerIDSet,
		sync.Mutex{},
		list.New(),
		peerGroupStateOpen,
	}, nil
}

type peerGroupState uint8

const (
	_ peerGroupState = iota
	peerGroupStateOpen
	peerGroupStateClosed
)

type peerGroup struct {
	reg  *endpointRegistration
	host *ragep2p.Host

	streamNamePrefix string
	peerIDSet        map[ragetypes.PeerID]struct{}

	mu            sync.Mutex
	openedStreams *list.List
	state         peerGroupState
}

// managedStream is a wrapper around ragep2p.Stream that removes the stream from
// peerGroup upon Close.
type managedStream struct {
	stream  *ragep2p.Stream
	onClose func()
}

func (m *managedStream) Close() error {
	m.onClose()
	return m.stream.Close()
}

func (m *managedStream) ReceiveMessages() <-chan []byte {
	return m.stream.ReceiveMessages()
}

func (m *managedStream) SendMessage(data []byte) {
	m.stream.SendMessage(data)
}

var _ Stream = &managedStream{}

func (f *peerGroup) NewStream(remotePeerID string, newStreamArgs NewStreamArgs) (Stream, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.state != peerGroupStateOpen {
		return nil, fmt.Errorf("peer group has been closed")
	}

	switch args := newStreamArgs.(type) {
	case NewStreamArgs1:
		var other ragetypes.PeerID
		if err := other.UnmarshalText([]byte(remotePeerID)); err != nil {
			return nil, fmt.Errorf("failed to decode remote peer ID %q: %w", remotePeerID, err)
		}

		if !strings.HasPrefix(args.StreamName, f.streamNamePrefix) {
			return nil, fmt.Errorf("stream name does not have expected prefix %q", f.streamNamePrefix)
		}
		if _, ok := f.peerIDSet[other]; !ok {
			return nil, fmt.Errorf("peer ID %s is not in the set of peer IDs registered with this peer group", other)
		}

		stream, err := f.host.NewStream(
			other,
			args.StreamName,
			args.OutgoingBufferSize,
			args.IncomingBufferSize,
			args.MaxMessageLength,
			args.MessagesLimit,
			args.BytesLimit,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create stream: %w", err)
		}

		element := f.openedStreams.PushBack(stream)

		managedStream := managedStream{
			stream,
			func() {
				f.mu.Lock()
				defer f.mu.Unlock()
				f.openedStreams.Remove(element)
			},
		}
		return &managedStream, nil
	default:
		panic("unreachable")
	}
}

func (f *peerGroup) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = peerGroupStateClosed

	var err error

	for e := f.openedStreams.Front(); e != nil; e = e.Next() {
		if e.Value == nil {
			// defensive
			continue
		}
		stream, ok := e.Value.(*ragep2p.Stream)
		if !ok {
			// defensive
			continue
		}
		// we don't really expect the first Close of a stream to error out but
		// let's be defensive
		err = multierr.Append(err, stream.Close())
	}
	f.openedStreams.Init()

	return multierr.Append(err, f.reg.Close())
}
