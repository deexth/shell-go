package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func handleTypeExecutable(path, msg string) {
	for dir := range strings.SplitSeq(path, ":") {
		fullPath := filepath.Join(dir, msg)
		f, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		if f.Mode().Perm()&0111 != 0 {
			fmt.Fprintf(os.Stdout, "%s is %s\n", msg, fullPath)
			return
		} else {
			continue
		}
	}

	fmt.Fprintf(os.Stderr, "%s: not found\n", msg)

}

func handleExternal(cmdName string, args []string, s *Shell) error {
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = s.Out
	cmd.Stderr = s.Err
	cmd.Stdin = s.In
	return cmd.Run()

}
