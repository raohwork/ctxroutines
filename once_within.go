// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ctxroutines

import "time"

func newOnceWithin(dur time.Duration, r StatefulRunner) (ret Runner) {
	return FuncRunner(r.Cancel, func() error {
		err, ran := r.TryRun()
		if err == nil && !ran {
			return nil
		}
		return err
	})
}

// OnceWithin ensures f is not run more than once within duration dur
//
// Additional calls are just ignored. Say you have an empty Runner f
//
//     r := OnceWithin(time.Second, f)
//     r.Run() // runs f
//     r.Run() // skipped
//     time.Sleep(time.Second)
//     r.Run() // runs f
func OnceWithin(dur time.Duration, f Runner) (ret Runner) {
	return newOnceWithin(
		dur,
		NewStatefulRunner(RunAtLeast(dur, f)),
	)
}

// OnceSuccessWithin is like OnceWithin, but only successful call counts.
func OnceSuccessWithin(dur time.Duration, f Runner) (ret Runner) {
	return newOnceWithin(
		dur,
		NewStatefulRunner(RunAtLeastSuccess(dur, f)),
	)
}

// OnceFailedWithin is like OnceWithin, but only failed call counts.
func OnceFailedWithin(dur time.Duration, f Runner) (ret Runner) {
	return newOnceWithin(
		dur,
		NewStatefulRunner(RunAtLeastFailed(dur, f)),
	)
}
