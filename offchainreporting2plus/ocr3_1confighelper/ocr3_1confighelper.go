package ocr3_1confighelper

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/crypto/curve25519"
)

// PublicConfig is identical to the internal type [ocr3_1config.PublicConfig]. See
// the documentation there for details. We intentionally duplicate the internal
// type to make potential future internal modifications easier.
type PublicConfig struct {
	OracleIdentities []confighelper.OracleIdentity
	F                int

	// pacemaker

	DeltaProgress time.Duration
	DeltaResend   *time.Duration

	// outcome generation

	DeltaInitial *time.Duration
	DeltaRound   time.Duration
	DeltaGrace   time.Duration
	RMax         uint64

	// report attestation

	DeltaReportsPlusPrecursorRequest *time.Duration

	// transmission

	DeltaStage time.Duration
	S          []int

	// state sync

	DeltaStateSyncSummaryInterval *time.Duration

	// block sync

	DeltaBlockSyncMinRequestToSameOracleInterval *time.Duration
	DeltaBlockSyncResponseTimeout                *time.Duration
	MaxBlocksPerBlockSyncResponse                *int
	MaxParallelRequestedBlocks                   *uint64

	// tree sync

	DeltaTreeSyncMinRequestToSameOracleInterval *time.Duration
	DeltaTreeSyncResponseTimeout                *time.Duration
	MaxTreeSyncChunkKeys                        *int
	MaxTreeSyncChunkKeysPlusValuesBytes         *int
	MaxParallelTreeSyncChunkFetches             *int

	// snapshotting

	SnapshotInterval               *uint64
	MaxHistoricalSnapshotsRetained *uint64

	// blobs

	DeltaBlobOfferMinRequestToSameOracleInterval *time.Duration
	DeltaBlobOfferResponseTimeout                *time.Duration
	DeltaBlobBroadcastGrace                      *time.Duration
	DeltaBlobChunkMinRequestToSameOracleInterval *time.Duration
	DeltaBlobChunkResponseTimeout                *time.Duration
	BlobChunkBytes                               *int

	// reporting plugin

	ReportingPluginConfig []byte
	OnchainConfig         []byte

	MaxDurationInitialization               time.Duration
	WarnDurationQuery                       time.Duration
	WarnDurationObservation                 time.Duration
	WarnDurationValidateObservation         time.Duration
	WarnDurationObservationQuorum           time.Duration
	WarnDurationStateTransition             time.Duration
	WarnDurationCommitted                   time.Duration
	MaxDurationShouldAcceptAttestedReport   time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration

	ConfigDigest types.ConfigDigest
}

func (pc PublicConfig) N() int {
	return len(pc.OracleIdentities)
}

type CheckPublicConfigLevel string

const (
	CheckPublicConfigLevelDefault                   CheckPublicConfigLevel = ""
	CheckPublicConfigLevelDangerInsaneForProduction CheckPublicConfigLevel = "danger_insane_for_production"
)

func (c CheckPublicConfigLevel) shouldSkipInsaneForProductionChecks() bool {
	switch c {
	case CheckPublicConfigLevelDefault:
		return false
	case CheckPublicConfigLevelDangerInsaneForProduction:
		return true
	default:
		panic(fmt.Sprintf("invalid checkPublicConfigLevel: %v", c))
	}
}

func PublicConfigFromContractConfig(checkPublicConfigLevel CheckPublicConfigLevel, change types.ContractConfig) (PublicConfig, error) {
	internalPublicConfig, err := ocr3_1config.PublicConfigFromContractConfig(checkPublicConfigLevel.shouldSkipInsaneForProductionChecks(), change)
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
		identities,
		internalPublicConfig.F,

		// outcome generation

		internalPublicConfig.DeltaProgress,
		internalPublicConfig.DeltaResend,
		internalPublicConfig.DeltaInitial,
		internalPublicConfig.DeltaRound,
		internalPublicConfig.DeltaGrace,
		internalPublicConfig.RMax,

		// report attestation

		internalPublicConfig.DeltaReportsPlusPrecursorRequest,

		// transmission

		internalPublicConfig.DeltaStage,
		internalPublicConfig.S,

		// state sync

		internalPublicConfig.DeltaStateSyncSummaryInterval,

		// block sync

		internalPublicConfig.DeltaBlockSyncMinRequestToSameOracleInterval,
		internalPublicConfig.DeltaBlockSyncResponseTimeout,
		internalPublicConfig.MaxBlocksPerBlockSyncResponse,
		internalPublicConfig.MaxParallelRequestedBlocks,

		// tree sync

		internalPublicConfig.DeltaTreeSyncMinRequestToSameOracleInterval,
		internalPublicConfig.DeltaTreeSyncResponseTimeout,
		internalPublicConfig.MaxTreeSyncChunkKeys,
		internalPublicConfig.MaxTreeSyncChunkKeysPlusValuesBytes,
		internalPublicConfig.MaxParallelTreeSyncChunkFetches,

		// snapshotting

		internalPublicConfig.SnapshotInterval,
		internalPublicConfig.MaxHistoricalSnapshotsRetained,

		// blobs

		internalPublicConfig.DeltaBlobOfferMinRequestToSameOracleInterval,
		internalPublicConfig.DeltaBlobOfferResponseTimeout,
		internalPublicConfig.DeltaBlobBroadcastGrace,
		internalPublicConfig.DeltaBlobChunkMinRequestToSameOracleInterval,
		internalPublicConfig.DeltaBlobChunkResponseTimeout,
		internalPublicConfig.BlobChunkBytes,

		// reporting plugin

		internalPublicConfig.ReportingPluginConfig,
		internalPublicConfig.OnchainConfig,

		internalPublicConfig.MaxDurationInitialization,
		internalPublicConfig.WarnDurationQuery,
		internalPublicConfig.WarnDurationObservation,
		internalPublicConfig.WarnDurationValidateObservation,
		internalPublicConfig.WarnDurationObservationQuorum,
		internalPublicConfig.WarnDurationStateTransition,
		internalPublicConfig.WarnDurationCommitted,
		internalPublicConfig.MaxDurationShouldAcceptAttestedReport,
		internalPublicConfig.MaxDurationShouldTransmitAcceptedReport,
		internalPublicConfig.ConfigDigest,
	}, nil
}

// ContractSetConfigArgsForTests generates setConfig args for testing. Only use
// this for testing, *not* for production.
// See [ocr3_1config.PublicConfig] for documentation of the arguments.
func ContractSetConfigArgsForTests(
	checkPublicConfigLevel CheckPublicConfigLevel,
	oracles []confighelper.OracleIdentityExtra,
	f int,
	deltaProgress time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	rMax uint64,
	deltaStage time.Duration,
	s []int,
	reportingPluginConfig []byte,
	onchainConfig []byte,
	maxDurationInitialization time.Duration,
	warnDurationQuery time.Duration,
	warnDurationObservation time.Duration,
	warnDurationValidateObservation time.Duration,
	warnDurationObservationQuorum time.Duration,
	warnDurationStateTransition time.Duration,
	warnDurationCommitted time.Duration,
	maxDurationShouldAcceptAttestedReport time.Duration,
	maxDurationShouldTransmitAcceptedReport time.Duration,
	// Optional configuration parameters
	optionalConfig ContractSetConfigArgsOptionalConfig,
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
		checkPublicConfigLevel,
		ephemeralSk,
		sharedSecret,
		oracles,
		f,
		deltaProgress,
		deltaRound,
		deltaGrace,
		rMax,
		deltaStage,
		s,
		reportingPluginConfig,
		onchainConfig,
		maxDurationInitialization,
		warnDurationQuery,
		warnDurationObservation,
		warnDurationValidateObservation,
		warnDurationObservationQuorum,
		warnDurationStateTransition,
		warnDurationCommitted,
		maxDurationShouldAcceptAttestedReport,
		maxDurationShouldTransmitAcceptedReport,
		optionalConfig,
	)
}

// ContractSetConfigArgsOptionalConfig holds optional parameters
// for ContractSetConfigArgsDeterministic.
type ContractSetConfigArgsOptionalConfig struct {
	DeltaResend                      *time.Duration
	DeltaInitial                     *time.Duration
	DeltaReportsPlusPrecursorRequest *time.Duration

	// state sync
	DeltaStateSyncSummaryInterval *time.Duration

	// block sync
	DeltaBlockSyncMinRequestToSameOracleInterval *time.Duration
	DeltaBlockSyncResponseTimeout                *time.Duration
	MaxBlocksPerBlockSyncResponse                *int
	MaxParallelRequestedBlocks                   *uint64

	// tree sync
	DeltaTreeSyncMinRequestToSameOracleInterval *time.Duration
	DeltaTreeSyncResponseTimeout                *time.Duration
	MaxTreeSyncChunkKeys                        *int
	MaxTreeSyncChunkKeysPlusValuesBytes         *int
	MaxParallelTreeSyncChunkFetches             *int

	// snapshotting
	SnapshotInterval               *uint64
	MaxHistoricalSnapshotsRetained *uint64

	// blobs
	DeltaBlobOfferMinRequestToSameOracleInterval *time.Duration
	DeltaBlobOfferResponseTimeout                *time.Duration
	DeltaBlobBroadcastGrace                      *time.Duration
	DeltaBlobChunkMinRequestToSameOracleInterval *time.Duration
	DeltaBlobChunkResponseTimeout                *time.Duration
	BlobChunkBytes                               *int
}

// This function may be used in production. If you use this as part of multisig
// tooling,  make sure that the input parameters are identical across all
// signers.
// See [ocr3_1config.PublicConfig] for documentation of the arguments.
func ContractSetConfigArgsDeterministic(
	checkPublicConfigLevel CheckPublicConfigLevel,
	// The ephemeral secret key used to encrypt the shared secret
	ephemeralSk [curve25519.ScalarSize]byte,
	// A secret shared between all oracles, enabling them to derive pseudorandom
	// leaders and transmitters. This is a low-value secret. An adversary who
	// learns this secret can only determine which oracle will be
	// leader/transmitter ahead of time, that's it.
	sharedSecret [config.SharedSecretSize]byte,
	// Check out the ocr3_1config package for documentation on the following args
	oracles []confighelper.OracleIdentityExtra,
	f int,
	deltaProgress time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	rMax uint64,
	deltaStage time.Duration,
	s []int,
	reportingPluginConfig []byte,
	onchainConfig []byte,
	maxDurationInitialization time.Duration,
	warnDurationQuery time.Duration,
	warnDurationObservation time.Duration,
	warnDurationValidateObservation time.Duration,
	warnDurationObservationQuorum time.Duration,
	warnDurationStateTransition time.Duration,
	warnDurationCommitted time.Duration,
	maxDurationShouldAcceptAttestedReport time.Duration,
	maxDurationShouldTransmitAcceptedReport time.Duration,
	// Optional configuration parameters
	optionalConfig ContractSetConfigArgsOptionalConfig,
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

	sharedConfig := ocr3_1config.SharedConfig{
		ocr3_1config.PublicConfig{
			deltaProgress,
			optionalConfig.DeltaResend,
			optionalConfig.DeltaInitial,
			deltaRound,
			deltaGrace,
			optionalConfig.DeltaReportsPlusPrecursorRequest,
			deltaStage,

			// state sync
			optionalConfig.DeltaStateSyncSummaryInterval,

			// block sync
			optionalConfig.DeltaBlockSyncMinRequestToSameOracleInterval,
			optionalConfig.DeltaBlockSyncResponseTimeout,
			optionalConfig.MaxBlocksPerBlockSyncResponse,
			optionalConfig.MaxParallelRequestedBlocks,

			// tree sync
			optionalConfig.DeltaTreeSyncMinRequestToSameOracleInterval,
			optionalConfig.DeltaTreeSyncResponseTimeout,
			optionalConfig.MaxTreeSyncChunkKeys,
			optionalConfig.MaxTreeSyncChunkKeysPlusValuesBytes,
			optionalConfig.MaxParallelTreeSyncChunkFetches,

			// snapshotting
			optionalConfig.SnapshotInterval,
			optionalConfig.MaxHistoricalSnapshotsRetained,

			// blobs
			optionalConfig.DeltaBlobOfferMinRequestToSameOracleInterval,
			optionalConfig.DeltaBlobOfferResponseTimeout,
			optionalConfig.DeltaBlobBroadcastGrace,
			optionalConfig.DeltaBlobChunkMinRequestToSameOracleInterval,
			optionalConfig.DeltaBlobChunkResponseTimeout,
			optionalConfig.BlobChunkBytes,

			rMax,
			s,
			identities,
			reportingPluginConfig,
			maxDurationInitialization,
			warnDurationQuery,
			warnDurationObservation,
			warnDurationValidateObservation,
			warnDurationObservationQuorum,
			warnDurationStateTransition,
			warnDurationCommitted,
			maxDurationShouldAcceptAttestedReport,
			maxDurationShouldTransmitAcceptedReport,
			f,
			onchainConfig,
			types.ConfigDigest{},
		},
		&sharedSecret,
	}

	if err := ocr3_1config.CheckPublicConfig(checkPublicConfigLevel.shouldSkipInsaneForProductionChecks(), sharedConfig.PublicConfig); err != nil {
		return nil, nil, 0, nil, 0, nil, fmt.Errorf("CheckPublicConfig: %w", err)
	}

	return ocr3_1config.ContractSetConfigArgsFromSharedConfigDeterministic(
		sharedConfig,
		configEncryptionPublicKeys,
		&ephemeralSk,
	)
}
