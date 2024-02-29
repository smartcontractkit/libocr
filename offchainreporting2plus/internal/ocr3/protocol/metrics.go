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
		Name: "ocr3_epoch",
		Help: "The total number of initialized epochs",
	})
	metricshelper.RegisterOrLogError(logger, registerer, epoch, "ocr3_epoch")

	leader := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_experimental_leader_oid",
		Help: "The leader oracle id",
	})
	metricshelper.RegisterOrLogError(logger, registerer, leader, "ocr3_experimental_leader_oid")

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

type outcomeGenerationMetrics struct {
	registerer     prometheus.Registerer
	committedSeqNr prometheus.Gauge
}

func newOutcomeGenerationMetrics(registerer prometheus.Registerer,
	logger commontypes.Logger) outcomeGenerationMetrics {

	committedSeqNr := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_committed_sequence_number",
		Help: "The total number of committed sequence numbers",
	})
	metricshelper.RegisterOrLogError(logger, registerer, committedSeqNr, "ocr3_committed_sequence_number")

	return outcomeGenerationMetrics{
		registerer,
		committedSeqNr,
	}
}

func (om *outcomeGenerationMetrics) Close() {
	om.registerer.Unregister(om.committedSeqNr)
}
