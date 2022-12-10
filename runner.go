// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package ctxroutines prevides some helpers to write common routines and handle
// gracefully shutdown procedure
package ctxroutines

import (
	"context"
)

func newNil() Runner {
	return CTXRunner(func(c context.Context) error { return c.Err() })
}

// Runner defines a cancelable function
type Runner interface {
	Context() context.Context
	Cancel()
	// Run SHOULD always return some error after canceled
	Run() error
}

func IsCanceled(r Runner) bool {
	select {
	case <-r.Context().Done():
		return true
	default:
		return false
	}
}

// NoCancelRunner represents a function that cannot be canceled.
//
// In other words, Cancel() is always ignored.
type NoCancelRunner func() error

func (f NoCancelRunner) Cancel()                  {}
func (f NoCancelRunner) Run() error               { return f() }
func (f NoCancelRunner) Context() context.Context { return context.Background() }

type funcRunner struct {
	ctx    context.Context
	cancel func()
	f      func() error
}

func (r *funcRunner) Context() context.Context { return r.ctx }
func (r *funcRunner) Cancel()                  { r.cancel() }
func (r *funcRunner) Run() error               { return r.f() }

// NonInterruptRunner creates a runner that calling Cancel() does not interrupt it
//
// In other words, Cancel() only affects further Run(), which always returns context.Canceled.
func NonInterruptRunner(f func() error) Runner {
	return CTXRunner(func(c context.Context) error {
		if err := c.Err(); err != nil {
			return err
		}

		return f()
	})
}

// FromRunner reuses context and cancel function from r, but runs different function
func FromRunner(r Runner, f func() error) Runner {
	return NewRunner(r.Context(), r.Cancel, f)
}

// NewRunner creates a basic runner
func NewRunner(ctx context.Context, cancel context.CancelFunc, f func() error) Runner {
	return &funcRunner{
		ctx:    ctx,
		cancel: cancel,
		f:      f,
	}
}

// FuncRunner creates a runner from function
//
// Typical usage is like FuncRunner(someStruct.Cancel, someStruct.Run) or
//
//     srv := &http.Server{Addr: ":8080"}
//     r := FuncRunner(srv.Shutdown, srv.ListenAndServe)
func FuncRunner(cancel context.CancelFunc, f func() error) Runner {
	ctx, cf := context.WithCancel(context.Background())
	return NewRunner(ctx, func() { cf(); cancel() }, f)
}

// CTXRunner creates a Runner from a context-controlled function
//
// Typical usage is to wrap a cancelable function for further use (like, passing to
// Loop()
//
// You have to call Cancel() to release resources.
func CTXRunner(f func(context.Context) error) Runner {
	return CTXRunnerWith(context.Background(), f)
}

// CTXRunnerWith creates a Runner from a context-controlled function with
// predefined context
//
// Typical usage is to wrap a cancelable function for further use (like, passing to
// Loop()
//
// You have to call Cancel() to release resources.
func CTXRunnerWith(ctx context.Context, f func(context.Context) error) Runner {
	ctx, cancel := context.WithCancel(ctx)
	return NewRunner(ctx, cancel, func() error { return f(ctx) })
}
