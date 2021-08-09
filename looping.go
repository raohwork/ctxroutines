// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import "context"

type loopRunner struct {
	ctx    context.Context
	cancel context.CancelFunc
	r      Runner
	cb     func(error)
	term   func(error) (err error, term bool)
}

func (r *loopRunner) Cancel() {
	r.cancel()
	r.r.Cancel()
}

func (r *loopRunner) Run() (err error) {
	for {
		var term bool
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
		}

		err = r.r.Run()

		err, term = r.term(err)

		if term {
			return
		}
		if err != nil {
			r.cb(err)
		}
	}
}

func Retry(r Runner) (ret Runner) {
	return RetryWithCB(r, func(error) {})
}

func RetryWithCB(r Runner, cb func(error)) (ret Runner) {
	ctx, cancel := context.WithCancel(context.Background())
	return &loopRunner{
		ctx:    ctx,
		cancel: cancel,
		r:      r,
		cb:     cb,
		term: func(e error) (err error, term bool) {
			if e == nil || e == context.Canceled {
				return e, true
			}
			return e, false
		},
	}
}

func RetryWithChan(r Runner, ch chan<- error) (ret Runner) {
	return RetryWithCB(r, func(e error) { ch <- e })
}

func TilErr(r Runner) (ret Runner) {
	ctx, cancel := context.WithCancel(context.Background())
	return &loopRunner{
		ctx:    ctx,
		cancel: cancel,
		r:      r,
		cb:     func(error) {},
		term: func(e error) (err error, term bool) {
			if e == nil {
				return
			}
			return e, true
		},
	}
}

func Loop(r Runner) (ret Runner) {
	ctx, cancel := context.WithCancel(context.Background())
	return &loopRunner{
		ctx:    ctx,
		cancel: cancel,
		r:      r,
		cb:     func(error) {},
		term: func(e error) (err error, term bool) {
			if e != context.Canceled {
				return
			}
			return e, true
		},
	}
}
