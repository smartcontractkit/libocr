// Package confighelper provides helpers for converting between the gethwrappers/OCR2Aggregator.SetConfig
// event and types.ContractConfig
package confighelper

import (
	"crypto/rand"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// OracleIdentity is identical to the internal type in package config.
// We intentionally make a copy to make potential future internal modifications easier.
type OracleIdentity struct {
	OnChainSigningAddress types.OnchainPublicKey
	TransmitAccount       types.Account
	OffchainPublicKey     types.OffchainPublicKey
	PeerID                string
}

// PublicConfig is identical to the internal type in package config.
// We intentionally make a copy to make potential future internal modifications easier.
type PublicConfig struct {
	DeltaProgress    time.Duration
	DeltaResend      time.Duration
	DeltaRound       time.Duration
	DeltaGrace       time.Duration
	DeltaStage       time.Duration
	RMax             uint8
	S                []int
	OracleIdentities []OracleIdentity

	ReportingPluginConfig []byte

	MaxDurationQuery                        time.Duration
	MaxDurationObservation                  time.Duration
	MaxDurationReport                       time.Duration
	MaxDurationShouldAcceptFinalizedReport  time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration

	F             int
	OnchainConfig []byte
	ConfigDigest  types.ConfigDigest
}

func (pc PublicConfig) N() int {
	return len(pc.OracleIdentities)
}

func PublicConfigFromContractConfig(chainID *big.Int, skipChainSpecificChecks bool, change types.ContractConfig) (PublicConfig, error) {
	internalPublicConfig, err := config.PublicConfigFromContractConfig(chainID, skipChainSpecificChecks, change)
	if err != nil {
		return PublicConfig{}, err
	}
	identities := []OracleIdentity{}
	for _, internalIdentity := range internalPublicConfig.OracleIdentities {
		identities = append(identities, OracleIdentity{
			internalIdentity.OnChainPublicKey,
			internalIdentity.TransmitAccount,
			internalIdentity.OffchainPublicKey,
			internalIdentity.PeerID,
		})
	}
	return PublicConfig{
		internalPublicConfig.DeltaProgress,
		internalPublicConfig.DeltaResend,
		internalPublicConfig.DeltaRound,
		internalPublicConfig.DeltaGrace,
		internalPublicConfig.DeltaStage,
		internalPublicConfig.RMax,
		internalPublicConfig.S,
		identities,
		internalPublicConfig.ReportingPluginConfig,
		internalPublicConfig.MaxDurationQuery,
		internalPublicConfig.MaxDurationObservation,
		internalPublicConfig.MaxDurationReport,
		internalPublicConfig.MaxDurationShouldAcceptFinalizedReport,
		internalPublicConfig.MaxDurationShouldTransmitAcceptedReport,
		internalPublicConfig.F,
		internalPublicConfig.OnchainConfig,
		internalPublicConfig.ConfigDigest,
	}, nil
}

type OracleIdentityExtra struct {
	OracleIdentity
	ConfigEncryptionPublicKey types.ConfigEncryptionPublicKey
}

// ContractSetConfigArgsForIntegrationTest generates setConfig args for integration tests in core
func ContractSetConfigArgsForIntegrationTest(
	oracles []OracleIdentityExtra,
	f int,
	alphaPPB uint64,
) (
	signers []common.Address,
	transmitters []common.Address,
	f_ uint8,
	onchainConfig []byte,
	encodedConfigVersion uint64,
	encodedConfig []byte,
	err error,
) {
	S := []int{}
	identities := []config.OracleIdentity{}
	sharedSecretEncryptionPublicKeys := []types.ConfigEncryptionPublicKey{}
	for _, oracle := range oracles {
		S = append(S, 1)
		identities = append(identities, config.OracleIdentity{
			oracle.PeerID,
			oracle.OffchainPublicKey,
			oracle.OnChainSigningAddress,
			oracle.TransmitAccount,
		})
		sharedSecretEncryptionPublicKeys = append(sharedSecretEncryptionPublicKeys, oracle.ConfigEncryptionPublicKey)
	}
	sharedConfig := config.SharedConfig{
		config.PublicConfig{
			2 * time.Second,
			1 * time.Second,
			1 * time.Second,
			500 * time.Millisecond,
			2 * time.Second,
			3,
			S,
			identities,
			median.Config{
				alphaPPB,
				0,
			}.Encode(),
			50 * time.Millisecond,
			50 * time.Millisecond,
			50 * time.Millisecond,
			50 * time.Millisecond,
			50 * time.Millisecond,
			f,
			[]byte{},
			types.ConfigDigest{},
		},
		&[config.SharedSecretSize]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
	}
	return config.XXXContractSetConfigArgsFromSharedConfig(sharedConfig, sharedSecretEncryptionPublicKeys)
}

// // ContractSetConfigArgs generates setConfig args from the relevant parameters.
// // Only use this for testing, *not* for production.
func ContractSetConfigArgs(
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaStage time.Duration,
	rMax uint8,
	s []int,
	oracles []OracleIdentityExtra,
	reportingPluginConfig []byte,
	maxDurationQuery time.Duration,
	maxDurationObservation time.Duration,
	maxDurationReport time.Duration,
	maxDurationShouldAcceptFinalizedReport time.Duration,
	maxDurationShouldTransmitAcceptedReport time.Duration,

	f int,
	onchainConfig []byte,
) (
	signers []common.Address,
	transmitters []common.Address,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	identities := []config.OracleIdentity{}
	configEncryptionPublicKeys := []types.ConfigEncryptionPublicKey{}
	for _, oracle := range oracles {
		identities = append(identities, config.OracleIdentity{
			oracle.PeerID,
			oracle.OffchainPublicKey,
			oracle.OnChainSigningAddress,
			oracle.TransmitAccount,
		})
		configEncryptionPublicKeys = append(configEncryptionPublicKeys, oracle.ConfigEncryptionPublicKey)
	}

	sharedSecret := [config.SharedSecretSize]byte{}
	if _, err := io.ReadFull(rand.Reader, sharedSecret[:]); err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}

	sharedConfig := config.SharedConfig{
		config.PublicConfig{
			deltaProgress,
			deltaResend,
			deltaRound,
			deltaGrace,
			deltaStage,
			rMax,
			s,
			identities,
			reportingPluginConfig,
			maxDurationQuery,
			maxDurationObservation,
			maxDurationReport,
			maxDurationShouldAcceptFinalizedReport,
			maxDurationShouldTransmitAcceptedReport,
			f,
			onchainConfig,
			types.ConfigDigest{},
		},
		&sharedSecret,
	}
	return config.XXXContractSetConfigArgsFromSharedConfig(sharedConfig, configEncryptionPublicKeys)
}
