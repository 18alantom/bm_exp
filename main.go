package main

import (
	"fmt"
	"test/bm_poc/bm"
)

func main() {
	run()
}

func run() {
	fmt.Println("\x1b[32;1mRunning BM\x1b[m")
	bm := bm.BM{
		Target: "/Users/alan/Desktop/code/test_go/bm_poc/bench",
		Cache:  "/Users/alan/Desktop/code/test_go/bm_poc/.cache",
		Config: GetBenchConfig(),
	}
	bm.SetupBench()
}
