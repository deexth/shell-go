package main

import (
	"io"
	"os"
	"strings"
)

type Shell struct {
	In                io.Reader
	Out               io.Writer
	Err               io.Writer
	EnvPath           string
	Home              string
	Builtins          map[string]Command
	CustomExecutables []string
}

type ShellCompleter struct {
	*Shell
}

func NewShell() *Shell {
	path, _ := os.LookupEnv("PATH")
	home, _ := os.LookupEnv("HOME")
	return &Shell{
		In:                os.Stdin,
		Out:               os.Stdout,
		Err:               os.Stderr,
		EnvPath:           path,
		Home:              home,
		Builtins:          make(map[string]Command),
		CustomExecutables: getPossileExecutables(path),
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
		if strings.HasPrefix(cmd, prefix) {
			suffix := cmd[len(prefix):]
			matches = append(matches, []rune(suffix+" "))
		}
	}

	if len(matches) == 0 {
		s.Out.Write([]byte{'\x07'})
	}

	for _, cmd := range s.CustomExecutables {
		if strings.HasPrefix(cmd, prefix) {
			suffix := cmd[len(prefix):]
			matches = append(matches, []rune(suffix+" "))
		}
	}

	if len(matches) == 0 {
		s.Out.Write([]byte{'\x07'})
	}

	matches = sortMatches(matches)

	return matches, len(prefix)
}

func getPossileExecutables(path string) []string {
	possibleExecutable := make([]string, 0)

	for p := range strings.SplitSeq(path, ":") {
		files, err := os.ReadDir(p)
		if err != nil {
			continue
		}

		seen := make(map[string]bool, 0)
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			name := file.Name()
			if seen[name] {
				continue
			}

			f, err := file.Info()
			if err != nil {
				continue
			}

			if f.Mode().Perm()&0111 != 0 {
				possibleExecutable = append(possibleExecutable, name)
				seen[name] = true
			}

		}

	}

	return possibleExecutable
}

func sortMatches(matches [][]rune) [][]rune {
	for i := 1; i < len(matches); i++ {
		for j := i; j > 0; j-- {
			if string(matches[j-1]) < string(matches[j]) {
				continue
			}

			matches[j-1], matches[j] = matches[j], matches[j-1]
		}
	}
	return matches
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
