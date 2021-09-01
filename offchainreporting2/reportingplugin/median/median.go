package median

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/golang/protobuf/proto"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/observationhelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/subprocesses"
	"go.uber.org/multierr"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

var i = big.NewInt

// Bounds on an ethereum int192
const byteWidth = 24
const bitWidth = byteWidth * 8

var MaxObservation = i(0).Sub(i(0).Lsh(i(1), bitWidth-1), i(1)) // 2**191 - 1
var MinObservation = i(0).Sub(i(0).Neg(MaxObservation), i(1))   // -2**191

var reportTypes = getReportTypes()

func getReportTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "rawObservers", Type: mustNewType("bytes32")},
		{Name: "observations", Type: mustNewType("int192[]")},
	})
}

var _ types.ReportingPlugin = (*numericalMedian)(nil)

type Config struct {
	AlphaPPB uint64
	DeltaC   time.Duration
}

func DecodeConfig(b []byte) (Config, error) {
	if len(b) != 16 {
		return Config{}, fmt.Errorf("pluggable logic configuration has wrong length, expected %v got %v", 16, len(b))
	}

	alphaPPB := binary.BigEndian.Uint64(b[0:8])
	deltaCUint := binary.BigEndian.Uint64(b[8:16])
	deltaC := time.Duration(int64(deltaCUint))
	return Config{
		alphaPPB,
		deltaC,
	}, nil
}

func (c Config) Encode() []byte {
	pluginConfig := [16]byte{}
	binary.BigEndian.PutUint64(pluginConfig[0:8], c.AlphaPPB)
	binary.BigEndian.PutUint64(pluginConfig[8:16], uint64(c.DeltaC))
	return pluginConfig[:]
}

type MedianContract interface {
	LatestTransmissionDetails(
		ctx context.Context,
	) (
		configDigest types.ConfigDigest,
		epoch uint32,
		round uint8,
		latestAnswer *big.Int,
		latestTimestamp time.Time,
		err error,
	)

	// LatestRoundRequested returns the configDigest, epoch, and round from the latest
	// RoundRequested event emitted by the contract. LatestRoundRequested may or may not
	// return a result if the latest such event was emitted in a block b such that
	// b.timestamp < tip.timestamp - lookback.
	//
	// If no event is found, LatestRoundRequested should return zero values, not an error.
	// An error should only be returned if an actual error occurred during execution,
	// e.g. because there was an error querying the blockchain or the database.
	//
	// As an optimization, this function may also return zero values, if no
	// RoundRequested event has been emitted after the latest NewTransmission event.
	LatestRoundRequested(
		ctx context.Context,
		lookback time.Duration,
	) (
		configDigest types.ConfigDigest,
		epoch uint32,
		round uint8,
		err error,
	)
}

// DataSource implementations must be thread-safe. Observe may be called by many different threads concurrently.
type DataSource interface {
	// Observe queries the data source. Returns a value or an error. Once the
	// context is expires, Observe may still do cheap computations and return a
	// result, but should return as quickly as possible.
	//
	// More details: In the current implementation, the context passed to
	// Observe will time out after LocalConfig.DataSourceTimeout. However,
	// Observe should *not* make any assumptions about context timeout behavior.
	// Once the context times out, Observe should prioritize returning as
	// quickly as possible, but may still perform fast computations to return a
	// result rather than errror. For example, if Observe medianizes a number
	// of data sources, some of which already returned a result to Observe prior
	// to the context's expiry, Observe might still compute their median, and
	// return it instead of an error.
	//
	// Important: Observe should not perform any potentially time-consuming
	// actions like database access, once the context passed has expired.
	Observe(context.Context) (*big.Int, error)
}

var _ types.ReportingPluginFactory = NumericalMedianFactory{}

type NumericalMedianFactory struct {
	ContractTransmitter MedianContract
	DataSource          DataSource
	Logger              commontypes.Logger
}

func (fac NumericalMedianFactory) NewReportingPlugin(configuration types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {

	c, err := DecodeConfig(configuration.OffchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	logger := loghelper.MakeRootLoggerWithContext(fac.Logger).MakeChild(commontypes.LogFields{
		"configDigest":    configuration.ConfigDigest,
		"reportingPlugin": "NumericalMedian",
	})

	return &numericalMedian{
			c,
			fac.ContractTransmitter,
			fac.DataSource,
			logger,

			configuration.ConfigDigest,
			epochRound{},
			new(big.Int),
		}, types.ReportingPluginInfo{
			"NumericalMedian",
			false,
			0,
			4 /* timestamp */ + byteWidth /* observation */ + 8, /* overhead */
			32 /* timestamp */ + 32 /* rawObservers */ + (2*32 + configuration.N*32 /*observations*/) + 32, /* overhead */
		}, nil
}

func deviates(thresholdPPB uint64, old *big.Int, new *big.Int) bool {
	if old.Cmp(i(0)) == 0 {
		if new.Cmp(i(0)) == 0 {
			return false // Both values are zero; no deviation
		}
		return true // Any deviation from 0 is significant
	}
	// ||new - old|| / ||old||, approximated by a float
	change := &big.Rat{}
	change.SetFrac(i(0).Sub(new, old), old)
	change.Abs(change)
	threshold := &big.Rat{}
	threshold.SetFrac(
		(&big.Int{}).SetUint64(thresholdPPB),
		(&big.Int{}).SetUint64(1000000000),
	)
	return change.Cmp(threshold) > 0
}

var _ types.ReportingPlugin = (*numericalMedian)(nil)

type numericalMedian struct {
	Config              Config
	ContractTransmitter MedianContract
	DataSource          DataSource
	Logger              loghelper.LoggerWithContext

	configDigest             types.ConfigDigest
	latestAcceptedEpochRound epochRound
	latestAcceptedMedian     *big.Int
}

func (nm *numericalMedian) Query(ctx context.Context, repts types.ReportTimestamp) (types.Query, error) {
	return nil, nil
}

func (nm *numericalMedian) Observation(ctx context.Context, repts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	if len(query) != 0 {
		return nil, fmt.Errorf("expected empty query")
	}
	value, err := nm.DataSource.Observe(ctx)
	if err != nil {
		return nil, fmt.Errorf("error during DataSource.Observe: %w", err)
	}
	if value == nil {
		return nil, fmt.Errorf("DataSource.Observe returned nil big.Int which should never happen")
	}
	return proto.Marshal(&ObservationProto{
		// zero-initialize protobuf built-ins
		protoimpl.MessageState{},
		0,
		nil,
		// fields
		uint32(time.Now().Unix()),
		observationhelper.ToBytes(value),
	})
}

type parsedAttributedObservation struct {
	Timestamp uint32
	Value     *big.Int
	Observer  commontypes.OracleID
}

func parseAttributedObservations(aos []types.AttributedObservation, skipBad bool) ([]parsedAttributedObservation, error) {
	paos := make([]parsedAttributedObservation, 0, len(aos))
	for i, ao := range aos {
		var observationProto ObservationProto
		if err := proto.Unmarshal(ao.Observation, &observationProto); err != nil {
			if skipBad {
				continue
			} else {
				return nil, fmt.Errorf("Error while unmarshalling %v-th attributed observation: %w", i, err)
			}
		}
		value, err := observationhelper.ToBigInt(observationProto.Value)
		if err != nil {
			if skipBad {
				continue
			} else {
				return nil, fmt.Errorf("Error while converting %v-th attributed observation to big.Int: %w", i, err)
			}
		}
		paos = append(paos, parsedAttributedObservation{
			observationProto.Timestamp,
			value,
			ao.Observer,
		})
	}
	return paos, nil
}

func (nm *numericalMedian) Report(ctx context.Context, repts types.ReportTimestamp, query types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	if len(query) != 0 {
		return false, nil, fmt.Errorf("expected empty query")
	}

	paos, err := parseAttributedObservations(aos, true)
	if err != nil {
		return false, nil, fmt.Errorf("error in parseAttributedObservations: %w", err)
	}

	should, err := nm.shouldReport(ctx, repts, paos)
	if err != nil {
		return false, nil, err
	}
	if !should {
		return false, nil, nil
	}
	report, err := nm.buildReport(ctx, paos)
	if err != nil {
		return false, nil, err
	}
	return true, report, nil
}

func (nm *numericalMedian) shouldReport(ctx context.Context, repts types.ReportTimestamp, paos []parsedAttributedObservation) (bool, error) {
	if len(paos) == 0 {
		return false, fmt.Errorf("cannot handle empty attributed observations")
	}

	var resultTransmissionDetails struct {
		configDigest    types.ConfigDigest
		epoch           uint32
		round           uint8
		latestAnswer    *big.Int
		latestTimestamp time.Time
		err             error
	}
	var resultRoundRequested struct {
		configDigest types.ConfigDigest
		epoch        uint32
		round        uint8
		err          error
	}

	var subs subprocesses.Subprocesses
	subs.Go(func() {
		resultTransmissionDetails.configDigest,
			resultTransmissionDetails.epoch,
			resultTransmissionDetails.round,
			resultTransmissionDetails.latestAnswer,
			resultTransmissionDetails.latestTimestamp,
			resultTransmissionDetails.err =
			nm.ContractTransmitter.LatestTransmissionDetails(ctx)
	})
	subs.Go(func() {
		resultRoundRequested.configDigest,
			resultRoundRequested.epoch,
			resultRoundRequested.round,
			resultRoundRequested.err =
			nm.ContractTransmitter.LatestRoundRequested(ctx, nm.Config.DeltaC)
	})
	subs.Wait()

	if err := multierr.Combine(resultTransmissionDetails.err, resultRoundRequested.err); err != nil {
		// Err on the side of creating too many reports. For instance, the Ethereum node
		// might be down, but that need not prevent us from still contributing to the
		// protocol.
		return true, fmt.Errorf("error during LatestTransmissionDetails/LatestRoundRequested: %w", err)
	}

	// sort by values
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Value.Cmp(paos[j].Value) < 0
	})

	answer := paos[len(paos)/2].Value

	initialRound := // Is this the first round for this configuration?
		resultTransmissionDetails.configDigest == repts.ConfigDigest &&
			resultTransmissionDetails.epoch == 0 &&
			resultTransmissionDetails.round == 0
	deviation := // Has the result changed enough to merit a new report?
		deviates(nm.Config.AlphaPPB, resultTransmissionDetails.latestAnswer, answer)

	deltaCTimeout := // Has enough time passed since the last report, to merit a new one?
		resultTransmissionDetails.latestTimestamp.Add(nm.Config.DeltaC).
			Before(time.Now())
	unfulfilledRequest := // Has a new report been requested explicitly?
		resultRoundRequested.configDigest == repts.ConfigDigest &&
			!(epochRound{resultRoundRequested.epoch, resultRoundRequested.round}).
				Less(epochRound{resultTransmissionDetails.epoch, resultTransmissionDetails.round})

	logger := nm.Logger.MakeChild(commontypes.LogFields{
		"timestamp":                 repts,
		"initialRound":              initialRound,
		"alphaPPB":                  nm.Config.AlphaPPB,
		"deviation":                 deviation,
		"deltaC":                    nm.Config.DeltaC,
		"deltaCTimeout":             deltaCTimeout,
		"lastTransmissionTimestamp": resultTransmissionDetails.latestTimestamp,
		"unfulfilledRequest":        unfulfilledRequest,
	})

	// The following is more succinctly expressed as a disjunction, but breaking
	// the branches up into their own conditions makes it easier to check that
	// each branch is tested, and also allows for more expressive log messages
	if initialRound {
		logger.Info("shouldReport: yes, because it's the first round of the first epoch", commontypes.LogFields{
			"result": true,
		})
		return true, nil
	}
	if deviation {
		logger.Info("shouldReport: yes, because new median deviates sufficiently from current onchain value", commontypes.LogFields{
			"result": true,
		})
		return true, nil
	}
	if deltaCTimeout {
		logger.Info("shouldReport: yes, because deltaC timeout since last onchain report", commontypes.LogFields{
			"result": true,
		})
		return true, nil
	}
	if unfulfilledRequest {
		logger.Info("shouldReport: yes, because a new report has been explicitly requested", commontypes.LogFields{
			"result": true,
		})
		return true, nil
	}
	logger.Info("shouldReport: no", commontypes.LogFields{"result": false})
	return false, nil
}

func (nm *numericalMedian) buildReport(ctx context.Context, paos []parsedAttributedObservation) (types.Report, error) {
	if len(paos) == 0 {
		return nil, fmt.Errorf("Cannot build report from empty attributed observations")
	}

	// it's okay to mutate aos2 here
	// get median timestamp
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Timestamp < paos[j].Timestamp
	})
	timestamp := paos[len(paos)/2].Timestamp

	// sort by values
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Value.Cmp(paos[j].Value) < 0
	})

	observers := [32]byte{}
	observations := []*big.Int{}

	for i, pao := range paos {
		observers[i] = byte(pao.Observer)
		observations = append(observations, pao.Value)
	}

	reportBytes, err := reportTypes.Pack(timestamp, observers, observations)
	return types.Report(reportBytes), err
}

func (nm *numericalMedian) ShouldAcceptFinalizedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	reportEpochRound := epochRound{repts.Epoch, repts.Round}
	if !nm.latestAcceptedEpochRound.Less(reportEpochRound) {
		nm.Logger.Debug("ShouldAcceptFinalizedReport() = false, report is stale", commontypes.LogFields{
			"latestAcceptedEpochRound": nm.latestAcceptedEpochRound,
			"reportEpochRound":         reportEpochRound,
		})
		return false, nil
	}

	contractConfigDigest, contractEpoch, contractRound, _, _, err := nm.ContractTransmitter.LatestTransmissionDetails(ctx)
	if err != nil {
		return false, err
	}

	contractEpochRound := epochRound{contractEpoch, contractRound}

	if contractConfigDigest != nm.configDigest {
		nm.Logger.Debug("ShouldAcceptFinalizedReport() = false, config digest mismatch", commontypes.LogFields{
			"contractConfigDigest": contractConfigDigest,
			"reportConfigDigest":   nm.configDigest,
			"reportEpochRound":     reportEpochRound,
		})
		return false, nil
	}

	if !contractEpochRound.Less(reportEpochRound) {
		nm.Logger.Debug("ShouldAcceptFinalizedReport() = false, report is stale", commontypes.LogFields{
			"contractEpochRound": contractEpochRound,
			"reportEpochRound":   reportEpochRound,
		})
		return false, nil
	}

	reportMedian, err := medianFromReport(report)
	if err != nil {
		return false, fmt.Errorf("error during medianFromReport: %w", err)
	}

	deviates := deviates(nm.Config.AlphaPPB, nm.latestAcceptedMedian, reportMedian)
	nothingPending := !contractEpochRound.Less(nm.latestAcceptedEpochRound)
	result := deviates || nothingPending

	nm.Logger.Debug("ShouldAcceptFinalizedReport() = result", commontypes.LogFields{
		"contractEpochRound":       contractEpochRound,
		"reportEpochRound":         reportEpochRound,
		"latestAcceptedEpochRound": nm.latestAcceptedEpochRound,
		"deviates":                 deviates,
		"result":                   result,
	})

	if result {
		nm.latestAcceptedEpochRound = reportEpochRound
		nm.latestAcceptedMedian = reportMedian
	}

	return result, nil
}

func medianFromReport(report types.Report) (*big.Int, error) {
	reportElems, err := reportTypes.Unpack(report)
	if err != nil {
		return nil, fmt.Errorf("error during unpack: %w", err)
	}

	if len(reportElems) != 3 {
		return nil, fmt.Errorf("length mismatch, expected 3, got %v", len(reportElems))
	}

	observations, ok := reportElems[2].([]*big.Int)
	if !ok {
		return nil, fmt.Errorf("cannot cast observations to []*big.Int, type is %T", reportElems[1])
	}

	if len(observations) == 0 {
		return nil, fmt.Errorf("observations are empty")
	}

	median := observations[len(observations)/2]
	if median == nil {
		return nil, fmt.Errorf("median is nil")
	}

	return median, nil
}

func (nm *numericalMedian) ShouldTransmitAcceptedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	reportEpochRound := epochRound{repts.Epoch, repts.Round}

	contractConfigDigest, contractEpoch, contractRound, _, _, err := nm.ContractTransmitter.LatestTransmissionDetails(ctx)
	if err != nil {
		return false, err
	}

	contractEpochRound := epochRound{contractEpoch, contractRound}

	if contractConfigDigest != nm.configDigest {
		nm.Logger.Debug("ShouldTransmitAcceptedReport() = false, config digest mismatch", commontypes.LogFields{
			"contractConfigDigest": contractConfigDigest,
			"reportConfigDigest":   nm.configDigest,
			"reportEpochRound":     reportEpochRound,
		})
		return false, nil
	}

	if !contractEpochRound.Less(reportEpochRound) {
		nm.Logger.Debug("ShouldTransmitAcceptedReport() = false, report is stale", commontypes.LogFields{
			"contractEpochRound": contractEpochRound,
			"reportEpochRound":   reportEpochRound,
		})
		return false, nil
	}

	return true, nil
}

func (nm *numericalMedian) Start() error {
	return nil
}

func (nm *numericalMedian) Close() error {
	return nil
}
