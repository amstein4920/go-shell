package main

import (
	"bufio"
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
			fmt.Println(config.typeCommandFunction(splitReadString[1]))
		case PWD:
			wd, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "reading working directory", err)
			}
			fmt.Println(wd)
		default:
			cmd := exec.Command(firstCommandReadString, arguments...)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Printf("%s: command not found\n", firstCommandReadString)
			} else {
				// I want to standardize all outputs to have exactly one \n. No more, no less
				fmt.Println(strings.Trim(string(stdout), "\n"))
			}
		}
		fmt.Fprint(os.Stdout, "$ ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input", err)
		os.Exit(1)
	}
}

func (config *shellConfig) typeCommandFunction(input string) string {
	switch strings.ToUpper(input) {
	case EXIT, ECHO, TYPE, PWD:
		return fmt.Sprintf("%s is a shell builtin", input)
	default:
		pathToCommand, err := exec.LookPath(input)
		if err != nil {
			return fmt.Sprintf("%s: not found", input)
		}
		return fmt.Sprintf("%s is %s", input, pathToCommand)
	}
}
