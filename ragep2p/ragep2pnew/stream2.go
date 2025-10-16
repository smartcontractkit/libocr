package ragep2pnew

import (
	"context"
	"fmt"
	"sync"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/demuxer"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/internaltypes"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/muxer"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/responselimit"
	"github.com/RoSpaceDev/libocr/ragep2p/ragep2pnew/internal/stream2types"
	"github.com/RoSpaceDev/libocr/ragep2p/types"
	"github.com/RoSpaceDev/libocr/subprocesses"
)

type StreamPriority = stream2types.StreamPriority

const (
	StreamPriorityLow     = stream2types.StreamPriorityLow
	StreamPriorityDefault = stream2types.StreamPriorityDefault
)

type Stream2 interface {
	Other() types.PeerID
	Name() string
	Send(msg OutboundBinaryMessage)
	Receive() <-chan InboundBinaryMessage
	UpdateLimits(limits Stream2Limits) error
	Close() error
}

var _ Stream2 = &stream2{}

// Stream2 is an over-the-network channel between two peers. Two peers may share
// multiple disjoint streams with different names. Streams are persistent and
// agnostic to the state of the connection. They abstract the underlying
// connection. Messages are delivered on a best effort basis.
type stream2 struct {
	closedMu sync.Mutex
	closed   bool

	name     string
	other    types.PeerID
	streamID internaltypes.StreamID

	maxOutgoingBufferedMessages int
	maxMessageLength            int

	host *Host

	subprocesses subprocesses.Subprocesses
	ctx          context.Context
	cancel       context.CancelFunc
	logger       loghelper.LoggerWithContext
	chReceive    chan InboundBinaryMessage

	mux   *muxer.Muxer
	demux *demuxer.Demuxer

	chStreamUpdateLimitsRequest  chan<- peerStreamUpdateLimitsRequest
	chStreamUpdateLimitsResponse <-chan peerStreamUpdateLimitsResponse

	chStreamCloseRequest  chan<- peerStreamCloseRequest
	chStreamCloseResponse <-chan peerStreamCloseResponse
}

// Other returns the peer ID of the stream's counterparty.
func (st *stream2) Other() types.PeerID {
	return st.other
}

// Name returns the name of the stream.
func (st *stream2) Name() string {
	return st.name
}

// Best effort sending of messages. May fail without returning an error.
func (st *stream2) Send(msg OutboundBinaryMessage) {
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
		ok = len(msg.Payload) <= types.MaxMessageLength
		payloadLength = len(msg.Payload)
	default:
		panic(fmt.Sprintf("unknown OutboundBinaryMessage type: %T", msg))
	}

	if !ok {
		st.logger.Warn("dropping outbound message that is too large", commontypes.LogFields{
			"messagePayloadLength":    payloadLength,
			"streamMaxMessageLength":  st.maxMessageLength,
			"ragep2pMaxMessageLength": types.MaxMessageLength,
		})
		return
	}

	_ = st.mux.PushEvict(st.streamID, msg)
}

// Best effort receiving of messages. The returned channel will be closed when
// the stream is closed. Note that this function may return the same channel
// across invocations.
func (st *stream2) Receive() <-chan InboundBinaryMessage {
	return st.chReceive
}

func (st *stream2) UpdateLimits(limits Stream2Limits) error {
	validatedLimits, err := limits.Validate()
	if err != nil {
		return err
	}

	select {
	case st.chStreamUpdateLimitsRequest <- peerStreamUpdateLimitsRequest{st.streamID, validatedLimits}:
		resp := <-st.chStreamUpdateLimitsResponse
		if resp.err != nil {
			return resp.err
		}
		return nil
	case <-st.ctx.Done():

		return fmt.Errorf("UpdateLimts: called after Stream internal context already expired")
	}
}

// Close the stream. This closes any channel returned by ReceiveMessages earlier.
// After close the stream cannot be reopened. If the stream is needed in the
// future it should be created again through NewStream2.
// After close, any messages passed to SendMessage will be dropped.
func (st *stream2) Close() error {
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

func (st *stream2) receiveLoop() {
	chSignalMaybePending := st.demux.SignalMaybePending(st.streamID)
	chDone := st.ctx.Done()
	for {
		select {
		case <-chSignalMaybePending:
			msg, popResult := st.demux.PopMessage(st.streamID)
			switch popResult {
			case demuxer.PopResultEmpty:
				st.logger.Debug("Demuxer buffer is empty", nil)
			case demuxer.PopResultUnknownStream:
				// Closing of streams does not happen in a single step, and so
				// it could be that in the process of closing, the stream has
				// been removed from demuxer, but receiveLoop has not stopped
				// yet (but should stop soon).
				st.logger.Info("Demuxer does not know of the stream, it is likely we are in the process of closing the stream", nil)
			case demuxer.PopResultSuccess:
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

////////////////////////////////////////////////////////
// Types for "new" messages
////////////////////////////////////////////////////////

type Stream2Limits = stream2types.Stream2Limits

type ResponsePolicy = responselimit.ResponsePolicy
type SingleUseSizedLimitedResponsePolicy = responselimit.SingleUseSizedLimitedResponsePolicy

type RequestHandle = stream2types.RequestHandle

type InboundBinaryMessage = stream2types.InboundBinaryMessage
type InboundBinaryMessagePlain = stream2types.InboundBinaryMessagePlain
type InboundBinaryMessageRequest = stream2types.InboundBinaryMessageRequest
type InboundBinaryMessageResponse = stream2types.InboundBinaryMessageResponse

type OutboundBinaryMessage = stream2types.OutboundBinaryMessage
type OutboundBinaryMessagePlain = stream2types.OutboundBinaryMessagePlain
type OutboundBinaryMessageRequest = stream2types.OutboundBinaryMessageRequest
type OutboundBinaryMessageResponse = stream2types.OutboundBinaryMessageResponse
