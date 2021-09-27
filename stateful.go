// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"sync"
)

// StatefulRunner represents a Runner that can inspect its running state
type StatefulRunner interface {
	Runner
	IsRunning() bool
	// Try to run the Runner if it's not running. ran will be true if it runs.
	TryRun() (err error, ran bool)
	// Lock the runner, pretending that it's running. Call release() to unlock.
	// It's safe to call release more than once.
	//
	// It will blocks until
	Lock() (release func())
}

type statefulRunner struct {
	token chan bool
	f     Runner
}

func (f *statefulRunner) IsRunning() (yes bool) {
	select {
	case <-f.token:
		f.token <- true
		return false
	default:
		return true
	}
}

func (f *statefulRunner) TryRun() (err error, ran bool) {
	select {
	case <-f.token:
		f.token <- true
		return f.Run(), true
	default:
		return
	}
}

func (f *statefulRunner) Cancel() {
	f.f.Cancel()
}

func (f *statefulRunner) Run() (err error) {
	<-f.token

	err = f.f.Run()
	f.token <- true
	return
}

func (f *statefulRunner) Lock() (release func()) {
	<-f.token
	once := &sync.Once{}

	return func() {
		once.Do(func() {
			f.token <- true
		})
	}
}

// NewStatefulRunner creates a StatefulRunner from existing Runner
func NewStatefulRunner(f Runner) (ret StatefulRunner) {
	x := &statefulRunner{
		token: make(chan bool, 1),
		f:     f,
	}
	x.token <- true

	return x
}
