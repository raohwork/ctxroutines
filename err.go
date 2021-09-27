// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"sync"
)

// CancelAll creates a function that calls Cancel() for every Runner of rs
func CancelAll(rs ...Runner) context.CancelFunc {
	return func() {
		for _, r := range rs {
			r.Cancel()
		}
	}
}

// Run runs every Runner of rs in separated goroutine, blocks til done, and returns all result
func Run(rs ...Runner) (err []error) {
	wg := sync.WaitGroup{}
	l := len(rs)
	wg.Add(l)
	err = make([]error, l)
	for idx, r := range rs {
		go func(idx int, r Runner) {
			err[idx] = r.Run()
			wg.Done()
		}(idx, r)
	}

	wg.Wait()
	return
}

// FirstErr creates a Runner that runs every Runner of rs in order, until first error occured
func FirstErr(rs ...Runner) (ret Runner) {
	return FuncRunner(CancelAll(rs...), func() (err error) {
		for _, r := range rs {
			if err = r.Run(); err != nil {
				return
			}
		}

		return
	})
}

// SomeErr creates a Runner runs every Runner of rs, and returns an error if there's one
//
//   - It checks error by the order of rs
//   - Returns first non-context.Canceled error
//   - Returns context.Canceled if no other errors
//   - Returns nil if everything's fine
func SomeErr(rs ...Runner) (ret Runner) {
	return FuncRunner(CancelAll(rs...), func() (err error) {
		errs := Run(rs...)
		canceled := false
		for _, err = range errs {
			if err == context.Canceled {
				canceled = true
				continue
			}
			if err != nil {
				return
			}
		}

		if canceled {
			err = context.Canceled
		}

		return
	})
}

// AnyErr creates a Runner that returns first known error.
func AnyErr(rs ...Runner) (ret Runner) {
	return FuncRunner(CancelAll(rs...), func() (err error) {
		ch := make(chan error)
		for _, r := range rs {
			go func(r Runner) {
				ch <- r.Run()
			}(r)
		}

		cur, max := 0, len(rs)
		for i := 0; i < max; i++ {
			err = <-ch
			cur++
			if err == nil {
				continue
			}

			go func() {
				for cur < max {
					<-ch
					cur++
				}
				close(ch)
			}()
			return
		}

		close(ch)
		return
	})
}
