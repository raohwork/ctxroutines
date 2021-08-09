// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"sync"
	"time"
)

type ctxFactory func() (context.Context, context.CancelFunc)

type bypass func(error) bool

func onlySuccess(err error) (skip bool) { return err != nil }
func onlyFail(err error) (skip bool)    { return err == nil }
func noBypass(err error) (skip bool)    { return }

type runAtLeast struct {
	dur time.Duration
	f   Runner
	bypass
	canceled chan struct{}
	closer   sync.Once
}

func (r *runAtLeast) Run() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.dur)
	defer cancel()

	select {
	case <-r.canceled:
		return context.Canceled
	default:
	}

	err = r.f.Run()

	if !r.bypass(err) {
		select {
		case <-r.canceled:
			err = context.Canceled
		case <-ctx.Done():
		}
	}

	return
}

func (r *runAtLeast) Cancel() {
	r.f.Cancel()
	r.closer.Do(func() { close(r.canceled) })
}

func newRAL(dur time.Duration, f Runner, b bypass) (ret *runAtLeast) {
	ret = &runAtLeast{
		dur:      dur,
		f:        f,
		bypass:   b,
		canceled: make(chan struct{}),
	}
	return
}

// RunAtLeast ensures the execution time of f is longer than dur
func RunAtLeast(dur time.Duration, f Runner) (ret Runner) {
	return newRAL(dur, f, noBypass)
}

func RunAtLeastSuccess(dur time.Duration, f Runner) (ret Runner) {
	return newRAL(dur, f, onlySuccess)
}

func RunAtLeastFailed(dur time.Duration, f Runner) (ret Runner) {
	return newRAL(dur, f, onlyFail)
}
