package ctxroutines

// WithPreRun creates a Runner that calls cb before executing r.Run()
func WithPreRun(r Runner, cb func()) Runner {
	return FromRunner(r, func() error {
		cb()
		return r.Run()
	})
}

// WithPostRun creates a Runner that calls cb after executing r.Run()
func WithPostRun(r Runner, cb func(error)) Runner {
	return FromRunner(r, func() error {
		err := r.Run()
		cb(err)
		return err
	})
}

// WithPreCancel creates a Runner that calls cb before executing r.Cancel()
func WithPreCancel(r Runner, cb func()) Runner {
	return NewRunner(r.Context(), func() {
		cb()
		r.Cancel()
	}, r.Run)
}

// WithPostCancel creates a Runner that calls cb after executing r.Cancel()
func WithPostCancel(r Runner, cb func()) Runner {
	return NewRunner(r.Context(), func() {
		r.Cancel()
		cb()
	}, r.Run)
}
