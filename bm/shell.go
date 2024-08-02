package bm

import (
	"os/exec"
	"strings"
)

type Shell struct {
	Out   chan Output
	Stage Stage
}

func (sh Shell) Run(cmd string) error {
	splits := strings.Split(cmd, " ")
	command := exec.Command(splits[0], splits[1:]...)

	command.Stdout = ChanWriter{sh.Out, sh.Stage}
	command.Stderr = ChanWriter{sh.Out, sh.Stage}

	return command.Run()
}

type ChanWriter struct {
	out   chan Output
	stage Stage
}

func (cw ChanWriter) Write(p []byte) (n int, err error) {
	cw.out <- Output{string(p), cw.stage}
	return len(p), nil
}
