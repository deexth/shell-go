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
	input = strings.TrimSpace(input)
	var args []string
	var currentArg strings.Builder
	inSingleQuote := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '\'' {
			inSingleQuote = !inSingleQuote
			continue
		}

		if char == ' ' && !inSingleQuote {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
			continue
		}

		currentArg.WriteByte(char)
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	if len(args) == 0 {
		return "", nil
	}

	return args[0], args[1:]
}
