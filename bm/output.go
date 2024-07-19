package bm

import (
	"fmt"
	"sync"
)

type ANSIColor = string

type Output struct {
	Data  string
	Stage Stage
}

type Out struct {
	Output chan Output
	Done   chan struct{}
	App    string
}

const (
	Red     ANSIColor = "\x1b[31m"
	Green   ANSIColor = "\x1b[32m"
	Yellow  ANSIColor = "\x1b[33m"
	Blue    ANSIColor = "\x1b[34m"
	Magenta ANSIColor = "\x1b[35m"
	Cyan    ANSIColor = "\x1b[36m"
	White   ANSIColor = "\x1b[37m"
	Purple  ANSIColor = "\x1b[38;5;177m"
	Orange  ANSIColor = "\x1b[38;5;214m"
	Salmon  ANSIColor = "\x1b[38;5;223m"
)

func getOuts(apps []App) []Out {
	outs := make([]Out, len(apps))
	for i, app := range apps {
		outs[i] = Out{
			Output: make(chan Output),
			Done:   make(chan struct{}),
			App:    fmt.Sprintf("%s/%s", app.User, app.Repo),
		}
	}

	return outs
}

func merge(outs []Out) {
	mux := make(chan string, 1024)
	var wgOutput sync.WaitGroup
	wgOutput.Add(len(outs))

	colorMap := map[Stage]ANSIColor{
		FetchRepo:     Salmon,
		Validate:      Orange,
		InstallJS:     Purple,
		BuildFrontend: Magenta,
		InstallPy:     Cyan,
		Completed:      Green,
		Stopped:       Yellow,
		Errored:       Red,
	}

	getOutput := func(out Out) {
		for {
			select {
			case output := <-out.Output:
				mux <- fmt.Sprintf(
					"%s%-26s\x1b[m \x1b[33m%-16s\x1b[m %s\r\n",
					colorMap[output.Stage],
					output.Stage,
					out.App, output.Data,
				)
			case <-out.Done:
				wgOutput.Done()
				return
			}
		}
	}

	for _, out := range outs {
		go getOutput(out)
	}

	var wgMux sync.WaitGroup
	wgMux.Add(1)
	go func() {
		for output := range mux {
			fmt.Print(output)
		}
		wgMux.Done()
	}()

	wgOutput.Wait()
	close(mux)
	wgMux.Wait()
}
