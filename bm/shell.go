package bm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TODO:
// Not all output is written. This is cause programs can tell if
// output stream is directed to a tty and change their output depending
// on that. Git can be forced to show the output by using --progress
// Not sure if this is important, will reconsider later.

type Shell struct {
	Output chan Output
	Stage  Stage
	Env    []string
}

func NewShell(output chan Output, stage Stage) *Shell {
	return &Shell{Output: output, Stage: stage, Env: []string{}}
}

func (sh *Shell) AppendEnv(env string) {
	sh.Env = append(sh.Env, env)
}

func (sh *Shell) Run(cmd string) error {
	sh.Output <- Output{
		fmt.Sprintf("$ %s", cmd),
		sh.Stage,
	}

	splits := strings.Split(cmd, " ")
	command := exec.Command(splits[0], splits[1:]...)
	command.Env = append(os.Environ(), sh.Env...)

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
