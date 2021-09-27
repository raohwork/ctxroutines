// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package ctxroutines prevides some helpers to write common routines and handle
// gracefully shutdown procedure
package ctxroutines

import "context"

// Runner defines a cancelable function
type Runner interface {
	Cancel()
	// Run SHOULD always return some error after canceled
	Run() error
}

// NoCancelRunner represents a function that cannot be canceled.
//
// In other words, Cancel() is always ignored.
type NoCancelRunner func() error

func (f NoCancelRunner) Cancel()    {}
func (f NoCancelRunner) Run() error { return f() }

type funcRunner struct {
	cancel func()
	f      func() error
}

func (r *funcRunner) Cancel()    { r.cancel() }
func (r *funcRunner) Run() error { return r.f() }

// FuncRunner creates a runner from function
//
// Typical usage is like FuncRunner(someStruct.Cancel, someStruct.Run) or
//
//     srv := &http.Server{Addr: ":8080"}
//     r := FuncRunner(srv.Shutdown, srv.ListenAndServe)
func FuncRunner(cancel context.CancelFunc, f func() error) Runner {
	return &funcRunner{
		cancel: cancel,
		f:      f,
	}
}

// CTXRunner creates a Runner from a context-controlled function
//
// Typical usage is to wrap a cancelable function for further use (like, passing to
// Loop()
//
// You have to call Cancel() to release resources.
func CTXRunner(f func(context.Context) error) Runner {
	ctx, cancel := context.WithCancel(context.Background())
	return FuncRunner(cancel, func() error { return f(ctx) })
}
