// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import "testing"

func TestStatefulIsRunning(t *testing.T) {
	start := make(chan int)
	defer close(start)
	wait := make(chan int)
	defer close(wait)
	done := make(chan int)
	defer close(done)

	f := NoCancelRunner(func() error {
		start <- 1
		<-wait
		return nil
	})
	r := NewStatefulRunner(f)

	if r.IsRunning() {
		t.Fatal("expected not running, but it is")
	}
	go func() { r.Run(); done <- 1 }()
	<-start
	if !r.IsRunning() {
		t.Fatal("expected running, but not")
	}
	wait <- 1
	<-done
	if r.IsRunning() {
		t.Fatal("expected not running, but it is")
	}
}

func TestStatefulTryRun(t *testing.T) {
	start := make(chan int)
	defer close(start)
	wait := make(chan int)
	defer close(wait)
	done := make(chan int)
	defer close(done)

	f := NoCancelRunner(func() error {
		start <- 1
		<-wait
		return nil
	})
	r := NewStatefulRunner(f)

	go func() { r.Run(); done <- 1 }()
	<-start
	err, ran := r.TryRun()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if ran {
		t.Fatal("expected skipped, but ran")
	}
	wait <- 1
	<-done

	go func() {
		<-start
		wait <- 1
	}()
	err, ran = r.TryRun()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if !ran {
		t.Fatal("expected ran, but not")
	}
}

func TestStatefulLockIsRunning(t *testing.T) {
	f := NoCancelRunner(func() error { return nil })
	r := NewStatefulRunner(f)

	release := r.Lock()
	if !r.IsRunning() {
		t.Fatal("expected running, but not")
	}

	release()
	if r.IsRunning() {
		t.Fatal("expected not running, but it is")
	}
}
