// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"sync"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRatelimitRunnerSync(t *testing.T) {
	r := RatelimitRunner(
		rate.NewLimiter(rate.Every(10*time.Millisecond), 2),
		NoCancelRunner(func() error { return nil }),
	)

	begin := time.Now()
	r.Run()
	r.Run()
	r.Run()
	r.Run()
	used := time.Since(begin)
	if used < 20*time.Millisecond || used >= 30*time.Millisecond {
		t.Fatal("unexpected run time: ", int64(used/time.Millisecond))
	}
}

func TestRatelimitRunnerAsync(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(4)
	r := RatelimitRunner(
		rate.NewLimiter(rate.Every(10*time.Millisecond), 2),
		NoCancelRunner(func() error { wg.Done(); return nil }),
	)

	begin := time.Now()
	go r.Run()
	go r.Run()
	go r.Run()
	go r.Run()

	wg.Wait()
	used := time.Since(begin)
	if used < 20*time.Millisecond || used >= 30*time.Millisecond {
		t.Fatal("unexpected run time: ", int64(used/time.Millisecond))
	}
}
