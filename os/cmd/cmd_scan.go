package cmd

import (
	"bufio"
	"fmt"
	"os"

	"utils/text/str"
)

// Scan prints <info> to stdout, reads and returns user input, which stops by '\n'.
func Scan(info ...interface{}) string {
	fmt.Print(info...)
	return readline()
}

// Scanf prints <info> to stdout with <format>, reads and returns user input, which stops by '\n'.
func Scanf(format string, info ...interface{}) string {
	fmt.Printf(format, info...)
	return readline()
}

func readline() string {
	var s string
	reader := bufio.NewReader(os.Stdin)
	s, _ = reader.ReadString('\n')
	s = str.Trim(s)
	return s
}
