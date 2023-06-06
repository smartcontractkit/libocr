package protocol

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const signedObservationDomainSeparator = "ocr3 SignedObservation"

type SignedObservation struct {
	Observation types.Observation
	Signature   []byte
}

func MakeSignedObservation(
	ocr3ts Timestamp,
	query types.Query,
	observation types.Observation,
	signer func(msg []byte) (sig []byte, err error),
) (
	SignedObservation,
	error,
) {
	payload := signedObservationMsg(ocr3ts, query, observation)
	sig, err := signer(payload)
	if err != nil {
		return SignedObservation{}, err
	}
	return SignedObservation{observation, sig}, nil
}

func (so SignedObservation) Verify(ocr3ts Timestamp, query types.Query, publicKey types.OffchainPublicKey) error {
	pk := ed25519.PublicKey(publicKey[:])
	// should never trigger since types.OffchainPublicKey is an array with length ed25519.PublicKeySize
	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, signedObservationMsg(ocr3ts, query, so.Observation), so.Signature)
	if !ok {
		return fmt.Errorf("SignedObservation has invalid signature")
	}

	return nil
}

func signedObservationMsg(ocr3ts Timestamp, query types.Query, observation types.Observation) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(prepareSignatureDomainSeparator))

	// ocr3ts
	_, _ = h.Write(ocr3ts.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ocr3ts.Epoch)

	// query
	_ = binary.Write(h, binary.BigEndian, uint64(len(query)))
	_, _ = h.Write(query)

	// observation
	_ = binary.Write(h, binary.BigEndian, uint64(len(observation)))
	_, _ = h.Write(observation)

	return h.Sum(nil)
}

type AttributedSignedObservation struct {
	SignedObservation SignedObservation
	Observer          commontypes.OracleID
}

const prepareSignatureDomainSeparator = "ocr3 PrepareSignature"

type OutcomeInputsDigest [32]byte

func MakeOutcomeInputsDigest(
	ocr3ts Timestamp,
	previousOutcome ocr3types.Outcome,
	seqNr uint64,
	query types.Query,
	attributedObservations []types.AttributedObservation,
) OutcomeInputsDigest {
	h := sha256.New()

	_, _ = h.Write(ocr3ts.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ocr3ts.Epoch)

	_ = binary.Write(h, binary.BigEndian, uint64(len(previousOutcome)))
	_, _ = h.Write(previousOutcome)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(query)))
	_, _ = h.Write(query)

	_ = binary.Write(h, binary.BigEndian, uint64(len(attributedObservations)))
	for _, ao := range attributedObservations {

		_ = binary.Write(h, binary.BigEndian, uint64(len(ao.Observation)))
		_, _ = h.Write(ao.Observation)

		_ = binary.Write(h, binary.BigEndian, uint64(ao.Observer))
	}

	var result OutcomeInputsDigest
	h.Sum(result[:0])

	if result == (OutcomeInputsDigest{}) {
		panic("wtf")
	}
	return result
}

type OutcomeDigest [32]byte

func MakeOutcomeDigest(outcome ocr3types.Outcome) OutcomeDigest {
	h := sha256.New()

	_, _ = h.Write(outcome)

	var result OutcomeDigest
	h.Sum(result[:0])
	if result == (OutcomeDigest{}) {
		panic("wtf")
	}
	return result
}

type PrepareSignature []byte

func MakePrepareSignature(
	ocr3ts Timestamp,
	seqNr uint64,
	outcomeInputsDigest OutcomeInputsDigest,
	outcomeDigest OutcomeDigest,
	signer func(msg []byte) ([]byte, error),
) (PrepareSignature, error) {

	return signer(prepareSignatureMsg(ocr3ts, seqNr, outcomeInputsDigest, outcomeDigest))
}

func (sig PrepareSignature) Verify(
	ocr3ts Timestamp,
	seqNr uint64,
	outcomeInputsDigest OutcomeInputsDigest,
	outcomeDigest OutcomeDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, prepareSignatureMsg(ocr3ts, seqNr, outcomeInputsDigest, outcomeDigest), sig)
	if !ok {
		return fmt.Errorf("PrepareSignature failed to verify")
	}

	return nil
}

func prepareSignatureMsg(
	ocr3ts Timestamp,
	seqNr uint64,
	outcomeInputsDigest OutcomeInputsDigest,
	outcomeDigest OutcomeDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(prepareSignatureDomainSeparator))

	_, _ = h.Write(ocr3ts.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ocr3ts.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_, _ = h.Write(outcomeInputsDigest[:])

	_, _ = h.Write(outcomeDigest[:])

	return h.Sum(nil)
}

type AttributedPrepareSignature struct {
	Signature PrepareSignature
	Signer    commontypes.OracleID
}

const commitSignatureDomainSeparator = "ocr3 CommitSignature"

type CommitSignature []byte

func MakeCommitSignature(
	ocr3ts Timestamp,
	seqNr uint64,
	outcomeDigest OutcomeDigest,
	signer func(msg []byte) ([]byte, error),
) (CommitSignature, error) {

	return signer(commitSignatureMsg(ocr3ts, seqNr, outcomeDigest))
}

func (sig CommitSignature) Verify(
	ocr3ts Timestamp,
	seqNr uint64,
	outcomeDigest OutcomeDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, commitSignatureMsg(ocr3ts, seqNr, outcomeDigest), sig)
	if !ok {
		return fmt.Errorf("CommitSignature failed to verify")
	}

	return nil
}

func commitSignatureMsg(
	ocr3ts Timestamp,
	seqNr uint64,
	outcomeDigest OutcomeDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(commitSignatureDomainSeparator))

	_, _ = h.Write(ocr3ts.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ocr3ts.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_, _ = h.Write(outcomeDigest[:])

	return h.Sum(nil)
}

type AttributedCommitSignature struct {
	Signature CommitSignature
	Signer    commontypes.OracleID
}

type HighestCertifiedTimestamp struct {
	SeqNr                 uint64
	CommittedElsePrepared bool
}

func (t HighestCertifiedTimestamp) Less(t2 HighestCertifiedTimestamp) bool {
	return t.SeqNr < t2.SeqNr || t.SeqNr == t2.SeqNr && !t.CommittedElsePrepared && t2.CommittedElsePrepared
}

const signedHighestCertifiedTimestampDomainSeparator = "ocr3 SignedHighestCertifiedTimestamp"

type SignedHighestCertifiedTimestamp struct {
	HighestCertifiedTimestamp HighestCertifiedTimestamp
	Signature                 []byte
}

func MakeSignedHighestCertifiedTimestamp(
	ocr3ts Timestamp,
	highestCertifiedTimestamp HighestCertifiedTimestamp,
	signer func(msg []byte) ([]byte, error),
) (SignedHighestCertifiedTimestamp, error) {
	sig, err := signer(signedHighestCertifiedTimestampMsg(ocr3ts, highestCertifiedTimestamp))
	if err != nil {
		return SignedHighestCertifiedTimestamp{}, err
	}

	return SignedHighestCertifiedTimestamp{
		highestCertifiedTimestamp,
		sig,
	}, nil
}

func (shct *SignedHighestCertifiedTimestamp) Verify(ocr3ts Timestamp, publicKey types.OffchainPublicKey) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, signedHighestCertifiedTimestampMsg(ocr3ts, shct.HighestCertifiedTimestamp), shct.Signature)
	if !ok {
		return fmt.Errorf("SignedHighestCertifiedTimestamp signature failed to verify")
	}

	return nil
}

func signedHighestCertifiedTimestampMsg(
	ocr3ts Timestamp,
	highestCertifiedTimestamp HighestCertifiedTimestamp,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(signedHighestCertifiedTimestampDomainSeparator))

	_, _ = h.Write(ocr3ts.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ocr3ts.Epoch)

	_ = binary.Write(h, binary.BigEndian, highestCertifiedTimestamp.SeqNr)

	var committedElsePreparedByte uint8
	if highestCertifiedTimestamp.CommittedElsePrepared {
		committedElsePreparedByte = 1
	} else {
		committedElsePreparedByte = 0
	}
	_, _ = h.Write([]byte{byte(committedElsePreparedByte)})

	return h.Sum(nil)
}

type AttributedSignedHighestCertifiedTimestamp struct {
	SignedHighestCertifiedTimestamp SignedHighestCertifiedTimestamp
	Signer                          commontypes.OracleID
}

type StartEpochProof struct {
	HighestCertified      CertifiedPrepareOrCommit
	HighestCertifiedProof []AttributedSignedHighestCertifiedTimestamp
}

func (qc *StartEpochProof) Verify(
	ocr3ts Timestamp,
	oracleIdentities []config.OracleIdentity,
	n int,
	f int,
) error {
	if ByzQuorumSize(n, f) != len(qc.HighestCertifiedProof) {
		return fmt.Errorf("wrong length of HighestCertifiedProof, expected %v for byz. quorum and got %v", ByzQuorumSize(n, f), len(qc.HighestCertifiedProof))
	}

	maximumTimestamp := qc.HighestCertifiedProof[0].SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp

	seen := make(map[commontypes.OracleID]bool)
	for i, ashct := range qc.HighestCertifiedProof {
		if seen[ashct.Signer] {
			return fmt.Errorf("duplicate signature by %v", ashct.Signer)
		}
		seen[ashct.Signer] = true
		if !(0 <= int(ashct.Signer) && int(ashct.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", ashct.Signer)
		}
		if err := ashct.SignedHighestCertifiedTimestamp.Verify(ocr3ts, oracleIdentities[ashct.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, ashct.Signer, oracleIdentities[ashct.Signer].OffchainPublicKey, err)
		}

		if maximumTimestamp.Less(ashct.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp) {
			maximumTimestamp = ashct.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp
		}
	}

	if qc.HighestCertified.Timestamp() != maximumTimestamp {
		return fmt.Errorf("mismatch between timestamp of HighestCertified (%v) and the max from HighestCertifiedProof (%v)", qc.HighestCertified.Timestamp(), maximumTimestamp)
	}

	if err := qc.HighestCertified.Verify(ocr3ts.ConfigDigest, oracleIdentities, n, f); err != nil {
		return fmt.Errorf("failed to verify HighestCertified: %w", err)
	}

	return nil
}

type CertifiedPrepareOrCommit interface {
	isCertifiedPrepareOrCommit()
	Epoch() uint64
	Timestamp() HighestCertifiedTimestamp
	IsGenesis() bool
	Verify(
		_ types.ConfigDigest,
		_ []config.OracleIdentity,
		n int,
		f int,
	) error
}

var _ CertifiedPrepareOrCommit = &CertifiedPrepareOrCommitPrepare{}

type CertifiedPrepareOrCommitPrepare struct {
	PrepareEpoch             uint64
	SeqNr                    uint64
	OutcomeInputsDigest      OutcomeInputsDigest
	Outcome                  ocr3types.Outcome
	PrepareQuorumCertificate []AttributedPrepareSignature
}

func (hc *CertifiedPrepareOrCommitPrepare) isCertifiedPrepareOrCommit() {}

func (hc *CertifiedPrepareOrCommitPrepare) Epoch() uint64 {
	return uint64(hc.PrepareEpoch)
}

func (hc *CertifiedPrepareOrCommitPrepare) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		hc.SeqNr,
		false,
	}
}

func (hc *CertifiedPrepareOrCommitPrepare) IsGenesis() bool {
	return false
}

func (hc *CertifiedPrepareOrCommitPrepare) Verify(
	configDigest types.ConfigDigest,
	oracleIdentities []config.OracleIdentity,
	n int,
	f int,
) error {
	if ByzQuorumSize(n, f) != len(hc.PrepareQuorumCertificate) {
		return fmt.Errorf("wrong number of signatures, expected %v for byz. quorum and got %v", ByzQuorumSize(n, f), len(hc.PrepareQuorumCertificate))
	}

	ocr3ts := Timestamp{
		configDigest,
		hc.PrepareEpoch,
	}

	seen := make(map[commontypes.OracleID]bool)
	for i, aps := range hc.PrepareQuorumCertificate {
		if seen[aps.Signer] {
			return fmt.Errorf("duplicate signature by %v", aps.Signer)
		}
		seen[aps.Signer] = true
		if !(0 <= int(aps.Signer) && int(aps.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", aps.Signer)
		}
		if err := aps.Signature.Verify(ocr3ts, hc.SeqNr, hc.OutcomeInputsDigest, MakeOutcomeDigest(hc.Outcome), oracleIdentities[aps.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, aps.Signer, oracleIdentities[aps.Signer].OffchainPublicKey, err)
		}
	}
	return nil
}

var _ CertifiedPrepareOrCommit = &CertifiedPrepareOrCommitCommit{}

type CertifiedPrepareOrCommitCommit struct {
	CommitEpoch             uint64
	SeqNr                   uint64
	Outcome                 ocr3types.Outcome
	CommitQuorumCertificate []AttributedCommitSignature
}

func (hc *CertifiedPrepareOrCommitCommit) isCertifiedPrepareOrCommit() {}

func (hc *CertifiedPrepareOrCommitCommit) Epoch() uint64 {
	return uint64(hc.CommitEpoch)
}

func (hc *CertifiedPrepareOrCommitCommit) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		hc.SeqNr,
		true,
	}
}

func (hc *CertifiedPrepareOrCommitCommit) IsGenesis() bool {
	return hc.CommitEpoch == 0 && hc.SeqNr == 0 && len(hc.Outcome) == 0 && len(hc.CommitQuorumCertificate) == 0
}

func (hc *CertifiedPrepareOrCommitCommit) Verify(
	configDigest types.ConfigDigest,
	oracleIdentities []config.OracleIdentity,
	n int,
	f int,
) error {

	if hc.IsGenesis() {
		return nil
	}

	if ByzQuorumSize(n, f) != len(hc.CommitQuorumCertificate) {

		return fmt.Errorf("wrong number of signatures, expected %v for byz. quorum and got %v. hc %+v", ByzQuorumSize(n, f), len(hc.CommitQuorumCertificate), hc)
	}

	ocr3ts := Timestamp{
		configDigest,
		hc.CommitEpoch,
	}

	seen := make(map[commontypes.OracleID]bool)
	for i, acs := range hc.CommitQuorumCertificate {
		if seen[acs.Signer] {
			return fmt.Errorf("duplicate signature by %v", acs.Signer)
		}
		seen[acs.Signer] = true
		if !(0 <= int(acs.Signer) && int(acs.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", acs.Signer)
		}
		if err := acs.Signature.Verify(ocr3ts, hc.SeqNr, MakeOutcomeDigest(hc.Outcome), oracleIdentities[acs.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, acs.Signer, oracleIdentities[acs.Signer].OffchainPublicKey, err)
		}
	}
	return nil
}
