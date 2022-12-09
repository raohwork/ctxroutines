package ctxroutines

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// RatelimitRunner creates a Runner that respects the rate limit.
//
// Say you have an empty Runner r with rate limit to once per second:
//
//     go r.Run() // returns immediatly
//     go r.Run() // returns after a second
//
// Once the Runner has been canceled, Run() always return context.Canceled.
//
// RatelimitRunner() looks like RunAtLeast,
func RatelimitRunner(l *rate.Limiter, r Runner) (ret Runner) {
	return &ratelimitRunner{
		lim:    l,
		Runner: r,
	}
}

// NoLessThan is a wrapper of RatelimitRunner
func NoLessThan(dur time.Duration, f Runner) Runner {
	return RatelimitRunner(rate.NewLimiter(rate.Every(dur), 1), f)
}

type ratelimitRunner struct {
	lim *rate.Limiter
	Runner
}

func (r *ratelimitRunner) sleep(timeout time.Duration) (canceled bool) {
	select {
	case <-r.Context().Done():
		return true
	case <-time.After(timeout):
		return false
	}
}

func (r *ratelimitRunner) Run() (err error) {
	if IsCanceled(r) {
		return context.Canceled
	}

	reserve := r.lim.Reserve()
	if r.sleep(reserve.Delay()) {
		reserve.Cancel()
		return context.Canceled
	}

	return r.Runner.Run()
}
