package protocol

import (
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type EventToPacemaker[RI any] interface {
	processPacemaker(pace *pacemakerState[RI])
}

type EventProgress[RI any] struct{}

var _ EventToPacemaker[struct{}] = (*EventProgress[struct{}])(nil) // implements EventToPacemaker

func (ev EventProgress[RI]) processPacemaker(pace *pacemakerState[RI]) {
	pace.eventProgress()
}

type EventNewEpochRequest[RI any] struct{}

var _ EventToPacemaker[struct{}] = (*EventNewEpochRequest[struct{}])(nil) // implements EventToPacemaker

func (ev EventNewEpochRequest[RI]) processPacemaker(pace *pacemakerState[RI]) {
	pace.eventNewEpochRequest()
}

type EventToOutcomeGeneration[RI any] interface {
	processOutcomeGeneration(outgen *outcomeGenerationState[RI])
}

type EventNewEpochStart[RI any] struct {
	Epoch uint64
}

var _ EventToOutcomeGeneration[struct{}] = EventNewEpochStart[struct{}]{}

func (ev EventNewEpochStart[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventNewEpochStart(ev)
}

type EventComputedQuery[RI any] struct {
	Epoch uint64
	SeqNr uint64
	Query types.Query
}

var _ EventToOutcomeGeneration[struct{}] = EventComputedQuery[struct{}]{}

func (ev EventComputedQuery[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventComputedQuery(ev)
}

type EventComputedValidateVerifyObservation[RI any] struct {
	Epoch  uint64
	SeqNr  uint64
	Sender commontypes.OracleID
}

var _ EventToOutcomeGeneration[struct{}] = EventComputedValidateVerifyObservation[struct{}]{}

func (ev EventComputedValidateVerifyObservation[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventComputedValidateVerifyObservation(ev)
}

type EventComputedObservationQuorumSuccess[RI any] struct {
	Epoch uint64
	SeqNr uint64
}

var _ EventToOutcomeGeneration[struct{}] = EventComputedObservationQuorumSuccess[struct{}]{}

func (ev EventComputedObservationQuorumSuccess[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventComputedObservationQuorumSuccess(ev)
}

type EventComputedObservation[RI any] struct {
	Epoch       uint64
	SeqNr       uint64
	Query       types.Query
	Observation types.Observation
}

var _ EventToOutcomeGeneration[struct{}] = EventComputedObservation[struct{}]{}

func (ev EventComputedObservation[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventComputedObservation(ev)
}

type EventKVTransactionRequest[RI any] struct {
	RoundCtx                    ocr3_1types.RoundContext
	Query                       types.Query
	Asos                        []AttributedSignedObservation
	Prepared                    bool
	StateTransitionOutputDigest StateTransitionOutputDigest
	ReportsPlusPrecursorDigest  ReportsPlusPrecursorDigest
	CommitQC                    []AttributedCommitSignature
}

var _ EventToStatePersistence[struct{}] = EventKVTransactionRequest[struct{}]{} // implements EventToStatePersistence

func (ev EventKVTransactionRequest[RI]) processStatePersistence(state *statePersistenceState[RI]) {
	state.eventKVTransactionRequest(ev)
}

type EventProduceStateTransition[RI any] struct {
	RoundCtx                    ocr3_1types.RoundContext
	Txn                         kvStoreTxn
	Query                       types.Query
	Asos                        []AttributedSignedObservation
	Prepared                    bool
	StateTransitionOutputDigest StateTransitionOutputDigest
	ReportsPlusPrecursorDigest  ReportsPlusPrecursorDigest
	CommitQC                    []AttributedCommitSignature
}

var _ EventToOutcomeGeneration[struct{}] = EventProduceStateTransition[struct{}]{}

func (ev EventProduceStateTransition[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventProduceStateTransition(ev)
}

type EventComputedStateTransition[RI any] struct {
	SeqNr uint64
	Epoch uint64
}

var _ EventToStatePersistence[struct{}] = EventComputedStateTransition[struct{}]{} // implements EventToStatePersistence

func (ev EventComputedStateTransition[RI]) processStatePersistence(state *statePersistenceState[RI]) {
	state.eventComputedStateTransition(ev)
}

type EventToReportAttestation[RI any] interface {
	processReportAttestation(repatt *reportAttestationState[RI])
}

type EventToStatePersistence[RI any] interface {
	processStatePersistence(state *statePersistenceState[RI])
}

type EventToTransmission[RI any] interface {
	processTransmission(t *transmissionState[RI])
}
type EventMissingOutcome[RI any] struct {
	SeqNr uint64
}

var _ EventToReportAttestation[struct{}] = EventMissingOutcome[struct{}]{} // implements EventToReportAttestation

func (ev EventMissingOutcome[RI]) processReportAttestation(repatt *reportAttestationState[RI]) {
	repatt.eventMissingOutcome(ev)
}

type EventCertifiedCommit[RI any] struct {
	CertifiedCommittedReports CertifiedCommittedReports[RI]
}

var _ EventToReportAttestation[struct{}] = EventCertifiedCommit[struct{}]{} // implements EventToReportAttestation

func (ev EventCertifiedCommit[RI]) processReportAttestation(repatt *reportAttestationState[RI]) {
	repatt.eventCertifiedCommit(ev)
}

type EventComputedReports[RI any] struct {
	SeqNr       uint64
	ReportsPlus []ocr3types.ReportPlus[RI]
}

var _ EventToReportAttestation[struct{}] = EventComputedReports[struct{}]{} // implements EventToReportAttestation

func (ev EventComputedReports[RI]) processReportAttestation(repatt *reportAttestationState[RI]) {
	repatt.eventComputedReports(ev)
}

type EventAttestedReport[RI any] struct {
	SeqNr                        uint64
	Index                        int
	AttestedReport               AttestedReportMany[RI]
	TransmissionScheduleOverride *ocr3types.TransmissionSchedule
}

var _ EventToTransmission[struct{}] = EventAttestedReport[struct{}]{} // implements EventToTransmission

func (ev EventAttestedReport[RI]) processTransmission(t *transmissionState[RI]) {
	t.eventAttestedReport(ev)
}

type EventAttestedStateTransitionBlock[RI any] struct {
	AttestedStateTransitionBlock AttestedStateTransitionBlock
}

var _ EventToStatePersistence[struct{}] = EventAttestedStateTransitionBlock[struct{}]{} // implements EventToStatePersistence

func (ev EventAttestedStateTransitionBlock[RI]) processStatePersistence(state *statePersistenceState[RI]) {
	state.eventAttestedStateTransitionBlock(ev)
}

type EventAcknowledgedCommittedKVTransaction[RI any] struct {
	SeqNr uint64
}

var _ EventToStatePersistence[struct{}] = EventAcknowledgedCommittedKVTransaction[struct{}]{} // implements EventToStatePersistence

func (ev EventAcknowledgedCommittedKVTransaction[RI]) processStatePersistence(state *statePersistenceState[RI]) {
	state.eventAcknowledgedCommittedKVTransaction(ev)
}

type EventAcknowledgedComputedStateTransition[RI any] struct {
	SeqNr uint64
}

var _ EventToOutcomeGeneration[struct{}] = EventAcknowledgedComputedStateTransition[struct{}]{} // implements EventToOutcomeGeneration

func (ev EventAcknowledgedComputedStateTransition[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventAcknowledgedComputedStateTransition(ev)
}

type EventCommittedKVTransaction[RI any] struct {
	SeqNr uint64
}

var _ EventToOutcomeGeneration[struct{}] = EventCommittedKVTransaction[struct{}]{} // implements EventToOutcomeGeneration

func (ev EventCommittedKVTransaction[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventCommittedKVTransaction(ev)
}

type EventStateSyncRequest[RI any] struct {
	SeqNr uint64
}

var _ EventToStatePersistence[struct{}] = EventStateSyncRequest[struct{}]{} // implements EventToStatePersistence

func (ev EventStateSyncRequest[RI]) processStatePersistence(state *statePersistenceState[RI]) {
	state.eventStateSyncRequest(ev)
}

type EventBlockSyncSummaryHeartbeat[RI any] struct{}

var _ EventToStatePersistence[struct{}] = EventBlockSyncSummaryHeartbeat[struct{}]{} // implements EventToStatePersistence

func (ev EventBlockSyncSummaryHeartbeat[RI]) processStatePersistence(state *statePersistenceState[RI]) {
	state.eventEventBlockSyncSummaryHeartbeat(ev)
}

type EventExpiredBlockSyncRequest[RI any] struct {
	RequestedFrom commontypes.OracleID
	Nonce         uint64
}

var _ EventToStatePersistence[struct{}] = EventExpiredBlockSyncRequest[struct{}]{} // implements EventToStatePersistence

func (ev EventExpiredBlockSyncRequest[RI]) processStatePersistence(state *statePersistenceState[RI]) {
	state.eventExpiredBlockSyncRequest(ev)
}

type EventReadyToSendNextBlockSyncRequest[RI any] struct{}

var _ EventToStatePersistence[struct{}] = EventReadyToSendNextBlockSyncRequest[struct{}]{} // implements EventToStatePersistence

func (ev EventReadyToSendNextBlockSyncRequest[RI]) processStatePersistence(state *statePersistenceState[RI]) {
	state.eventReadyToSendNextBlockSyncRequest(ev)
}

type EventReplayVerifiedStateTransition[RI any] struct {
	AttestedStateTransitionBlock AttestedStateTransitionBlock
}

var _ EventToOutcomeGeneration[struct{}] = EventReplayVerifiedStateTransition[struct{}]{} // implements EventToOutcomeGeneration

func (ev EventReplayVerifiedStateTransition[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventReplayVerifiedStateTransition(ev)
}

type EventDiscardKVTransaction[RI any] struct {
	SeqNr uint64
}

var _ EventToStatePersistence[struct{}] = EventDiscardKVTransaction[struct{}]{} // implements EventToStatePersistence

func (ev EventDiscardKVTransaction[RI]) processStatePersistence(state *statePersistenceState[RI]) {
	state.eventDiscardKVTransaction(ev)
}
