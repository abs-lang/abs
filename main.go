package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/abs-lang/abs/repl"
)

var VERSION = "1.1.0"

// The ABS interpreter
func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	args := os.Args

	if len(args) == 2 && args[1] == "--version" {
		fmt.Println(VERSION)
		return
	}

	// if we're called without arguments,
	// launch the REPL
	if len(args) == 1 || strings.HasPrefix(args[1], "-") {
		fmt.Printf("Hello %s, welcome to the ABS (%s) programming language!\n", user.Username, VERSION)
		fmt.Printf("Type 'quit' when you're done, 'help' if you get lost!\n")
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	// let's parse our argument as a file
	code, err := ioutil.ReadFile(args[1])

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	repl.Run(string(code), false)
}
