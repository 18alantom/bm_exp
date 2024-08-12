package bm

import (
	"fmt"
	"path"
)

// Each Stage is associated with an Action
type Stage string
type Action func(ctx Context, stage Stage, app App, out Out) error

const (
	Bench         Stage = "Bench"
	FetchRepo     Stage = "FetchRepository"
	Validate      Stage = "Validate"
	InstallJS     Stage = "InstallJSDependencies"
	BuildFrontend Stage = "BuildFrontend"
	InstallPy     Stage = "InstallPythonDependencies"
	Completed     Stage = "Completed"
	Stopped       Stage = "Stopped"
)

func validate(ctx Context, stage Stage, app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("Validating %s", app.Name()),
		Stage: stage,
	}
	return nil
}

func installPy(ctx Context, stage Stage, app App, out Out) error {
	appPath := GetAppPath(ctx, app)
	command := fmt.Sprintf("python -m pip install --upgrade -e %s", appPath)
	if ctx.NoCache {
		command = fmt.Sprintf("%s --no-cache-dir", command)
	}
	shell := NewShell(out.Output, stage)

	cacheFolder := path.Join(ctx.Cache, "pip")
	cacheFolderEnv := fmt.Sprintf("PIP_CACHE_FOLDER=%s", cacheFolder)
	shell.AppendEnv(cacheFolderEnv)

	return shell.Run(command)
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
