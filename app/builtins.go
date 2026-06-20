package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type PwdCommand struct {
	flag *pflag.FlagSet
}

func (c *PwdCommand) Execute(args []string, s *Shell) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(s.Out, cwd)
	return nil
}

type EchoCommand struct {
	flag *pflag.FlagSet
}

func (c *EchoCommand) Execute(args []string, s *Shell) error {
	fmt.Fprintln(s.Out, strings.Join(args, " "))
	return nil
}

type ExitCommand struct {
	flag *pflag.FlagSet
}

func (c *ExitCommand) Execute(args []string, s *Shell) error {
	defer os.Exit(0)
	return nil
}

type TypeCommand struct {
	flag *pflag.FlagSet
}

func (c *TypeCommand) Execute(args []string, s *Shell) error {
	if len(args) > 1 {
		fmt.Fprintln(s.Err, args, "not found")
	}
	if _, ok := s.Builtins[args[0]]; ok {
		fmt.Fprintln(s.Out, args[0], "is a shell builtin")
	} else {
		handleTypeExecutable(s.EnvPath, args[0])
	}

	return nil
}

type CdCommand struct {
	flag *pflag.FlagSet
}

func (c *CdCommand) Execute(args []string, s *Shell) error {
	path := args[0]

	if args[0] == "~" {
		path = s.Home
	}
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("cd: %s: No such file or directory", args[0])
	}
	return nil
}

type CompleteCommand struct {
	flag   *pflag.FlagSet
	values map[string]string
}

func (c *CompleteCommand) Execute(args []string, s *Shell) error {
	c.flag = pflag.NewFlagSet("complete", pflag.ContinueOnError)
	print := c.flag.BoolP("print", "p", true, "prints the completion specification registered for the command")
	register := c.flag.BoolP("register", "C", true, "registers a completer script for the command")

	err := c.flag.Parse(args)
	if err != nil {
		return err
	}

	// c.values = make(map[string]string)

	remaining := c.flag.Args()
	switch {
	case *register:
		c.values = map[string]string{
			remaining[1]: remaining[0],
		}
	case *print:
		if *register {
			fmt.Fprintf(s.Out, "complete -C '%s' %s", remaining[0], c.values[remaining[0]])
		} else {
			fmt.Fprintf(s.Out, "complete: %s: no completion specification\n", remaining[0])
		}
	}

	// if *print {
	// 	fmt.Fprintf(s.Out, "complete: %s: no completion specification\n", remaining[0])
	// }
	return nil
}
