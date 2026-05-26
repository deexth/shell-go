package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func handleType(msg, path string) {
	if _, ok := commands[msg]; ok {
		fmt.Fprintln(os.Stdout, msg, "is a shell builtin")
	} else {
		handleTypeExecutable(path, msg)
	}
}

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

func handleExecutable(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}
