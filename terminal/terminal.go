package terminal

import (
	"crypto/rand"
	"fmt"
	"io"
	"maps"
	"math/big"
	mrand "math/rand"
	"os"
	"os/user"
	"slices"
	"strings"

	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/runner"
	"github.com/abs-lang/abs/util"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO
// autocompleter
// maybe only save incrementally in history https://stackoverflow.com/questions/7151261/append-to-a-file-in-go ?
// worth renaming repl to runner? and maybe terminal back to repl
// add prompt formatting tests
// more example statements

var debug = os.Getenv("DEBUG") == "1"

func NewTerminal(env *object.Environment, stdinRelay io.Writer) *tea.Program {
	historyFile, maxLines := getHistoryConfiguration(env)
	history := getHistory(historyFile, maxLines)

	// Setup the input line of our terminal
	prompt := func() string {
		return getPrompt(env)
	}
	in := textinput.New()
	in.Prompt = prompt()
	in.Placeholder = exampleStatements[mrand.Intn(len(exampleStatements))] + " # just something you can run... (tab + enter)"
	in.Focus()

	m := Model{
		in:              in,
		env:             env,
		stdinRelay:      stdinRelay,
		prompt:          prompt,
		history:         history,
		historyIndex:    len(history) - 1,
		historyFile:     historyFile,
		historyMaxLInes: maxLines,
	}

	p := tea.NewProgram(m)

	return p
}

// Terminal state
type Model struct {
	// environment code should be ran on
	env *object.Environment
	// our terminal hijacks OS stdio
	// so functions like ABS reading
	// from stdin won't work by default
	// (because bubbletea hogs stdin).
	// We instead create a relay used to
	// forward stdin events from terminal
	// to abs' stdin.
	stdinRelay io.Writer
	// flag to know whether ABS is executing
	// code or not -- for example, this is used
	// to determine that while ABS is executing,
	// we should relay stdin from the terminal
	isEvaluating bool
	// function to print the prompt 'prefix'
	prompt func() string
	// dirty input -- input I may have typed on
	// the terminal but not yet submitted -- this
	// is primarily used to make sure you can navigate
	// to history and come back to the command you
	// were about to type
	dirtyInput string
	// input field to type all of ABS' goodness!
	in              textinput.Model
	history         []string
	historyIndex    int
	historyFile     string
	historyMaxLInes int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("abs-repl"),
		textarea.Blink,
		tea.Sequence(m.welcome()...),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)

	m.in, tiCmd = m.in.Update(msg)

	switch msg := msg.(type) {
	case doneEval:
		return m.onDoneEval(msg)
	case tea.KeyMsg:
		// the REPL is evaluating ABS code,
		// so if we type during this time,
		// we should forward this to ABS' stdin
		if m.isEvaluating {
			return m.interceptStdin(msg)
		}
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlD:
			return m.quit()
		case tea.KeyCtrlC:
			m = m.resetHistory()
			return m.interrupt()
		case tea.KeyEnter:
			// Let's get rid of the placeholder
			// first time user submits something
			m.in.Placeholder = ""

			// The user submitted empty code.
			// Just print a new line and continue...
			if m.in.Value() == "" {
				return m, tea.Println(m.prompt())
			}

			// We have something submitted, let's add
			// it to the history, only if it's a duplicate
			// of the last entry
			if m.maxHistoryIndex() < 0 || m.history[m.historyIndex] != m.in.Value() {
				m.history = append(m.history, m.in.Value())
			}

			m = m.resetHistory()

			switch m.in.Value() {
			case "quit":
				return m.quit()
			case "help":
				return m.help()
			default:
				return m.eval()
			}
		case tea.KeyTab:
			// If the placeholder code is shown,
			// allow the user to run it by tabbing
			if m.in.Placeholder != "" {
				return m.engagePlaceholder()
			}
		case tea.KeyCtrlL:
			return m.clear()
		case tea.KeyUp:
			m = m.prevHistory()
		case tea.KeyDown:
			m = m.nextHistory()
		}

	}

	return m, tiCmd
}

func (m Model) maxHistoryIndex() int {
	return len(m.history) - 1
}

func (m Model) prevHistory() Model {
	if m.historyIndex < 0 {
		return m
	}

	// Only save dirty state on the first
	// up press
	if m.historyIndex == m.maxHistoryIndex() {
		m.dirtyInput = m.in.Value()
	}

	m.in.SetValue(m.history[m.historyIndex])

	ix := m.historyIndex - 1
	m.historyIndex = ix

	return m
}

func (m Model) nextHistory() Model {
	newPoint := m.historyIndex + 1

	if newPoint <= m.maxHistoryIndex() {
		m.historyIndex = newPoint
		m.in.SetValue(m.history[m.historyIndex])

		return m
	}

	// We reached the end of history,
	// if we had a dirty value, let's use it
	m.in.SetValue(m.dirtyInput)
	return m
}

func (m Model) resetHistory() Model {
	m.dirtyInput = ""
	m.historyIndex = m.maxHistoryIndex()

	return m
}

func (m Model) View() string {
	view := m.in.View()

	if debug {
		m := m.asMap()
		view += "\n"
		for _, k := range slices.Sorted(maps.Keys(m)) {
			view += fmt.Sprintf(faintStyle("\n  %s: %v"), k, m[k])
		}
	}

	return view
}

func (m Model) asMap() map[string]any {
	return map[string]any{
		"history_index":     m.historyIndex,
		"max_history_index": m.maxHistoryIndex(),
		"dirty_input":       m.dirtyInput,
		"is_evaluating":     m.isEvaluating,
	}
}

func faintStyle(s string) string {
	return lipgloss.NewStyle().Faint(true).Render(s)
}

func (m Model) welcome() []tea.Cmd {
	u, err := user.Current()
	username := u.Username

	if err != nil {
		username = "there"
	}

	lines := Lines{}
	lines.Add(fmt.Sprintf("Hello %s, welcome to the ABS (%s) programming language!", username, m.env.Version))
	lines.Add("Type 'quit' when you're done, 'help' if you get lost!")

	// check for new version about 10% of the time,
	// to avoid too many hangups
	if r, e := rand.Int(rand.Reader, big.NewInt(100)); e == nil && r.Int64() < 10 {
		if newver, update := util.UpdateAvailable(m.env.Version); update {
			lines.Add(faintStyle(fmt.Sprintf(
				"\n*** Update available: %s (your version is %s) ***",
				newver,
				m.env.Version,
			)))
		}
	}

	return lines
}

func (m Model) onDoneEval(res doneEval) (Model, tea.Cmd) {
	errfmt := func(s string) string { return lipgloss.NewStyle().Foreground(lipgloss.Color("#ed4747")).Render(s) }
	m.isEvaluating = false

	lines := Lines{}
	lines.Add(m.prompt() + m.in.Value())

	if len(res.parseErrors) > 0 {
		lines.Add(errfmt(fmt.Sprintf(
			"encountered %d syntax errors:\n",
			len(res.parseErrors),
		)))

		for _, e := range res.parseErrors {
			ls := strings.Split(e, "\n")

			for i, l := range ls {
				prefix := ""

				if i == 0 {
					prefix = fmt.Sprintf("%d) ", i+1)
				}
				lines.Add(errfmt("  " + prefix + l))
			}
		}
	}

	if res.out != object.NULL {
		out := res.out.Inspect()

		if !res.ok {
			out = errfmt(out)
		}

		lines.Add(out)
	}

	m.in.Reset()

	return m, tea.Sequence(lines...)
}

func (m Model) interceptStdin(msg tea.KeyMsg) (Model, tea.Cmd) {
	if msg.String() == "enter" {
		m.stdinRelay.Write([]byte{'\n'})
		return m, nil
	}

	m.stdinRelay.Write([]byte(string(msg.Runes)))
	return m, nil
}

func (m Model) clear() (Model, tea.Cmd) {
	m.in.Placeholder = ""

	return m, tea.ClearScreen
}

func (m Model) quit() (Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	err := saveHistory(m.historyFile, m.historyMaxLInes, m.history)

	if err != nil {
		cmds = append(cmds, tea.Println(fmt.Sprintf(
			"Cannot write to ABS history file (%s): %s",
			m.historyFile,
			err.Error(),
		)))
	}

	cmds = append(cmds, tea.Quit)

	return m, tea.Sequence(cmds...)
}

func (m Model) currentLine() string {
	return m.prompt() + m.in.Value()
}

func (m Model) help() (Model, tea.Cmd) {
	lines := Lines{}
	prompt := m.prompt()
	help := func(s string) string { return lipgloss.NewStyle().Faint(true).Render(s) }
	code := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Faint(true).Render(s)
	}

	lines.Add(m.currentLine())
	lines.Add(help("Try typing something along the lines of:\n"))
	lines.Add("  " + prompt + code("current_date = `date`\n"))
	lines.Add(help("A command should be triggered in your system. Then try printing the result of that command with:\n"))
	lines.Add("  " + prompt + code("current_date\n"))
	lines.Add(help("Here some other valid examples of ABS code:\n"))

	for i := 0; i < 5; i++ {
		ix := mrand.Intn(len(exampleStatements))
		lines.Add("  " + prompt + code(exampleStatements[ix]+"\n"))
	}

	m.in.Reset()

	return m, tea.Sequence(lines...)
}

func (m Model) engagePlaceholder() (Model, tea.Cmd) {
	m.in.SetValue(m.in.Placeholder)

	return m, nil
}

type doneEval struct {
	out         object.Object
	ok          bool
	parseErrors []string
}

func (m Model) eval() (Model, tea.Cmd) {
	m.isEvaluating = true
	done := make(chan doneEval)

	go func() {
		out, ok, parseErrors := runner.Run(m.in.Value(), m.env)
		done <- doneEval{out, ok, parseErrors}
	}()

	return m, func() tea.Msg {
		return <-done
	}
}

func (m Model) interrupt() (Model, tea.Cmd) {
	l := m.currentLine()
	m.in.Reset()

	return m, tea.Println(l)
}
