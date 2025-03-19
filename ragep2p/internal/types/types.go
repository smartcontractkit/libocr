package types

import (
	"encoding/hex"
	"fmt"
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
