package networking

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/networking/internal/ocrendpointv3/responselimit"
	ocrendpointv3types "github.com/smartcontractkit/libocr/networking/internal/ocrendpointv3/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/libocr/subprocesses"

	"github.com/smartcontractkit/libocr/internal/loghelper"
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
	host                  *ragep2p.Host
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

	responseChecker *responselimit.ResponseChecker

	// recv is exposed to clients of this network endpoint
	recv chan ocr2types.InboundBinaryMessageWithSender

	logger loghelper.LoggerWithContext
}

type priorityStreamGroup struct {
	Low     *ragep2p.Stream
	Default *ragep2p.Stream
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

	o := &ocrEndpointV3{
		defaultPriorityConfig,
		lowPriorityConfig,
		peerMapping,
		peer.host,
		configDigest,
		ownOracleID,
		chSendToSelf,
		make(chan struct{}),
		make(map[commontypes.OracleID]priorityStreamGroup),
		registration,
		ocrEndpointUnstarted,
		sync.RWMutex{},
		subprocesses.Subprocesses{},
		responselimit.NewResponseChecker(),
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

		noLimitsMaxMessageLength := ragep2p.MaxMessageLength
		noLimitsTokenBucketParams := ragep2p.TokenBucketParams{
			math.MaxFloat64,
			math.MaxUint32,
		}

		// Initialize the underlying streams, one stream per priority level.
		lowPriorityStream, err := o.host.NewStream(
			pid,
			streamNameFromConfigDigestAndPriority(o.configDigest, ocr2types.BinaryMessagePriorityLow),
			o.lowPriorityConfig.OverrideOutgoingMessageBufferSize,
			o.lowPriorityConfig.OverrideIncomingMessageBufferSize,
			noLimitsMaxMessageLength,
			noLimitsTokenBucketParams,
			noLimitsTokenBucketParams,
		)
		if err != nil {
			return fmt.Errorf("failed to create (low priority) stream for oracle %v (peer id: %q): %w", oid, pid, err)
		}

		defaultPriorityStream, err := o.host.NewStream(
			pid,
			streamNameFromConfigDigestAndPriority(o.configDigest, ocr2types.BinaryMessagePriorityDefault),
			o.defaultPriorityConfig.OverrideOutgoingMessageBufferSize,
			o.defaultPriorityConfig.OverrideIncomingMessageBufferSize,
			noLimitsMaxMessageLength,
			noLimitsTokenBucketParams,
			noLimitsTokenBucketParams,
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
	chRecv1 := o.streams[oid].Default.ReceiveMessages()
	chRecv2 := o.streams[oid].Low.ReceiveMessages()
	for {
		var (
			msg      []byte
			priority ocr2types.BinaryMessageOutboundPriority
		)

		select {
		case msg = <-chRecv1:
			priority = ocr2types.BinaryMessagePriorityDefault
		case msg = <-chRecv2:
			priority = ocr2types.BinaryMessagePriorityLow
		case <-o.chClose:
			return
		}

		inMsg, err := o.translateInboundMessage(msg, priority, oid)
		if err != nil {
			o.logger.Warn("Invalid inbound message", commontypes.LogFields{
				"remoteOracleID": oid,
				"priority":       priority,
				"reason":         err,
			})
			continue
		}
		select {
		case o.recv <- ocr2types.InboundBinaryMessageWithSender{inMsg, oid}:
			continue
		case <-o.chClose:
			return
		}
	}
}

type ocrEndpointV3PayloadType byte

const (
	_ ocrEndpointV3PayloadType = iota
	ocrEndpointV3PayloadTypePlain
	ocrEndpointV3PayloadTypeRequest
	ocrEndpointV3PayloadTypeResponse
)

type ocrEndpointV3Payload struct {
	sumType ocrEndpointV3PayloadSumType
}

func (o *ocrEndpointV3Payload) MarshalBinary() ([]byte, error) {
	var prefix byte
	switch o.sumType.(type) {
	case *ocrEndpointV3PayloadPlain:
		prefix = byte(ocrEndpointV3PayloadTypePlain)
	case *ocrEndpointV3PayloadRequest:
		prefix = byte(ocrEndpointV3PayloadTypeRequest)
	case *ocrEndpointV3PayloadResponse:
		prefix = byte(ocrEndpointV3PayloadTypeResponse)
	}
	sumTypeBytes, err := o.sumType.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append([]byte{prefix}, sumTypeBytes...), nil
}

func (o *ocrEndpointV3Payload) UnmarshalBinary(data []byte) error {
	if len(data) < 1 {
		return fmt.Errorf("data is too short to contain prefix")
	}
	prefix := ocrEndpointV3PayloadType(data[0])
	data = data[1:]
	switch prefix {
	case ocrEndpointV3PayloadTypePlain:
		o.sumType = &ocrEndpointV3PayloadPlain{}
	case ocrEndpointV3PayloadTypeRequest:
		o.sumType = &ocrEndpointV3PayloadRequest{}
	case ocrEndpointV3PayloadTypeResponse:
		o.sumType = &ocrEndpointV3PayloadResponse{}
	}
	return o.sumType.UnmarshalBinary(data)
}

//go-sumtype:decl ocrEndpointV3PayloadSumType

type ocrEndpointV3PayloadSumType interface {
	isOCREndpointV3PayloadSumType()
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type ocrEndpointV3PayloadPlain struct {
	payload []byte
}

func (op ocrEndpointV3PayloadPlain) isOCREndpointV3PayloadSumType() {}

func (op *ocrEndpointV3PayloadPlain) MarshalBinary() ([]byte, error) {
	return op.payload, nil
}

func (op *ocrEndpointV3PayloadPlain) UnmarshalBinary(data []byte) error {
	op.payload = data
	return nil
}

type ocrEndpointV3PayloadRequest struct {
	requestID ocrendpointv3types.RequestID
	payload   []byte
}

func (oreq ocrEndpointV3PayloadRequest) isOCREndpointV3PayloadSumType() {}

func (oreq *ocrEndpointV3PayloadRequest) MarshalBinary() ([]byte, error) {
	return append(oreq.requestID[:], oreq.payload...), nil
}

func (oreq *ocrEndpointV3PayloadRequest) UnmarshalBinary(data []byte) error {
	if len(data) < len(oreq.requestID) {
		return fmt.Errorf("data is too short to contain requestID")
	}
	oreq.requestID = ocrendpointv3types.RequestID(data[:len(oreq.requestID)])
	oreq.payload = data[len(oreq.requestID):]
	return nil
}

type ocrEndpointV3PayloadResponse struct {
	requestID ocrendpointv3types.RequestID
	payload   []byte
}

func (ores ocrEndpointV3PayloadResponse) isOCREndpointV3PayloadSumType() {}

func (ores *ocrEndpointV3PayloadResponse) MarshalBinary() ([]byte, error) {
	return append(ores.requestID[:], ores.payload...), nil
}

func (ores *ocrEndpointV3PayloadResponse) UnmarshalBinary(data []byte) error {
	if len(data) < len(ores.requestID) {
		return fmt.Errorf("data is too short to contain requestID")
	}
	ores.requestID = ocrendpointv3types.RequestID(data[:len(ores.requestID)])
	ores.payload = data[len(ores.requestID):]
	return nil
}

func (o *ocrEndpointV3) translateInboundMessage(ragepayload []byte, priority ocr2types.BinaryMessageOutboundPriority, from commontypes.OracleID) (ocr2types.InboundBinaryMessage, error) {
	var payload ocrEndpointV3Payload
	if err := payload.UnmarshalBinary(ragepayload); err != nil {
		return nil, err
	}

	switch msg := payload.sumType.(type) {
	case *ocrEndpointV3PayloadPlain:
		return ocr2types.InboundBinaryMessagePlain{msg.payload, priority}, nil

	case *ocrEndpointV3PayloadRequest:
		rid := msg.requestID

		return ocr2types.InboundBinaryMessageRequest{
			ocrEndpointV3RequestHandle{priority, rid},
			msg.payload,
			priority,
		}, nil

	case *ocrEndpointV3PayloadResponse:
		sid := ocrendpointv3types.StreamID{from, priority}
		rid := msg.requestID

		checkResult := o.responseChecker.CheckResponse(sid, rid, len(msg.payload))
		switch checkResult {
		case responselimit.ResponseCheckResultReject:
			return nil, fmt.Errorf("rejected response")
		case responselimit.ResponseCheckResultAllow:
			return ocr2types.InboundBinaryMessageResponse{msg.payload, priority}, nil
		}

		panic(fmt.Sprintf("unexpected responselimit.ResponseCheckResult: %#v", checkResult))
	}

	panic("unknown ocrEndpointV3PayloadType")
}

func (o *ocrEndpointV3) translateOutboundMessage(outMsg ocr2types.OutboundBinaryMessage, to commontypes.OracleID) (
	ragepayload []byte,
	priority ocr2types.BinaryMessageOutboundPriority,
	err error,
) {
	var payload ocrEndpointV3Payload
	switch msg := outMsg.(type) {
	case ocr2types.OutboundBinaryMessagePlain:
		payload.sumType = &ocrEndpointV3PayloadPlain{msg.Payload}
		priority = msg.Priority

	case ocr2types.OutboundBinaryMessageRequest:
		var ocrendpointv3responsepolicy responselimit.ResponsePolicy
		switch responsepolicy := msg.ResponsePolicy.(type) {
		case ocr2types.SingleUseSizedLimitedResponsePolicy:
			ocrendpointv3responsepolicy = &responselimit.SingleUseSizedLimitedResponsePolicy{
				responsepolicy.MaxSize,
				responsepolicy.ExpiryTimestamp,
			}
		}

		priority = msg.Priority

		sid := ocrendpointv3types.StreamID{to, priority}
		rid := ocrendpointv3types.GetRandomRequestID()
		o.responseChecker.SetPolicy(sid, rid, ocrendpointv3responsepolicy)

		payload.sumType = &ocrEndpointV3PayloadRequest{rid, msg.Payload}

	case ocr2types.OutboundBinaryMessageResponse:
		requestHandle, ok := ocr2types.MustGetOutboundBinaryMessageResponseRequestHandle(msg).(ocrEndpointV3RequestHandle)
		if !ok {
			o.logger.Error(
				"dropping OutboundBinaryMessageResponse with requestHandle of unknown type",
				commontypes.LogFields{},
			)
			return
		}

		requestID := requestHandle.requestID
		payload.sumType = &ocrEndpointV3PayloadResponse{requestID, msg.Payload}
		priority = msg.Priority

	default:
		panic("unknown type of ocr2types.OutboundBinaryMessage")
	}
	ragepayload, err = payload.MarshalBinary()
	return ragepayload, priority, err
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
		{
			sid := ocrendpointv3types.StreamID{oid, ocr2types.BinaryMessagePriorityDefault}
			o.responseChecker.ClearPoliciesForStream(sid)
			if err := priorityGroupStream.Default.Close(); err != nil {
				allErrors = errors.Join(allErrors, fmt.Errorf("error while closing (default priority) stream with oracle %v: %w", oid, err))
			}
		}
		{
			sid := ocrendpointv3types.StreamID{oid, ocr2types.BinaryMessagePriorityLow}
			o.responseChecker.ClearPoliciesForStream(sid)
			if err := priorityGroupStream.Low.Close(); err != nil {
				allErrors = errors.Join(allErrors, fmt.Errorf("error while closing (low priority) stream with oracle %v: %w", oid, err))
			}
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

	ragemsg, priority, err := o.translateOutboundMessage(msg, to)
	if err != nil {
		o.logger.Error("Failed to translate outbound message", commontypes.LogFields{
			"error": err,
		})
		return
	}
	switch priority {
	case ocr2types.BinaryMessagePriorityDefault:
		o.streams[to].Default.SendMessage(ragemsg)
	case ocr2types.BinaryMessagePriorityLow:
		o.streams[to].Low.SendMessage(ragemsg)
	}
}

type ocrEndpointV3RequestHandle struct {
	priority  ocr2types.BinaryMessageOutboundPriority
	requestID ocrendpointv3types.RequestID
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

func streamNameFromConfigDigestAndPriority(cd ocr2types.ConfigDigest, priority ocr2types.BinaryMessageOutboundPriority) string {
	switch priority {
	case ocr2types.BinaryMessagePriorityLow:
		return fmt.Sprintf("ocr/%s/priority=low", cd)
	case ocr2types.BinaryMessagePriorityDefault:
		return fmt.Sprintf("ocr/%s", cd)
	}
	panic("case implementation for ragep2p.StreamPriority missing")
}
