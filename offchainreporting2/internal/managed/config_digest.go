package managed

import (
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// Convenience wrapper around OffchainConfigDigester
type prefixCheckConfigDigester struct {
	offchainConfigDigester types.OffchainConfigDigester
}

// ConfigDigest method that checks that the computed ConfigDigest's prefix is
// consistent with OffchainConfigDigester.ConfigDigestPrefix
func (d prefixCheckConfigDigester) ConfigDigest(cc types.ContractConfig) (types.ConfigDigest, error) {
	prefix := d.offchainConfigDigester.ConfigDigestPrefix()
	prefixBytes := [2]byte{}
	if types.ConfigDigestPrefix(uint16(prefix)) != prefix {
		return types.ConfigDigest{}, fmt.Errorf("did somebody change the size of ConfigDigestPrefix? this should never happen!")
	}
	binary.BigEndian.PutUint16(prefixBytes[:], uint16(prefix))

	cd, err := d.offchainConfigDigester.ConfigDigest(cc)
	if err != nil {
		return types.ConfigDigest{}, err
	}

	if cd[0] != prefixBytes[0] || cd[1] != prefixBytes[1] {
		return types.ConfigDigest{}, fmt.Errorf("ConfigDigest has prefix %v, but wanted prefix %v", cd[:2], prefixBytes)
	}

	return cd, nil
}

// Check that the ContractConfig's ConfigDigest matches the one computed
// offchain
func (d prefixCheckConfigDigester) CheckContractConfig(cc types.ContractConfig) error {
	goodConfigDigest, err := d.ConfigDigest(cc)
	if err != nil {
		return err
	}

	if goodConfigDigest != cc.ConfigDigest {
		return fmt.Errorf("ConfigDigest mismatch. Expected %x but got %x", goodConfigDigest, cc.ConfigDigest)
	}

	return nil
}
