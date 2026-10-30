package list

import (
	stdlist "container/list"
	"iter"
)

// List is a type-safe wrapper around the standard library's container/list,
// scoped to what we need: insertion at the back, removal of an arbitrary
// element, and access to the front. All operations are O(1).
type List[T any] struct {
	internal stdlist.List
}

func NewList[T any]() *List[T] {
	return &List[T]{}
}

// Element is a handle to an item in a List, needed to Remove it.
type Element[T any] struct {
	internal *stdlist.Element
}

func (e Element[T]) Value() T {
	return e.internal.Value.(T)
}

func (l *List[T]) PushBack(item T) Element[T] {
	return Element[T]{l.internal.PushBack(item)}
}

// Remove removes e from the list. Removing an element that has already been
// removed is a no-op.
func (l *List[T]) Remove(e Element[T]) {
	l.internal.Remove(e.internal)
}

func (l *List[T]) Front() (T, bool) {
	if e := l.internal.Front(); e != nil {
		return e.Value.(T), true
	} else {
		var zero T
		return zero, false
	}
}

func (l *List[T]) Len() int {
	return l.internal.Len()
}

// All iterates from front to back. The list must not be modified during
// iteration.
func (l *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.internal.Front(); e != nil; e = e.Next() {
			if !yield(e.Value.(T)) {
				return
			}
		}
	}
}
