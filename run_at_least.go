// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"time"
)

type bypass func(error) bool

func onlySuccess(err error) (skip bool) { return err != nil }
func onlyFail(err error) (skip bool)    { return err == nil }
func noBypass(err error) (skip bool)    { return }

type runAtLeast struct {
	dur time.Duration
	Runner
	bypass
}

func (r *runAtLeast) Run() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.dur)
	defer cancel()

	err = r.Runner.Run()

	if !r.bypass(err) {
		<-ctx.Done()
	}

	return
}

func newRAL(dur time.Duration, f Runner, b bypass) (ret *runAtLeast) {
	ret = &runAtLeast{
		dur:    dur,
		Runner: f,
		bypass: b,
	}
	return
}

// RunAtLeast ensures the execution time of f is longer than dur
//
// Say you have an empty Runner f
//
//     r := RunAtLeast(time.Second, f)
//     r.Run() // runs f immediately, blocks 1s
func RunAtLeast(dur time.Duration, f Runner) (ret Runner) {
	return newRAL(dur, f, noBypass)
}

// RunAtLeastSuccess is like RunAtLeast, but only successful call counts
func RunAtLeastSuccess(dur time.Duration, f Runner) (ret Runner) {
	return newRAL(dur, f, onlySuccess)
}

// RunAtLeastFailed is like RunAtLeast, but only failed call counts
func RunAtLeastFailed(dur time.Duration, f Runner) (ret Runner) {
	return newRAL(dur, f, onlyFail)
}
