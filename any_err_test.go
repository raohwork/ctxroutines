// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"errors"
	"testing"
)

func TestAnyErrOK(t *testing.T) {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	f1 := NoCancelRunner(func() error { close(ch1); return nil })
	f2 := NoCancelRunner(func() error { close(ch2); return nil })

	err := AnyErr(f1, f2).Run()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	select {
	case <-ch1:
	default:
		t.Fatal("f1 not stopped")
	}
	select {
	case <-ch2:
	default:
		t.Fatal("f2 not stopped")
	}
}

func TestAnyErr(t *testing.T) {
	wait := make(chan struct{})
	started := make(chan struct{})
	e := errors.New("")
	f1 := NoCancelRunner(func() error { <-started; return e })
	f2 := NewStatefulRunner(NoCancelRunner(func() error {
		close(started)
		<-wait
		return nil
	}))

	err := AnyErr(f1, f2).Run()
	<-started
	if err != e {
		t.Fatal("unexpected error:", err)
	}
	if !f2.IsRunning() {
		t.Fatal("f2 is not running")
	}
	close(wait)
}
