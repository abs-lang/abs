package repl

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
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
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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
	env = object.NewEnvironment(os.Stdout, d, "")
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

func getPromptPrefix() string {
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

	return promptPrefix
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

	// create and start the command prompt run loop
	p := prompt.New(
		executor,
		completer,
		// prompt.OptionPrefix(promptPrefix),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("abs-repl"),
		prompt.OptionHistory(history),
		prompt.OptionAddKeyBind(reverseSearch()),           // ControlR: reverse search
		prompt.OptionBreakLineCallback(clearReverseSearch), // at every line break clear the reverse search
	)

	p.Run()
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
func Run(code string, interactive bool) string {
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

	// TODO this should be removed and injected in the environment
	// when the module is initialized
	env.Version = version
	env.Set("ABS_VERSION", &object.String{Value: version})

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

		p := tea.NewProgram(initialConsole(user.Username, version))

		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
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

type (
	errMsg error
)

type model struct {
	user            string
	version         string
	history         []string
	historyPoint    int
	historyFile     string
	historyMaxLInes int
	dirty           string
	messages        []string
	in              textinput.Model
	// senderStyle lipgloss.Style
	err error
}

func initialConsole(user string, version string) model {
	in := textinput.New()
	in.Prompt = getPromptPrefix()
	in.Placeholder = "`date`"
	historyFile, maxLines = getHistoryConfiguration()
	history = getHistory(historyFile, maxLines)
	in.Focus()
	// ti.CharLimit = 156
	// ti.Width = 20
	messages := []string{}
	messages = append(messages, fmt.Sprintf("Hello %s, welcome to the ABS (%s) programming language!", user, version))
	// check for new version about 10% of the time,
	// to avoid too many hangups
	if r, e := rand.Int(rand.Reader, big.NewInt(100)); e == nil && r.Int64() < 10 {
		if newver, update := util.UpdateAvailable(version); update {
			messages = append(messages, fmt.Sprintf(
				"*** Update available: %s (your version is %s) ***",
				newver,
				version,
			))
		}
	}
	messages = append(messages, "Type 'quit' when you're done, 'help' if you get lost!")

	return model{
		in:              in,
		history:         history,
		historyPoint:    len(history),
		historyFile:     historyFile,
		historyMaxLInes: maxLines,
		dirty:           "",
		messages:        messages,
		user:            user,
		version:         version,
		// senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err: nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)

	m.in, tiCmd = m.in.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlD:
			saveHistory(m.historyFile, m.historyMaxLInes, m.history)
			return m, tea.Quit
		case tea.KeyEnter, tea.KeyCtrlC:
			m.dirty = ""
			m.messages = append(m.messages, m.in.Prompt+" "+m.in.Value())

			if m.in.Value() == "" {
				break
			}

			res := Run(m.in.Value(), true)
			m.history = append(m.history, m.in.Value())
			m.historyPoint = len(m.history)

			if res != "" {
				m.messages = append(m.messages, res)
			}
			m.in.Reset()
		case tea.KeyUp:
			// oly store dirty state on first key up
			if m.dirty == "" {
				m.dirty = m.in.Value()
			}

			if m.historyPoint <= 0 {
				break
			}

			newPoint := m.historyPoint - 1

			if len(m.history) < newPoint {
				break
			}

			m.historyPoint = newPoint
			m.in.SetValue(m.history[m.historyPoint])
		case tea.KeyDown:
			newPoint := m.historyPoint + 1

			if newPoint <= len(m.history)-1 {
				m.historyPoint = newPoint
				m.in.SetValue(m.history[m.historyPoint])
				break
			}

			m.in.SetValue(m.dirty)
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tiCmd
}

func formatMsg(msg string) string {
	return lipgloss.NewStyle().Render(msg)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		formatMsg(strings.Join(m.messages, "\n")),
		m.in.View(),
	)
}

// TODO clear working
// history
// cursor up and down
// quit command
// help
// autocompleter
// TODO don't make environment global, not needed anymore, instead pass it to model
// navigate commands up
// print parse errors / generale errors correctly
// > noexist (for example)
// blue color prompt
// stdin() not working
// ctrlc?
// move under terminal package
// remove deprecated terminat package
// maybe only save incrementally in history https://stackoverflow.com/questions/7151261/append-to-a-file-in-go
// remove deprecated ioutil methods
// live prefix
