package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Decodes permissively from various hex representations, with and without
// leading 0x, and is case insensitive.
type HexBytes []byte

func (hb *HexBytes) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	dec, err := HexBytesFromString(s)
	if err != nil {
		return err
	}
	*hb = dec
	return nil
}

func HexBytesFromString(s string) (HexBytes, error) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}

	decoded, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}
	return decoded, nil
}

func parseContractConfigJSON(data []byte) (types.ContractConfig, error) {
	type HumanContractConfig struct {
		ConfigDigest          HexBytes
		ConfigCount           uint64
		Signers               []HexBytes
		Transmitters          []string
		F                     uint8
		OnchainConfig         HexBytes
		OffchainConfigVersion uint64
		OffchainConfig        HexBytes
	}
	var hcc HumanContractConfig
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&hcc); err != nil {
		return types.ContractConfig{}, fmt.Errorf("failed to parse contract config JSON: %w", err)
	}

	configDigest, err := types.BytesToConfigDigest(hcc.ConfigDigest)
	if err != nil {
		return types.ContractConfig{}, fmt.Errorf("invalid ConfigDigest: %w", err)
	}

	var signers []types.OnchainPublicKey
	for _, signer := range hcc.Signers {
		signers = append(signers, types.OnchainPublicKey(signer))
	}

	var transmitters []types.Account
	for _, transmitter := range hcc.Transmitters {
		transmitters = append(transmitters, types.Account(transmitter))
	}

	return types.ContractConfig{
		configDigest,
		hcc.ConfigCount,
		signers,
		transmitters,
		hcc.F,
		[]byte(hcc.OnchainConfig),
		hcc.OffchainConfigVersion,
		[]byte(hcc.OffchainConfig),
	}, nil
}
