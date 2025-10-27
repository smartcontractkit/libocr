package protocol

import (
	"crypto/ed25519"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/byzquorum"
	"github.com/smartcontractkit/libocr/internal/jmt"
	"github.com/smartcontractkit/libocr/internal/mt"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3_1config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

//go-sumtype:decl Message

// Message is the interface used to pass an inter-oracle message to the local
// oracle process.
type Message[RI any] interface {
	// CheckSize checks whether the given message conforms to the limits imposed by
	// reportingPluginLimits
	CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, config ocr3_1config.PublicConfig) bool

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

type MessageToStateSync[RI any] interface {
	Message[RI]

	processStateSync(stasy *stateSyncState[RI], sender commontypes.OracleID)
}

type MessageToStateSyncWithSender[RI any] struct {
	msg    MessageToStateSync[RI]
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

func (msg MessageNewEpochWish[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int, _ ocr3_1config.PublicConfig) bool {
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

func (msg MessageEpochStartRequest[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
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
	Abdicate        bool
}

var _ MessageToOutcomeGeneration[struct{}] = (*MessageEpochStart[struct{}])(nil)

func (msg MessageEpochStart[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
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

func (msg MessageRoundStart[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
	return len(msg.Query) <= limits.MaxQueryBytes
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

func (msg MessageObservation[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
	return len(msg.SignedObservation.Observation) <= limits.MaxObservationBytes && len(msg.SignedObservation.Signature) == ed25519.SignatureSize
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

func (msg MessageProposal[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
	if len(msg.AttributedSignedObservations) > n {
		return false
	}
	for _, aso := range msg.AttributedSignedObservations {
		if len(aso.SignedObservation.Observation) > limits.MaxObservationBytes {
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

func (msg MessagePrepare[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
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

func (msg MessageCommit[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
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
	SeqNr                      uint64
	ReportSignatures           [][]byte
	ReportsPlusPrecursorDigest ReportsPlusPrecursorDigest
}

var _ MessageToReportAttestation[struct{}] = MessageReportSignatures[struct{}]{}

func (msg MessageReportSignatures[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
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

type MessageReportsPlusPrecursorRequest[RI any] struct {
	SeqNr uint64
}

var _ MessageToReportAttestation[struct{}] = MessageReportsPlusPrecursorRequest[struct{}]{}

func (msg MessageReportsPlusPrecursorRequest[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
	return true
}

func (msg MessageReportsPlusPrecursorRequest[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportAttestation <- MessageToReportAttestationWithSender[RI]{msg, sender}
}

func (msg MessageReportsPlusPrecursorRequest[RI]) processReportAttestation(repatt *reportAttestationState[RI], sender commontypes.OracleID) {
	repatt.messageReportsPlusPrecursorRequest(msg, sender)
}

type MessageReportsPlusPrecursor[RI any] struct {
	SeqNr                uint64
	ReportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor
}

var _ MessageToReportAttestation[struct{}] = MessageReportsPlusPrecursor[struct{}]{}

func (msg MessageReportsPlusPrecursor[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
	if len(msg.ReportsPlusPrecursor) > limits.MaxReportsPlusPrecursorBytes {
		return false
	}
	return true
}

func (msg MessageReportsPlusPrecursor[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportAttestation <- MessageToReportAttestationWithSender[RI]{msg, sender}
}

func (msg MessageReportsPlusPrecursor[RI]) processReportAttestation(repatt *reportAttestationState[RI], sender commontypes.OracleID) {
	repatt.messageReportsPlusPrecursor(msg, sender)
}

type MessageBlockSyncRequest[RI any] struct {
	RequestHandle types.RequestHandle // actual handle for outbound message, sentinel for inbound
	RequestInfo   *types.RequestInfo
	StartSeqNr    uint64 // a successful response must contain at least the block with this sequence number
	EndExclSeqNr  uint64 // the response may only contain sequence numbers less than this
}

var _ MessageToStateSync[struct{}] = MessageBlockSyncRequest[struct{}]{}

func (msg MessageBlockSyncRequest[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int, _ ocr3_1config.PublicConfig) bool {
	return true
}

func (msg MessageBlockSyncRequest[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStateSync <- MessageToStateSyncWithSender[RI]{msg, sender}
}

func (msg MessageBlockSyncRequest[RI]) processStateSync(stasy *stateSyncState[RI], sender commontypes.OracleID) {
	stasy.messageBlockSyncRequest(msg, sender)
}

type MessageStateSyncSummary[RI any] struct {
	LowestPersistedSeqNr  uint64
	HighestCommittedSeqNr uint64
}

var _ MessageToStateSync[struct{}] = MessageStateSyncSummary[struct{}]{}

func (msg MessageStateSyncSummary[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int, _ ocr3_1config.PublicConfig) bool {
	return true
}

func (msg MessageStateSyncSummary[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStateSync <- MessageToStateSyncWithSender[RI]{msg, sender}
}

func (msg MessageStateSyncSummary[RI]) processStateSync(stasy *stateSyncState[RI], sender commontypes.OracleID) {
	stasy.messageStateSyncSummary(msg, sender)
}

type MessageBlockSyncResponse[RI any] struct {
	RequestHandle                 types.RequestHandle // actual handle for outbound message, sentinel for inbound
	RequestStartSeqNr             uint64
	RequestEndExclSeqNr           uint64
	AttestedStateTransitionBlocks []AttestedStateTransitionBlock // must be contiguous and (if non-empty) starting at RequestStartSeqNr
}

var _ MessageToStateSync[struct{}] = MessageBlockSyncResponse[struct{}]{}

func (msg MessageBlockSyncResponse[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, config ocr3_1config.PublicConfig) bool {
	if len(msg.AttestedStateTransitionBlocks) > config.GetMaxBlocksPerBlockSyncResponse() {
		return false
	}
	for _, astb := range msg.AttestedStateTransitionBlocks {
		if !astb.CheckSize(n, f, limits) {
			return false
		}
	}
	return true
}

func (msg MessageBlockSyncResponse[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStateSync <- MessageToStateSyncWithSender[RI]{msg, sender}
}

func (msg MessageBlockSyncResponse[RI]) processStateSync(stasy *stateSyncState[RI], sender commontypes.OracleID) {
	stasy.messageBlockSyncResponse(msg, sender)
}

type MessageTreeSyncChunkRequest[RI any] struct {
	RequestHandle types.RequestHandle // actual handle for outbound message, sentinel for inbound
	RequestInfo   *types.RequestInfo
	ToSeqNr       uint64
	StartIndex    jmt.Digest
	EndInclIndex  jmt.Digest
}

var _ MessageToStateSync[struct{}] = MessageTreeSyncChunkRequest[struct{}]{}

func (msg MessageTreeSyncChunkRequest[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, _ ocr3_1config.PublicConfig) bool {
	return true
}

func (msg MessageTreeSyncChunkRequest[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStateSync <- MessageToStateSyncWithSender[RI]{msg, sender}
}

func (msg MessageTreeSyncChunkRequest[RI]) processStateSync(stasy *stateSyncState[RI], sender commontypes.OracleID) {
	stasy.messageTreeSyncChunkRequest(msg, sender)
}

type MessageTreeSyncChunkResponse[RI any] struct {
	RequestHandle       types.RequestHandle // actual handle for outbound message, sentinel for inbound
	ToSeqNr             uint64
	StartIndex          jmt.Digest
	RequestEndInclIndex jmt.Digest
	GoAway              bool
	EndInclIndex        jmt.Digest
	KeyValues           []KeyValuePair
	BoundingLeaves      []jmt.BoundingLeaf
}

var _ MessageToStateSync[struct{}] = MessageTreeSyncChunkResponse[struct{}]{}

func (msg MessageTreeSyncChunkResponse[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int, config ocr3_1config.PublicConfig) bool {
	if len(msg.BoundingLeaves) > jmt.MaxBoundingLeaves {
		return false
	}
	for _, bl := range msg.BoundingLeaves {
		if len(bl.Siblings) > jmt.MaxProofLength {
			return false
		}
	}
	if len(msg.KeyValues) > config.GetMaxTreeSyncChunkKeys() {
		return false
	}
	treeSyncChunkLeavesSize := 0
	for _, kv := range msg.KeyValues {
		if len(kv.Key) > ocr3_1types.MaxMaxKeyValueKeyBytes {
			return false
		}
		if len(kv.Value) > ocr3_1types.MaxMaxKeyValueValueBytes {
			return false
		}
		treeSyncChunkLeavesSize += len(kv.Key) + len(kv.Value)
	}
	if treeSyncChunkLeavesSize > config.GetMaxTreeSyncChunkKeysPlusValuesBytes() {
		return false
	}
	return true
}

func (msg MessageTreeSyncChunkResponse[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStateSync <- MessageToStateSyncWithSender[RI]{msg, sender}
}

func (msg MessageTreeSyncChunkResponse[RI]) processStateSync(stasy *stateSyncState[RI], sender commontypes.OracleID) {
	stasy.messageTreeSyncChunkResponse(msg, sender)
}

type MessageBlobOffer[RI any] struct {
	RequestHandle    types.RequestHandle // actual handle for outbound message, sentinel for inbound
	RequestInfo      *types.RequestInfo
	ChunkDigestsRoot mt.Digest
	PayloadLength    uint64
	ExpirySeqNr      uint64
}

var _ MessageToBlobExchange[struct{}] = MessageBlobOffer[struct{}]{}

func (msg MessageBlobOffer[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, _ int, config ocr3_1config.PublicConfig) bool {
	return true
}

func (msg MessageBlobOffer[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToBlobExchange <- MessageToBlobExchangeWithSender[RI]{msg, sender}
}

func (msg MessageBlobOffer[RI]) processBlobExchange(bex *blobExchangeState[RI], sender commontypes.OracleID) {
	bex.messageBlobOffer(msg, sender)
}

type MessageBlobOfferResponse[RI any] struct {
	RequestHandle types.RequestHandle // actual handle for outbound message, sentinel for inbound
	BlobDigest    BlobDigest
	RejectOffer   bool
	Signature     BlobAvailabilitySignature
}

var _ MessageToBlobExchange[struct{}] = MessageBlobOfferResponse[struct{}]{}

func (msg MessageBlobOfferResponse[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int, _ ocr3_1config.PublicConfig) bool {
	if msg.RejectOffer {
		return len(msg.Signature) == 0
	} else {
		return len(msg.Signature) == ed25519.SignatureSize
	}
}

func (msg MessageBlobOfferResponse[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToBlobExchange <- MessageToBlobExchangeWithSender[RI]{msg, sender}
}

func (msg MessageBlobOfferResponse[RI]) processBlobExchange(bex *blobExchangeState[RI], sender commontypes.OracleID) {
	bex.messageBlobOfferResponse(msg, sender)
}

type MessageBlobChunkRequest[RI any] struct {
	RequestHandle types.RequestHandle // actual handle for outbound message, sentinel for inbound
	RequestInfo   *types.RequestInfo
	BlobDigest    BlobDigest
	ChunkIndex    uint64
}

var _ MessageToBlobExchange[struct{}] = MessageBlobChunkRequest[struct{}]{}

func (msg MessageBlobChunkRequest[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int, _ ocr3_1config.PublicConfig) bool {
	return true
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
	GoAway     bool
	Chunk      []byte
	Proof      []mt.Digest
}

var _ MessageToBlobExchange[struct{}] = MessageBlobChunkResponse[struct{}]{}

func (msg MessageBlobChunkResponse[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int, config ocr3_1config.PublicConfig) bool {
	if len(msg.Chunk) > config.GetBlobChunkBytes() {
		return false
	}
	return true
}

func (msg MessageBlobChunkResponse[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToBlobExchange <- MessageToBlobExchangeWithSender[RI]{msg, sender}
}

func (msg MessageBlobChunkResponse[RI]) processBlobExchange(bex *blobExchangeState[RI], sender commontypes.OracleID) {
	bex.messageBlobChunkResponse(msg, sender)
}
