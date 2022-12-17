// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import "context"

type loopRunner struct {
	Runner
	cb   func(error)
	term func(error) (err error, term bool)
}

func (r *loopRunner) Run() (err error) {
	for {
		var term bool
		select {
		case <-r.Context().Done():
			return r.Context().Err()
		default:
		}

		err = r.Runner.Run()

		err, term = r.term(err)

		if term {
			return
		}
		if err != nil {
			r.cb(err)
		}
	}
}

// Retry creates a Runner runs r until it returns nil
func Retry(r Runner) (ret Runner) {
	return RetryWithCB(r, func(error) {})
}

// RetryWithCB creates a Runner runs r until it returns nil
//
// It calls cb if r returns error.
//
// You have to call Cancel() to release resources.
func RetryWithCB(r Runner, cb func(error)) (ret Runner) {
	return &loopRunner{
		Runner: r,
		cb:     cb,
		term: func(e error) (err error, term bool) {
			if e == nil || e == context.Canceled {
				return e, true
			}
			return e, false
		},
	}
}

// RetryWithChan is shortcut for RetryWithCB(r, func(e error) { ch <- e })
//
// You have to call Cancel() to release resources.
func RetryWithChan(r Runner, ch chan<- error) (ret Runner) {
	return RetryWithCB(r, func(e error) { ch <- e })
}

// TilErr creates a Runner runs r until it returns any error
//
// You have to call Cancel() to release resources.
func TilErr(r Runner) (ret Runner) {
	return &loopRunner{
		Runner: r,
		cb:     func(error) {},
		term: func(e error) (err error, term bool) {
			if e == nil {
				return
			}
			return e, true
		},
	}
}

// Loop creates a Runner that runs r until canceled
//
// You have to call Cancel() to release resources.
func Loop(r Runner) (ret Runner) {
	return &loopRunner{
		Runner: r,
		cb:     func(error) {},
		term: func(e error) (err error, term bool) {
			if e != context.Canceled {
				return
			}
			return e, true
		},
	}
}
