// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"sync/atomic"
)

// Counter represents a cancelable function that receives how many time it has been run
type Counter interface {
	// cancel this counter, further Run() calls always return context.Cancel
	Cancel()
	// run the counter once. n represents how many times have been run, so
	// first call will be Run(0)
	Run(n uint64) error
}

// NoCancelCounter returns a Counter that cannot be canceled
type NoCancelCounter func(uint64) error

func (f NoCancelCounter) Cancel()            {}
func (f NoCancelCounter) Run(n uint64) error { return f(n) }

type funcCounter struct {
	cancel func()
	f      func(uint64) error
}

func (r *funcCounter) Cancel()            { r.cancel() }
func (r *funcCounter) Run(n uint64) error { return r.f(n) }

// FuncCounter wraps a counter function into a Counter
//
// f SHOULD always return error after cancel is called.
func FuncCounter(cancel context.CancelFunc, f func(uint64) error) (ret Counter) {
	return &funcCounter{cancel: cancel, f: f}
}

type recorded struct {
	n *uint64
	c Counter
}

func (r *recorded) Cancel() { r.c.Cancel() }
func (r *recorded) Run() error {
	n := atomic.AddUint64(r.n, 1)
	return r.c.Run(n - 1)
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
func Recorded(c Counter) (ret RecordedRunner) {
	var n uint64
	return &recorded{n: &n, c: c}
}

// TryAtMostWithChan creates a Runner that runs f for at most n times before it returns nil
//
// Every error ever tried are send through ch, it will not run f before sending the
// error into ch. nil is also sent if f returns it.
//
// It will not close ch since you might want to call ret.Run() many times.
func TryAtMostWithChan(n uint64, f Runner, ch chan error) (ret Runner) {
	return FuncRunner(f.Cancel, func() (err error) {
		r := Recorded(FuncCounter(f.Cancel, func(idx uint64) error {
			return f.Run()
		}))
		for r.Count() < n {
			err = r.Run()
			ch <- err
			if err == nil {
				return
			}
		}

		return
	})
}

// TryAtMost creates a Runner that runs f for at most n times before it returns nil
//
// Say you have a f only success on third try:
//
//     r := FuncRunner(5, f)
//     r.Run() // run f 3 times, and returs nil
//
//     r := FuncRunner(2, f)
//     r.Run() // run f 2 times, and returs the error from second run
//
//     r := FuncRunner(3, f)
//     r.Run() // run f 3 times, and returs nil
func TryAtMost(n uint64, f Runner) (ret Runner) {
	return FuncRunner(f.Cancel, func() (err error) {
		r := Recorded(FuncCounter(f.Cancel, func(idx uint64) error {
			return f.Run()
		}))
		for r.Count() < n {
			err = r.Run()
			if err == nil {
				return
			}
		}

		return
	})
}
