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

	// Guard: Only autocomplete if the user is typing the command name
	fields := strings.Fields(prefix)
	if len(fields) > 1 || strings.HasSuffix(prefix, " ") || strings.TrimSpace(prefix) == "" {
		return nil, 0
	}

	current := fields[0]
	var matchNames []string
	seen := make(map[string]bool)

	// Gather matching builtins
	for cmd := range s.Builtins {
		if strings.HasPrefix(cmd, current) {
			if !seen[cmd] {
				suffix := cmd[len(current):]
				seen[cmd] = true
				matchNames = append(matchNames, suffix)
			}
		}
	}

	// Gather matching custom executables
	for _, cmd := range s.CustomExecutables {
		if strings.HasPrefix(cmd, current) {
			if !seen[cmd] {
				suffix := cmd[len(current):]
				seen[cmd] = true
				matchNames = append(matchNames, suffix)
			}
		}
	}

	// Case 0: No matches found
	if len(matchNames) == 0 {
		s.Out.Write([]byte{'\x07'}) // Ring bell
		return nil, 0
	}

	// Sort matching items alphabetically
	sort.Strings(matchNames)

	// Case 1: Single match -> Autocomplete immediately
	if len(matchNames) == 1 {
		s.LastLine = ""
		s.TabCount = 0
		choice := matchNames[0] + " "
		return [][]rune{[]rune(choice)}, len(current)
	}

	// Case 2: Multiple matches -> Manage 1st vs 2nd TAB press
	state := string(line[:pos])
	if state == s.LastLine {
		s.TabCount++
	} else {
		s.LastLine = state
		s.TabCount = 1
	}

	if s.TabCount == 1 {
		// First TAB press: Ring the bell and suppress default menu behavior
		s.Out.Write([]byte{'\x07'})
		return nil, 0
	}

	// Second TAB press: Print options on a new line
	fmt.Fprint(s.Out, "\n")
	fmt.Fprintln(s.Out, strings.Join(matchNames, "  ")) // Separated by two spaces

	// Redraw the prompt and re-populate the original prefix on the new line
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
