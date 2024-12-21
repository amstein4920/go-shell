package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, "$ ")
	for scanner.Scan() {
		readString := scanner.Text()

		strippedReadString := strings.Trim(readString, "\n")
		splitReadString := strings.Split(strippedReadString, " ")
		firstCommandReadString := splitReadString[0]

		switch strings.ToUpper(firstCommandReadString) {
		case "EXIT":
			os.Exit(0)
		case "ECHO":
			fmt.Println(strings.Join(splitReadString[1:], " "))
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
