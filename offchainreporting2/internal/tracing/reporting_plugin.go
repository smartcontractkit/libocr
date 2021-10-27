package tracing

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type pluginFactory struct {
	backend  types.ReportingPluginFactory
	tracer   *Tracer
	oracleID OracleID
}

func MakeReportingPluginFactory(tracer *Tracer, oracleID OracleID, backend types.ReportingPluginFactory) types.ReportingPluginFactory {
	return pluginFactory{backend, tracer, oracleID}
}

func (p pluginFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	backend, info, err := p.backend.NewReportingPlugin(config)
	wrappedPlugin := &plugin{backend, p.oracleID, info, p.tracer}
	return wrappedPlugin, info, err
}

type plugin struct {
	backend  types.ReportingPlugin
	oracleID OracleID
	info     types.ReportingPluginInfo
	tracer   *Tracer
}

func (p *plugin) Query(ctx context.Context, ts types.ReportTimestamp) (types.Query, error) {
	query, err := p.backend.Query(ctx, ts)
	p.tracer.Append(NewQuery(p.oracleID, ts, query, err))
	return query, err
}

func (p *plugin) Observation(ctx context.Context, ts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	obs, err := p.backend.Observation(ctx, ts, query)
	p.tracer.Append(NewObservation(p.oracleID, ts, query, obs, err))
	return obs, err
}

func (p *plugin) Report(ctx context.Context, ts types.ReportTimestamp, query types.Query, obss []types.AttributedObservation) (bool, types.Report, error) {
	ok, report, err := p.backend.Report(ctx, ts, query, obss)
	p.tracer.Append(NewReport(p.oracleID, ts, query, obss, ok, report, err))
	return ok, report, err
}

func (p *plugin) ShouldAcceptFinalizedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	ok, err := p.backend.ShouldAcceptFinalizedReport(ctx, ts, report)
	p.tracer.Append(NewShouldAcceptFinalizedReport(p.oracleID, ts, report, ok, err))
	return ok, err
}

func (p *plugin) ShouldTransmitAcceptedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	ok, err := p.backend.ShouldTransmitAcceptedReport(ctx, ts, report)
	p.tracer.Append(NewShouldAcceptFinalizedReport(p.oracleID, ts, report, ok, err))
	return ok, err
}

func (p *plugin) Start() error {
	p.tracer.Append(NewPluginStart(p.oracleID))
	return p.backend.Start()
}

func (p *plugin) Close() error {
	p.tracer.Append(NewPluginClose(p.oracleID))
	return p.backend.Close()
}
