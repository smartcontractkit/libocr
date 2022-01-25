package networking

import (
	"crypto/ed25519"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/networking/ragedisco"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

// concretePeerV2 represents a ragep2p peer with one peer ID listening on one port
type concretePeerV2 struct {
	peerID         ragetypes.PeerID
	host           *ragep2p.Host
	discoverer     *ragedisco.Ragep2pDiscoverer
	registrantsMu  *sync.Mutex
	registrants    map[ocr2types.ConfigDigest]struct{}
	logger         loghelper.LoggerWithContext
	endpointConfig EndpointConfigV2
}

type registrantV2 interface {
	getConfigDigest() ocr2types.ConfigDigest
	getV2Bootstrappers() []ragetypes.PeerInfo
	getV2Oracles() []ragetypes.PeerID
}

func newPeerV2(c PeerConfig) (*concretePeerV2, error) {

	rawPriv, err := c.PrivKey.Raw()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get raw private key to use for v2")
	}
	ed25519Priv := ed25519.PrivateKey(rawPriv)
	if err := ed25519SanityCheck(ed25519Priv); err != nil {
		return nil, errors.Wrap(err, "ed25519 sanity check failed")
	}

	peerID, err := ragetypes.PeerIDFromPrivateKey(ed25519Priv)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting v2 peer ID from private key")
	}

	logger := loghelper.MakeRootLoggerWithContext(c.Logger).MakeChild(commontypes.LogFields{
		"id":     "PeerV2",
		"peerID": peerID.String(),
	})

	announceAddresses := c.V2AnnounceAddresses
	if len(c.V2AnnounceAddresses) == 0 {
		announceAddresses = c.V2ListenAddresses
	}
	discoverer := ragedisco.NewRagep2pDiscoverer(c.V2DeltaReconcile, announceAddresses, c.V2DiscovererDatabase)
	host, err := ragep2p.NewHost(
		ragep2p.HostConfig{c.V2DeltaDial},
		ed25519Priv,
		c.V2ListenAddresses,
		discoverer,
		c.Logger,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct ragep2p host")
	}
	err = host.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start ragep2p host")
	}

	logger.Info("PeerV2: ragep2p host booted", nil)

	return &concretePeerV2{
		peerID,
		host,
		discoverer,
		&sync.Mutex{},
		make(map[ocr2types.ConfigDigest]struct{}),
		logger,
		c.V2EndpointConfig,
	}, nil
}

func (p2 *concretePeerV2) register(r registrantV2) error {
	configDigest := r.getConfigDigest()
	p2.registrantsMu.Lock()
	defer p2.registrantsMu.Unlock()

	p2.logger.Debug("PeerV2: registering v2 protocol handler", commontypes.LogFields{
		"configDigest": configDigest.Hex(),
	})

	if _, ok := p2.registrants[configDigest]; ok {
		return fmt.Errorf("v2 endpoint with config digest %s has already been registered", configDigest.Hex())
	}
	p2.registrants[configDigest] = struct{}{}
	return p2.discoverer.AddGroup(
		r.getConfigDigest(),
		r.getV2Oracles(),
		r.getV2Bootstrappers(),
	)
}

func (p2 *concretePeerV2) deregister(r registrantV2) error {
	configDigest := r.getConfigDigest()
	p2.registrantsMu.Lock()
	defer p2.registrantsMu.Unlock()
	p2.logger.Debug("PeerV2: deregistering v2 protocol handler", commontypes.LogFields{
		"ProtocolID": configDigest.Hex(),
	})

	if _, ok := p2.registrants[configDigest]; !ok {
		return fmt.Errorf("v2 endpoint with config digest %s is not currently registered", configDigest.Hex())
	}
	delete(p2.registrants, configDigest)

	return p2.discoverer.RemoveGroup(configDigest)
}

func (p2 *concretePeerV2) PeerID() string {
	return p2.peerID.String()
}

func (p2 *concretePeerV2) Close() error {
	return p2.host.Close()
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

func (p2 *concretePeerV2) newEndpoint(
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
	limits BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	if f <= 0 {
		return nil, errors.New("can't set F to 0 or smaller")
	}

	if len(v2bootstrappers) < 1 {
		return nil, errors.New("requires at least one v2 bootstrapper")
	}

	decodedv2PeerIDs, err := decodev2PeerIDs(v2peerIDs)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v2 peer IDs")
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v2 bootstrappers")
	}

	return newOCREndpointV2(
		p2.logger,
		configDigest,
		p2,
		decodedv2PeerIDs,
		decodedv2Bootstrappers,
		EndpointConfigV2{
			p2.endpointConfig.IncomingMessageBufferSize,
			p2.endpointConfig.OutgoingMessageBufferSize,
		},
		f,
		limits,
	)
}

func (p2 *concretePeerV2) newBootstrapper(
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	f int,
) (commontypes.Bootstrapper, error) {
	if f <= 0 {
		return nil, errors.New("can't set f to zero or smaller")
	}

	decodedv2PeerIDs, err := decodev2PeerIDs(v2peerIDs)
	if err != nil {
		return nil, err
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode v2 bootstrappers")
	}

	return newBootstrapperV2(p2.logger, configDigest, p2, decodedv2PeerIDs, decodedv2Bootstrappers, f)
}
