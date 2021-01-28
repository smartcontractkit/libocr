package networking

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const (
	// MaxMsgLength is the maximum allowed length for a data payload in bytes
	// NOTE: This is slightly larger than 2x of the largest message we can possibly send, assuming N=31.
	MaxMsgLength = 10000
)

func wireEncode(b []byte) []byte {
	length := len(b)
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(length))
	b = append(lengthBytes, b...)
	return b
}

// NOTE: This can block indefinitely if not enough bytes are forthcoming
// It can error if the stream unexpectedly closes, or the provided data is invalid
func readOneFromWire(r io.Reader) (payload []byte, err error) {
	lenBuf := make([]byte, 4)
	_, err = io.ReadFull(r, lenBuf)
	if err != nil {
		return nil, errors.Wrap(err, "error reading message length")
	}

	msgLength := binary.BigEndian.Uint32(lenBuf)
	if msgLength > MaxMsgLength {
		// This does not need to skip the reader pointer because the returned error will trigger a reconnection.
		return nil, errors.Errorf("message length of %v exceeds max allowed message length of %v", msgLength, MaxMsgLength)
	}

	payload = make([]byte, msgLength)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, errors.Wrap(err, "error reading blob from wire")
	}
	return payload, nil
}

// isNextMessageAllowed will check if the next message is permitted by the rate limiter.
// It will wait for a new message to be available on the stream reader by peeking
// at the first 4 bytes representing the new message's length.
// If the rate limiter rejects the request, the rejected message is consumed from
// the reader and discarded. This way the sync with the sender is not broken.
func isNextMessageAllowed(r *bufio.Reader, l limiter) (bool, error) {
	lenBuf, err := r.Peek(4)
	if err != nil {
		return false, errors.Wrap(err, "error reading the next message's length")
	}
	if l.Allow() {
		return true, nil
	}
	msgLength := binary.BigEndian.Uint32(lenBuf)
	_, err = r.Discard(4 + int(msgLength))
	return false, err
}
