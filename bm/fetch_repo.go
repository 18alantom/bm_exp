package bm

import (
	"fmt"
)

func fetchRepo(ctx Context, app App, out Out) error {
	// TODO:
	// - [ ] check if repo is present in cache
	// - [ ] cache repo after cloning
	// - [ ] fix git output capture
	// - [ ] clone repo to correct destination

	out.Output <- Output{
		Data:  fmt.Sprintf("Fetching %s", app.Name()),
		Stage: FetchRepo,
	}

	// shell := Shell{Out: out.Output, Stage: FetchRepo}
	// url := fmt.Sprintf("https://github.com/%s/%s", app.User, app.Repo)
	// command := fmt.Sprintf("git clone %s --depth 1", url)

	// shell := Shell{Out: out.Output, Stage: FetchRepo}
	// command := fmt.Sprintf("git clone %s --depth 1", app.Repo)
	// return shell.Run(command)
	return nil
}
