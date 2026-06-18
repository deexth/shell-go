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
	In       io.Reader
	Out      io.Writer
	Err      io.Writer
	EnvPath  string
	Home     string
	Builtins map[string]Command
	Files    struct {
		CustomExecutables []string
		OtherFiles        [][]rune
	}
}

type ShellCompleter struct {
	*Shell
	RlInstance *readline.Instance
	LastLine   string
	TabCount   int
}

func NewShell() *Shell {
	path, _ := os.LookupEnv("PATH")
	home, _ := os.LookupEnv("HOME")
	cExec, Ofiles := getPossileExecutables(path)
	return &Shell{
		In:       os.Stdin,
		Out:      os.Stdout,
		Err:      os.Stderr,
		EnvPath:  path,
		Home:     home,
		Builtins: make(map[string]Command),
		Files: struct {
			CustomExecutables []string
			OtherFiles        [][]rune
		}{
			CustomExecutables: cExec,
			OtherFiles:        Ofiles,
		},
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

	if strings.TrimSpace(prefix) == "" {
		return nil, 0
	}

	current := fields[0]
	seen := make(map[string]bool, 0)
	var matches []string

	if len(fields) > 1 {
		file := fields[len(fields)-1]
		for _, f := range s.Files.OtherFiles {
			if strings.HasPrefix(string(f), file) {
				if !seen[string(f)] {
					seen[string(f)] = true
					matches = append(matches, string(f))
				}
			}
		}
	}

	for cmd := range s.Builtins {
		if strings.HasPrefix(cmd, current) {
			if !seen[cmd] {
				seen[cmd] = true
				matches = append(matches, cmd)
			}
		}

	}

	for _, cmd := range s.Files.CustomExecutables {
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
		suffix := strings.TrimPrefix(matches[0], current)
		choice := suffix + " "
		return [][]rune{[]rune(choice)}, len(current)
	}

	lcp := longestCommonPrefix(matches)
	if len(lcp) > len(current) {
		s.LastLine = ""
		s.TabCount = 0
		suffix := lcp[len(current):]
		return [][]rune{[]rune(suffix)}, len(current)
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

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		for !strings.HasPrefix(strs[i], prefix) {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

func getPossileExecutables(path string) ([]string, [][]rune) {
	possibleExecutable := make([]string, 0)
	allFiles := make([][]rune, 0)

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
			allFiles = append(allFiles, []rune(name))

		}

	}

	return possibleExecutable, allFiles
}
