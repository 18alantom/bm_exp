package bm

import (
	"fmt"
	"sync"
	"time"
)

type Exec struct {
	Ctx Context

	apps []App
	outs []Out
	stop *Stop

	err_chan  chan string
	time_chan chan TimeTuple
}

type ActionTuple struct {
	action Action
	stage  Stage
}

// TODO:
// - Stage start and stage end message (whether stage failed)
// - Put shell output message ($ ...) in shell.run itself
// - Cleanup on error
// - Timestamp all outputs
// - App is not erroring out soon enough
// - Timeouts

func (exec *Exec) Execute(
	apps []App, outs []Out,
	err_chan chan string, time_chan chan TimeTuple,
	concurrently bool,
) {
	exec.apps = apps
	exec.outs = outs
	exec.stop = NewStop()
	exec.err_chan = err_chan
	exec.time_chan = time_chan

	defer close(exec.time_chan)
	defer close(exec.err_chan)
	defer exec.done()
	defer exec.stop.stop()

	benchOut := outs[len(outs)-1]

	start := time.Now()
	if err := exec.initBench(benchOut.Output); err != nil {
		benchOut.Output <- doneErrorOutput(err, time.Since(start), InitBench)
		exec.err_chan <- fmt.Sprintf("%s :: %s", "bench", err.Error())
		return
	}
	benchOut.Output <- doneOutput(time.Since(start), InitBench)

	exec.executeActions(concurrently)
}

func (exec *Exec) executeActions(concurrently bool) {
	concurrentActions := []ActionTuple{
		{fetchRepo, FetchRepo},
		{validate, Validate},
		{installJS, InstallJS},
		{buildFrontend, BuildFrontend},
	}

	if concurrently {
		exec.concurrentSequence(concurrentActions)
	} else {
		exec.sequentialSequence(concurrentActions)
	}

	sequentialActions := []ActionTuple{
		{installPy, InstallPy},
		{completed, Completed},
	}
	exec.sequentialSequence(sequentialActions)
}

func (exec *Exec) concurrentSequence(actions []ActionTuple) {
	var wg sync.WaitGroup
	wg.Add(len(exec.apps))

	runSequential := func(app App, out Out, actions []ActionTuple) {
		exec.sequential(app, out, actions)
		wg.Done()
	}

	for i, app := range exec.apps {
		go runSequential(app, exec.outs[i], actions)
	}
	wg.Wait()
}

func (exec *Exec) sequentialSequence(actions []ActionTuple) error {
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

func (exec *Exec) sequential(app App, out Out, actions []ActionTuple) error {
	for _, t := range actions {
		if exec.stop.Stopped() {
			stopped(exec.Ctx, Stopped, app, out)
			return nil
		}

		start := time.Now()
		err := t.action(exec.Ctx, t.stage, app, out)
		end := time.Since(start)
		exec.time_chan <- TimeTuple{app.Name(), t.stage, end}

		if err != nil {
			out.Output <- doneErrorOutput(err, end, t.stage)
			exec.stop.stop()
			exec.err_chan <- fmt.Sprintf("%s :: %s", app.Name(), err.Error())
			return err
		}

		out.Output <- doneOutput(end, t.stage)
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

func doneOutput(end time.Duration, stage Stage) Output {
	return Output{
		Data:  fmt.Sprintf("Done (%.3fs)", end.Seconds()),
		Stage: stage,
	}
}

func doneErrorOutput(err error, end time.Duration, stage Stage) Output {
	return Output{
		Data:  fmt.Sprintf("Error: %s (%.3fs)", err.Error(), end.Seconds()),
		Stage: stage,
	}
}
