package metricshelper

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
)

// LazyGauge registers its underlying Gauge the first time Set is called. This
// avoids exposing a default zero value before we have an appropriate value to
// expose.
type LazyGauge struct {
	// in
	registerer prometheus.Registerer
	logger     commontypes.Logger
	name       string

	// local
	gauge prometheus.Gauge

	mu             sync.Mutex
	shouldRegister bool
}

func NewLazyGauge(
	registerer prometheus.Registerer,
	logger commontypes.Logger,
	opts prometheus.GaugeOpts,
) *LazyGauge {
	return &LazyGauge{
		registerer,
		logger,
		opts.Name,
		prometheus.NewGauge(opts),

		sync.Mutex{},
		true,
	}
}

func (g *LazyGauge) Set(value float64) {
	g.gauge.Set(value)

	var shouldRegister bool
	g.mu.Lock()
	if g.shouldRegister {
		shouldRegister = true
		g.shouldRegister = false
	}
	g.mu.Unlock()

	if shouldRegister {
		RegisterOrLogError(g.logger, g.registerer, g.gauge, g.name)
		g.gauge.Set(value) // To ensure the value is detected on collection, even if collection only sees values set after registration
	}
}

// Unregister must be used instead of [prometheus.Registerer.Unregister].
// Otherwise, it would be possible for the LazyGauge to have never been Set and
// thus registered, and then after invocation of
// [prometheus.Registerer.Unregister], the next Set would cause the gauge to be
// registered, which is almost definitely unexpected behavior. In fact,
// LazyGauge intentionally does not implement prometheus.Collector to avoid
// confusion. The meaning of the return value matches that of
// [prometheus.Registerer.Unregister].
func (g *LazyGauge) Unregister() bool {
	g.mu.Lock()
	g.shouldRegister = false
	g.mu.Unlock()

	return g.registerer.Unregister(g.gauge)
}
