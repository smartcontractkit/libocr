package tracing

import (
	"strconv"

	"github.com/smartcontractkit/libocr/commontypes"
)

// OracleID pairs commontypes.OracleID with a twin flag.
// OracleID is a structure that allows an oracle and its twin to share the same id.
// While some implementations in this package use -x as the twin of x,
// this is transformed into an OracleID internally to avoid confusion!
type OracleID struct {
	OracleID commontypes.OracleID
	IsTwin   bool
}

func FromInt(id int) OracleID {
	if id < 0 {
		return OracleID{commontypes.OracleID(-id), true}
	}
	return OracleID{commontypes.OracleID(id), false}
}

func (oid OracleID) Twin() OracleID {
	return OracleID{oid.OracleID, !oid.IsTwin}
}

func (oid OracleID) Int() int {
	if oid.IsTwin {
		return -(int(oid.OracleID))
	}
	return int(oid.OracleID)
}

func (oid OracleID) String() string {
	return strconv.Itoa(oid.Int())
}
