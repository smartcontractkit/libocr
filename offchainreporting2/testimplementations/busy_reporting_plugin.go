package testimplementations

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type busyReportingPluginFactory struct {
	observationTimeout time.Duration
	uniqueReports      bool
}

func NewBusyReportingPluginFactory(observationTimeout time.Duration, uniqueReports bool) types.ReportingPluginFactory {
	return &busyReportingPluginFactory{observationTimeout, uniqueReports}
}

func (b *busyReportingPluginFactory) NewReportingPlugin(cfg types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	info := types.ReportingPluginInfo{
		"busy reporting plugin", // Name:
		b.uniqueReports,         // UniqueReports:
		1024,                    // 1KiBi - MaxQueryLen:
		1024,                    // 1KiBi - MaxObservationLen:
		100 * 1024,              // 100KiBi - MaxReportLen:
	}
	return &busyReportingPlugin{cfg, b.observationTimeout}, info, nil
}

type busyReportingPlugin struct {
	config             types.ReportingPluginConfig
	observationTimeout time.Duration
}

// Query produces a buffer of 10 random bytes.
// All instances of this plugin will produce the same query for a given timestamp.
func (b *busyReportingPlugin) Query(_ context.Context, ts types.ReportTimestamp) (types.Query, error) {
	hmd5 := md5.New()
	hmd5.Write(ts.ConfigDigest[:])
	_ = binary.Write(hmd5, binary.BigEndian, ts.Epoch)
	_ = binary.Write(hmd5, binary.BigEndian, ts.Round)
	return types.Query(hmd5.Sum(nil)), nil
}

// Observation concatenates the query and 10 random bytes to produce a types.Objservation
func (b *busyReportingPlugin) Observation(_ context.Context, _ types.ReportTimestamp, query types.Query) (types.Observation, error) {
	time.Sleep(b.observationTimeout)
	obs := make([]byte, 10)
	if _, err := rand.Read(obs); err != nil {
		return nil, err
	}
	out := append(append([]byte{}, query...), obs...)
	return types.Observation(out), nil
}

// Report concatenates all observations.
func (b *busyReportingPlugin) Report(_ context.Context, _ types.ReportTimestamp, _ types.Query, obss []types.AttributedObservation) (bool, types.Report, error) {
	if len(obss) == 0 {
		return false, nil, fmt.Errorf("no observations received")
	}
	report := []byte{}
	for _, obs := range obss {
		report = append(report, []byte(obs.Observation)...)
	}
	return true, types.Report(report), nil
}

// noops

func (b *busyReportingPlugin) ShouldAcceptFinalizedReport(_ context.Context, _ types.ReportTimestamp, _ types.Report) (bool, error) {
	return true, nil
}

func (b *busyReportingPlugin) ShouldTransmitAcceptedReport(_ context.Context, _ types.ReportTimestamp, _ types.Report) (bool, error) {
	return true, nil
}

func (b *busyReportingPlugin) Start() error {
	return nil
}

func (b *busyReportingPlugin) Close() error {
	return nil
}
