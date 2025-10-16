package networking

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/RoSpaceDev/libocr/commontypes"
	ocr2types "github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew"
	ragetypes "github.com/RoSpaceDev/libocr/ragep2p/types"
	"github.com/RoSpaceDev/libocr/subprocesses"

	"github.com/RoSpaceDev/libocr/internal/loghelper"
)

var (
	_ ocr2types.BinaryNetworkEndpoint2 = &ocrEndpointV3{}
)

// ocrEndpointV3 represents a member of a particular feed oracle group
type ocrEndpointV3 struct {
	// configuration and settings
	defaultPriorityConfig ocr2types.BinaryNetworkEndpoint2Config
	lowPriorityConfig     ocr2types.BinaryNetworkEndpoint2Config
	peerMapping           map[commontypes.OracleID]ragetypes.PeerID
	host                  *ragep2pnew.Host
	configDigest          ocr2types.ConfigDigest
	ownOracleID           commontypes.OracleID

	// internal and state management
	chSendToSelf chan ocr2types.InboundBinaryMessageWithSender
	chClose      chan struct{}
	streams      map[commontypes.OracleID]priorityStreamGroup
	registration io.Closer
	state        ocrEndpointState

	stateMu sync.RWMutex
	subs    subprocesses.Subprocesses

	// recv is exposed to clients of this network endpoint
	recv chan ocr2types.InboundBinaryMessageWithSender

	logger loghelper.LoggerWithContext
}

type priorityStreamGroup struct {
	Low     ragep2pnew.Stream2
	Default ragep2pnew.Stream2
}

//nolint:unused
func newOCREndpointV3(
	logger loghelper.LoggerWithContext,
	configDigest ocr2types.ConfigDigest,
	peer *concretePeerV2,
	peerIDs []ragetypes.PeerID,
	v3bootstrappers []ragetypes.PeerInfo,
	defaultPriorityConfig ocr2types.BinaryNetworkEndpoint2Config,
	lowPriorityConfig ocr2types.BinaryNetworkEndpoint2Config,
	registration io.Closer,
) (*ocrEndpointV3, error) {
	peerMapping := make(map[commontypes.OracleID]ragetypes.PeerID)
	for i, peerID := range peerIDs {
		peerMapping[commontypes.OracleID(i)] = peerID
	}
	reversedPeerMapping := reverseMappingV2(peerMapping)
	ownOracleID, ok := reversedPeerMapping[peer.peerID]
	if !ok {
		return nil, fmt.Errorf("host peer ID %s is not present in given peerMapping", peer.PeerID())
	}

	chSendToSelf := make(chan ocr2types.InboundBinaryMessageWithSender, sendToSelfBufferSize)

	logger = logger.MakeChild(commontypes.LogFields{
		"configDigest": configDigest.Hex(),
		"oracleID":     ownOracleID,
		"id":           "OCREndpointV3",
	})

	logger.Info("OCREndpointV3: Initialized", commontypes.LogFields{
		"bootstrappers": v3bootstrappers,
		"oracles":       peerIDs,
	})

	if len(v3bootstrappers) == 0 {
		logger.Warn("OCREndpointV3: No bootstrappers were provided. Peer discovery might not work reliably for this instance.", nil)
	}

	host, ok := peer.host.RawWrappee().(*ragep2pnew.Host)
	if !ok {
		return nil, fmt.Errorf("host is not a wrapped ragep2pnew.Host. Please set the appropriate value for PeerConfig.EnableExperimentalRageP2P")
	}

	o := &ocrEndpointV3{
		defaultPriorityConfig,
		lowPriorityConfig,
		peerMapping,
		host,
		configDigest,
		ownOracleID,
		chSendToSelf,
		make(chan struct{}),
		make(map[commontypes.OracleID]priorityStreamGroup),
		registration,
		ocrEndpointUnstarted,
		sync.RWMutex{},
		subprocesses.Subprocesses{},
		make(chan ocr2types.InboundBinaryMessageWithSender),
		logger,
	}
	err := o.start()
	return o, err
}

// Start the ocrEndpointV3. Called once at the end of the initialization code.
func (o *ocrEndpointV3) start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			o.Close()
		}
	}()

	o.stateMu.Lock()
	defer o.stateMu.Unlock()

	if o.state != ocrEndpointUnstarted {
		return fmt.Errorf("cannot start ocrEndpointV3 that is not unstarted, state was: %d", o.state)
	}
	o.state = ocrEndpointStarted

	for oid, pid := range o.peerMapping {
		if oid == o.ownOracleID {
			continue
		}

		// Initialize the underlying streams, one stream per priority level.
		lowPriorityStream, err := o.host.NewStream2(
			pid,
			streamNameFromConfigDigestAndPriority(o.configDigest, ragep2pnew.StreamPriorityLow),
			ragep2pnew.StreamPriorityLow,
			ragep2pnew.Stream2Limits{o.lowPriorityConfig.OverrideOutgoingMessageBufferSize,
				o.lowPriorityConfig.OverrideIncomingMessageBufferSize,
				o.lowPriorityConfig.MaxMessageLength,
				ragetypes.TokenBucketParams{
					o.lowPriorityConfig.MessagesRatePerOracle,
					uint32(o.lowPriorityConfig.MessagesCapacityPerOracle),
				},
				ragetypes.TokenBucketParams{
					o.lowPriorityConfig.BytesRatePerOracle,
					uint32(o.lowPriorityConfig.BytesCapacityPerOracle),
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create (low priority) stream for oracle %v (peer id: %q): %w", oid, pid, err)
		}

		defaultPriorityStream, err := o.host.NewStream2(
			pid,
			streamNameFromConfigDigestAndPriority(o.configDigest, ragep2pnew.StreamPriorityDefault),
			ragep2pnew.StreamPriorityDefault,
			ragep2pnew.Stream2Limits{
				o.defaultPriorityConfig.OverrideOutgoingMessageBufferSize,
				o.defaultPriorityConfig.OverrideIncomingMessageBufferSize,
				o.defaultPriorityConfig.MaxMessageLength,
				ragetypes.TokenBucketParams{
					o.defaultPriorityConfig.MessagesRatePerOracle,
					uint32(o.defaultPriorityConfig.MessagesCapacityPerOracle),
				},
				ragetypes.TokenBucketParams{
					o.defaultPriorityConfig.BytesRatePerOracle,
					uint32(o.defaultPriorityConfig.BytesCapacityPerOracle),
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create (default priority) stream for oracle %v (peer id: %q): %w", oid, pid, err)
		}

		o.streams[oid] = priorityStreamGroup{lowPriorityStream, defaultPriorityStream}
	}

	for oid := range o.streams {
		o.subs.Go(func() {
			o.runRecv(oid)
		})
	}
	o.subs.Go(func() {
		o.runSendToSelf()
	})

	o.logger.Info("OCREndpointV3: Started listening", nil)
	succeeded = true
	return nil
}

// Receive runloop is per-remote
// This means that each remote gets its own buffered channel, so even if one
// remote goes mad and sends us thousands of messages, we don't drop any
// messages from good remotes
func (o *ocrEndpointV3) runRecv(oid commontypes.OracleID) {
	chRecv1 := o.streams[oid].Default.Receive()
	chRecv2 := o.streams[oid].Low.Receive()
	for {
		select {
		case msg := <-chRecv1:
			select {
			case o.recv <- ocr2types.InboundBinaryMessageWithSender{o.translateInboundMessage(msg, ocr2types.BinaryMessagePriorityDefault), oid}:
				continue
			case <-o.chClose:
				return
			}
		case msg := <-chRecv2:
			select {
			case o.recv <- ocr2types.InboundBinaryMessageWithSender{o.translateInboundMessage(msg, ocr2types.BinaryMessagePriorityLow), oid}:
				continue
			case <-o.chClose:
				return
			}
		case <-o.chClose:
			return
		}
	}
}

func (o *ocrEndpointV3) translateInboundMessage(inMsg ragep2pnew.InboundBinaryMessage, priority ocr2types.BinaryMessageOutboundPriority) ocr2types.InboundBinaryMessage {
	switch msg := inMsg.(type) {
	case ragep2pnew.InboundBinaryMessagePlain:
		return ocr2types.InboundBinaryMessagePlain{msg.Payload, priority}

	case ragep2pnew.InboundBinaryMessageRequest:
		return ocr2types.InboundBinaryMessageRequest{
			ocrEndpointV3RequestHandle{priority, msg.RequestHandle},
			msg.Payload,
			priority,
		}

	case ragep2pnew.InboundBinaryMessageResponse:
		return ocr2types.InboundBinaryMessageResponse{msg.Payload, priority}
	}
	panic("unknown type of ragep2pnew.InboundBinaryMessage")
}

func (o *ocrEndpointV3) translateOutboundMessage(outMsg ocr2types.OutboundBinaryMessage) (
	ragemsg ragep2pnew.OutboundBinaryMessage,
	priority ocr2types.BinaryMessageOutboundPriority,
) {
	switch msg := outMsg.(type) {
	case ocr2types.OutboundBinaryMessagePlain:
		ragemsg = ragep2pnew.OutboundBinaryMessagePlain{msg.Payload}
		priority = msg.Priority

	case ocr2types.OutboundBinaryMessageRequest:
		var rageresponsepolicy ragep2pnew.ResponsePolicy
		switch responsepolicy := msg.ResponsePolicy.(type) {
		case ocr2types.SingleUseSizedLimitedResponsePolicy:
			rageresponsepolicy = &ragep2pnew.SingleUseSizedLimitedResponsePolicy{
				responsepolicy.MaxSize,
				responsepolicy.ExpiryTimestamp,
			}
		}
		ragemsg = ragep2pnew.OutboundBinaryMessageRequest{rageresponsepolicy, msg.Payload}
		priority = msg.Priority

	case ocr2types.OutboundBinaryMessageResponse:
		requestHandle, ok := ocr2types.MustGetOutboundBinaryMessageResponseRequestHandle(msg).(ocrEndpointV3RequestHandle)
		if !ok {
			o.logger.Error(
				"dropping OutboundBinaryMessageResponse with requestHandle of unknown type",
				commontypes.LogFields{},
			)
			return
		}
		ragemsg = requestHandle.rageRequestHandle.MakeResponse(msg.Payload)
		priority = msg.Priority

	default:
		panic("unknown type of commontypes.OutboundBinaryMessage")
	}

	return ragemsg, priority
}

func (o *ocrEndpointV3) runSendToSelf() {
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

// Close should be called to clean up even if Start is never called.
func (o *ocrEndpointV3) Close() error {
	o.stateMu.Lock()
	defer o.stateMu.Unlock()
	if o.state != ocrEndpointStarted {
		return fmt.Errorf("cannot close ocrEndpointV3 that is not started, state was: %d", o.state)
	}
	o.state = ocrEndpointClosed

	o.logger.Debug("OCREndpointV3: Closing", nil)

	o.logger.Debug("OCREndpointV3: Closing streams", nil)
	close(o.chClose)
	o.subs.Wait()

	var allErrors error
	for oid, priorityGroupStream := range o.streams {
		if err := priorityGroupStream.Default.Close(); err != nil {
			allErrors = errors.Join(allErrors, fmt.Errorf("error while closing (default priority) stream with oracle %v: %w", oid, err))
		}
		if err := priorityGroupStream.Low.Close(); err != nil {
			allErrors = errors.Join(allErrors, fmt.Errorf("error while closing (low priority) stream with oracle %v: %w", oid, err))
		}
	}

	o.logger.Debug("OCREndpointV3: Deregister", nil)
	if err := o.registration.Close(); err != nil {
		allErrors = errors.Join(allErrors, fmt.Errorf("error closing OCREndpointV3: could not deregister: %w", err))
	}

	o.logger.Debug("OCREndpointV3: Closing o.recv", nil)
	close(o.recv)

	o.logger.Info("OCREndpointV3: Closed", nil)
	return allErrors
}

// SendTo sends a message to the given oracle.
// It makes a best effort delivery. If stream is unavailable for any
// reason, it will fill up to outgoingMessageBufferSize then drop
// old messages until the stream becomes available again.
//
// NOTE: If a stream connection is lost, the buffer will keep only the newest
// messages and drop older ones until the stream opens again.
func (o *ocrEndpointV3) SendTo(msg ocr2types.OutboundBinaryMessage, to commontypes.OracleID) {
	o.stateMu.RLock()
	state := o.state
	o.stateMu.RUnlock()
	if state != ocrEndpointStarted {
		o.logger.Error("Send on non-started ocrEndpointV3", commontypes.LogFields{"state": state})
		return
	}

	if to == o.ownOracleID {
		o.sendToSelf(msg)
		return
	}

	ragemsg, priority := o.translateOutboundMessage(msg)
	switch priority {
	case ocr2types.BinaryMessagePriorityDefault:
		o.streams[to].Default.Send(ragemsg)
	case ocr2types.BinaryMessagePriorityLow:
		o.streams[to].Low.Send(ragemsg)
	}
}

type ocrEndpointV3RequestHandle struct {
	priority          ocr2types.BinaryMessageOutboundPriority
	rageRequestHandle ragep2pnew.RequestHandle
}

func (rh ocrEndpointV3RequestHandle) MakeResponse(payload []byte) ocr2types.OutboundBinaryMessageResponse {
	return ocr2types.MustMakeOutboundBinaryMessageResponse(
		rh,
		payload,
		rh.priority,
	)
}

type ocrEndpointV3RequestHandleSelf struct {
	priority ocr2types.BinaryMessageOutboundPriority
}

func (rh ocrEndpointV3RequestHandleSelf) MakeResponse(payload []byte) ocr2types.OutboundBinaryMessageResponse {
	return ocr2types.MustMakeOutboundBinaryMessageResponse(
		rh,
		payload,
		rh.priority,
	)
}

func (o *ocrEndpointV3) sendToSelf(outboundMessage ocr2types.OutboundBinaryMessage) {
	var inboundMessage ocr2types.InboundBinaryMessage

	switch msg := outboundMessage.(type) {
	case ocr2types.OutboundBinaryMessagePlain:
		inboundMessage = ocr2types.InboundBinaryMessagePlain(msg)
	case ocr2types.OutboundBinaryMessageRequest:
		inboundMessage = ocr2types.InboundBinaryMessageRequest{
			ocrEndpointV3RequestHandleSelf{msg.Priority},
			msg.Payload,
			msg.Priority,
		}
	case ocr2types.OutboundBinaryMessageResponse:
		// TODO: We may want to reconsider how self forwarding works in case of requests/responses, because
		// with the updates to Stream2 and ResponseChecker extra checks are performed. Therefore, this "alternative"
		// code path may not behave fully equivalent to how request/responses are handled in the normal flow.
		inboundMessage = ocr2types.InboundBinaryMessageResponse{
			msg.Payload,
			msg.Priority,
		}
	}

	select {
	case o.chSendToSelf <- ocr2types.InboundBinaryMessageWithSender{inboundMessage, o.ownOracleID}:
	default:
		o.logger.Error("Send-to-self buffer is full, dropping message", commontypes.LogFields{
			"remoteOracleID": o.ownOracleID,
		})
	}
}

// Broadcast sends a msg to all oracles in the peer mapping
func (o *ocrEndpointV3) Broadcast(msg ocr2types.OutboundBinaryMessage) {
	var subs subprocesses.Subprocesses
	defer subs.Wait()
	for oracleID := range o.peerMapping {
		subs.Go(func() {
			o.SendTo(msg, oracleID)
		})
	}
}

// Receive gives the channel to receive messages
func (o *ocrEndpointV3) Receive() <-chan ocr2types.InboundBinaryMessageWithSender {
	return o.recv
}

func streamNameFromConfigDigestAndPriority(cd ocr2types.ConfigDigest, priority ragep2pnew.StreamPriority) string {
	switch priority {
	case ragep2pnew.StreamPriorityLow:
		return fmt.Sprintf("ocr/%s/priority=low", cd)
	case ragep2pnew.StreamPriorityDefault:
		return fmt.Sprintf("ocr/%s", cd)
	}
	panic("case implementation for ragep2pnew.StreamPriority missing")
}
