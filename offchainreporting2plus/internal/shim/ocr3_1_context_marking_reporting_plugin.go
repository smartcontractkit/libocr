package shim

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// ContextMarkingOCR3_1ReportingPlugin marks the context it passes into the
// wrapped plugin's methods, so that facilities the plugin was handed for use
// *outside* its methods can tell when they are being misused from *inside* one.
//
// Currently the only such facility is [ocr3_1types.ReadOnlyKeyValueState]: the
// methods that are handed a KeyValueStateReader must use that one, since it
// reads at the round's sequence number, and reading anything else would
// silently lose determinism and break agreement between oracles. See
// withOutsideKeyValueStateReadsDisallowed for which methods are marked.
//

type ContextMarkingOCR3_1ReportingPlugin[RI any] struct {
	Plugin ocr3_1types.ReportingPlugin[RI]
}

var _ ocr3_1types.ReportingPlugin[struct{}] = ContextMarkingOCR3_1ReportingPlugin[struct{}]{}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) Query(ctx context.Context, seqNr uint64, kvReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Query, error) {
	return rp.Plugin.Query(withOutsideKeyValueStateReadsDisallowed(ctx), seqNr, kvReader, blobBroadcastFetcher)
}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) ObservationQuorum(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, kvReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) (bool, error) {
	return rp.Plugin.ObservationQuorum(withOutsideKeyValueStateReadsDisallowed(ctx), seqNr, aq, aos, kvReader, blobFetcher)
}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) Observation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, kvReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Observation, error) {
	return rp.Plugin.Observation(withOutsideKeyValueStateReadsDisallowed(ctx), seqNr, aq, kvReader, blobBroadcastFetcher)
}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) ValidateObservation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, ao types.AttributedObservation, kvReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) error {
	return rp.Plugin.ValidateObservation(withOutsideKeyValueStateReadsDisallowed(ctx), seqNr, aq, ao, kvReader, blobFetcher)
}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) StateTransition(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, kvReadWriter ocr3_1types.KeyValueStateReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	return rp.Plugin.StateTransition(withOutsideKeyValueStateReadsDisallowed(ctx), seqNr, aq, aos, kvReadWriter, blobFetcher)
}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) Committed(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader) error {
	return rp.Plugin.Committed(withOutsideKeyValueStateReadsDisallowed(ctx), seqNr, keyValueReader)
}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[RI], error) {
	return rp.Plugin.Reports(withOutsideKeyValueStateReadsDisallowed(ctx), seqNr, reportsPlusPrecursor)
}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	return rp.Plugin.ShouldAcceptAttestedReport(ctx, seqNr, reportWithInfo)
}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	return rp.Plugin.ShouldTransmitAcceptedReport(ctx, seqNr, reportWithInfo)
}

func (rp ContextMarkingOCR3_1ReportingPlugin[RI]) Close() error {
	return rp.Plugin.Close()
}

// outsideKeyValueStateReadsDisallowedKey marks a context as belonging to the
// execution of a ReportingPlugin method that must not read the KeyValueState
// through anything but the KeyValueStateReader it was given. The type is
// unexported so that nobody but libocr can produce the marker.
type outsideKeyValueStateReadsDisallowedKey struct{}

// withOutsideKeyValueStateReadsDisallowed marks ctx so that
// [ReadOnlyKeyValueStateWrapper] refuses to open a read transaction from it.
//
// Marked: Query, Observation, ValidateObservation, ObservationQuorum,
// StateTransition, Committed and Reports. They run at a specific sequence
// number, and those documented as pure would otherwise silently lose
// determinism.
//
// Deliberately not marked: ShouldAcceptAttestedReport and
// ShouldTransmitAcceptedReport. They are node-local transmission decisions that
// are not required to be pure, and they receive no KeyValueStateReader of their
// own, so the outside ReadOnlyKeyValueState is their only way to consult the
// replicated state. Consulting it there is legitimate.
func withOutsideKeyValueStateReadsDisallowed(ctx context.Context) context.Context {

	return context.WithValue(ctx, outsideKeyValueStateReadsDisallowedKey{}, struct{}{})
}

func outsideKeyValueStateReadsDisallowed(ctx context.Context) bool {
	return ctx.Value(outsideKeyValueStateReadsDisallowedKey{}) != nil
}
