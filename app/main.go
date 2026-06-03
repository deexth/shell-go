package main

import (
	"bufio"
	"fmt"
	"os"
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

		args, filename, isStderr, isAppend := parseRedirection(args)

		origOut := sh.Out
		origErr := sh.Err

		var file *os.File
		if filename != "" {
			flags := os.O_CREATE | os.O_WRONLY
			if isAppend {
				flags |= os.O_APPEND
			} else {
				flags |= os.O_TRUNC
			}

			file, err = os.OpenFile(filename, flags, 0644)
			if err != nil {
				fmt.Fprintln(sh.Err, err)
				continue
			}

			if isStderr {
				sh.Err = file
			} else {
				sh.Out = file
			}
		}

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

		if file != nil {
			file.Close()
		}

		sh.Out = origOut
		sh.Err = origErr
	}

}
