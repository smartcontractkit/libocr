package ringbuffer

import "fmt"

// RingBuffer implements a fixed capacity ring buffer for items of type T.
// NOTE: THIS IMPLEMENTATION IS NOT SAFE FOR CONCURRENT USE.
type RingBuffer[T any] struct {
	first int // index of the front (=oldest) element
	size  int // number of elements currently stored in this ring buffer
	items []T // fixed size buffer holding the elements
}

func NewRingBuffer[T any](cap int) *RingBuffer[T] {
	if cap <= 0 {
		panic(fmt.Sprintf("NewRingBuffer: cap must be positive, got %d", cap))
	}
	return &RingBuffer[T]{
		0,
		0,
		make([]T, cap),
	}
}

func (rb *RingBuffer[T]) Size() int {
	return rb.size
}

func (rb *RingBuffer[T]) Cap() int {
	return len(rb.items)
}

func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.size == 0
}

func (rb *RingBuffer[T]) IsFull() bool {
	return rb.size == len(rb.items)
}

// Peek returns the front (=oldest) item without removing it.
// Return false as second argument if there are no items in the ring buffer.
func (rb *RingBuffer[T]) Peek() (result T, ok bool) {
	if rb.size > 0 {
		ok = true
		result = rb.items[rb.first]
	}
	return result, ok
}

// Pop removes and returns the front (=oldest) item.
// Return false as second argument if there are no items in the ring buffer.
func (rb *RingBuffer[T]) Pop() (result T, ok bool) {
	result, ok = rb.Peek()
	if ok {
		var zero T
		rb.items[rb.first] = zero
		rb.first = (rb.first + 1) % len(rb.items)
		rb.size--
	}
	return result, ok
}

// Try to push a new item to the back of the ring buffer.
// Returns
//   - true if the item was added, or
//   - false if the item cannot be added because the buffer is currently full.
func (rb *RingBuffer[T]) TryPush(item T) (ok bool) {
	if rb.IsFull() {
		return false
	}
	rb.items[(rb.first+rb.size)%len(rb.items)] = item
	rb.size++
	return true
}

// Push new item to the back of the ring buffer.
// If the buffer is currently full, the front (=oldest) item is evicted and returned to make space for the new item.
func (rb *RingBuffer[T]) PushEvict(item T) (evicted T, didEvict bool) {
	if rb.IsFull() {
		// Evict the oldest item to be returned.
		evicted = rb.items[rb.first]
		didEvict = true

		// Push the new item to new empty space and update the first index to the next (oldest) item.
		rb.items[rb.first] = item
		rb.first = (rb.first + 1) % len(rb.items)
	} else {
		// Perform a normal push operation (which is known to be successful as the buffer is not full).
		rb.items[(rb.first+rb.size)%len(rb.items)] = item
		rb.size++
	}
	return evicted, didEvict
}
