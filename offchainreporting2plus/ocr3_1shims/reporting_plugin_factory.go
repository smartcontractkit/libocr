package ocr3_1shims

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

type reportingPluginFactory2Shim[RI any] struct {
	wrapped ocr3_1types.ReportingPluginFactory[RI]
}

func (f *reportingPluginFactory2Shim[RI]) NewReportingPlugin(
	ctx context.Context,
	config ocr3types.ReportingPluginConfig,
	blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher,
	_ ocr3_1types.ReadOnlyKeyValueState,
) (ocr3_1types.ReportingPlugin[RI], ocr3_1types.ReportingPluginInfo, error) {
	return f.wrapped.NewReportingPlugin(ctx, config, blobBroadcastFetcher)
}

// ReportingPluginFactoryAsReportingPluginFactory2 wraps an
// [ocr3_1types.ReportingPluginFactory] and returns an
// [ocr3_1types.ReportingPluginFactory2] that drops the
// [ocr3_1types.ReadOnlyKeyValueState].
//
// Unlike the other shims there is no "already implements the newer interface"
// fast path: both interfaces declare NewReportingPlugin with different
// signatures, so no concrete type can satisfy both.
func ReportingPluginFactoryAsReportingPluginFactory2[RI any](
	f ocr3_1types.ReportingPluginFactory[RI],
) ocr3_1types.ReportingPluginFactory2[RI] {
	return &reportingPluginFactory2Shim[RI]{f}
}
