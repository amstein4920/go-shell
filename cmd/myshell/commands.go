package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func typeBuiltinFunction(input string) string {
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
func commandExecutionFunction(inputStrings []string) string {
	cmd := exec.Command(inputStrings[0], inputStrings[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("%s: command not found\n", inputStrings[0])
	} else {
		// I want to standardize all outputs to have exactly one \n. No more, no less
		return fmt.Sprintln(strings.Trim(string(stdout), "\n"))
	}
}

func cdBuiltinFunction(input string) error {
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
