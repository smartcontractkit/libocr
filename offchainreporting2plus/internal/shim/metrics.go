package shim

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
)

type serializingEndpointMetrics struct {
	registerer           prometheus.Registerer
	sentMessagesTotal    prometheus.Counter
	droppedMessagesTotal prometheus.Counter
}

func newSerializingEndpointMetrics(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
) *serializingEndpointMetrics {
	sentMessagesTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr_telemetry_sent_messages_total",
		Help: "The number of telemetry messages sent.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, sentMessagesTotal, "ocr_telemetry_sent_messages_total")

	droppedMessagesTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr_telemetry_dropped_messages_total",
		Help: "The number of telemetry messages dropped.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, droppedMessagesTotal, "ocr_telemetry_dropped_messages_total")

	return &serializingEndpointMetrics{
		registerer,
		sentMessagesTotal,
		droppedMessagesTotal,
	}
}

func (m *serializingEndpointMetrics) Close() {
	m.registerer.Unregister(m.sentMessagesTotal)
	m.registerer.Unregister(m.droppedMessagesTotal)
}
