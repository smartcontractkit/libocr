// Package confighelper provides helpers for converting between the gethwrappers/OCR2Aggregator.SetConfig
// event and types.ContractConfig
package confighelper

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// OracleIdentity is identical to the internal type in package config.
// We intentionally make a copy to make potential future internal modifications easier.
type OracleIdentity = confighelper.OracleIdentity

// PublicConfig is identical to the internal type in package config.
// We intentionally make a copy to make potential future internal modifications easier.
type PublicConfig = confighelper.PublicConfig

func PublicConfigFromContractConfig(skipResourceExhaustionChecks bool, change types.ContractConfig) (PublicConfig, error) {
	return confighelper.PublicConfigFromContractConfig(skipResourceExhaustionChecks, change)
}

type OracleIdentityExtra = confighelper.OracleIdentityExtra

// ContractSetConfigArgsForIntegrationTest generates setConfig args for integration tests in core.
// Only use this for testing, *not* for production.
func ContractSetConfigArgsForEthereumIntegrationTest(
	oracles []OracleIdentityExtra,
	f int,
	alphaPPB uint64,
) (
	signers []common.Address,
	transmitters []common.Address,
	f_ uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	return confighelper.ContractSetConfigArgsForEthereumIntegrationTest(oracles, f, alphaPPB)
}

// ContractSetConfigArgsForTestsWithAuxiliaryArgs generates setConfig args from
// the relevant parameters. Only use this for testing, *not* for production.
func ContractSetConfigArgsForTestsWithAuxiliaryArgs(
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
	auxiliaryArgs AuxiliaryArgs,
) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	return confighelper.ContractSetConfigArgsForTestsWithAuxiliaryArgs(
		deltaProgress,
		deltaResend,
		deltaRound,
		deltaGrace,
		deltaStage,
		rMax,
		s,
		oracles,
		reportingPluginConfig,
		maxDurationQuery,
		maxDurationObservation,
		maxDurationReport,
		maxDurationShouldAcceptFinalizedReport,
		maxDurationShouldTransmitAcceptedReport,

		f,
		onchainConfig,
		auxiliaryArgs,
	)
}

type AuxiliaryArgs = confighelper.AuxiliaryArgs

// ContractSetConfigArgsForTests generates setConfig args from the relevant
// parameters. Only use this for testing, *not* for production.
func ContractSetConfigArgsForTests(
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
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	return confighelper.ContractSetConfigArgsForTests(
		deltaProgress,
		deltaResend,
		deltaRound,
		deltaGrace,
		deltaStage,
		rMax,
		s,
		oracles,
		reportingPluginConfig,
		maxDurationQuery,
		maxDurationObservation,
		maxDurationReport,
		maxDurationShouldAcceptFinalizedReport,
		maxDurationShouldTransmitAcceptedReport,
		f,
		onchainConfig,
	)
}

// ContractSetConfigArgsForTestsWithAuxiliaryArgsMercuryV02 generates setConfig
// args for mercury v0.2. Only use this for testing, *not* for production.
func ContractSetConfigArgsForTestsMercuryV02(
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaStage time.Duration,
	rMax uint8,
	s []int,
	oracles []OracleIdentityExtra,
	reportingPluginConfig []byte,
	maxDurationObservation time.Duration,
	f int,
	onchainConfig []byte,
) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	return confighelper.ContractSetConfigArgsForTestsMercuryV02(
		deltaProgress,
		deltaResend,
		deltaRound,
		deltaGrace,
		deltaStage,
		rMax,
		s,
		oracles,
		reportingPluginConfig,
		maxDurationObservation,
		f,
		onchainConfig,
	)
}
