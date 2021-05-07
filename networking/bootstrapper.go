package networking

import (
	"context"
	"fmt"
	"sync"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	dhtrouter "github.com/smartcontractkit/libocr/networking/dht-router"
	"github.com/smartcontractkit/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

var (
	_ types.Bootstrapper = &bootstrapper{}
)

type bootstrapper struct {
	peer                 *concretePeer
	peerAllowlist        map[p2ppeer.ID]struct{}
	bootstrappers        []p2ppeer.AddrInfo
	routing              dhtrouter.PeerDiscoveryRouter
	logger               loghelper.LoggerWithContext
	configDigest         types.ConfigDigest
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

func newBootstrapper(logger loghelper.LoggerWithContext, configDigest types.ConfigDigest,
	peer *concretePeer, peerIDs []p2ppeer.ID, bootstrappers []p2ppeer.AddrInfo, F int) (*bootstrapper, error) {

	lowerBandwidthLimits := increaseBandwidthLimits(peer.bandwidthLimiters, peerIDs, bootstrappers,
		bootstrapNodeTokenBucketRefillRate, bootstrapNodeTokenBucketSize, logger)

	allowlist := make(map[p2ppeer.ID]struct{})
	for _, pid := range peerIDs {
		allowlist[pid] = struct{}{}
	}
	for _, b := range bootstrappers {
		allowlist[b.ID] = struct{}{}
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger = logger.MakeChild(types.LogFields{
		"id":           "OCREndpoint",
		"configDigest": configDigest.Hex(),
	})
	return &bootstrapper{
		peer,
		allowlist,
		bootstrappers,
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

func (b *bootstrapper) Start() error {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state != bootstrapperUnstarted {
		panic(fmt.Sprintf("cannot start bootstrapper that is not unstarted, state was: %d", b.state))
	}

	b.state = bootstrapperStarted

	if err := b.peer.register(b); err != nil {
		return err
	}

	if err := b.setupDHT(); err != nil {
		return errors.Wrap(err, "error setting up DHT")
	}

	b.logger.Info("Bootstrapper: Started listening", nil)

	return nil
}

func (b *bootstrapper) setupDHT() (err error) {
	config := dhtrouter.BuildConfig(
		b.bootstrappers,
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
	aclHost := dhtrouter.WrapACL(b.peer, acl, b.logger)

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
		panic(fmt.Sprintf("cannot close bootstrapper that is not started, state was: %d", b.state))
	}
	b.state = bootstrapperClosed
	b.stateMu.Unlock()

	b.logger.Debug("Bootstrapper: lowering bandwidth limits when closing the bootstrap node", nil)
	b.lowerBandwidthLimits()

	if err := b.routing.Close(); err != nil {
		return errors.Wrap(err, "could not close dht router")
	}

	b.ctxCancel()

	return errors.Wrap(b.peer.deregister(b), "could not unregister bootstrapper")
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

func (b *bootstrapper) getConfigDigest() types.ConfigDigest {
	return b.configDigest
}
