package bm

import (
	"fmt"
	"math/rand"
	"time"
)

type Stage string
type Action func(app App, out Out)

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

const (
	FetchRepo     Stage = "FetchRepository"
	Validate      Stage = "Validate"
	InstallJS     Stage = "InstallJSDependencies"
	BuildFrontend Stage = "BuildFrontend"
	InstallPy     Stage = "InstallPythonDependencies"
	Complete      Stage = "Complete"
)

func RandSleep(max float64) {
	duration := time.Duration(max*rand.Float64()) * time.Millisecond
	time.Sleep(duration)
}
