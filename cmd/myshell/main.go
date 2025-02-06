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
			fmt.Print(commandExecutionFunction(parsed))
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

	for _, char := range input {
		switch char {
		case ' ':
			if singleQuoted || doubleQuoted {
				temp.WriteRune(char)
			} else if temp.Len() > 0 {
				result = append(result, temp.String())
				temp.Reset()
			}
		case '\'':
			if doubleQuoted {
				temp.WriteRune(char)
				break
			}
			if singleQuoted {
				singleQuoted = false
			} else {
				singleQuoted = true
			}
		case '"':
			if doubleQuoted {
				doubleQuoted = false
			} else {
				doubleQuoted = true
			}
		default:
			temp.WriteRune(char)
		}
	}

	if temp.Len() > 0 {
		result = append(result, temp.String())
	}
	return result
}
