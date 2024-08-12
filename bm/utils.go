package bm

import (
	"fmt"
	"math/rand"
	"path"
	"time"
)

type Config struct {
	Apps []App
	Args Args
}

type Args struct {
	Sequential bool     // run installation sequentially
	NoCache    bool     // skip cache for yarn, pip, etc
	Apps       []string // only frappe/app apps allowed
}

type App struct {
	User string
	Repo string
}

type PackageJSON struct {
	Scripts struct {
		Build string `json:"build"`
	} `json:"scripts"`
}

func (app *App) Name() string {
	return fmt.Sprintf("%s/%s", app.User, app.Repo)
}

func RandSleep(max float64) {
	duration := time.Duration(max*rand.Float64()) * time.Millisecond
	time.Sleep(duration)
}

func GetAppPath(ctx Context, app App) string {
	return path.Join(ctx.Target, "apps", app.Repo)
}
