package networking

import (
	"context"
	"crypto/ed25519"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/transport"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	"github.com/libp2p/go-tcp-transport"
	"github.com/smartcontractkit/libocr/commontypes"
	inhousedisco "github.com/smartcontractkit/libocr/networking/inhouse-disco"
	"github.com/smartcontractkit/libocr/networking/knockingtls"

	"github.com/smartcontractkit/libocr/internal/loghelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/libp2p/go-libp2p"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	mplex "github.com/libp2p/go-libp2p-mplex"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

const (
	dhtPrefix = "/cl_peer_discovery_dht"
)

type GroupedDiscoverer interface {
	AddGroup(digest ocr2types.ConfigDigest, onodes []ragetypes.PeerID, bnodes []ragetypes.PeerInfo) error
	RemoveGroup(digest ocr2types.ConfigDigest) error
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
	Logger          commontypes.Logger

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
	libp2pHost   p2phost.Host
	libp2pPeerID p2ppeer.ID

	ragep2pHost       *ragep2p.Host
	ragep2pPeerID     ragetypes.PeerID
	ragep2pDiscoverer GroupedDiscoverer

	tls            *knockingtls.KnockingTLSTransport
	gater          *connectionGater
	logger         loghelper.LoggerWithContext
	endpointConfig EndpointConfig
	v1registrants  map[ocr2types.ConfigDigest]struct{}
	v2registrants  map[ocr2types.ConfigDigest]struct{}
	registrantsMu  *sync.Mutex

	dhtAnnouncementCounterUserPrefix uint32

	networkingStack NetworkingStack

	// list of bandwidth limiters, one for each connection to a remote peer.
	bandwidthLimiters *knockingtls.Limiters
}

// registrant is an endpoint pinned to a particular config digest that is attached to this peer
// There may only be one registrant per config digest
type registrant interface {
	getConfigDigest() ocr2types.ConfigDigest
}

type registrantV1 interface {
	registrant
	allower
}

type registrantV2 interface {
	registrant
	getV2Bootstrappers() []ragetypes.PeerInfo
	getV2Oracles() []ragetypes.PeerID
}

// NewPeer creates a new peer
func NewPeer(c PeerConfig) (*concretePeer, error) {

	rawPriv, err := c.PrivKey.Raw()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get raw private key to use for v2")
	}
	ed25519Priv := ed25519.PrivateKey(rawPriv)

	libp2pPeerID, err := p2ppeer.IDFromPrivateKey(c.PrivKey)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting v1 peer ID from private key")
	}

	ragep2pPeerID, err := ragetypes.PeerIDFromPrivateKey(ed25519Priv)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting v2 peer ID from private key")
	}

	logger := loghelper.MakeRootLoggerWithContext(c.Logger).MakeChild(commontypes.LogFields{
		"id":     "Peer",
		"peerID": libp2pPeerID.Pretty(),
	})

	gater, err := newConnectionGater(logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not create gater")
	}

	var libp2pHost p2phost.Host
	var tls *knockingtls.KnockingTLSTransport
	var bandwidthLimiters *knockingtls.Limiters
	if c.NetworkingStack.needsv1() {
		if c.V1ListenPort == 0 {
			return nil, errors.New("NewPeer requires a non-zero listen port")
		}

		listenAddr, err := makeMultiaddr(c.V1ListenIP, c.V1ListenPort)
		if err != nil {
			return nil, errors.Wrap(err, "could not make listen multiaddr")
		}
		logger = logger.MakeChild(commontypes.LogFields{
			"v1listenPort": c.V1ListenPort,
			"v1listenIP":   c.V1ListenIP.String(),
			"v1listenAddr": listenAddr.String(),
		})

		bandwidthLimiters = knockingtls.NewLimiters(logger)

		tlsID := knockingtls.ID
		tls, err = knockingtls.NewKnockingTLS(logger, c.PrivKey, bandwidthLimiters)
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

		libp2pHost, err = libp2p.New(context.Background(), opts...)
		if err != nil {
			return nil, err
		}

		logger.Info("Peer: libp2p host booted", nil)
	}

	var ragep2pHost *ragep2p.Host
	var gDiscoverer GroupedDiscoverer
	if c.NetworkingStack.needsv2() {
		discoverer := inhousedisco.NewRagep2pDiscoverer(c.V2DeltaReconcile, c.V2DiscovererDatabase)
		gDiscoverer = discoverer
		ragep2pHost, err = ragep2p.NewHost(
			ragep2p.HostConfig{c.V2DeltaDial},
			ed25519Priv,
			c.V2ListenAddresses,
			c.V2AnnounceAddresses,
			discoverer,
			c.Logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to construct ragep2p host")
		}
		err = ragep2pHost.Start()
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ragep2p host")
		}
	}

	return &concretePeer{
		libp2pHost:                       libp2pHost,
		libp2pPeerID:                     libp2pPeerID,
		ragep2pHost:                      ragep2pHost,
		ragep2pPeerID:                    ragep2pPeerID,
		ragep2pDiscoverer:                gDiscoverer,
		gater:                            gater,
		tls:                              tls,
		logger:                           logger,
		endpointConfig:                   c.EndpointConfig,
		v1registrants:                    make(map[ocr2types.ConfigDigest]struct{}),
		v2registrants:                    make(map[ocr2types.ConfigDigest]struct{}),
		registrantsMu:                    &sync.Mutex{},
		dhtAnnouncementCounterUserPrefix: c.V1DHTAnnouncementCounterUserPrefix,
		networkingStack:                  c.NetworkingStack,
		bandwidthLimiters:                bandwidthLimiters,
	}, nil
}

// newEndpoint returns an appropriate OCR endpoint depending on the networking stack used
func (p *concretePeer) newEndpoint(
	networkingStack NetworkingStack,
	configDigest ocr2types.ConfigDigest,
	pids []string,
	v1bootstrappers []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	failureThreshold int,
	limits BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	if failureThreshold <= 0 {
		return nil, errors.New("can't set F to 0 or smaller")
	}

	if networkingStack.needsv1() && len(v1bootstrappers) < 1 {
		return nil, errors.New("requires at least one v1 bootstrapper")
	}
	if networkingStack.needsv2() && len(v2bootstrappers) < 1 {
		return nil, errors.New("requires at least one v2 bootstrapper")
	}

	v1peerIDs, err := decodev1PeerIDs(pids)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v1 peer IDs")
	}

	v2peerIDs, err := decodev2PeerIDs(pids)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v2 peer IDs")
	}

	bnAddrs, err := decodev1Bootstrappers(v1bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v1 bootstrappers")
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v2 bootstrappers")
	}

	return newOCREndpoint(
		networkingStack,
		p.logger,
		configDigest,
		p,
		v1peerIDs,
		v2peerIDs,
		bnAddrs,
		decodedv2Bootstrappers,
		p.endpointConfig,
		failureThreshold,
		limits,
	)
}

func decodev1Bootstrappers(bootstrappers []string) (bnAddrs []p2ppeer.AddrInfo, err error) {
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

func decodev2Bootstrappers(v2bootstrappers []commontypes.BootstrapperLocator) (infos []ragetypes.PeerInfo, err error) {
	for _, b := range v2bootstrappers {
		addrs := make([]ragetypes.Address, len(b.Addrs))
		for i, a := range b.Addrs {
			addrs[i] = ragetypes.Address(a)
		}
		var rageID ragetypes.PeerID
		err := rageID.UnmarshalText([]byte(b.PeerID))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal v2 peer ID (%s) from BootstrapperLocator", b.PeerID)
		}
		infos = append(infos, ragetypes.PeerInfo{
			rageID,
			addrs,
		})
	}
	return
}

func decodev1PeerIDs(pids []string) ([]p2ppeer.ID, error) {
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

func decodev2PeerIDs(pids []string) ([]ragetypes.PeerID, error) {
	peerIDs := make([]ragetypes.PeerID, len(pids))
	for i, pid := range pids {
		var rid ragetypes.PeerID
		err := rid.UnmarshalText([]byte(pid))
		if err != nil {
			return nil, errors.Wrapf(err, "error decoding v2 peer ID: %s", pid)
		}
		peerIDs[i] = rid
	}
	return peerIDs, nil
}

func (p *concretePeer) newBootstrapper(
	networkingStack NetworkingStack,
	configDigest ocr2types.ConfigDigest,
	pids []string,
	v1bootstrappers []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	F int,
) (commontypes.Bootstrapper, error) {
	if F <= 0 {
		return nil, errors.New("can't set F to zero or smaller")
	}
	v1peerIDs, err := decodev1PeerIDs(pids)
	if err != nil {
		return nil, err
	}

	v2peerIDs, err := decodev2PeerIDs(pids)
	if err != nil {
		return nil, err
	}

	bnAddrs, err := decodev1Bootstrappers(v1bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v1 bootstrappers")
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v2 bootstrappers")
	}

	return newBootstrapper(networkingStack, p.logger, configDigest, p, v1peerIDs, v2peerIDs, bnAddrs, decodedv2Bootstrappers, F)
}

func (p *concretePeer) registerV1(r registrantV1) error {
	if !p.needsv1() {
		return nil
	}
	configDigest := r.getConfigDigest()
	p.registrantsMu.Lock()
	defer p.registrantsMu.Unlock()

	p.logger.Debug("Peer: registering v1 protocol handler", commontypes.LogFields{
		"configDigest": configDigest.Hex(),
	})

	if _, ok := p.v1registrants[configDigest]; ok {
		return errors.Errorf("v1 endpoint with config digest %s has already been registered", configDigest.Hex())
	}
	p.v1registrants[configDigest] = struct{}{}
	p.gater.add(r)
	p.tls.UpdateAllowlist(p.gater.allowlist())
	return nil
}

func (p *concretePeer) registerV2(r registrantV2) error {
	if !p.needsv2() {
		return nil
	}
	configDigest := r.getConfigDigest()
	p.registrantsMu.Lock()
	defer p.registrantsMu.Unlock()

	p.logger.Debug("Peer: registering v2 protocol handler", commontypes.LogFields{
		"configDigest": configDigest.Hex(),
	})

	if _, ok := p.v2registrants[configDigest]; ok {
		return errors.Errorf("v2 endpoint with config digest %s has already been registered", configDigest.Hex())
	}
	p.v2registrants[configDigest] = struct{}{}
	return p.ragep2pDiscoverer.AddGroup(
		r.getConfigDigest(),
		r.getV2Oracles(),
		r.getV2Bootstrappers(),
	)
}

func (p *concretePeer) deregisterV1(r registrantV1) error {
	if !p.needsv1() {
		return nil
	}
	configDigest := r.getConfigDigest()
	p.registrantsMu.Lock()
	defer p.registrantsMu.Unlock()
	p.logger.Debug("Peer: deregistering v1 protocol handler", commontypes.LogFields{
		"ProtocolID": configDigest.Hex(),
	})

	if _, ok := p.v1registrants[configDigest]; !ok {
		return errors.Errorf("v1 endpoint with config digest %s is not currently registered", configDigest.Hex())
	}
	delete(p.v1registrants, configDigest)

	p.gater.remove(r)
	p.tls.UpdateAllowlist(p.gater.allowlist())
	return nil
}

func (p *concretePeer) deregisterV2(r registrantV2) error {
	if !p.needsv2() {
		return nil
	}
	configDigest := r.getConfigDigest()
	p.registrantsMu.Lock()
	defer p.registrantsMu.Unlock()
	p.logger.Debug("Peer: deregistering v2 protocol handler", commontypes.LogFields{
		"ProtocolID": configDigest.Hex(),
	})

	if _, ok := p.v2registrants[configDigest]; !ok {
		return errors.Errorf("v2 endpoint with config digest %s is not currently registered", configDigest.Hex())
	}
	delete(p.v2registrants, configDigest)

	return p.ragep2pDiscoverer.RemoveGroup(configDigest)
}

func (p *concretePeer) needsv1() bool {
	return p.networkingStack.needsv1()
}

func (p *concretePeer) needsv2() bool {
	return p.networkingStack.needsv2()
}

func (p *concretePeer) PeerID() string {
	return p.ragep2pPeerID.String()
}

// backwards compatibility with libp2p provided method
func (p *concretePeer) ID() p2ppeer.ID {
	return p.libp2pPeerID
}

func (p *concretePeer) Close() error {
	if p.needsv2() {
		if err := p.ragep2pHost.Close(); err != nil {
			return err
		}
	}
	if p.needsv1() {
		return p.libp2pHost.Close()
	}
	return nil
}

func (c *concretePeer) OCRBinaryNetworkEndpointFactory() *ocrBinaryNetworkEndpointFactory {
	return &ocrBinaryNetworkEndpointFactory{c}
}

func (c *concretePeer) GenOCRBinaryNetworkEndpointFactory() *genocrBinaryNetworkEndpointFactory {
	return &genocrBinaryNetworkEndpointFactory{c}
}

func (c *concretePeer) OCRBootstrapperFactory() *ocrBootstrapperFactory {
	return &ocrBootstrapperFactory{c}
}

func (c *concretePeer) GenOCRBootstrapperFactory() *genocrBootstrapperFactory {
	return &genocrBootstrapperFactory{c}
}
