package config

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

// PublicConfig is the configuration disseminated through the smart contract
// It's public, because anybody can read it from the blockchain
type PublicConfig struct {
	DeltaProgress    time.Duration
	DeltaResend      time.Duration
	DeltaRound       time.Duration
	DeltaGrace       time.Duration
	DeltaC           time.Duration
	AlphaPPB         uint64
	DeltaStage       time.Duration
	RMax             uint8
	S                []int
	OracleIdentities []OracleIdentity

	F            int
	ConfigDigest types.ConfigDigest
}

type OracleIdentity struct {
	PeerID                string
	OffchainPublicKey     types.OffchainPublicKey
	OnChainSigningAddress types.OnChainSigningAddress
	TransmitAddress       common.Address
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

func PublicConfigFromContractConfig(change types.ContractConfig) (PublicConfig, error) {
	pubcon, _, err := publicConfigFromContractConfig(change)
	return pubcon, err
}

func publicConfigFromContractConfig(change types.ContractConfig) (PublicConfig, SharedSecretEncryptions, error) {
	oc, err := decodeContractSetConfigEncodedComponents(change.Encoded)
	if err != nil {
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
			oc.PeerIDs[i],
			oc.OffchainPublicKeys[i],
			types.OnChainSigningAddress(change.Signers[i]),
			change.Transmitters[i],
		})
	}

	cfg := PublicConfig{
		oc.DeltaProgress,
		oc.DeltaResend,
		oc.DeltaRound,
		oc.DeltaGrace,
		oc.DeltaC,
		oc.AlphaPPB,
		oc.DeltaStage,
		oc.RMax,
		oc.S,
		identities,
		int(change.Threshold),
		change.ConfigDigest,
	}

	if err := checkPublicConfigParameters(cfg); err != nil {
		return PublicConfig{}, SharedSecretEncryptions{}, err
	}

	return cfg, oc.SharedSecretEncryptions, nil
}

func checkIdentityListsHaveTheSameLength(
	change types.ContractConfig, oc setConfigEncodedComponents,
) error {
	expectedLength := len(change.Signers)
	errorMsg := "%s list must have same length as onchain signers list: %d ≠ " +
		strconv.Itoa(expectedLength)
	for _, identityList := range []struct {
		length int
		name   string
	}{
		{len(oc.PeerIDs) /*                       */, "peer ID"},
		{len(oc.OffchainPublicKeys) /*            */, "offchain public keys"},
		{len(change.Transmitters) /*              */, "transmitter address"},
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
	// be made when you change these values!
	/////////////////////////////////////////////////////////////////

	if !(0 <= cfg.DeltaC) {
		return fmt.Errorf("DeltaC (%v) must be non-negative",
			cfg.DeltaC)
	}

	if !(1*time.Second < cfg.DeltaStage) {
		return fmt.Errorf("DeltaStage (%v) must be greater than 1s",
			cfg.DeltaStage)
	}

	if !(500*time.Millisecond < cfg.DeltaRound) {
		return fmt.Errorf("DeltaRound (%v) must be greater than 500ms",
			cfg.DeltaRound)
	}

	if !(500*time.Millisecond < cfg.DeltaProgress) {
		return fmt.Errorf("DeltaProgress (%v) must be greater than 500ms",
			cfg.DeltaProgress)
	}

	if !(500*time.Millisecond < cfg.DeltaResend) {
		return fmt.Errorf("DeltaResend (%v) must be greater than 500ms",
			cfg.DeltaResend)
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

	if !(cfg.DeltaGrace < cfg.DeltaRound) {
		return fmt.Errorf("DeltaGrace (%v) must be less than DeltaRound (%v)",
			cfg.DeltaGrace, cfg.DeltaRound)
	}

	if !(cfg.DeltaRound < cfg.DeltaProgress) {
		return fmt.Errorf("DeltaRound (%v) must be less than DeltaProgress (%v)",
			cfg.DeltaRound, cfg.DeltaProgress)
	}

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
