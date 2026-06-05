package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
)

func main() {
	sh := NewShell()
	sh.Register("pwd", &PwdCommand{})
	sh.Register("echo", &EchoCommand{})
	sh.Register("exit", &ExitCommand{})
	sh.Register("type", &TypeCommand{})
	sh.Register("cd", &CdCommand{})

	completer := sh.BuildCompleter()

	cfg := &readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
	}

	l, err := readline.NewEx(cfg)
	if err != nil {
		fmt.Fprintln(sh.Err, err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	// reader := bufio.NewReader(sh.In)

	for {
		// fmt.Fprint(sh.Out, "$ ")
		// input, err := reader.ReadString('\n')
		input, err := l.Readline()
		if err != nil {
			break
		}

		cmdName, args := parseInput(input)

		args, filename, isStderr, isAppendErr, isAppend := parseRedirection(args)

		origOut := sh.Out
		origErr := sh.Err

		var file *os.File
		if filename != "" {
			flags := os.O_CREATE | os.O_WRONLY
			if isAppend || isAppendErr {
				flags |= os.O_APPEND
			} else {
				flags |= os.O_TRUNC
			}

			file, err = os.OpenFile(filename, flags, 0644)
			if err != nil {
				fmt.Fprintln(sh.Err, err)
				continue
			}

			if isStderr || isAppendErr {
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
			handleExternal(cmdName, args, sh)
		}

		if file != nil {
			file.Close()
		}

		sh.Out = origOut
		sh.Err = origErr
	}

}
