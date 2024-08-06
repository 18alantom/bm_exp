package bm

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

// TODO:
// - handle non built apps
// - handle apps that dont have a frontend

func installJS(ctx Context, app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("Installing JS Dependencies for %s", app.Name()),
		Stage: InstallJS,
	}

	targetPath := getTargetPath(ctx, app)
	exists, err := ensureTarget(targetPath, "package.json")
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	command := fmt.Sprintf("yarn --cwd %s install", targetPath)
	out.Output <- Output{
		Data:  fmt.Sprintf("$ %s", command),
		Stage: InstallJS,
	}

	return Shell{Out: out.Output, Stage: InstallJS}.Run(command)
}

func buildFrontend(ctx Context, app App, out Out) error {
	out.Output <- Output{
		Data:  fmt.Sprintf("Building Frontend for %s", app.Name()),
		Stage: BuildFrontend,
	}

	targetPath := getTargetPath(ctx, app)
	exists, err := ensureTarget(targetPath, "package.json")
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	command := fmt.Sprintf("yarn --cwd %s build", targetPath)
	out.Output <- Output{
		Data:  fmt.Sprintf("$ %s", command),
		Stage: BuildFrontend,
	}

	return Shell{Out: out.Output, Stage: BuildFrontend}.Run(command)
}

func ensureTarget(targetPath string, filename string) (bool, error) {
	s, err := os.Stat(
		path.Join(targetPath, filename),
	)

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if s.Mode().IsRegular() {
		return true, nil
	}

	return false, errors.New("not regular file")
}
