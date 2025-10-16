package ragep2pnew

import (
	"fmt"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/loghelper"
	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

// Deprecated: Please switch to Stream2.
type Stream struct {
	stream2   *stream2
	chReceive chan []byte
}

// Helper function for initializing a legacy Stream using a new Stream2.
func newStreamFromStream2(wrappedStream2 Stream2) (*Stream, error) {
	rawStream2, ok := wrappedStream2.(*stream2)
	if !ok {
		return nil, fmt.Errorf("assumption violation: wrappedStream2 is not of type *stream2")
	}
	stream := &Stream{
		rawStream2,
		make(chan []byte, 5),
	}
	go stream.receiveForwardingLoop()
	return stream, nil
}

// Other returns the peer ID of the stream counterparty.
func (st *Stream) Other() types.PeerID {
	return st.stream2.Other()
}

// Name returns the name of the stream.
func (st *Stream) Name() string {
	return st.stream2.Name()
}

// Best effort sending of messages. May fail without returning an error.
func (st *Stream) SendMessage(data []byte) {
	st.stream2.Send(OutboundBinaryMessagePlain{data})
}

// Best effort receiving of messages. The returned channel will be closed when
// the stream is closed. Note that this function may return the same channel
// across invocations.
func (st *Stream) ReceiveMessages() <-chan []byte {
	// Here, return the st.chReceive (instead of type incompatible st.stream2.chReceive).
	// See NewStream(...), which starts a go-routines forwarding all messages from st.stream2.chReceive to st.chReceive.
	return st.chReceive
}

// Close the stream. This closes any channel returned by ReceiveMessages earlier.
// After close the stream cannot be reopened. If the stream is needed in the
// future it should be created again through NewStream.
// After close, any messages passed to SendMessage will be dropped.
func (st *Stream) Close() error {
	return st.stream2.Close()
}

// Implements forwarding of received messages as workaround for incompatibility of the chReceive types.
func (st *Stream) receiveForwardingLoop() {
	defer close(st.chReceive)
	logTaper := loghelper.LogarithmicTaper{}

	for {
		select {
		case msg := <-st.stream2.Receive():
			if msg == nil {
				// stream2 was closed, so we stop the forwarding loop
				return
			}

			if msg, ok := msg.(InboundBinaryMessagePlain); ok {
				select {
				case st.chReceive <- msg.Payload:
				case <-st.stream2.ctx.Done():
					return
				}
			} else {
				logTaper.Trigger(func(newCount uint64) {
					st.stream2.logger.Warn(
						"Stream: dropping InboundBinaryMessage that is not InboundBinaryMessagePlain. Use Stream2 for support of these.",
						commontypes.LogFields{},
					)
				})
			}

		case <-st.stream2.ctx.Done():
			return
		}
	}
}
