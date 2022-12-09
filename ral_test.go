// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
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
}

type ralCase struct{}

func (c *ralCase) ral(f Runner) (ret *runAtLeast) {
	return newRAL(20*time.Millisecond, f, noBypass)
}

func (c *ralCase) Sync(t *testing.T) {
	const z = 4
	f := NoCancelRunner(func() error {
		return nil
	})

	r := c.ral(f)

	for i := 0; i < z; i++ {
		dur := int64(cost(func() { r.Run() }) / time.Millisecond)
		if dur < 20 {
			t.Fatalf("run time #%d less than 20ms: %d", i, dur)
		}
	}
}

func (c *ralCase) Async(t *testing.T) {
	const z = 4

	f := NoCancelRunner(func() error { return nil })
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
}
