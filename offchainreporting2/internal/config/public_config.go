package config

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// PublicConfig is the configuration disseminated through the smart contract
// It's public, because anybody can read it from the blockchain
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

	F int

	OnchainConfig []byte

	ConfigDigest types.ConfigDigest
}

type OracleIdentity struct {
	OffchainPublicKey types.OffchainPublicKey
	OnchainPublicKey  types.OnchainPublicKey
	PeerID            string
	TransmitAccount   types.Account
}

// N is the number of oracles participating in the protocol
func (c *PublicConfig) N() int {
	return len(c.OracleIdentities)
}

func (c *PublicConfig) CheckParameterBounds() error {
	if c.F < 0 || c.F > math.MaxUint8 {
		return errors.Errorf("number of potentially faulty oracles must fit in 8 bits.")
	}
	return nil
}

func PublicConfigFromContractConfig(skipResourceExhaustionChecks bool, change types.ContractConfig) (PublicConfig, error) {
	pubcon, _, err := publicConfigFromContractConfig(skipResourceExhaustionChecks, change)
	return pubcon, err
}

func publicConfigFromContractConfig(skipResourceExhaustionChecks bool, change types.ContractConfig) (PublicConfig, SharedSecretEncryptions, error) {
	if change.OffchainConfigVersion != OffchainConfigVersion {
		return PublicConfig{}, SharedSecretEncryptions{}, fmt.Errorf("unsuppported OffchainConfigVersion %v, supported OffchainConfigVersion is %v", change.OffchainConfigVersion, OffchainConfigVersion)
	}

	oc, err := deserializeOffchainConfig(change.OffchainConfig)
	if err != nil {
		return PublicConfig{}, SharedSecretEncryptions{}, err
	}

	if err := checkIdentityListsHaveNoDuplicates(change, oc); err != nil {
		return PublicConfig{}, SharedSecretEncryptions{}, err
	}

	// must check that all lists have the same length, or bad input could crash
	// the following for loop.
	if err := checkIdentityListsHaveTheSameLength(change, oc); err != nil {
		return PublicConfig{}, SharedSecretEncryptions{}, err
	}

	identities := []OracleIdentity{}
	for i := range change.Signers {
		identities = append(identities, OracleIdentity{
			oc.OffchainPublicKeys[i],
			types.OnchainPublicKey(change.Signers[i][:]),
			oc.PeerIDs[i],
			change.Transmitters[i],
		})
	}

	cfg := PublicConfig{
		oc.DeltaProgress,
		oc.DeltaResend,
		oc.DeltaRound,
		oc.DeltaGrace,
		oc.DeltaStage,
		oc.RMax,
		oc.S,
		identities,
		oc.ReportingPluginConfig,
		oc.MaxDurationQuery,
		oc.MaxDurationObservation,
		oc.MaxDurationReport,
		oc.MaxDurationShouldAcceptFinalizedReport,
		oc.MaxDurationShouldTransmitAcceptedReport,

		int(change.F),
		change.OnchainConfig,
		change.ConfigDigest,
	}

	if err := checkPublicConfigParameters(cfg); err != nil {
		return PublicConfig{}, SharedSecretEncryptions{}, err
	}

	if !skipResourceExhaustionChecks {
		if err := checkResourceExhaustion(cfg); err != nil {
			return PublicConfig{}, SharedSecretEncryptions{}, err
		}
	}

	return cfg, oc.SharedSecretEncryptions, nil
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

	if !(0 <= cfg.DeltaStage) {
		return fmt.Errorf("DeltaStage (%v) must be non-negative", cfg.DeltaStage)
	}

	if !(0 <= cfg.DeltaRound) {
		return fmt.Errorf("DeltaRound (%v) must be non-negative", cfg.DeltaRound)
	}

	if !(0 <= cfg.DeltaProgress) {
		return fmt.Errorf("DeltaProgress (%v) must be non-negative", cfg.DeltaProgress)
	}

	if !(0 <= cfg.DeltaResend) {
		return fmt.Errorf("DeltaResend (%v) must be non-negative", cfg.DeltaResend)
	}

	if !(0 <= cfg.F && cfg.F*3 < cfg.N()) {
		return fmt.Errorf("F (%v) must be non-negative and less than N/3 (N = %v)",
			cfg.F, cfg.N())
	}

	if !(cfg.N() <= types.MaxOracles) {
		return fmt.Errorf("N (%v) must be less than or equal MaxOracles (%v)",
			cfg.N(), types.MaxOracles)
	}

	if !(0 <= cfg.DeltaGrace) {
		return fmt.Errorf("DeltaGrace (%v) must be non-negative",
			cfg.DeltaGrace)
	}

	if !(0 <= cfg.MaxDurationQuery) {
		return fmt.Errorf("MaxDurationQuery (%v) must be non-negative", cfg.MaxDurationQuery)
	}

	if !(0 <= cfg.MaxDurationObservation) {
		return fmt.Errorf("MaxDurationObservation (%v) must be non-negative", cfg.MaxDurationObservation)
	}

	if !(0 <= cfg.MaxDurationReport) {
		return fmt.Errorf("MaxDurationReport (%v) must be non-negative", cfg.MaxDurationReport)
	}

	if !(0 <= cfg.MaxDurationShouldAcceptFinalizedReport) {
		return fmt.Errorf("MaxDurationShouldAcceptFinalizedReport (%v) must be non-negative", cfg.MaxDurationShouldAcceptFinalizedReport)
	}

	if !(0 <= cfg.MaxDurationShouldTransmitAcceptedReport) {
		return fmt.Errorf("MaxDurationShouldTransmitAcceptedReport (%v) must be non-negative", cfg.MaxDurationShouldTransmitAcceptedReport)
	}

	if !(cfg.DeltaRound < cfg.DeltaProgress) {
		return fmt.Errorf("DeltaRound (%v) must be less than DeltaProgress (%v)",
			cfg.DeltaRound, cfg.DeltaProgress)
	}

	sumMaxDurationsReportGeneration := cfg.MaxDurationQuery + cfg.MaxDurationObservation + cfg.MaxDurationReport
	if !(sumMaxDurationsReportGeneration < cfg.DeltaProgress) {
		return fmt.Errorf("sum of MaxDurationQuery/Observation/Report (%v) must be less than DeltaProgress (%v)",
			sumMaxDurationsReportGeneration, cfg.DeltaProgress)
	}

	// We cannot easily add a similar check for the MaxDuration variables used
	// in the transmission protocol (MaxDurationShouldAcceptFinalizedReport,
	// MaxDurationShouldTransmitAcceptedReport), because we don't know how often
	// they will be triggered. But if we assume that there is one transmission
	// for each round, we should have MaxDurationShouldAcceptFinalizedReport +
	// MaxDurationShouldTransmitAcceptedReport < round duration.

	// *less* than 255 is intentional!
	// In report_generation_leader.go, we add 1 to a round number that can equal RMax.
	if !(0 < cfg.RMax && cfg.RMax < 255) {
		return fmt.Errorf("RMax (%v) must be greater than zero and less than 255", cfg.RMax)
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

func checkResourceExhaustion(cfg PublicConfig) error {
	// Sending a NewEpoch more than every 200ms shouldn't be necessary in any
	// realistic WAN deployment and could cause resource exhaustion
	const safeInterval = 200 * time.Millisecond
	if cfg.DeltaProgress < safeInterval {
		return fmt.Errorf("DeltaProgress (%v) is set below the resource exhaustion safe interval (%v)", cfg.DeltaProgress, safeInterval)
	}
	if cfg.DeltaResend < safeInterval {
		return fmt.Errorf("DeltaResend (%v) is set below the resource exhaustion safe interval (%v)", cfg.DeltaResend, safeInterval)
	}
	// We don't check DeltaGrace, DeltaRound, DeltaStage since none of them
	// would exhaust the oracle's resources even if they are all set to 0.
	return nil
}
