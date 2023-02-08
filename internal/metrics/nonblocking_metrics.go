package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/subprocesses"
)

var _ commontypes.Metrics = (*NonblockingMetricsWrapper)(nil)

type NonblockingMetricsWrapper struct {
	subs          subprocesses.Subprocesses
	ctx           context.Context
	ctxCancel     context.CancelFunc
	logger        commontypes.Logger
	flushInterval time.Duration
	metricsImpl   commontypes.Metrics
	metricIndex   metricIndexMap // metricFingerprint(name, labels) => nonblockingMetricVec
}

// NewNonblockingMetricsWrapper is a non-blocking implementation of  the commontypes.Metrics interface.
// The provided metricsImpl is used to instantiate wrapper vectors of type commontypes.MetricVec with NewMetricVec
// as well as wrapper metrics with GetMerticWhith on the wrapper vectors.
// All method calls on the wrappers are non-blocking, even if the respective calls for the underlying wrapped
// metricsImpl could block.
// To achieve this, NewNonblockingMetricsWrapper starts a flush loop which periodically, according to the provided interval,
// applies collectively the result of all the calls on the wrappers methods to the respective wrapped metricsImpl instances.
// The order of the method calls is reflected on the applied result.
func NewNonblockingMetricsWrapper(logger commontypes.Logger, metricsImpl commontypes.Metrics, flushInterval time.Duration) *NonblockingMetricsWrapper {
	logger.Info("NonblockingMetricsWrapper: New non-blocking metrics wrapper", commontypes.LogFields{"flushInterval": flushInterval.String()})
	ctx, ctxCancel := context.WithCancel(context.Background())
	nbmw := &NonblockingMetricsWrapper{
		subprocesses.Subprocesses{},
		ctx,
		ctxCancel,
		logger,
		flushInterval,
		metricsImpl,
		metricIndexMap{},
	}

	nbmw.subs.Go(nbmw.flushLoop)
	return nbmw
}

// NewMetricVec returns a commontypes.MetricVec without blocking, even if the invocation of NewMetricVec of the
// underlying metrics implementation is blocking. NewMetricVec for the underlying metrics implementation is only invoked
// in the flush loop of the NonblockingMetricsWrapper the first time a value for a metric of the returned metric vector
// is encountered.
// The name and label names may contain ASCII letters, numbers, as well as underscores.
// An error is returned if either the name or some label names is of invalid format.
func (nbmw *NonblockingMetricsWrapper) NewMetricVec(name string, help string, labelNames ...string) (commontypes.MetricVec, error) {
	if err := validName(name); err != nil {
		return nil, fmt.Errorf("could not create metric vector with invalid name: %w", err)
	}

	if err := validLabelNames(labelNames); err != nil {
		return nil, fmt.Errorf("could not create metric vector with invalid label names: %w", err)
	}

	return &nonblockingMetricVec{
		nbmw,
		name,
		help,
		labelNames,
		nil, // initialized in flush loop
	}, nil
}

func validLabelNames(labelNames []string) error {
	for _, name := range labelNames {
		if err := validName(name); err != nil {
			return fmt.Errorf("format of label name %q is not valid: %w", name, err)
		}
	}
	return nil
}

// Close shuts down the flush loop and waits for the go routines of the NonblockingMetricsWrapper to finish
// Close does not close the underlying metrics implementation. It should be the responsibility of the underlying
// metrics implementation to shut down without go routines leakage
func (nbmw *NonblockingMetricsWrapper) Close() error {
	nbmw.logger.Debug("NonblockingMetricsWrapper: Closing non-blocking metrics wrapper", nil)
	nbmw.ctxCancel()
	nbmw.subs.Wait()
	nbmw.logger.Info("NonblockingMetricsWrapper: Closed non-blocking metrics wrapper", nil)
	return nil
}

func (nbmw *NonblockingMetricsWrapper) flushLoop() {
	nbmw.logger.Debug("NonblockingMetricsWrapper: Starting metric reporter flush loop", nil)
	ticker := time.NewTicker(nbmw.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-nbmw.ctx.Done():
			nbmw.logger.Debug("NonblockingMetricsWrapper: Canceling metric reporter flush loop", nil)
			return
		case <-ticker.C:
			nbmw.flush()
		}
	}
}

func (nbmw *NonblockingMetricsWrapper) flush() {
	nbmw.logger.Trace("NonblockingMetricsWrapper: Flushing metricsImpl", nil)
	nbmw.metricIndex.Range(
		func(fingerprint string, nonblockingMetric *nonblockingMetric) bool {
			// The underlying metric implementation might not have been initialized yet
			if nonblockingMetric.metricImpl == nil {
				// Get the parent metric vector
				nonblockingMetricVec := nonblockingMetric.nonblockingMetricVec
				// The underlying metric vector implementation might not have been initialized yet
				if nonblockingMetricVec.metricVecImpl == nil {
					var err error
					nonblockingMetricVec.metricVecImpl, err = nbmw.metricsImpl.NewMetricVec(nonblockingMetricVec.name, nonblockingMetricVec.help, nonblockingMetricVec.labelNames...)
					if err != nil {
						// TODO: This will potentially log the same error in every flush loop iteration, consider deduplicating.
						nbmw.logger.Error("NonblockingMetricsWrapper: Could not create metric vector", commontypes.LogFields{"err": err.Error()})
						return true
					}
				}
				var err error
				nonblockingMetric.metricImpl, err = nonblockingMetricVec.metricVecImpl.GetMetricWith(nonblockingMetric.labels)
				if err != nil {
					// TODO: This will potentially log the same error in every flush loop iteration, consider deduplicating.
					nbmw.logger.Error("NonblockingMetricsWrapper: Could not create metric with labels", commontypes.LogFields{"err": err.Error()})
					return true
				}
			}

			err := nonblockingMetric.Flush()
			if err != nil {
				// TODO: This will potentially log the same error in every flush loop iteration, consider deduplicating.
				nbmw.logger.Error("NonblockingMetricsWrapper: Could not set metric value", commontypes.LogFields{"err": err.Error()})
				return true
			}

			return true
		})
}

type nonblockingMetricVec struct {
	nonblockingMetricsWrapper *NonblockingMetricsWrapper
	name                      string
	help                      string
	labelNames                []string
	metricVecImpl             commontypes.MetricVec
}

// GetMetricWith returns a commontypes.Metric without blocking, even if the invocation of GetMetricWith of the
// underlying metricsImpl is blocking. GetMetricWith for the underlying metrics implementation is only invoked in the
// flush loop of the corresponding NonblockingMetricsWrapper the first time a value for the returned metric is
// encountered.
// An error is returned if the number and names of the labels are
// inconsistent with those of the variable labels of the MetricVec
// or if some labels' value is empty
// or if some labels' value is not a valid UTF-8 string.
// Label names are not validated since they are already validated when the MetricVec is constructed.
func (nbv *nonblockingMetricVec) GetMetricWith(labels map[string]string) (commontypes.Metric, error) {
	if err := nbv.validLabels(labels); err != nil {
		return nil, err
	}

	fingerprint := metricFingerprint(nbv.name, labels)

	nbm := &nonblockingMetric{
		0,
		sync.Mutex{},
		nbv,
		labels,
		nil, // initialized in flush loop
	}

	nbm, _ = nbv.nonblockingMetricsWrapper.metricIndex.LoadOrStore(fingerprint, nbm)

	return nbm, nil
}

func (nbv *nonblockingMetricVec) validLabels(labels map[string]string) error {
	if len(labels) != len(nbv.labelNames) {
		return fmt.Errorf(
			"expected %d label values but got %d in %#v",
			len(nbv.labelNames),
			len(labels), labels,
		)
	}

	for name, val := range labels {
		if len(val) == 0 {
			return fmt.Errorf("label %q: label value is empty", name)
		}
		if !utf8.ValidString(val) {
			return fmt.Errorf("label %q: value %q is not a valid UTF-8 sting", name, val)
		}
	}

	for _, label := range nbv.labelNames {
		if _, ok := labels[label]; !ok {
			return fmt.Errorf("label name %q missing in label map", label)
		}
	}
	return nil
}

type nonblockingMetric struct {
	value     float64
	valueLock sync.Mutex

	nonblockingMetricVec *nonblockingMetricVec
	labels               map[string]string
	metricImpl           commontypes.Metric
}

func (nbm *nonblockingMetric) Set(f float64) {
	nbm.valueLock.Lock()
	defer nbm.valueLock.Unlock()
	nbm.value = f
}

func (nbm *nonblockingMetric) Inc() {
	nbm.Add(1)
}

func (nbm *nonblockingMetric) Dec() {
	nbm.Add(-1)
}

func (nbm *nonblockingMetric) Add(f float64) {
	nbm.valueLock.Lock()
	defer nbm.valueLock.Unlock()
	nbm.value = nbm.value + f
}

func (nbm *nonblockingMetric) Sub(f float64) {
	nbm.Add(f * -1)
}

func (nbm *nonblockingMetric) Flush() error {
	if nbm.metricImpl == nil {
		return fmt.Errorf("metric implementation is not initialized yet")
	}
	nbm.valueLock.Lock()
	value := nbm.value
	nbm.valueLock.Unlock()
	nbm.metricImpl.Set(value)
	return nil
}
