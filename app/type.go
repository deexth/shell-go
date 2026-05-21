package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func handleType(msg, path string) {
	if _, ok := commands[msg]; ok {
		fmt.Fprintln(os.Stdout, msg, "is a shell builtin")
	} else {
		handlePath(path, msg)
	}
}

func handlePath(path, msg string) {
	for dir := range strings.SplitSeq(path, ":") {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && d.Type().IsRegular() && d.Name() == msg {
				fmt.Fprintln(os.Stdin, msg, " is ", path)
			}

			fmt.Fprintf(os.Stderr, "%s: not found\n", msg)
			return nil

		})
		if err != nil {
			continue
		}
	}
}
