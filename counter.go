// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"sync/atomic"
)

// Counter creates RecordedRunner from function
func Counter(f func(uint64) error) RecordedRunner {
	n := uint64(0)
	return &recorded{
		n: &n,
		Runner: CTXRunner(func(c context.Context) error {
			if c.Err() != nil {
				return c.Err()
			}

			return f(atomic.LoadUint64(&n))
		}),
	}
}

type recorded struct {
	n *uint64
	Runner
}

func (r *recorded) Run() (err error) {
	err = r.Runner.Run()
	atomic.AddUint64(r.n, 1)
	return
}
func (r *recorded) Count() uint64 {
	return atomic.LoadUint64(r.n)
}

// RecordedRunner is a Runner remembers how many times has been run
type RecordedRunner interface {
	Runner
	Count() uint64
}

// Recorded creates a RecordedRunner
func Recorded(r Runner) (ret RecordedRunner) {
	var n uint64
	return &recorded{n: &n, Runner: r}
}

// TryAtMost creates a Runner that runs f for at most n times before it returns nil
//
// Say you have a f only success on third try:
//
//     r := TryAtMost(5, f)
//     r.Run() // run f 3 times, and returs nil
//
//     r := TryAtMost(2, f)
//     r.Run() // run f 2 times, and returs the error from second run
//
//     r := TryAtMost(3, f)
//     r.Run() // run f 3 times, and returs nil
func TryAtMost(n uint64, f Runner) (ret Runner) {
	r := Recorded(f)
	return FromRunner(r, func() (err error) {
		for r.Count() < n {
			if IsCanceled(r) {
				return context.Canceled
			}

			err = r.Run()
			if err == nil {
				return
			}
		}

		return
	})
}
