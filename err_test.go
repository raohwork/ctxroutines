// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"errors"
	"testing"
)

func testRun(expect []error) func(*testing.T) {
	return func(t *testing.T) {
		l := len(expect)
		ran := make([]bool, l)
		rs := make([]Runner, l)
		for x := 0; x < l; x++ {
			y := x
			rs[x] = NoCancelRunner(func() error {
				ran[y] = true
				return expect[y]
			})
		}

		actual := Run(rs...)

		for x := 0; x < l; x++ {
			if !ran[x] {
				t.Fatalf("runner #%d not ran", x)
			}
		}
		if v := len(actual); l != v {
			t.Fatalf("expected %d items, got %d", l, v)
		}

		for x := 0; x < l; x++ {
			if actual[x] != expect[x] {
				t.Errorf(
					"expected runner #%d returns %v, got %v",
					x, expect[x], actual[x],
				)
			}
		}
	}
}

func TestRunRunners(t *testing.T) {
	t.Run("ok1", testRun([]error{nil}))
	t.Run("ok3", testRun([]error{nil, nil, nil}))
	t.Run("err1", testRun([]error{errors.New("1")}))
	t.Run("err3", testRun([]error{
		errors.New("1"),
		errors.New("2"),
		errors.New("3"),
	}))

	t.Run("mixed", testRun([]error{
		errors.New("1"),
		nil,
		errors.New("3"),
	}))
}

func TestFirstErr(t *testing.T) {
	e1 := errors.New("1")
	e2 := errors.New("2")
	t.Run("ok1", testFirstErr(
		[]bool{true},
		nil,
		[]error{nil},
	))
	t.Run("ok3", testFirstErr(
		[]bool{true, true, true},
		nil,
		[]error{nil, nil, nil},
	))

	t.Run("err1", testFirstErr(
		[]bool{true},
		e1,
		[]error{e1},
	))
	t.Run("err3", testFirstErr(
		[]bool{true, false, false},
		e1,
		[]error{e1, e2, e2},
	))

	t.Run("mixed1", testFirstErr(
		[]bool{true, true, true},
		e1,
		[]error{nil, nil, e1},
	))
	t.Run("mixed2", testFirstErr(
		[]bool{true, true, false},
		e1,
		[]error{nil, e1, e2},
	))
}

func testFirstErr(expectRan []bool, expectErr error, runner []error) func(*testing.T) {
	return func(t *testing.T) {
		l := len(expectRan)
		if x := len(runner); x != l {
			t.Fatal("incorrect test case!")
		}

		ran := make([]bool, l)
		rs := make([]Runner, l)
		for x := 0; x < l; x++ {
			y := x
			rs[x] = NoCancelRunner(func() error {
				ran[y] = true
				return runner[y]
			})
		}

		actual := FirstErr(rs...).Run()
		if actual != expectErr {
			t.Fatal("unexpected error:", actual)
		}

		for x := 0; x < l; x++ {
			if ran[x] != expectRan[x] {
				s := "skipped"
				if expectRan[x] {
					s = "ran"
				}
				v := "skipped"
				if ran[x] {
					v = "ran"
				}
				t.Errorf(
					"expected runner #%d %s, but %s",
					x, s, v,
				)
			}
		}
	}
}

func testSomeErr(expect error, errs []error) func(*testing.T) {
	return func(t *testing.T) {
		rs := make([]Runner, len(errs))
		for idx, e := range errs {
			er := e
			rs[idx] = NoCancelRunner(func() error { return er })
		}
		if e := SomeErr(rs...).Run(); expect != e {
			t.Fatal("unexpected error:", e)
		}
	}
}

func TestSomeErr(t *testing.T) {
	e := errors.New("1")
	f := context.Canceled
	t.Run("ok", testSomeErr(nil, []error{nil, nil, nil}))
	t.Run("err", testSomeErr(e, []error{e, e}))
	t.Run("err1", testSomeErr(e, []error{e, nil, nil}))
	t.Run("err2", testSomeErr(e, []error{nil, e, nil}))
	t.Run("cancel", testSomeErr(f, []error{f, f}))
	t.Run("cancel1", testSomeErr(f, []error{f, nil, nil}))
	t.Run("cancel2", testSomeErr(f, []error{nil, f, nil}))
	t.Run("mixed1", testSomeErr(e, []error{e, f, nil}))
	t.Run("mixed2", testSomeErr(e, []error{f, e, nil}))
	t.Run("mixed3", testSomeErr(e, []error{nil, f, e}))
}
