package repl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/abs-lang/abs/evaluator"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"
	"github.com/abs-lang/abs/terminal"
	"github.com/abs-lang/abs/util"
)

// support for ABS init file
const ABS_INIT_FILE = "~/.absrc"

func getAbsInitFile(env *object.Environment) {
	// get ABS_INIT_FILE from OS environment or default
	initFile := os.Getenv("ABS_INIT_FILE")
	if len(initFile) == 0 {
		initFile = ABS_INIT_FILE
	}
	// expand the ABS_INIT_FILE to the user's HomeDir
	filePath, err := util.ExpandPath(initFile)
	if err != nil {
		fmt.Printf("Unable to expand ABS init file path: %s\nError: %s\n", initFile, err.Error())
		os.Exit(99)
	}
	initFile = filePath
	// read and eval the abs init file
	code, err := ioutil.ReadFile(initFile)
	if err != nil {
		// abs init file is optional -- nothing to do here
		return
	}
	Run(string(code), env)
}

// Core of the REPL.
//
// This function takes code and evaluates
// it, spitting out the result.
func Run(code string, env *object.Environment) string {
	// let's check if this REPL is interactive
	v, _ := env.Get("ABS_INTERACTIVE")
	interactive := v == object.TRUE

	lex := lexer.New(code)
	p := parser.New(lex)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		if !interactive {
			os.Exit(99)
		}
		return ""
	}

	// invoke BeginEval() passing in the program, env, and lexer for error position
	// NB. Eval(node, env) is recursive so we can't call it directly
	evaluated := evaluator.BeginEval(program, env, lex)

	if evaluated != nil {
		isError := evaluated.Type() == object.ERROR_OBJ

		if isError {
			fmt.Printf("%s", evaluated.Inspect())
			fmt.Println("")

			if !interactive {
				os.Exit(99)
			}
			return ""
		}

		if interactive && evaluated.Type() != object.NULL_OBJ {
			return evaluated.Inspect()
		}
	}

	return ""
}

func printParserErrors(errors []string) {
	fmt.Printf("%s", " parser errors:\n")
	for _, msg := range errors {
		fmt.Printf("%s", "\t"+msg+"\n")
	}
}

// BeginRepl (args) -- the REPL, both interactive and script modes begin here
// This allows us to prime the global env with ABS_INTERACTIVE = true/false,
// load the builtin Fns names for the use of command completion, and
// load the ABS_INIT_FILE into the global env
func BeginRepl(args []string, version string) {
	d, _ := os.Getwd()
	env := object.NewEnvironment(os.Stdout, d, "")
	env.Version = version
	env.Set("ABS_VERSION", &object.String{Value: version})

	// if we're called without arguments, this is interactive REPL, otherwise a script
	var interactive bool
	if len(args) == 1 || strings.HasPrefix(args[1], "-") {
		interactive = true
		env.Set("ABS_INTERACTIVE", evaluator.TRUE)
	} else {
		interactive = false
		env.Set("ABS_INTERACTIVE", evaluator.FALSE)
		// Make sure we set the right Dir when evaluating a script,
		// so that the script thinks it's running from its location
		// and things like relative require() calls work.
		env.Dir = filepath.Dir(args[1])
	}

	// get abs init file
	// user may test ABS_INTERACTIVE to decide what code to run
	getAbsInitFile(env)

	// This is a terminal / actual REPL
	if interactive {
		// launch the interactive terminal
		user, err := user.Current()
		if err != nil {
			panic(err)
		}

		term := terminal.New(
			user.Username,
			version,
			env,
			func(code string) string {
				return Run(code, env)
			},
		)

		if err := term.Run(); err != nil {
			log.Fatal(err)
		}

		return
	}

	// this is a script
	// let's parse our argument as a file and run it
	code, err := ioutil.ReadFile(args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(99)
	}

	Run(string(code), env)
}
