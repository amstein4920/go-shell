package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	EXIT = "EXIT"
	ECHO = "ECHO"
	TYPE = "TYPE"
	PWD  = "PWD"
	CD   = "CD"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, "$ ")
	for scanner.Scan() {
		input := scanner.Text()

		if input == "" {
			fmt.Fprint(os.Stdout, "$ ")
			continue
		}

		stripped := strings.Trim(input, "\n")
		parsed := parse(stripped)
		command := parsed[0]
		arguments := parsed[1:]

		switch strings.ToUpper(command) {
		case EXIT:
			os.Exit(0)
		case ECHO:
			fmt.Println(strings.Join(arguments, " "))
		case TYPE:
			fmt.Println(typeBuiltinFunction(parsed[1]))
		case PWD:
			wd, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "reading working directory", err)
			}
			fmt.Println(wd)
		case CD:
			err := cdBuiltinFunction(arguments[0])
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Print(commandExecutionFunction(command, arguments))
		}
		fmt.Fprint(os.Stdout, "$ ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input", err)
		os.Exit(1)
	}
}

func parse(input string) []string {
	var result []string
	var temp strings.Builder
	singleQuoted := false
	doubleQuoted := false
	escaped := false

	for _, char := range input {
		// if escaped {
		// 	if doubleQuoted || singleQuoted {
		// 		temp.WriteRune('\\')
		// 		temp.WriteRune(char)
		// 	} else {
		// 		temp.WriteRune(char)
		// 	}
		// 	escaped = false
		// 	continue
		// }
		switch char {
		case '\\':
			if escaped || singleQuoted {
				temp.WriteRune(char)
				escaped = false
			} else {
				escaped = true
			}
		case ' ':
			if escaped && doubleQuoted {
				temp.WriteRune('\\')
			}
			if singleQuoted || doubleQuoted || escaped {
				temp.WriteRune(char)
			} else if temp.Len() > 0 {
				result = append(result, temp.String())
				temp.Reset()
			}
			escaped = false
		case '\'':
			if escaped && doubleQuoted {
				temp.WriteRune('\\')
			}
			if doubleQuoted || escaped {
				temp.WriteRune(char)
			} else {
				singleQuoted = !singleQuoted
			}
			escaped = false
		case '"':
			if singleQuoted || escaped {
				temp.WriteRune(char)
			} else {
				doubleQuoted = !doubleQuoted
			}
			escaped = false
		default:
			if doubleQuoted && escaped {
				temp.WriteRune('\\')
			}
			temp.WriteRune(char)
			escaped = false
		}
	}

	if temp.Len() > 0 {
		result = append(result, temp.String())
	}
	return result
}
