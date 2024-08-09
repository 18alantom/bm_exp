package bm

import (
	"fmt"
)

// Each Stage is associated with an Action
type Stage string
type Context struct {
	Target string
	Cache  string
}
type Action func(ctx Context, stage Stage, app App, out Out) error

const (
	InitBench     Stage = "InitBench"
	FetchRepo     Stage = "FetchRepository"
	Validate      Stage = "Validate"
	InstallJS     Stage = "InstallJSDependencies"
	BuildFrontend Stage = "BuildFrontend"
	InstallPy     Stage = "InstallPythonDependencies"
	Completed     Stage = "Completed"
	Stopped       Stage = "Stopped"
)

func validate(ctx Context, stage Stage, app App, out Out) error {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Validating %s", app.Name()),
		Stage: stage,
	}

	return nil
}

func installPy(ctx Context, stage Stage, app App, out Out) error {
	RandSleep(1000)
	out.Output <- Output{
		Data:  fmt.Sprintf("Installing Python dependencies %s", app.Name()),
		Stage: stage,
	}

	return nil
}

func completed(ctx Context, stage Stage, app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("Installation Completed %s", app.Name()),
		Stage: stage,
	}

	return nil
}

func stopped(_ Context, stage Stage, app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("App installation stopped %s", app.Name()),
		Stage: stage,
	}

	return nil
}
