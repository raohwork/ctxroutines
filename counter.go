// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"sync/atomic"
)

type Counter interface {
	Cancel()
	Run(n uint64) error
}

type NoCancelCounter func(uint64) error

func (f NoCancelCounter) Cancel()            {}
func (f NoCancelCounter) Run(n uint64) error { return f(n) }

type funcCounter struct {
	cancel func()
	f      func(uint64) error
}

func (r *funcCounter) Cancel()            { r.cancel() }
func (r *funcCounter) Run(n uint64) error { return r.f(n) }

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

type RecordedRunner interface {
	Runner
	Count() uint64
}

func Recorded(c Counter) (ret RecordedRunner) {
	var n uint64
	return &recorded{n: &n, c: c}
}

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
