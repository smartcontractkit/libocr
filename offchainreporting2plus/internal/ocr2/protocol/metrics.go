package protocol

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
)

type pacemakerMetrics struct {
	registerer prometheus.Registerer
	epoch      prometheus.Gauge
	leader     prometheus.Gauge
}

func newPacemakerMetrics(registerer prometheus.Registerer,
	oracleID commontypes.OracleID,
	logger commontypes.Logger) pacemakerMetrics {
	newEpochsNum := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_epoch",
		Help: "The total number of initialized epochs",
		ConstLabels: prometheus.Labels{
			"oracleID": fmt.Sprintf("%d", oracleID),
		},
	})

	metricshelper.RegisterOrLogError(logger, registerer, newEpochsNum, "ocr2_epoch")

	leader := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_experimental_leader_oid",
		Help: "The leader oracle id",
		ConstLabels: prometheus.Labels{
			"oracleID": fmt.Sprintf("%d", oracleID),
		},
	})

	metricshelper.RegisterOrLogError(logger, registerer, leader, "ocr2_experimental_leader_oid")

	return pacemakerMetrics{
		registerer,
		newEpochsNum,
		leader,
	}
}

func (pm *pacemakerMetrics) Close() {
	pm.registerer.Unregister(pm.epoch)
}
