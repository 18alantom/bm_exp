package main

import (
	"test/bm_poc/bm"
)

func main() {
	run()
}

func run() {
	bm := bm.BM{
		Target: "/Users/alan/Desktop/code/test_go/bm_poc/bench",
		Cache:  "/Users/alan/Desktop/code/test_go/bm_poc/.cache",
		Config: GetBenchConfig(),
	}
	bm.SetupBench()
}
