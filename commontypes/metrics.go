package commontypes

type Metrics interface {
	// NewMetricVec creates a new MetricVec with the provided name and help
	// and partitioned by the given label names
	// The name and label names may contain ASCII letters, numbers, as well as underscores.
	// An error is returned if either the name or some label names is of invalid format.
	NewMetricVec(name string, help string, labelNames ...string) (MetricVec, error)
}

// MetricVec must be thread safe
type MetricVec interface {
	// GetMetricWith returns a Metric for the given MetricVec with the given labels map.
	// If that label map is assessed for the first time, a new Metric is created.
	// Label values may contain any Unicode character (UTF-8 encoded).
	// An error is returned if the number and names of the labels are
	// inconsistent with those of the variable labels of the MetricVec
	// or if some labels' value is empty
	// or if some labels' value is not a valid UTF-8 string.
	// Label names are not validated since they are already validated when the MetricVec is constructed.
	// The implementation of this method must guarantee that different invocations
	// on the same MetricVec with the same label map returns the same Metric.
	// I.e. for
	// m1, err1 := metricVec.MetricWithLabels(L1)
	// m2, err2 := metricVec.MetricWithLabels(L2)
	// s.t. reflect.DeepEqual(L1, L2)
	// then (err1 == nil && err2 == nil) implies m1 == m2
	GetMetricWith(labels map[string]string) (Metric, error)
}

// Metric must be thread safe
type Metric interface {
	// Set sets the Metric to an arbitrary value.
	Set(float64)
	// Inc increments the Metric by 1. Use Add to increment it by arbitrary
	// values.
	Inc()
	// Dec decrements the Metric by 1. Use Sub to decrement it by arbitrary
	// values.
	Dec()
	// Add adds the given value to the Metric. (The value can be negative,
	// resulting in a decrease of the Metric.)
	Add(float64)
	// Sub subtracts the given value from the Metric. (The value can be
	// negative, resulting in an increase of the Metric.)
	Sub(float64)
}
