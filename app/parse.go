package main

import (
	"strings"
)

func parseInput(input string) (string, []string) {
	input = strings.TrimSpace(input)
	var args []string
	var currentArg strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}

		if char == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if char == ' ' && (!inSingleQuote && !inDoubleQuote) {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
			continue
		}

		if char == '\\' && (!inSingleQuote && !inDoubleQuote) {
			if i+2 <= len(input) {
				currentArg.WriteByte(input[i+1])
				i = i + 1
			}
			continue
		}

		if char == '\\' && inDoubleQuote {
			if (i+2 < len(input)) && (input[i+1] == '\\' || input[i+1] == '"') {
				currentArg.WriteByte(input[i+1])
				i = i + 1
			} else {
				currentArg.WriteByte(char)
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

func parseRedirection(args []string) (cleanedArgs []string, filename string, isStderr, isAppend bool) {
	for i := range args {
		sep := args[i]
		if i+1 < len(args) {
			if sep == "1>" || sep == ">" || sep == "2>" || sep == "2>>" || sep == ">>" || sep == "1>>" {
				filename = args[i+1]
				cleanedArgs = append(cleanedArgs, args[:i]...)
				isAppend = (sep == "2>>" || sep == "1>>" || sep == ">>")
				isStderr = sep == "2>"

				return cleanedArgs, filename, isStderr, isAppend
			}
		}
	}
	return args, "", false, false
}
