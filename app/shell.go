package main

import (
	"io"
	"os"

	"github.com/chzyer/readline"
)

type Shell struct {
	In       io.Reader
	Out      io.Writer
	Err      io.Writer
	EnvPath  string
	Home     string
	Builtins map[string]Command
}

func NewShell() *Shell {
	path, _ := os.LookupEnv("PATH")
	home, _ := os.LookupEnv("HOME")
	return &Shell{
		In:       os.Stdin,
		Out:      os.Stdout,
		Err:      os.Stderr,
		EnvPath:  path,
		Home:     home,
		Builtins: make(map[string]Command),
	}
}

func (s *Shell) Register(name string, cmd Command) {
	s.Builtins[name] = cmd
}

func (s *Shell) BuildCompleter() *readline.PrefixCompleter {
	cmds := make([]readline.PrefixCompleterInterface, 0)

	for builtin := range s.Builtins {
		cmds = append(cmds, readline.PcItem(builtin))
	}

	return readline.NewPrefixCompleter(cmds...)
}
