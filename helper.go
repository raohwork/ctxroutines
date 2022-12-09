// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import (
	"os"
	"os/signal"
)

// CancelOnSigal runs r and calls r.Cancel when receiving first signal in sig
func CancelOnSignal(r Runner, sig ...os.Signal) (err error) {
	ch := make(chan os.Signal)
	signal.Notify(ch, sig...)
	go func() {
		<-ch
		signal.Reset(sig...)
		close(ch)
		r.Cancel()
		for range ch {
		}
	}()

	return r.Run()
}
