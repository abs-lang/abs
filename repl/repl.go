package repl

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/abs-lang/abs/evaluator"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"
	"github.com/abs-lang/abs/util"
	"github.com/c-bata/go-prompt"

	"golang.org/x/crypto/ssh/terminal"
)

// Global environment for the REPL.
//
// We want the environment to be persistent
// across invokations (else how useful would
// the REPL be?), but we also want it to
// be available here so that other features,
// such as suggestions, work by inspecting
// the environment.
var env *object.Environment

// Support for persistent history in interactive REPL
var (
	historyFile string
	maxLines    int
	history     []string
)

func init() {
	d, _ := os.Getwd()
	env = object.NewEnvironment(os.Stdout, d)
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}

	for _, key := range env.GetKeys() {
		s = append(s, prompt.Suggest{Text: key})
	}

	if len(d.GetWordBeforeCursor()) == 0 {
		return nil
	}

	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func changeLivePrefix() (string, bool) {
	livePrefix := formatLivePrefix(LivePrefixState.LivePrefix)
	return livePrefix, LivePrefixState.IsEnable
}

func Start(in io.Reader, out io.Writer) {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println("unable to start the ABS repl (no terminal detected)")
		os.Exit(1)
	}

	// get history file only when interactive REPL is running
	historyFile, maxLines = getHistoryConfiguration()
	history = getHistory(historyFile, maxLines)
	// once we load the history we can setup reverse search
	// which will need to go through the history itself
	initReverseSearch()
	// get prompt prefix template string
	promptPrefix := util.GetEnvVar(env, "ABS_PROMPT_PREFIX", ABS_PROMPT_PREFIX)
	// get live prompt boolean
	livePrompt := util.GetEnvVar(env, "ABS_PROMPT_LIVE_PREFIX", "false")
	if livePrompt == "true" {
		LivePrefixState.LivePrefix = promptPrefix
		LivePrefixState.IsEnable = true
	} else {
		if promptPrefix != formatLivePrefix(promptPrefix) {
			// we have a template string when livePrompt mode is turned off
			// use default static prompt instead
			promptPrefix = ABS_PROMPT_PREFIX
		}
	}

	// create and start the command prompt run loop
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(promptPrefix),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("abs-repl"),
		prompt.OptionHistory(history),
		prompt.OptionAddKeyBind(reverseSearch()),           // ControlR: reverse search
		prompt.OptionBreakLineCallback(clearReverseSearch), // at every line break clear the reverse search
	)

	p.Run()
	// we get here on ^D from the prompt
	saveHistory(historyFile, maxLines, history)
}

// The executor simply reads what
// the user has entered and decides
// what to do with it:
// either quit, print the help message
// or evaluate code.
func executor(line string) {
	if line == "quit" {
		fmt.Printf("%s\n", "Adios!")
		saveHistory(historyFile, maxLines, history)
		os.Exit(0)
	}

	if line == "help" {
		fmt.Printf("Try typing something along the lines of:")
		fmt.Printf("%s", "\n")
		fmt.Printf("%s", "\n")
		fmt.Print("  ⧐  current_date = $(date)")
		fmt.Printf("%s", "\n")
		fmt.Printf("%s", "\n")
		fmt.Print("A command should be triggered in your system. Then try printing the result of that command with:")
		fmt.Printf("%s", "\n")
		fmt.Printf("%s", "\n")
		fmt.Printf("  ⧐  current_date")
		fmt.Printf("%s", "\n")
		return
	}

	// record this line for posterity
	history = addToHistory(history, maxLines, line)
	Run(line, true)
}

// Core of the REPL.
//
// This function takes code and evaluates
// it, spitting out the result.
func Run(code string, interactive bool) {
	lex := lexer.New(code)
	p := parser.New(lex)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		if !interactive {
			os.Exit(99)
		}
		return
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
			return
		}

		if interactive && evaluated.Type() != object.NULL_OBJ {
			fmt.Printf("%s", evaluated.Inspect())
			fmt.Println("")
		}
	}
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
	getAbsInitFile(interactive)

	if interactive {
		// preload the ABS global env with the builtin Fns names
		for k, v := range evaluator.Fns {
			env.Set(k, v)
		}
		// launch the interactive REPL
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s, welcome to the ABS (%s) programming language!\n", user.Username, version)
		fmt.Printf("Type 'quit' when you're done, 'help' if you get lost!\n")
		Start(os.Stdin, os.Stdout)
	} else {
		// this is a script
		// let's parse our argument as a file and run it
		code, err := ioutil.ReadFile(args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(99)
		}

		Run(string(code), false)
	}

}
