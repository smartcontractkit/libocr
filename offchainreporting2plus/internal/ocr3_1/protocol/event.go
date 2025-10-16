package protocol

import (
	"github.com/RoSpaceDev/libocr/commontypes"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/ocr3types"
	"github.com/RoSpaceDev/libocr/offchainreporting2plus/types"
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
	Epoch           uint64
	SeqNr           uint64
	AttributedQuery types.AttributedQuery
	Observation     types.Observation
}

var _ EventToOutcomeGeneration[struct{}] = EventComputedObservation[struct{}]{}

func (ev EventComputedObservation[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventComputedObservation(ev)
}

type EventComputedProposalStateTransition[RI any] struct {
	Epoch                                uint64
	SeqNr                                uint64
	KeyValueDatabaseReadWriteTransaction KeyValueDatabaseReadWriteTransaction
	stateTransitionInfo                  stateTransitionInfo
}

var _ EventToOutcomeGeneration[struct{}] = EventComputedProposalStateTransition[struct{}]{}

func (ev EventComputedProposalStateTransition[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventComputedProposalStateTransition(ev)
}

type EventComputedCommitted[RI any] struct {
	Epoch uint64
	SeqNr uint64
}

var _ EventToOutcomeGeneration[struct{}] = EventComputedCommitted[struct{}]{}

func (ev EventComputedCommitted[RI]) processOutcomeGeneration(outgen *outcomeGenerationState[RI]) {
	outgen.eventComputedCommitted(ev)
}

type EventToReportAttestation[RI any] interface {
	processReportAttestation(repatt *reportAttestationState[RI])
}

type EventToStateSync[RI any] interface {
	processStateSync(stasy *stateSyncState[RI])
}

type EventToBlobExchange[RI any] interface {
	processBlobExchange(bex *blobExchangeState[RI])
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

type EventStateSyncRequest[RI any] struct {
	SeqNr uint64
}

var _ EventToStateSync[struct{}] = EventStateSyncRequest[struct{}]{} // implements EventToStateSync

func (ev EventStateSyncRequest[RI]) processStateSync(stasy *stateSyncState[RI]) {
	stasy.eventStateSyncRequest(ev)
}

type EventBlobBroadcastRequestRespond[RI any] struct {
	BlobDigest BlobDigest
	Request    blobBroadcastRequest
}

var _ EventToBlobExchange[struct{}] = EventBlobBroadcastRequestRespond[struct{}]{} // implements EventToBlobExchange

func (ev EventBlobBroadcastRequestRespond[RI]) processBlobExchange(bex *blobExchangeState[RI]) {
	bex.eventBlobBroadcastRequestRespond(ev)
}

type EventBlobBroadcastRequestDone[RI any] struct {
	BlobDigest BlobDigest
}

var _ EventToBlobExchange[struct{}] = EventBlobBroadcastRequestDone[struct{}]{} // implements EventToBlobExchange

func (ev EventBlobBroadcastRequestDone[RI]) processBlobExchange(bex *blobExchangeState[RI]) {
	bex.eventBlobBroadcastRequestDone(ev)
}

type EventBlobFetchRequestRespond[RI any] struct {
	BlobDigest BlobDigest
	Request    blobFetchRequest
}

var _ EventToBlobExchange[struct{}] = EventBlobFetchRequestRespond[struct{}]{} // implements EventToBlobExchange

func (ev EventBlobFetchRequestRespond[RI]) processBlobExchange(bex *blobExchangeState[RI]) {
	bex.eventBlobFetchRequestRespond(ev)
}

type EventBlobFetchRequestDone[RI any] struct {
	BlobDigest BlobDigest
}

var _ EventToBlobExchange[struct{}] = EventBlobFetchRequestDone[struct{}]{} // implements EventToBlobExchange

func (ev EventBlobFetchRequestDone[RI]) processBlobExchange(bex *blobExchangeState[RI]) {
	bex.eventBlobFetchRequestDone(ev)
}

type EventBlobBroadcastGraceTimeout[RI any] struct {
	BlobDigest BlobDigest
}

var _ EventToBlobExchange[struct{}] = EventBlobBroadcastGraceTimeout[struct{}]{} // implements EventToBlobExchange

func (ev EventBlobBroadcastGraceTimeout[RI]) processBlobExchange(bex *blobExchangeState[RI]) {
	bex.eventBlobBroadcastGraceTimeout(ev)
}
