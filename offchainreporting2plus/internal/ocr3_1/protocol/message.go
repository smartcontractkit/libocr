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

	CheckPriority(priority types.BinaryMessageOutboundPriority) bool

	CheckMessageType(inboundMessageType MessageType) bool // check in here with some switch if it is the expected type

	GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage

	// process passes this Message instance to the oracle o, as a message from
	// oracle with the given sender index
	process(o *oracleState[RI], sender commontypes.OracleID)
}

type SerializableRequestMessage[RI any] interface {
	Message[RI]
	NewInboundRequestMessage(handle types.RequestHandle) Message[RI]
}

type SerializableResponseMessage[RI any] interface {
	Message[RI]
	NewInboundResponseMessage() Message[RI]
}

type RequestMessage[RI any] interface {
	Message[RI]
	GetSerializableRequestMessage() Message[RI]
	GetRequestHandle() types.RequestHandle
}

type ResponseMessage[RI any] interface {
	Message[RI]
	GetSerializableResponseMessage() Message[RI]
}

type MessageType int

const (
	MessageTypePlain    = 1
	MessageTypeRequest  = 2
	MessageTypeResponse = 3
)

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

type MessageNewEpochWish[RI any] struct {
	Epoch uint64
}

var _ MessageToPacemaker[struct{}] = (*MessageNewEpochWish[struct{}])(nil)

func (msg MessageNewEpochWish[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (msg MessageNewEpochWish[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageNewEpochWish[RI]) CheckMessageType(messageType MessageType) bool {
	return messageType == MessageTypePlain
}

func (msg MessageNewEpochWish[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessageEpochStartRequest[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageEpochStartRequest[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageEpochStartRequest[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessageEpochStart[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageEpochStart[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageEpochStart[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessageRoundStart[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageRoundStart[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageRoundStart[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessageObservation[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageObservation[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageObservation[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessageProposal[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageProposal[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageProposal[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessagePrepare[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessagePrepare[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessagePrepare[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessageCommit[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageCommit[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageCommit[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessageReportSignatures[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageReportSignatures[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageReportSignatures[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessageCertifiedCommitRequest[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageCertifiedCommitRequest[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageCertifiedCommitRequest[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
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

func (msg MessageCertifiedCommit[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageCertifiedCommit[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageCertifiedCommit[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
}

func (msg MessageCertifiedCommit[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportAttestation <- MessageToReportAttestationWithSender[RI]{msg, sender}
}

func (msg MessageCertifiedCommit[RI]) processReportAttestation(repatt *reportAttestationState[RI], sender commontypes.OracleID) {
	repatt.messageCertifiedCommit(msg, sender)
}

type MessageBlockSyncRequestWrapper[RI any] struct {
	SerializableMessage Message[RI]
	RequestHandle       types.RequestHandle
	ResponsePolicy      types.ResponsePolicy
}

type MessageBlockSyncRequest[RI any] struct {
	HighestCommittedSeqNr uint64
	Nonce                 uint64
}

var _ MessageToStatePersistence[struct{}] = (*MessageBlockSyncRequestWrapper[struct{}])(nil)
var _ Message[struct{}] = (*MessageBlockSyncRequestWrapper[struct{}])(nil)
var _ RequestMessage[struct{}] = (*MessageBlockSyncRequestWrapper[struct{}])(nil)

var _ SerializableRequestMessage[struct{}] = (*MessageBlockSyncRequest[struct{}])(nil)

func (msg MessageBlockSyncRequest[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (msg MessageBlockSyncRequest[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityLow
}

func (msg MessageBlockSyncRequest[RI]) CheckMessageType(messageType MessageType) bool {
	return messageType == MessageTypeRequest
}

func (msg MessageBlockSyncRequest[RI]) GetOutboundBinaryMessage([]byte) types.OutboundBinaryMessage {
	panic("should have called the process method on the MessageBlockSyncRequestWrapper instead")
}

func (msg MessageBlockSyncRequest[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	panic("should have called the process method on the MessageBlockSyncRequestWrapper instead")
}

func (msg MessageBlockSyncRequest[RI]) NewInboundRequestMessage(handle types.RequestHandle) Message[RI] {
	return MessageBlockSyncRequestWrapper[RI]{
		SerializableMessage: msg,
		RequestHandle:       handle,
	}
}

func (msg MessageBlockSyncRequestWrapper[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxReportSigLen int) bool {
	return msg.SerializableMessage.CheckSize(n, f, limits, maxReportSigLen)
}

func (msg MessageBlockSyncRequestWrapper[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return msg.SerializableMessage.CheckPriority(priority)
}

func (msg MessageBlockSyncRequestWrapper[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return msg.SerializableMessage.CheckMessageType(inboundMessageType)
}

func (msg MessageBlockSyncRequestWrapper[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessageRequest{
		msg.ResponsePolicy,
		sMsg,
		types.BinaryMessagePriorityLow,
	}
}

func (msg MessageBlockSyncRequestWrapper[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStatePersistence <- MessageToStatePersistenceWithSender[RI]{msg, sender}
}

func (msg MessageBlockSyncRequestWrapper[RI]) processStatePersistence(state *statePersistenceState[RI], sender commontypes.OracleID) {
	state.messageBlockSyncReq(msg, sender)
}

func (msg MessageBlockSyncRequestWrapper[RI]) GetRequestHandle() types.RequestHandle {
	return msg.RequestHandle
}

func (msg MessageBlockSyncRequestWrapper[RI]) GetSerializableRequestMessage() Message[RI] {
	return msg.SerializableMessage
}

type MessageBlockSyncSummary[RI any] struct {
	LowestPersistedSeqNr uint64
}

var _ MessageToStatePersistence[struct{}] = (*MessageBlockSyncSummary[struct{}])(nil)

func (msg MessageBlockSyncSummary[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (msg MessageBlockSyncSummary[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityDefault
}

func (msg MessageBlockSyncSummary[RI]) CheckMessageType(inboundMessageType MessageType) bool {
	return inboundMessageType == MessageTypePlain
}

func (msg MessageBlockSyncSummary[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.OutboundBinaryMessagePlain{
		sMsg,
		types.BinaryMessagePriorityDefault,
	}
}

func (msg MessageBlockSyncSummary[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStatePersistence <- MessageToStatePersistenceWithSender[RI]{msg, sender}
}

func (msg MessageBlockSyncSummary[RI]) processStatePersistence(state *statePersistenceState[RI], sender commontypes.OracleID) {
	state.messageBlockSyncSummary(msg, sender)
}

type MessageBlockSyncWrapper[RI any] struct {
	SerializableMessage Message[RI]
	RequestHandle       types.RequestHandle
}

type MessageBlockSync[RI any] struct {
	AttestedStateTransitionBlocks []AttestedStateTransitionBlock
	Nonce                         uint64
}

var _ MessageToStatePersistence[struct{}] = (*MessageBlockSyncWrapper[struct{}])(nil)
var _ Message[struct{}] = (*MessageBlockSyncWrapper[struct{}])(nil)
var _ ResponseMessage[struct{}] = (*MessageBlockSyncWrapper[struct{}])(nil)

var _ SerializableResponseMessage[struct{}] = (*MessageBlockSync[struct{}])(nil)

func (msg MessageBlockSync[RI]) CheckSize(n int, f int, _ ocr3_1types.ReportingPluginLimits, _ int) bool {
	return true
}

func (msg MessageBlockSync[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return priority == types.BinaryMessagePriorityLow
}

func (msg MessageBlockSync[RI]) CheckMessageType(messageType MessageType) bool {
	return messageType == MessageTypeResponse
}

func (msg MessageBlockSync[RI]) GetOutboundBinaryMessage([]byte) types.OutboundBinaryMessage {
	panic("should have called the process method on MessageBlockSyncWrapper instead")
}

func (msg MessageBlockSync[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	panic("should have called the process method on MessageBlockSyncWrapper instead")
}

func (msg MessageBlockSync[RI]) NewInboundResponseMessage() Message[RI] {
	return MessageBlockSyncWrapper[RI]{
		SerializableMessage: msg,
	}
}

func (msg MessageBlockSyncWrapper[RI]) CheckSize(n int, f int, limits ocr3_1types.ReportingPluginLimits, maxSigLen int) bool {
	return msg.SerializableMessage.CheckSize(n, f, limits, maxSigLen)
}

func (msg MessageBlockSyncWrapper[RI]) CheckPriority(priority types.BinaryMessageOutboundPriority) bool {
	return msg.SerializableMessage.CheckPriority(priority)
}

func (msg MessageBlockSyncWrapper[RI]) CheckMessageType(messageType MessageType) bool {
	return msg.SerializableMessage.CheckMessageType(messageType)
}

func (msg MessageBlockSyncWrapper[RI]) GetOutboundBinaryMessage(sMsg []byte) types.OutboundBinaryMessage {
	return types.MustMakeOutboundBinaryMessageResponse(
		msg.RequestHandle,
		sMsg,
		types.BinaryMessagePriorityLow,
	)
}

func (msg MessageBlockSyncWrapper[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToStatePersistence <- MessageToStatePersistenceWithSender[RI]{msg, sender}
}

func (msg MessageBlockSyncWrapper[RI]) processStatePersistence(state *statePersistenceState[RI], sender commontypes.OracleID) {
	state.messageBlockSync(msg, sender)
}

func (msg MessageBlockSyncWrapper[RI]) GetSerializableResponseMessage() Message[RI] {
	return msg.SerializableMessage
}
