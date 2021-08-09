// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import "context"

type Runner interface {
	Cancel()
	Run() error
}

type NoCancelRunner func() error

func (f NoCancelRunner) Cancel()    {}
func (f NoCancelRunner) Run() error { return f() }

type funcRunner struct {
	cancel func()
	f      func() error
}

func (r *funcRunner) Cancel()    { r.cancel() }
func (r *funcRunner) Run() error { return r.f() }

func FuncRunner(cancel context.CancelFunc, f func() error) Runner {
	return &funcRunner{
		cancel: cancel,
		f:      f,
	}
}

func CTXRunner(f func(context.Context) error) Runner {
	ctx, cancel := context.WithCancel(context.Background())
	return FuncRunner(cancel, func() error { return f(ctx) })
}
