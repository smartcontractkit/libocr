package wire

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/time/rate"
)

// A Wire is allows communicating length-delimited bytes to and from the
// network or a buffer, up until some predefined maximum message length.
type Wire struct {
	maxMsgLength uint32
}

func NewWire(maxMsgLength uint32) *Wire {
	return &Wire{maxMsgLength: maxMsgLength}
}

func (w *Wire) WireEncode(b []byte) []byte {
	length := len(b)
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(length))
	b = append(lengthBytes, b...)
	return b
}

// NOTE: This can block indefinitely if not enough bytes are forthcoming
// It can error if the stream unexpectedly closes, or the provided data is invalid
func (w *Wire) ReadOneFromWire(r io.Reader) (payload []byte, err error) {
	lenBuf := make([]byte, 4)
	_, err = io.ReadFull(r, lenBuf)
	if err != nil {
		return nil, fmt.Errorf("error reading message length: %w", err)
	}

	msgLength := binary.BigEndian.Uint32(lenBuf)
	if msgLength > w.maxMsgLength {
		// This does not need to skip the reader pointer because the returned error will trigger a reconnection.
		return nil, fmt.Errorf("message length of %v exceeds max allowed message length of %v", msgLength, w.maxMsgLength)
	}

	payload = make([]byte, msgLength)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, fmt.Errorf("error reading blob from wire: %w", err)
	}
	return payload, nil
}

// IsNextMessageAllowed will check if the next message is permitted by the rate limiter.
// It will wait for a new message to be available on the stream reader by peeking
// at the first 4 bytes representing the new message's length.
// If the rate limiter rejects the request, the rejected message is consumed from
// the reader and discarded. This way the sync with the sender is not broken.
func (w *Wire) IsNextMessageAllowed(r *bufio.Reader, l *rate.Limiter) (bool, error) {
	lenBuf, err := r.Peek(4)
	if err != nil {
		return false, fmt.Errorf("error reading the next message's length: %w", err)
	}
	if l.Allow() {
		return true, nil
	}
	msgLength := binary.BigEndian.Uint32(lenBuf)
	_, err = r.Discard(4 + int(msgLength))
	return false, err
}
