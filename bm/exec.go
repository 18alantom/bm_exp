package bm

import (
	"fmt"
	"os"
	"sync"
)

type Exec struct {
	Ctx Context

	apps     []App
	outs     []Out
	stop     *Stop
	err_chan chan string
}

// TODO:
// - Measure per stage timing
// - Stage start and stage end message (whether stage failed)
// - Put shell output message ($ ...) in shell.run itself
// - Cleanup on error
// - Timestamp all outputs
// - App is not erroring out soon enough

func (exec *Exec) Execute(apps []App, outs []Out, err_chan chan string, concurrently bool) {
	exec.apps = apps
	exec.outs = outs
	exec.stop = NewStop()
	exec.err_chan = err_chan

	defer close(exec.err_chan)
	defer exec.done()
	defer exec.stop.stop()

	if err := exec.initBench(); err != nil {
		err_chan <- err.Error()
		return
	}

	exec.executeActions(concurrently)
}

func (exec *Exec) initBench() error {
	// TODO: Should probably be under BM, bm.go
	return os.RemoveAll(exec.Ctx.Target)
}

func (exec *Exec) executeActions(concurrently bool) {
	concurrentActions := []Action{
		fetchRepo,
		validate,
		installJS,
		buildFrontend,
	}

	if concurrently {
		exec.concurrentSequence(concurrentActions)
	} else {
		exec.sequentialSequence(concurrentActions)
	}

	sequentialActions := []Action{
		installPy,
		completed,
	}
	exec.sequentialSequence(sequentialActions)
}

func (exec *Exec) concurrentSequence(actions []Action) {
	var wg sync.WaitGroup
	wg.Add(len(exec.apps))

	runSequential := func(app App, out Out, actions []Action) {
		exec.sequential(app, out, actions)
		wg.Done()
	}

	for i, app := range exec.apps {
		go runSequential(app, exec.outs[i], actions)
	}
	wg.Wait()
}

func (exec *Exec) sequentialSequence(actions []Action) error {
	for i, app := range exec.apps {
		out := exec.outs[i]

		if err := exec.sequential(app, out, actions); err != nil {
			return err
		}

		if exec.stop.Stopped() {
			return nil
		}
	}

	return nil
}

func (exec *Exec) sequential(app App, out Out, actions []Action) error {
	for _, action := range actions {
		if exec.stop.Stopped() {
			stopped(exec.Ctx, app, out)
			return nil
		}

		if err := action(exec.Ctx, app, out); err != nil {
			errored(app, out, err.Error())
			exec.stop.stop()
			exec.err_chan <- fmt.Sprintf("%s :: %s", app.Name(), err.Error())
			return err
		}
	}

	return nil
}

func (exec *Exec) done() {
	for _, out := range exec.outs {
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
