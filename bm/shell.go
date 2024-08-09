package bm

import (
	"fmt"
	"os/exec"
	"strings"
)

type Shell struct {
	Output chan Output
	Stage  Stage
}

func (sh Shell) Run(cmd string) error {
	sh.Output <- Output{
		fmt.Sprintf("$ %s", cmd),
		sh.Stage,
	}

	splits := strings.Split(cmd, " ")
	command := exec.Command(splits[0], splits[1:]...)

	command.Stdout = ChanWriter{sh.Output, sh.Stage}
	command.Stderr = ChanWriter{sh.Output, sh.Stage}

	return command.Run()
}

type ChanWriter struct {
	output chan Output
	stage  Stage
}

func (cw ChanWriter) Write(p []byte) (n int, err error) {
	cw.output <- Output{string(p), cw.stage}
	return len(p), nil
}
