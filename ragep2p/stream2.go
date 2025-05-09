package ragep2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/internal/ringbuffer"
	"github.com/smartcontractkit/libocr/ragep2p/internal/responselimit"
	internaltypes "github.com/smartcontractkit/libocr/ragep2p/internal/types"
	"github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type StreamPriority byte

const (
	_ StreamPriority = iota
	StreamPriorityLow
	StreamPriorityDefault
)

// Stream2 is an over-the-network channel between two peers. Two peers may share
// multiple disjoint streams with different names. Streams are persistent and
// agnostic to the state of the connection. They abstract the underlying
// connection. Messages are delivered on a best effort basis.
type Stream2 struct {
	closedMu sync.Mutex
	closed   bool

	name     string
	other    types.PeerID
	streamID streamID

	outgoingBufferSize int
	maxMessageLength   int

	host *Host

	subprocesses subprocesses.Subprocesses
	ctx          context.Context
	cancel       context.CancelFunc
	logger       loghelper.LoggerWithContext
	chSend       chan OutboundBinaryMessage
	chReceive    chan InboundBinaryMessage

	chStreamToConn chan<- streamIDAndMessage
	demux          *demuxer
	chStreamOnOff  <-chan bool

	chStreamCloseRequest  chan<- peerStreamCloseRequest
	chStreamCloseResponse <-chan peerStreamCloseResponse
}

// Other returns the peer ID of the stream's counterparty.
func (st *Stream2) Other() types.PeerID {
	return st.other
}

// Name returns the name of the stream.
func (st *Stream2) Name() string {
	return st.name
}

// Best effort sending of messages. May fail without returning an error.
func (st *Stream2) Send(msg OutboundBinaryMessage) {
	var (
		ok            bool
		payloadLength int
	)
	switch msg := msg.(type) {
	case OutboundBinaryMessagePlain:
		ok = len(msg.Payload) <= st.maxMessageLength
		payloadLength = len(msg.Payload)
	case OutboundBinaryMessageRequest:
		ok = len(msg.Payload) <= st.maxMessageLength
		payloadLength = len(msg.Payload)
	case OutboundBinaryMessageResponse:
		// Response size is limited by the policy of the corresponding request
		// and may exceed the stream's default max message length.
		// Responses must never exceed the global ragep2p max message length.
		ok = len(msg.Payload) <= MaxMessageLength
		payloadLength = len(msg.Payload)
	default:
		panic(fmt.Sprintf("unknown OutboundBinaryMessage type: %T", msg))
	}

	if !ok {
		st.logger.Warn("dropping outbound message that is too large", commontypes.LogFields{
			"messagePayloadLength":    payloadLength,
			"streamMaxMessageLength":  st.maxMessageLength,
			"ragep2pMaxMessageLength": MaxMessageLength,
		})
		return
	}

	select {
	case st.chSend <- msg:
	case <-st.ctx.Done():
	}
}

// Best effort receiving of messages. The returned channel will be closed when
// the stream is closed. Note that this function may return the same channel
// across invocations.
func (st *Stream2) Receive() <-chan InboundBinaryMessage {
	return st.chReceive
}

// Close the stream. This closes any channel returned by ReceiveMessages earlier.
// After close the stream cannot be reopened. If the stream is needed in the
// future it should be created again through NewStream2.
// After close, any messages passed to SendMessage will be dropped.
func (st *Stream2) Close() error {
	st.closedMu.Lock()
	defer st.closedMu.Unlock()
	host := st.host

	if st.closed {
		return fmt.Errorf("already closed stream")
	}

	st.logger.Info("Stream winding down", nil)

	err := func() error {
		// Grab peersMu in case the peer has no streams left and we need to
		// delete it
		host.peersMu.Lock()
		defer host.peersMu.Unlock()

		select {
		case st.chStreamCloseRequest <- peerStreamCloseRequest{st.streamID}:
			resp := <-st.chStreamCloseResponse
			if resp.err != nil {
				st.logger.Error("Unexpected error during stream Close()", commontypes.LogFields{
					"error": resp.err,
				})
				return resp.err
			}
			if resp.peerHasNoStreams {
				st.logger.Trace("Garbage collecting peer", nil)
				peer := host.peers[st.other]
				host.subprocesses.Go(func() {
					peer.connLifeCycleMu.Lock()
					defer peer.connLifeCycleMu.Unlock()
					peer.connLifeCycle.connCancel()
					peer.connLifeCycle.connSubs.Wait()
				})
				delete(host.peers, st.other)
			}
		case <-st.ctx.Done():
		}
		return nil
	}()
	if err != nil {
		return err
	}

	st.closed = true
	st.cancel()
	st.subprocesses.Wait()
	close(st.chReceive)
	st.logger.Info("Stream exiting", nil)
	return nil
}

func (st *Stream2) receiveLoop() {
	chSignalMaybePending := st.demux.SignalMaybePending(st.streamID)
	chDone := st.ctx.Done()
	for {
		select {
		case <-chSignalMaybePending:
			msg, popResult := st.demux.PopMessage(st.streamID)
			switch popResult {
			case popResultEmpty:
				st.logger.Debug("Demuxer buffer is empty", nil)
			case popResultUnknownStream:
				// Closing of streams does not happen in a single step, and so
				// it could be that in the process of closing, the stream has
				// been removed from demuxer, but receiveLoop has not stopped
				// yet (but should stop soon).
				st.logger.Info("Demuxer does not know of the stream, it is likely we are in the process of closing the stream", nil)
			case popResultSuccess:
				if msg != nil {
					select {
					case st.chReceive <- msg:
					case <-chDone:
					}
				} else {
					st.logger.Error("Demuxer indicated success but we received nil msg, this should not happen", nil)
				}
			}
		case <-chDone:
			return
		}
	}
}

func (st *Stream2) sendLoop() {
	var chStreamToPeerOrNil chan<- streamIDAndMessage
	var pending streamIDAndMessage // invariant: `pending` equals the item returned by `ringBuffer.Peek()`
	var onOff bool
	pendingFilled := false

	ringBuffer := ringbuffer.NewRingBuffer[OutboundBinaryMessage](st.outgoingBufferSize)

	for {
		select {
		case onOff = <-st.chStreamOnOff:
			if onOff {
				if pendingFilled {
					chStreamToPeerOrNil = st.chStreamToConn
				}
				st.logger.Info("Turned on stream", nil)
			} else {
				chStreamToPeerOrNil = nil
				st.logger.Info("Turned off stream", nil)
			}

		case msg := <-st.chSend:
			if _, didEvict := ringBuffer.PushEvict(msg); didEvict || !pendingFilled {
				pendingMsg, _ := ringBuffer.Peek()
				pending = streamIDAndMessage{st.streamID, pendingMsg}
				pendingFilled = true
				if onOff {
					chStreamToPeerOrNil = st.chStreamToConn
				}
			}

		case chStreamToPeerOrNil <- pending:
			ringBuffer.Pop()
			if p, ok := ringBuffer.Peek(); ok {
				pending = streamIDAndMessage{st.streamID, p}
			} else {
				pendingFilled = false
				chStreamToPeerOrNil = nil
			}

		case <-st.ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////
// Types for "new" messages
////////////////////////////////////////////////////////

type ResponsePolicy = responselimit.ResponsePolicy
type SingleUseSizedLimitedResponsePolicy = responselimit.SingleUseSizedLimitedResponsePolicy

type RequestHandle internaltypes.RequestID

func (r *RequestHandle) MakeResponse(payload []byte) OutboundBinaryMessageResponse {
	return OutboundBinaryMessageResponse{
		internaltypes.RequestID(*r),
		payload,
	}
}

type InboundBinaryMessage interface {
	isInboundBinaryMessage()
}

var _ InboundBinaryMessage = InboundBinaryMessagePlain{}
var _ InboundBinaryMessage = InboundBinaryMessageRequest{}
var _ InboundBinaryMessage = InboundBinaryMessageResponse{}

type InboundBinaryMessagePlain struct {
	Payload []byte
}

func (InboundBinaryMessagePlain) isInboundBinaryMessage() {}

type InboundBinaryMessageRequest struct {
	RequestHandle RequestHandle
	Payload       []byte
}

func (InboundBinaryMessageRequest) isInboundBinaryMessage() {}

type InboundBinaryMessageResponse struct {
	Payload []byte
}

func (InboundBinaryMessageResponse) isInboundBinaryMessage() {}

type OutboundBinaryMessage interface {
	isOutboundBinaryMessage()
}

var _ OutboundBinaryMessage = OutboundBinaryMessagePlain{}
var _ OutboundBinaryMessage = OutboundBinaryMessageRequest{}
var _ OutboundBinaryMessage = OutboundBinaryMessageResponse{}

type OutboundBinaryMessagePlain struct {
	Payload []byte
}

func (OutboundBinaryMessagePlain) isOutboundBinaryMessage() {}

type OutboundBinaryMessageRequest struct {
	ResponsePolicy ResponsePolicy
	Payload        []byte
}

func (OutboundBinaryMessageRequest) isOutboundBinaryMessage() {}

type OutboundBinaryMessageResponse struct {
	requestID internaltypes.RequestID
	Payload   []byte
}

func (OutboundBinaryMessageResponse) isOutboundBinaryMessage() {}
