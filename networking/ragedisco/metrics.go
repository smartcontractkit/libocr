package ragedisco

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
)

type discoveryProtocolMetrics struct {
	registerer      prometheus.Registerer
	peers           prometheus.Gauge
	peersDiscovered prometheus.Gauge
	bootstrappers   prometheus.Gauge
}

func newDiscoveryProtocolMetrics(registerer prometheus.Registerer, peerID string, logger commontypes.Logger) discoveryProtocolMetrics {
	labels := map[string]string{"peerID": peerID}

	peers := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ragedisco_peers",
		Help:        "The total number of peers in network discovery",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, peers, "ragedisco_peers")

	peersDiscovered := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ragedisco_peers_discovered",
		Help:        "The number of discovered peers in network discovery",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, peersDiscovered, "ragedisco_peers_discovered")

	bootstappers := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ragedisco_bootstappers",
		Help:        "The number of undiscovered peers in network discover",
		ConstLabels: labels,
	})

	metricshelper.RegisterOrLogError(logger, registerer, bootstappers, "ragedisco_bootstappers")

	return discoveryProtocolMetrics{
		registerer,
		peers,
		peersDiscovered,
		bootstappers,
	}
}

func (dpm *discoveryProtocolMetrics) Close() {
	dpm.registerer.Unregister(dpm.peers)
	dpm.registerer.Unregister(dpm.bootstrappers)
	dpm.registerer.Unregister(dpm.peersDiscovered)
}
