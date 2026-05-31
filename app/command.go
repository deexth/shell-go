package main

type Command interface {
	Execute(args []string, s *Shell) error
}
