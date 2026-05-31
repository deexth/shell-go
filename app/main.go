package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	sh := NewShell()
	sh.Register("pwd", &PwdCommand{})
	sh.Register("echo", &EchoCommand{})
	sh.Register("exit", &ExitCommand{})
	sh.Register("type", &TypeCommand{})
	sh.Register("cd", &CdCommand{})

	reader := bufio.NewReader(sh.In)

	for {
		fmt.Fprint(sh.Out, "$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		cmdName, args := parseInput(input)

		if cmd, exists := sh.Builtins[cmdName]; exists {
			err = cmd.Execute(args, sh)
			if err != nil {
				fmt.Fprintln(sh.Err, err)
			}

		} else {
			if err = handleExternal(cmdName, args, sh); err != nil {
				fmt.Fprintf(sh.Err, "%s: command not found\n", cmdName)
			}
		}
	}

}

func parseInput(input string) (string, []string) {
	cmdArgs := strings.Fields(input)
	return cmdArgs[0], cmdArgs[1:]
}
