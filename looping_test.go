// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	ch := make(chan string, 3)
	cnt := 0
	f := NoCancelRunner(func() error {
		if cnt < 2 {
			ch <- "err"
			cnt++
			return errors.New("")
		}
		ch <- "ok"
		return nil
	})

	err := Retry(f).Run()
	close(ch)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	check := func(expect, actual string) {
		if expect != actual {
			t.Fatalf("expected '%s', go '%s'", expect, actual)
		}
	}
	check("err", <-ch)
	check("err", <-ch)
	check("ok", <-ch)
	check("", <-ch)
}

func TestRetryCancel(t *testing.T) {
	f := CTXRunner(func(c context.Context) error {
		return errors.New("")
	})

	r := Retry(f)
	go func() { time.Sleep(time.Millisecond); r.Cancel() }()
	err := r.Run()
	if err != context.Canceled {
		t.Fatal("unexpected error:", err)
	}
}

func TestTilErr(t *testing.T) {
	ch := make(chan string, 3)
	cnt := 0
	e := errors.New("")
	f := CTXRunner(func(c context.Context) error {
		if cnt < 2 {
			ch <- "ok"
			cnt++
			return nil
		}
		ch <- "err"
		return e
	})

	err := TilErr(f).Run()
	close(ch)
	if err != e {
		t.Fatal("unexpected error:", err)
	}
	check := func(expect, actual string) {
		if expect != actual {
			t.Fatalf("expected '%s', go '%s'", expect, actual)
		}
	}
	check("ok", <-ch)
	check("ok", <-ch)
	check("err", <-ch)
	check("", <-ch)
}

func TestTilErrCancel(t *testing.T) {
	f := CTXRunner(func(c context.Context) error {
		return nil
	})

	r := TilErr(f)
	go func() { time.Sleep(time.Millisecond); r.Cancel() }()
	err := r.Run()
	if err != context.Canceled {
		t.Fatal("unexpected error:", err)
	}
}

func TestLoopCancel(t *testing.T) {
	f := CTXRunner(func(c context.Context) error {
		return context.Canceled
	})

	r := Loop(f)
	err := r.Run()
	if err != context.Canceled {
		t.Fatal("unexpected error:", err)
	}
}

func TestLoopCancelNil(t *testing.T) {
	f := CTXRunner(func(c context.Context) error {
		return nil
	})

	r := Loop(f)
	go func() { time.Sleep(time.Millisecond); r.Cancel() }()
	err := r.Run()
	if err != context.Canceled {
		t.Fatal("unexpected error:", err)
	}
}

func TestLoopCancelErr(t *testing.T) {
	f := CTXRunner(func(c context.Context) error {
		return errors.New("")
	})

	r := Loop(f)
	go func() { time.Sleep(time.Millisecond); r.Cancel() }()
	err := r.Run()
	if err != context.Canceled {
		t.Fatal("unexpected error:", err)
	}
}
