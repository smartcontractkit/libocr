package tracing

import (
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Trace interface {
	isTrace()
	String() string
}

type TraceType int

const (
	// commontypes.BinaryNetworkEndpoint
	TypeSendTo TraceType = iota
	TypeBroadcast
	TypeReceive
	TypeDrop
	TypeEndpointStart
	TypeEndpointClose

	// types.Database
	TypeReadState
	TypeWriteState
	TypeReadConfig
	TypeWriteConfig
	TypeStorePendingTransmission
	TypePendingTransmissionsWithConfigDigest
	TypeDeletePendingTransmission
	TypeDeletePendingTransmissionsOlderThan

	// types.ContractConfigTracker
	TypeNotify
	TypeLatestConfigDetails
	TypeLatestConfig
	TypeLatestBlockHeight

	// Transmission contract operations
	TypeTransmit
	TypeLatestConfigDigestAndEpoch
	TypeFromAccount

	// types.ReportingPlugin
	TypeQuery
	TypeObservation
	TypeReport
	TypeShouldAcceptFinalizedReport
	TypeShouldTransmitAcceptedReport
	TypePluginStart
	TypePluginClose

	// types.OffchainConfigDigester
	TypeConfigDigest
	TypeConfigDigestPrefix
)

type Common struct {
	Typ        TraceType
	Originator OracleID
	Timestamp  time.Time
}

// commontypes.BinaryNetworkEndpoint

type SendTo struct {
	Common
	Src     OracleID
	Dst     OracleID
	Message protocol.Message
}

type Broadcast struct {
	Common
	Src     OracleID
	Message protocol.Message
}

type Receive struct {
	Common
	Src     OracleID
	Dst     OracleID
	Message protocol.Message
}

type Drop struct {
	Common
	Src     OracleID
	Dst     OracleID
	Message protocol.Message
}

type EndpointStart struct {
	Common
}

type EndpointClose struct {
	Common
}

// types.Database

type ReadState struct {
	Common
	Digest types.ConfigDigest
	State  types.PersistentState
	Err    error
}

type WriteState struct {
	Common
	Digest types.ConfigDigest
	State  types.PersistentState
	Err    error
}

type ReadConfig struct {
	Common
	Config types.ContractConfig
	Err    error
}

type WriteConfig struct {
	Common
	Config types.ContractConfig
	Err    error
}

type StorePendingTransmission struct {
	Common
	Timestamp    types.ReportTimestamp
	Transmission types.PendingTransmission
	Err          error
}

type PendingTransmissionsWithConfigDigest struct {
	Common
	Digest types.ConfigDigest
	Err    error
}

type DeletePendingTransmission struct {
	Common
	Timestamp types.ReportTimestamp
	Err       error
}

type DeletePendingTransmissionsOlderThan struct {
	Common
	Cutoff time.Time
	Err    error
}

// types.ContractConfigTracker

type Notify struct {
	Common
}

type LatestConfigDetails struct {
	Common
	ChangedInBlock uint64
	Digest         types.ConfigDigest
	Err            error
}

type LatestConfig struct {
	Common
	ChangedInBlock uint64
	Config         types.ContractConfig
	Err            error
}

type LatestBlockHeight struct {
	Common
	BlockHeight uint64
	Err         error
}

// types.ContractTransmitter

type Transmit struct {
	Common
	ReportContext types.ReportContext
	Report        types.Report
	Signatures    []types.AttributedOnchainSignature
	Err           error
}

type LatestConfigDigestAndEpoch struct {
	Common
	Digest types.ConfigDigest
	Epoch  uint32
	Err    error
}

type FromAccount struct {
	Common
	Account types.Account
}

// types.ReportingPlugin

type Query struct {
	Common
	Timestamp types.ReportTimestamp
	Query     types.Query
	Err       error
}

type Observation struct {
	Common
	Timestamp   types.ReportTimestamp
	Query       types.Query
	Observation types.Observation
	Err         error
}

type Report struct {
	Common
	Timestamp    types.ReportTimestamp
	Query        types.Query
	Observations []types.AttributedObservation
	OK           bool
	Report       types.Report
	Err          error
}

type ShouldAcceptFinalizedReport struct {
	Common
	Timestamp types.ReportTimestamp
	Report    types.Report
	OK        bool
	Err       error
}

type ShouldTransmitAcceptedReport struct {
	Common
	Timestamp types.ReportTimestamp
	Report    types.Report
	OK        bool
	Err       error
}

type PluginStart struct {
	Common
}

type PluginClose struct {
	Common
}

// types.OffchainConfigDigester

type ConfigDigest struct {
	Common
	ContractConfig types.ContractConfig
	Digest         types.ConfigDigest
	Err            error
}

type ConfigDigestPrefix struct {
	Common
	Prefix types.ConfigDigestPrefix
}

// type checking

var _ Trace = (*SendTo)(nil)
var _ Trace = (*Broadcast)(nil)
var _ Trace = (*Receive)(nil)
var _ Trace = (*Drop)(nil)
var _ Trace = (*EndpointStart)(nil)
var _ Trace = (*EndpointClose)(nil)

var _ Trace = (*ReadState)(nil)
var _ Trace = (*WriteState)(nil)
var _ Trace = (*ReadConfig)(nil)
var _ Trace = (*WriteConfig)(nil)
var _ Trace = (*StorePendingTransmission)(nil)
var _ Trace = (*PendingTransmissionsWithConfigDigest)(nil)
var _ Trace = (*DeletePendingTransmission)(nil)
var _ Trace = (*DeletePendingTransmissionsOlderThan)(nil)

var _ Trace = (*Notify)(nil)
var _ Trace = (*LatestConfigDetails)(nil)
var _ Trace = (*LatestConfig)(nil)
var _ Trace = (*LatestBlockHeight)(nil)

var _ Trace = (*Transmit)(nil)
var _ Trace = (*LatestConfigDigestAndEpoch)(nil)
var _ Trace = (*FromAccount)(nil)

var _ Trace = (*Query)(nil)
var _ Trace = (*Observation)(nil)
var _ Trace = (*Report)(nil)
var _ Trace = (*ShouldAcceptFinalizedReport)(nil)
var _ Trace = (*ShouldTransmitAcceptedReport)(nil)
var _ Trace = (*PluginStart)(nil)
var _ Trace = (*PluginClose)(nil)

var _ Trace = (*ConfigDigest)(nil)
var _ Trace = (*ConfigDigestPrefix)(nil)
