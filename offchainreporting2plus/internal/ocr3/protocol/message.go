package protocol //

import (
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// EventToPacemaker is the interface used to pass in-process events to the
// leader-election protocol.
type EventToPacemaker[RI any] interface {
	// processPacemaker is called when the local oracle process invokes an event
	// intended for the leader-election protocol.
	processPacemaker(pace *pacemakerState[RI])
}

// EventProgress is used to process the "progress" event passed by the local
// oracle from its the reporting protocol to the leader-election protocol. It is
// sent by the reporting protocol when the leader has produced a valid new
// report.
type EventProgress[RI any] struct{}

var _ EventToPacemaker[struct{}] = (*EventProgress[struct{}])(nil) // implements EventToPacemaker

func (ev EventProgress[RI]) processPacemaker(pace *pacemakerState[RI]) {
	pace.eventProgress()
}

// EventChangeLeader is used to process the "change-leader" event passed by the
// local oracle from its the reporting protocol to the leader-election protocol
type EventChangeLeader[RI any] struct{}

var _ EventToPacemaker[struct{}] = (*EventChangeLeader[struct{}])(nil) // implements EventToPacemaker

func (ev EventChangeLeader[RI]) processPacemaker(pace *pacemakerState[RI]) {
	pace.eventChangeLeader()
}

type EventToReportGeneration[RI any] interface {
	processReportGeneration(repgen *reportGenerationState[RI])
}

type EventStartNewEpoch[RI any] struct {
	Epoch uint64
}

var _ EventToReportGeneration[struct{}] = EventStartNewEpoch[struct{}]{}

func (ev EventStartNewEpoch[RI]) processReportGeneration(repgen *reportGenerationState[RI]) {
	repgen.eventStartNewEpoch(ev)
}

type EventToReportFinalization[RI any] interface {
	processReportFinalization(repfin *reportFinalizationState[RI])
}

// EventToTransmission is the interface used to pass a completed report to the
// protocol which will transmit it to the on-chain smart contract.
type EventToTransmission[RI any] interface {
	processTransmission(t *transmissionState[RI])
}

// Message is the interface used to pass an inter-oracle message to the local
// oracle process.
type Message[RI any] interface {
	// CheckSize checks whether the given message conforms to the limits imposed by
	// reportingPluginLimits
	CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool

	// process passes this Message instance to the oracle o, as a message from
	// oracle with the given sender index
	process(o *oracleState[RI], sender commontypes.OracleID)
}

// MessageWithSender records a msg with the index of the sender oracle
type MessageWithSender[RI any] struct {
	Msg    Message[RI]
	Sender commontypes.OracleID
}

// MessageToPacemaker is the interface used to pass a message to the local
// leader-election protocol
type MessageToPacemaker[RI any] interface {
	Message[RI]

	// process passes this MessageToPacemaker instance to the oracle o, as a
	// message from oracle with the given sender index
	processPacemaker(pace *pacemakerState[RI], sender commontypes.OracleID)
}

// MessageToPacemakerWithSender records a msg with the idx of the sender oracle
type MessageToPacemakerWithSender[RI any] struct {
	msg    MessageToPacemaker[RI]
	sender commontypes.OracleID
}

// MessageToReportGeneration is the interface used to pass an inter-oracle message
// to the local oracle reporting process.
type MessageToReportGeneration[RI any] interface {
	Message[RI]

	// processReportGeneration is called to send this message to the local oracle
	// reporting process.
	processReportGeneration(repgen *reportGenerationState[RI], sender commontypes.OracleID)

	epoch() uint64
}

// MessageToReportGenerationWithSender records a message destined for the oracle
// reporting
type MessageToReportGenerationWithSender[RI any] struct {
	msg    MessageToReportGeneration[RI]
	sender commontypes.OracleID
}

type MessageToReportFinalization[RI any] interface {
	Message[RI]

	processReportFinalization(repfin *reportFinalizationState[RI], sender commontypes.OracleID)
}

type MessageToReportFinalizationWithSender[RI any] struct {
	msg    MessageToReportFinalization[RI]
	sender commontypes.OracleID
}

// MessageNewEpoch corresponds to the "newepoch(epoch_number)" message from alg.
// 1. It indicates that the node believes the protocol should move to the
// specified epoch.
type MessageNewEpoch[RI any] struct {
	Epoch uint64
}

var _ MessageToPacemaker[struct{}] = (*MessageNewEpoch[struct{}])(nil)

func (msg MessageNewEpoch[RI]) CheckSize(types.ReportingPluginLimits) bool {
	return true
}

func (msg MessageNewEpoch[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToPacemaker <- MessageToPacemakerWithSender[RI]{msg, sender}
}

func (msg MessageNewEpoch[RI]) processPacemaker(pace *pacemakerState[RI], sender commontypes.OracleID) {
	pace.messageNewepoch(msg, sender)
}

type MessageReconcile[RI any] struct {
	Epoch                           uint64
	HighestCertified                CertifiedPrepareOrCommit
	SignedHighestCertifiedTimestamp SignedHighestCertifiedTimestamp
}

var _ MessageToReportGeneration[struct{}] = (*MessageReconcile[struct{}])(nil)

func (msg MessageReconcile[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true
}

func (msg MessageReconcile[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageReconcile[RI]) processReportGeneration(repgen *reportGenerationState[RI], sender commontypes.OracleID) {
	repgen.messageReconcile(msg, sender)
}

func (msg MessageReconcile[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessageStartEpoch[RI any] struct {
	Epoch           uint64
	StartEpochProof StartEpochProof
}

var _ MessageToReportGeneration[struct{}] = (*MessageStartEpoch[struct{}])(nil)

func (msg MessageStartEpoch[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true
}

func (msg MessageStartEpoch[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageStartEpoch[RI]) processReportGeneration(repgen *reportGenerationState[RI], sender commontypes.OracleID) {
	repgen.messageStartEpoch(msg, sender)
}

func (msg MessageStartEpoch[RI]) epoch() uint64 {
	return msg.Epoch
}

// MessageStartRound corresponds to the "observe-req" message from alg. 2. The
// leader transmits this to request observations from participating oracles, so
// that it can collate them into a report.
type MessageStartRound[RI any] struct {
	Epoch uint64
	SeqNr uint64
	Query types.Query
}

var _ MessageToReportGeneration[struct{}] = (*MessageStartRound[struct{}])(nil)

func (msg MessageStartRound[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true
	// return len(msg.Query) <= reportingPluginLimits.MaxQueryLength
}

func (msg MessageStartRound[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageStartRound[RI]) processReportGeneration(repgen *reportGenerationState[RI], sender commontypes.OracleID) {
	repgen.messageStartRound(msg, sender)
}

func (msg MessageStartRound[RI]) epoch() uint64 {
	return msg.Epoch
}

// MessageObserve corresponds to the "observe" message from alg. 2.
// Participating oracles send this back to the leader in response to
// MessageStartRound's.
type MessageObserve[RI any] struct {
	Epoch             uint64
	SeqNr             uint64
	SignedObservation SignedObservation
}

var _ MessageToReportGeneration[struct{}] = (*MessageObserve[struct{}])(nil)

func (msg MessageObserve[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true
	// return len(msg.SignedObservation.Observation) <= reportingPluginLimits.MaxObservationLength
}

func (msg MessageObserve[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageObserve[RI]) processReportGeneration(repgen *reportGenerationState[RI], sender commontypes.OracleID) {
	repgen.messageObserve(msg, sender)
}

func (msg MessageObserve[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessagePropose[RI any] struct {
	Epoch                        uint64
	SeqNr                        uint64
	AttributedSignedObservations []AttributedSignedObservation
}

var _ MessageToReportGeneration[struct{}] = MessagePropose[struct{}]{}

func (msg MessagePropose[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true
}

func (msg MessagePropose[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessagePropose[RI]) processReportGeneration(repgen *reportGenerationState[RI], sender commontypes.OracleID) {
	repgen.messagePropose(msg, sender)
}

func (msg MessagePropose[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessagePrepare[RI any] struct {
	Epoch     uint64
	SeqNr     uint64
	Signature PrepareSignature
}

var _ MessageToReportGeneration[struct{}] = MessagePrepare[struct{}]{}

func (msg MessagePrepare[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true
}

func (msg MessagePrepare[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessagePrepare[RI]) processReportGeneration(repgen *reportGenerationState[RI], sender commontypes.OracleID) {
	repgen.messagePrepare(msg, sender)
}

func (msg MessagePrepare[RI]) epoch() uint64 {
	return msg.Epoch
}

type MessageCommit[RI any] struct {
	Epoch     uint64
	SeqNr     uint64
	Signature CommitSignature
}

var _ MessageToReportGeneration[struct{}] = MessageCommit[struct{}]{}

func (msg MessageCommit[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true
}

func (msg MessageCommit[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportGeneration <- MessageToReportGenerationWithSender[RI]{
		msg,
		sender,
	}
}

func (msg MessageCommit[RI]) processReportGeneration(repgen *reportGenerationState[RI], sender commontypes.OracleID) {
	repgen.messageCommit(msg, sender)
}

func (msg MessageCommit[RI]) epoch() uint64 {
	return msg.Epoch
}

// MessageFinal corresponds to the "final" message in alg. 2. It is sent by the
// current leader with the aggregated signature(s) to all participating oracles,
// for them to participate in the subsequent transmission of the report to the
// on-chain contract.
type MessageFinal[RI any] struct {
	SeqNr            uint64
	ReportSignatures [][]byte
}

var _ MessageToReportFinalization[struct{}] = MessageFinal[struct{}]{}

func (msg MessageFinal[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true

	// return len(msg.AttestedReport.AttributedSignatures) <= types.MaxOracles &&
	// 	len(msg.AttestedReport.Report) <= reportingPluginLimits.MaxReportLength
}

func (msg MessageFinal[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportFinalization <- MessageToReportFinalizationWithSender[RI]{msg, sender}
}

func (msg MessageFinal[RI]) processReportFinalization(repfin *reportFinalizationState[RI], sender commontypes.OracleID) {
	repfin.messageFinal(msg, sender)
}

type MessageRequestCertifiedCommit[RI any] struct {
	SeqNr uint64
}

var _ MessageToReportFinalization[struct{}] = MessageRequestCertifiedCommit[struct{}]{}

func (msg MessageRequestCertifiedCommit[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true

}

func (msg MessageRequestCertifiedCommit[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportFinalization <- MessageToReportFinalizationWithSender[RI]{msg, sender}
}

func (msg MessageRequestCertifiedCommit[RI]) processReportFinalization(repfin *reportFinalizationState[RI], sender commontypes.OracleID) {
	repfin.messageRequestCertifiedCommit(msg, sender)
}

type MessageSupplyCertifiedCommit[RI any] struct {
	CertifiedCommit CertifiedPrepareOrCommitCommit
}

var _ MessageToReportFinalization[struct{}] = MessageSupplyCertifiedCommit[struct{}]{}

func (msg MessageSupplyCertifiedCommit[RI]) CheckSize(reportingPluginLimits types.ReportingPluginLimits) bool {
	return true

}

func (msg MessageSupplyCertifiedCommit[RI]) process(o *oracleState[RI], sender commontypes.OracleID) {
	o.chNetToReportFinalization <- MessageToReportFinalizationWithSender[RI]{msg, sender}
}

func (msg MessageSupplyCertifiedCommit[RI]) processReportFinalization(repfin *reportFinalizationState[RI], sender commontypes.OracleID) {
	repfin.messageSupplyCertifiedCommit(msg, sender)
}

type EventMissingOutcome[RI any] struct {
	SeqNr uint64
}

var _ EventToReportFinalization[struct{}] = EventMissingOutcome[struct{}]{} // implements EventToReportFinalization

func (ev EventMissingOutcome[RI]) processReportFinalization(repfin *reportFinalizationState[RI]) {
	repfin.eventMissingOutcome(ev)
}

type EventDeliver[RI any] struct {
	CertifiedCommit CertifiedPrepareOrCommitCommit
}

var _ EventToReportFinalization[struct{}] = EventDeliver[struct{}]{} // implements EventToReportFinalization

func (ev EventDeliver[RI]) processReportFinalization(repfin *reportFinalizationState[RI]) {
	repfin.eventDeliver(ev)
}

// EventTransmit is used to process the "transmit" event passed by the local
// reporting protocol to to the local transmit-to-the-onchain-smart-contract
// protocol.
type EventTransmit[RI any] struct {
	SeqNr          uint64
	Index          int
	AttestedReport AttestedReportMany[RI]
}

var _ EventToTransmission[struct{}] = EventTransmit[struct{}]{} // implements EventToTransmission

func (ev EventTransmit[RI]) processTransmission(t *transmissionState[RI]) {
	t.eventTransmit(ev)
}
