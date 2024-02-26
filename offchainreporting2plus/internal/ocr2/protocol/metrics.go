package protocol

import (
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
	logger commontypes.Logger) pacemakerMetrics {
	epoch := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_epoch",
		Help: "The total number of initialized epochs",
	})
	metricshelper.RegisterOrLogError(logger, registerer, epoch, "ocr2_epoch")

	leader := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_experimental_leader_oid",
		Help: "The leader oracle id",
	})
	metricshelper.RegisterOrLogError(logger, registerer, leader, "ocr2_experimental_leader_oid")

	return pacemakerMetrics{
		registerer,
		epoch,
		leader,
	}
}

func (pm *pacemakerMetrics) Close() {
	pm.registerer.Unregister(pm.epoch)
	pm.registerer.Unregister(pm.leader)
}
