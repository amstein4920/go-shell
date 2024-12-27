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
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, "$ ")
	for scanner.Scan() {
		readString := scanner.Text()

		strippedReadString := strings.Trim(readString, "\n")
		splitReadString := strings.Split(strippedReadString, " ")
		firstCommandReadString := splitReadString[0]

		switch strings.ToUpper(firstCommandReadString) {
		case EXIT:
			os.Exit(0)
		case ECHO:
			fmt.Println(strings.Join(splitReadString[1:], " "))
		case TYPE:
			fmt.Println(typeCommandFunction(splitReadString[1]))
		default:
			fmt.Printf("%s: command not found\n", strippedReadString)
		}
		fmt.Fprint(os.Stdout, "$ ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input", err)
		os.Exit(1)
	}
}

func typeCommandFunction(input string) string {
	switch strings.ToUpper(input) {
	case EXIT, ECHO, TYPE:
		return fmt.Sprintf("%s is a shell builtin", input)
	default:
		return fmt.Sprintf("%s: not found", input)
	}
}
