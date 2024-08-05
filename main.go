package main

import (
	"test/bm_poc/bm"
)

func main() {
	run()
}

func run() {
	// TODO:
	// - --no-cache [what]
	maker := bm.BM{Config: GetBenchConfig()}
	ctx := bm.Context{
		Target: "/Users/alan/Desktop/code/test_go/bm_poc/temp/bench",
		Cache:  "/Users/alan/Desktop/code/test_go/bm_poc/temp/.cache",
	}
	maker.SetupBench(ctx)
}
