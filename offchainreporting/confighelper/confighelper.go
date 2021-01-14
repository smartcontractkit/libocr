

package confighelper

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

func makeConfigDigestArgs() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "contractAddress", Type: mustNewType("address")},
		{Name: "configCount", Type: mustNewType("uint64")},
		{Name: "signers", Type: mustNewType("address[]")},
		{Name: "transmitters", Type: mustNewType("address[]")},
		{Name: "threshold", Type: mustNewType("uint8")},
		{Name: "encodedConfigVersion", Type: mustNewType("uint64")},
		{Name: "encodedConfig", Type: mustNewType("bytes")},
	})
}

var configDigestArgs = makeConfigDigestArgs()

func configDigest(
	contractAddress common.Address,
	configCount uint64,
	oracles []common.Address,
	transmitters []common.Address,
	threshold uint8,
	encodedConfigVersion uint64,
	config []byte,
) types.ConfigDigest {
	msg, err := configDigestArgs.Pack(
		contractAddress,
		configCount,
		oracles,
		transmitters,
		threshold,
		encodedConfigVersion,
		config,
	)
	if err != nil {
		
		panic(err)
	}
	rawHash := crypto.Keccak256(msg)
	configDigest := types.ConfigDigest{}
	if n := copy(configDigest[:], rawHash); n != len(configDigest) {
		
		panic("copy too little data")
	}
	return configDigest
}

func ContractConfigFromConfigSetEvent(changed offchainaggregator.OffchainAggregatorConfigSet) types.ContractConfig {
	return types.ContractConfig{
		configDigest(
			changed.Raw.Address,
			changed.ConfigCount,
			changed.Signers,
			changed.Transmitters,
			changed.Threshold,
			changed.EncodedConfigVersion,
			changed.Encoded,
		),
		changed.Signers,
		changed.Transmitters,
		changed.Threshold,
		changed.EncodedConfigVersion,
		changed.Encoded,
	}
}

type OracleIdentity struct {
	OnChainSigningAddress           types.OnChainSigningAddress
	TransmitAddress                 common.Address
	OffchainPublicKey               types.OffchainPublicKey
	PeerID                          string
	SharedSecretEncryptionPublicKey types.SharedSecretEncryptionPublicKey
}

func ContractSetConfigArgsForIntegrationTest(
	oracles []OracleIdentity,
	f int,
	alphaPPB uint64,
) (
	signers []common.Address,
	transmitters []common.Address,
	threshold uint8,
	encodedConfigVersion uint64,
	encodedConfig []byte,
	err error,
) {
	S := []int{}
	identities := []config.OracleIdentity{}
	sharedSecretEncryptionPublicKeys := []types.SharedSecretEncryptionPublicKey{}
	for _, oracle := range oracles {
		S = append(S, 1)
		identities = append(identities, config.OracleIdentity{
			oracle.PeerID,
			oracle.OffchainPublicKey,
			oracle.OnChainSigningAddress,
			oracle.TransmitAddress,
		})
		sharedSecretEncryptionPublicKeys = append(sharedSecretEncryptionPublicKeys, oracle.SharedSecretEncryptionPublicKey)
	}
	sharedConfig := config.SharedConfig{
		config.PublicConfig{
			2 * time.Second,
			1 * time.Second,
			1 * time.Second,
			500 * time.Millisecond,
			0,
			alphaPPB,
			2 * time.Second,
			3,
			S,
			identities,
			f,
			types.ConfigDigest{},
		},
		&[config.SharedSecretSize]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
	}
	return config.XXXContractSetConfigArgsFromSharedConfig(sharedConfig, sharedSecretEncryptionPublicKeys)
}
