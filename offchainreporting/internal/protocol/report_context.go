package protocol

import (
	"encoding/binary"

	"github.com/smartcontractkit/libocr/offchainreporting/types"
)








type DomainSeparationTag [32]byte

type ReportContext struct {
	ConfigDigest types.ConfigDigest
	Epoch        uint32
	Round        uint8
}

func (r ReportContext) DomainSeparationTag() (d DomainSeparationTag) {
	serialization := r.ConfigDigest[:]
	serialization = append(serialization, []byte{0, 0, 0, 0}...)
	binary.BigEndian.PutUint32(serialization[len(serialization)-4:], r.Epoch)
	serialization = append(serialization, byte(r.Round))
	copy(d[11:], serialization)
	return d
}

func (r ReportContext) Equal(r2 ReportContext) bool {
	return r.ConfigDigest == r2.ConfigDigest && r.Epoch == r2.Epoch && r.Round == r2.Round
}
