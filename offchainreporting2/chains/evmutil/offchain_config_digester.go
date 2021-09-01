package evmutil

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ types.OffchainConfigDigester = EVMOffchainConfigDigester{}

type EVMOffchainConfigDigester struct {
	ChainID         uint64
	ContractAddress common.Address
}

func (d EVMOffchainConfigDigester) ConfigDigest(cc types.ContractConfig) (types.ConfigDigest, error) {
	signers := []common.Address{}
	for _, signer := range cc.Signers {
		a := common.BytesToAddress(signer)
		signers = append(signers, a)
	}
	transmitters := []common.Address{}
	for _, transmitter := range cc.Transmitters {
		a := common.HexToAddress(string(transmitter))
		transmitters = append(transmitters, a)
	}

	return configDigest(
		d.ChainID,
		d.ContractAddress,
		cc.ConfigCount,
		signers,
		transmitters,
		cc.F,
		cc.OnchainConfig,
		cc.OffchainConfigVersion,
		cc.OffchainConfig,
	), nil
}

func (d EVMOffchainConfigDigester) ConfigDigestPrefix() types.ConfigDigestPrefix {
	return types.ConfigDigestPrefixEVM
}
