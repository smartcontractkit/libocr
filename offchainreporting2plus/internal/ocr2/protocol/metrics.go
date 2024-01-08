package protocol

import (
	"fmt"
	"github.com/smartcontractkit/libocr/commontypes"
)

type noopMetric struct{}

func (t noopMetric) Set(f float64) {}

func (t noopMetric) Inc() {}

func (t noopMetric) Dec() {}

func (t noopMetric) Add(f float64) {}

func (t noopMetric) Sub(f float64) {}

type pacemakerMetrics struct {
	newEpochsCount commontypes.Metric
}

func noopPacemakerMetrics() pacemakerMetrics {
	return pacemakerMetrics{
		noopMetric{},
	}
}

func newPacemakerMetrics(metrics commontypes.Metrics, oracleID commontypes.OracleID) (pacemakerMetrics, error) {
	newEpochsCountVec, err := metrics.NewMetricVec("ocr2_new_epochs_count", "The total number of initialized epochs", "oracleID")
	if err != nil {
		return pacemakerMetrics{}, err
	}
	newEpochsCount, err := newEpochsCountVec.GetMetricWith(map[string]string{"oracleID": fmt.Sprintf("%d", oracleID)})
	if err != nil {
		return pacemakerMetrics{}, err
	}
	return pacemakerMetrics{
		newEpochsCount,
	}, nil
}

func newPacemakerMetricsOrNoop(metrics commontypes.Metrics, oracleID commontypes.OracleID, logger commontypes.Logger) pacemakerMetrics {
	pm, err := newPacemakerMetrics(metrics, oracleID)
	if err != nil {
		logger.Warn("Pacemaker: failed to instantiate pacemaker metrics", commontypes.LogFields{"error": err})
		return noopPacemakerMetrics()
	}
	return pm
}
