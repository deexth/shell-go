package main

import (
	"io"
	"os"
	"strings"
)

type Shell struct {
	In       io.Reader
	Out      io.Writer
	Err      io.Writer
	EnvPath  string
	Home     string
	Builtins map[string]Command
}

type ShellCompleter struct {
	*Shell
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

func NewShellCompleter(s *Shell) *ShellCompleter {
	return &ShellCompleter{
		Shell: s,
	}
}

func (s *ShellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	prefix := string(line[:pos])

	var matches [][]rune

	for cmd := range s.Builtins {
		suffix, ok := strings.CutPrefix(cmd, prefix)
		if ok {
			matches = append(matches, []rune(suffix))
		}
	}

	if len(matches) == 0 {
		s.Out.Write([]byte{'\x07'})
	}

	return matches, 0
}

// func (s *Shell) BuildCompleter() *readline.PrefixCompleter {
// 	cmds := make([]readline.PrefixCompleterInterface, 0)
//
// 	for builtin := range s.Builtins {
// 		cmds = append(cmds, readline.PcItem(builtin))
// 	}
//
// 	return readline.NewPrefixCompleter(cmds...)
// }
