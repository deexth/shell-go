package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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

		if msg, ok := strings.CutPrefix(cmd, "echo "); ok {
			fmt.Fprintln(os.Stdout, msg)
			continue
		}

		fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd)
	}
}
