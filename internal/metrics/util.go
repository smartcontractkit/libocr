package metrics

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
)

// separatorByte is a byte that cannot occur in valid UTF-8 sequences
const separatorByte byte = 255

// Making metricIndexMap thread safe, since it might be accessed concurrently
// Using sync.Map instead of a map with a mutex to avoid lock contention
// Importantly, iterating over a sync.Map with Range does not block any method on the receiver
type metricIndexMap struct {
	mm sync.Map
}

func (m *metricIndexMap) Load(key string) (*nonblockingMetric, bool) {
	if v, ok := m.mm.Load(key); ok {
		if v == nil {
			return nil, ok
		}
		return v.(*nonblockingMetric), ok
	}
	return nil, false
}

func (m *metricIndexMap) Store(key string, value *nonblockingMetric) {
	m.mm.Store(key, value)
}

func (m *metricIndexMap) Delete(key string) {
	m.mm.Delete(key)
}

func (m *metricIndexMap) LoadAndDelete(key string) (*nonblockingMetric, bool) {
	if v, ok := m.mm.LoadAndDelete(key); ok {
		if v == nil {
			return nil, ok
		}
		return v.(*nonblockingMetric), ok
	}
	return nil, false
}

func (m *metricIndexMap) LoadOrStore(key string, value *nonblockingMetric) (*nonblockingMetric, bool) {
	v, ok := m.mm.LoadOrStore(key, value)
	if v == nil {
		return nil, ok
	}
	return v.(*nonblockingMetric), ok
}

func (m *metricIndexMap) Range(f func(key string, value *nonblockingMetric) bool) {
	m.mm.Range(func(k, v any) bool {
		if v == nil {
			return f(k.(string), nil)
		}
		return f(k.(string), v.(*nonblockingMetric))
	})
}

// metricFingerprint returns a unique signature for a given name and label set
func metricFingerprint(name string, labels map[string]string) string {
	var b bytes.Buffer
	b.WriteString(name)
	b.WriteByte(separatorByte)

	labelNames := make([]string, 0, len(labels))
	for labelName := range labels {
		labelNames = append(labelNames, labelName)
	}
	sort.Strings(labelNames)

	for _, labelName := range labelNames {
		b.WriteString(labelName)
		b.WriteByte(separatorByte)
		b.WriteString(labels[labelName])
		b.WriteByte(separatorByte)
	}
	return b.String()
}

// validName checks that the name only contains ASCII letters and digits, as well as underscores and colons
func validName(name string) error {
	if len(name) < 1 {
		return fmt.Errorf("empty string is not allowed")
	}

	for _, c := range name {
		if !('0' <= c && c <= '9') && !('a' <= c && c <= 'z') && !('A' <= c && c <= 'Z') && c != ':' && c != '_' {
			return fmt.Errorf("the string contains an invalid character; it should only contain only ASCII letters and digits, underscores and colons")
		}
	}

	return nil
}
