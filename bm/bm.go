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
	var wg sync.WaitGroup
	defer bm.wrapUp(&wg, time.Now())

	fmt.Println("\x1b[34;1mSetting up bench\x1b[m")
	outs := getOuts(bm.Config.Apps)

	wg.Add(1)
	go func() {
		bm.executeActions(outs, true)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		merge(outs)
		wg.Done()
	}()

}

func (bm *BM) wrapUp(wg *sync.WaitGroup, start time.Time) {
	wg.Wait()
	fmt.Println("\x1b[32;1mBench setup completed\x1b[m")
	fmt.Printf("%d apps installed in %.3fs\n", len(bm.Config.Apps), time.Since(start).Seconds())
}
