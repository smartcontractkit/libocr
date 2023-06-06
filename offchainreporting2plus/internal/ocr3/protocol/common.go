package protocol

import (
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const ReportingPluginTimeoutWarningGracePeriod = 100 * time.Millisecond

func ByzQuorumSize(n, f int) int {
	return (n+f)/2 + 1
}

type Timestamp struct {
	ConfigDigest types.ConfigDigest
	Epoch        uint64
}
