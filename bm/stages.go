package bm

import (
	"fmt"
)

// TODO: Compose sequential and concurrent stages
// TODO: Debug why not all expected output is being printed

func fetchRepo(app App, out Out) {
	stub := fmt.Sprintf("%s/%s", app.User, app.Repo)

	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Fetching %s", stub),
		Stage: FetchRepo,
	}

	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Done fetching %s", stub),
		Stage: FetchRepo,
	}

	complete(app, out)
}

func installJS(app App, out Out) {
	stub := fmt.Sprintf("%s/%s", app.User, app.Repo)

	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Installing JS dependencies %s", stub),
		Stage: InstallJS,
	}

	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Building frontend %s", stub),
		Stage: BuildFrontend,
	}
}

func installPy(app App, out Out) {
	stub := fmt.Sprintf("%s/%s", app.User, app.Repo)

	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Installing Python dependencies %s", stub),
		Stage: InstallPy,
	}
}

func complete(app App, out Out) {
	stub := fmt.Sprintf("%s/%s", app.User, app.Repo)
	out.Output <- Output{
		Data:  fmt.Sprintf("Installation Complete %s", stub),
		Stage: Complete,
	}

	out.Done <- struct{}{}
	close(out.Output)
	close(out.Done)
}
