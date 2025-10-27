package ocr3_1config

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/internal/byzquorum"
	"github.com/smartcontractkit/libocr/internal/util"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// PublicConfig is the configuration disseminated through the smart contract.
// It's public, because anybody can read it from the blockchain.
// The various parameters (e.g. Delta*, MaxDuration*) have some dependencies
// on each other, so be sure to consider the holistic impact of changes to them.
type PublicConfig struct {
	// If an epoch (driven by a leader) fails to achieve progress (generate a
	// report) after DeltaProgress, we enter a new epoch. This parameter must be
	// chosen carefully. If the duration is too short, we may keep prematurely
	// switching epochs without ever achieving any progress, resulting in a
	// liveness failure!
	DeltaProgress time.Duration
	// DeltaResend determines how often Pacemaker messages should be
	// resent, allowing oracles that had crashed and are recovering to rejoin
	// the protocol more quickly.
	DeltaResend *time.Duration
	// If no message from the leader has been received after the epoch start plus
	// DeltaInitial, we enter a new epoch. This parameter must be
	// chosen carefully. If the duration is too short, we may keep prematurely
	// switching epochs without ever achieving any progress, resulting in a
	// liveness failure!
	DeltaInitial *time.Duration
	// DeltaRound determines the minimal amount of time that should pass between
	// the start of outcome generation rounds. With OCR3 and higher versions (not OCR1!)
	// you can set this value very aggressively. Note that this only provides a lower
	// bound on the round interval; actual rounds might take longer.
	DeltaRound time.Duration
	// Once the leader of a outcome generation round has collected sufficiently
	// many observations, it will wait for DeltaGrace to pass to allow slower
	// oracles to still contribute an observation before moving on to generating
	// the report. Consequently, rounds driven by correct leaders will always
	// take at least DeltaGrace.
	DeltaGrace time.Duration
	// DeltaReportsPlusPrecursorRequest determines the duration between requests for
	// reports plus precursor after we have received f+1 signatures in the report
	// attestation protocol but are still missing the reports plus precursor
	// required for validating the report signatures.
	DeltaReportsPlusPrecursorRequest *time.Duration
	// DeltaStage determines the duration between stages of the transmission
	// protocol. In each stage, a certain number of oracles (determined by S)
	// will attempt to transmit, assuming that no other oracle has yet
	// successfully transmitted a report.
	DeltaStage time.Duration

	// === State Synchronization ===

	// DeltaStateSyncSummaryInterval defines how frequently an oracle
	// broadcasts a summary of its current state for synchronization purposes.
	DeltaStateSyncSummaryInterval *time.Duration

	// === Block Synchronization ===

	// DeltaBlockSyncMinRequestToSameOracleInterval specifies the minimum
	// duration between two consecutive block synchronization requests
	// sent to the same oracle.
	DeltaBlockSyncMinRequestToSameOracleInterval *time.Duration
	// DeltaBlockSyncResponseTimeout specifies the maximum time to wait
	// for a response to a specific block synchronization request.
	// If no response is received within this duration,
	// the protocol retries with another oracle.
	DeltaBlockSyncResponseTimeout *time.Duration
	// MaxBlocksPerBlockSyncResponse defines the maximum number of blocks
	// that can be included in a single block synchronization response.
	MaxBlocksPerBlockSyncResponse *int
	// MaxParallelRequestedBlocks upper bounds the number of blocks being
	// requested in parallel. Multiple blocks might be fetched as part of a one
	// request.
	MaxParallelRequestedBlocks *uint64

	// === Tree Synchronization ===

	// DeltaTreeSyncMinRequestToSameOracleInterval specifies the minimum
	// duration between two consecutive requests for tree synchronization
	// sent to the same oracle.
	DeltaTreeSyncMinRequestToSameOracleInterval *time.Duration
	// DeltaTreeSyncResponseTimeout specifies the maximum amount of time to
	// wait for a response to a specific tree synchronization request.
	// If no response is received within this duration,
	// the protocol retries with another oracle.
	DeltaTreeSyncResponseTimeout *time.Duration
	// MaxTreeSyncChunkKeys defines the maximum number of key-value pairs
	// that an oracle includes in a single tree synchronization response chunk.
	MaxTreeSyncChunkKeys *int
	// MaxTreeSyncChunkKeysPlusValuesBytes defines the maximum combined
	// size (in bytes) of all keys and values in a single tree
	// synchronization response chunk.
	// The protocol ensures that each chunk includes as many key-value pairs
	// as possible without exceeding either this byte-size limit or
	// MaxTreeSyncChunkKeys.
	// A chunk must always fit at least one maximally sized (using maxmax)
	// key-value pair.
	MaxTreeSyncChunkKeysPlusValuesBytes *int
	// MaxParallelTreeSyncChunkFetches defines the maximum number of tree
	// synchronization requests that can be performed in parallel.
	MaxParallelTreeSyncChunkFetches *int

	// === Snapshotting ===
	// SnapshotInterval is defined such that the committed sequence number of
	// any snapshot must be a multiple of SnapshotInterval. Decreasing this
	// value increases the max historical snapshots retained.
	SnapshotInterval *uint64
	// MaxHistoricalSnapshotsRetained defines how many complete historical
	// snapshots are retained. Retained snapshots enable other oracles to
	// synchronize against the committed state from previous snapshot sequence
	// numbers. All blocks from the highest block of the earliest retained
	// snapshot onward will be kept available for synchronization purposes.
	MaxHistoricalSnapshotsRetained *uint64

	// === Blob Synchronization ===
	//
	// DeltaBlobOfferMinRequestToSameOracleInterval defines the minimum
	// duration between two consecutive blob offer requests sent to the same
	// oracle.
	DeltaBlobOfferMinRequestToSameOracleInterval *time.Duration
	// DeltaBlobOfferResponseTimeout specifies the maximum duration to wait
	// for a response to a blob offer before resending the blob offer.
	DeltaBlobOfferResponseTimeout *time.Duration
	// DeltaBlobBroadcastGrace defines the additional grace period to wait
	// after receiving the minimum number of accepting blob offer responses.
	// This allows more oracles a final opportunity to be included in the
	// availability certificate.
	DeltaBlobBroadcastGrace *time.Duration
	// DeltaBlobChunkMinRequestToSameOracleInterval defines the minimum
	// duration between two consecutive blob chunk requests sent to the same
	// oracle.
	DeltaBlobChunkMinRequestToSameOracleInterval *time.Duration
	// DeltaBlobChunkResponseTimeout specifies the maximum duration to wait
	// for a blob chunk response. If no response is received within this
	// time, the protocol retries with another oracle.
	DeltaBlobChunkResponseTimeout *time.Duration
	// BlobChunkBytes defines the size of blob chunks in bytes.
	BlobChunkBytes *int

	// The maximum number of rounds during an epoch.
	RMax uint64
	// S is the transmission schedule. For example, S = [1,2,3] indicates that
	// in the first stage of transmission one oracle will attempt to transmit,
	// in the second stage two more will attempt to transmit (if in their view
	// the first stage didn't succeed), and in the third stage three more will
	// attempt to transmit (if in their view the first and second stage didn't
	// succeed).
	//
	// sum(S) should equal n.
	S []int
	// Identities (i.e. public keys) of the oracles participating in this
	// protocol instance.
	OracleIdentities []config.OracleIdentity

	// Binary blob containing configuration passed through to the
	// ReportingPlugin.
	ReportingPluginConfig []byte

	// MaxDurationX is the maximum duration a ReportingPlugin should spend
	// performing X. Reasonable values for these will be specific to each
	// ReportingPlugin. Be sure to not set these too short, or the corresponding
	// ReportingPlugin function may always time out. The logic for
	// WarnDurationQuery and WarnDurationObservation has changed since these
	// values were first introduced. Unlike the other MaxDurationX values,
	// exceeding WarnDurationQuery and WarnDurationObservation will only cause
	// warnings to be logged, but will *not* cause X to time out.
	//
	// These values are passed to the ReportingPlugin during initialization.
	// Consequently, the ReportingPlugin may exhibit specific behaviors based on
	// these values. For instance, the MercuryReportingPlugin uses
	// WarnDurationObservation to set context timeouts.
	MaxDurationInitialization               time.Duration // Context deadline passed to NewReportingPlugin.
	WarnDurationQuery                       time.Duration // If the Query function takes longer than this, a warning will be logged.
	WarnDurationObservation                 time.Duration // If the Observation function takes longer than this, a warning will be logged.
	WarnDurationValidateObservation         time.Duration // If the ValidateObservation function takes longer than this, a warning will be logged.
	WarnDurationObservationQuorum           time.Duration // If the ObservationQuorum function takes longer than this, a warning will be logged.
	WarnDurationStateTransition             time.Duration // If the StateTransition function takes longer than this, a warning will be logged.
	WarnDurationCommitted                   time.Duration // If the Committed function takes longer than this, a warning will be logged.
	MaxDurationShouldAcceptAttestedReport   time.Duration // Context deadline passed to ShouldAcceptAttestedReport.
	MaxDurationShouldTransmitAcceptedReport time.Duration // Context deadline passed to ShouldTransmitAcceptedReport.

	// The maximum number of oracles that are assumed to be faulty while the
	// protocol can retain liveness and safety. Unless you really know what
	// you’re doing, be sure to set this to floor((n-1)/3) where n is the total
	// number of oracles.
	F int

	// Binary blob containing configuration passed through to the
	// ReportingPlugin, and also available to the contract. (Unlike
	// ReportingPluginConfig which is only available offchain.)
	OnchainConfig []byte

	ConfigDigest types.ConfigDigest
}

// N is the number of oracles participating in the protocol
func (c *PublicConfig) N() int {
	return len(c.OracleIdentities)
}

func (c *PublicConfig) ByzQuorumSize() int {
	return byzquorum.Size(c.N(), c.F)
}

func (c *PublicConfig) GetDeltaResend() time.Duration {
	return util.NilCoalesce(c.DeltaResend, DefaultDeltaResend)
}

func (c *PublicConfig) GetDeltaInitial() time.Duration {
	return util.NilCoalesce(c.DeltaInitial, DefaultDeltaInitial())
}

func (c *PublicConfig) GetDeltaReportsPlusPrecursorRequest() time.Duration {
	return util.NilCoalesce(c.DeltaReportsPlusPrecursorRequest, DefaultDeltaReportsPlusPrecursorRequest())
}

func (c *PublicConfig) GetDeltaStateSyncSummaryInterval() time.Duration {
	return util.NilCoalesce(c.DeltaStateSyncSummaryInterval, DefaultDeltaStateSyncSummaryInterval)
}

func (c *PublicConfig) GetDeltaBlockSyncMinRequestToSameOracleInterval() time.Duration {
	return util.NilCoalesce(c.DeltaBlockSyncMinRequestToSameOracleInterval, DefaultDeltaBlockSyncMinRequestToSameOracleInterval)
}

func (c *PublicConfig) GetDeltaBlockSyncResponseTimeout() time.Duration {
	return util.NilCoalesce(c.DeltaBlockSyncResponseTimeout, DefaultDeltaBlockSyncResponseTimeout())
}

func (c *PublicConfig) GetMaxBlocksPerBlockSyncResponse() int {
	return util.NilCoalesce(c.MaxBlocksPerBlockSyncResponse, DefaultMaxBlocksPerBlockSyncResponse)
}

func (c *PublicConfig) GetMaxParallelRequestedBlocks() uint64 {
	return util.NilCoalesce(c.MaxParallelRequestedBlocks, DefaultMaxParallelRequestedBlocks)
}

func (c *PublicConfig) GetDeltaTreeSyncMinRequestToSameOracleInterval() time.Duration {
	return util.NilCoalesce(c.DeltaTreeSyncMinRequestToSameOracleInterval, DefaultDeltaTreeSyncMinRequestToSameOracleInterval)
}

func (c *PublicConfig) GetDeltaTreeSyncResponseTimeout() time.Duration {
	return util.NilCoalesce(c.DeltaTreeSyncResponseTimeout, DefaultDeltaTreeSyncResponseTimeout())
}

func (c *PublicConfig) GetMaxTreeSyncChunkKeys() int {
	return util.NilCoalesce(c.MaxTreeSyncChunkKeys, DefaultMaxTreeSyncChunkKeys)
}

func (c *PublicConfig) GetMaxTreeSyncChunkKeysPlusValuesBytes() int {
	return util.NilCoalesce(c.MaxTreeSyncChunkKeysPlusValuesBytes, DefaultMaxTreeSyncChunkKeysPlusValuesBytes)
}

func (c *PublicConfig) GetMaxParallelTreeSyncChunkFetches() int {
	return util.NilCoalesce(c.MaxParallelTreeSyncChunkFetches, DefaultMaxParallelTreeSyncChunkFetches)
}

func (c *PublicConfig) GetSnapshotInterval() uint64 {
	return util.NilCoalesce(c.SnapshotInterval, DefaultSnapshotInterval)
}

func (c *PublicConfig) GetMaxHistoricalSnapshotsRetained() uint64 {
	return util.NilCoalesce(c.MaxHistoricalSnapshotsRetained, DefaultMaxHistoricalSnapshotsRetained)
}

func (c *PublicConfig) GetDeltaBlobOfferMinRequestToSameOracleInterval() time.Duration {
	return util.NilCoalesce(c.DeltaBlobOfferMinRequestToSameOracleInterval, DefaultDeltaBlobOfferMinRequestToSameOracleInterval)
}

func (c *PublicConfig) GetDeltaBlobOfferResponseTimeout() time.Duration {
	return util.NilCoalesce(c.DeltaBlobOfferResponseTimeout, DefaultDeltaBlobOfferResponseTimeout)
}

func (c *PublicConfig) GetDeltaBlobBroadcastGrace() time.Duration {
	return util.NilCoalesce(c.DeltaBlobBroadcastGrace, DefaultDeltaBlobBroadcastGrace)
}

func (c *PublicConfig) GetDeltaBlobChunkMinRequestToSameOracleInterval() time.Duration {
	return util.NilCoalesce(c.DeltaBlobChunkMinRequestToSameOracleInterval, DefaultDeltaBlobChunkMinRequestToSameOracleInterval)
}

func (c *PublicConfig) GetDeltaBlobChunkResponseTimeout() time.Duration {
	return util.NilCoalesce(c.DeltaBlobChunkResponseTimeout, DefaultDeltaBlobChunkResponseTimeout())
}

func (c *PublicConfig) GetBlobChunkBytes() int {
	return util.NilCoalesce(c.BlobChunkBytes, DefaultBlobChunkBytes)
}

// The minimum interval between round starts.
// This is not a guaranteed lower bound. For example, a malicious leader could
// violate this bound.
func (c *PublicConfig) MinRoundInterval() time.Duration {
	if c.DeltaRound > c.DeltaGrace {
		return c.DeltaRound
	}
	return c.DeltaGrace
}

func (c *PublicConfig) CheckParameterBounds() error {
	if c.F < 0 || c.F > math.MaxUint8 {
		return errors.Errorf("number of potentially faulty oracles must fit in 8 bits.")
	}
	return nil
}

func PublicConfigFromContractConfig(skipInsaneForProductionChecks bool, change types.ContractConfig) (PublicConfig, error) {
	pubcon, _, err := publicConfigFromContractConfig(skipInsaneForProductionChecks, change)
	return pubcon, err
}

func publicConfigFromContractConfig(skipInsaneForProductionChecks bool, change types.ContractConfig) (PublicConfig, config.SharedSecretEncryptions, error) {
	if change.OffchainConfigVersion != config.OCR3_1OffchainConfigVersion {
		return PublicConfig{}, config.SharedSecretEncryptions{}, fmt.Errorf("unsuppported OffchainConfigVersion %v, supported OffchainConfigVersion is %v", change.OffchainConfigVersion, config.OCR3_1OffchainConfigVersion)
	}

	oc, err := deserializeOffchainConfig(change.OffchainConfig)
	if err != nil {
		return PublicConfig{}, config.SharedSecretEncryptions{}, err
	}

	if err := checkIdentityListsHaveNoDuplicates(change, oc); err != nil {
		return PublicConfig{}, config.SharedSecretEncryptions{}, err
	}

	// must check that all lists have the same length, or bad input could crash
	// the following for loop.
	if err := checkIdentityListsHaveTheSameLength(change, oc); err != nil {
		return PublicConfig{}, config.SharedSecretEncryptions{}, err
	}

	identities := []config.OracleIdentity{}
	for i := range change.Signers {
		identities = append(identities, config.OracleIdentity{
			oc.OffchainPublicKeys[i],
			types.OnchainPublicKey(change.Signers[i][:]),
			oc.PeerIDs[i],
			change.Transmitters[i],
		})
	}

	cfg := PublicConfig{
		oc.DeltaProgress,
		oc.DeltaResend,
		oc.DeltaInitial,
		oc.DeltaRound,
		oc.DeltaGrace,
		oc.DeltaReportsPlusPrecursorRequest,
		oc.DeltaStage,

		// state sync
		oc.DeltaStateSyncSummaryInterval,

		// block sync
		oc.DeltaBlockSyncMinRequestToSameOracleInterval,
		oc.DeltaBlockSyncResponseTimeout,
		oc.MaxBlocksPerBlockSyncResponse,
		oc.MaxParallelRequestedBlocks,

		// tree sync
		oc.DeltaTreeSyncMinRequestToSameOracleInterval,
		oc.DeltaTreeSyncResponseTimeout,
		oc.MaxTreeSyncChunkKeys,
		oc.MaxTreeSyncChunkKeysPlusValuesBytes,
		oc.MaxParallelTreeSyncChunkFetches,

		// snapshotting
		oc.SnapshotInterval,
		oc.MaxHistoricalSnapshotsRetained,

		// blobs
		oc.DeltaBlobOfferMinRequestToSameOracleInterval,
		oc.DeltaBlobOfferResponseTimeout,
		oc.DeltaBlobBroadcastGrace,
		oc.DeltaBlobChunkMinRequestToSameOracleInterval,
		oc.DeltaBlobChunkResponseTimeout,
		oc.BlobChunkBytes,

		oc.RMax,
		oc.S,
		identities,
		oc.ReportingPluginConfig,
		oc.MaxDurationInitialization,
		oc.WarnDurationQuery,
		oc.WarnDurationObservation,
		oc.WarnDurationValidateObservation,
		oc.WarnDurationObservationQuorum,
		oc.WarnDurationStateTransition,
		oc.WarnDurationCommitted,
		oc.MaxDurationShouldAcceptAttestedReport,
		oc.MaxDurationShouldTransmitAcceptedReport,

		int(change.F),
		change.OnchainConfig,
		change.ConfigDigest,
	}

	if err := CheckPublicConfig(skipInsaneForProductionChecks, cfg); err != nil {
		return PublicConfig{}, config.SharedSecretEncryptions{}, err
	}

	return cfg, oc.SharedSecretEncryptions, nil
}

func CheckPublicConfig(skipInsaneForProductionChecks bool, publicConfig PublicConfig) error {
	if err := checkPublicConfigParameters(publicConfig); err != nil {
		return fmt.Errorf("checkPublicConfigParameters: %w", err)
	}

	if err := checkIdentityListsHaveNoDuplicatesPublicConfig(publicConfig); err != nil {
		return fmt.Errorf("checkIdentityListsHaveNoDuplicatesPublicConfig: %w", err)
	}

	if !skipInsaneForProductionChecks {
		if err := checkNotInsaneForProduction(publicConfig); err != nil {
			return fmt.Errorf("checkNotInsaneForProduction: %w", err)
		}
	}
	return nil
}

func checkIdentityListsHaveNoDuplicates(change types.ContractConfig, oc offchainConfig) error {
	// inefficient, but it doesn't matter
	for i := range change.Signers {
		for j := range change.Signers {
			if i != j && bytes.Equal(change.Signers[i], change.Signers[j]) {
				return fmt.Errorf("%v-th and %v-th signer are identical: %x", i, j, change.Signers[i])
			}
		}
	}

	{
		uniquePeerIDs := map[string]struct{}{}
		for _, peerID := range oc.PeerIDs {
			if _, ok := uniquePeerIDs[peerID]; ok {
				return fmt.Errorf("duplicate PeerID '%v'", peerID)
			}
			uniquePeerIDs[peerID] = struct{}{}
		}
	}

	{
		uniqueOffchainPublicKeys := map[types.OffchainPublicKey]struct{}{}
		for _, ocpk := range oc.OffchainPublicKeys {
			if _, ok := uniqueOffchainPublicKeys[ocpk]; ok {
				return fmt.Errorf("duplicate OffchainPublicKey %x", ocpk)
			}
			uniqueOffchainPublicKeys[ocpk] = struct{}{}
		}
	}

	{
		// this isn't strictly necessary, but since we don't intend to run
		// with duplicate transmitters at this time, we might as well check
		uniqueTransmitters := map[types.Account]struct{}{}
		for _, transmitter := range change.Transmitters {
			if _, ok := uniqueTransmitters[transmitter]; ok {
				return fmt.Errorf("duplicate transmitter '%v'", transmitter)
			}
			uniqueTransmitters[transmitter] = struct{}{}
		}
	}

	// no point in checking SharedSecretEncryptions for uniqueness

	return nil
}

func checkIdentityListsHaveNoDuplicatesPublicConfig(publicConfig PublicConfig) error {
	// inefficient, but it doesn't matter
	for i := range publicConfig.OracleIdentities {
		for j := range publicConfig.OracleIdentities {
			if i != j && bytes.Equal(publicConfig.OracleIdentities[i].OnchainPublicKey, publicConfig.OracleIdentities[j].OnchainPublicKey) {
				return fmt.Errorf("%v-th and %v-th OnchainPublicKey are identical: %x", i, j, publicConfig.OracleIdentities[i].OnchainPublicKey)
			}
		}
	}

	{
		uniquePeerIDs := map[string]struct{}{}
		for _, oid := range publicConfig.OracleIdentities {
			if _, ok := uniquePeerIDs[oid.PeerID]; ok {
				return fmt.Errorf("duplicate PeerID '%v'", oid.PeerID)
			}
			uniquePeerIDs[oid.PeerID] = struct{}{}
		}
	}

	{
		uniqueOffchainPublicKeys := map[types.OffchainPublicKey]struct{}{}
		for _, oid := range publicConfig.OracleIdentities {
			if _, ok := uniqueOffchainPublicKeys[oid.OffchainPublicKey]; ok {
				return fmt.Errorf("duplicate OffchainPublicKey %x", oid.OffchainPublicKey)
			}
			uniqueOffchainPublicKeys[oid.OffchainPublicKey] = struct{}{}
		}
	}

	{
		// this isn't strictly necessary, but since we don't intend to run
		// with duplicate transmitters at this time, we might as well check
		uniqueTransmitters := map[types.Account]struct{}{}
		for _, oid := range publicConfig.OracleIdentities {
			if _, ok := uniqueTransmitters[oid.TransmitAccount]; ok {
				return fmt.Errorf("duplicate TransmitAccount '%v'", oid.TransmitAccount)
			}
			uniqueTransmitters[oid.TransmitAccount] = struct{}{}
		}
	}

	return nil
}

func checkIdentityListsHaveTheSameLength(
	change types.ContractConfig, oc offchainConfig,
) error {
	expectedLength := len(change.Signers)
	errorMsg := "%s list must have same length as onchain signers list: %d ≠ " +
		strconv.Itoa(expectedLength)
	for _, identityList := range []struct {
		length int
		name   string
	}{
		{len(oc.PeerIDs) /*                       */, "peer ids"},
		{len(oc.OffchainPublicKeys) /*            */, "offchain public keys"},
		{len(change.Transmitters) /*              */, "transmitters"},
		{len(oc.SharedSecretEncryptions.Encryptions), "shared-secret encryptions"},
	} {
		if identityList.length != expectedLength {
			return errors.Errorf(errorMsg, identityList.name, identityList.length)
		}
	}
	return nil
}

// Sanity check on parameters:
// (1) violations of fundamental constraints like 3*f<n;
// (2) configurations that would trivially exhaust all of a node's resources;
// (3) (some) simple mistakes

func checkPublicConfigParameters(cfg PublicConfig) error {
	/////////////////////////////////////////////////////////////////
	// Be sure to think about changes to other tooling that need to
	// be made when you change this function!
	/////////////////////////////////////////////////////////////////

	if !(0 <= cfg.F && cfg.F*3 < cfg.N()) {
		return fmt.Errorf("F (%v) must be non-negative and less than N/3 (N = %v)",
			cfg.F, cfg.N())
	}

	if !(cfg.N() <= types.MaxOracles) {
		return fmt.Errorf("N (%v) must be less than or equal MaxOracles (%v)",
			cfg.N(), types.MaxOracles)
	}

	if !(0 < cfg.DeltaProgress) {
		return fmt.Errorf("DeltaProgress (%v) must be positive", cfg.DeltaProgress)
	}

	if !(0 < cfg.GetDeltaResend()) {
		return fmt.Errorf("DeltaResend (%v) must be positive", cfg.GetDeltaResend())
	}

	if !(0 < cfg.GetDeltaInitial()) {
		return fmt.Errorf("DeltaInitial (%v) must be positive", cfg.GetDeltaInitial())
	}

	if !(0 <= cfg.DeltaRound) {
		return fmt.Errorf("DeltaRound (%v) must be non-negative", cfg.DeltaRound)
	}

	if !(0 <= cfg.DeltaGrace) {
		return fmt.Errorf("DeltaGrace (%v) must be non-negative",
			cfg.DeltaGrace)
	}

	if !(0 < cfg.GetDeltaReportsPlusPrecursorRequest()) {
		return fmt.Errorf("DeltaReportsPlusPrecursorRequest (%v) must be positive", cfg.GetDeltaReportsPlusPrecursorRequest())
	}

	if !(0 <= cfg.DeltaStage) {
		return fmt.Errorf("DeltaStage (%v) must be non-negative", cfg.DeltaStage)
	}

	// state sync
	if !(0 < cfg.GetDeltaStateSyncSummaryInterval()) {
		return fmt.Errorf("DeltaStateSyncSummaryInterval (%v) must be positive", cfg.GetDeltaStateSyncSummaryInterval())
	}

	// block sync
	if !(0 < cfg.GetDeltaBlockSyncMinRequestToSameOracleInterval()) {
		return fmt.Errorf("DeltaBlockSyncMinRequestToSameOracleInterval (%v) must be positive", cfg.GetDeltaBlockSyncMinRequestToSameOracleInterval())
	}
	if !(0 < cfg.GetDeltaBlockSyncResponseTimeout()) {
		return fmt.Errorf("DeltaBlockSyncResponseTimeout (%v) must be positive", cfg.GetDeltaBlockSyncResponseTimeout())
	}
	if !(0 < cfg.GetMaxBlocksPerBlockSyncResponse()) {
		return fmt.Errorf("MaxBlocksPerBlockSyncResponse (%v) must be positive", cfg.GetMaxBlocksPerBlockSyncResponse())
	}
	if !(0 < cfg.GetMaxParallelRequestedBlocks()) {
		return fmt.Errorf("MaxParallelRequestedBlocks (%v) must be positive", cfg.GetMaxParallelRequestedBlocks())
	}

	// tree sync
	if !(0 < cfg.GetDeltaTreeSyncMinRequestToSameOracleInterval()) {
		return fmt.Errorf("DeltaTreeSyncMinRequestToSameOracleInterval (%v) must be positive", cfg.GetDeltaTreeSyncMinRequestToSameOracleInterval())
	}
	if !(0 < cfg.GetDeltaTreeSyncResponseTimeout()) {
		return fmt.Errorf("DeltaTreeSyncResponseTimeout (%v) must be positive", cfg.GetDeltaTreeSyncResponseTimeout())
	}
	if !(0 < cfg.GetMaxTreeSyncChunkKeys()) {
		return fmt.Errorf("MaxTreeSyncChunkKeys (%v) must be positive", cfg.GetMaxTreeSyncChunkKeys())
	}

	if !(ocr3_1types.MaxMaxKeyValueKeyBytes+ocr3_1types.MaxMaxKeyValueValueBytes <= cfg.GetMaxTreeSyncChunkKeysPlusValuesBytes()) {
		return fmt.Errorf("MaxTreeSyncChunkKeysPlusValuesBytes (%v) must be greater than or equal to MaxMaxKeyValueKeyBytes (%v) + MaxMaxKeyValueValueBytes (%v)", cfg.GetMaxTreeSyncChunkKeysPlusValuesBytes(), ocr3_1types.MaxMaxKeyValueKeyBytes, ocr3_1types.MaxMaxKeyValueValueBytes)
	}
	if !(0 < cfg.GetMaxParallelTreeSyncChunkFetches()) {
		return fmt.Errorf("MaxParallelTreeSyncChunkFetches (%v) must be positive", cfg.GetMaxParallelTreeSyncChunkFetches())
	}

	// snapshotting
	if !(0 < cfg.GetSnapshotInterval()) {
		return fmt.Errorf("SnapshotInterval (%v) must be positive", cfg.GetSnapshotInterval())
	}
	if !(0 < cfg.GetMaxHistoricalSnapshotsRetained()) {
		return fmt.Errorf("MaxHistoricalSnapshotsRetained (%v) must be positive", cfg.GetMaxHistoricalSnapshotsRetained())
	}

	// blobs
	if !(0 < cfg.GetDeltaBlobOfferMinRequestToSameOracleInterval()) {
		return fmt.Errorf("DeltaBlobOfferMinRequestToSameOracleInterval (%v) must be positive", cfg.GetDeltaBlobOfferMinRequestToSameOracleInterval())
	}
	if !(0 < cfg.GetDeltaBlobOfferResponseTimeout()) {
		return fmt.Errorf("DeltaBlobOfferResponseTimeout (%v) must be positive", cfg.GetDeltaBlobOfferResponseTimeout())
	}
	if !(0 <= cfg.GetDeltaBlobBroadcastGrace()) {
		return fmt.Errorf("DeltaBlobBroadcastGrace (%v) must be non-negative", cfg.GetDeltaBlobBroadcastGrace())
	}
	if !(0 < cfg.GetDeltaBlobChunkMinRequestToSameOracleInterval()) {
		return fmt.Errorf("DeltaBlobChunkMinRequestToSameOracleInterval (%v) must be positive", cfg.GetDeltaBlobChunkMinRequestToSameOracleInterval())
	}
	if !(0 < cfg.GetDeltaBlobChunkResponseTimeout()) {
		return fmt.Errorf("DeltaBlobChunkResponseTimeout (%v) must be positive", cfg.GetDeltaBlobChunkResponseTimeout())
	}

	if !(0 < cfg.GetBlobChunkBytes()) {
		return fmt.Errorf("BlobChunkBytes (%v) must be positive", cfg.GetBlobChunkBytes())
	}

	if !(0 < cfg.MaxDurationInitialization) {
		return fmt.Errorf("MaxDurationInitialization (%v) must be positive", cfg.MaxDurationInitialization)
	}

	if !(0 < cfg.WarnDurationQuery) {
		return fmt.Errorf("WarnDurationQuery (%v) must be positive", cfg.WarnDurationQuery)
	}

	if !(0 < cfg.WarnDurationObservation) {
		return fmt.Errorf("WarnDurationObservation (%v) must be positive", cfg.WarnDurationObservation)
	}

	if !(0 < cfg.WarnDurationValidateObservation) {
		return fmt.Errorf("WarnDurationValidateObservation (%v) must be positive", cfg.WarnDurationValidateObservation)
	}

	if !(0 < cfg.WarnDurationObservationQuorum) {
		return fmt.Errorf("WarnDurationObservationQuorum (%v) must be positive", cfg.WarnDurationObservationQuorum)
	}

	if !(0 < cfg.WarnDurationStateTransition) {
		return fmt.Errorf("WarnDurationStateTransition (%v) must be positive", cfg.WarnDurationStateTransition)
	}

	if !(0 < cfg.WarnDurationCommitted) {
		return fmt.Errorf("WarnDurationCommitted (%v) must be positive", cfg.WarnDurationCommitted)
	}

	if !(0 < cfg.MaxDurationShouldAcceptAttestedReport) {
		return fmt.Errorf("MaxDurationShouldAcceptAttestedReport (%v) must be positive", cfg.MaxDurationShouldAcceptAttestedReport)
	}

	if !(0 < cfg.MaxDurationShouldTransmitAcceptedReport) {
		return fmt.Errorf("MaxDurationShouldTransmitAcceptedReport (%v) must be positive", cfg.MaxDurationShouldTransmitAcceptedReport)
	}

	if !(cfg.DeltaRound < cfg.DeltaProgress) {
		return fmt.Errorf("DeltaRound (%v) must be less than DeltaProgress (%v)",
			cfg.DeltaRound, cfg.DeltaProgress)
	}

	if !(cfg.DeltaGrace < cfg.DeltaProgress) {
		return fmt.Errorf("DeltaGrace (%v) must be less than DeltaProgress (%v)",
			cfg.DeltaGrace, cfg.DeltaProgress)
	}

	// We cannot easily add a similar check for the MaxDuration variables used
	// in the transmission protocol (MaxDurationShouldAcceptAttestedReport,
	// MaxDurationShouldTransmitAcceptedReport), because we don't know how often
	// they will be triggered. But if we assume that there is one transmission
	// for each round, we should have MaxDurationShouldAcceptAttestedReport +
	// MaxDurationShouldTransmitAcceptedReport < round duration.

	if !(0 < cfg.RMax) {
		return fmt.Errorf("RMax (%v) must be greater than zero", cfg.RMax)
	}

	// This prevents possible overflows adding up the elements of S. We should never
	// hit this.
	if !(len(cfg.S) < 1000) {
		return fmt.Errorf("len(S) (%v) must be less than 1000", len(cfg.S))
	}

	for i, s := range cfg.S {
		if !(0 <= s && s <= types.MaxOracles) {
			return fmt.Errorf("S[%v] (%v) must be between 0 and types.MaxOracles (%v)", i, s, types.MaxOracles)
		}
	}

	return nil
}

func checkNotInsaneForProduction(cfg PublicConfig) error {
	// Sending messages related to epoch changes and missing certified commits
	// shouldn't be necessary in any realistic WAN deployment and could cause
	// resource exhaustion
	const safeInterval = 100 * time.Millisecond
	if !(safeInterval <= cfg.DeltaProgress) {
		return fmt.Errorf("DeltaProgress (%v) is set below the resource exhaustion safe interval (%v)", cfg.DeltaProgress, safeInterval)
	}
	if !(safeInterval <= cfg.GetDeltaResend()) {
		return fmt.Errorf("DeltaResend (%v) is set below the resource exhaustion safe interval (%v)", cfg.GetDeltaResend(), safeInterval)
	}
	if !(safeInterval <= cfg.GetDeltaInitial()) {
		return fmt.Errorf("DeltaInitial (%v) is set below the resource exhaustion safe interval (%v)", cfg.GetDeltaInitial(), safeInterval)
	}
	const safeRoundInterval = 10 * time.Millisecond
	if !(safeRoundInterval <= cfg.DeltaRound) {
		return fmt.Errorf("DeltaRound (%v) is set below the safe round interval (%v)", cfg.DeltaRound, safeRoundInterval)
	}
	if !(safeRoundInterval <= cfg.MinRoundInterval()) {
		return fmt.Errorf("MinRoundInterval (%v) is set below the safe round interval (%v)", cfg.MinRoundInterval(), safeRoundInterval)
	}
	{
		prodNs := new(big.Int).SetInt64(1)
		prodNs = prodNs.Mul(prodNs, new(big.Int).SetUint64(cfg.GetSnapshotInterval()))
		prodNs = prodNs.Mul(prodNs, new(big.Int).SetUint64(cfg.GetMaxHistoricalSnapshotsRetained()))
		prodNs = prodNs.Mul(prodNs, new(big.Int).SetInt64(int64(cfg.MinRoundInterval())))
		oneHourNs := new(big.Int).SetInt64(int64(time.Hour))
		if !(prodNs.Cmp(oneHourNs) >= 0) {
			return fmt.Errorf("SnapshotInterval (%v) * MaxHistoricalSnapshotsRetained (%v) * MinRoundInterval (%v) must be greater than or equal to 1 hour", cfg.GetSnapshotInterval(), cfg.GetMaxHistoricalSnapshotsRetained(), cfg.MinRoundInterval())
		}
		sevenDaysNs := new(big.Int).SetInt64(int64(7 * 24 * time.Hour))
		if !(prodNs.Cmp(sevenDaysNs) <= 0) {
			return fmt.Errorf("SnapshotInterval (%v) * MaxHistoricalSnapshotsRetained (%v) * MinRoundInterval (%v) must be less than or equal to 7 days", cfg.GetSnapshotInterval(), cfg.GetMaxHistoricalSnapshotsRetained(), cfg.MinRoundInterval())
		}
	}
	const maxMaxHistoricalSnapshotsRetained = 1_000
	if !(cfg.GetMaxHistoricalSnapshotsRetained() <= maxMaxHistoricalSnapshotsRetained) {
		return fmt.Errorf("MaxHistoricalSnapshotsRetained (%v) must be less than or equal to %v", cfg.GetMaxHistoricalSnapshotsRetained(), maxMaxHistoricalSnapshotsRetained)
	}

	if !(safeInterval <= cfg.GetDeltaStateSyncSummaryInterval()) {
		return fmt.Errorf("DeltaStateSyncSummaryInterval (%v) is set below the resource exhaustion safe interval (%v)", cfg.GetDeltaStateSyncSummaryInterval(), safeInterval)
	}
	const safeRequestInterval = 10 * time.Millisecond
	// request intervals
	if !(safeRequestInterval <= cfg.GetDeltaReportsPlusPrecursorRequest()) {
		return fmt.Errorf("DeltaReportsPlusPrecursorRequest (%v) is set below the safe request interval (%v)", cfg.GetDeltaReportsPlusPrecursorRequest(), safeRequestInterval)
	}
	if !(safeRequestInterval <= cfg.GetDeltaBlockSyncMinRequestToSameOracleInterval()) {
		return fmt.Errorf("DeltaBlockSyncMinRequestToSameOracleInterval (%v) is set below the safe request interval (%v)", cfg.GetDeltaBlockSyncMinRequestToSameOracleInterval(), safeRequestInterval)
	}
	if !(safeRequestInterval <= cfg.GetDeltaTreeSyncMinRequestToSameOracleInterval()) {
		return fmt.Errorf("DeltaTreeSyncMinRequestToSameOracleInterval (%v) is set below the safe request interval (%v)", cfg.GetDeltaTreeSyncMinRequestToSameOracleInterval(), safeRequestInterval)
	}
	if !(safeRequestInterval <= cfg.GetDeltaBlobChunkMinRequestToSameOracleInterval()) {
		return fmt.Errorf("DeltaBlobChunkMinRequestToSameOracleInterval (%v) is set below the safe request interval (%v)", cfg.GetDeltaBlobChunkMinRequestToSameOracleInterval(), safeRequestInterval)
	}

	// response timeouts
	if !(safeInterval <= cfg.GetDeltaBlockSyncResponseTimeout()) {
		return fmt.Errorf("DeltaBlockSyncResponseTimeout (%v) is set below the resource exhaustion safe interval (%v)", cfg.GetDeltaBlockSyncResponseTimeout(), safeInterval)
	}
	if !(safeInterval <= cfg.GetDeltaTreeSyncResponseTimeout()) {
		return fmt.Errorf("DeltaTreeSyncResponseTimeout (%v) is set below the resource exhaustion safe interval (%v)", cfg.GetDeltaTreeSyncResponseTimeout(), safeInterval)
	}
	if !(safeInterval <= cfg.GetDeltaBlobChunkResponseTimeout()) {
		return fmt.Errorf("DeltaBlobChunkResponseTimeout (%v) is set below the resource exhaustion safe interval (%v)", cfg.GetDeltaBlobChunkResponseTimeout(), safeInterval)
	}

	if !(cfg.GetMaxBlocksPerBlockSyncResponse() <= MaxMaxBlocksPerBlockSyncResponse) {
		return fmt.Errorf("MaxBlocksPerBlockSyncResponse (%v) must be less than or equal to %v", cfg.GetMaxBlocksPerBlockSyncResponse(), MaxMaxBlocksPerBlockSyncResponse)
	}
	if !(cfg.GetMaxTreeSyncChunkKeys() <= MaxMaxTreeSyncChunkKeys) {
		return fmt.Errorf("MaxTreeSyncChunkKeys (%v) must be less than or equal to %v", cfg.GetMaxTreeSyncChunkKeys(), MaxMaxTreeSyncChunkKeys)
	}
	if !(cfg.GetMaxTreeSyncChunkKeysPlusValuesBytes() <= MaxMaxTreeSyncChunkKeysPlusValuesBytes) {
		return fmt.Errorf("MaxTreeSyncChunkKeysPlusValuesBytes (%v) must be less than or equal to %v", cfg.GetMaxTreeSyncChunkKeysPlusValuesBytes(), MaxMaxTreeSyncChunkKeysPlusValuesBytes)
	}

	// MaxParallelTreeSyncChunkFetches is already upper bounded in protocol code.

	if !(cfg.GetBlobChunkBytes() <= MaxMaxBlobChunkBytes) {
		return fmt.Errorf("BlobChunkBytes (%v) must be less than or equal to %v", cfg.GetBlobChunkBytes(), MaxMaxBlobChunkBytes)
	}

	const (
		minMaxDurationPluginCall = 10 * time.Millisecond
		maxMaxDurationPluginCall = 10 * time.Minute
	)
	if !(minMaxDurationPluginCall <= cfg.MaxDurationInitialization) {
		return fmt.Errorf("MaxDurationInitialization (%v) must be greater than or equal to %v", cfg.MaxDurationInitialization, minMaxDurationPluginCall)
	}
	if !(cfg.MaxDurationInitialization <= maxMaxDurationPluginCall) {
		return fmt.Errorf("MaxDurationInitialization (%v) must be less than or equal to %v", cfg.MaxDurationInitialization, maxMaxDurationPluginCall)
	}
	if !(minMaxDurationPluginCall <= cfg.MaxDurationShouldAcceptAttestedReport) {
		return fmt.Errorf("MaxDurationShouldAcceptAttestedReport (%v) must be greater than or equal to %v", cfg.MaxDurationShouldAcceptAttestedReport, minMaxDurationPluginCall)
	}
	if !(cfg.MaxDurationShouldAcceptAttestedReport <= maxMaxDurationPluginCall) {
		return fmt.Errorf("MaxDurationShouldAcceptAttestedReport (%v) must be less than or equal to %v", cfg.MaxDurationShouldAcceptAttestedReport, maxMaxDurationPluginCall)
	}
	if !(minMaxDurationPluginCall <= cfg.MaxDurationShouldTransmitAcceptedReport) {
		return fmt.Errorf("MaxDurationShouldTransmitAcceptedReport (%v) must be greater than or equal to %v", cfg.MaxDurationShouldTransmitAcceptedReport, minMaxDurationPluginCall)
	}
	if !(cfg.MaxDurationShouldTransmitAcceptedReport <= maxMaxDurationPluginCall) {
		return fmt.Errorf("MaxDurationShouldTransmitAcceptedReport (%v) must be less than or equal to %v", cfg.MaxDurationShouldTransmitAcceptedReport, maxMaxDurationPluginCall)
	}

	// We don't check DeltaGrace, DeltaStage since none of them would exhaust
	// the oracle's resources even if they are all set to 0.
	return nil
}
