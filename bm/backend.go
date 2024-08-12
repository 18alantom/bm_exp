package bm

import (
	"fmt"
	"path"
)

func installPy(ctx Context, stage Stage, app App, out Out) error {
	appPath := GetAppPath(ctx, app)

	shell := NewShell(out.Output, stage)
	shell.AppendEnv(getPipCacheFolderEnv(ctx))

	command := fmt.Sprintf("python -m pip install --upgrade -e %s", appPath)
	if ctx.NoCache {
		command = fmt.Sprintf("%s --no-cache-dir", command)
	}

	return shell.Run(command)
}

func getPipCacheFolderEnv(ctx Context) string {
	cacheFolder := path.Join(ctx.Cache, "pip")
	if ctx.NoCache {
		cacheFolder = path.Join("/tmp", "junk", "pip")
	}

	return fmt.Sprintf("PIP_CACHE_DIR=%s", cacheFolder)
}
