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

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		stdOutput := ""
		stdError := ""

		readString, _ := reader.ReadString('\n')

		if readString == "" {
			fmt.Fprint(os.Stdout, "$ ")
			continue
		}

		strippedReadString := strings.Trim(readString, "\n")
		splitReadString := strings.Split(strippedReadString, " ")
		firstCommandReadString := splitReadString[0]
		arguments := splitReadString[1:]
		var secondaryCommand []string

		for i, arg := range arguments {
			if arg == "1>" || arg == ">" || arg == "2>" {
				arguments = splitReadString[1 : i+1]
				secondaryCommand = splitReadString[i+1:]
			}
		}

		switch strings.ToUpper(firstCommandReadString) {
		case EXIT:
			os.Exit(0)
		case ECHO:
			stdOutput = strings.Trim(strings.Join(arguments, " "), "'") + "\n"
		case TYPE:
			var err error
			stdOutput, err = config.typeBuiltinFunction(splitReadString[1])
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
			err := config.cdBuiltinFunction(arguments[0])
			if err != nil {
				stdError = err.Error() + "\n"
			}
		default:
			var err error
			stdOutput, err = config.commandExecutionFunction(firstCommandReadString, arguments)
			if err != nil {
				stdError = err.Error()
			}
		}
		if len(secondaryCommand) > 0 {
			secondCommand := secondaryCommand[0]
			secondArgs := secondaryCommand[1:]
			switch secondCommand {
			case ">", "1>":
				err := config.writeToFile(secondArgs[0], stdOutput)
				if err != nil {
					fmt.Println("Error writing to file:", err)
				}
				fmt.Fprint(os.Stderr, stdError)
			case "2>":
				err := config.writeToFile(secondArgs[0], stdError)
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

func (config *shellConfig) typeBuiltinFunction(input string) (string, error) {
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

// Executes provided command and returns the standard output string. The command not found message also catches for builtins
func (config *shellConfig) commandExecutionFunction(command string, arguments []string) (string, error) {
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

func (config *shellConfig) writeToFile(file string, output string) error {
	err := os.WriteFile(file, []byte(output), 0o777)
	return err
}
