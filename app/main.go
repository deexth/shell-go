package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Print("$ ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fmt.Fprintf(os.Stderr, "%s: command not found\n", scanner.Text())
}
