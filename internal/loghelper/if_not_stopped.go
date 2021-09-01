package loghelper

import (
	"time"

	"github.com/smartcontractkit/libocr/subprocesses"
)

type IfNotStopped struct {
	chStop chan struct{}
	subs   subprocesses.Subprocesses
}

// If Stop is called prior to expiry of d, f won't be executed. Otherwise, f
// will be executed and Stop will block until f returns.
func NewIfNotStopped(d time.Duration, f func()) *IfNotStopped {
	ins := IfNotStopped{
		make(chan struct{}, 1),
		subprocesses.Subprocesses{},
	}
	ins.subs.Go(func() {
		t := time.NewTimer(d)
		defer t.Stop()
		select {
		case <-t.C:
			f()
		case <-ins.chStop:
		}
	})
	return &ins
}

func (ins *IfNotStopped) Stop() {
	select {
	case <-ins.chStop:
		// chStopp has been closed, don't close again
	default:
		close(ins.chStop)
	}

	ins.subs.Wait()
}
