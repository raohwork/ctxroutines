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

func OnceWithin(dur time.Duration, f Runner) (ret Runner) {
	return newOnceWithin(
		dur,
		NewStatefulRunner(RunAtLeast(dur, f)),
	)
}

func OnceSuccessWithin(dur time.Duration, f Runner) (ret Runner) {
	return newOnceWithin(
		dur,
		NewStatefulRunner(RunAtLeastSuccess(dur, f)),
	)
}

func OnceFailedWithin(dur time.Duration, f Runner) (ret Runner) {
	return newOnceWithin(
		dur,
		NewStatefulRunner(RunAtLeastFailed(dur, f)),
	)
}
