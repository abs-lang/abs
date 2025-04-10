package terminal

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"

	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/util"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO
// history
// cursor up and down
// autocompleter
// navigate commands up
// print parse errors / generale errors correctly
// > noexist (for example)
// maybe only save incrementally in history https://stackoverflow.com/questions/7151261/append-to-a-file-in-go
// remove deprecated ioutil methods
// remove dependencies
// worth renaming repl to runner? and maybe terminal back to repl
// add prompt formatting tests
// more example statements
// suggestions
// review entire code org

const (
	SIGNAL_TERMINAL_SUSPEND float32 = iota + 1
	SIGNAL_TERMINAL_RESTORE
)

// Runner takes in ABS code and returns
// the programs output after evaluating
// it
type Runner func(string) string

// Channel that can be used to communicate
// with the terminal
type signals chan float32

func NewTerminal(user string, env *object.Environment, runner Runner) *tea.Program {
	signal := make(signals)
	p := tea.NewProgram(getInitialModel(signal, user, env, runner))

	go func() {
		for s := range signal {
			switch s {
			case SIGNAL_TERMINAL_SUSPEND:
				p.ReleaseTerminal()
			case SIGNAL_TERMINAL_RESTORE:
				p.RestoreTerminal()
			}
		}
	}()

	return p
}

func getInitialModel(sigs signals, user string, env *object.Environment, r Runner) Model {
	historyFile, maxLines := getHistoryConfiguration(env)
	history := getHistory(historyFile, maxLines)

	m := Model{
		signals:         sigs,
		user:            user,
		runner:          r,
		env:             env,
		history:         history,
		historyPoint:    len(history),
		historyFile:     historyFile,
		historyMaxLInes: maxLines,
		dirty:           "",
	}

	m.prompt = func() string {
		return getPrompt(m.env)
	}

	// Setup the input line of our terminal
	in := textinput.New()
	in.Prompt = m.prompt()
	in.Placeholder = exampleStatements[mrand.Intn(len(exampleStatements))] + " # just something you can run... (tab + enter)"
	in.Focus()

	m.in = in

	return m
}

type Model struct {
	signals         signals
	user            string
	runner          Runner
	env             *object.Environment
	prompt          func() string
	history         []string
	historyPoint    int
	historyFile     string
	historyMaxLInes int
	dirty           string
	in              textinput.Model
}

func (m Model) Init() tea.Cmd {
	lines := Lines{}
	lines.Add(fmt.Sprintf("Hello %s, welcome to the ABS (%s) programming language!", m.user, m.env.Version))
	lines.Add("Type 'quit' when you're done, 'help' if you get lost!")

	// check for new version about 10% of the time,
	// to avoid too many hangups
	if r, e := rand.Int(rand.Reader, big.NewInt(100)); e == nil && r.Int64() < 10 {
		if newver, update := util.UpdateAvailable(m.env.Version); update {
			lines.Add(lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(
				"\n*** Update available: %s (your version is %s) ***",
				newver,
				m.env.Version,
			)))
		}
	}

	return tea.Batch(
		tea.SetWindowTitle("abs-repl"),
		textarea.Blink,
		tea.Sequence(lines...),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)

	m.in, tiCmd = m.in.Update(msg)

	switch msg := msg.(type) {
	case evalCmd:
		return m.Print(msg)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlD:
			return m.Quit()
		case tea.KeyCtrlC:
			return m.Interrupt()
		case tea.KeyEnter:
			m.in.Placeholder = ""

			if m.in.Value() == "" {
				return m, tea.Println(m.prompt())
			}

			m.history = append(m.history, m.in.Value())
			m.historyPoint = len(m.history)

			switch m.in.Value() {
			case "quit":
				return m.Quit()
			case "help":
				return m.Help()
			default:
				return m.Eval()
			}

		case tea.KeyTab:
			if m.in.Placeholder != "" {
				return m.EngagePlaceholder()
			}
		case tea.KeyCtrlL:
			return m.Clear()
		case tea.KeyUp:
			// oly store dirty Model on first key up
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

	}

	return m, tiCmd
}

func (m Model) View() string {
	return m.in.View()
}

func (m Model) Clear() (Model, tea.Cmd) {
	m.in.Placeholder = ""

	return m, tea.ClearScreen
}

func (m Model) Quit() (Model, tea.Cmd) {
	saveHistory(m.historyFile, m.historyMaxLInes, m.history)

	return m, tea.Quit
}

func (m Model) currentLine() string {
	return m.prompt() + m.in.Value()
}

func (m Model) Help() (Model, tea.Cmd) {
	lines := Lines{}
	prompt := m.prompt()
	help := func(s string) string { return lipgloss.NewStyle().Faint(true).Render(s) }

	lines.Add(m.currentLine())
	lines.Add(help("Try typing something along the lines of:\n"))
	lines.Add("  " + prompt + help("current_date = `date`\n"))
	lines.Add(help("A command should be triggered in your system. Then try printing the result of that command with:\n"))
	lines.Add("  " + prompt + help("current_date\n"))
	lines.Add(help("Here some other valid examples of ABS code:\n"))

	for i := 0; i < 5; i++ {
		ix := mrand.Intn(len(exampleStatements))
		lines.Add("  " + prompt + help(exampleStatements[ix]+"\n"))
	}

	m.in.Reset()

	return m, tea.Sequence(lines...)
}

func (m Model) EngagePlaceholder() (Model, tea.Cmd) {
	m.in.SetValue(m.in.Placeholder)

	return m, nil
}

type evalCmd struct {
	code   string
	result string
}

var ch = make(chan evalCmd)

func (m Model) Eval() (Model, tea.Cmd) {
	m.dirty = ""

	return m, func() tea.Msg {
		m.signals <- SIGNAL_TERMINAL_SUSPEND
		go func() {
			res := evalCmd{m.in.Value(), m.runner(m.in.Value())}
			m.signals <- SIGNAL_TERMINAL_RESTORE
			ch <- res
		}()
		return <-ch
	}
}

func (m Model) Print(msg evalCmd) (Model, tea.Cmd) {
	lines := Lines{}
	lines.Add(m.prompt() + msg.code)

	if msg.result != "" {
		lines.Add(msg.result)
	}

	m.in.Reset()

	return m, tea.Sequence(lines...)
}

func (m Model) Interrupt() (Model, tea.Cmd) {
	l := m.currentLine()
	m.in.Reset()

	return m, tea.Println(l)
}
