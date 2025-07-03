package queue

type Queue[T any] struct {
	elements    []T
	maxCapacity *int
}

// NewQueue returns a queue with infinite maxCapacity
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		elements: make([]T, 0),
	}
}

// NewQueueWithMaxCapacity returns queue with maxCapacity cap.
// If the maxCapacity is reached the queue does not accept more elements.
func NewQueueWithMaxCapacity[T any](cap int) *Queue[T] {
	return &Queue[T]{
		elements:    make([]T, 0),
		maxCapacity: &cap,
	}
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.elements) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.elements)
}

// Push returns false if the queue is at maxCapacity and the element is not added
func (q *Queue[T]) Push(element T) bool {
	if q.maxCapacity == nil || len(q.elements) < *q.maxCapacity {
		q.elements = append(q.elements, element)
		return true
	}
	return false
}

// Peek returns the first element without removing it. It returns false if the queue is empty.
func (q *Queue[T]) Peek() (*T, bool) {
	if len(q.elements) == 0 {
		return nil, false
	}
	return &q.elements[0], true
}

// Pop returns the first element after removing it. It returns false if the queue is empty.
func (q *Queue[T]) Pop() (T, bool) {
	if len(q.elements) == 0 {
		var zero T
		return zero, false
	}
	first := q.elements[0]

	q.elements = q.elements[1:len(q.elements)]
	return first, true
}

// PeekLast returns the last element without removing it. It returns false if the queue is empty.
func (q *Queue[T]) PeekLast() (T, bool) {
	if len(q.elements) == 0 {
		var zero T
		return zero, false
	}
	return q.elements[len(q.elements)-1], true
}
