package bm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Shell struct {
	Output chan Output
	Stage  Stage
	Env    []string
}

func NewShell(output chan Output, stage Stage) Shell {
	return Shell{Output: output, Stage: stage, Env: []string{}}
}

func (sh *Shell) AppendEnv(env string) {
	sh.Env = append(sh.Env, env)
}

func (sh Shell) Run(cmd string) error {
	sh.Output <- Output{
		fmt.Sprintf("$ %s", cmd),
		sh.Stage,
	}

	splits := strings.Split(cmd, " ")
	command := exec.Command(splits[0], splits[1:]...)
	command.Env = append(os.Environ(), command.Env...)

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
