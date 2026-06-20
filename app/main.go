package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/spf13/pflag"
)

func main() {
	sh := NewShell()
	sh.Register("pwd", &PwdCommand{
		flag: pflag.NewFlagSet("pwd", pflag.ContinueOnError),
	})
	sh.Register("echo", &EchoCommand{
		flag: pflag.NewFlagSet("pwd", pflag.ContinueOnError),
	})
	sh.Register("exit", &ExitCommand{
		flag: pflag.NewFlagSet("pwd", pflag.ContinueOnError),
	})
	sh.Register("type", &TypeCommand{
		flag: pflag.NewFlagSet("pwd", pflag.ContinueOnError),
	})
	sh.Register("cd", &CdCommand{
		flag: pflag.NewFlagSet("pwd", pflag.ContinueOnError),
	})
	sh.Register("complete", &CompleteCommand{
		flag: pflag.NewFlagSet("pwd", pflag.ContinueOnError),
	})

	completer := NewShellCompleter(sh)

	cfg := &readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
	}

	reader, err := readline.NewEx(cfg)
	if err != nil {
		fmt.Fprintln(sh.Err, err)
		os.Exit(1)
	}
	defer reader.Close()
	reader.CaptureExitSignal()

	completer.RlInstance = reader

	for {
		input, err := reader.Readline()
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
