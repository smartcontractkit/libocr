package internaltypes

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/RoSpaceDev/libocr/ragep2p/types"
)

const (
	StreamIDSize  = 32
	RequestIDSize = 32
)

type StreamID [StreamIDSize]byte
type RequestID [RequestIDSize]byte

var _ fmt.Stringer = StreamID{}
var _ fmt.Stringer = RequestID{}

func (s StreamID) String() string {
	return hex.EncodeToString(s[:])
}

func (r RequestID) String() string {
	return hex.EncodeToString(r[:])
}

func MakeStreamID(self types.PeerID, other types.PeerID, name string) StreamID {
	if bytes.Compare(self[:], other[:]) < 0 {
		return MakeStreamID(other, self, name)
	}

	h := sha256.New()
	_, _ = h.Write(self[:])
	_, _ = h.Write(other[:])
	// this is fine because self and other are of constant length. if more than
	// one variable length item is ever added here, we should also hash lengths
	// to prevent collisions.
	_, _ = h.Write([]byte(name))

	var result StreamID
	_ = h.Sum(result[:0])
	return result
}

func MakeRandomRequestID() (RequestID, error) {
	requestID := RequestID{}
	_, err := rand.Read(requestID[:])
	return requestID, err
}
