package main

import (
	"fmt"
	"os"
)

type PwdCommand struct{}

func (c *PwdCommand) Execute(args []string, s *Shell) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(s.Out, cwd)
	return nil
}

type EchoCommand struct{}

func (c *EchoCommand) Execute(args []string, s *Shell) error {
	fmt.Fprintln(s.Out, args)
	return nil
}

type ExitCommand struct{}

func (c *ExitCommand) Execute(args []string, s *Shell) error {
	defer os.Exit(0)
	return nil
}

type TypeCommand struct{}

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

type CdCommand struct{}

func (c *CdCommand) Execute(args []string, s *Shell) error {
	return nil
}
