package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func typeBuiltinFunction(input string) (string, error) {
	switch strings.ToUpper(input) {
	case EXIT, ECHO, TYPE, PWD, CD:
		return fmt.Sprintf("%s is a shell builtin\n", input), nil
	default:
		pathToCommand, err := exec.LookPath(input)
		if err != nil {
			return "", fmt.Errorf("%s: not found", input)
		}
		return fmt.Sprintf("%s is %s\n", input, pathToCommand), nil
	}
}

func commandExecutionFunction(command string, arguments []string) (string, error) {
	var stdout []byte
	var returnError error

	_, err := exec.LookPath(command)
	if err != nil {
		// This really should probably be output to stdErr, but I struggled to find a way to get that to function
		return command + ": command not found\n", nil
	}
	cmd := exec.Command(command, arguments...)
	stdout, err = cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			returnError = errors.New(string(exitErr.Stderr))
		} else {
			returnError = fmt.Errorf("error:%s", err)
		}
	}
	return string(stdout), returnError
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
