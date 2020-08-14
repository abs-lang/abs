package main

import (
	"fmt"
	"os"

	"github.com/abs-lang/abs/install"
	"github.com/abs-lang/abs/repl"
	"github.com/abs-lang/abs/util"
)

// Version of the ABS interpreter
var Version = "dev"

// The ABS interpreter
func main() {
	args := os.Args
	if len(args) == 2 && args[1] == "--version" {
		fmt.Println(Version)
		return
	}

	if len(args) == 2 && args[1] == "--check-update" {
		if newver, update := util.UpdateAvailable(Version); update {
			fmt.Printf("Update available: %s (your version is %s)\n", newver, Version)
			os.Exit(1)
		} else {
			return
		}
	}

	if len(args) == 3 && args[1] == "get" {
		install.Install(args[2])
		return
	}

	// begin the REPL
	repl.BeginRepl(args, Version)
}
