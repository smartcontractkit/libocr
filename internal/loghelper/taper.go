package loghelper

import "sync"

// LogarithmicTaper provides logarithmic tapering of an event sequence.
// For example, if the taper is Triggered 50 times with a function that
// simply prints the provided count, the output would be 1,2,4,8,16,32.
type LogarithmicTaper struct {
	count   uint64
	countMu sync.Mutex
}

// Trigger increments a count and calls f iff the new count is a power of two
func (tap *LogarithmicTaper) Trigger(f func(newCount uint64)) {
	tap.countMu.Lock()
	tap.count++
	newCount := tap.count
	tap.countMu.Unlock()
	if f != nil && isPowerOfTwo(newCount) {
		f(newCount)
	}
}

// Count returns the internal count of the taper
func (tap *LogarithmicTaper) Count() uint64 {
	tap.countMu.Lock()
	defer tap.countMu.Unlock()
	return tap.count
}

// Reset resets the count to 0 and then calls f with the previous count
// iff it wasn't already 0
func (tap *LogarithmicTaper) Reset(f func(oldCount uint64)) {
	tap.countMu.Lock()
	oldCount := tap.count
	tap.count = 0
	tap.countMu.Unlock()
	if oldCount != 0 {
		f(oldCount)
	}
}

func isPowerOfTwo(num uint64) bool {
	return num != 0 && (num&(num-1)) == 0
}
