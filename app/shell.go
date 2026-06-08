package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/chzyer/readline"
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
	RlInstance *readline.Instance // Reference to refresh the prompt line
	LastLine   string             // Tracks the current input text state
	TabCount   int
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

	fields := strings.Fields(prefix)

	if len(prefix) > 1 || strings.HasPrefix(prefix, " ") || strings.TrimSpace(prefix) == "" {
		return nil, 0
	}

	current := fields[0]
	seen := make(map[string]bool, 0)
	var matches []string

	for cmd := range s.Builtins {
		if strings.HasPrefix(cmd, current) {
			if !seen[cmd] {
				seen[cmd] = true
				matches = append(matches, cmd)
			}
		}

	}

	for _, cmd := range s.CustomExecutables {
		if strings.HasPrefix(cmd, current) {
			if !seen[cmd] {
				seen[cmd] = true
				matches = append(matches, cmd)
			}
		}
	}

	if len(matches) == 0 {
		s.Out.Write([]byte{'\x07'})
		return nil, 0
	}

	sort.Strings(matches)

	if len(matches) == 1 {
		s.LastLine = ""
		s.TabCount = 0
		choice := matches[0] + " "
		return [][]rune{[]rune(choice)}, len(current)
	}

	state := string(line[:pos])
	if state == s.LastLine {
		s.TabCount++
	} else {
		s.LastLine = state
		s.TabCount = 1
	}

	if s.TabCount == 1 {
		s.Out.Write([]byte{'\x07'})
		return nil, 0
	}

	fmt.Fprint(s.Out, "\n")
	fmt.Fprintln(s.Out, strings.Join(matches, " "))

	if s.RlInstance != nil {
		s.RlInstance.Refresh()
	}

	return nil, 0
}

func getPossileExecutables(path string) []string {
	possibleExecutable := make([]string, 0)

	seen := make(map[string]bool, 0)

	for p := range strings.SplitSeq(path, ":") {
		files, err := os.ReadDir(p)
		if err != nil {
			continue
		}

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

// func sortMatches(matches [][]rune) [][]rune {
// 	for i := 1; i < len(matches); i++ {
// 		for j := i; j > 0; j-- {
// 			if string(matches[j-1]) <= string(matches[j]) {
// 				break
// 			}
//
// 			matches[j-1], matches[j] = matches[j], matches[j-1]
// 		}
// 	}
// 	return matches
// }

// func (s *Shell) BuildCompleter() *readline.PrefixCompleter {
// 	cmds := make([]readline.PrefixCompleterInterface, 0)
//
// 	for builtin := range s.Builtins {
// 		cmds = append(cmds, readline.PcItem(builtin))
// 	}
//
// 	return readline.NewPrefixCompleter(cmds...)
// }
