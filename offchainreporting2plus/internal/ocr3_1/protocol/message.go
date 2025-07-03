package protocol //

import (
	"crypto/ed25519"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/byzquorum"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Message is the interface used to pass an inter-oracle message to the local
// oracle process.
type Message[RI any] interface {
	// CheckSize checks whether the given message conforms to the limits imposed by
	// reportingPluginLimits
	CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool

	// process passes this Message instance to the oracle o, as a message from
	// oracle with the given sender index
	process(o *oracleState[RI], sender commontypes.OracleID)
}

// MessageWithSender records a msg with the index of the sender oracle
type MessageWithSender[RI any] struct {
	Msg    Message[RI]
	Sender commontypes.OracleID
}

type MessageToPacemaker[RI any] interface {
	Message[RI]

	processPacemaker(pace *pacemakerState[RI], sender commontypes.OracleID)
}

type MessageToPacemakerWithSender[RI any] struct {
	msg    MessageToPacemaker[RI]
	sender commontypes.OracleID
}

type MessageToOutcomeGeneration[RI any] interface {
	Message[RI]

	processOutcomeGeneration(outgen *outcomeGenerationState[RI], sender commontypes.OracleID)

	epoch() uint64
}

type MessageToOutcomeGenerationWithSender[RI any] struct {
	msg    MessageToOutcomeGeneration[RI]
	sender commontypes.OracleID
}

type MessageToReportAttestation[RI any] interface {
	Message[RI]

	processReportAttestation(repatt *reportAttestationState[RI], sender commontypes.OracleID)
}

type MessageToReportAttestationWithSender[RI any] struct {
	msg    MessageToReportAttestation[RI]
	sender commontypes.OracleID
}

type MessageToStatePersistence[RI any] interface {
	Message[RI]

	processStatePersistence(state *statePersistenceState[RI], sender commontypes.OracleID)
}

type MessageToStatePersistenceWithSender[RI any] struct {
	msg    MessageToStatePersistence[RI]
	sender commontypes.OracleID
}

type MessageToBlobExchange[RI any] interface {
	Message[RI]

	processBlobExchange(bex *blobExchangeState[RI], sender commontypes.OracleID)
}

type MessageToBlobExchangeWithSender[RI any] struct {
	msg    MessageToBlobExchange[RI]
	sender commontypes.OracleID
}

type MessageNewEpochWish[RI any] struct {
	Epoch uint64
}

var _ MessageToPacemaker[struct{}] = (*MessageNewEpochWish[struct{}])(nil)

func (msg MessageNewEpochWish[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (msg MessageNewEpochWish[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToPacemaker <- MessageToPacemakerWithSender[RI]{msg, sender}
}

func (msg MessageNewEpochWish[RI]) processPacemaker(pace *pacemakerState[RI], sender commontypes.OracleID) {
	pace.messageNewEpochWish(msg, sender)
}

type MessageEpochStartRequest[RI any] struct {
	Epoch                           uint64
	HighestCertified                CertifiedPrepareOrCommit
	SignedHighestCertifiedTimestamp SignedHighestCertifiedTimestamp
}

var _ MessageToOutcomeGeneration[struct{}] = (*MessageEpochStartRequest[struct{}])(nil)

func (msg MessageEpochStartRequest[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	if !msg.HighestCertified.CheckSize(n, f, limits, maxReportSigLen) {
		return false
	}
	if len(msg.SignedHighestCertifiedTimestamp.Signature) != ed25519.SignatureSize {
		return false
	}
	return true
}

func (msg MessageEpochStartRequest[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToOutcomeGeneration <- MessageToOutcomeGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageEpochStartRequest[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI], sender commontypes.OracleID) {
	outgen.messageEpochStartRequest(msg, sender)
}

func (msg MessageEpochStartRequest[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessageEpochStart[RI any] struct {
	Epoch           uint64
	EpochStartProof EpochStartProof
}

var _ MessageToOutcomeGeneration[struct{}] = (*MessageEpochStart[struct{}])(nil)

func (msg MessageEpochStart[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	if !msg.EpochStartProof.HighestCertified.CheckSize(n, f, limits, maxReportSigLen) {
		return false
	}
	if len(msg.EpochStartProof.HighestCertifiedProof) != byzquorum.Size(n, f) {
		return false
	}
	for _, ashct := range msg.EpochStartProof.HighestCertifiedProof {
		if len(ashct.SignedHighestCertifiedTimestamp.Signature) != ed25519.SignatureSize {
			return false
		}
	}
	return true
}

func (msg MessageEpochStart[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToOutcomeGeneration <- MessageToOutcomeGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageEpochStart[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI], sender commontypes.OracleID) {
	outgen.messageEpochStart(msg, sender)
}

func (msg MessageEpochStart[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessageRoundStart[RI any] struct {
	Epoch uint64
	SeqNr uint64
	Query types.Query
}

var _ MessageToOutcomeGeneration[struct{}] = (*MessageRoundStart[struct{}])(nil)

func (msg MessageRoundStart[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	return len(msg.Query) <= limits.MaxQueryLength
}

func (msg MessageRoundStart[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToOutcomeGeneration <- MessageToOutcomeGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageRoundStart[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI], sender commontypes.OracleID) {
	outgen.messageRoundStart(msg, sender)
}

func (msg MessageRoundStart[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessageObservation[RI any] struct {
	Epoch             uint64
	SeqNr             uint64
	SignedObservation SignedObservation
}

var _ MessageToOutcomeGeneration[struct{}] = (*MessageObservation[struct{}])(nil)

func (msg MessageObservation[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	return len(msg.SignedObservation.Observation) <= limits.MaxObservationLength && len(msg.SignedObservation.Signature) == ed25519.SignatureSize
}

func (msg MessageObservation[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToOutcomeGeneration <- MessageToOutcomeGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageObservation[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI], sender commontypes.OracleID) {
	outgen.messageObservation(msg, sender)
}

func (msg MessageObservation[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessageProposal[RI any] struct {
	Epoch                        uint64
	SeqNr                        uint64
	AttributedSignedObservations []AttributedSignedObservation
}

var _ MessageToOutcomeGeneration[struct{}] = MessageProposal[struct{}]{}

func (msg MessageProposal[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	if len(msg.AttributedSignedObservations) > n {
		return false
	}
	for _, aso := range msg.AttributedSignedObservations {
		if len(aso.SignedObservation.Observation) > limits.MaxObservationLength {
			return false
		}
		if len(aso.SignedObservation.Signature) != ed25519.SignatureSize {
			return false
		}
	}
	return true
}

func (msg MessageProposal[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToOutcomeGeneration <- MessageToOutcomeGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageProposal[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI], sender commontypes.OracleID) {
	outgen.messageProposal(msg, sender)
}

func (msg MessageProposal[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessagePrepare[RI any] struct {
	Epoch     uint64
	SeqNr     uint64
	Signature PrepareSignature
}

var _ MessageToOutcomeGeneration[struct{}] = MessagePrepare[struct{}]{}

func (msg MessagePrepare[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	return len(msg.Signature) == ed25519.SignatureSize
}

func (msg MessagePrepare[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToOutcomeGeneration <- MessageToOutcomeGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessagePrepare[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI], sender commontypes.OracleID) {
	outgen.messagePrepare(msg, sender)
}

func (msg MessagePrepare[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessageCommit[RI any] struct {
	Epoch     uint64
	SeqNr     uint64
	Signature CommitSignature
}

var _ MessageToOutcomeGeneration[struct{}] = MessageCommit[struct{}]{}

func (msg MessageCommit[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	return len(msg.Signature) == ed25519.SignatureSize
}

func (msg MessageCommit[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToOutcomeGeneration <- MessageToOutcomeGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageCommit[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI], sender commontypes.OracleID) {
	outgen.messageCommit(msg, sender)
}

func (msg MessageCommit[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessageReportSignatures[RI any] struct {
	SeqNr            uint64
	ReportSignatures [][]byte
}

var _ MessageToReportAttestation[struct{}] = MessageReportSignatures[struct{}]{}

func (msg MessageReportSignatures[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	if len(msg.ReportSignatures) > limits.MaxReportCount {
		return false
	}
	for _, sig := range msg.ReportSignatures {
		if len(sig) > maxReportSigLen {
			return false
		}
	}

	return true
}

func (msg MessageReportSignatures[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportAttestation <- MessageToReportAttestationWithSender[RI]{msg, sender}
}

func (msg MessageReportSignatures[RI]) processReportAttestation(repatt *reportAttestationState[RI], sender commontypes.OracleID) {
	repatt.messageReportSignatures(msg, sender)
}

type MessageCertifiedCommitRequest[RI any] struct {
	SeqNr uint64
}

var _ MessageToReportAttestation[struct{}] = MessageCertifiedCommitRequest[struct{}]{}

func (msg MessageCertifiedCommitRequest[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	return true
}

func (msg MessageCertifiedCommitRequest[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportAttestation <- MessageToReportAttestationWithSender[RI]{msg, sender}
}

func (msg MessageCertifiedCommitRequest[RI]) processReportAttestation(repatt *reportAttestationState[RI], sender commontypes.OracleID) {
	repatt.messageCertifiedCommitRequest(msg, sender)
}

type MessageCertifiedCommit[RI any] struct {
	CertifiedCommittedReports CertifiedCommittedReports[RI]
}

var _ MessageToReportAttestation[struct{}] = MessageCertifiedCommit[struct{}]{}

func (msg MessageCertifiedCommit[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	return msg.CertifiedCommittedReports.CheckSize(n, f, limits, maxReportSigLen)
}

func (msg MessageCertifiedCommit[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportAttestation <- MessageToReportAttestationWithSender[RI]{msg, sender}
}

func (msg MessageCertifiedCommit[RI]) processReportAttestation(repatt *reportAttestationState[RI], sender commontypes.OracleID) {
	repatt.messageCertifiedCommit(msg, sender)
}

type MessageBlockSyncRequest[RI any] struct {
	RequestHandle         types.RequestHandle // actual handle for outbound message, sentinel for inbound
	HighestCommittedSeqNr uint64
	Nonce                 uint64
}

var _ MessageToStatePersistence[struct{}] = MessageBlockSyncRequest[struct{}]{}

func (msg MessageBlockSyncRequest[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (msg MessageBlockSyncRequest[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStatePersistence <- MessageToStatePersistenceWithSender[RI]{msg, sender}
}

func (msg MessageBlockSyncRequest[RI]) processStatePersistence(state *statePersistenceState[RI], sender commontypes.OracleID) {
	state.messageBlockSyncReq(msg, sender)
}

type MessageBlockSyncSummary[RI any] struct {
	LowestPersistedSeqNr uint64
}

var _ MessageToStatePersistence[struct{}] = MessageBlockSyncSummary[struct{}]{}

func (msg MessageBlockSyncSummary[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (msg MessageBlockSyncSummary[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStatePersistence <- MessageToStatePersistenceWithSender[RI]{msg, sender}
}

func (msg MessageBlockSyncSummary[RI]) processStatePersistence(state *statePersistenceState[RI], sender commontypes.OracleID) {
	state.messageBlockSyncSummary(msg, sender)
}

type MessageBlockSync[RI any] struct {
	RequestHandle                 types.RequestHandle // actual handle for outbound message, sentinel for inbound
	AttestedStateTransitionBlocks []AttestedStateTransitionBlock
	Nonce                         uint64
}

var _ MessageToStatePersistence[struct{}] = MessageBlockSync[struct{}]{}

func (msg MessageBlockSync[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (msg MessageBlockSync[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStatePersistence <- MessageToStatePersistenceWithSender[RI]{msg, sender}
}

func (msg MessageBlockSync[RI]) processStatePersistence(state *statePersistenceState[RI], sender commontypes.OracleID) {
	state.messageBlockSync(msg, sender)
}

type MessageBlobOffer[RI any] struct {
	ChunkDigests  []BlobChunkDigest
	PayloadLength uint64
	ExpirySeqNr   uint64
	Submitter     commontypes.OracleID
}

var _ MessageToBlobExchange[struct{}] = MessageBlobOffer[struct{}]{}

func (msg MessageBlobOffer[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true // TODO: add proper size checks
}

func (msg MessageBlobOffer[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToBlobExchange <- MessageToBlobExchangeWithSender[RI]{msg, sender}
}

func (msg MessageBlobOffer[RI]) processBlobExchange(bex *blobExchangeState[RI], sender commontypes.OracleID) {
	bex.messageBlobOffer(msg, sender)
}

type MessageBlobChunkRequest[RI any] struct {
	RequestHandle types.RequestHandle // actual handle for outbound message, sentinel for inbound
	BlobDigest    BlobDigest
	ChunkIndex    uint64
}

var _ MessageToBlobExchange[struct{}] = MessageBlobChunkRequest[struct{}]{}

func (msg MessageBlobChunkRequest[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true // TODO: add proper size checks
}

func (msg MessageBlobChunkRequest[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToBlobExchange <- MessageToBlobExchangeWithSender[RI]{msg, sender}
}

func (msg MessageBlobChunkRequest[RI]) processBlobExchange(bex *blobExchangeState[RI], sender commontypes.OracleID) {
	bex.messageBlobChunkRequest(msg, sender)
}

type MessageBlobChunkResponse[RI any] struct {
	RequestHandle types.RequestHandle // actual handle for outbound message, sentinel for inbound

	BlobDigest BlobDigest
	ChunkIndex uint64
	Chunk      []byte
}

var _ MessageToBlobExchange[struct{}] = MessageBlobChunkResponse[struct{}]{}

func (msg MessageBlobChunkResponse[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true // TODO: add proper size checks
}

func (msg MessageBlobChunkResponse[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToBlobExchange <- MessageToBlobExchangeWithSender[RI]{msg, sender}
}

func (msg MessageBlobChunkResponse[RI]) processBlobExchange(bex *blobExchangeState[RI], sender commontypes.OracleID) {
	bex.messageBlobChunkResponse(msg, sender)
}

type MessageBlobAvailable[RI any] struct {
	BlobDigest BlobDigest
	Signature  BlobAvailabilitySignature
}

var _ MessageToBlobExchange[struct{}] = MessageBlobAvailable[struct{}]{}

func (msg MessageBlobAvailable[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true // TODO: add proper size checks
}

func (msg MessageBlobAvailable[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToBlobExchange <- MessageToBlobExchangeWithSender[RI]{msg, sender}
}

func (msg MessageBlobAvailable[RI]) processBlobExchange(bex *blobExchangeState[RI], sender commontypes.OracleID) {
	bex.messageBlobAvailable(msg, sender)
}
