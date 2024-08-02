package bm

import (
	"fmt"
	"sync"
)

// TODO:
// - Measure per stage timing

func Execute(apps []App, outs []Out, err_chan chan string, concurrently bool) {
	stop := NewStop()

	defer close(err_chan)
	defer done(outs)
	defer stop.stop()

	concurrentActions := []Action{
		fetchRepo,
		validate,
		installJS,
		buildFrontend,
	}

	if concurrently {
		concurrentSequence(apps, outs, err_chan, stop, concurrentActions)
	} else {
		sequentialSequence(apps, outs, err_chan, stop, concurrentActions)
	}

	sequentialActions := []Action{
		installPy,
		completed,
	}
	sequentialSequence(apps, outs, err_chan, stop, sequentialActions)
}

func concurrentSequence(apps []App, outs []Out, err_chan chan string, stop *Stop, actions []Action) {
	var wg sync.WaitGroup
	wg.Add(len(apps))

	runSequential := func(app App, out Out, actions []Action) {
		sequential(app, out, err_chan, stop, actions)
		wg.Done()
	}

	for i, app := range apps {
		go runSequential(app, outs[i], actions)
	}
	wg.Wait()
}

func sequentialSequence(apps []App, outs []Out, err_chan chan string, stop *Stop, actions []Action) error {
	for i, app := range apps {
		out := outs[i]

		if err := sequential(app, out, err_chan, stop, actions); err != nil {
			return err
		}

		if stop.Stopped() {
			return nil
		}
	}

	return nil
}

func sequential(app App, out Out, err_chan chan string, stop *Stop, actions []Action) error {
	for _, action := range actions {
		if stop.Stopped() {
			stopped(app, out)
			return nil
		}

		if err := action(app, out); err != nil {
			errored(app, out, err.Error())
			stop.stop()
			err_chan <- fmt.Sprintf("%s :: %s", app.Name(), err.Error())
			return err
		}
	}

	return nil
}

func done(outs []Out) {
	for _, out := range outs {
		out.Done <- struct{}{}
		close(out.Output)
		close(out.Done)
	}
}

// Used to broadcast a stop signal on error
type Stop struct {
	ch   chan struct{}
	stop func()
}

func NewStop() *Stop {
	s := new(Stop)
	s.ch = make(chan struct{})
	s.stop = sync.OnceFunc(func() {
		close(s.ch)
	})
	return s
}

func (s *Stop) Stopped() bool {
	select {
	case <-s.ch:
		return true
	default:
		return false
	}
}
