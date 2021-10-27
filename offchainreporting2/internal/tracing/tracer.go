package tracing

import (
	"fmt"
	"strings"
	"sync"
)

type Tracer struct {
	mu      sync.Mutex
	raw     []Trace
	hs      *hooks
	started bool
}

func NewTracer() *Tracer {
	return &Tracer{
		sync.Mutex{},
		[]Trace{},
		newHooks(),
		false,
	}
}

func (t *Tracer) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.started = true
}

func (t *Tracer) Append(frame Trace) {
	t.mu.Lock()
	if !t.started {
		t.mu.Unlock()
		return
	}
	t.raw = append(t.raw, frame)
	t.mu.Unlock()
	t.hs.execute(frame)
}

func (t *Tracer) String() string {
	var b strings.Builder
	for _, frame := range t.raw {
		b.WriteString(frame.String())
		b.WriteByte('\n')
	}
	return b.String()
}

func (t *Tracer) RegisterHook(fn Hook) {
	t.hs.register(fn)
}

func (t *Tracer) GetTraces() []Trace {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.raw
}

type TraceStats struct {
	Aggregates map[TraceType]map[OracleID]uint
}

func (t *TraceStats) Inc(typ TraceType, orig OracleID) {
	if _, ok := t.Aggregates[typ]; !ok {
		t.Aggregates[typ] = make(map[OracleID]uint)
	}
	if _, ok := t.Aggregates[typ][orig]; !ok {
		t.Aggregates[typ][orig] = 0
	}
	t.Aggregates[typ][orig] += 1
}

func (t *TraceStats) String() string {
	var b strings.Builder
	for typ, aggs := range t.Aggregates {
		fmt.Fprintf(&b, "<%s>: ", printType(typ))
		for oid, counts := range aggs {
			fmt.Fprintf(&b, "(%s) %d, ", oid, counts)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func (t *Tracer) Stats() *TraceStats {
	stats := &TraceStats{make(map[TraceType]map[OracleID]uint)}
	t.mu.Lock()
	for _, raw := range t.raw {
		switch f := raw.(type) {
		case *SendTo:
			stats.Inc(f.Typ, f.Originator)
		case *Broadcast:
			stats.Inc(f.Typ, f.Originator)
		case *Receive:
			stats.Inc(f.Typ, f.Originator)
		case *Drop:
			stats.Inc(f.Typ, f.Originator)
		case *EndpointStart:
			stats.Inc(f.Typ, f.Originator)
		case *EndpointClose:
			stats.Inc(f.Typ, f.Originator)

		case *ReadState:
			stats.Inc(f.Typ, f.Originator)
		case *WriteState:
			stats.Inc(f.Typ, f.Originator)
		case *ReadConfig:
			stats.Inc(f.Typ, f.Originator)
		case *WriteConfig:
			stats.Inc(f.Typ, f.Originator)
		case *StorePendingTransmission:
			stats.Inc(f.Typ, f.Originator)
		case *PendingTransmissionsWithConfigDigest:
			stats.Inc(f.Typ, f.Originator)
		case *DeletePendingTransmission:
			stats.Inc(f.Typ, f.Originator)
		case *DeletePendingTransmissionsOlderThan:
			stats.Inc(f.Typ, f.Originator)

		case *Notify:
			stats.Inc(f.Typ, f.Originator)
		case *LatestConfigDetails:
			stats.Inc(f.Typ, f.Originator)
		case *LatestConfig:
			stats.Inc(f.Typ, f.Originator)
		case *LatestBlockHeight:
			stats.Inc(f.Typ, f.Originator)

		case *Transmit:
			stats.Inc(f.Typ, f.Originator)
		case *LatestConfigDigestAndEpoch:
			stats.Inc(f.Typ, f.Originator)
		case *FromAccount:
			stats.Inc(f.Typ, f.Originator)

		case *Query:
			stats.Inc(f.Typ, f.Originator)
		case *Observation:
			stats.Inc(f.Typ, f.Originator)
		case *Report:
			stats.Inc(f.Typ, f.Originator)
		case *ShouldAcceptFinalizedReport:
			stats.Inc(f.Typ, f.Originator)
		case *ShouldTransmitAcceptedReport:
			stats.Inc(f.Typ, f.Originator)
		case *PluginStart:
			stats.Inc(f.Typ, f.Originator)
		case *PluginClose:
			stats.Inc(f.Typ, f.Originator)

		case *ConfigDigest:
			stats.Inc(f.Typ, f.Originator)
		case *ConfigDigestPrefix:
			stats.Inc(f.Typ, f.Originator)

		default:
			panic(fmt.Sprintf("unrecognised type %T", raw))
		}
	}
	t.mu.Unlock()
	return stats
}
