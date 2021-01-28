// Package subprocesses keeps track of concurrent processes,
// for coordination of cleanly shutting down systems of goroutines. This is a
// stripped-down version of errgroup.Group, motivated by the fact that allowing
// a single process to shut down the entire system by returning an error is
// quite fragile.
package subprocesses

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Subprocesses struct {
	wg sync.WaitGroup
}

// Wait blocks until all function calls from the Go method have returned.
func (s *Subprocesses) Wait() {
	s.wg.Wait()
}

// Go calls the given function in a new goroutine.
func (s *Subprocesses) Go(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}

// BlockForAtMost invokes f and blocks for at most duration d before returning,
// regardless of whether f finished or not, or the passed in ctx is cancelled.
// If f finished, returns true.
// Otherwise, returns false.
func (s *Subprocesses) BlockForAtMost(ctx context.Context, d time.Duration, f func(context.Context)) (ok bool) {
	done := make(chan struct{})
	childCtx, childCancel := context.WithTimeout(ctx, d)
	defer childCancel()
	s.Go(func() {
		f(childCtx)
		close(done)
	})
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-done:
		return true
	case <-t.C:
		return false
	}
}

// RepeatWithCancel repeats f with the specified interval. Cancel if ctx.Done is signaled
func (s *Subprocesses) RepeatWithCancel(name string, interval time.Duration, ctx context.Context, f func()) {
	s.Go(func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("canceling", name)
				return
			case <-ticker.C:
				f()
			}
		}
	})
}
