package networking

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/internal/metricshelper"
	"github.com/RoSpaceDev/libocr/internal/peerkeyringhelper"
	"github.com/RoSpaceDev/libocr/networking/ragedisco"
	"github.com/RoSpaceDev/libocr/networking/ragep2pwrapper"
	"github.com/RoSpaceDev/libocr/networking/rageping"
	nettypes "github.com/RoSpaceDev/libocr/networking/types"
	ocr2types "github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/ragep2p"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew"
	ragetypes "github.com/RoSpaceDev/libocr/ragep2p/types"
	"github.com/prometheus/client_golang/prometheus"
)

var DangerDangerEnableExperimentalRageP2P = "I promise I know what I'm doing, give me the experimental ragep2p"

// Exactly one of PrivKey (deprecated) or PeerKeyring must be provided.
type PeerConfig struct {
	// Exactly one of PrivKey (deprecated) or PeerKeyring must be provided.
	//
	// Deprecated: Use PeerKeyring instead. This field is maintained for
	// backwards compatibility and will be removed in a future release.
	PrivKey ed25519.PrivateKey
	// Exactly one of PrivKey (deprecated) or PeerKeyring must be provided.
	PeerKeyring ragetypes.PeerKeyring

	Logger commontypes.Logger

	// V2ListenAddresses contains the addresses the peer will listen to on the network in <ip>:<port> form as
	// accepted by net.Listen.
	V2ListenAddresses []string

	// V2AnnounceAddresses contains the addresses the peer will advertise on the network in <ip>:<port> form as
	// accepted by net.Dial. The addresses should be reachable by peers of interest.
	// May be left unspecified, in which case the announce addresses are auto-detected based on V2ListenAddresses.
	V2AnnounceAddresses []string

	// Every V2DeltaReconcile a Reconcile message is sent to every peer.
	V2DeltaReconcile time.Duration

	// Dial attempts will be at least V2DeltaDial apart.
	V2DeltaDial time.Duration

	V2DiscovererDatabase nettypes.DiscovererDatabase

	V2EndpointConfig EndpointConfigV2

	MetricsRegisterer prometheus.Registerer

	LatencyMetricsServiceConfigs []*rageping.LatencyMetricsServiceConfig

	// Set this to DangerDangerEnableExperimentalRageP2P to use the experimental ragep2p stack.
	// If set to any other value, the default production ragep2p stack is used.
	// Note that the experimental ragep2p stack is not yet ready for production use, but should
	// in principle be backwards compatible with the current ragep2p stack.
	EnableExperimentalRageP2P string
}

func (c *PeerConfig) keyring() (ragetypes.PeerKeyring, error) {
	switch {
	case c.PeerKeyring != nil && c.PrivKey == nil:
		return c.PeerKeyring, nil
	case c.PrivKey != nil && c.PeerKeyring == nil:
		return peerkeyringhelper.NewPeerKeyringWithPrivateKey(c.PrivKey)
	default:
		return nil, fmt.Errorf("exactly one of PrivKey (deprecated) or PeerKeyring must be provided")
	}
}

// concretePeerV2 represents a ragep2p peer with one peer ID listening on one port
type concretePeerV2 struct {
	peerID                ragetypes.PeerID
	host                  ragep2pwrapper.Host
	discoverer            *ragedisco.Ragep2pDiscoverer
	metricsRegisterer     prometheus.Registerer
	logger                loghelper.LoggerWithContext
	endpointConfig        EndpointConfigV2
	latencyMetricsService rageping.LatencyMetricsService
}

// Users are expected to create (using the OCR*Factory() methods) and close endpoints and bootstrappers before calling
// Close() on the peer itself.
func NewPeer(c PeerConfig) (*concretePeerV2, error) {
	keyring, err := c.keyring()
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate keyring: %w", err)
	}

	peerID := ragetypes.PeerIDFromKeyring(keyring)

	logger := loghelper.MakeRootLoggerWithContext(c.Logger).MakeChild(commontypes.LogFields{
		"id":     "PeerV2",
		"peerID": peerID.String(),
	})

	if c.EnableExperimentalRageP2P == DangerDangerEnableExperimentalRageP2P {
		logger = logger.MakeChild(commontypes.LogFields{
			"ragep2p": "experimental",
		})
	}

	announceAddresses := c.V2AnnounceAddresses
	if len(c.V2AnnounceAddresses) == 0 {
		announceAddresses = c.V2ListenAddresses
	}

	metricsRegistererWrapper := metricshelper.NewPrometheusRegistererWrapper(c.MetricsRegisterer, c.Logger)

	discoverer := ragedisco.NewRagep2pDiscoverer(c.V2DeltaReconcile, announceAddresses, c.V2DiscovererDatabase, metricsRegistererWrapper)
	var host ragep2pwrapper.Host
	if c.EnableExperimentalRageP2P == DangerDangerEnableExperimentalRageP2P {
		h, err := ragep2pnew.NewHost(
			ragep2pnew.HostConfig{c.V2DeltaDial},
			keyring,
			c.V2ListenAddresses,
			discoverer,
			c.Logger,
			metricsRegistererWrapper,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to construct ragep2pnew host: %w", err)
		}
		host = ragep2pnew.Wrapped(h)
	} else {
		h, err := ragep2p.NewHost(
			ragep2p.HostConfig{c.V2DeltaDial},
			keyring,
			c.V2ListenAddresses,
			discoverer,
			c.Logger,
			metricsRegistererWrapper,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to construct ragep2pnew host: %w", err)
		}
		host = ragep2p.Wrapped(h)
	}
	err = host.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start ragep2p host: %w", err)
	}

	logger.Info("PeerV2: ragep2p host booted", nil)

	latencyMetricsService := rageping.NewLatencyMetricsService(
		host, metricsRegistererWrapper, logger, c.LatencyMetricsServiceConfigs,
	)

	return &concretePeerV2{
		peerID,
		host,
		discoverer,
		metricsRegistererWrapper,
		logger,
		c.V2EndpointConfig,
		latencyMetricsService,
	}, nil
}

// An endpointRegistration is held by an endpoint which services a particular configDigest. The invariant is that only
// there can be at most a single active (ie. not closed) endpointRegistration for some configDigest, and thus only at
// most one endpoint can service a particular configDigest at any given point in time. The endpoint is responsible for
// calling Close on the registration.
type endpointRegistration struct {
	deregisterFunc func() error
	once           sync.Once
}

func newEndpointRegistration(deregisterFunc func() error) *endpointRegistration {
	return &endpointRegistration{deregisterFunc, sync.Once{}}
}

func (r *endpointRegistration) Close() (err error) {
	r.once.Do(func() {
		err = r.deregisterFunc()
	})
	return err
}

func (p2 *concretePeerV2) register(configDigest ocr2types.ConfigDigest, oracles []ragetypes.PeerID, bootstrappers []ragetypes.PeerInfo) (*endpointRegistration, error) {
	if err := p2.discoverer.AddGroup(configDigest, oracles, bootstrappers); err != nil {
		p2.logger.Warn("PeerV2: Failed to register endpoint", commontypes.LogFields{"configDigest": configDigest})
		return nil, err
	}

	bootstrappersIDs := make([]ragetypes.PeerID, 0, len(bootstrappers))
	for _, b := range bootstrappers {
		bootstrappersIDs = append(bootstrappersIDs, b.ID)
	}

	p2.latencyMetricsService.RegisterPeers(oracles)
	p2.latencyMetricsService.RegisterPeers(bootstrappersIDs)

	return newEndpointRegistration(func() error {
		// Discoverer will not be closed until concretePeerV2.Close() is called.
		// By the time concretePeerV2.Close() is called all endpoints/bootstrappers should have already been closed.
		// Even if this weren't true, RemoveGroup() is a no-op if the discoverer is closed.

		p2.latencyMetricsService.UnregisterPeers(oracles)
		p2.latencyMetricsService.UnregisterPeers(bootstrappersIDs)

		return p2.discoverer.RemoveGroup(configDigest)
	}), nil
}

func (p2 *concretePeerV2) PeerID() string {
	return p2.peerID.String()
}

func (p2 *concretePeerV2) Close() error {
	p2.latencyMetricsService.Close()
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
			return nil, fmt.Errorf("failed to unmarshal v2 peer ID (%q) from BootstrapperLocator: %w", b.PeerID, err)
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
			return nil, fmt.Errorf("error decoding v2 peer ID (%q): %w", pid, err)
		}
		peerIDs[i] = rid
	}
	return peerIDs, nil
}

func (p2 *concretePeerV2) newEndpoint(
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	limits BinaryNetworkEndpointLimits,
) (commontypes.BinaryNetworkEndpoint, error) {
	decodedv2PeerIDs, err := decodev2PeerIDs(v2peerIDs)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 peer IDs: %w", err)
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 bootstrappers: %w", err)
	}

	registration, err := p2.register(configDigest, decodedv2PeerIDs, decodedv2Bootstrappers)
	if err != nil {
		return nil, err
	}

	endpoint, err := newOCREndpointV2(
		p2.logger,
		configDigest,
		p2,
		decodedv2PeerIDs,
		decodedv2Bootstrappers,
		EndpointConfigV2{
			p2.endpointConfig.IncomingMessageBufferSize,
			p2.endpointConfig.OutgoingMessageBufferSize,
		},
		limits,
		registration,
	)
	if err != nil {
		// Important: we close registration in case newOCREndpointV2 failed to prevent zombie registrations.
		return nil, errors.Join(err, registration.Close())
	}
	return endpoint, nil
}

func (p2 *concretePeerV2) newEndpoint3_1(
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	defaultPriorityConfig ocr2types.BinaryNetworkEndpoint2Config,
	lowPriorityConfig ocr2types.BinaryNetworkEndpoint2Config,
) (ocr2types.BinaryNetworkEndpoint2, error) {
	decodedv2PeerIDs, err := decodev2PeerIDs(v2peerIDs)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 peer IDs: %w", err)
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 bootstrappers: %w", err)
	}

	registration, err := p2.register(configDigest, decodedv2PeerIDs, decodedv2Bootstrappers)
	if err != nil {
		return nil, err
	}

	endpoint, err := newOCREndpointV3(
		p2.logger,
		configDigest,
		p2,
		decodedv2PeerIDs,
		decodedv2Bootstrappers,
		defaultPriorityConfig,
		lowPriorityConfig,
		registration,
	)
	if err != nil {
		// Important: we close registration in case newOCREndpointV2 failed to prevent zombie registrations.
		return nil, errors.Join(err, registration.Close())
	}
	return endpoint, nil
}

func (p2 *concretePeerV2) newBootstrapper(
	configDigest ocr2types.ConfigDigest,
	v2peerIDs []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
) (commontypes.Bootstrapper, error) {
	decodedv2PeerIDs, err := decodev2PeerIDs(v2peerIDs)
	if err != nil {
		return nil, err
	}

	decodedv2Bootstrappers, err := decodev2Bootstrappers(v2bootstrappers)
	if err != nil {
		return nil, fmt.Errorf("could not decode v2 bootstrappers: %w", err)
	}

	registration, err := p2.register(configDigest, decodedv2PeerIDs, decodedv2Bootstrappers)
	if err != nil {
		return nil, err
	}

	bootstrapper, err := newBootstrapperV2(p2.logger, configDigest, decodedv2PeerIDs, decodedv2Bootstrappers, registration)
	if err != nil {
		// Important: we close registration in case newBootstrapperV2 failed to prevent zombie registrations.
		return nil, errors.Join(err, registration.Close())
	}
	return bootstrapper, nil
}

func (p2 *concretePeerV2) OCR1BinaryNetworkEndpointFactory() *ocr1BinaryNetworkEndpointFactory {
	return &ocr1BinaryNetworkEndpointFactory{p2}
}

func (p2 *concretePeerV2) OCR2BinaryNetworkEndpointFactory() *ocr2BinaryNetworkEndpointFactory {
	return &ocr2BinaryNetworkEndpointFactory{p2}
}

func (p2 *concretePeerV2) OCR3_1BinaryNetworkEndpointFactory() *ocr3_1BinaryNetworkEndpointFactory {
	return &ocr3_1BinaryNetworkEndpointFactory{p2}
}

func (p2 *concretePeerV2) OCR1BootstrapperFactory() *ocr1BootstrapperFactory {
	return &ocr1BootstrapperFactory{p2}
}

func (p2 *concretePeerV2) OCR2BootstrapperFactory() *ocr2BootstrapperFactory {
	return &ocr2BootstrapperFactory{p2}
}

func (p2 *concretePeerV2) PeerGroupFactory() *peerGroupFactory {
	return &peerGroupFactory{p2}
}
