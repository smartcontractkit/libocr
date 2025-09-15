package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type StreamID struct {
	OracleID commontypes.OracleID
	Priority ocr2types.BinaryMessageOutboundPriority
}

const RequestIDSize = 32

type RequestID [RequestIDSize]byte

var _ fmt.Stringer = RequestID{}

func (r RequestID) String() string {
	return hex.EncodeToString(r[:])
}

func GetRandomRequestID() RequestID {
	var b [RequestIDSize]byte
	_, _ = rand.Read(b[:])
	return RequestID(b)
}
