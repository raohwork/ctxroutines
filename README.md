some helpers to write common routines

[![GoDoc](https://godoc.org/github.com/raohwork/ctxroutines?status.svg)](https://godoc.org/github.com/raohwork/ctxroutines)
[![Go Report Card](https://goreportcard.com/badge/github.com/raohwork/ctxroutines)](https://goreportcard.com/report/github.com/raohwork/ctxroutines)

# Work in progress

This library is still WORK IN PROGRESS. Codes here are used in several small projects in production for few months, should be safe I think. But it still need some refinement like writting docs, better test cases and benchmarks.

# Graceful shutdown made easy

```go
func myCrawler(ctx context.Context) (err error) {
	req, err := http.NewRequestWithContext(
		ctx, "GET", "https://google.com", nil,
	)
	if err != nil {
		// log the error and
		return err
	}

	data, err := grab(req)
	if err != nil {
		// log the error and
		return
	}

	err = saveDB(ctx, data)
	// log the error and
	return
}

func main() {
	r := NewStatefulRunner(
		Loop( // infinite loop
			RunAtLeast( // do not run crawler too fast
				10*time.Second,
				CTXRunner(myCrawler),
			)))

	// listen to signal and start main prog
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go r.Run()

	// wait signal
	<-ch

    // got signal, clear up
	signal.Reset(os.Interrupt, os.Kill)
	close(ch)

	// cancel job
	r.Cancel()

	// waits for shutdown
	r.Lock()()
}
```

# Race conditions

Codes in this library are thread-safe unless specified. However, thread-safety of external function is not covered.

Considering this example:

```go
f := Loop(FuncRunner(cancel, yourFunc))
```

`f` is thread-safe iff `cancel` and `yourFunc` are thread-safe.

# License

Copyright Chung-Ping Jen <ronmi.ren@gmail.com> 2021-

MPL v2.0
