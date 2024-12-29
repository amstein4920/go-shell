package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	EXIT = "EXIT"
	ECHO = "ECHO"
	TYPE = "TYPE"
	PWD  = "PWD"
	CD   = "CD"
)

type shellConfig struct {
}

func main() {
	config := shellConfig{}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, "$ ")
	for scanner.Scan() {
		readString := scanner.Text()

		if readString == "" {
			fmt.Fprint(os.Stdout, "$ ")
			continue
		}

		strippedReadString := strings.Trim(readString, "\n")
		splitReadString := strings.Split(strippedReadString, " ")
		firstCommandReadString := splitReadString[0]
		arguments := splitReadString[1:]

		switch strings.ToUpper(firstCommandReadString) {
		case EXIT:
			os.Exit(0)
		case ECHO:
			fmt.Println(strings.Join(arguments, " "))
		case TYPE:
			fmt.Println(config.typeBuiltinFunction(splitReadString[1]))
		case PWD:
			wd, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "reading working directory", err)
			}
			fmt.Println(wd)
		case CD:
			err := config.cdBuiltinFunction(arguments[0])
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Print(config.commandExecutionFunction(splitReadString))
		}
		fmt.Fprint(os.Stdout, "$ ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input", err)
		os.Exit(1)
	}
}

func (config *shellConfig) typeBuiltinFunction(input string) string {
	switch strings.ToUpper(input) {
	case EXIT, ECHO, TYPE, PWD, CD:
		return fmt.Sprintf("%s is a shell builtin", input)
	default:
		pathToCommand, err := exec.LookPath(input)
		if err != nil {
			return fmt.Sprintf("%s: not found", input)
		}
		return fmt.Sprintf("%s is %s", input, pathToCommand)
	}
}

// Executes provided command and returns the standard output string. The command not found message also catches for builtins
func (config *shellConfig) commandExecutionFunction(inputStrings []string) string {
	cmd := exec.Command(inputStrings[0], inputStrings[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("%s: command not found\n", inputStrings[0])
	} else {
		// I want to standardize all outputs to have exactly one \n. No more, no less
		return fmt.Sprintln(strings.Trim(string(stdout), "\n"))
	}
}

func (config *shellConfig) cdBuiltinFunction(input string) error {
	inputCopy := strings.Clone(input)
	if inputCopy == "~" {
		var err error
		inputCopy, err = os.UserHomeDir()
		if err != nil {
			return errors.New("cd: invalid home")
		}
	}
	err := os.Chdir(inputCopy)
	if err != nil {
		return fmt.Errorf("cd: %s: No such file or directory", inputCopy)
	}
	return nil
}
