package bm

import (
	"sync"
	"time"
)

// TODO:
// - [ ] On Error out bench setup failed message
// - [ ] Aggregate errors and return to the calling process

func Execute(apps []App, outs []Out, concurrently bool) {
	stop := NewStop()

	defer done(outs)
	defer stop.stop()

	// TODO: Remove (used to check termination)
	go func() {
		time.Sleep(500 * time.Millisecond)
		stop.stop()
	}()

	concurrentActions := []Action{
		fetchRepo,
		validate,
		installJS,
		buildFrontend,
	}

	if concurrently {
		concurrentSequence(apps, outs, stop, concurrentActions)
	} else {
		sequentialSequence(apps, outs, stop, concurrentActions)
	}

	sequentialActions := []Action{
		installPy,
		completed,
	}
	sequentialSequence(apps, outs, stop, sequentialActions)
}

func concurrentSequence(apps []App, outs []Out, stop *Stop, actions []Action) {
	var wg sync.WaitGroup
	wg.Add(len(apps))

	runSequential := func(app App, out Out, actions []Action) {
		sequential(app, out, stop, actions)
		wg.Done()
	}

	for i, app := range apps {
		go runSequential(app, outs[i], actions)
	}
	wg.Wait()
}

func sequentialSequence(apps []App, outs []Out, stop *Stop, actions []Action) error {
	for i, app := range apps {
		out := outs[i]

		if err := sequential(app, out, stop, actions); err != nil {
			return err
		}

		if stop.Stopped() {
			return nil
		}
	}

	return nil
}

func sequential(app App, out Out, stop *Stop, actions []Action) error {
	for _, action := range actions {
		if stop.Stopped() {
			stopped(app, out)
			return nil
		}

		if err := action(app, out); err != nil {
			errored(app, out)
			stop.stop()
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
