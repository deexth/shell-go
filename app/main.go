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

		cmds := strings.Split(cmd, " ")

		if err = handleExecutable(cmds[0], cmds[1:]...); err != nil {

			fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd)
		}
	}
}

func checkPrefix(input, cmd string) (string, bool) {
	return strings.CutPrefix(input, cmd)
}
