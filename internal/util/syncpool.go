package util

import "sync"

// SyncPool is a type-safe wrapper around sync.SyncPool using generics.
type SyncPool[T any] struct {
	pool sync.Pool
}

// NewSyncPool creates a new type-safe pool.
// If newFunc is non-nil, it is used to generate new values when Get is called on an empty pool.
// If newFunc is nil, Get will return the zero value of T when the pool is empty.
func NewSyncPool[T any](newFunc func() T) *SyncPool[T] {
	var newF func() any
	if newFunc != nil {
		newF = func() any {
			return newFunc()
		}
	}
	return &SyncPool[T]{
		sync.Pool{
			New: newF,
		},
	}
}

// Get returns a value from the pool. If the pool is empty and no New function
// was provided, it returns the zero value of T.
func (p *SyncPool[T]) Get() T {
	v := p.pool.Get()
	if v == nil {
		var zero T
		return zero
	}
	return v.(T)
}

func (p *SyncPool[T]) Put(x T) {
	p.pool.Put(x)
}

func (p *SyncPool[T]) WithPoolItem(f func(T)) {
	item := p.Get()
	defer p.Put(item)
	f(item)
}
