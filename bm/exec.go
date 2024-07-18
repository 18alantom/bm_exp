package bm

import "sync"

func (bm *BM) executeActions(outs []Out, concurrently bool) {
	concurrentActions := []Action{
		fetchRepo,
		validate,
		installJS,
		buildFrontend,
	}

	if concurrently {
		concurrentSequence(bm.Config.Apps, outs, concurrentActions)
	} else {
		sequentialSequence(bm.Config.Apps, outs, concurrentActions)
	}

	sequentialActions := []Action{
		installPy,
		complete,
	}
	sequentialSequence(bm.Config.Apps, outs, sequentialActions)
}

func concurrentSequence(apps []App, outs []Out, actions []Action) {
	var wg sync.WaitGroup
	wg.Add(len(apps))

	runSequential := func(app App, out Out, actions []Action) {
		sequential(app, out, actions)
		wg.Done()
	}

	for i, app := range apps {
		go runSequential(app, outs[i], actions)
	}
	wg.Wait()
}

func sequentialSequence(apps []App, outs []Out, actions []Action) {
	for i, app := range apps {
		sequential(app, outs[i], actions)
	}
}

func sequential(app App, out Out, actions []Action) {
	for _, action := range actions {
		action(app, out)
	}
}
