package bm

import (
	"fmt"
	"os"
	"path"
	"test/bm_poc/utils"
)

func fetchRepo(ctx Context, stage Stage, app App, out Out) error {
	// TODO:
	// - fix git output capture
	// - repo uid, guid?

	out.Output <- Output{
		Data:  fmt.Sprintf("Fetching %s", app.Name()),
		Stage: stage,
	}
	cachePath := getCachePath(ctx, app)
	targetPath := GetAppPath(ctx, app)

	// Ensure repo exists in the cache folder
	var err error = nil
	if !hasCache(cachePath) {
		if err = cloneRepo(stage, app, out, cachePath); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	// Make sure repo is at the target
	if err = initializeTarget(cachePath, targetPath); err != nil {
		return err
	}

	return nil
}

func cloneRepo(stage Stage, app App, out Out, cachePath string) error {
	url := fmt.Sprintf("https://github.com/%s/%s", app.User, app.Repo)

	// TODO: clone to cache or clone to target then copy to cache?
	command := fmt.Sprintf("git clone %s --depth 1 %s", url, cachePath)
	return Shell{Out: out.Output, Stage: stage}.Run(command)
}

func getCachePath(ctx Context, app App) string {
	// TODO: Cache path should be `$base/$user/$repo/${hash|branch|tag}?`
	return path.Join(ctx.Cache, app.User, app.Repo)
}

func hasCache(cachePath string) bool {
	// TODO:
	// - Use tar or something for caching
	// - File locking?

	stat, err := os.Stat(cachePath)
	// TODO: Handle error properly, err doesn't always mean path doesn't exist
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func initializeTarget(cachePath string, targetPath string) error {
	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		return err
	}

	// TODO: run output that dir is being copied to target
	return utils.CopyDir(cachePath, targetPath)
}
