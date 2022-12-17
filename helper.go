// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

// ErrSignalReceived indicates the signal that SignalRunner received
type ErrSignalReceived struct {
	os.Signal
}

func (s ErrSignalReceived) Error() string {
	return "signal received: " + s.String()
}

// SignalRunner creates a Runner that waits first signal in sig and returns it as error
func SignalRunner(sig ...os.Signal) Runner {
	ch := make(chan os.Signal)
	signal.Notify(ch, sig...)
	once := sync.Once{}
	f := func() {
		signal.Reset(sig...)
		close(ch)
		for range ch {
		}
	}
	return FuncRunner(func() {
		once.Do(f)
	}, func() error {
		s := <-ch
		once.Do(f)
		if s == nil {
			return context.Canceled
		}

		return ErrSignalReceived{
			Signal: s,
		}
	})
}

// CancelOnSignal runs r and calls r.Cancel when receiving first signal in sig
func CancelOnSignal(r Runner, sig ...os.Signal) (err error) {
	x := Skip(r, SignalRunner(sig...))
	defer x.Cancel()

	return x.Run()
}
