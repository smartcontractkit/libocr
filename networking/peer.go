package networking

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/transport"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	"github.com/libp2p/go-tcp-transport"
	"github.com/smartcontractkit/libocr/networking/knockingtls"

	"github.com/smartcontractkit/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/libp2p/go-libp2p"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	mplex "github.com/libp2p/go-libp2p-mplex"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

var (
	_ types.BinaryNetworkEndpointFactory = &concretePeer{}
	_ types.BootstrapperFactory          = &concretePeer{}
)

const (
	dhtPrefix = "/cl_peer_discovery_dht"
)

type NetworkingStack uint8

const (
	_ NetworkingStack = iota
	NetworkingStackV1
	NetworkingStackV2
	NetworkingStackPreferV2
)

func (n NetworkingStack) needsv2() bool {
	return n == NetworkingStackV2 || n == NetworkingStackPreferV2
}

func (n NetworkingStack) needsv1() bool {
	return n == NetworkingStackV1 || n == NetworkingStackPreferV2
}

type DiscovererDatabase interface {
	// StoreAnnouncement has key-value-store semantics and stores a peerID (key) and an associated serialized
	//announcement (value).
	StoreAnnouncement(ctx context.Context, peerID string, ann []byte) error

	// ReadAnnouncements returns one serialized announcement (if available) for each of the peerIDs in the form of a map
	// keyed by each announcement's corresponding peer ID.
	ReadAnnouncements(ctx context.Context, peerIDs []string) (map[string][]byte, error)
}

// PeerConfig configures the peer. A peer can operate with the v1 or v2 or both networking stacks, depending on
// the NetworkingStack set. The options for each stack are clearly marked, those for v1 start with V1 and those for v2
// start with V2. Only the options for the desired stack(s) need to be set.
type PeerConfig struct {
	// NetworkingStack declares which network stack will be used: v1, v2 or both (prefer v2).
	NetworkingStack NetworkingStack
	PrivKey         p2pcrypto.PrivKey
	Logger          types.Logger

	V1ListenIP     net.IP
	V1ListenPort   uint16
	V1AnnounceIP   net.IP
	V1AnnouncePort uint16
	V1Peerstore    p2ppeerstore.Peerstore

	// This should be 0 most of times, but when needed (eg when counter is somehow rolled back)
	// users can bump this value to manually bump the counter.
	V1DHTAnnouncementCounterUserPrefix uint32

	// V2ListenAddresses contains the addresses the peer will listen to on the network in <host>:<port> form as
	// accepted by net.Listen, but host and port must be fully specified and cannot be empty.
	V2ListenAddresses []string

	// V2AnnounceAddresses contains the addresses the peer will advertise on the network in <host>:<port> form as
	// accepted by net.Dial. The addresses should be reachable by peers of interest.
	V2AnnounceAddresses []string

	// Every V2DeltaReconcile a Reconcile message is sent to every peer.
	V2DeltaReconcile time.Duration

	// Dial attempts will be at least V2DeltaDial apart.
	V2DeltaDial time.Duration

	V2DiscovererDatabase DiscovererDatabase

	EndpointConfig EndpointConfig
}

// concretePeer represents a libp2p peer with one peer ID listening on one port
type concretePeer struct {
	p2phost.Host
	tls            *knockingtls.KnockingTLSTransport
	gater          *connectionGater
	logger         loghelper.LoggerWithContext
	endpointConfig EndpointConfig
	registrants    map[types.ConfigDigest]struct{}
	registrantsMu  *sync.Mutex

	dhtAnnouncementCounterUserPrefix uint32

	// list of bandwidth limiters, one for each connection to a remote peer.
	bandwidthLimiters *knockingtls.Limiters
}

var _ types.BinaryNetworkEndpointFactory = (*concretePeer)(nil)
var _ types.BootstrapperFactory = (*concretePeer)(nil)

// registrant is an endpoint pinned to a particular config digest that is attached to this peer
// There may only be one registrant per config digest
type registrant interface {
	allower
	getConfigDigest() types.ConfigDigest
}

// NewPeer creates a new peer
func NewPeer(c PeerConfig) (*concretePeer, error) {
	if c.V1ListenPort == 0 {
		return nil, errors.New("NewPeer requires a non-zero listen port")
	}

	peerID, err := p2ppeer.IDFromPrivateKey(c.PrivKey)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting peer ID from private key")
	}

	listenAddr, err := makeMultiaddr(c.V1ListenIP, c.V1ListenPort)
	if err != nil {
		return nil, errors.Wrap(err, "could not make listen multiaddr")
	}

	logger := loghelper.MakeRootLoggerWithContext(c.Logger).MakeChild(types.LogFields{
		"id":         "Peer",
		"peerID":     peerID.Pretty(),
		"listenPort": c.V1ListenPort,
		"listenIP":   c.V1ListenIP.String(),
		"listenAddr": listenAddr.String(),
	})

	gater, err := newConnectionGater(logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not create gater")
	}

	bandwidthLimiters := knockingtls.NewLimiters(logger)

	tlsID := knockingtls.ID
	tls, err := knockingtls.NewKnockingTLS(logger, c.PrivKey, bandwidthLimiters)
	if err != nil {
		return nil, errors.Wrap(err, "could not create knocking tls")
	}

	addrsFactory, err := makeAddrsFactory(c.V1AnnounceIP, c.V1AnnouncePort)
	if err != nil {
		return nil, errors.Wrap(err, "could not make addrs factory")
	}

	// build a custom upgrader that overrides the default secure transport with knocking TLS
	transportCon := func(upgrader *tptu.Upgrader) transport.Transport {
		betterUpgrader := tptu.Upgrader{
			upgrader.PSK,
			tls,
			upgrader.Muxer,
			upgrader.ConnGater,
		}

		return tcp.NewTCPTransport(&betterUpgrader)
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrs(listenAddr),
		libp2p.Identity(c.PrivKey),
		libp2p.DisableRelay(),
		libp2p.Security(tlsID, tls),
		libp2p.ConnectionGater(gater),
		libp2p.Peerstore(c.V1Peerstore),
		libp2p.AddrsFactory(addrsFactory),
		libp2p.Transport(transportCon),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	logger.Info("Peer: libp2p host booted", nil)

	return &concretePeer{
		Host:                             basicHost,
		gater:                            gater,
		tls:                              tls,
		logger:                           logger,
		endpointConfig:                   c.EndpointConfig,
		registrants:                      make(map[types.ConfigDigest]struct{}),
		registrantsMu:                    &sync.Mutex{},
		dhtAnnouncementCounterUserPrefix: c.V1DHTAnnouncementCounterUserPrefix,
		bandwidthLimiters:                bandwidthLimiters,
	}, nil
}

// NewEndpoint returns a new ocrEndpoint
func (p *concretePeer) NewEndpoint(
	configDigest types.ConfigDigest,
	pids []string,
	v1bootstrappers []string,
	v2bootstrappers []types.BootstrapperLocator,
	failureThreshold int,
	// number of messages allowed to be consumed by the peer per second.
	tokenBucketRefillRate float64,
	// number of allowed requests in a burst.
	tokenBucketSize int,
) (types.BinaryNetworkEndpoint, error) {
	if failureThreshold <= 0 {
		return nil, errors.New("can't set F to 0 or smaller")
	}

	if len(v1bootstrappers) < 1 {
		return nil, errors.New("requires at least one bootstrapper")
	}
	peerIDs, err := decodePeerIDs(pids)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode peer IDs")
	}

	bnAddrs, err := decodeBootstrappers(v1bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode bootstrappers")
	}

	return newOCREndpoint(p.logger, configDigest, p, peerIDs, bnAddrs, p.endpointConfig,
		failureThreshold, tokenBucketRefillRate, tokenBucketSize)
}

func decodeBootstrappers(bootstrappers []string) (bnAddrs []p2ppeer.AddrInfo, err error) {
	bnMAddrs := make([]ma.Multiaddr, len(bootstrappers))
	for i, bNode := range bootstrappers {
		bnMAddr, err := ma.NewMultiaddr(bNode)
		if err != nil {
			return bnAddrs, errors.Wrapf(err, "could not decode peer address %s", bNode)
		}
		bnMAddrs[i] = bnMAddr
	}
	bnAddrs, err = p2ppeer.AddrInfosFromP2pAddrs(bnMAddrs...)
	if err != nil {
		return bnAddrs, errors.Wrap(err, "could not get addrinfos")
	}
	return
}

func decodePeerIDs(pids []string) ([]p2ppeer.ID, error) {
	peerIDs := make([]p2ppeer.ID, len(pids))
	for i, pid := range pids {
		peerID, err := p2ppeer.Decode(pid)
		if err != nil {
			return nil, errors.Wrapf(err, "error decoding peer ID: %s", pid)
		}
		peerIDs[i] = peerID
	}
	return peerIDs, nil
}

func (p *concretePeer) NewBootstrapper(
	configDigest types.ConfigDigest,
	pids []string,
	v1bootstrappers []string,
	v2bootstrappers []types.BootstrapperLocator,
	F int,
) (types.Bootstrapper, error) {
	if F <= 0 {
		return nil, errors.New("can't set F to zero or smaller")
	}
	peerIDs, err := decodePeerIDs(pids)
	if err != nil {
		return nil, err
	}

	bnAddrs, err := decodeBootstrappers(v1bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode bootstrappers")
	}

	return newBootstrapper(p.logger, configDigest, p, peerIDs, bnAddrs, F)
}

func (p *concretePeer) register(r registrant) error {
	configDigest := r.getConfigDigest()
	p.logger.Debug("Peer: registering protocol handler", types.LogFields{
		"configDigest": configDigest.Hex(),
	})

	p.registrantsMu.Lock()
	defer p.registrantsMu.Unlock()

	if _, ok := p.registrants[configDigest]; ok {
		return errors.Errorf("endpoint with config digest %s has already been registered", configDigest.Hex())
	}
	p.registrants[configDigest] = struct{}{}

	p.gater.add(r)

	p.tls.UpdateAllowlist(p.gater.allowlist())

	return nil
}

func (p *concretePeer) deregister(r registrant) error {
	configDigest := r.getConfigDigest()
	p.logger.Debug("Peer: deregistering protocol handler", types.LogFields{
		"ProtocolID": configDigest.Hex(),
	})

	p.registrantsMu.Lock()
	defer p.registrantsMu.Unlock()

	if _, ok := p.registrants[configDigest]; !ok {
		return errors.Errorf("endpoint with config digest %s is not currently registered", configDigest.Hex())
	}
	delete(p.registrants, configDigest)

	p.gater.remove(r)

	p.tls.UpdateAllowlist(p.gater.allowlist())

	return nil
}

func (p *concretePeer) PeerID() string {
	return p.ID().Pretty()
}
