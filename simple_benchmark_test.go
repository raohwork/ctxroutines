// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"errors"
	"testing"
)

func BenchmarkRunAtLeast(b *testing.B) {
	f := NoCancelRunner(func() error { return nil })
	r := RunAtLeast(0, f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Run()
	}
}

func BenchmarkOnceWithin(b *testing.B) {
	f := NoCancelRunner(func() error { return nil })
	r := OnceWithin(0, f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Run()
	}
}

func BenchmarkRetry(b *testing.B) {
	f := NoCancelRunner(func() error { return nil })
	r := Retry(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Run()
	}
}
func BenchmarkTilErr(b *testing.B) {
	e := errors.New("")
	f := NoCancelRunner(func() error { return e })
	r := TilErr(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Run()
	}
}

func BenchmarkStatefulRunner(b *testing.B) {
	f := NoCancelRunner(func() error { return nil })
	r := NewStatefulRunner(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Run()
	}
}

func BenchmarkRecorded(b *testing.B) {
	f := NoCancelCounter(func(n uint64) error { return nil })
	r := Recorded(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Run()
	}
}

func BenchmarkTryAtMost(b *testing.B) {
	r := NoCancelRunner(func() error { return nil })
	x := TryAtMost(1, r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Run()
	}
}

func BenchmarkTryAtMost2(b *testing.B) {
	e := errors.New("")
	r := NoCancelRunner(func() error { return e })
	x := TryAtMost(2, r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Run()
	}
}

func BenchmarkAnyErr(b *testing.B) {
	r := NoCancelRunner(func() error { return nil })
	x := AnyErr(r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Run()
	}
}

func BenchmarkSomeErr(b *testing.B) {
	r := NoCancelRunner(func() error { return nil })
	x := SomeErr(r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Run()
	}
}

func BenchmarkFirstErr(b *testing.B) {
	r := NoCancelRunner(func() error { return nil })
	x := FirstErr(r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Run()
	}
}
