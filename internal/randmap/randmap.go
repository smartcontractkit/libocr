package randmap

import "math/rand/v2"

type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// Like a regular map, but with a GetRandom operation
// that returns a uniformly random entry from the map.
//
// Limitation: The memory consumed by this data structure is evergrowing, it never shrinks.

type Map[K comparable, V any] struct {
	indices map[K]int
	entries []MapEntry[K, V]
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		make(map[K]int),
		make([]MapEntry[K, V], 0),
	}
}

func (m *Map[K, V]) Set(key K, value V) {
	if _, ok := m.indices[key]; !ok {
		m.indices[key] = len(m.entries)
		m.entries = append(m.entries, MapEntry[K, V]{key, value})
	} else {
		m.entries[m.indices[key]] = MapEntry[K, V]{key, value}
	}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	if _, ok := m.indices[key]; !ok {
		var zero V
		return zero, false
	}
	return m.entries[m.indices[key]].Value, true
}

func (m *Map[K, V]) Delete(key K) {
	if _, ok := m.indices[key]; !ok {
		return
	}

	index := m.indices[key]
	delete(m.indices, key)

	if index == len(m.entries)-1 {
		// the element we want to delete is already in last place
		m.entries = m.entries[0 : len(m.entries)-1]
		return
	} else {
		m.entries[index] = m.entries[len(m.entries)-1]
		m.indices[m.entries[index].Key] = index
		m.entries = m.entries[0 : len(m.entries)-1]
	}
}

func (m *Map[K, V]) GetRandom() (MapEntry[K, V], bool) {
	if len(m.entries) == 0 {
		var zero MapEntry[K, V]
		return zero, false
	}
	randomEntry := m.entries[rand.IntN(len(m.entries))]
	return randomEntry, true
}

func (m *Map[K, V]) Size() int {
	return len(m.entries)
}
