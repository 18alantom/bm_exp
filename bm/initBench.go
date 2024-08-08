package bm

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path"
)

type CommonSiteConfig struct {
	SocketIOPort  int `json:"socketio_port"`
	WebServerPort int `json:"webserver_port"`
}

func (exec *Exec) initBench() error {
	// This output should go into the common bench output
	if err := os.RemoveAll(exec.Ctx.Target); err != nil {
		return err
	}

	return writeCommonSiteConfig(exec)
}

func writeCommonSiteConfig(exec *Exec) error {
	sites := path.Join(exec.Ctx.Target, "sites")
	if err := os.MkdirAll(sites, 0o755); err != nil {
		return err
	}

	csc := CommonSiteConfig{9000, 8000}
	cscBytes, err := json.Marshal(csc)
	if err != nil {
		return err
	}

	cscPath := path.Join(sites, "common_site_config.json")
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
