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

	errStrs := make([]string, 0)
	errChan := make(chan string, len(bm.Config.Apps))

	timeMap := make(TimeMap)
	timeChan := make(chan TimeTuple, len(bm.Config.Apps))

	wg := sync.WaitGroup{}
	exec := Exec{Ctx: ctx}

	// Handle err chan messages
	wg.Add(1)
	go func() {
		for err := range errChan {
			errStrs = append(errStrs, err)
		}
		wg.Done()
	}()

	// Handle time chan messages
	wg.Add(1)
	go func() {
		for t := range timeChan {
			m, ok := timeMap[t.App]
			if !ok {
				m = make(map[Stage]time.Duration)
				timeMap[t.App] = m
			}

			m[t.Stage] = t.Duration
		}
		wg.Done()
	}()

	// Handle execute
	wg.Add(1)
	go func() {
		exec.Execute(bm.Config.Apps, outs, errChan, timeChan)
		wg.Done()
	}()

	// Handle output merging
	wg.Add(1)
	go func() {
		merge(outs)
		wg.Done()
	}()

	wg.Wait()

	end := time.Since(start).Seconds()
	bm.wrapUp(errStrs, end)
	printTimeBreakdown(timeMap, end)
}

func (bm *BM) wrapUp(errs []string, end float64) {
	if len(errs) > 0 {
		fmt.Println("\x1b[31;1mBench setup failed\x1b[m")
	} else {
		fmt.Println("\x1b[32;1mBench setup succeeded\x1b[m")
	}

	fmt.Printf("\nWall time taken: %.3fs\n", end)

	if len(errs) > 0 {
		fmt.Println("Errors:")
		for _, err := range errs {
			fmt.Printf("- %s\n", err)
		}
	}
}

func printTimeBreakdown(timeMap TimeMap, wallTime float64) {
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

	concTotal := 0.0 // time taken executing concurrent steps, less than wall time
	seqTotal := 0.0  // time taken executing sequential steps, same as wall time
	appTotal := 0.0

	// Print header
	fmt.Printf("\n\nTime Breakdown:\n| %-16s ", "org/repo")
	for _, s := range seq {
		fmt.Printf("| %9s ", s.string)
	}
	fmt.Printf("| %9s |\n", "total")

	// Print data
	for key, val := range timeMap {
		if key == "bench" {
			continue
		}

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

			if In(s.Stage, FetchRepo, Validate, InstallJS, BuildFrontend) {
				concTotal += sec
			} else {
				seqTotal += sec
			}
			appTotal += sec
		}

		fmt.Printf("| %8.3fs |\n", total)
	}

	benchDur := timeMap["bench"][Bench].Seconds()
	total := appTotal + benchDur

	fmt.Printf("\nTotals:\n")
	fmt.Printf("Bench init            : %8.3fs\n", benchDur)
	fmt.Printf("Concurrent app stages : %8.3fs\n", concTotal)
	fmt.Printf("Sequential app stages : %8.3fs\n", seqTotal)
	fmt.Printf("---------------------------------\n")
	fmt.Printf("Total app             : %8.3fs\n", appTotal)
	fmt.Printf("Total app + bench     : %8.3fs\n", appTotal+benchDur)
	fmt.Printf("---------------------------------\n")
	fmt.Printf("Total wall time       : %8.3fs\n", wallTime)
	fmt.Printf("Time saved            : %8.3fs\n", total-wallTime)
}

func In(s Stage, values ...Stage) bool {
	for _, v := range values {
		if v == s {
			return true
		}
	}
	return false
}
