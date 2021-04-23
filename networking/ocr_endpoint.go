package networking

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/peerstore"

	p2pnetwork "github.com/libp2p/go-libp2p-core/network"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2pprotocol "github.com/libp2p/go-libp2p-core/protocol"
	swarm "github.com/libp2p/go-libp2p-swarm"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/pkg/errors"
	dhtrouter "github.com/smartcontractkit/libocr/networking/dht-router"
	"github.com/smartcontractkit/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

var (
	_ types.BinaryNetworkEndpoint = &ocrEndpoint{}
)

type EndpointConfig struct {
	// IncomingMessageBufferSize is the per-remote number of incoming
	// messages to buffer. Any additional messages received on top of those
	// already in the queue will be dropped.
	IncomingMessageBufferSize int

	// OutgoingMessageBufferSize is the per-remote number of outgoing
	// messages to buffer. Any additional messages send on top of those
	// already in the queue will displace the oldest.
	// NOTE: OutgoingMessageBufferSize should be comfortably smaller than remote's
	// IncomingMessageBufferSize to give the remote enough space to process
	// them all in case we regained connection and now send a bunch at once
	OutgoingMessageBufferSize int

	// NewStreamTimeout is the maximum length of time to wait to open a
	// stream before we give up.
	// We shouldn't hit this in practice since libp2p will give up fast if
	// it can't get a connection, but it is here anyway as a failsafe.
	// Set to 0 to disable any timeout on top of what libp2p gives us by default.
	NewStreamTimeout time.Duration

	// DHTLookupInterval is the interval between which we do the expensive peer
	// lookup using DHT.
	//
	// Every DHTLookupInterval failures to open a stream to a peer, we will
	// attempt to lookup its IP from DHT
	DHTLookupInterval int

	// Interval at which nodes check connections to bootstrap nodes and reconnect if any of them is lost.
	// Setting this to a small value would allow newly joined bootstrap nodes to get more connectivity
	// more quickly, which helps to make bootstrap process faster. The cost of this operation is relatively
	// cheap. We set this to 1 minute during our test.
	BootstrapCheckInterval time.Duration
}

type ocrEndpointState int

// ocrEndpoint represents a member of a particular feed oracle group
type ocrEndpoint struct {
	// configuration and settings
	config              EndpointConfig
	peerMapping         map[types.OracleID]p2ppeer.ID
	reversedPeerMapping map[p2ppeer.ID]types.OracleID
	peerAllowlist       map[p2ppeer.ID]struct{}
	peer                *concretePeer
	rhost               *rhost.RoutedHost
	routing             dhtrouter.PeerDiscoveryRouter
	configDigest        types.ConfigDigest
	protocolID          p2pprotocol.ID
	bootstrapperAddrs   []p2ppeer.AddrInfo
	failureThreshold    int
	ownOracleID         types.OracleID

	// internal and state management
	chRecvs      map[types.OracleID](chan []byte)
	chSends      map[types.OracleID](chan []byte)
	muSends      map[types.OracleID]*sync.Mutex
	chSendToSelf chan types.BinaryMessageWithSender
	chClose      chan struct{}
	state        ocrEndpointState
	stateMu      *sync.RWMutex
	wg           *sync.WaitGroup
	ctx          context.Context
	ctxCancel    context.CancelFunc

	// recv is exposed to clients of this network endpoint
	recv chan types.BinaryMessageWithSender

	logger loghelper.LoggerWithContext

	// a map of rate limiters, one for each peer
	recvRateLimiters map[types.OracleID]limiter
}

const (
	ocrEndpointUnstarted = iota
	ocrEndpointStarted
	ocrEndpointClosed

	// sendToSelfBufferSize is how many messages we will keep in memory that
	// are sent to ourself before we start dropping
	sendToSelfBufferSize = 20

	protocolBaseName = "cl_offchainreporting"
	protocolVersion  = "1.0.0"
)

func newOCREndpoint(
	logger loghelper.LoggerWithContext,
	configDigest types.ConfigDigest,
	peer *concretePeer,
	peerIDs []p2ppeer.ID,
	bootstrappers []p2ppeer.AddrInfo,
	config EndpointConfig,

	failureThreshold int,
	tokenBucketRefillRate float64,
	tokenBucketSize int,
) (*ocrEndpoint, error) {
	peerMapping := make(map[types.OracleID]p2ppeer.ID)
	for i, peerID := range peerIDs {
		peerMapping[types.OracleID(i)] = peerID
	}
	reversedPeerMapping := reverseMapping(peerMapping)
	ownOracleID, ok := reversedPeerMapping[peer.ID()]
	if !ok {
		return nil, errors.Errorf("host peer ID 0x%x is not present in given peerMapping", peer.ID())
	}

	chRecvs := make(map[types.OracleID]chan []byte)
	chSends := make(map[types.OracleID]chan []byte)
	muSends := make(map[types.OracleID]*sync.Mutex)
	recvRateLimiters := make(map[types.OracleID]limiter)
	for oid := range peerMapping {
		if oid != ownOracleID {
			chRecvs[oid] = make(chan []byte, config.IncomingMessageBufferSize)
			chSends[oid] = make(chan []byte, config.OutgoingMessageBufferSize)
			muSends[oid] = new(sync.Mutex)
			recvRateLimiters[oid] = newLimiter(tokenBucketRefillRate, tokenBucketSize)
		}
	}

	chSendToSelf := make(chan types.BinaryMessageWithSender, sendToSelfBufferSize)

	protocolID := genProtocolID(configDigest)

	logger = logger.MakeChild(types.LogFields{
		"protocolID":   protocolID,
		"configDigest": configDigest.Hex(),
		"oracleID":     ownOracleID,
		"id":           "OCREndpoint",
	})

	ctx, cancel := context.WithCancel(context.Background())

	allowlist := make(map[p2ppeer.ID]struct{})
	for pid := range reversedPeerMapping {
		allowlist[pid] = struct{}{}
	}
	for _, b := range bootstrappers {
		allowlist[b.ID] = struct{}{}
	}

	return &ocrEndpoint{
		config,
		peerMapping,
		reversedPeerMapping,
		allowlist,
		peer,
		// Will be set in Start()
		nil,
		// Will be set in Start()
		nil,
		configDigest,
		protocolID,
		bootstrappers,
		failureThreshold,
		ownOracleID,
		chRecvs,
		chSends,
		muSends,
		chSendToSelf,
		make(chan struct{}),
		ocrEndpointUnstarted,
		new(sync.RWMutex),
		new(sync.WaitGroup),
		ctx,
		cancel,
		make(chan types.BinaryMessageWithSender),
		logger,
		recvRateLimiters,
	}, nil
}

func reverseMapping(m map[types.OracleID]p2ppeer.ID) map[p2ppeer.ID]types.OracleID {
	n := make(map[p2ppeer.ID]types.OracleID)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func genProtocolID(configDigest types.ConfigDigest) p2pprotocol.ID {
	// configDigest is namespaced under version but libp2p standard specifies a
	// trailing version, hence the dummy 1.0.0
	return p2pprotocol.ID(fmt.Sprintf("/%s/%s/%x/1.0.0", protocolBaseName, protocolVersion, configDigest))
}

// Start the ocrEndpoint. Should only be called once.
func (o *ocrEndpoint) Start() (err error) {
	o.stateMu.Lock()
	defer o.stateMu.Unlock()

	if o.state != ocrEndpointUnstarted {
		panic("ocrEndpoint has already been started")
	}
	o.state = ocrEndpointStarted

	if err := o.peer.register(o); err != nil {
		return err
	}

	if err := o.setupDHT(); err != nil {
		return errors.Wrap(err, "error setting up DHT")
	}

	o.rhost.SetStreamHandler(o.protocolID, o.streamReceiver)

	o.wg.Add(len(o.chRecvs))
	for oid := range o.chRecvs {
		go o.runRecv(oid)
	}
	o.wg.Add(len(o.chSends))
	for oid := range o.chSends {
		go o.runSend(oid)
	}
	o.wg.Add(1)
	go o.runSendToSelf()

	o.logger.Info("OCREndpoint: Started listening", nil)

	return nil
}

func (o *ocrEndpoint) setupDHT() (err error) {
	config := dhtrouter.BuildConfig(
		o.bootstrapperAddrs,
		dhtPrefix,
		o.configDigest,
		o.logger,
		o.config.BootstrapCheckInterval,
		o.failureThreshold,
		true, // XXX: set to true for debug. Should not do this in production.
		o.peer.dhtAnnouncementCounterUserPrefix,
	)

	acl := dhtrouter.NewPermitListACL(o.logger)

	acl.Activate(config.ProtocolID(), o.allowlist()...)
	aclHost := dhtrouter.WrapACL(o.peer, acl, o.logger)

	o.routing, err = dhtrouter.NewDHTRouter(
		o.ctx,
		config,
		aclHost,
	)
	if err != nil {
		return errors.Wrap(err, "could not initialize DHTRouter")
	}

	// Async
	o.routing.Start()

	o.rhost = rhost.Wrap(o.peer, o.routing)

	return nil
}

// Receive runloop is per-remote
// This means that each remote gets its own buffered channel, so even if one
// remote goes mad and sends us thousands of messages, we don't drop any
// messages from good remotes
func (o *ocrEndpoint) runRecv(oid types.OracleID) {
	defer o.wg.Done()
	var chRecv <-chan []byte = o.chRecvs[oid]
	for {
		select {
		case payload := <-chRecv:
			msg := types.BinaryMessageWithSender{
				Msg:    payload,
				Sender: oid,
			}
			select {
			case o.recv <- msg:
				continue
			case <-o.chClose:
				return
			}
		case <-o.chClose:
			return
		}
	}
}

func (o *ocrEndpoint) runSend(oid types.OracleID) {
	defer o.wg.Done()

	var chSend <-chan []byte = o.chSends[oid]
	destPeerID, err := o.oracleID2PeerID(oid)
	if err != nil {
		panic("error getting destination peer ID")
	}

	for {
		shouldRetry := o.sendOnStream(destPeerID, chSend)
		if !shouldRetry {
			return
		}
	}
}

func (o *ocrEndpoint) sendOnStream(destPeerID p2ppeer.ID, chSend <-chan []byte) (shouldRetry bool) {
	// Open a new stream to the destination peer
	var stream p2pnetwork.Stream

	nRetry := 0

	// Get a stream open by any means necessary, retry for as long as it takes
	for {
		var err error
		stream, err = func() (p2pnetwork.Stream, error) {
			var ctx context.Context
			if o.config.NewStreamTimeout == 0 {
				ctx = o.ctx
			} else {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(o.ctx, o.config.NewStreamTimeout)
				defer cancel()
			}
			return o.peer.NewStream(ctx, destPeerID, o.protocolID)
		}()

		if err == nil {
			break
		}

		// Exit early if the context was canceled because of Close
		if errors.Is(err, context.Canceled) {
			select {
			case <-o.chClose:
				return false
			default:
			}
		}

		// Libp2p automatically handles dial backoff for us, just ignore this error
		if !errors.Is(err, swarm.ErrDialBackoff) {
			o.logger.Debug("Peer unreachable", types.LogFields{
				"err":            err,
				"remotePeerID":   destPeerID,
				"nRetry":         nRetry,
				"remoteOracleId": o.reversedPeerMapping[destPeerID],
			})
		}

		// Fallback to the authenticated DHT peer discovery periodically if we can't connect right away
		if nRetry > 0 && nRetry%o.config.DHTLookupInterval == 0 {
			pAddr, err := o.routing.FindPeer(o.ctx, destPeerID)
			switch {
			case err == nil:
				o.logger.Debug("DHT lookup finished", types.LogFields{
					"result": pAddr,
					"nRetry": nRetry,
				})
				o.peer.Peerstore().AddAddrs(destPeerID, pAddr.Addrs, peerstore.TempAddrTTL)
			case errors.Is(err, context.Canceled):
				// Exit early if the context was canceled by the Close function
				return false
			default:
				o.logger.Warn("DHT lookup failed", types.LogFields{
					"err":            err,
					"remoteOracleId": o.reversedPeerMapping[destPeerID],
					"nRetry":         nRetry,
					"remotePeerID":   destPeerID,
				})
			}
		}

		nRetry++

		// Wait about 5 seconds before trying again
		// With some jitter to try and prevent simultaneous TCP dials
		// between hosts
		waitms := time.Duration(int64((4+rand.Float64()*2)*1000)) * time.Millisecond
		waitCh := time.After(waitms)

		select {
		case <-waitCh:
			// sleep here
		case <-o.chClose:
			return false
		}
	}

	defer stream.Reset()

	o.logger.Debug("Opened stream", types.LogFields{
		"remotePeerID": destPeerID,
	})

	for {
		select {
		case <-o.chClose:
			// All necessary cleanup has already been deferred by this point
			return false
		case payload := <-chSend:
			b := wireEncode(payload)
			_, err := stream.Write(b)
			if err != nil {
				// NOTE: We do the safest thing which is to exit.
				// This will close the stream for this write and restart this
				// function from the top.
				// Probably the connection got broken. No point in even trying
				// to resend the message.
				o.logger.Debug("Could not write to stream", types.LogFields{
					"err":          err,
					"remotePeerID": destPeerID,
				})

				return true
			}
		}
	}
}

func (o *ocrEndpoint) runSendToSelf() {
	defer o.wg.Done()
	for {
		select {
		case <-o.chClose:
			return
		case m := <-o.chSendToSelf:
			select {
			case o.recv <- m:
			case <-o.chClose:
				return
			}
		}
	}
}

func (o *ocrEndpoint) Close() error {
	o.stateMu.Lock()
	if o.state != ocrEndpointStarted {
		o.stateMu.Unlock()
		panic("cannot close ocrEndpoint that is not started")
	}
	o.state = ocrEndpointClosed
	o.stateMu.Unlock()

	o.logger.Debug("OCREndpoint: Closing", nil)

	o.logger.Debug("OCREndpoint: Removing stream handler", nil)
	o.peer.RemoveStreamHandler(o.protocolID)

	o.logger.Debug("OCREndpoint: Closing streams", nil)
	close(o.chClose)
	o.ctxCancel()
	o.wg.Wait()

	o.logger.Debug("OCREndpoint: Closing dht", nil)
	err := o.routing.Close()
	if err != nil {
		return errors.Wrap(err, "error closing OCREndpoint: could not close dht")
	}

	o.logger.Debug("OCREndpoint: Deregister", nil)
	if err := o.peer.deregister(o); err != nil {
		return errors.Wrap(err, "error closing OCREndpoint: could not deregister")
	}

	o.logger.Debug("OCREndpoint: Closing o.recv", nil)
	close(o.recv)

	o.logger.Info("OCREndpoint: Closed", nil)
	return nil
}

func (o *ocrEndpoint) streamReceiver(s p2pnetwork.Stream) {
	exit := make(chan struct{})
	defer close(exit)

	// Force stream reset on our side if close signal is received or if this function exits
	go func() {
		defer s.Reset()
		select {
		case <-o.chClose:
		case <-exit:
		}
	}()

	remotePeerID := s.Conn().RemotePeer()

	o.logger.Debug("Got incoming stream", types.LogFields{
		"remotePeerID":    remotePeerID,
		"remoteMultiaddr": s.Conn().RemoteMultiaddr(),
	})

	sender, err := o.peerID2OracleID(remotePeerID)
	if err != nil {
		o.logger.Error("Error getting sender", types.LogFields{
			"err":             err,
			"remotePeerID":    remotePeerID,
			"remoteMultiaddr": s.Conn().RemoteMultiaddr(),
		})
		return
	}
	r := bufio.NewReader(s)
	l := o.recvRateLimiters[sender]
	var countDropped uint64
	for {
		// Apply the rate limiter.
		isAllowed, err := isNextMessageAllowed(r, l)
		if err != nil {
			o.logger.Debug("Unable to peek at the next message from peer", types.LogFields{
				"err":             err,
				"remotePeerID":    remotePeerID,
				"remoteOracleID":  sender,
				"remoteMultiaddr": s.Conn().RemoteMultiaddr(),
			})
			return
		}
		if !isAllowed {
			countDropped += 1
			if isPowerOfTwo(countDropped) {
				o.logger.Info("Messages were dropped by the rate limiter", types.LogFields{
					"remotePeerID":         remotePeerID,
					"remoteOracleID":       sender,
					"remoteMultiaddr":      s.Conn().RemoteMultiaddr(),
					"messagesDroppedSoFar": countDropped,
				})
			}
			continue
		}
		if countDropped != 0 {
			o.logger.Info("Rate limiter is allowing messages to pass through again. Resetting dropped counter to zero.", types.LogFields{
				"remotePeerID":         remotePeerID,
				"remoteOracleID":       sender,
				"remoteMultiaddr":      s.Conn().RemoteMultiaddr(),
				"messagesDroppedSoFar": countDropped,
			})
			countDropped = 0
		}
		payload, err := readOneFromWire(r)
		if err != nil {
			o.logger.Debug("Lost connection to peer", types.LogFields{
				"err":             err,
				"remotePeerID":    remotePeerID,
				"remoteOracleID":  sender,
				"remoteMultiaddr": s.Conn().RemoteMultiaddr(),
			})
			// Safest thing to do on any error is to kill the stream and give up
			// A new one will automatically be opened next time we want to send a message
			return
		}

		chRecv := o.chRecvs[sender]
		select {
		case chRecv <- payload:
			continue
		default:
			o.logger.Warn("Incoming buffer is full, dropping message", types.LogFields{
				"remotePeerID":    remotePeerID,
				"remoteOracleID":  sender,
				"remoteMultiaddr": s.Conn().RemoteMultiaddr(),
			})
		}
	}
}

func (o *ocrEndpoint) peerID2OracleID(peerID p2ppeer.ID) (types.OracleID, error) {
	oracleID, ok := o.reversedPeerMapping[peerID]
	if !ok {
		return 0, errors.New("peer ID not found")
	}
	return oracleID, nil
}

func (o *ocrEndpoint) oracleID2PeerID(oracleID types.OracleID) (p2ppeer.ID, error) {
	peerID, ok := o.peerMapping[oracleID]
	if !ok {
		return "", errors.New("oracle ID not found")
	}
	return peerID, nil
}

func (o *ocrEndpoint) isStarted() bool {
	o.stateMu.RLock()
	defer o.stateMu.RUnlock()
	return o.state == ocrEndpointStarted
}

// SendTo sends a message to the given oracle
// It makes a best effort delivery. If stream is unavailable for any
// reason, it will fill up to outgoingMessageBufferSize then drop messages
// until the stream becomes available again
//
// NOTE: If a stream connection is lost, the buffer will keep only the newest
// messages and drop older ones until the stream opens again.
func (o *ocrEndpoint) SendTo(payload []byte, to types.OracleID) {
	if !o.isStarted() {
		panic("send on non-running ocrEndpoint")
	}

	if to == o.ownOracleID {
		o.sendToSelf(payload)
		return
	}

	var chSend chan []byte = o.chSends[to]

	// Must not allow concurrent sends on the same channel since it could cause
	// the simple ringbuffer below to block
	mu := o.muSends[to]
	mu.Lock()
	defer mu.Unlock()

	select {
	case chSend <- payload:
	default:
		select {
		case <-chSend:
			o.logger.Warn("Send buffer full, dropping oldest message", types.LogFields{
				"remoteOracleID": to,
			})
			chSend <- payload
		default:
			chSend <- payload
		}
	}
}

func (o *ocrEndpoint) sendToSelf(payload []byte) {
	m := types.BinaryMessageWithSender{
		Msg:    payload,
		Sender: o.ownOracleID,
	}

	select {
	case o.chSendToSelf <- m:
	default:
		o.logger.Error("Send-to-self buffer is full, dropping message", types.LogFields{
			"remoteOracleID": o.ownOracleID,
		})
	}
}

// Broadcast sends a msg to all oracles in the peer mapping
func (o *ocrEndpoint) Broadcast(payload []byte) {
	var wg sync.WaitGroup
	for oracleID := range o.peerMapping {
		wg.Add(1)
		go func(oid types.OracleID) {
			o.SendTo(payload, oid)
			wg.Done()
		}(oracleID)
	}
	wg.Wait()
}

// Receive gives the channel to receive messages
func (o *ocrEndpoint) Receive() <-chan types.BinaryMessageWithSender {
	return o.recv
}

// Conform to allower interface
func (o *ocrEndpoint) isAllowed(id p2ppeer.ID) bool {
	_, ok := o.peerAllowlist[id]
	return ok
}

// Conform to allower interface
func (o *ocrEndpoint) allowlist() (allowlist []p2ppeer.ID) {
	for k := range o.peerAllowlist {
		allowlist = append(allowlist, k)
	}
	return
}

func (o *ocrEndpoint) getConfigDigest() types.ConfigDigest {
	return o.configDigest
}

func isPowerOfTwo(num uint64) bool {
	return num != 0 && (num&(num-1)) == 0
}
