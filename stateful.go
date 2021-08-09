// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"sync"
)

type StatefulRunner interface {
	Runner
	IsRunning() bool
	TryRun() (err error, ran bool)
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

func NewStatefulRunner(f Runner) (ret StatefulRunner) {
	x := &statefulRunner{
		token: make(chan bool, 1),
		f:     f,
	}
	x.token <- true

	return x
}
