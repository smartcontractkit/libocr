package ocr3plugintesthelpers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocr3types.ReportingPlugin[any] = (*reportingPluginDeterminismWrapper[any])(nil)

// Result is the (value, error) returned by a single wrapped call.
type Result struct {
	Value any
	Err   error
}

// Equal reports whether r and other are equivalent for the purposes of determinism checking:
// errors must match (by message), and when there is no error the values must be deeply equal.
// When an error is present the values are ignored, mirroring the rule used by checkDeterminism.
func (r Result) Equal(other Result) bool {
	if !errorsEqualByString(r.Err, other.Err) {
		return false
	}
	if r.Err != nil {
		return true
	}
	return reflect.DeepEqual(r.Value, other.Value)
}

// Call is a single invocation of a wrapped method, tagged with its zero-based index within the
// repetition sequence and the Result it returned.
type Call struct {
	Index  int    // zero-based index of the call within the repetition sequence
	Result Result // value and error returned by that call
}

// ReportingPluginNondeterminismDetected is the value panicked with when a wrapped method returns different
// results across repetitions. It carries the method name, the arguments of the call, and the two
// Calls that disagreed — CallA is always the first call (index 0) and CallB is the first
// mismatching call — so the failure can be diagnosed without re-running. Args is the call's
// non-context arguments as a []any, in the order they were passed to the wrapped method.
type ReportingPluginNondeterminismDetected struct {
	FuncName string // name of the ReportingPlugin method that was nondeterministic
	Args     []any  // the non-context arguments passed to the wrapped call, in order
	CallA    Call   // the first call (index 0)
	CallB    Call   // the first mismatching call
}

func (n ReportingPluginNondeterminismDetected) String() string {
	return fmt.Sprintf(
		"non-deterministic %s: args=%+v; call %d = (%+v, %v); call %d = (%+v, %v)",
		n.FuncName, n.Args,
		n.CallA.Index, n.CallA.Result.Value, n.CallA.Result.Err,
		n.CallB.Index, n.CallB.Result.Value, n.CallB.Result.Err,
	)
}

type reportingPluginDeterminismWrapper[RI any] struct {
	underlying                   ocr3types.ReportingPlugin[RI]
	numAdditionalCallRepetitions int
}

// NewReportingPluginDeterminismWrapperForTests wraps the four ReportingPlugin functions documented
// as pure to ensure determinism, calling each one numAdditionalCallRepetitions extra times and
// comparing the results. Upon detecting nondeterminism it panics with a
// ReportingPluginNondeterminismDetected value, which fails the surrounding test.
//
// Not all forms of nondeterminism can be captured by this wrapper; for example race conditions
// or nondeterminism derived from node state are not always caught. It targets nondeterminism
// such as that caused by iterating a map, which is common in serialization.
// Returned values are not snapshotted, so a plugin that mutates memory backing a previously
// returned value may evade detection.
//
// Typical usage in a test:
//
//	wrapped := ocr3plugintesthelpers.NewReportingPluginDeterminismWrapperForTests(plugin, 10)
//	_, err := wrapped.Outcome(ctx, ocr3types.OutcomeContext{}, types.Query{}, nil)
//	// If Outcome is nondeterministic, the call panics with a ReportingPluginNondeterminismDetected
//	// value, failing the test.
//
// ONLY FOR TESTING. NEVER IN PRODUCTION.
func NewReportingPluginDeterminismWrapperForTests[RI any](underlying ocr3types.ReportingPlugin[RI], numAdditionalCallRepetitions int) ocr3types.ReportingPlugin[RI] {
	if numAdditionalCallRepetitions < 1 {
		panic(fmt.Sprintf("numAdditionalCallRepetitions must be >= 1, got %d", numAdditionalCallRepetitions))
	}
	return &reportingPluginDeterminismWrapper[RI]{
		underlying,
		numAdditionalCallRepetitions,
	}
}

func (w *reportingPluginDeterminismWrapper[RI]) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
	return w.underlying.Query(ctx, outctx)
}

func (w *reportingPluginDeterminismWrapper[RI]) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query) (types.Observation, error) {
	return w.underlying.Observation(ctx, outctx, query)
}

func (w *reportingPluginDeterminismWrapper[RI]) ValidateObservation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, ao types.AttributedObservation) error {
	_, err := checkDeterminism("ValidateObservation", []any{outctx, query, ao}, w.numAdditionalCallRepetitions, func() (struct{}, error) {
		return struct{}{}, w.underlying.ValidateObservation(ctx, outctx, query, ao)
	})
	return err
}

func (w *reportingPluginDeterminismWrapper[RI]) ObservationQuorum(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (bool, error) {
	return checkDeterminism("ObservationQuorum", []any{outctx, query, aos}, w.numAdditionalCallRepetitions, func() (bool, error) {
		return w.underlying.ObservationQuorum(ctx, outctx, query, aos)
	})
}

func (w *reportingPluginDeterminismWrapper[RI]) Outcome(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	return checkDeterminism("Outcome", []any{outctx, query, aos}, w.numAdditionalCallRepetitions, func() (ocr3types.Outcome, error) {
		return w.underlying.Outcome(ctx, outctx, query, aos)
	})
}

func (w *reportingPluginDeterminismWrapper[RI]) Reports(ctx context.Context, seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportPlus[RI], error) {
	return checkDeterminism("Reports", []any{seqNr, outcome}, w.numAdditionalCallRepetitions, func() ([]ocr3types.ReportPlus[RI], error) {
		return w.underlying.Reports(ctx, seqNr, outcome)
	})
}

func (w *reportingPluginDeterminismWrapper[RI]) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	return w.underlying.ShouldAcceptAttestedReport(ctx, seqNr, reportWithInfo)
}

func (w *reportingPluginDeterminismWrapper[RI]) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	return w.underlying.ShouldTransmitAcceptedReport(ctx, seqNr, reportWithInfo)
}

func (w *reportingPluginDeterminismWrapper[RI]) Close() error {
	return w.underlying.Close()
}

type reportingPluginFactoryDeterminismWrapper[RI any] struct {
	underlying                   ocr3types.ReportingPluginFactory[RI]
	numAdditionalCallRepetitions int
}

// NewReportingPluginFactoryDeterminismWrapperForTests wraps the given ReportingPluginFactory so
// that each plugin instance it creates verifies determinism by performing
// numAdditionalCallRepetitions additional calls to each pure method and comparing their results.
// Detected nondeterminism is surfaced as a ReportingPluginNondeterminismDetected panic.
//
// Use it to produce a ReportingPluginFactory that can be wired into OCR3OracleArgs (or
// OCR3OracleArgs2), so a real OCR instance drives the plugin and the determinism check runs
// against actual protocol traffic — handy for smoke testing. See NewReportingPluginDeterminismWrapperForTests
// for a usage example of the underlying wrapper.
//
// ONLY FOR TESTING. NEVER IN PRODUCTION.
func NewReportingPluginFactoryDeterminismWrapperForTests[RI any](underlying ocr3types.ReportingPluginFactory[RI], numAdditionalCallRepetitions int) ocr3types.ReportingPluginFactory[RI] {
	if numAdditionalCallRepetitions < 1 {
		panic(fmt.Sprintf("numAdditionalCallRepetitions must be >= 1, got %d", numAdditionalCallRepetitions))
	}
	return &reportingPluginFactoryDeterminismWrapper[RI]{
		underlying,
		numAdditionalCallRepetitions,
	}
}

func (f *reportingPluginFactoryDeterminismWrapper[RI]) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[RI], ocr3types.ReportingPluginInfo, error) {
	plugin, info, err := f.underlying.NewReportingPlugin(ctx, config)
	if err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, err
	}
	return &reportingPluginDeterminismWrapper[RI]{
		plugin,
		f.numAdditionalCallRepetitions,
	}, info, nil
}

// checkDeterminism calls the function numAdditionalCallRepetitions extra times and compares each
// result against the first via Result.Equal. On a mismatch it panics with a
// ReportingPluginNondeterminismDetected value carrying the method name, the call's arguments, and the two
// Calls that disagreed. On success, the first call's (result, error) is returned unchanged so
// regular plugin errors propagate as usual.
func checkDeterminism[T any](funcName string, args []any, numAdditionalCallRepetitions int, call func() (T, error)) (T, error) {
	v0, e0 := call()
	a := Call{0, Result{v0, e0}}
	for i := 1; i <= numAdditionalCallRepetitions; i++ {
		vi, ei := call()
		b := Call{i, Result{vi, ei}}
		if !a.Result.Equal(b.Result) {
			panic(ReportingPluginNondeterminismDetected{
				funcName,
				args,
				a,
				b,
			})
		}
	}
	return v0, e0
}

func errorsEqualByString(err1, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}

	if err1 == nil || err2 == nil {
		return false
	}

	return err1.Error() == err2.Error()
}
