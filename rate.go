package ctxroutines

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RatelimitRunner creates a Runner that respects the rate limit.
//
// Say you have an empty Runner r with rate limit to once per second:
//
//     r.Run() // returns immediatly
//     r.Run() // returns after a second
//
// Once the Runner has been canceled, Run() always return context.Canceled.
//
// RatelimitRunner(rate.NewLimiter(rate.Every(time.Second), 1), r) is identical with
// RunAtLeast(time.Second, r). However:
//
//   - It's impossible to provide RatelimitSuccess (or Failed) due to restrictions
//     of golang.org/x/time/rate.
//   - RunAtLeast does not support bursting.
func RatelimitRunner(l *rate.Limiter, r Runner) (ret Runner) {
	return &ratelimitRunner{
		lim:      l,
		f:        r,
		canceled: make(chan struct{}),
	}
}

type ratelimitRunner struct {
	lim      *rate.Limiter
	f        Runner
	canceled chan struct{}
	closer   sync.Once
}

func (r *ratelimitRunner) sleep(timeout time.Duration) (canceled bool) {
	select {
	case <-r.canceled:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (r *ratelimitRunner) Run() (err error) {
	select {
	case <-r.canceled:
		return context.Canceled
	default:
	}

	reserve := r.lim.Reserve()
	if r.sleep(reserve.Delay()) {
		reserve.Cancel()
		return context.Canceled
	}

	return r.f.Run()
}

func (r *ratelimitRunner) Cancel() {
	r.closer.Do(func() {
		r.f.Cancel()
		close(r.canceled)
	})
}
