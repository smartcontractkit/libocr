package scheduler

import (
	"context"
	"time"

	"github.com/RoSpaceDev/libocr/offchainreporting2plus/internal/common/minheap"
	"github.com/RoSpaceDev/libocr/subprocesses"
)

type itemWithDeadline[T any] struct {
	Item     T
	Deadline time.Time
}

type Scheduler[T any] struct {
	subs   subprocesses.Subprocesses
	ctx    context.Context
	cancel context.CancelFunc

	in  chan<- itemWithDeadline[T]
	out <-chan T
}

func NewScheduler[T any]() *Scheduler[T] {
	ctx, cancel := context.WithCancel(context.Background())

	in := make(chan itemWithDeadline[T])
	out := make(chan T)

	scheduler := &Scheduler[T]{
		subprocesses.Subprocesses{},
		ctx,
		cancel,

		in,
		out,
	}

	scheduler.subs.Go(func() {
		// create an expired timer
		timer := time.NewTimer(0)
		defer timer.Stop()
		<-timer.C

		heap := minheap.NewMinHeap(func(a, b itemWithDeadline[T]) bool {
			return a.Deadline.Before(b.Deadline)
		})

		var pendingItem T
		var maybeOut chan<- T

		for {
			select {
			case item := <-in:
				if maybeOut == nil {
					peeked, ok := heap.Peek()
					if !ok {
						// the timer must be stopped already
						timer.Reset(time.Until(item.Deadline))
					} else if peeked.Deadline.After(item.Deadline) {
						// we're dealing with the new minimum
						if timer.Stop() {
							// timer hasn't fired yet
							timer.Reset(time.Until(item.Deadline))
						} // else: timer has fired. no need to do anything since
						//   we will handle <-timer.C in an upcoming loop iteration
					}
				}
				heap.Push(item)
			case <-timer.C:
				popped, ok := heap.Pop()
				if ok {
					pendingItem = popped.Item
					maybeOut = out
				} else { //nolint:staticcheck
					// We should never enter this else branch. But if we did, it's
					// better to ignore the spurious firing of the timer than
					// to panic.
					// Tests should still pass with the panic not commented out.

					// panic("timer fired despite heap being empty, this should never happen")
				}
			case maybeOut <- pendingItem:
				maybeOut = nil
				peeked, ok := heap.Peek()
				if ok {
					timer.Reset(time.Until(peeked.Deadline))
				}
			case <-ctx.Done():
				return
			}
		}
	})

	return scheduler
}

func (s *Scheduler[T]) ScheduleDeadline(item T, deadline time.Time) {
	select {
	case s.in <- itemWithDeadline[T]{item, deadline}:
	case <-s.ctx.Done():
	}
}

func (s *Scheduler[T]) ScheduleDelay(item T, delay time.Duration) {
	s.ScheduleDeadline(item, time.Now().Add(delay))
}

func (s *Scheduler[T]) Scheduled() <-chan T {
	return s.out
}

func (s *Scheduler[T]) Close() {
	s.cancel()
	s.subs.Wait()
}
