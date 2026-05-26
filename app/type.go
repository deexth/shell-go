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

func handleExecutable(cmds []string) error {
	cmd := exec.Command(cmds[0], cmds[1:]...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Program was passed %d args (including program name).\n", len(cmds))
	fmt.Fprintf(os.Stdout, "Arg #0 (program name): %s\n", cmds[0])
	for i, c := range cmds[1:] {
		fmt.Fprintf(os.Stdout, "Arg #%d: %s\n", i+1, c)
	}
	fmt.Fprintln(os.Stdout, "Program signature: 1")
	return nil
}
