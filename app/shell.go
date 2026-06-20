package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chzyer/readline"
)

type Shell struct {
	In          io.Reader
	Out         io.Writer
	Err         io.Writer
	EnvPath     string
	Home        string
	Builtins    map[string]Command
	Executables []string
}

type ShellCompleter struct {
	*Shell
	RlInstance *readline.Instance
	lastLine   string
	tabCount   int
}

func NewShell() *Shell {
	path, _ := os.LookupEnv("PATH")
	home, _ := os.LookupEnv("HOME")
	return &Shell{
		In:          os.Stdin,
		Out:         os.Stdout,
		Err:         os.Stderr,
		EnvPath:     path,
		Home:        home,
		Builtins:    make(map[string]Command),
		Executables: findExecutables(path),
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

	if strings.TrimSpace(prefix) == "" {
		return nil, 0
	}

	fields := strings.Fields(prefix)

	completeArg := len(fields) > 1 || strings.HasSuffix(prefix, " ")

	if completeArg {
		partial := ""
		if !strings.HasSuffix(prefix, " ") {
			partial = fields[len(fields)-1]
		}
		return s.completeFilename(line, pos, partial)
	}

	return s.completeCommand(line, pos, fields[0])
}

func (s *ShellCompleter) completeFilename(line []rune, pos int, partial string) ([][]rune, int) {
	dir, filePrefix := filepath.Split(partial)
	readDir := dir
	if readDir == "" {
		readDir = "."
	}

	entries, err := os.ReadDir(readDir)
	if err != nil {
		s.ringBell()
		return nil, 0
	}
	var matches []string

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, filePrefix) {
			continue
		}

		candidate := dir + name
		if entry.IsDir() {
			candidate += "/"
		}

		matches = append(matches, candidate)
	}

	return s.applyMatches(line, pos, partial, matches)
}

func (s *ShellCompleter) completeCommand(line []rune, pos int, partial string) ([][]rune, int) {
	seen := make(map[string]bool)
	var matches []string

	for cmd := range s.Builtins {
		if strings.HasPrefix(cmd, partial) && !seen[cmd] {
			seen[cmd] = true
			matches = append(matches, cmd)
		}
	}

	for _, cmd := range s.Executables {
		if strings.HasPrefix(cmd, partial) && !seen[cmd] {
			seen[cmd] = true
			matches = append(matches, cmd)
		}
	}

	return s.applyMatches(line, pos, partial, matches)
}

func (s *ShellCompleter) applyMatches(line []rune, pos int, partial string, matches []string) ([][]rune, int) {
	if len(matches) == 0 {
		s.ringBell()
		return nil, 0
	}

	sort.Strings(matches)

	if len(matches) == 1 {
		s.resetTabState()
		suffix := strings.TrimPrefix(matches[0], partial)
		if !strings.HasSuffix(matches[0], "/") {
			suffix += " "
		}
		return [][]rune{[]rune(suffix)}, len(partial)
	}

	lcp := longestCommonPrefix(matches)
	if len(lcp) > len(partial) {
		s.resetTabState()
		return [][]rune{[]rune(lcp[len(partial):])}, len(partial)
	}

	state := string(line[:pos])
	if state == s.lastLine {
		s.tabCount += 1
	} else {
		s.lastLine = state
		s.tabCount = 1
	}

	if s.tabCount < 2 {
		s.ringBell()
		return nil, 0
	}

	w := s.Out
	if s.RlInstance != nil {
		w = s.RlInstance.Stdout()
	}
	fmt.Fprintf(w, "\r\n%s\r\n", strings.Join(matches, "  "))

	return nil, 0
}

func (s *ShellCompleter) ringBell() {
	s.Out.Write([]byte{'\x07'})
}

func (s *ShellCompleter) resetTabState() {
	s.lastLine = ""
	s.tabCount = 0
}

func longestCommonPrefix(strs []string) string {
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}

func findExecutables(path string) []string {
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
