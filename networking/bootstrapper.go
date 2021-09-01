package networking

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	dhtrouter "github.com/smartcontractkit/libocr/networking/dht-router"
)

var (
	_ commontypes.Bootstrapper = &bootstrapper{}
)

type bootstrapper struct {
	networkingStack      NetworkingStack
	peer                 *concretePeer
	peerAllowlist        map[p2ppeer.ID]struct{}
	peerIDs              []p2ppeer.ID
	v2peerIDs            []ragetypes.PeerID // same as peerIDs but converted for v2
	v1bootstrappers      []p2ppeer.AddrInfo
	v2bootstrappers      []ragetypes.PeerInfo
	routing              dhtrouter.PeerDiscoveryRouter
	logger               loghelper.LoggerWithContext
	configDigest         ocr2types.ConfigDigest
	ctx                  context.Context
	ctxCancel            context.CancelFunc
	state                bootstrapperState
	stateMu              *sync.Mutex
	failureThreshold     int
	lowerBandwidthLimits func()
}

type bootstrapperState int

const (
	bootstrapperUnstarted = iota
	bootstrapperStarted
	bootstrapperClosed
	// Bandwidth rate limiter parameters for the bootstrap node.
	// Bootstrap nodes are contacted to fetch the mapping between peer IDs and peer IPs.
	// This bootstrapping is supposed to happen relatively rarely. Also, the full mapping is only a few KiB.
	bootstrapNodeTokenBucketRefillRate = 20 * 1024 // 20 KiB/s
	bootstrapNodeTokenBucketSize       = 50 * 1024 // 50 KiB/s
)

func newBootstrapper(
	networkingStack NetworkingStack,
	logger loghelper.LoggerWithContext,
	configDigest ocr2types.ConfigDigest,
	peer *concretePeer,
	v1peerIDs []p2ppeer.ID,
	v2peerIDs []ragetypes.PeerID,
	v1bootstrappers []p2ppeer.AddrInfo,
	v2bootstrappers []ragetypes.PeerInfo,
	F int,
) (*bootstrapper, error) {
	if !networkingStack.subsetOf(peer.networkingStack) {
		return nil, fmt.Errorf("newBootstrapper called with incompatible networking stack (peer has: %s, you want: %s)", peer.networkingStack, networkingStack)
	}
	lowerBandwidthLimits := func() {}
	if networkingStack.needsv1() {
		lowerBandwidthLimits = increaseBandwidthLimits(peer.bandwidthLimiters, v1peerIDs, v1bootstrappers,
			bootstrapNodeTokenBucketRefillRate, bootstrapNodeTokenBucketSize, logger)
	}

	allowlist := make(map[p2ppeer.ID]struct{})
	if networkingStack.needsv1() {
		for _, pid := range v1peerIDs {
			allowlist[pid] = struct{}{}
		}
		for _, b := range v1bootstrappers {
			allowlist[b.ID] = struct{}{}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger = logger.MakeChild(commontypes.LogFields{
		"id":           "OCREndpoint",
		"configDigest": configDigest.Hex(),
	})
	return &bootstrapper{
		networkingStack,
		peer,
		allowlist,
		v1peerIDs,
		v2peerIDs,
		v1bootstrappers,
		v2bootstrappers,
		nil,
		logger,
		configDigest,
		ctx,
		cancel,
		bootstrapperUnstarted,
		new(sync.Mutex),
		F,
		lowerBandwidthLimits,
	}, nil
}

// Start the bootstrapper. Should only be called once. Even in case of error Close() _should_ be called afterwards for cleanup.
func (b *bootstrapper) Start() error {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state != bootstrapperUnstarted {
		return fmt.Errorf("cannot start bootstrapper that is not unstarted, state was: %d", b.state)
	}

	b.state = bootstrapperStarted

	if b.networkingStack.needsv1() {
		if err := b.peer.registerV1(b); err != nil {
			return err
		}
	}
	if b.networkingStack.needsv2() {
		if err := b.peer.registerV2(b); err != nil {
			return err
		}
	}

	if b.networkingStack.needsv1() {
		if err := b.setupDHT(); err != nil {
			return errors.Wrap(err, "error setting up DHT")
		}
	}

	b.logger.Info("Bootstrapper: Started listening", nil)

	return nil
}

func (b *bootstrapper) setupDHT() (err error) {
	config := dhtrouter.BuildConfig(
		b.v1bootstrappers,
		dhtPrefix,
		b.configDigest,
		b.logger,
		b.peer.endpointConfig.BootstrapCheckInterval,
		b.failureThreshold,
		true,
		b.peer.dhtAnnouncementCounterUserPrefix,
	)

	acl := dhtrouter.NewPermitListACL(b.logger)

	acl.Activate(config.ProtocolID(), b.allowlist()...)
	aclHost := dhtrouter.WrapACL(b.peer.libp2pHost, acl, b.logger)

	b.routing, err = dhtrouter.NewDHTRouter(
		b.ctx,
		config,
		aclHost,
	)
	if err != nil {
		return errors.Wrap(err, "could not initialize DHTRouter")
	}

	// Async
	b.routing.Start()

	return nil
}

func (b *bootstrapper) Close() error {
	b.stateMu.Lock()
	if b.state != bootstrapperStarted {
		defer b.stateMu.Unlock()
		return fmt.Errorf("cannot close bootstrapper that is not started, state was: %d", b.state)
	}
	b.state = bootstrapperClosed
	b.stateMu.Unlock()

	var allErrors error
	if b.networkingStack.needsv1() {
		b.logger.Debug("Bootstrapper: lowering bandwidth limits when closing the bootstrap node", nil)
		b.lowerBandwidthLimits()

		allErrors = multierr.Append(allErrors, errors.Wrap(b.routing.Close(), "could not close dht router"))
	}

	b.ctxCancel()

	if b.networkingStack.needsv1() {
		allErrors = multierr.Append(allErrors, errors.Wrap(b.peer.deregisterV1(b), "could not unregister v1 bootstrapper"))
	}
	if b.networkingStack.needsv2() {
		allErrors = multierr.Append(allErrors, errors.Wrap(b.peer.deregisterV2(b), "could not unregister v2 bootstrapper"))
	}
	return allErrors
}

// Conform to allower interface
func (b *bootstrapper) isAllowed(id p2ppeer.ID) bool {
	_, ok := b.peerAllowlist[id]
	return ok
}

// Conform to allower interface
func (b *bootstrapper) allowlist() (allowlist []p2ppeer.ID) {
	for k := range b.peerAllowlist {
		allowlist = append(allowlist, k)
	}
	return
}

func (b *bootstrapper) getConfigDigest() ocr2types.ConfigDigest {
	return b.configDigest
}

func (b *bootstrapper) getV2Oracles() []ragetypes.PeerID {
	return b.v2peerIDs
}

func (b *bootstrapper) getV2Bootstrappers() []ragetypes.PeerInfo {
	return b.v2bootstrappers
}
