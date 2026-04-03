package ocr3shims

import (
	"bytes"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type onchainKeyring2Shim[RI any] struct {
	ocr3types.OnchainKeyring[RI]
}

func (k *onchainKeyring2Shim[RI]) Has(onchainPublicKey types.OnchainPublicKey) bool {
	return bytes.Equal(k.OnchainKeyring.PublicKey(), onchainPublicKey)
}

func (k *onchainKeyring2Shim[RI]) DebugIdentifier() string {
	return fmt.Sprintf("%x", k.OnchainKeyring.PublicKey())
}

// OnchainKeyringAsOnchainKeyring2 wraps an [ocr3types.OnchainKeyring] and returns an [ocr3types.OnchainKeyring2].
// If k already implements OnchainKeyring2, it is returned as-is.
func OnchainKeyringAsOnchainKeyring2[RI any](k ocr3types.OnchainKeyring[RI]) ocr3types.OnchainKeyring2[RI] {
	if c, ok := any(k).(ocr3types.OnchainKeyring2[RI]); ok {
		return c
	}
	return &onchainKeyring2Shim[RI]{k}
}
