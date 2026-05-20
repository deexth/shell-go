package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input: ", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "%s: command not found\n", command[:len(command)-1])
	}
}
