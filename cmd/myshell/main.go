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
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	readString, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	strippedReadString := strings.Trim(readString, "\n")

	switch {
	default:
		fmt.Printf("%s: command not found\n", strippedReadString)
	}
}
