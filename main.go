package main

import (
	"abs/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s, welcome to the Abs programming language!\n", user.Username)
	fmt.Printf("Type 'quit' when you're done, 'help' if you get lost!\n")
	repl.Start(os.Stdin, os.Stdout)
}
