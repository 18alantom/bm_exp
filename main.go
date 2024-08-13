package main

import (
	"os"
	"test/bm_poc/bm"
)

// TODO:
// --no-cache [no-cache option]

func main() {
	run()
}

func run() {
	args := getArgs()
	maker := bm.BM{Config: getBenchConfig(args)}
	ctx := bm.Context{
		NoCache:    maker.Config.Args.NoCache,
		Sequential: maker.Config.Args.Sequential,
		Target:     "./temp/bench",
		Cache:      "./temp/.cache",
	}
	maker.SetupBench(ctx)
}

func getBenchConfig(args bm.Args) bm.Config {
	apps := []bm.App{
		{User: "frappe", Repo: "frappe"},
	}

	for _, app := range args.Apps {
		if app == "frappe" {
			continue
		}

		apps = append(apps, bm.App{User: "frappe", Repo: app})
	}

	return bm.Config{Apps: apps, Args: args}
}

func getArgs() bm.Args {
	opt := bm.Args{Sequential: false, NoCache: false, Apps: []string{}}

	inApps := false
	for _, arg := range os.Args {
		if arg == "--no-cache" {
			opt.NoCache = true
		} else if arg == "--seq" {
			opt.Sequential = true
		} else if arg == "--apps" {
			inApps = true
			continue
		} else if inApps && len(arg) >= 2 && arg[:2] == "--" {
			inApps = false
		} else if inApps {
			opt.Apps = append(opt.Apps, arg)
		}
	}

	return opt
}
