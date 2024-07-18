package bm

import (
	"fmt"
)

func fetchRepo(app App, out Out) {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Fetching %s", app.Name()),
		Stage: FetchRepo,
	}
}

func validate(app App, out Out) {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Validating %s", app.Name()),
		Stage: Validate,
	}
}

func installJS(app App, out Out) {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Installing JS dependencies %s", app.Name()),
		Stage: InstallJS,
	}
}

func buildFrontend(app App, out Out) {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Building frontend %s", app.Name()),
		Stage: BuildFrontend,
	}
}

func installPy(app App, out Out) {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Installing Python dependencies %s", app.Name()),
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
