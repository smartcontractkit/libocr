package networking

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/internal/loghelper"
)

var (
	_ commontypes.BinaryNetworkEndpoint = &ocrEndpointV2{}
)

type EndpointConfigV2 struct {
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
}

// ocrEndpointV2 represents a member of a particular feed oracle group
type ocrEndpointV2 struct {
	// configuration and settings
	config              EndpointConfigV2
	peerIDs             []ragetypes.PeerID
	peerMapping         map[commontypes.OracleID]ragetypes.PeerID
	reversedPeerMapping map[ragetypes.PeerID]commontypes.OracleID
	peer                *concretePeer
	host                *ragep2p.Host
	configDigest        ocr2types.ConfigDigest
	bootstrappers       []ragetypes.PeerInfo
	failureThreshold    int
	ownOracleID         commontypes.OracleID

	// internal and state management
	chSendToSelf chan commontypes.BinaryMessageWithSender
	chClose      chan struct{}
	streams      map[commontypes.OracleID]*ragep2p.Stream
	state        ocrEndpointState
	stateMu      *sync.RWMutex
	wg           *sync.WaitGroup
	ctx          context.Context
	ctxCancel    context.CancelFunc

	// recv is exposed to clients of this network endpoint
	recv chan commontypes.BinaryMessageWithSender

	logger loghelper.LoggerWithContext

	limits BinaryNetworkEndpointLimits
}

func reverseMappingV2(m map[commontypes.OracleID]ragetypes.PeerID) map[ragetypes.PeerID]commontypes.OracleID {
	n := make(map[ragetypes.PeerID]commontypes.OracleID)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func newOCREndpointV2(
	logger loghelper.LoggerWithContext,
	configDigest ocr2types.ConfigDigest,
	peer *concretePeer,
	peerIDs []ragetypes.PeerID,
	v2bootstrappers []ragetypes.PeerInfo,
	config EndpointConfigV2,
	failureThreshold int,
	limits BinaryNetworkEndpointLimits,
) (*ocrEndpointV2, error) {
	peerMapping := make(map[commontypes.OracleID]ragetypes.PeerID)
	for i, peerID := range peerIDs {
		peerMapping[commontypes.OracleID(i)] = peerID
	}
	reversedPeerMapping := reverseMappingV2(peerMapping)
	ownOracleID, ok := reversedPeerMapping[peer.ragep2pHost.ID()]
	if !ok {
		return nil, errors.Errorf("host peer ID 0x%x is not present in given peerMapping", peer.ID())
	}

	chSendToSelf := make(chan commontypes.BinaryMessageWithSender, sendToSelfBufferSize)

	logger = logger.MakeChild(commontypes.LogFields{
		"configDigest": configDigest.Hex(),
		"oracleID":     ownOracleID,
		"id":           "OCREndpointV2",
	})

	ctx, cancel := context.WithCancel(context.Background())

	return &ocrEndpointV2{
		config,
		peerIDs,
		peerMapping,
		reversedPeerMapping,
		peer,
		peer.ragep2pHost,
		configDigest,
		v2bootstrappers,
		failureThreshold,
		ownOracleID,
		chSendToSelf,
		make(chan struct{}),
		make(map[commontypes.OracleID]*ragep2p.Stream),
		ocrEndpointUnstarted,
		new(sync.RWMutex),
		new(sync.WaitGroup),
		ctx,
		cancel,
		make(chan commontypes.BinaryMessageWithSender),
		logger,
		limits,
	}, nil
}

// Start the ocrEndpointV2. Should only be called once. Even in case of error Close() _should_ be called afterwards for cleanup.
func (o *ocrEndpointV2) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			o.Close()
		}
	}()

	o.stateMu.Lock()
	defer o.stateMu.Unlock()

	if o.state != ocrEndpointUnstarted {
		return fmt.Errorf("cannot start ocrEndpointV2 that is not unstarted, state was: %d", o.state)
	}
	o.state = ocrEndpointStarted

	if err := o.peer.registerV2(o); err != nil {
		return err
	}

	for oid, pid := range o.peerMapping {
		if oid == o.ownOracleID {
			continue
		}
		stream, err := o.host.NewStream(
			pid,
			fmt.Sprintf("ocr/%x", o.configDigest),
			o.config.OutgoingMessageBufferSize,
			o.config.IncomingMessageBufferSize,
			o.limits.MaxMessageLength,
			ragep2p.TokenBucketParams{
				o.limits.MessagesRatePerOracle,
				uint32(o.limits.MessagesCapacityPerOracle),
			},
			ragep2p.TokenBucketParams{
				o.limits.BytesRatePerOracle,
				uint32(o.limits.BytesCapacityPerOracle),
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to create stream for oracle %v (peer id: %s)", oid, pid)
		}
		o.streams[oid] = stream
	}

	o.wg.Add(len(o.streams))
	for oid := range o.streams {
		go o.runRecv(oid)
	}
	o.wg.Add(1)
	go o.runSendToSelf()

	o.logger.Info("OCREndpointV2: Started listening", nil)

	succeeded = true
	return nil
}

// Receive runloop is per-remote
// This means that each remote gets its own buffered channel, so even if one
// remote goes mad and sends us thousands of messages, we don't drop any
// messages from good remotes
func (o *ocrEndpointV2) runRecv(oid commontypes.OracleID) {
	defer o.wg.Done()
	chRecv := o.streams[oid].ReceiveMessages()
	for {
		select {
		case payload := <-chRecv:
			msg := commontypes.BinaryMessageWithSender{
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

func (o *ocrEndpointV2) runSendToSelf() {
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

func (o *ocrEndpointV2) Close() error {
	o.stateMu.Lock()
	if o.state != ocrEndpointStarted {
		defer o.stateMu.Unlock()
		return fmt.Errorf("cannot close ocrEndpointV2 that is not started, state was: %d", o.state)
	}
	o.state = ocrEndpointClosed
	o.stateMu.Unlock()

	o.logger.Debug("OCREndpointV2: Closing", nil)

	o.logger.Debug("OCREndpointV2: Closing streams", nil)
	close(o.chClose)
	o.ctxCancel()
	o.wg.Wait()

	var allErrors error
	for oid, stream := range o.streams {
		allErrors = multierr.Append(allErrors, errors.Wrapf(stream.Close(), "error while closing stream with oracle %v", oid))
	}

	o.logger.Debug("OCREndpointV2: Deregister", nil)
	allErrors = multierr.Append(allErrors, errors.Wrap(o.peer.deregisterV2(o), "error closing OCREndpointV2: could not deregister"))

	o.logger.Debug("OCREndpointV2: Closing o.recv", nil)
	close(o.recv)

	o.logger.Info("OCREndpointV2: Closed", nil)
	return allErrors
}

// SendTo sends a message to the given oracle
// It makes a best effort delivery. If stream is unavailable for any
// reason, it will fill up to outgoingMessageBufferSize then drop messages
// until the stream becomes available again
//
// NOTE: If a stream connection is lost, the buffer will keep only the newest
// messages and drop older ones until the stream opens again.
func (o *ocrEndpointV2) SendTo(payload []byte, to commontypes.OracleID) {
	o.stateMu.RLock()
	state := o.state
	o.stateMu.RUnlock()
	if state != ocrEndpointStarted {
		panic(fmt.Sprintf("send on non-running ocrEndpointV2, state was: %d", state))
	}

	if to == o.ownOracleID {
		o.sendToSelf(payload)
		return
	}

	o.streams[to].SendMessage(payload)
}

func (o *ocrEndpointV2) sendToSelf(payload []byte) {
	m := commontypes.BinaryMessageWithSender{
		Msg:    payload,
		Sender: o.ownOracleID,
	}

	select {
	case o.chSendToSelf <- m:
	default:
		o.logger.Error("Send-to-self buffer is full, dropping message", commontypes.LogFields{
			"remoteOracleID": o.ownOracleID,
		})
	}
}

// Broadcast sends a msg to all oracles in the peer mapping
func (o *ocrEndpointV2) Broadcast(payload []byte) {
	var wg sync.WaitGroup
	for oracleID := range o.peerMapping {
		wg.Add(1)
		go func(oid commontypes.OracleID) {
			o.SendTo(payload, oid)
			wg.Done()
		}(oracleID)
	}
	wg.Wait()
}

// Receive gives the channel to receive messages
func (o *ocrEndpointV2) Receive() <-chan commontypes.BinaryMessageWithSender {
	return o.recv
}

func (o *ocrEndpointV2) getConfigDigest() ocr2types.ConfigDigest {
	return o.configDigest
}

func (o *ocrEndpointV2) getV2Oracles() []ragetypes.PeerID {
	return o.peerIDs
}

func (o *ocrEndpointV2) getV2Bootstrappers() []ragetypes.PeerInfo {
	return o.bootstrappers
}
