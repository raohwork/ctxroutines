// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"errors"
	"testing"
)

func TestTryAtMost1OK(t *testing.T) {
	ch := make(chan bool, 10)
	f := NoCancelRunner(func() error { ch <- true; return nil })
	r := TryAtMost(10, f)
	err := r.Run()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if l := len(ch); l != 1 {
		t.Fatalf("expected 1 item, got %d", l)
	}
}

func TestTryAtMost1OK2Fail(t *testing.T) {
	e := errors.New("")
	ch := make(chan bool, 10)
	f := Recorded(NoCancelCounter(func(n uint64) error {
		ch <- true
		if n < 2 {
			return e
		}
		return nil
	}))
	r := TryAtMost(10, f)
	err := r.Run()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if l := len(ch); l != 3 {
		t.Fatalf("expected 3 items, got %d", l)
	}
}

func TestTryAtMost1OK9Fail(t *testing.T) {
	e := errors.New("")
	ch := make(chan bool, 10)
	f := Recorded(NoCancelCounter(func(n uint64) error {
		ch <- true
		if n < 9 {
			return e
		}
		return nil
	}))
	r := TryAtMost(10, f)
	err := r.Run()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if l := len(ch); l != 10 {
		t.Fatalf("expected 10 items, got %d", l)
	}
}

func TestTryAtMostFail(t *testing.T) {
	e := errors.New("")
	ch := make(chan bool, 11)
	f := NoCancelRunner(func() error {
		ch <- true
		return e
	})
	r := TryAtMost(10, f)
	err := r.Run()
	if err != e {
		t.Fatal("unexpected error:", err)
	}
	if l := len(ch); l != 10 {
		t.Fatalf("expected 10 items, got %d", l)
	}
}
