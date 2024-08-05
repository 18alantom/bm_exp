package bm

import (
	"fmt"
	"sync"
	"time"
)

type BM struct {
	Target string
	Cache  string
	Config Config
}

func (bm *BM) SetupBench() {
	fmt.Println("\x1b[34;1mSetting up bench\x1b[m")
	start := time.Now()
	outs := getOuts(bm.Config.Apps)

	err_strs := make([]string, 0)
	err_chan := make(chan string, len(bm.Config.Apps))

	wg := sync.WaitGroup{}
	exec := Exec{Ctx: bm.context()}

	wg.Add(1)
	go func() {
		for err := range err_chan {
			err_strs = append(err_strs, err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		exec.Execute(bm.Config.Apps, outs, err_chan, true)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		merge(outs)
		wg.Done()
	}()

	wg.Wait()
	bm.wrapUp(err_strs, start)
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

func (bm *BM) context() Context {
	return Context{
		Target: bm.Target,
		Cache:  bm.Cache,
	}
}
