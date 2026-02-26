package median

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type mockMedianContract struct {
	configDigest    types.ConfigDigest
	epoch           uint32
	round           uint8
	latestAnswer    *big.Int
	latestTimestamp time.Time
}

func (m *mockMedianContract) LatestTransmissionDetails(ctx context.Context) (
	types.ConfigDigest, uint32, uint8, *big.Int, time.Time, error,
) {
	return m.configDigest, m.epoch, m.round, m.latestAnswer, m.latestTimestamp, nil
}

func (m *mockMedianContract) LatestRoundRequested(ctx context.Context, lookback time.Duration) (
	types.ConfigDigest, uint32, uint8, error,
) {
	return types.ConfigDigest{}, 0, 0, nil
}

type mockReportCodec struct{}

func (m *mockReportCodec) BuildReport(_ context.Context, paos []ParsedAttributedObservation) (types.Report, error) {
	median := paos[len(paos)/2].Value
	encoded, err := EncodeValue(median)
	if err != nil {
		return nil, err
	}
	return types.Report(encoded), nil
}

func (m *mockReportCodec) MedianFromReport(_ context.Context, report types.Report) (*big.Int, error) {
	return DecodeValue(report)
}

func (m *mockReportCodec) MaxReportLength(_ context.Context, _ int) (int, error) {
	return byteWidth, nil
}

type mockDataSource struct {
	value *big.Int
}

func (m *mockDataSource) Observe(_ context.Context, _ types.ReportTimestamp) (*big.Int, error) {
	return m.value, nil
}

type noopLogger struct{}

func (noopLogger) Trace(_ string, _ commontypes.LogFields)    {}
func (noopLogger) Debug(_ string, _ commontypes.LogFields)    {}
func (noopLogger) Info(_ string, _ commontypes.LogFields)     {}
func (noopLogger) Warn(_ string, _ commontypes.LogFields)     {}
func (noopLogger) Error(_ string, _ commontypes.LogFields)    {}
func (noopLogger) Critical(_ string, _ commontypes.LogFields) {}

var testConfigDigest = types.ConfigDigest{0x01}

func buildReport(t *testing.T, value *big.Int) types.Report {
	t.Helper()
	encoded, err := EncodeValue(value)
	if err != nil {
		t.Fatalf("failed to encode value: %v", err)
	}
	return types.Report(encoded)
}

func newTestPlugin(contract *mockMedianContract, deltaC time.Duration) *numericalMedian {
	return &numericalMedian{
		offchainConfig: OffchainConfig{
			AlphaReportInfinite: true,
			AlphaReportPPB:      0,
			AlphaAcceptInfinite: true,
			AlphaAcceptPPB:      0,
			DeltaC:              deltaC,
		},
		onchainConfig: OnchainConfig{
			Min: big.NewInt(-1e18),
			Max: big.NewInt(1e18),
		},
		contractTransmitter:                  contract,
		dataSource:                           &mockDataSource{big.NewInt(100)},
		juelsPerFeeCoinDataSource:            &mockDataSource{big.NewInt(1)},
		gasPriceSubunitsDataSource:           &mockDataSource{big.NewInt(1)},
		includeGasPriceSubunitsInObservation: false,
		logger:                               loghelper.MakeRootLoggerWithContext(noopLogger{}),
		reportCodec:                          &mockReportCodec{},
		deviationFunc:                        DefaultDeviationFunc,

		configDigest:             testConfigDigest,
		f:                        1,
		latestAcceptedEpochRound: epochRound{},
		latestAcceptedMedian:     new(big.Int),
		latestAcceptedAt:         time.Time{},
		maxReportLength:          byteWidth,
	}
}

// TestShouldAcceptFinalizedReport_PendingTooOldRecovery verifies the fix:
// after a silent TX failure, once DeltaC has elapsed the plugin treats the
// pending report as expired and accepts new reports again.
func TestShouldAcceptFinalizedReport_PendingTooOldRecovery(t *testing.T) {
	contract := &mockMedianContract{
		configDigest:    testConfigDigest,
		epoch:           100,
		round:           5,
		latestAnswer:    big.NewInt(100),
		latestTimestamp: time.Now().Add(-60 * time.Second),
	}

	// Use a tiny DeltaC so the test doesn't need to sleep
	plugin := newTestPlugin(contract, 50*time.Millisecond)
	ctx := context.Background()
	report := buildReport(t, big.NewInt(100))

	// Step 1: Accept first report
	repts1 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 101, Round: 1}
	accepted, err := plugin.ShouldAcceptFinalizedReport(ctx, repts1, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !accepted {
		t.Fatal("expected first report to be accepted")
	}
	if plugin.latestAcceptedEpochRound != (epochRound{101, 1}) {
		t.Fatalf("expected latestAcceptedEpochRound={101,1}, got %+v", plugin.latestAcceptedEpochRound)
	}

	// Step 2: TX fails silently -- contract state stays at {100, 5}

	// Step 3: Immediately, new report is rejected (pending not yet expired)
	repts2 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 105, Round: 1}
	accepted, err = plugin.ShouldAcceptFinalizedReport(ctx, repts2, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accepted {
		t.Fatal("expected report to be rejected while pending is fresh")
	}

	// Step 4: Wait for DeltaC to expire
	time.Sleep(60 * time.Millisecond)

	// Step 5: Now the pending report is too old -- should accept
	repts3 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 110, Round: 1}
	accepted, err = plugin.ShouldAcceptFinalizedReport(ctx, repts3, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !accepted {
		t.Fatal("expected report to be accepted after pendingTooOld timeout")
	}
	t.Log("Recovery confirmed: stable feed accepted after DeltaC expiry")
}

// TestShouldAcceptFinalizedReport_PendingCheckStillWorksWithinDeltaC verifies
// that the nothingPending optimization still prevents duplicate transmissions
// when a report has been accepted recently (within DeltaC).
func TestShouldAcceptFinalizedReport_PendingCheckStillWorksWithinDeltaC(t *testing.T) {
	contract := &mockMedianContract{
		configDigest:    testConfigDigest,
		epoch:           100,
		round:           5,
		latestAnswer:    big.NewInt(100),
		latestTimestamp: time.Now().Add(-60 * time.Second),
	}

	// Large DeltaC so the pending check stays active for the test
	plugin := newTestPlugin(contract, 10*time.Second)
	ctx := context.Background()
	report := buildReport(t, big.NewInt(100))

	// Accept first report
	repts1 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 101, Round: 1}
	accepted, err := plugin.ShouldAcceptFinalizedReport(ctx, repts1, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !accepted {
		t.Fatal("expected first report to be accepted")
	}

	// Second report (same value, within DeltaC) -- should be rejected
	repts2 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 102, Round: 1}
	accepted, err = plugin.ShouldAcceptFinalizedReport(ctx, repts2, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accepted {
		t.Fatal("expected duplicate report to be rejected within DeltaC window")
	}
	t.Log("Pending check still active: duplicate report correctly rejected within DeltaC")
}

// TestShouldAcceptFinalizedReport_VolatileFeedCanRecover verifies that volatile
// feeds (with price deviation) can still bypass the pending check immediately
// without waiting for DeltaC.
func TestShouldAcceptFinalizedReport_VolatileFeedCanRecover(t *testing.T) {
	contract := &mockMedianContract{
		configDigest:    testConfigDigest,
		epoch:           100,
		round:           5,
		latestAnswer:    big.NewInt(100),
		latestTimestamp: time.Now().Add(-60 * time.Second),
	}

	plugin := newTestPlugin(contract, 10*time.Second)
	plugin.offchainConfig.AlphaAcceptInfinite = false
	plugin.offchainConfig.AlphaAcceptPPB = 10_000_000 // 1%
	ctx := context.Background()

	// Accept initial report at value=100
	report1 := buildReport(t, big.NewInt(100))
	repts1 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 101, Round: 1}
	accepted, err := plugin.ShouldAcceptFinalizedReport(ctx, repts1, report1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !accepted {
		t.Fatal("expected first report to be accepted")
	}

	// TX fails -- contract stays at {100, 5}

	// Same-value report rejected (no deviation, pending active)
	report2 := buildReport(t, big.NewInt(100))
	repts2 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 105, Round: 1}
	accepted, err = plugin.ShouldAcceptFinalizedReport(ctx, repts2, report2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accepted {
		t.Fatal("expected same-value report to be rejected")
	}

	// 5% deviated report accepted immediately
	report3 := buildReport(t, big.NewInt(105))
	repts3 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 106, Round: 1}
	accepted, err = plugin.ShouldAcceptFinalizedReport(ctx, repts3, report3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !accepted {
		t.Fatal("expected deviated report to be accepted immediately")
	}
	t.Log("Volatile feed recovery via deviation still works")
}

// TestShouldAcceptFinalizedReport_NormalFlow verifies the happy path where TX
// lands on-chain and subsequent reports are accepted normally.
func TestShouldAcceptFinalizedReport_NormalFlow(t *testing.T) {
	contract := &mockMedianContract{
		configDigest:    testConfigDigest,
		epoch:           100,
		round:           5,
		latestAnswer:    big.NewInt(100),
		latestTimestamp: time.Now().Add(-60 * time.Second),
	}

	plugin := newTestPlugin(contract, 180*time.Second)
	ctx := context.Background()
	report := buildReport(t, big.NewInt(100))

	// Accept first report
	repts1 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 101, Round: 1}
	accepted, err := plugin.ShouldAcceptFinalizedReport(ctx, repts1, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !accepted {
		t.Fatal("expected first report to be accepted")
	}

	// TX lands on-chain
	contract.epoch = 101
	contract.round = 1

	// Next report accepted (nothingPending = true)
	repts2 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 105, Round: 1}
	accepted, err = plugin.ShouldAcceptFinalizedReport(ctx, repts2, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !accepted {
		t.Fatal("expected second report to be accepted after TX landed on-chain")
	}
	t.Log("Normal flow: on-chain confirmation unblocks next report")
}

// TestShouldAcceptFinalizedReport_StaleReportRejected verifies that genuinely
// stale reports (epoch/round <= latestAccepted) are still rejected.
func TestShouldAcceptFinalizedReport_StaleReportRejected(t *testing.T) {
	contract := &mockMedianContract{
		configDigest:    testConfigDigest,
		epoch:           100,
		round:           5,
		latestAnswer:    big.NewInt(100),
		latestTimestamp: time.Now().Add(-60 * time.Second),
	}

	plugin := newTestPlugin(contract, 180*time.Second)
	ctx := context.Background()
	report := buildReport(t, big.NewInt(100))

	// Accept report at epoch 101
	repts1 := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 101, Round: 1}
	accepted, err := plugin.ShouldAcceptFinalizedReport(ctx, repts1, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !accepted {
		t.Fatal("expected first report to be accepted")
	}

	// Try to accept a stale report (epoch 100, before what we already accepted)
	reptsStale := types.ReportTimestamp{ConfigDigest: testConfigDigest, Epoch: 100, Round: 1}
	accepted, err = plugin.ShouldAcceptFinalizedReport(ctx, reptsStale, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accepted {
		t.Fatal("expected stale report to be rejected")
	}
	t.Log("Stale report correctly rejected")
}
