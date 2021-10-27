package tracing

import (
	"sync"
)

// Hook is a function that gets called every time a specific event is triggered.
// Hooks should be small, quick, synchronous and thread-safe. They should FIRST filter
// out all the events they are not interested in, then do their job.
type Hook func(t Trace)

type hooks struct {
	mu        sync.Mutex
	callbacks []Hook
}

func newHooks() *hooks {
	return &hooks{
		sync.Mutex{},
		[]Hook{},
	}
}

func (hs *hooks) register(newHook Hook) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.callbacks = append(hs.callbacks, newHook)
}

func (hs *hooks) execute(trace Trace) {
	for _, cb := range hs.callbacks {
		cb(trace)
	}
}

// Helpers

// TriggerUntilTrue will keep executing the fn hook as long as it returns false.
// When it returns true, it will stop executing fn.
func TriggerUntilTrue(fn func(t Trace) bool) Hook {
	var done bool = false
	var doneMu sync.Mutex
	return func(t Trace) {
		doneMu.Lock()
		defer doneMu.Unlock()
		if done {
			return
		}
		done = done || fn(t)
	}
}
