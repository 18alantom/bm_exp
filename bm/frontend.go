package bm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

func installJS(ctx Context, stage Stage, app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("Installing JS Dependencies for %s", app.Name()),
		Stage: stage,
	}

	appPath := GetAppPath(ctx, app)

	_, err := readPackageJSON(appPath)
	if errors.Is(err, fs.ErrNotExist) {
		// App doesn't have a frontend
		return nil
	}

	command := fmt.Sprintf("yarn --cwd %s install", appPath)
	return Shell{Output: out.Output, Stage: stage}.Run(command)
}

func buildFrontend(ctx Context, stage Stage, app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("Building Frontend for %s", app.Name()),
		Stage: stage,
	}

	appPath := GetAppPath(ctx, app)

	pj, err := readPackageJSON(appPath)
	if errors.Is(err, fs.ErrNotExist) || len(pj.Scripts.Build) == 0 {
		// App doesn't have a frontend or doesn't need building
		return nil
	}

	if err != nil {
		return err
	}

	command := fmt.Sprintf("yarn --cwd %s build", appPath)
	return Shell{Output: out.Output, Stage: stage}.Run(command)
}

func readPackageJSON(appPath string) (PackageJSON, error) {
	pjPath := path.Join(appPath, "package.json")
	pj := PackageJSON{}

	data, err := os.ReadFile(pjPath)
	if err != nil {
		return pj, err
	}

	json.Unmarshal(data, &pj)
	return pj, nil
}
