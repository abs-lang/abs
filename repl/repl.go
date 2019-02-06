package repl

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/abs-lang/abs/evaluator"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"

	prompt "github.com/c-bata/go-prompt"
)

var env *object.Environment

// Global environment for the REPL.
//
// We want the environment to be persistent
// across invokations (else how useful would
// the REPL be?), but we also want it to
// be available here so that other features,
// such as suggestions, work by inspecting
// the environment.
func init() {
	env = object.NewEnvironment()
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}

	for _, key := range env.GetKeys() {
		s = append(s, prompt.Suggest{Text: key})
	}

	if len(d.GetWordBeforeCursor()) == 0 {
		return nil
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}

func Start(in io.Reader, out io.Writer) {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("⧐  "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("abs-repl"),
		prompt.OptionHistory(getHistory()),
	)
	p.Run()
	// we get here on ^D from the prompt
	saveHistory()
}

/*
Support for abs history file in the interactive REPL:

1) The current ABS_HISTORY_FILE is loaded into the prompt.Run() cycle
   using prompt.OptionHistory(getHistory()). This also loads the local history as well.
   Default ABS_HISTORY_FILE is "~/.abs_history".
2) Append each non-null, unique next line passed from prompt to the executor() to the local history.
   NB. the live prompt history shows duplicate next lines, but they are not saved to the local history.
3) Save the local history whenever the prompt.Run() loop exits (^D) or the executor() exits (on quit).
   Write the local history to the ABS_HISTORY_FILE up to ABS_MAX_HISTORY_LINES (default 1000 lines).

Note that ABS_HISTORY_FILE and ABS_MAX_HISTORY_LINES variables may come from the OS environment.
*/

const (
	ABS_HISTORY_FILE      = "~/.abs_history"
	ABS_MAX_HISTORY_LINES = "1000"
)

var (
	historyFilePath string   // fully expanded ABS_HISTORY_FILE path
	historyMaxLines int      // ABS_MAX_HISTORY_LINES
	history         []string // local history
)

// expand full path to ABS_HISTORY_FILE for current user and get ABS_MAX_HISTORY_LINES
func expandHistoryFilePath() string {
	// obtain any OS environment variablse
	maxHistoryLines := os.Getenv("ABS_MAX_HISTORY_LINES")
	if len(maxHistoryLines) == 0 {
		maxHistoryLines = ABS_MAX_HISTORY_LINES
	}
	historyMaxLines, _ = strconv.Atoi(maxHistoryLines)
	//
	historyFile := os.Getenv("ABS_HISTORY_FILE")
	if len(historyFile) == 0 {
		historyFile = ABS_HISTORY_FILE
	}
	// identify the user's homeDir
	user, ok := user.Current()
	if ok != nil {
		fmt.Println("Unable to resolve current userId for ~/.abs_history")
		os.Exit(99)
	}
	// Default path to user's homeDir
	path := user.HomeDir
	// expand bash-style ~/ path prefix to homeDir (also works in windows)
	if strings.HasPrefix(historyFile, "~/") {
		path = strings.Replace(historyFile, "~/", path+"/", 1)
	} else if len(historyFile) > 0 {
		path = historyFile
	}
	return path
}

// read the history file and split it into the local history[...] slice
func getHistory() []string {
	// consult the OS environment
	historyFilePath = expandHistoryFilePath()
	if historyMaxLines == 0 {
		// do not open a history file for zero max lines
		return history
	}
	// verify the expanded historyFile exists, if not create it now
	fd, _ := os.OpenFile(historyFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	fd.Close()
	// read the file and split the lines into history[...]
	bytes, err := ioutil.ReadFile(historyFilePath)
	if err != nil {
		return history
	}
	// fill the local history from the file
	if len(bytes) > 0 {
		history = strings.Split(string(bytes), "\n")
	}
	return history
}

// append unique next line to local history[...]
func addToHistory(line string) {
	if historyMaxLines == 0 {
		// do not save history for zero max lines
		return
	}
	// do not save null lines nor duplicate the previous line in local history
	// NB. this is not the prompt.history which shows all added lines
	if len(line) > 0 {
		if len(history) == 0 {
			history = append(history, line)
		} else if line != history[len(history)-1] {
			history = append(history, line)
		}
	}
}

// save the local history containing ABS_MAX_HISTORY_LINES to historyFilePath
func saveHistory() {
	if historyMaxLines == 0 {
		// do not save a history file for zero max lines
		return
	}
	if len(history) > historyMaxLines {
		// remove the excess lines from the front of the history slice
		history = history[len(history)-historyMaxLines:]
	}
	// write the augmented local history back out to the file
	historyStr := strings.Join(history, "\n")
	ioutil.WriteFile(historyFilePath, []byte(historyStr), 0664)
}

// The executor simply reads what
// the user has entered and decides
// what to do with it:
// either quit, print the help message
// or evaluate code.
func executor(line string) {
	if line == "quit" {
		fmt.Printf("%s\n", "Adios!")
		saveHistory()
		os.Exit(0)
	}

	if line == "help" {
		fmt.Printf("Try typing something along the lines of:\n\n")
		fmt.Print("  ⧐  current_date = $(date)\n\n")
		fmt.Print("A command should be triggered in your system. Then try printing the result of that command with:\n\n")
		fmt.Printf("  ⧐  current_date\n")
		return
	}

	// record this line for posterity
	addToHistory(line)

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
