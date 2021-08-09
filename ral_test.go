// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"sync"
	"testing"
	"time"
)

func cost(f func()) (ret time.Duration) {
	begin := time.Now()
	f()
	return time.Since(begin)
}

func TestRunAtLeast(t *testing.T) {
	c := &ralCase{}
	t.Run("sync", c.Sync)
	t.Run("async", c.Async)
	t.Run("cancelBefore", c.CancelBefore)
	t.Run("cancelAfter", c.CancelAfter)
}

type ralCase struct{}

func (c *ralCase) ral(f Runner) (ret *runAtLeast) {
	return newRAL(20*time.Millisecond, f, noBypass)
}

func (c *ralCase) Sync(t *testing.T) {
	const z = 4
	f := NoCancelRunner(func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	r := c.ral(f)

	for i := 0; i < z; i++ {
		dur := int64(cost(func() { r.Run() }) / time.Millisecond)
		if dur < 10 {
			t.Fatalf("run time #%d less than 10ms: %d", i, dur)
		}
	}
}

func (c *ralCase) Async(t *testing.T) {
	const z = 4

	f := NoCancelRunner(func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	r := c.ral(f)
	wg := &sync.WaitGroup{}
	wg.Add(z)

	for i := 0; i < z; i++ {
		go func() {
			r.Run()
			wg.Done()
		}()
	}

	dur := int64(cost(wg.Wait) / time.Millisecond)
	if dur < 20 {
		t.Fatalf("f should run at least 20ms, got %dms", dur)
	}
	if dur >= 22 {
		t.Fatalf("f should run about 20ms, got %dms", dur)
	}
}

func (c *ralCase) CancelBefore(t *testing.T) {
	done := false
	f := NoCancelRunner(func() error {
		time.Sleep(10 * time.Millisecond)
		done = true
		return nil
	})
	r := c.ral(f)
	r.Cancel()

	var err error
	dur := int64(cost(func() { err = r.Run() }) / time.Millisecond)
	if err != context.Canceled {
		t.Fatal("unexpected error: ", err)
	}
	if done {
		t.Fatalf("expected it doesn't run, detected run time %dms", dur)
	}
}

func (c *ralCase) CancelAfter(t *testing.T) {
	done := make(chan struct{})
	f := RunAtLeast(time.Second, NoCancelRunner(func() error {
		close(done)
		time.Sleep(10 * time.Millisecond)
		return nil
	}))
	r := c.ral(f)

	var err error
	var dur time.Duration
	detected := make(chan struct{})
	go func() {
		dur = cost(func() {
			err = r.Run()
			close(detected)
		})
	}()

	<-done
	r.Cancel()
	<-detected

	if err != context.Canceled {
		t.Fatal("unexpected error: ", err)
	}
	if dur >= time.Second {
		t.Fatalf("expected it doesn't run till end, detected run time %dms", dur/time.Millisecond)
	}
}
