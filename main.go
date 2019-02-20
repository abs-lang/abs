package main

import (
	"fmt"
	"os"

	"github.com/abs-lang/abs/repl"
)

var VERSION = "1.3.0"

// The ABS interpreter
func main() {
	args := os.Args
	if len(args) == 2 && args[1] == "--version" {
		fmt.Println(VERSION)
		return
	}
	// begin the REPL
	repl.BeginRepl(args, VERSION)
}
