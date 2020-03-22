package main

import (
	"fmt"
	"os"

	"github.com/abs-lang/abs/install"
	"github.com/abs-lang/abs/repl"
)

var Version = "1.11.3"

// The ABS interpreter
func main() {
	args := os.Args
	if len(args) == 2 && args[1] == "--version" {
		fmt.Println(Version)
		return
	}

	if len(args) == 3 && args[1] == "get" {
		install.Install(args[2])
		return
	}

	// begin the REPL
	repl.BeginRepl(args, Version)
}
