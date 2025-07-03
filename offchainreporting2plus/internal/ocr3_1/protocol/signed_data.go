package protocol

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/byzquorum"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Returns a byte slice whose first four bytes are the string "ocr3" and the rest
// of which is the sum returned by h. Used for domain separation vs ocr2, where
// we just directly sign sha256 hashes.
//
// Any signatures made with the OffchainKeyring should use ocr3_1DomainSeparatedSum!
func ocr3_1DomainSeparatedSum(h hash.Hash) []byte {
	result := make([]byte, 0, 6+32)
	result = append(result, []byte("ocr3.1")...)
	return h.Sum(result)
}

const signedObservationDomainSeparator = "ocr3.1 SignedObservation"

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

	return ocr3_1DomainSeparatedSum(h)
}

type AttributedSignedObservation struct {
	SignedObservation SignedObservation
	Observer          commontypes.OracleID
}

type StateTransitionInputsDigest [32]byte

func MakeStateTransitionInputsDigest(
	ogid OutcomeGenerationID,
	seqNr uint64,
	query types.Query,
	attributedObservations []types.AttributedObservation,
) StateTransitionInputsDigest {
	h := sha256.New()

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(query)))
	_, _ = h.Write(query)

	_ = binary.Write(h, binary.BigEndian, uint64(len(attributedObservations)))
	for _, ao := range attributedObservations {

		_ = binary.Write(h, binary.BigEndian, uint64(len(ao.Observation)))
		_, _ = h.Write(ao.Observation)

		_ = binary.Write(h, binary.BigEndian, uint64(ao.Observer))
	}

	var result StateTransitionInputsDigest
	h.Sum(result[:0])
	return result
}

type StateTransitionOutputDigest [32]byte

func MakeStateTransitionOutputDigest(ogid OutcomeGenerationID, seqNr uint64, output []KeyValuePair) StateTransitionOutputDigest {
	h := sha256.New()

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(output)))
	for _, o := range output {

		_ = binary.Write(h, binary.BigEndian, uint64(len(o.Key)))
		_, _ = h.Write(o.Key)

		_ = binary.Write(h, binary.BigEndian, uint64(len(o.Value)))
		_, _ = h.Write(o.Value)
	}

	var result StateTransitionOutputDigest
	h.Sum(result[:0])
	return result
}

type ReportsPlusPrecursorDigest [32]byte

func MakeReportsPlusPrecursorDigest(ogid OutcomeGenerationID, seqNr uint64, precursor ocr3_1types.ReportsPlusPrecursor) ReportsPlusPrecursorDigest {
	h := sha256.New()

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(precursor)))
	_, _ = h.Write(precursor)

	var result ReportsPlusPrecursorDigest
	h.Sum(result[:0])
	return result
}

const prepareSignatureDomainSeparator = "ocr3.1 PrepareSignature"

type PrepareSignature []byte

func MakePrepareSignature(
	ogid OutcomeGenerationID,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	outputDigest StateTransitionOutputDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	signer func(msg []byte) ([]byte, error),
) (PrepareSignature, error) {
	return signer(prepareSignatureMsg(ogid, seqNr, inputsDigest, outputDigest, reportsPlusPrecursorDigest))
}

func (sig PrepareSignature) Verify(
	ogid OutcomeGenerationID,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	outputDigest StateTransitionOutputDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}
	msg := prepareSignatureMsg(ogid, seqNr, inputsDigest, outputDigest, reportsPlusPrecursorDigest)
	ok := ed25519.Verify(pk, prepareSignatureMsg(ogid, seqNr, inputsDigest, outputDigest, reportsPlusPrecursorDigest), sig)
	if !ok {
		// Other less common causes include leader equivocation or actually invalid signatures.
		return fmt.Errorf("PrepareSignature failed to verify. This is commonly caused by non-determinism in the ReportingPlugin msg: %x, sig: %x", msg, sig)
	}

	return nil
}

func prepareSignatureMsg(
	ogid OutcomeGenerationID,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	outputDigest StateTransitionOutputDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(prepareSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_, _ = h.Write(inputsDigest[:])

	_, _ = h.Write(outputDigest[:])

	_, _ = h.Write(reportsPlusPrecursorDigest[:])

	return ocr3_1DomainSeparatedSum(h)
}

type AttributedPrepareSignature struct {
	Signature PrepareSignature
	Signer    commontypes.OracleID
}

const commitSignatureDomainSeparator = "ocr3.1 CommitSignature"

type CommitSignature []byte

func MakeCommitSignature(
	ogid OutcomeGenerationID,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	outputDigest StateTransitionOutputDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	signer func(msg []byte) ([]byte, error),
) (CommitSignature, error) {
	return signer(commitSignatureMsg(ogid, seqNr, inputsDigest, outputDigest, reportsPlusPrecursorDigest))
}

func (sig CommitSignature) Verify(
	ogid OutcomeGenerationID,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	outputDigest StateTransitionOutputDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, commitSignatureMsg(ogid, seqNr, inputsDigest, outputDigest, reportsPlusPrecursorDigest), sig)
	if !ok {
		return fmt.Errorf("CommitSignature failed to verify")
	}

	return nil
}

func commitSignatureMsg(
	ogid OutcomeGenerationID,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	outputDigest StateTransitionOutputDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(commitSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_, _ = h.Write(inputsDigest[:])

	_, _ = h.Write(outputDigest[:])

	_, _ = h.Write(reportsPlusPrecursorDigest[:])

	return ocr3_1DomainSeparatedSum(h)
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

const signedHighestCertifiedTimestampDomainSeparator = "ocr3.1 SignedHighestCertifiedTimestamp"

type SignedHighestCertifiedTimestamp struct {
	HighestCertifiedTimestamp HighestCertifiedTimestamp
	Signature                 []byte
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

	return SignedHighestCertifiedTimestamp{
		highestCertifiedTimestamp,
		sig,
	}, nil
}

func (shct *SignedHighestCertifiedTimestamp) Verify(ogid OutcomeGenerationID, publicKey types.OffchainPublicKey) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, signedHighestCertifiedTimestampMsg(ogid, shct.HighestCertifiedTimestamp), shct.Signature)
	if !ok {
		return fmt.Errorf("SignedHighestCertifiedTimestamp signature failed to verify")
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

	_ = binary.Write(h, binary.BigEndian, highestCertifiedTimestamp.Epoch)

	return ocr3_1DomainSeparatedSum(h)
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

	if qc.HighestCertified.Timestamp() != maximumTimestamp {
		return fmt.Errorf("mismatch between timestamp of HighestCertified (%v) and the max from HighestCertifiedProof (%v)", qc.HighestCertified.Timestamp(), maximumTimestamp)
	}

	if err := qc.HighestCertified.Verify(ogid.ConfigDigest, oracleIdentities, byzQuorumSize); err != nil {
		return fmt.Errorf("failed to verify HighestCertified: %w", err)
	}

	return nil
}

type CertifiedPrepareOrCommit interface {
	isCertifiedPrepareOrCommit()
	Epoch() uint64
	SeqNr() uint64
	Timestamp() HighestCertifiedTimestamp
	IsGenesis() bool
	Verify(
		_ types.ConfigDigest,
		_ []config.OracleIdentity,
		byzQuorumSize int,
	) error
	CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool
}

var _ CertifiedPrepareOrCommit = &CertifiedPrepare{}

type CertifiedPrepare struct {
	PrepareEpoch                uint64
	PrepareSeqNr                uint64
	StateTransitionInputsDigest StateTransitionInputsDigest
	StateTransitionOutputs      StateTransitionOutputs
	ReportsPlusPrecursor        ocr3_1types.ReportsPlusPrecursor
	PrepareQuorumCertificate    []AttributedPrepareSignature
}

func (hc *CertifiedPrepare) isCertifiedPrepareOrCommit() {}

func (hc *CertifiedPrepare) Epoch() uint64 {
	return hc.PrepareEpoch
}

func (hc *CertifiedPrepare) SeqNr() uint64 {
	return hc.PrepareSeqNr
}

func (hc *CertifiedPrepare) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		hc.SeqNr(),
		false,
		hc.Epoch(),
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
		hc.Epoch(),
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
		outputDigest := MakeStateTransitionOutputDigest(
			ogid,
			hc.SeqNr(),
			hc.StateTransitionOutputs.WriteSet,
		)
		reportsPlusPrecursorDigest := MakeReportsPlusPrecursorDigest(
			ogid,
			hc.SeqNr(),
			hc.ReportsPlusPrecursor,
		)
		if err := aps.Signature.Verify(
			ogid, hc.SeqNr(), hc.StateTransitionInputsDigest, outputDigest,
			reportsPlusPrecursorDigest, oracleIdentities[aps.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w", i, aps.Signer, oracleIdentities[aps.Signer].OffchainPublicKey, err)
		}
	}
	return nil
}
func (hc *CertifiedPrepare) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {

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

var _ CertifiedPrepareOrCommit = &CertifiedCommit{}

// The empty CertifiedCommit{} is the genesis value
type CertifiedCommit struct {
	CommitEpoch                 uint64
	CommitSeqNr                 uint64
	StateTransitionInputsDigest StateTransitionInputsDigest
	StateTransitionOutputs      StateTransitionOutputs
	ReportsPlusPrecursor        ocr3_1types.ReportsPlusPrecursor
	CommitQuorumCertificate     []AttributedCommitSignature
}

func (hc *CertifiedCommit) isCertifiedPrepareOrCommit() {}

func (hc *CertifiedCommit) Epoch() uint64 {
	return uint64(hc.CommitEpoch)
}

func (hc *CertifiedCommit) SeqNr() uint64 {
	return uint64(hc.CommitSeqNr)
}

func (hc *CertifiedCommit) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		hc.SeqNr(),
		true,
		hc.Epoch(),
	}
}

func (hc *CertifiedCommit) IsGenesis() bool {
	// We intentionally don't just compare with CertifiedCommit{}, because after
	// protobuf deserialization, we might end up with hc.Outcome = []byte{}
	return hc.Epoch() == uint64(0) &&
		hc.SeqNr() == uint64(0) &&
		hc.StateTransitionInputsDigest == StateTransitionInputsDigest{} &&
		len(hc.ReportsPlusPrecursor) == 0 &&
		len(hc.CommitQuorumCertificate) == 0
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
		hc.Epoch(),
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
		outputDigest := MakeStateTransitionOutputDigest(
			ogid,
			hc.SeqNr(),
			hc.StateTransitionOutputs.WriteSet,
		)
		reportsPlusPrecursorDigest := MakeReportsPlusPrecursorDigest(
			ogid,
			hc.SeqNr(),
			hc.ReportsPlusPrecursor,
		)
		if err := acs.Signature.Verify(ogid,
			hc.SeqNr(),
			hc.StateTransitionInputsDigest,
			outputDigest,
			reportsPlusPrecursorDigest,
			oracleIdentities[acs.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle does not verify: %w", i, acs.Signer, err)
		}
	}
	return nil
}

func (hc *CertifiedCommit) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	if hc.IsGenesis() {
		return true
	}

	if len(hc.CommitQuorumCertificate) != byzquorum.Size(n, f) {
		return false
	}
	for _, aps := range hc.CommitQuorumCertificate {
		if len(aps.Signature) != ed25519.SignatureSize {
			return false
		}
	}
	return true
}

type AttestedStateTransitionBlock struct {
	StateTransitionBlock StateTransitionBlock
	AttributedSignatures []AttributedCommitSignature
}

func (astb *AttestedStateTransitionBlock) Verify(
	configDigest types.ConfigDigest,
	oracleIdentities []config.OracleIdentity,
	byzQuorumSize int,
) error {
	if byzQuorumSize != len(astb.AttributedSignatures) {
		return fmt.Errorf("wrong number of signatures, expected %d for byz. quorum but got %d", byzQuorumSize, len(astb.AttributedSignatures))
	}

	ogid := OutcomeGenerationID{
		configDigest,
		astb.StateTransitionBlock.Epoch,
	}
	seqNr := astb.StateTransitionBlock.SeqNr()

	seen := make(map[commontypes.OracleID]bool)
	for i, sig := range astb.AttributedSignatures {
		if seen[sig.Signer] {
			return fmt.Errorf("duplicate signature by %v", sig.Signer)
		}
		seen[sig.Signer] = true
		if !(0 <= int(sig.Signer) && int(sig.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", sig.Signer)
		}
		if err := sig.Signature.Verify(
			ogid,
			seqNr,
			astb.StateTransitionBlock.StateTransitionInputsDigest,
			MakeStateTransitionOutputDigest(
				ogid,
				seqNr,
				astb.StateTransitionBlock.StateTransitionOutputs.WriteSet,
			),
			MakeReportsPlusPrecursorDigest(
				ogid,
				seqNr,
				astb.StateTransitionBlock.ReportsPlusPrecursor,
			),
			oracleIdentities[sig.Signer].OffchainPublicKey,
		); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle with pubkey %x does not verify: %w",
				i, sig.Signer, oracleIdentities[sig.Signer].OffchainPublicKey, err)
		}
	}
	return nil
}

type CertifiedCommittedReports[RI any] struct {
	CommitEpoch                 uint64
	SeqNr                       uint64
	StateTransitionInputsDigest StateTransitionInputsDigest
	StateTransitionOutputDigest StateTransitionOutputDigest
	ReportsPlusPrecursor        ocr3_1types.ReportsPlusPrecursor
	CommitQuorumCertificate     []AttributedCommitSignature
}

func (ccrs *CertifiedCommittedReports[RI]) isGenesis() bool {
	return ccrs.CommitEpoch == uint64(0) && ccrs.SeqNr == uint64(0) &&
		ccrs.StateTransitionInputsDigest == StateTransitionInputsDigest{} &&
		ccrs.StateTransitionOutputDigest == StateTransitionOutputDigest{} &&
		len(ccrs.ReportsPlusPrecursor) == 0 &&
		len(ccrs.CommitQuorumCertificate) == 0
}

func (ccrs *CertifiedCommittedReports[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	// Check if we this is a genesis certificate
	if ccrs.isGenesis() {
		return true
	}

	if len(ccrs.ReportsPlusPrecursor) > limits.MaxReportsPlusPrecursorLength {
		return false
	}

	if len(ccrs.CommitQuorumCertificate) != byzquorum.Size(n, f) {
		return false
	}
	for _, acs := range ccrs.CommitQuorumCertificate {
		if len(acs.Signature) != ed25519.SignatureSize {
			return false
		}
	}
	return true
}

func (ccrs *CertifiedCommittedReports[RI]) Verify(
	configDigest types.ConfigDigest,
	oracleIdentities []config.OracleIdentity,
	byzQuorumSize int,
) error {
	if byzQuorumSize != len(ccrs.CommitQuorumCertificate) {
		return fmt.Errorf("wrong number of signatures, expected %d for byz. quorum but got %d", byzQuorumSize, len(ccrs.CommitQuorumCertificate))
	}

	ogid := OutcomeGenerationID{
		configDigest,
		ccrs.CommitEpoch,
	}

	seen := make(map[commontypes.OracleID]bool)
	for i, acs := range ccrs.CommitQuorumCertificate {
		if seen[acs.Signer] {
			return fmt.Errorf("duplicate signature by %v", acs.Signer)
		}
		seen[acs.Signer] = true
		if !(0 <= int(acs.Signer) && int(acs.Signer) < len(oracleIdentities)) {
			return fmt.Errorf("signer out of bounds: %v", acs.Signer)
		}
		reportsPlusPrecursorDigest := MakeReportsPlusPrecursorDigest(ogid, ccrs.SeqNr, ccrs.ReportsPlusPrecursor)
		if err := acs.Signature.Verify(ogid,
			ccrs.SeqNr,
			ccrs.StateTransitionInputsDigest,
			ccrs.StateTransitionOutputDigest,
			reportsPlusPrecursorDigest,
			oracleIdentities[acs.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle does not verify: %w", i, acs.Signer, err)
		}
	}
	return nil
}

type BlobDigest = blobtypes.BlobDigest
type BlobChunkDigest = blobtypes.BlobChunkDigest
type BlobAvailabilitySignature = blobtypes.BlobAvailabilitySignature
type AttributedBlobAvailabilitySignature = blobtypes.AttributedBlobAvailabilitySignature
type LightCertifiedBlob = blobtypes.LightCertifiedBlob
