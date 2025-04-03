package shim

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// LimitCheckOCR3_1ReportingPlugin wraps another plugin and checks that its outputs respect
// limits. We use it to surface violations to authors of plugins as early as
// possible.
//
// It does not check inputs since those are checked by the SerializingEndpoint.
type LimitCheckOCR3_1ReportingPlugin[RI any] struct {
	Plugin ocr3_1types.ReportingPlugin[RI]
	Limits ocr3_1types.ReportingPluginLimits
}

var _ ocr3_1types.ReportingPlugin[struct{}] = LimitCheckOCR3_1ReportingPlugin[struct{}]{}

func (rp LimitCheckOCR3_1ReportingPlugin[RI]) Query(ctx context.Context, roundCtx ocr3_1types.RoundContext, kvReader ocr3_1types.KeyValueReader) (types.Query, error) {
	query, err := rp.Plugin.Query(ctx, roundCtx, kvReader)
	if err != nil {
		return nil, err
	}
	if !(len(query) <= rp.Limits.MaxQueryLength) {
		return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned oversize query (%v vs %v)", len(query), rp.Limits.MaxQueryLength)
	}
	return query, nil
}

func (rp LimitCheckOCR3_1ReportingPlugin[RI]) ObservationQuorum(ctx context.Context, roundContext ocr3_1types.RoundContext, query types.Query, aos []types.AttributedObservation, kvReader ocr3_1types.KeyValueReader) (bool, error) {
	return rp.Plugin.ObservationQuorum(ctx, roundContext, query, aos, kvReader)
}

func (rp LimitCheckOCR3_1ReportingPlugin[RI]) Observation(ctx context.Context, roundCtx ocr3_1types.RoundContext, query types.Query, kvReader ocr3_1types.KeyValueReader) (types.Observation, error) {
	observation, err := rp.Plugin.Observation(ctx, roundCtx, query, kvReader)
	if err != nil {
		return nil, err
	}
	if !(len(observation) <= rp.Limits.MaxObservationLength) {
		return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned oversize observation (%v vs %v)", len(observation), rp.Limits.MaxObservationLength)
	}
	return observation, nil
}

func (rp LimitCheckOCR3_1ReportingPlugin[RI]) ValidateObservation(ctx context.Context, roundCtx ocr3_1types.RoundContext, query types.Query, ao types.AttributedObservation, kvReader ocr3_1types.KeyValueReader) error {
	return rp.Plugin.ValidateObservation(ctx, roundCtx, query, ao, kvReader)
}

func (rp LimitCheckOCR3_1ReportingPlugin[RI]) StateTransition(ctx context.Context, roundCtx ocr3_1types.RoundContext, query types.Query, aos []types.AttributedObservation, kvReadWriter ocr3_1types.KeyValueReadWriter) (ocr3_1types.ReportsPlusPrecursor, error) {
	reportsPlusPrecursor, err := rp.Plugin.StateTransition(ctx, roundCtx, query, aos, kvReadWriter)
	if err != nil {
		return nil, err
	}

	//if !(len(reportsPlusPrecursor) <= rp.Limits.MaxReportsPlusPrecursorLength) {
	//	return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned oversize reportsPlus (%v vs %v)", len(reportsPlusPrecursor), rp.Limits.MaxReportsPlusPrecursorLength)
	//}
	return reportsPlusPrecursor, nil
}

func (rp LimitCheckOCR3_1ReportingPlugin[RI]) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[RI], error) {
	reports, err := rp.Plugin.Reports(ctx, seqNr, reportsPlusPrecursor)
	if err != nil {
		return nil, err
	}
	if !(len(reports) <= rp.Limits.MaxReportCount) {
		return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned too many reports (%v vs %v)", len(reports), rp.Limits.MaxReportCount)
	}
	for i, reportPlus := range reports {
		if !(len(reportPlus.ReportWithInfo.Report) <= rp.Limits.MaxReportLength) {
			return nil, fmt.Errorf("LimitCheckOCR3Plugin: underlying plugin returned oversize report at index %v (%v vs %v)", i, len(reportPlus.ReportWithInfo.Report), rp.Limits.MaxReportLength)
		}
	}
	return reports, nil
}

func (rp LimitCheckOCR3_1ReportingPlugin[RI]) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, report ocr3types.ReportWithInfo[RI]) (bool, error) {
	return rp.Plugin.ShouldAcceptAttestedReport(ctx, seqNr, report)
}

func (rp LimitCheckOCR3_1ReportingPlugin[RI]) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, report ocr3types.ReportWithInfo[RI]) (bool, error) {
	return rp.Plugin.ShouldTransmitAcceptedReport(ctx, seqNr, report)
}

func (rp LimitCheckOCR3_1ReportingPlugin[RI]) Close() error {
	return rp.Plugin.Close()
}
