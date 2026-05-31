package main

import (
	"io"
	"os"
)

type Shell struct {
	In       io.Reader
	Out      io.Writer
	Err      io.Writer
	EnvPath  string
	Builtins map[string]Command
}

func NewShell() *Shell {
	path, _ := os.LookupEnv("PATH")
	return &Shell{
		In:       os.Stdin,
		Out:      os.Stdout,
		Err:      os.Stderr,
		EnvPath:  path,
		Builtins: make(map[string]Command),
	}
}

func (s *Shell) Register(name string, cmd Command) {
	s.Builtins[name] = cmd
}
