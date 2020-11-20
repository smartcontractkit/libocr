package networking

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const (
	
	
	MaxMsgLength = 10000
)

func wireEncode(b []byte) []byte {
	length := len(b)
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(length))
	b = append(lengthBytes, b...)
	return b
}



func readOneFromWire(r io.Reader) (payload []byte, err error) {
	lenBuf := make([]byte, 4)
	_, err = io.ReadFull(r, lenBuf)
	if err != nil {
		return nil, errors.Wrap(err, "error reading message length")
	}

	msgLength := binary.BigEndian.Uint32(lenBuf)
	if msgLength > MaxMsgLength {
		
		return nil, errors.Errorf("message length of %v exceeds max allowed message length of %v", msgLength, MaxMsgLength)
	}

	payload = make([]byte, msgLength)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, errors.Wrap(err, "error reading blob from wire")
	}
	return payload, nil
}






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
