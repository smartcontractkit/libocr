package protocol

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
)

type pacemakerMetrics struct {
	registerer                  prometheus.Registerer
	epoch                       prometheus.Gauge
	leader                      prometheus.Gauge
	receivedNewEpochWishesTotal []prometheus.Counter
}

func newPacemakerMetrics(registerer prometheus.Registerer,
	logger commontypes.Logger, n int) *pacemakerMetrics {

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

	receivedNewEpochWishesTotal := make([]prometheus.Counter, n)
	for i := range n {
		receivedNewEpochWishesTotal[i] = prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ocr3_received_new_epoch_wishes_total",
			Help: "The total number of NewEpochWish messages received from the sender oracle. Note that " +
				"oracles periodically resend their latest NewEpochWish every DeltaResend. This metric is " +
				"useful for checking whether the sender oracle is up and participating in the protocol.",
			ConstLabels: prometheus.Labels{"sender_oracle_id": fmt.Sprint(i)},
		})
		metricshelper.RegisterOrLogError(logger, registerer, receivedNewEpochWishesTotal[i], "ocr3_received_new_epoch_wishes_total")
	}

	return &pacemakerMetrics{
		registerer,
		epoch,
		leader,
		receivedNewEpochWishesTotal,
	}
}

func (pm *pacemakerMetrics) Close() {
	pm.registerer.Unregister(pm.epoch)
	pm.registerer.Unregister(pm.leader)
	for _, c := range pm.receivedNewEpochWishesTotal {
		pm.registerer.Unregister(c)
	}
}

type outcomeGenerationMetrics struct {
	registerer                     prometheus.Registerer
	committedSeqNr                 prometheus.Gauge
	sentObservationsTotal          prometheus.Counter
	includedObservationsTotal      prometheus.Counter
	ledCommittedRoundsTotal        prometheus.Counter
	theirIncludedObservationsTotal []prometheus.Counter
}

func newOutcomeGenerationMetrics(registerer prometheus.Registerer,
	logger commontypes.Logger, n int) *outcomeGenerationMetrics {

	committedSeqNr := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr3_committed_sequence_number",
		Help: "The committed sequence number",
	})
	metricshelper.RegisterOrLogError(logger, registerer, committedSeqNr, "ocr3_committed_sequence_number")

	sentObservationsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_sent_observations_total",
		Help: "The total number of observations by this oracle sent to the leader. Note that a " +
			"sent observation might not arrive at the leader in time, or not be included in a " +
			"proposal for other reasons. This metric is useful for checking an oracle's ability " +
			"to make observations.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, sentObservationsTotal, "ocr3_sent_observations_total")

	includedObservationsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_included_observations_total",
		Help: "The total number of (valid) observations by this oracle included in a proposal " +
			"from the leader. Note that there is no guarantee that the proposal will actually get " +
			"committed; for instance, because the leader crashes or maliciously equivocates to " +
			"make this oracle believe that an observation was included. This metric is useful in " +
			"comparison with ocr3_sent_observations_total to check whether an oracle is able to " +
			"regularly make observations that are included in proposals.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, includedObservationsTotal, "ocr3_included_observations_total")

	ledCommittedRoundsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_led_committed_rounds_total",
		Help: "The total number of rounds committed that were led by this oracle. This metric is " +
			"useful for checking an oracle's ability to act as leader.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, ledCommittedRoundsTotal, "ocr3_led_committed_rounds_total")

	theirIncludedObservationsTotal := make([]prometheus.Counter, n)
	for i := range n {
		theirIncludedObservationsTotal[i] = prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ocr3_their_included_observations_total",
			Help: "The total number of (valid) observations by the observer oracle included in a proposal " +
				"from the leader. Note that there is no guarantee that the proposal will actually get " +
				"committed. This metric is useful for checking whether the observer oracle is able to " +
				"regularly make observations that are included in proposals.",
			ConstLabels: prometheus.Labels{"observer_oracle_id": fmt.Sprint(i)},
		})
		metricshelper.RegisterOrLogError(logger, registerer, theirIncludedObservationsTotal[i], "ocr3_their_included_observations_total")
	}

	return &outcomeGenerationMetrics{
		registerer,
		committedSeqNr,
		sentObservationsTotal,
		includedObservationsTotal,
		ledCommittedRoundsTotal,
		theirIncludedObservationsTotal,
	}
}

func (om *outcomeGenerationMetrics) Close() {
	om.registerer.Unregister(om.committedSeqNr)
	om.registerer.Unregister(om.sentObservationsTotal)
	om.registerer.Unregister(om.includedObservationsTotal)
	om.registerer.Unregister(om.ledCommittedRoundsTotal)
	for _, c := range om.theirIncludedObservationsTotal {
		om.registerer.Unregister(c)
	}
}
