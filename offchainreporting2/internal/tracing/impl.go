package tracing

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// commontypes.BinaryNetworkEndpoint

func NewSendTo(orig OracleID, src, dst OracleID, msg []byte) *SendTo {
	return &SendTo{Common{TypeSendTo, orig, time.Now()}, src, dst, parseMessage(msg)}
}

func NewBroadcast(orig OracleID, src OracleID, msg []byte) *Broadcast {
	return &Broadcast{Common{TypeBroadcast, orig, time.Now()}, src, parseMessage(msg)}
}

func NewReceive(orig OracleID, src, dst OracleID, msg []byte) *Receive {
	return &Receive{Common{TypeReceive, orig, time.Now()}, src, dst, parseMessage(msg)}
}

func NewDrop(orig OracleID, src, dst OracleID, msg []byte) *Drop {
	return &Drop{Common{TypeDrop, orig, time.Now()}, src, dst, parseMessage(msg)}
}

func NewEndpointStart(orig OracleID) *EndpointStart {
	return &EndpointStart{Common{TypeEndpointStart, orig, time.Now()}}
}
func NewEndpointClose(orig OracleID) *EndpointClose {
	return &EndpointClose{Common{TypeEndpointClose, orig, time.Now()}}
}

// types.Database

func NewReadState(orig OracleID, digest types.ConfigDigest, state types.PersistentState, err error) *ReadState {
	return &ReadState{Common{TypeReadState, orig, time.Now()}, digest, copyState(state), err}
}

func NewWriteState(orig OracleID, digest types.ConfigDigest, state types.PersistentState, err error) *WriteState {
	return &WriteState{Common{TypeWriteState, orig, time.Now()}, digest, copyState(state), err}
}

func NewReadConfig(orig OracleID, cfg types.ContractConfig, err error) *ReadConfig {
	return &ReadConfig{Common{TypeReadConfig, orig, time.Now()}, copyContractConfig(cfg), err}
}

func NewWriteConfig(orig OracleID, cfg types.ContractConfig, err error) *WriteConfig {
	return &WriteConfig{Common{TypeWriteConfig, orig, time.Now()}, copyContractConfig(cfg), err}
}

func NewStorePendingTransmission(orig OracleID, ts types.ReportTimestamp, tx types.PendingTransmission, err error) *StorePendingTransmission {
	return &StorePendingTransmission{Common{TypeStorePendingTransmission, orig, time.Now()}, ts, copyTransmission(tx), err}
}

func NewPendingTransmissionsWithConfigDigest(orig OracleID, digest types.ConfigDigest, err error) *PendingTransmissionsWithConfigDigest {
	return &PendingTransmissionsWithConfigDigest{Common{TypePendingTransmissionsWithConfigDigest, orig, time.Now()}, digest, err}
}

func NewDeletePendingTransmission(orig OracleID, ts types.ReportTimestamp, err error) *DeletePendingTransmission {
	return &DeletePendingTransmission{Common{TypeDeletePendingTransmission, orig, time.Now()}, ts, err}
}

func NewDeletePendingTransmissionsOlderThan(orig OracleID, cutoff time.Time, err error) *DeletePendingTransmissionsOlderThan {
	return &DeletePendingTransmissionsOlderThan{Common{TypeDeletePendingTransmissionsOlderThan, orig, time.Now()}, cutoff, err}
}

// types.ContractConfigTracker

func NewNotify(orig OracleID) *Notify {
	return &Notify{Common{TypeNotify, orig, time.Now()}}
}

func NewLatestConfigDetails(orig OracleID, changedInBlock uint64, digest types.ConfigDigest, err error) *LatestConfigDetails {
	return &LatestConfigDetails{Common{TypeLatestConfigDetails, orig, time.Now()}, changedInBlock, digest, err}
}

func NewLatestConfig(orig OracleID, changedInBlock uint64, cfg types.ContractConfig, err error) *LatestConfig {
	return &LatestConfig{Common{TypeLatestConfig, orig, time.Now()}, changedInBlock, copyContractConfig(cfg), err}
}

func NewLatestBlockHeight(orig OracleID, blockHeight uint64, err error) *LatestBlockHeight {
	return &LatestBlockHeight{Common{TypeLatestBlockHeight, orig, time.Now()}, blockHeight, err}
}

// types.ContractTransmitter

func NewTransmit(orig OracleID, repCtx types.ReportContext, report types.Report, sigs []types.AttributedOnchainSignature, err error) *Transmit {
	return &Transmit{Common{TypeTransmit, orig, time.Now()}, repCtx, copyReport(report), sigs, err}
}

func NewLatestConfigDigestAndEpoch(orig OracleID, digest types.ConfigDigest, epoch uint32, err error) *LatestConfigDigestAndEpoch {
	return &LatestConfigDigestAndEpoch{Common{TypeLatestConfigDigestAndEpoch, orig, time.Now()}, digest, epoch, err}
}

func NewFromAccount(orig OracleID, account types.Account) *FromAccount {
	return &FromAccount{Common{TypeFromAccount, orig, time.Now()}, account}
}

// types.ReportingPlugin

func NewQuery(orig OracleID, ts types.ReportTimestamp, query types.Query, err error) *Query {
	return &Query{Common{TypeQuery, orig, time.Now()}, ts, query, err}
}

func NewObservation(orig OracleID, ts types.ReportTimestamp, query types.Query, obs types.Observation, err error) *Observation {
	return &Observation{Common{TypeObservation, orig, time.Now()}, ts, query, obs, err}
}

func NewReport(orig OracleID, ts types.ReportTimestamp, query types.Query, obss []types.AttributedObservation, ok bool, report types.Report, err error) *Report {
	return &Report{Common{TypeReport, orig, time.Now()}, ts, query, obss, ok, report, err}
}

func NewShouldAcceptFinalizedReport(orig OracleID, ts types.ReportTimestamp, report types.Report, ok bool, err error) *ShouldAcceptFinalizedReport {
	return &ShouldAcceptFinalizedReport{Common{TypeShouldAcceptFinalizedReport, orig, time.Now()}, ts, report, ok, err}
}

func NewShouldTransmitAcceptedReport(orig OracleID, ts types.ReportTimestamp, report types.Report, ok bool, err error) *ShouldTransmitAcceptedReport {
	return &ShouldTransmitAcceptedReport{Common{TypeShouldTransmitAcceptedReport, orig, time.Now()}, ts, report, ok, err}
}

func NewPluginStart(orig OracleID) *PluginStart {
	return &PluginStart{Common{TypePluginStart, orig, time.Now()}}
}

func NewPluginClose(orig OracleID) *PluginClose {
	return &PluginClose{Common{TypePluginClose, orig, time.Now()}}
}

// types.OffchainConfigDigester

func NewConfigDigest(orig OracleID, cfg types.ContractConfig, digest types.ConfigDigest, err error) *ConfigDigest {
	return &ConfigDigest{Common{TypeConfigDigest, orig, time.Now()}, cfg, digest, err}
}

func NewConfigDigestPrefix(orig OracleID, prefix types.ConfigDigestPrefix) *ConfigDigestPrefix {
	return &ConfigDigestPrefix{Common{TypeConfigDigestPrefix, orig, time.Now()}, prefix}
}

// isTrace groups the tracing types and makes them instances of tracing.Trace

func (_ *SendTo) isTrace()        {}
func (_ *Broadcast) isTrace()     {}
func (_ *Receive) isTrace()       {}
func (_ *Drop) isTrace()          {}
func (_ *EndpointStart) isTrace() {}
func (_ *EndpointClose) isTrace() {}

func (_ *ReadState) isTrace()                            {}
func (_ *WriteState) isTrace()                           {}
func (_ *ReadConfig) isTrace()                           {}
func (_ *WriteConfig) isTrace()                          {}
func (_ *StorePendingTransmission) isTrace()             {}
func (_ *PendingTransmissionsWithConfigDigest) isTrace() {}
func (_ *DeletePendingTransmission) isTrace()            {}
func (_ *DeletePendingTransmissionsOlderThan) isTrace()  {}

func (_ *Notify) isTrace()              {}
func (_ *LatestConfigDetails) isTrace() {}
func (_ *LatestConfig) isTrace()        {}
func (_ *LatestBlockHeight) isTrace()   {}

func (_ *Transmit) isTrace()                   {}
func (_ *LatestConfigDigestAndEpoch) isTrace() {}
func (_ *FromAccount) isTrace()                {}

func (_ *Query) isTrace()                        {}
func (_ *Observation) isTrace()                  {}
func (_ *Report) isTrace()                       {}
func (_ *ShouldAcceptFinalizedReport) isTrace()  {}
func (_ *ShouldTransmitAcceptedReport) isTrace() {}
func (_ *PluginStart) isTrace()                  {}
func (_ *PluginClose) isTrace()                  {}

func (_ *ConfigDigest) isTrace()       {}
func (_ *ConfigDigestPrefix) isTrace() {}

// String() methods for all tracing types.

func printTrace(namespace string, c Common) string {
	return fmt.Sprintf("[%s] (%s) <%s:%s> :", c.Timestamp.Format(time.StampMilli), c.Originator, namespace, printType(c.Typ))
}

// commontypes.BinaryNetworkEndpoint

func (s *SendTo) String() string {
	msg := messageAsStr(s.Message)
	return fmt.Sprintf("%s (%s)->(%s) %s", printTrace("Net", s.Common), s.Src, s.Dst, msg)
}

func (b *Broadcast) String() string {
	msg := messageAsStr(b.Message)
	return fmt.Sprintf("%s (%s)->(*) %s", printTrace("Net", b.Common), b.Src, msg)
}

func (r *Receive) String() string {
	msg := messageAsStr(r.Message)
	return fmt.Sprintf("%s (%s)<-(%s) %s", printTrace("Net", r.Common), r.Dst, r.Src, msg)
}

func (d *Drop) String() string {
	msg := messageAsStr(d.Message)
	return fmt.Sprintf("%s (%s)-/->(%s) %s", printTrace("Net", d.Common), d.Src, d.Dst, msg)
}

func (e *EndpointStart) String() string {
	return printTrace("Net", e.Common)
}

func (e *EndpointClose) String() string {
	return printTrace("Net", e.Common)
}

// types.Database

func (r *ReadState) String() string {
	var msg string
	if r.Err != nil {
		msg = "Failed! "
	}
	cd := configDigestAsStr(r.Digest)
	msg += fmt.Sprintf("{%s=%s, err=%v}", cd, stateAsStr(r.State), r.Err)
	return fmt.Sprintf("%s %s", printTrace("DB", r.Common), msg)
}

func (w *WriteState) String() string {
	cd := configDigestAsStr(w.Digest)
	state := stateAsStr(w.State)
	var msg string
	if w.Err != nil {
		msg = "Failed! "
	}
	msg += fmt.Sprintf("{%s=%s, err=%v}", cd, state, w.Err)
	return fmt.Sprintf("%s %s", printTrace("DB", w.Common), msg)
}

func (r *ReadConfig) String() string {
	var msg string
	if r.Err != nil {
		msg = "Failed! "
	}
	msg += fmt.Sprintf("{cfg=%s, err=%v}", configAsStr(r.Config), r.Err)
	return fmt.Sprintf("%s %s", printTrace("DB", r.Common), msg)
}

func (w *WriteConfig) String() string {
	var msg string
	if w.Err != nil {
		msg = "Failed! "
	}
	msg += fmt.Sprintf("{cfg=%s, err=%v}", configAsStr(w.Config), w.Err)
	return fmt.Sprintf("%s cfg=%s", printTrace("DB", w.Common), msg)
}

func (s *StorePendingTransmission) String() string {
	var msg string
	if s.Err != nil {
		msg = fmt.Sprintf("{key=%s, tx=%s, err=%v}", timestampAsStr(s.Timestamp), transmissionAsStr(s.Transmission), s.Err)
	} else {
		msg = fmt.Sprintf("{key=%s, tx=%s}", timestampAsStr(s.Timestamp), transmissionAsStr(s.Transmission))
	}
	return fmt.Sprintf("%s %s", printTrace("DB", s.Common), msg)
}

func (s *StorePendingTransmission) EqualIgnoringTimestamp(p *StorePendingTransmission) bool {
	// Note that this does not check that the transmission timestamps are equal as these might differe slightly:
	// Ie. s.Transmission.Time.Equal(p.Transmission.Time)
	return (s.Timestamp.ConfigDigest == p.Timestamp.ConfigDigest &&
		s.Timestamp.Epoch == p.Timestamp.Epoch &&
		s.Timestamp.Round == p.Timestamp.Round &&
		s.Transmission.ExtraHash == p.Transmission.ExtraHash &&
		bytes.Equal(s.Transmission.Report, p.Transmission.Report) &&
		equalAttributedSignatures(s.Transmission.AttributedSignatures, p.Transmission.AttributedSignatures))
}

func equalAttributedSignatures(sigs1, sigs2 []types.AttributedOnchainSignature) bool {
	if len(sigs1) != len(sigs2) {
		return false
	}
	for i := 0; i < len(sigs1); i++ {
		sig1, sig2 := sigs1[i], sigs2[i]
		if !bytes.Equal(sig1.Signature, sig2.Signature) || sig1.Signer != sig2.Signer {
			return false
		}
	}
	return true
}

func (p *PendingTransmissionsWithConfigDigest) String() string {
	var msg string
	if p.Err != nil {
		msg = "Failed! "
	}
	msg += fmt.Sprintf("{digest=%s, err=%v}", configDigestAsStr(p.Digest), p.Err)
	return fmt.Sprintf("%s %s", printTrace("DB", p.Common), msg)
}

func (d *DeletePendingTransmission) String() string {
	var msg string
	if d.Err != nil {
		msg = "Failed! "
	}
	msg += fmt.Sprintf("{key=%s, err=%v}", timestampAsStr(d.Timestamp), d.Err)
	return fmt.Sprintf("%s %s", printTrace("DB", d.Common), msg)
}

func (d *DeletePendingTransmissionsOlderThan) String() string {
	var msg string
	if d.Err != nil {
		msg = "Failed! "
	}
	msg += fmt.Sprintf("{cutoff=%s, err=%v}", d.Cutoff.Format(time.StampMilli), d.Err)
	return fmt.Sprintf("%s %s", printTrace("DB", d.Common), msg)
}

// types.ContractConfigTracker

func (n *Notify) String() string {
	return printTrace("ConfigTracker", n.Common)
}

func (l *LatestConfigDetails) String() string {
	var msg string
	if l.Err != nil {
		msg = fmt.Sprintf("failed to get the latest config details: %s", l.Err)
	} else {
		msg = fmt.Sprintf("{height=%d, digest=%s}", l.ChangedInBlock, configDigestAsStr(l.Digest))
	}
	return fmt.Sprintf("%s %s", printTrace("ConfigTracker", l.Common), msg)
}

func (l *LatestConfig) String() string {
	var msg string
	if l.Err != nil {
		msg = fmt.Sprintf("failed to get the latest contract config: %s", l.Err)
	} else {
		msg = fmt.Sprintf("{config=%s}", configAsStr(l.Config))
	}
	return fmt.Sprintf("%s %s", printTrace("ConfigTracker", l.Common), msg)
}

func (l *LatestBlockHeight) String() string {
	var msg string
	if l.Err != nil {
		msg = fmt.Sprintf("failed to get block height: %s", l.Err)
	} else {
		msg = fmt.Sprintf("{height=%d}", l.BlockHeight)
	}
	return fmt.Sprintf("%s %s", printTrace("ConfigTracker", l.Common), msg)
}

// types.ContractTransmitter

func (t *Transmit) String() string {
	var msg string
	if t.Err != nil {
		msg = fmt.Sprintf("failed to transmit report: %s", t.Err)
	} else {
		msg = fmt.Sprintf("Report transmitted {context={epoch=%v, round=%v, digest=%s}, report=%s}",
			t.ReportContext.Epoch, t.ReportContext.Round, configDigestAsStr(t.ReportContext.ConfigDigest),
			reportAsStr(t.Report))
	}
	return fmt.Sprintf("%s %s", printTrace("Transmitter", t.Common), msg)
}

func (l *LatestConfigDigestAndEpoch) String() string {
	var msg string
	if l.Err != nil {
		msg = fmt.Sprintf("failed to get latest config digest and epoch: %s", l.Err)
	} else {
		msg = fmt.Sprintf("{digest=%s, epoch=%d}", configDigestAsStr(l.Digest), l.Epoch)
	}
	return fmt.Sprintf("%s %s", printTrace("Transmitter", l.Common), msg)
}

func (f *FromAccount) String() string {
	return fmt.Sprintf("%s {account=%s}", printTrace("Transmitter", f.Common), accountAsStr(f.Account))
}

// types.ReportingPlugin

func (q *Query) String() string {
	var msg string
	if q.Err != nil {
		msg = fmt.Sprintf("call to query failed with error: %s", q.Err)
	} else {
		msg = fmt.Sprintf("{ts=%s, query=%s}", timestampAsStr(q.Timestamp), queryAsStr(q.Query))
	}
	return fmt.Sprintf("%s %s", printTrace("Plugin", q.Common), msg)
}

func (o *Observation) String() string {
	var msg string
	if o.Err != nil {
		msg = fmt.Sprintf("failed to compute observation with error: %s", o.Err)
	} else {
		msg = fmt.Sprintf("{ts=%s, query=%s, obs=%s}", timestampAsStr(o.Timestamp), queryAsStr(o.Query), observationAsStr(o.Observation))
	}
	return fmt.Sprintf("%s %s", printTrace("Plugin", o.Common), msg)
}

func (r *Report) String() string {
	var msg string
	if r.Err != nil {
		msg = fmt.Sprintf("failed to build report: %s", r.Err)
	} else {
		msg = fmt.Sprintf("{ts=%s, query=%s, ok=%t, report=%s}", timestampAsStr(r.Timestamp), queryAsStr(r.Query), r.OK, reportAsStr(r.Report))
	}
	return fmt.Sprintf("%s %s", printTrace("Plugin", r.Common), msg)
}

func (s *ShouldAcceptFinalizedReport) String() string {
	var msg string
	if s.Err != nil {
		msg = fmt.Sprintf("ShouldAcceptFinalizedReport() failed: %s", s.Err)
	} else {
		msg = fmt.Sprintf("{ts=%s, report=%s, ok=%t}", timestampAsStr(s.Timestamp), reportAsStr(s.Report), s.OK)
	}
	return fmt.Sprintf("%s %s", printTrace("Plugin", s.Common), msg)
}

func (s *ShouldTransmitAcceptedReport) String() string {
	var msg string
	if s.Err != nil {
		msg = fmt.Sprintf("ShouldTransmitAcceptedReport() failed: %s", s.Err)
	} else {
		msg = fmt.Sprintf("{ts=%s, report=%s, ok=%t}", timestampAsStr(s.Timestamp), reportAsStr(s.Report), s.OK)
	}
	return fmt.Sprintf("%s %s", printTrace("Plugin", s.Common), msg)
}

func (p *PluginStart) String() string {
	return printTrace("Plugin", p.Common)
}

func (p *PluginClose) String() string {
	return printTrace("Plugin", p.Common)
}

// types.OffchainConfigDigester

func (c *ConfigDigest) String() string {
	var msg string
	if c.Err != nil {
		msg = fmt.Sprintf("ConfigDigest() failed: %s", c.Err)
	} else {
		msg = fmt.Sprintf("{cfg=%s, digest=%s}", configAsStr(c.ContractConfig), configDigestAsStr(c.Digest))
	}
	return fmt.Sprintf("%s %s", printTrace("Digester", c.Common), msg)
}

func (c *ConfigDigestPrefix) String() string {
	return fmt.Sprintf("%s prefix=%d", printTrace("Digester", c.Common), c.Prefix)
}

// Helpers

func printType(typ TraceType) string {
	typeName, ok := map[TraceType]string{
		TypeSendTo:        "SendTo",
		TypeBroadcast:     "Broadcast",
		TypeReceive:       "Receive",
		TypeDrop:          "Drop",
		TypeEndpointStart: "EndpointStart",
		TypeEndpointClose: "EndpointClose",

		// types.Database
		TypeReadState:                            "ReadState",
		TypeWriteState:                           "WriteState",
		TypeReadConfig:                           "ReadConfig",
		TypeWriteConfig:                          "WriteConfig",
		TypeStorePendingTransmission:             "StorePendingTransmission",
		TypePendingTransmissionsWithConfigDigest: "PendintTransmissionsWithConfigDigest",
		TypeDeletePendingTransmission:            "DeletePendingTransmission",
		TypeDeletePendingTransmissionsOlderThan:  "DeletePendingTransmissionsOlderThan",

		// types.ContractConfigTracker
		TypeNotify:              "Notify",
		TypeLatestConfigDetails: "LatestConfigDetails",
		TypeLatestConfig:        "LatestConfig",
		TypeLatestBlockHeight:   "LatestBlockHeight",

		// Transmission contract operations
		TypeTransmit:                   "Transmit",
		TypeLatestConfigDigestAndEpoch: "LatestConfigDigestAndEpoch",
		TypeFromAccount:                "FromAccount",

		// types.ReportingPlugin
		TypeQuery:                        "Query",
		TypeObservation:                  "Observation",
		TypeReport:                       "Report",
		TypeShouldAcceptFinalizedReport:  "ShouldAcceptFinalizedReport",
		TypeShouldTransmitAcceptedReport: "ShouldTransmitAcceptedReport",
		TypePluginStart:                  "PluginStart",
		TypePluginClose:                  "PluginClose",

		// types.OffchainConfigDigester
		TypeConfigDigest:       "ConfigDigest",
		TypeConfigDigestPrefix: "ConfigDigestPrefix",
	}[typ]
	if !ok {
		return fmt.Sprintf("unknown TraceType %d", typ)
	}
	return typeName
}

func parseMessage(buf []byte) protocol.Message {
	msg, _, err := serialization.Deserialize(buf)
	if err != nil {
		panic(err)
	}
	return msg
}

func observationAsStr(obs types.Observation) string {
	src := []byte(obs)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	o := string(dst)
	return o[:2] + ".." + o[len(o)-4:]
}

func messageAsStr(raw protocol.Message) string {
	switch msg := raw.(type) {
	case protocol.MessageNewEpoch:
		return fmt.Sprintf("NewEpoch{epoch=%d}", msg.Epoch)
	case protocol.MessageObserveReq:
		return fmt.Sprintf("ObserveReq{epoch=%d, round=%d}", msg.Epoch, msg.Round)
	case protocol.MessageObserve:
		return fmt.Sprintf("Observe{epoch=%d, round=%d, obs=%s}", msg.Epoch, msg.Round, observationAsStr(msg.SignedObservation.Observation))
	case protocol.MessageReportReq:
		return fmt.Sprintf("ReportReq{epoch=%d, round=%d, obss=%s}", msg.Epoch, msg.Round, observationsAsStr(msg.AttributedSignedObservations))
	case protocol.MessageReport:
		return fmt.Sprintf("Report{epoch=%d, round=%d, report=%s}", msg.Epoch, msg.Round, reportAsStr(msg.AttestedReport.Report))
	case protocol.MessageFinal:
		return fmt.Sprintf("Final{epoch=%d, round=%d, report=%s}", msg.Epoch, msg.Round, reportAsStr(types.Report(msg.AttestedReport.Report)))
	case protocol.MessageFinalEcho:
		return fmt.Sprintf("FinalEcho{epoch=%d, round=%d, report=%s}", msg.Epoch, msg.Round, reportAsStr(types.Report(msg.AttestedReport.Report)))
	default:
		panic(fmt.Sprintf("unsupported message type %T", raw))
	}
}

func observationsAsStr(os []protocol.AttributedSignedObservation) string {
	var b strings.Builder
	b.WriteByte('{')
	for _, o := range os {
		fmt.Fprintf(&b, "%d=%s,", o.Observer, observationAsStr(o.SignedObservation.Observation))
	}
	b.WriteByte('}')
	return b.String()
}

func configDigestAsStr(cd types.ConfigDigest) string {
	h := cd.Hex()
	if len(h) > 6 {
		return h[:2] + ".." + h[len(h)-4:]
	}
	return h
}

func stateAsStr(state types.PersistentState) string {
	return fmt.Sprintf("{epoch=%d, highestSent=%d, highestReceived=%v}", state.Epoch, state.HighestSentEpoch, state.HighestReceivedEpoch)
}

func configAsStr(cfg types.ContractConfig) string {
	return fmt.Sprintf("{digest=%s, count=%d, f=%d, version=%d}",
		configDigestAsStr(cfg.ConfigDigest), cfg.ConfigCount, cfg.F, cfg.OffchainConfigVersion)
}

func timestampAsStr(ts types.ReportTimestamp) string {
	return fmt.Sprintf("{digest=%s, epoch=%d, round=%d}",
		configDigestAsStr(ts.ConfigDigest), ts.Epoch, ts.Round)
}

func reportAsStr(report types.Report) string {
	src := []byte(report)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	r := string(dst)
	if len(r) > 6 {
		return r[:2] + ".." + r[len(r)-4:]
	}
	if len(r) == 0 {
		return "<empty>"
	}
	return r
}

func transmissionAsStr(tr types.PendingTransmission) string {
	return fmt.Sprintf("{time=%s, report=%s}", tr.Time.Format(time.StampMilli), reportAsStr(tr.Report))
}

func accountAsStr(acc types.Account) string {
	if len(string(acc)) > 6 {
		return string(acc)[:6] + ".."
	}
	if len(string(acc)) == 0 {
		return "<empty>"
	}
	return string(acc)
}

func queryAsStr(query types.Query) string {
	if len(string(query)) > 6 {
		return string(query)[:6] + ".."
	}
	if len(string(query)) == 0 {
		return "<empty>"
	}
	return string(query)
}

func copyState(st types.PersistentState) types.PersistentState {
	return types.PersistentState{
		Epoch:                st.Epoch,
		HighestSentEpoch:     st.HighestSentEpoch,
		HighestReceivedEpoch: st.HighestReceivedEpoch[:],
	}
}

func copyContractConfig(cfg types.ContractConfig) types.ContractConfig {
	return types.ContractConfig{
		ConfigDigest:          cfg.ConfigDigest,
		ConfigCount:           cfg.ConfigCount,
		Signers:               cfg.Signers[:],
		Transmitters:          cfg.Transmitters[:],
		F:                     cfg.F,
		OnchainConfig:         cfg.OnchainConfig[:],
		OffchainConfigVersion: cfg.OffchainConfigVersion,
		OffchainConfig:        cfg.OffchainConfig[:],
	}
}

func copyTransmission(tx types.PendingTransmission) types.PendingTransmission {
	return types.PendingTransmission{
		Time:                 tx.Time,
		ExtraHash:            tx.ExtraHash,
		Report:               types.Report([]byte(tx.Report)[:]),
		AttributedSignatures: tx.AttributedSignatures[:],
	}
}

func copyReport(rep types.Report) []byte {
	return types.Report([]byte(rep)[:])
}
