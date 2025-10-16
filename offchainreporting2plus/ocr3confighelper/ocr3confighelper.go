package ocr3confighelper

import (
	"crypto/rand"
	"time"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/confighelper"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
	"golang.org/x/crypto/curve25519"
)

// PublicConfig is identical to the internal type [ocr3config.PublicConfig]. See
// the documentation there for details. We intentionally duplicate the internal
// type to make potential future internal modifications easier.
type PublicConfig struct {
	DeltaProgress               time.Duration
	DeltaResend                 time.Duration
	DeltaInitial                time.Duration
	DeltaRound                  time.Duration
	DeltaGrace                  time.Duration
	DeltaCertifiedCommitRequest time.Duration
	DeltaStage                  time.Duration
	RMax                        uint64
	S                           []int
	OracleIdentities            []confighelper.OracleIdentity

	ReportingPluginConfig []byte

	MaxDurationInitialization               *time.Duration
	MaxDurationQuery                        time.Duration
	MaxDurationObservation                  time.Duration
	MaxDurationShouldAcceptAttestedReport   time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration

	F             int
	OnchainConfig []byte
	ConfigDigest  types.ConfigDigest
}

func (pc PublicConfig) N() int {
	return len(pc.OracleIdentities)
}

func PublicConfigFromContractConfig(skipResourceExhaustionChecks bool, change types.ContractConfig) (PublicConfig, error) {
	internalPublicConfig, err := ocr3config.PublicConfigFromContractConfig(skipResourceExhaustionChecks, change)
	if err != nil {
		return PublicConfig{}, err
	}
	identities := []confighelper.OracleIdentity{}
	for _, internalIdentity := range internalPublicConfig.OracleIdentities {
		identities = append(identities, confighelper.OracleIdentity{
			internalIdentity.OffchainPublicKey,
			internalIdentity.OnchainPublicKey,
			internalIdentity.PeerID,
			internalIdentity.TransmitAccount,
		})
	}
	return PublicConfig{
		internalPublicConfig.DeltaProgress,
		internalPublicConfig.DeltaResend,
		internalPublicConfig.DeltaInitial,
		internalPublicConfig.DeltaRound,
		internalPublicConfig.DeltaGrace,
		internalPublicConfig.DeltaCertifiedCommitRequest,
		internalPublicConfig.DeltaStage,
		internalPublicConfig.RMax,
		internalPublicConfig.S,
		identities,
		internalPublicConfig.ReportingPluginConfig,
		internalPublicConfig.MaxDurationInitialization,
		internalPublicConfig.MaxDurationQuery,
		internalPublicConfig.MaxDurationObservation,
		internalPublicConfig.MaxDurationShouldAcceptAttestedReport,
		internalPublicConfig.MaxDurationShouldTransmitAcceptedReport,
		internalPublicConfig.F,
		internalPublicConfig.OnchainConfig,
		internalPublicConfig.ConfigDigest,
	}, nil
}

// ContractSetConfigArgsForTestsWithAuxiliaryArgsMercuryV02 generates setConfig
// args for mercury v0.2. Only use this for testing, *not* for production.
// See [ocr3config.PublicConfig] for documentation of the arguments.
func ContractSetConfigArgsForTestsMercuryV02(
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaInitial time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaCertifiedCommitRequest time.Duration,
	deltaStage time.Duration,
	rMax uint8,
	s []int,
	oracles []confighelper.OracleIdentityExtra,
	reportingPluginConfig []byte,
	maxDurationInitialization *time.Duration,
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
	return ContractSetConfigArgsForTests(
		deltaProgress,
		deltaResend,
		deltaInitial,
		deltaRound,
		deltaGrace,
		deltaCertifiedCommitRequest,
		deltaStage,
		uint64(rMax),
		s,
		oracles,
		reportingPluginConfig,
		maxDurationInitialization,
		0,
		maxDurationObservation,
		0,
		0,
		f,
		onchainConfig,
	)
}

// ContractSetConfigArgsForTestsOCR3 generates setConfig args for OCR3. Only use
// this for testing, *not* for production.
// See [ocr3config.PublicConfig] for documentation of the arguments.
func ContractSetConfigArgsForTests(
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaInitial time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaCertifiedCommitRequest time.Duration,
	deltaStage time.Duration,
	rMax uint64,
	s []int,
	oracles []confighelper.OracleIdentityExtra,
	reportingPluginConfig []byte,
	maxDurationInitialization *time.Duration,
	maxDurationQuery time.Duration,
	maxDurationObservation time.Duration,
	maxDurationShouldAcceptAttestedReport time.Duration,
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
	ephemeralSk := [curve25519.ScalarSize]byte{}
	if _, err := rand.Read(ephemeralSk[:]); err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}

	sharedSecret := [config.SharedSecretSize]byte{}
	if _, err := rand.Read(sharedSecret[:]); err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}

	return ContractSetConfigArgsDeterministic(
		ephemeralSk,
		sharedSecret,

		deltaProgress,
		deltaResend,
		deltaInitial,
		deltaRound,
		deltaGrace,
		deltaCertifiedCommitRequest,
		deltaStage,
		rMax,
		s,
		oracles,
		reportingPluginConfig,
		maxDurationInitialization,
		maxDurationQuery,
		maxDurationObservation,
		maxDurationShouldAcceptAttestedReport,
		maxDurationShouldTransmitAcceptedReport,
		f,
		onchainConfig,
	)
}

// This function may be used in production. If you use this as part of multisig
// tooling,  make sure that the input parameters are identical across all
// signers.
// See [ocr3config.PublicConfig] for documentation of the arguments.
func ContractSetConfigArgsDeterministic(
	// The ephemeral secret key used to encrypt the shared secret
	ephemeralSk [curve25519.ScalarSize]byte,
	// A secret shared between all oracles, enabling them to derive pseudorandom
	// leaders and transmitters. This is a low-value secret. An adversary who
	// learns this secret can only determine which oracle will be
	// leader/transmitter ahead of time, that's it.
	sharedSecret [config.SharedSecretSize]byte,
	// Check out the ocr3config package for documentation on the following args
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaInitial time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaCertifiedCommitRequest time.Duration,
	deltaStage time.Duration,
	rMax uint64,
	s []int,
	oracles []confighelper.OracleIdentityExtra,
	reportingPluginConfig []byte,
	maxDurationInitialization *time.Duration,
	maxDurationQuery time.Duration,
	maxDurationObservation time.Duration,
	maxDurationShouldAcceptAttestedReport time.Duration,
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
	identities := []config.OracleIdentity{}
	configEncryptionPublicKeys := []types.ConfigEncryptionPublicKey{}
	for _, oracle := range oracles {
		identities = append(identities, config.OracleIdentity{
			oracle.OffchainPublicKey,
			oracle.OnchainPublicKey,
			oracle.PeerID,
			oracle.TransmitAccount,
		})
		configEncryptionPublicKeys = append(configEncryptionPublicKeys, oracle.ConfigEncryptionPublicKey)
	}

	sharedConfig := ocr3config.SharedConfig{
		ocr3config.PublicConfig{
			deltaProgress,
			deltaResend,
			deltaInitial,
			deltaRound,
			deltaGrace,
			deltaCertifiedCommitRequest,
			deltaStage,
			rMax,
			s,
			identities,
			reportingPluginConfig,
			maxDurationInitialization,
			maxDurationQuery,
			maxDurationObservation,
			maxDurationShouldAcceptAttestedReport,
			maxDurationShouldTransmitAcceptedReport,
			f,
			onchainConfig,
			types.ConfigDigest{},
		},
		&sharedSecret,
	}
	return ocr3config.ContractSetConfigArgsFromSharedConfigDeterministic(
		sharedConfig,
		configEncryptionPublicKeys,
		&ephemeralSk,
	)
}
