package bm

import (
	"fmt"
)

// Each Stage is associated with an Action
type Stage string
type Action func(app App, out Out) error

const (
	FetchRepo     Stage = "FetchRepository"
	Validate      Stage = "Validate"
	InstallJS     Stage = "InstallJSDependencies"
	BuildFrontend Stage = "BuildFrontend"
	InstallPy     Stage = "InstallPythonDependencies"
	Completed      Stage = "Completed"
	Stopped       Stage = "Stopped"
	Errored       Stage = "Errored"
)

func fetchRepo(app App, out Out) error {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Fetching %s", app.Name()),
		Stage: FetchRepo,
	}

	return nil
}

func validate(app App, out Out) error {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Validating %s", app.Name()),
		Stage: Validate,
	}

	return nil
}

func installJS(app App, out Out) error {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Installing JS dependencies %s", app.Name()),
		Stage: InstallJS,
	}

	return nil
}

func buildFrontend(app App, out Out) error {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Building frontend %s", app.Name()),
		Stage: BuildFrontend,
	}

	return nil
}

func installPy(app App, out Out) error {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Installing Python dependencies %s", app.Name()),
		Stage: InstallPy,
	}

	return nil
}

func completed(app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("Installation Completed %s", app.Name()),
		Stage: Completed,
	}
	return nil
}

func errored(app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("App errored out %s", app.Name()),
		Stage: Errored,
	}
	return nil
}

func stopped(app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("App installation stopped %s", app.Name()),
		Stage: Stopped,
	}
	return nil
}
