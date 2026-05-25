package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commands = map[string]struct{}{
	"echo": {},
	"type": {},
	"exit": {},
}

func main() {
	path, ok := os.LookupEnv("PATH")
	if !ok {
		fmt.Fprintln(os.Stderr, "PATH not provided")
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input: ", err)
			os.Exit(1)
		}

		cmd := strings.TrimSpace(command)

		if strings.ToLower(cmd) == "exit" {
			break
		}

		if msg, ok := checkPrefix(cmd, "echo "); ok {
			fmt.Fprintln(os.Stdout, msg)
			continue
		}

		if msg, ok := checkPrefix(cmd, "type "); ok {
			handleType(msg, path)
			continue
		}

		handleExecutable(path, strings.ToLower(cmd))

		fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd)
	}
}

func checkPrefix(input, cmd string) (string, bool) {
	return strings.CutPrefix(input, cmd)
}

// func handleType(msg, path string) {
// 	if _, ok := commands[msg]; ok {
// 		fmt.Fprintln(os.Stdout, msg, "is a shell builtin")
// 	} else {
// 		handleExecutable(path, msg)
// 	}
// }
//
// func handleExecutable(path, msg string) {
// 	for dir := range strings.SplitSeq(path, ":") {
// 		fullPath := filepath.Join(dir, msg)
// 		f, err := os.Stat(fullPath)
// 		if err != nil {
// 			continue
// 		}
//
// 		if f.Mode().Perm()&0111 != 0 {
// 			fmt.Fprintf(os.Stdout, "%s is %s\n", msg, fullPath)
// 			return
// 		} else {
// 			continue
// 		}
// 	}
//
// 	fmt.Fprintf(os.Stderr, "%s: not found\n", msg)
//
// }
