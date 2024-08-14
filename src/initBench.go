package bm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

type CommonSiteConfig struct {
	SocketIOPort  int `json:"socketio_port"`
	WebServerPort int `json:"webserver_port"`
}

var pathsInBench = []string{
	"apps", "sites", "sites/assets", "config", "logs", "config/pids",
}

func (exec *Exec) initBench(output chan Output) error {
	output <- Output{"Initializing directories", Bench, time.Now()}
	// This output should go into the common bench output
	if err := exec.initDirs(); err != nil {
		return err
	}

	output <- Output{"Initializing python env", Bench, time.Now()}
	exec.initPythonEnv(output)

	output <- Output{"Initializing config", Bench, time.Now()}
	return exec.initFiles()
}

func (exec *Exec) initDirs() error {
	// Remove prior bench if it exists
	// TODO: Probably should fail? match frappe/bench behavior
	if err := os.RemoveAll(exec.Ctx.Target); err != nil {
		return err
	}

	// Create bench folder
	if err := os.MkdirAll(exec.Ctx.Target, 0o755); err != nil {
		return err
	}

	for _, sub := range pathsInBench {
		subDir := path.Join(exec.Ctx.Target, sub)
		os.MkdirAll(subDir, 0o755)
	}

	return nil
}

func (exec *Exec) initPythonEnv(output chan Output) error {
	envPath := path.Join(exec.Ctx.Target, "env")
	command := fmt.Sprintf("python -m venv %s", envPath)
	return NewShell(output, Bench).Run(command)
}

func (exec *Exec) initFiles() error {
	if err := exec.initConfigJson(); err != nil {
		return err
	}

	if err := exec.initAppsTxt(); err != nil {
		return err
	}

	return exec.initAssetsJson()
}

func (exec *Exec) initAssetsJson() error {
	json, err := json.Marshal(struct{}{})
	if err != nil {
		return err
	}

	assets := path.Join(exec.Ctx.Target, "sites", "assets", "assets.json")

	f, err := os.Create(assets)
	if err != nil {
		return err
	}

	l, err := f.Write(json)
	if err != nil {
		return err
	}

	if l != len(json) {
		return errors.New("could not write assets.json")
	}

	return nil
}

func (exec *Exec) initAppsTxt() error {
	apps := path.Join(exec.Ctx.Target, "sites", "apps.txt")

	f, err := os.Create(apps)
	if err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return nil
}

// Dummy functions writes common_site_config fields only to the extent that it's
// required by app installs
func (exec *Exec) initConfigJson() error {
	csc := CommonSiteConfig{9000, 8000}
	cscBytes, err := json.Marshal(csc)
	if err != nil {
		return err
	}

	cscPath := path.Join(exec.Ctx.Target, "sites", "common_site_config.json")
	f, err := os.Create(cscPath)
	if err != nil {
		return err
	}

	cscBB := bytes.Buffer{}
	err = json.Indent(&cscBB, cscBytes, "", " ")
	if err == nil {
		cscBytes = cscBB.Bytes()
	}

	l, err := f.Write(cscBytes)
	if err != nil {
		return err
	}

	if l != len(cscBytes) {
		return errors.New("could not write common_site_config.json")
	}

	return nil
}
