package main

import (
	"bufio"
	"errors"
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
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		stdOutput := ""
		stdError := ""

		input, _ := reader.ReadString('\n')

		if input == "" {
			fmt.Fprint(os.Stdout, "$ ")
			continue
		}

		stripped := strings.Trim(input, "\n")
		parsed := parse(stripped)
		command := parsed[0]
		arguments := parsed[1:]
		var secondaryCommand []string

		for i, arg := range arguments {
			if arg == "1>" || arg == "1>>" || arg == ">" || arg == ">>" || arg == "2>" || arg == "2>>" {
				arguments = parsed[1 : i+1]
				secondaryCommand = parsed[i+1:]
			}
		}

		switch strings.ToUpper(command) {
		case EXIT:
			os.Exit(0)
		case ECHO:
			stdOutput = strings.Join(arguments, " ") + "\n"
		case TYPE:
			var err error
			stdOutput, err = typeBuiltinFunction(parsed[1])
			if err != nil {
				stdError = err.Error() + "\n"
			}
		case PWD:
			wd, err := os.Getwd()
			if err != nil {
				stdError = fmt.Sprintln("reading working directory", err)
			}
			stdOutput = wd + "\n"
		case CD:
			err := cdBuiltinFunction(arguments[0])
			if err != nil {
				stdError = err.Error() + "\n"
			}
		default:
			var err error
			stdOutput, err = commandExecutionFunction(command, arguments)
			if err != nil {
				stdError = err.Error()
			}
		}
		if len(secondaryCommand) > 0 {
			secondCommand := secondaryCommand[0]
			secondArgs := secondaryCommand[1:]
			switch secondCommand {
			case ">", "1>", ">>", "1>>":
				err := writeToFile(secondArgs[0], stdOutput)
				if err != nil {
					fmt.Println("Error writing to file:", err)
				}
				fmt.Fprint(os.Stderr, stdError)
			case "2>", "2>>":
				err := writeToFile(secondArgs[0], stdError)
				if err != nil {
					fmt.Println("Error writing to file:", err)
				}
				fmt.Fprint(os.Stdout, stdOutput)
			}
		} else {
			if len(stdOutput) > 0 {
				fmt.Fprint(os.Stdout, stdOutput)
			}
			if len(stdError) > 0 {
				fmt.Fprint(os.Stderr, stdError)
			}
		}
	}
}

func parse(input string) []string {
	var result []string
	var temp strings.Builder
	singleQuoted := false
	doubleQuoted := false
	escaped := false

	for _, char := range input {
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

func writeToFile(file string, output string) error {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(file, []byte(output), 0o777)
	} else {
		fileContent, fileErr := os.ReadFile(file)
		if fileErr != nil {
			return fileErr
		}
		newFileContent := append(fileContent, []byte(output)...)
		return os.WriteFile(file, newFileContent, 0o777)
	}
}
