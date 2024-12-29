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
		stdOutput := ""

		readString := scanner.Text()

		if readString == "" {
			fmt.Fprint(os.Stdout, "$ ")
			continue
		}

		strippedReadString := strings.Trim(readString, "\n")
		splitReadString := strings.Split(strippedReadString, " ")
		firstCommandReadString := splitReadString[0]
		arguments := splitReadString[1:]
		var secondaryCommand []string

		for index, arg := range arguments {
			if arg == "1>" || arg == ">" {
				arguments = splitReadString[1 : index+1]
				secondaryCommand = splitReadString[index+1:]
			}
		}

		switch strings.ToUpper(firstCommandReadString) {
		case EXIT:
			os.Exit(0)
		case ECHO:
			stdOutput = strings.Trim(strings.Join(arguments, " "), "'") + "\n"
		case TYPE:
			fmt.Println(config.typeBuiltinFunction(splitReadString[1]))
		case PWD:
			wd, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "reading working directory", err)
			}
			stdOutput = wd + "\n"
		case CD:
			err := config.cdBuiltinFunction(arguments[0])
			if err != nil {
				fmt.Println(err)
			}
		default:
			var err error
			stdOutput, err = config.commandExecutionFunction(firstCommandReadString, arguments)
			if err != nil {
				fmt.Print(err)
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
			}
		} else {
			fmt.Fprint(os.Stdout, stdOutput)
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
func (config *shellConfig) commandExecutionFunction(command string, arguments []string) (string, error) {
	var stdout []byte
	var returnError error

	_, err := exec.LookPath(command)
	if err != nil {
		returnError = fmt.Errorf("%s: command not found", command)
	} else {
		cmd := exec.Command(command, arguments...)
		stdout, err = cmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				returnError = errors.New(string(exitErr.Stderr))
			} else {
				returnError = fmt.Errorf("error:%s", err)
			}
		}
	}
	// I want to standardize all outputs to have exactly one \n. No more, no less
	return fmt.Sprintln(strings.Trim(string(stdout), "\n")), returnError
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
