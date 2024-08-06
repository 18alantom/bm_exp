package bm

import (
	"fmt"
	"math/rand"
	"path"
	"time"
)


type Config struct {
	Apps []App
}

type App struct {
	User   string
	Repo   string
	Branch string
}

func (app *App) Name() string {
	return fmt.Sprintf("%s/%s", app.User, app.Repo)
}


func RandSleep(max float64) {
	duration := time.Duration(max*rand.Float64()) * time.Millisecond
	time.Sleep(duration)
}

func getTargetPath(ctx Context, app App) string {
	return path.Join(ctx.Target, "apps", app.Repo)
}
