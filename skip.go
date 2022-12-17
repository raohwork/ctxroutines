package ctxroutines

// Skip creates a Runner that runs every Runner of rs in separated goroutine,
// returns first result and cancels others.
func Skip(rs ...Runner) Runner {
	c := CancelAll(rs...)
	return FuncRunner(c, func() error {
		ch := make(chan error, 1)

		for _, r := range rs {
			go func(r Runner) {
				ch <- r.Run()
			}(r)
		}

		ret := <-ch
		c()
		for range rs[1:] {
			<-ch
		}
		return ret
	})
}
