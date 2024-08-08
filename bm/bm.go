package bm

import (
	"fmt"
	"sync"
	"time"
)

// TODO:
// - Differentiate output between BM output and subcommand output

type BM struct {
	Config Config
}

type TimeMap map[string]map[Stage]time.Duration
type TimeTuple struct {
	App      string
	Stage    Stage
	Duration time.Duration
}

func (bm *BM) SetupBench(ctx Context) {
	fmt.Println("\x1b[34;1mSetting up bench\x1b[m")
	start := time.Now()
	outs := getOuts(bm.Config.Apps)

	err_strs := make([]string, 0)
	err_chan := make(chan string, len(bm.Config.Apps))

	time_map := make(TimeMap)
	time_chan := make(chan TimeTuple, len(bm.Config.Apps))

	wg := sync.WaitGroup{}
	exec := Exec{Ctx: ctx}

	// Handle err chan messages
	wg.Add(1)
	go func() {
		for err := range err_chan {
			err_strs = append(err_strs, err)
		}
		wg.Done()
	}()

	// Handle time chan messages
	wg.Add(1)
	go func() {
		for t := range time_chan {
			m, ok := time_map[t.App]
			if !ok {
				m = make(map[Stage]time.Duration)
				time_map[t.App] = m
			}

			m[t.Stage] = t.Duration
		}
		wg.Done()
	}()

	// Handle execute
	wg.Add(1)
	go func() {
		exec.Execute(bm.Config.Apps, outs, err_chan, time_chan, true)
		wg.Done()
	}()

	// Handle output merging
	wg.Add(1)
	go func() {
		merge(outs)
		wg.Done()
	}()

	wg.Wait()
	bm.wrapUp(err_strs, start)
	printTimeBreakdown(time_map)
}

func (bm *BM) wrapUp(errs []string, start time.Time) {
	end := time.Since(start).Seconds()
	if len(errs) > 0 {
		fmt.Println("\x1b[31;1mBench setup failed\x1b[m")
	} else {
		fmt.Println("\x1b[32;1mBench setup succeeded\x1b[m")
	}

	fmt.Printf("\nTime taken: %.3fs\n", end)

	if len(errs) > 0 {
		fmt.Println("Errors:")
		for _, err := range errs {
			fmt.Printf("- %s\n", err)
		}
	}
}

func printTimeBreakdown(time_map TimeMap) {
	seq := []struct {
		Stage
		string
	}{
		{FetchRepo, "clone"},
		{Validate, "validate"},
		{InstallJS, "ins js"},
		{BuildFrontend, "build"},
		{InstallPy, "ins py"},
		{Completed, "complete"},
		{Stopped, "stop"},
	}
	grand_total := 0.0

	// Print header
	fmt.Printf("\n\nTime Breakdown\n| %-16s ", "org/repo")
	for _, s := range seq {
		fmt.Printf("| %9s ", s.string)
	}
	fmt.Printf("| %9s |\n", "total")

	// Print data
	for key, val := range time_map {
		total := 0.0
		fmt.Printf("| %-16s ", key)
		for _, s := range seq {
			dur, ok := val[s.Stage]
			sec := 0.0
			if ok {
				sec = dur.Seconds()
			}

			total += sec
			fmt.Printf("| %8.3fs ", sec)
		}

		fmt.Printf("| %8.3fs |\n", total)
		grand_total += total
	}

	fmt.Printf("Grand total: %0.3fs\n", grand_total)
}
