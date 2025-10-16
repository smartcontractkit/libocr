package protocol

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/internal/byzquorum"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/config"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
)

// Returns a byte slice that starts with the string "ocr3" and the rest
// of which is the sum returned by h. Used for domain separation as per the comment
// on offchainreporting2plus/types.OffchainKeyring.
//
// Any signatures made with the OffchainKeyring should use ocr3DomainSeparatedSum!
func ocr3DomainSeparatedSum(h hash.Hash) []byte {
	result := make([]byte, 0, 4+32)
	result = append(result, []byte("ocr3")...)
	return h.Sum(result)
}

const signedObservationDomainSeparator = "ocr3 SignedObservation"

type SignedObservation struct {
	Observation types.Observation
	Signature   []byte
}

func MakeSignedObservation(
	ogid OutcomeGenerationID,
	seqNr uint64,
	query types.Query,
	observation types.Observation,
	signer func(msg []byte) (sig []byte, err error),
) (
	SignedObservation,
	error,
) {
	payload := signedObservationMsg(ogid, seqNr, query, observation)
	sig, err := signer(payload)
	if err != nil {
		return SignedObservation{}, err
	}
	return SignedObservation{observation, sig}, nil
}

func (so SignedObservation) Verify(ogid OutcomeGenerationID, seqNr uint64, query types.Query, publicKey types.OffchainPublicKey) error {
	pk := ed25519.PublicKey(publicKey[:])
	// should never trigger since types.OffchainPublicKey is an array with length ed25519.PublicKeySize
	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, signedObservationMsg(ogid, seqNr, query, so.Observation), so.Signature)
	if !ok {
		return fmt.Errorf("SignedObservation has invalid signature")
	}

	return nil
}

func signedObservationMsg(ogid OutcomeGenerationID, seqNr uint64, query types.Query, observation types.Observation) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(signedObservationDomainSeparator))

	// ogid
	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	// seqNr
	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, seqNr)

	// query
	_ = binary.Write(h, binary.BigEndian, uint64(len(query)))
	_, _ = h.Write(query)

	// observation
	_ = binary.Write(h, binary.BigEndian, uint64(len(observation)))
	_, _ = h.Write(observation)

	return ocr3DomainSeparatedSum(h)
}

type AttributedSignedObservation struct {
	SignedObservation SignedObservation
	Observer          commontypes.OracleID
}

type OutcomeInputsDigest [32]byte

func MakeOutcomeInputsDigest(
	ogid OutcomeGenerationID,
	previousOutcome ocr3types.Outcome,
	seqNr uint64,
	query types.Query,
	attributedObservations []types.AttributedObservation,
) OutcomeInputsDigest {
	h := sha256.New()

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

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
	return result
}

type OutcomeDigest [32]byte

func MakeOutcomeDigest(outcome ocr3types.Outcome) OutcomeDigest {
	h := sha256.New()

	_, _ = h.Write(outcome)

	var result OutcomeDigest
	h.Sum(result[:0])
	return result
}

const prepareSignatureDomainSeparator = "ocr3 PrepareSignature"

type PrepareSignature []byte

func MakePrepareSignature(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeInputsDigest OutcomeInputsDigest,
	outcomeDigest OutcomeDigest,
	signer func(msg []byte) ([]byte, error),
) (PrepareSignature, error) {
	return signer(prepareSignatureMsg(ogid, seqNr, outcomeInputsDigest, outcomeDigest))
}

func (sig PrepareSignature) Verify(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeInputsDigest OutcomeInputsDigest,
	outcomeDigest OutcomeDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, prepareSignatureMsg(ogid, seqNr, outcomeInputsDigest, outcomeDigest), sig)
	if !ok {
		// Other less common causes include leader equivocation or actually invalid signatures.
		return fmt.Errorf("PrepareSignature failed to verify. This is commonly caused by non-determinism in the ReportingPlugin")
	}

	return nil
}

func prepareSignatureMsg(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeInputsDigest OutcomeInputsDigest,
	outcomeDigest OutcomeDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(prepareSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_, _ = h.Write(outcomeInputsDigest[:])

	_, _ = h.Write(outcomeDigest[:])

	return ocr3DomainSeparatedSum(h)
}

type AttributedPrepareSignature struct {
	Signature PrepareSignature
	Signer    commontypes.OracleID
}

const commitSignatureDomainSeparator = "ocr3 CommitSignature"

type CommitSignature []byte

func MakeCommitSignature(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeDigest OutcomeDigest,
	signer func(msg []byte) ([]byte, error),
) (CommitSignature, error) {
	return signer(commitSignatureMsg(ogid, seqNr, outcomeDigest))
}

func (sig CommitSignature) Verify(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeDigest OutcomeDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, commitSignatureMsg(ogid, seqNr, outcomeDigest), sig)
	if !ok {
		return fmt.Errorf("CommitSignature failed to verify")
	}

	return nil
}

func commitSignatureMsg(
	ogid OutcomeGenerationID,
	seqNr uint64,
	outcomeDigest OutcomeDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(commitSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_, _ = h.Write(outcomeDigest[:])

	return ocr3DomainSeparatedSum(h)
}

type AttributedCommitSignature struct {
	Signature CommitSignature
	Signer    commontypes.OracleID
}

type HighestCertifiedTimestamp struct {
	SeqNr                 uint64
	CommittedElsePrepared bool
	Epoch                 uint64
}

func (t HighestCertifiedTimestamp) Less(t2 HighestCertifiedTimestamp) bool {
	return t.SeqNr < t2.SeqNr ||
		t.SeqNr == t2.SeqNr && !t.CommittedElsePrepared && t2.CommittedElsePrepared ||
		t.SeqNr == t2.SeqNr && t.CommittedElsePrepared == t2.CommittedElsePrepared && t.Epoch < t2.Epoch
}

func (t HighestCertifiedTimestamp) Equal(t2 HighestCertifiedTimestamp) bool {
	return t.SeqNr == t2.SeqNr && t.CommittedElsePrepared == t2.CommittedElsePrepared && t.Epoch == t2.Epoch
}

const signedHighestCertifiedTimestampDomainSeparator = "ocr3 SignedHighestCertifiedTimestamp"

type SignedHighestCertifiedTimestamp struct {
	HighestCertifiedTimestamp HighestCertifiedTimestamp
	Signature                 []byte
	Signature31               []byte
}

func MakeSignedHighestCertifiedTimestamp(
	ogid OutcomeGenerationID,
	highestCertifiedTimestamp HighestCertifiedTimestamp,
	signer func(msg []byte) ([]byte, error),
) (SignedHighestCertifiedTimestamp, error) {
	sig, err := signer(signedHighestCertifiedTimestampMsg(ogid, highestCertifiedTimestamp))
	if err != nil {
		return SignedHighestCertifiedTimestamp{}, err
	}

	sig31, err := signer(signedHighestCertifiedTimestamp31Msg(ogid, highestCertifiedTimestamp))
	if err != nil {
		return SignedHighestCertifiedTimestamp{}, err
	}

	return SignedHighestCertifiedTimestamp{
		highestCertifiedTimestamp,
		sig,
		sig31,
	}, nil
}

func (shct *SignedHighestCertifiedTimestamp) Verify(ogid OutcomeGenerationID, publicKey types.OffchainPublicKey) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, signedHighestCertifiedTimestampMsg(ogid, shct.HighestCertifiedTimestamp), shct.Signature)
	ok2 := ed25519.Verify(pk, signedHighestCertifiedTimestamp31Msg(ogid, shct.HighestCertifiedTimestamp), shct.Signature31)

	if !ok || !ok2 {
		return fmt.Errorf("SignedHighestCertifiedTimestamp signature failed to verify (ok: %v, ok2: %v)", ok, ok2)
	}

	return nil
}

func signedHighestCertifiedTimestampMsg(
	ogid OutcomeGenerationID,
	highestCertifiedTimestamp HighestCertifiedTimestamp,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(signedHighestCertifiedTimestampDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, highestCertifiedTimestamp.SeqNr)

	var committedElsePreparedByte uint8
	if highestCertifiedTimestamp.CommittedElsePrepared {
		committedElsePreparedByte = 1
	} else {
		committedElsePreparedByte = 0
	}
	_, _ = h.Write([]byte{byte(committedElsePreparedByte)})

	return ocr3DomainSeparatedSum(h)
}

func signedHighestCertifiedTimestamp31Msg(
	ogid OutcomeGenerationID,
	highestCertifiedTimestamp HighestCertifiedTimestamp,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(signedHighestCertifiedTimestampDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, highestCertifiedTimestamp.SeqNr)

	var committedElsePreparedByte uint8
	if highestCertifiedTimestamp.CommittedElsePrepared {
		committedElsePreparedByte = 1
	} else {
		committedElsePreparedByte = 0
	}
	_, _ = h.Write([]byte{byte(committedElsePreparedByte)})

	_ = binary.Write(h, binary.BigEndian, highestCertifiedTimestamp.Epoch)

	return ocr3DomainSeparatedSum(h)
}

type AttributedSignedHighestCertifiedTimestamp struct {
	SignedHighestCertifiedTimestamp SignedHighestCertifiedTimestamp
	Signer                          commontypes.OracleID
}

type EpochStartProof struct {
	HighestCertified      CertifiedPrepareOrCommit
	HighestCertifiedProof []AttributedSignedHighestCertifiedTimestamp
}

func (qc *EpochStartProof) Verify(
	ogid OutcomeGenerationID,
	oracleIdentities []config.OracleIdentity,
	byzQuorumSize int,
) error {
	if byzQuorumSize != len(qc.HighestCertifiedProof) {
		return fmt.Errorf("wrong length of HighestCertifiedProof, expected %v for byz. quorum and got %v", byzQuorumSize, len(qc.HighestCertifiedProof))
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
		if err := ashct.SignedHighestCertifiedTimestamp.Verify(ogid, oracleIdentities[ashct.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, ashct.Signer, oracleIdentities[ashct.Signer].OffchainPublicKey, err)
		}

		if maximumTimestamp.Less(ashct.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp) {
			maximumTimestamp = ashct.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp
		}
	}

	if !qc.HighestCertified.Timestamp().Equal(maximumTimestamp) {
		return fmt.Errorf("mismatch between timestamp of HighestCertified (%v) and the max from HighestCertifiedProof (%v)", qc.HighestCertified.Timestamp(), maximumTimestamp)
	}

	if err := qc.HighestCertified.Verify(ogid.ConfigDigest, oracleIdentities, byzQuorumSize); err != nil {
		return fmt.Errorf("failed to verify HighestCertified: %w", err)
	}

	return nil
}

type CertifiedPrepareOrCommitDigest [32]byte

type CertifiedPrepareOrCommit interface {
	isCertifiedPrepareOrCommit()
	Epoch() uint64
	Timestamp() HighestCertifiedTimestamp
	IsGenesis() bool
	Verify(
		_ types.ConfigDigest,
		_ []config.OracleIdentity,
		byzQuorumSize int,
	) error
	CheckSize(n int, f int, limits ocr3types.ReportingPluginLimits, maxReportSigLen int) bool
	Digest() CertifiedPrepareOrCommitDigest
}

var _ CertifiedPrepareOrCommit = &CertifiedPrepare{}

type CertifiedPrepare struct {
	PrepareEpoch             uint64
	SeqNr                    uint64
	OutcomeInputsDigest      OutcomeInputsDigest
	Outcome                  ocr3types.Outcome
	PrepareQuorumCertificate []AttributedPrepareSignature
}

func (hc *CertifiedPrepare) isCertifiedPrepareOrCommit() {}

func (hc *CertifiedPrepare) Epoch() uint64 {
	return uint64(hc.PrepareEpoch)
}

func (hc *CertifiedPrepare) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		hc.SeqNr,
		false,
		hc.PrepareEpoch,
	}
}

func (hc *CertifiedPrepare) IsGenesis() bool {
	return false
}

func (hc *CertifiedPrepare) Verify(
	configDigest types.ConfigDigest,
	oracleIdentities []config.OracleIdentity,
	byzQuorumSize int,
) error {
	if byzQuorumSize != len(hc.PrepareQuorumCertificate) {
		return fmt.Errorf("wrong number of signatures, expected %v for byz. quorum and got %v", byzQuorumSize, len(hc.PrepareQuorumCertificate))
	}

	ogid := OutcomeGenerationID{
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
		if err := aps.Signature.Verify(ogid, hc.SeqNr, hc.OutcomeInputsDigest, MakeOutcomeDigest(hc.Outcome), oracleIdentities[aps.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, aps.Signer, oracleIdentities[aps.Signer].OffchainPublicKey, err)
		}
	}
	return nil
}

func (hc *CertifiedPrepare) CheckSize(n int, f int, limits ocr3types.ReportingPluginLimits, maxReportSigLen int) bool {
	if len(hc.Outcome) > limits.MaxOutcomeLength {
		return false
	}
	if len(hc.PrepareQuorumCertificate) != byzquorum.Size(n, f) {
		return false
	}
	for _, aps := range hc.PrepareQuorumCertificate {
		if len(aps.Signature) != ed25519.SignatureSize {
			return false
		}
	}
	return true
}

// domain separation from other variants of CertifiedPrepareOrCommit
const certifiedPrepareDigestDomainSeparator = "ocr3 CertifiedPrepareDigest"

func (hc *CertifiedPrepare) Digest() CertifiedPrepareOrCommitDigest {
	h := sha256.New()

	_, _ = h.Write([]byte(certifiedPrepareDigestDomainSeparator))

	_ = binary.Write(h, binary.BigEndian, hc.PrepareEpoch)

	_ = binary.Write(h, binary.BigEndian, hc.SeqNr)

	_, _ = h.Write(hc.OutcomeInputsDigest[:])

	_ = binary.Write(h, binary.BigEndian, uint64(len(hc.Outcome)))
	_, _ = h.Write(hc.Outcome)

	_ = binary.Write(h, binary.BigEndian, uint64(len(hc.PrepareQuorumCertificate)))
	for _, aps := range hc.PrepareQuorumCertificate {

		_ = binary.Write(h, binary.BigEndian, uint64(len(aps.Signature)))
		_, _ = h.Write(aps.Signature)

		_ = binary.Write(h, binary.BigEndian, uint64(aps.Signer))
	}

	var result CertifiedPrepareOrCommitDigest
	h.Sum(result[:0])
	return result
}

var _ CertifiedPrepareOrCommit = &CertifiedCommit{}

// The empty CertifiedCommit{} is the genesis value
type CertifiedCommit struct {
	CommitEpoch             uint64
	SeqNr                   uint64
	Outcome                 ocr3types.Outcome
	CommitQuorumCertificate []AttributedCommitSignature
}

func (hc *CertifiedCommit) isCertifiedPrepareOrCommit() {}

func (hc *CertifiedCommit) Epoch() uint64 {
	return uint64(hc.CommitEpoch)
}

func (hc *CertifiedCommit) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		hc.SeqNr,
		true,
		hc.CommitEpoch,
	}
}

func (hc *CertifiedCommit) IsGenesis() bool {
	// We intentionally don't just compare with CertifiedCommit{}, because after
	// protobuf deserialization, we might end up with hc.Outcome = []byte{}
	return hc.CommitEpoch == 0 && hc.SeqNr == 0 && len(hc.Outcome) == 0 && len(hc.CommitQuorumCertificate) == 0
}

func (hc *CertifiedCommit) Verify(
	configDigest types.ConfigDigest,
	oracleIdentities []config.OracleIdentity,
	byzQuorumSize int,
) error {
	if hc.IsGenesis() {
		return nil
	}

	if byzQuorumSize != len(hc.CommitQuorumCertificate) {
		return fmt.Errorf("wrong number of signatures, expected %d for byz. quorum but got %d", byzQuorumSize, len(hc.CommitQuorumCertificate))
	}

	ogid := OutcomeGenerationID{
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
		if err := acs.Signature.Verify(ogid, hc.SeqNr, MakeOutcomeDigest(hc.Outcome), oracleIdentities[acs.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, acs.Signer, oracleIdentities[acs.Signer].OffchainPublicKey, err)
		}
	}
	return nil
}

func (hc *CertifiedCommit) CheckSize(n int, f int, limits ocr3types.ReportingPluginLimits, maxReportSigLen int) bool {
	if hc.IsGenesis() {
		return true
	}

	if len(hc.Outcome) > limits.MaxOutcomeLength {
		return false
	}
	if len(hc.CommitQuorumCertificate) != byzquorum.Size(n, f) {
		return false
	}
	for _, acs := range hc.CommitQuorumCertificate {
		if len(acs.Signature) != ed25519.SignatureSize {
			return false
		}
	}
	return true
}

// domain separation from other variants of CertifiedPrepareOrCommit
const certifiedCommitDigestDomainSeparator = "ocr3 CertifiedCommitDigest"

func (hc *CertifiedCommit) Digest() CertifiedPrepareOrCommitDigest {
	h := sha256.New()

	_, _ = h.Write([]byte(certifiedCommitDigestDomainSeparator))

	_ = binary.Write(h, binary.BigEndian, hc.CommitEpoch)

	_ = binary.Write(h, binary.BigEndian, hc.SeqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(hc.Outcome)))
	_, _ = h.Write(hc.Outcome)

	_ = binary.Write(h, binary.BigEndian, uint64(len(hc.CommitQuorumCertificate)))
	for _, acs := range hc.CommitQuorumCertificate {

		_ = binary.Write(h, binary.BigEndian, uint64(len(acs.Signature)))
		_, _ = h.Write(acs.Signature)

		_ = binary.Write(h, binary.BigEndian, uint64(acs.Signer))
	}

	var result CertifiedPrepareOrCommitDigest
	h.Sum(result[:0])
	return result
}

const epochStartSignatureDomainSeparator = "ocr3 EpochStartSignature"

type EpochStartSignature31 []byte

func MakeEpochStartSignature31(
	ogid OutcomeGenerationID,
	epochStartProof EpochStartProof,
	signer func(msg []byte) ([]byte, error),
) (EpochStartSignature31, error) {
	return signer(epochStartSignatureMsg(ogid, epochStartProof))
}

func (sig EpochStartSignature31) Verify(
	ogid OutcomeGenerationID,
	epochStartProof EpochStartProof,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, epochStartSignatureMsg(ogid, epochStartProof), sig)
	if !ok {
		return fmt.Errorf("EpochStartSignature31 failed to verify")
	}

	return nil
}

func epochStartSignatureMsg(ogid OutcomeGenerationID, esp EpochStartProof) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(epochStartSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	var highestCertifiedDigest CertifiedPrepareOrCommitDigest
	if esp.HighestCertified != nil {
		highestCertifiedDigest = esp.HighestCertified.Digest()
	}
	_, _ = h.Write(highestCertifiedDigest[:])

	_ = binary.Write(h, binary.BigEndian, uint64(len(esp.HighestCertifiedProof)))
	for _, ashct := range esp.HighestCertifiedProof {

		_ = binary.Write(h, binary.BigEndian, ashct.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp.SeqNr)

		var committedElsePreparedByte uint8
		if ashct.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp.CommittedElsePrepared {
			committedElsePreparedByte = 1
		} else {
			committedElsePreparedByte = 0
		}
		_, _ = h.Write([]byte{byte(committedElsePreparedByte)})

		_ = binary.Write(h, binary.BigEndian, ashct.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp.Epoch)

		_ = binary.Write(h, binary.BigEndian, uint64(len(ashct.SignedHighestCertifiedTimestamp.Signature31)))
		_, _ = h.Write(ashct.SignedHighestCertifiedTimestamp.Signature31)

		_ = binary.Write(h, binary.BigEndian, uint64(ashct.Signer))
	}

	return ocr3DomainSeparatedSum(h)
}

const roundStartSignatureDomainSeparator = "ocr3 RoundStartSignature"

type RoundStartSignature31 []byte

func MakeRoundStartSignature31(
	ogid OutcomeGenerationID,
	seqNr uint64,
	query types.Query,
	signer func(msg []byte) ([]byte, error),
) (RoundStartSignature31, error) {
	return signer(roundStartSignatureMsg(ogid, seqNr, query))
}

func (sig RoundStartSignature31) Verify(
	ogid OutcomeGenerationID,
	seqNr uint64,
	query types.Query,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, roundStartSignatureMsg(ogid, seqNr, query), sig)
	if !ok {
		return fmt.Errorf("RoundStartSignature31 failed to verify")
	}
	return nil
}

func roundStartSignatureMsg(
	ogid OutcomeGenerationID,
	seqNr uint64,
	query types.Query,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(roundStartSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(query)))
	_, _ = h.Write(query)

	return ocr3DomainSeparatedSum(h)
}

const proposalSignatureDomainSeparator = "ocr3 ProposalSignature"

type ProposalSignature31 []byte

func MakeProposalSignature31(
	ogid OutcomeGenerationID,
	seqNr uint64,
	attributedSignedObservations []AttributedSignedObservation,
	signer func(msg []byte) ([]byte, error),
) (ProposalSignature31, error) {
	return signer(proposalSignatureMsg(ogid, seqNr, attributedSignedObservations))
}

func (sig ProposalSignature31) Verify(
	ogid OutcomeGenerationID,
	seqNr uint64,
	attributedSignedObservations []AttributedSignedObservation,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, proposalSignatureMsg(ogid, seqNr, attributedSignedObservations), sig)
	if !ok {
		return fmt.Errorf("ProposalSignature31 failed to verify")
	}

	return nil
}

func proposalSignatureMsg(
	ogid OutcomeGenerationID,
	seqNr uint64,
	attributedSignedObservations []AttributedSignedObservation,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(proposalSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(attributedSignedObservations)))
	for _, aso := range attributedSignedObservations {

		_ = binary.Write(h, binary.BigEndian, uint64(len(aso.SignedObservation.Observation)))
		_, _ = h.Write(aso.SignedObservation.Observation)

		_ = binary.Write(h, binary.BigEndian, uint64(len(aso.SignedObservation.Signature)))
		_, _ = h.Write(aso.SignedObservation.Signature)

		_ = binary.Write(h, binary.BigEndian, uint64(aso.Observer))
	}

	return ocr3DomainSeparatedSum(h)
}
