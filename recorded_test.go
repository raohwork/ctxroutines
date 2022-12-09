// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import "testing"

func TestRecorded(t *testing.T) {
	ch := make(chan uint64, 1)
	defer close(ch)

	f := Counter(func(n uint64) error { ch <- n; return nil })
	r := Recorded(f)
	err := r.Run()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if l := len(ch); l != 1 {
		t.Fatal("expected 1 items in chan, got ", l)
	}
	if u := <-ch; u != 0 {
		t.Fatal("expected '0' in queue, got ", u)
	}
	if l := r.Count(); l != 1 {
		t.Fatal("expected ran 1 time, got ", l)
	}

	err = r.Run()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if l := len(ch); l != 1 {
		t.Fatal("expected 1 items in chan, got ", l)
	}
	if u := <-ch; u != 1 {
		t.Fatal("expected '1' in queue, got ", u)
	}
	if l := r.Count(); l != 2 {
		t.Fatal("expected ran 2 times, got ", l)
	}
}
