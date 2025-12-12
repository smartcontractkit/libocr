package protocol

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/byzquorum"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3_1/blobtypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Returns a byte slice whose first six bytes are the string "ocr3.1" and the rest
// of which is the sum returned by h. Used for domain separation as per the comment
// on offchainreporting2plus/types.OffchainKeyring.
//
// Any signatures made with the OffchainKeyring should use ocr3_1DomainSeparatedSum!
func ocr3_1DomainSeparatedSum(h hash.Hash) []byte {
	const domainSeparator = "ocr3.1"
	result := make([]byte, 0, len(domainSeparator)+sha256.Size)
	result = append(result, []byte(domainSeparator)...)
	return h.Sum(result)
}

const signedObservationDomainSeparator = "ocr3.1/SignedObservation/"

type SignedObservation struct {
	Observation types.Observation
	Signature   []byte
}

func MakeSignedObservation(
	ogid OutcomeGenerationID,
	seqNr uint64,
	aq types.AttributedQuery,
	observation types.Observation,
	signer func(msg []byte) (sig []byte, err error),
) (
	SignedObservation,
	error,
) {
	payload := signedObservationMsg(ogid, seqNr, aq, observation)
	sig, err := signer(payload)
	if err != nil {
		return SignedObservation{}, err
	}
	return SignedObservation{observation, sig}, nil
}

func (so SignedObservation) Verify(ogid OutcomeGenerationID, seqNr uint64, aq types.AttributedQuery, publicKey types.OffchainPublicKey) error {
	pk := ed25519.PublicKey(publicKey[:])
	// should never trigger since types.OffchainPublicKey is an array with length ed25519.PublicKeySize
	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, signedObservationMsg(ogid, seqNr, aq, so.Observation), so.Signature)
	if !ok {
		return fmt.Errorf("SignedObservation has invalid signature")
	}

	return nil
}

func signedObservationMsg(ogid OutcomeGenerationID, seqNr uint64, attributedQuery types.AttributedQuery, observation types.Observation) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(signedObservationDomainSeparator))

	// ogid
	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	// seqNr
	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, seqNr)

	// attributedQuery.Query
	_ = binary.Write(h, binary.BigEndian, uint64(len(attributedQuery.Query)))
	_, _ = h.Write(attributedQuery.Query)

	// attributedQuery.Proposer
	_ = binary.Write(h, binary.BigEndian, uint64(attributedQuery.Proposer))

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
	configDigest types.ConfigDigest,
	seqNr uint64,
	attributedQuery types.AttributedQuery,
	attributedObservations []types.AttributedObservation,
) StateTransitionInputsDigest {
	h := sha256.New()

	_, _ = h.Write(configDigest[:])

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(attributedQuery.Query)))
	_, _ = h.Write(attributedQuery.Query)

	_ = binary.Write(h, binary.BigEndian, uint64(attributedQuery.Proposer))

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

type StateWriteSetDigest [32]byte

func MakeStateWriteSetDigest(configDigest types.ConfigDigest, seqNr uint64, writeSet []KeyValuePairWithDeletions) StateWriteSetDigest {
	h := sha256.New()

	_, _ = h.Write(configDigest[:])

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(writeSet)))
	for _, o := range writeSet {

		_ = binary.Write(h, binary.BigEndian, uint64(len(o.Key)))
		_, _ = h.Write(o.Key)

		_ = binary.Write(h, binary.BigEndian, uint64(len(o.Value)))
		_, _ = h.Write(o.Value)
	}

	var result StateWriteSetDigest
	h.Sum(result[:0])
	return result
}

type ReportsPlusPrecursorDigest [32]byte

func MakeReportsPlusPrecursorDigest(configDigest types.ConfigDigest, seqNr uint64, precursor ocr3_1types.ReportsPlusPrecursor) ReportsPlusPrecursorDigest {
	h := sha256.New()

	_, _ = h.Write(configDigest[:])

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_ = binary.Write(h, binary.BigEndian, uint64(len(precursor)))
	_, _ = h.Write(precursor)

	var result ReportsPlusPrecursorDigest
	h.Sum(result[:0])
	return result
}

const prepareSignatureDomainSeparator = "ocr3.1/PrepareSignature/"

type PrepareSignature []byte

func MakePrepareSignature(
	ogid OutcomeGenerationID,
	prevHistoryDigest HistoryDigest,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	writeSetDigest StateWriteSetDigest,
	rootDigest StateRootDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	signer func(msg []byte) ([]byte, error),
) (PrepareSignature, error) {
	return signer(prepareSignatureMsg(ogid, prevHistoryDigest, seqNr, inputsDigest, writeSetDigest, rootDigest, reportsPlusPrecursorDigest))
}

func (sig PrepareSignature) Verify(
	ogid OutcomeGenerationID,
	prevHistoryDigest HistoryDigest,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	writeSetDigest StateWriteSetDigest,
	rootDigest StateRootDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}
	msg := prepareSignatureMsg(ogid, prevHistoryDigest, seqNr, inputsDigest, writeSetDigest, rootDigest, reportsPlusPrecursorDigest)
	ok := ed25519.Verify(pk, msg, sig)
	if !ok {
		// Other less common causes include leader equivocation or actually invalid signatures.
		return fmt.Errorf("PrepareSignature failed to verify. This is commonly caused by non-determinism in the ReportingPlugin msg: %x, sig: %x", msg, sig)
	}

	return nil
}

func prepareSignatureMsg(
	ogid OutcomeGenerationID,
	prevHistoryDigest HistoryDigest,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	writeSetDigest StateWriteSetDigest,
	rootDigest StateRootDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(prepareSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	historyDigest := MakeHistoryDigest(
		ogid.ConfigDigest,
		prevHistoryDigest,
		seqNr,
		inputsDigest,
		writeSetDigest,
		rootDigest,
		reportsPlusPrecursorDigest,
	)
	_, _ = h.Write(historyDigest[:])

	return ocr3_1DomainSeparatedSum(h)
}

type AttributedPrepareSignature struct {
	Signature PrepareSignature
	Signer    commontypes.OracleID
}

const commitSignatureDomainSeparator = "ocr3.1/CommitSignature/"

type CommitSignature []byte

func MakeCommitSignature(
	ogid OutcomeGenerationID,
	prevHistoryDigest HistoryDigest,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	writeSetDigest StateWriteSetDigest,
	rootDigest StateRootDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	signer func(msg []byte) ([]byte, error),
) (CommitSignature, error) {
	return signer(commitSignatureMsg(ogid, prevHistoryDigest, seqNr, inputsDigest, writeSetDigest, rootDigest, reportsPlusPrecursorDigest))
}

func (sig CommitSignature) Verify(
	ogid OutcomeGenerationID,
	prevHistoryDigest HistoryDigest,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	writeSetDigest StateWriteSetDigest,
	rootDigest StateRootDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
	publicKey types.OffchainPublicKey,
) error {
	pk := ed25519.PublicKey(publicKey[:])

	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, commitSignatureMsg(ogid, prevHistoryDigest, seqNr, inputsDigest, writeSetDigest, rootDigest, reportsPlusPrecursorDigest), sig)
	if !ok {
		return fmt.Errorf("CommitSignature failed to verify")
	}

	return nil
}

type HistoryDigest = types.HistoryDigest

const historyDigestDomainSeparator = "ocr3.1/HistoryDigest/"

func MakeHistoryDigest(
	configDigest types.ConfigDigest,
	prevHistoryDigest HistoryDigest,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	writeSetDigest StateWriteSetDigest,
	rootDigest StateRootDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
) HistoryDigest {
	h := sha256.New()

	_, _ = h.Write([]byte(historyDigestDomainSeparator))

	_, _ = h.Write(configDigest[:])

	_, _ = h.Write(prevHistoryDigest[:])

	_ = binary.Write(h, binary.BigEndian, seqNr)

	_, _ = h.Write(inputsDigest[:])

	_, _ = h.Write(writeSetDigest[:])

	_, _ = h.Write(rootDigest[:])

	_, _ = h.Write(reportsPlusPrecursorDigest[:])

	var result HistoryDigest
	h.Sum(result[:0])
	return result
}

func commitSignatureMsg(
	ogid OutcomeGenerationID,
	prevHistoryDigest HistoryDigest,
	seqNr uint64,
	inputsDigest StateTransitionInputsDigest,
	writeSetDigest StateWriteSetDigest,
	rootDigest StateRootDigest,
	reportsPlusPrecursorDigest ReportsPlusPrecursorDigest,
) []byte {
	h := sha256.New()

	_, _ = h.Write([]byte(commitSignatureDomainSeparator))

	_, _ = h.Write(ogid.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, ogid.Epoch)

	historyDigest := MakeHistoryDigest(
		ogid.ConfigDigest,
		prevHistoryDigest,
		seqNr,
		inputsDigest,
		writeSetDigest,
		rootDigest,
		reportsPlusPrecursorDigest,
	)
	_, _ = h.Write(historyDigest[:])

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

const signedHighestCertifiedTimestampDomainSeparator = "ocr3.1/SignedHighestCertifiedTimestamp/"

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
	config ocr3_1config.PublicConfig,
) error {
	oracleIdentities := config.OracleIdentities
	byzQuorumSize := config.ByzQuorumSize()
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

	if err := qc.HighestCertified.Verify(config); err != nil {
		return fmt.Errorf("failed to verify HighestCertified: %w", err)
	}

	return nil
}

//go-sumtype:decl CertifiedPrepareOrCommit

type CertifiedPrepareOrCommit interface {
	isCertifiedPrepareOrCommit()
	Epoch() uint64
	SeqNr() uint64
	Timestamp() HighestCertifiedTimestamp
	HistoryDigest(_ types.ConfigDigest) HistoryDigest
	Verify(_ ocr3_1config.PublicConfig) error
	CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool
}

var _ CertifiedPrepareOrCommit = &CertifiedPrepare{}

type CertifiedPrepare struct {
	PrevHistoryDigest           HistoryDigest
	PrepareEpoch                uint64
	PrepareSeqNr                uint64
	StateTransitionInputsDigest StateTransitionInputsDigest
	StateWriteSetDigest         StateWriteSetDigest
	StateRootDigest             StateRootDigest
	ReportsPlusPrecursorDigest  ReportsPlusPrecursorDigest
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

func (hc *CertifiedPrepare) HistoryDigest(configDigest types.ConfigDigest) HistoryDigest {
	return MakeHistoryDigest(
		configDigest,
		hc.PrevHistoryDigest,
		hc.PrepareSeqNr,
		hc.StateTransitionInputsDigest,
		hc.StateWriteSetDigest,
		hc.StateRootDigest,
		hc.ReportsPlusPrecursorDigest,
	)
}

func (hc *CertifiedPrepare) Verify(config ocr3_1config.PublicConfig) error {
	configDigest := config.ConfigDigest
	oracleIdentities := config.OracleIdentities
	byzQuorumSize := config.ByzQuorumSize()
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
		if err := aps.Signature.Verify(
			ogid, hc.PrevHistoryDigest,
			hc.SeqNr(), hc.StateTransitionInputsDigest, hc.StateWriteSetDigest, hc.StateRootDigest,
			hc.ReportsPlusPrecursorDigest, oracleIdentities[aps.Signer].OffchainPublicKey); err != nil {
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

type CertifiedCommit struct {
	PrevHistoryDigest           HistoryDigest
	CommitEpoch                 uint64
	CommitSeqNr                 uint64
	StateTransitionInputsDigest StateTransitionInputsDigest
	StateWriteSetDigest         StateWriteSetDigest
	StateRootDigest             StateRootDigest
	ReportsPlusPrecursorDigest  ReportsPlusPrecursorDigest
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

func (hc *CertifiedCommit) Verify(config ocr3_1config.PublicConfig) error {
	configDigest := config.ConfigDigest
	oracleIdentities := config.OracleIdentities
	byzQuorumSize := config.ByzQuorumSize()
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
		if err := acs.Signature.Verify(
			ogid,
			hc.PrevHistoryDigest,
			hc.SeqNr(),
			hc.StateTransitionInputsDigest,
			hc.StateWriteSetDigest,
			hc.StateRootDigest,
			hc.ReportsPlusPrecursorDigest,
			oracleIdentities[acs.Signer].OffchainPublicKey); err != nil {
			return fmt.Errorf("%v-th signature by %v-th oracle does not verify: %w", i, acs.Signer, err)
		}
	}
	return nil
}

func (hc *CertifiedCommit) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
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

func (hc *CertifiedCommit) HistoryDigest(configDigest types.ConfigDigest) HistoryDigest {
	return MakeHistoryDigest(
		configDigest,
		hc.PrevHistoryDigest,
		hc.SeqNr(),
		hc.StateTransitionInputsDigest,
		hc.StateWriteSetDigest,
		hc.StateRootDigest,
		hc.ReportsPlusPrecursorDigest,
	)
}

type GenesisFromScratch struct{}

var _ CertifiedPrepareOrCommit = &GenesisFromScratch{}

func (gscratch *GenesisFromScratch) CheckSize(_ int, _ int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (gscratch *GenesisFromScratch) Epoch() uint64 {
	return 0
}

func (gscratch *GenesisFromScratch) SeqNr() uint64 {
	return 0
}

func (gscratch *GenesisFromScratch) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		0,
		true,
		0,
	}
}

func (gscratch *GenesisFromScratch) Verify(config ocr3_1config.PublicConfig) error {
	if _, ok := config.GetPrevFields(); ok {
		return fmt.Errorf("GenesisFromScratch is invalid when previous instance is specified in PublicConfig")
	}
	return nil
}

func (gscratch *GenesisFromScratch) HistoryDigest(configDigest types.ConfigDigest) HistoryDigest {
	return MakeHistoryDigest(
		configDigest,
		HistoryDigest{},
		0,
		StateTransitionInputsDigest{},
		StateWriteSetDigest{},
		StateRootDigest{},
		ReportsPlusPrecursorDigest{},
	)
}

func (gscratch *GenesisFromScratch) isCertifiedPrepareOrCommit() {}

type GenesisFromPrevInstance struct {
	PrevHistoryDigest HistoryDigest
	PrevSeqNr         uint64
}

var _ CertifiedPrepareOrCommit = &GenesisFromPrevInstance{}

func (gprev *GenesisFromPrevInstance) CheckSize(_ int, _ int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (gprev *GenesisFromPrevInstance) Epoch() uint64 {
	return 0
}

func (gprev *GenesisFromPrevInstance) SeqNr() uint64 {
	return gprev.PrevSeqNr
}

func (gprev *GenesisFromPrevInstance) Timestamp() HighestCertifiedTimestamp {
	return HighestCertifiedTimestamp{
		gprev.SeqNr(),
		true,
		gprev.Epoch(),
	}
}

func (gprev *GenesisFromPrevInstance) Verify(config ocr3_1config.PublicConfig) error {
	prev, ok := config.GetPrevFields()
	if !ok {
		return fmt.Errorf("GenesisFromPrevInstance is invalid when previous instance is not specified in PublicConfig")
	}
	if gprev.PrevHistoryDigest != prev.PrevHistoryDigest {
		return fmt.Errorf("GenesisFromPrevInstance.PrevHistoryDigest (%x) does not match config.PrevHistoryDigest (%x)", gprev.PrevHistoryDigest, prev.PrevHistoryDigest)
	}
	if gprev.PrevSeqNr != prev.PrevSeqNr {
		return fmt.Errorf("GenesisFromPrevInstance.PrevSeqNr (%d) does not match config.PrevSeqNr (%d)", gprev.PrevSeqNr, prev.PrevSeqNr)
	}
	return nil
}

func (gprev *GenesisFromPrevInstance) HistoryDigest(configDigest types.ConfigDigest) HistoryDigest {
	return MakeHistoryDigest(
		configDigest,
		gprev.PrevHistoryDigest,
		gprev.SeqNr(),
		StateTransitionInputsDigest{},
		StateWriteSetDigest{},
		StateRootDigest{},
		ReportsPlusPrecursorDigest{},
	)
}

func (gprev *GenesisFromPrevInstance) isCertifiedPrepareOrCommit() {}

func GenesisCertifiedPrepareOrCommit(cfg ocr3_1config.PublicConfig) CertifiedPrepareOrCommit {
	prev, ok := cfg.GetPrevFields()
	if !ok {
		return &GenesisFromScratch{}
	}
	return &GenesisFromPrevInstance{
		prev.PrevHistoryDigest,
		prev.PrevSeqNr,
	}
}

type StateTransitionBlock struct {
	PrevHistoryDigest           HistoryDigest
	Epoch                       uint64
	BlockSeqNr                  uint64
	StateTransitionInputsDigest StateTransitionInputsDigest
	StateWriteSet               StateWriteSet
	StateRootDigest             StateRootDigest
	ReportsPlusPrecursorDigest  ReportsPlusPrecursorDigest
}

func (stb *StateTransitionBlock) SeqNr() uint64 {
	return stb.BlockSeqNr
}

func checkWriteSetSize(writeSet []KeyValuePairWithDeletions, limits ocr3_1types.ReportingPluginLimits) bool {
	if len(writeSet) > limits.MaxKeyValueModifiedKeys {
		return false
	}

	modifiedKeysPlusValuesLength := 0
	for _, kvPair := range writeSet {
		if len(kvPair.Key) > ocr3_1types.MaxMaxKeyValueKeyBytes {
			return false
		}
		if kvPair.Deleted && len(kvPair.Value) > 0 {
			return false
		}
		if len(kvPair.Value) > ocr3_1types.MaxMaxKeyValueValueBytes {
			return false
		}
		modifiedKeysPlusValuesLength += len(kvPair.Key) + len(kvPair.Value)
	}
	if modifiedKeysPlusValuesLength > limits.MaxKeyValueModifiedKeysPlusValuesBytes {
		return false
	}
	return true
}

func (stb *StateTransitionBlock) CheckSize(limits ocr3_1types.ReportingPluginLimits) bool {
	if !checkWriteSetSize(stb.StateWriteSet.Entries, limits) {
		return false
	}
	return true
}

type AttestedStateTransitionBlock struct {
	StateTransitionBlock       StateTransitionBlock
	AttributedCommitSignatures []AttributedCommitSignature
}

func (astb *AttestedStateTransitionBlock) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits) bool {
	if len(astb.AttributedCommitSignatures) != byzquorum.Size(n, f) {
		return false
	}
	for _, acs := range astb.AttributedCommitSignatures {
		if len(acs.Signature) != ed25519.SignatureSize {
			return false
		}
	}
	return astb.StateTransitionBlock.CheckSize(limits)
}

func (astb *AttestedStateTransitionBlock) Verify(config ocr3_1config.PublicConfig) error {
	certifiedCommit := astb.ToCertifiedCommit(config.ConfigDigest)
	return certifiedCommit.Verify(config)
}

func (astb *AttestedStateTransitionBlock) ToCertifiedCommit(configDigest types.ConfigDigest) CertifiedCommit {
	stb := astb.StateTransitionBlock
	stateWriteSetDigest := MakeStateWriteSetDigest(
		configDigest,
		stb.SeqNr(),
		stb.StateWriteSet.Entries,
	)
	return CertifiedCommit{
		stb.PrevHistoryDigest,
		stb.Epoch,
		stb.SeqNr(),
		stb.StateTransitionInputsDigest,
		stateWriteSetDigest,
		stb.StateRootDigest,
		stb.ReportsPlusPrecursorDigest,
		astb.AttributedCommitSignatures,
	}
}

type GenesisStateTransitionBlock struct {
	PrevHistoryDigest           HistoryDigest
	SeqNr                       uint64
	StateTransitionInputsDigest StateTransitionInputsDigest
	StateWriteSetDigest         StateWriteSetDigest
	StateRootDigest             StateRootDigest
	ReportsPlusPrecursorDigest  ReportsPlusPrecursorDigest
}

//go-sumtype:decl AttestedOrGenesisStateTransitionBlock

type AttestedOrGenesisStateTransitionBlock interface {
	isAttestedOrGenesisStateTransitionBlock()
	seqNr() uint64
	stateRootDigest() StateRootDigest
}

func (astb *AttestedStateTransitionBlock) isAttestedOrGenesisStateTransitionBlock() {}
func (astb *AttestedStateTransitionBlock) seqNr() uint64 {
	return astb.StateTransitionBlock.SeqNr()
}
func (astb *AttestedStateTransitionBlock) stateRootDigest() StateRootDigest {
	return astb.StateTransitionBlock.StateRootDigest
}

func (gstb *GenesisStateTransitionBlock) isAttestedOrGenesisStateTransitionBlock() {}
func (gstb *GenesisStateTransitionBlock) seqNr() uint64 {
	return gstb.SeqNr
}
func (gstb *GenesisStateTransitionBlock) stateRootDigest() StateRootDigest {
	return gstb.StateRootDigest
}

type BlobDigest = blobtypes.BlobDigest
type BlobChunkDigest = blobtypes.BlobChunkDigest
type BlobChunkDigestsRoot = blobtypes.BlobChunkDigestsRoot
type BlobAvailabilitySignature = blobtypes.BlobAvailabilitySignature
type AttributedBlobAvailabilitySignature = blobtypes.AttributedBlobAvailabilitySignature
type LightCertifiedBlob = blobtypes.LightCertifiedBlob
