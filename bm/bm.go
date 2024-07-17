package bm

import (
	"fmt"
	"sync"
)

type BM struct {
	Target string
	Cache  string
	Config Config
}

type Output struct {
	Data  string
	Stage Stage
}

// Used for multiplexing output
type Out struct {
	Output chan Output
	Done   chan struct{}
	App    string
}

func (bm *BM) SetupBench() {
	outs := bm.getOuts()
	defer bm.merge(outs)

	// Not structured concurrency
	bm.getRepos(outs)
}

func (bm *BM) getRepos(outs []Out) {
	for i, app := range bm.Config.Apps {
		go fetchRepo(app, outs[i])
	}
}

func (bm *BM) merge(outs []Out) {
	mux := make(chan string, 1024)

	var wg sync.WaitGroup
	wg.Add(len(outs))

	for _, out := range outs {
		go func() {
			for {
				select {
				case output := <-out.Output:
					mux <- fmt.Sprintf("\x1b[33m%s\x1b[m(\x1b[34m%s\x1b[m) :: %s\r\n", output.Stage, out.App, output.Data)
				case <-out.Done:
					wg.Done()
					return
				}
			}
		}()
	}

	go func() {
		for output := range mux {
			fmt.Print(output)
		}
	}()

	wg.Wait()
	close(mux)
}

// func (bm *BM) installJS(app App) {}

// func (bm *BM) buildJS(app App) {}

// func (bm *BM) installPy(app App) {}

// func mux() {}

func (bm *BM) getOuts() []Out {
	outs := make([]Out, len(bm.Config.Apps))
	for i, app := range bm.Config.Apps {
		outs[i] = Out{
			Output: make(chan Output),
			Done:   make(chan struct{}),
			App:    fmt.Sprintf("%s/%s", app.User, app.Repo),
		}
	}

	return outs
}
