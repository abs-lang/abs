package terminal

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strings"

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
// help
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
	in := textinput.New()
	in.Prompt = getPrompt(env)
	in.Placeholder = exampleStatements[mrand.Intn(len(exampleStatements))] + " # just something you can run... (tab + enter)"
	historyFile, maxLines := getHistoryConfiguration(env)
	history := getHistory(historyFile, maxLines)
	in.Focus()
	messages := []string{}
	messages = append(messages, fmt.Sprintf("Hello %s, welcome to the ABS (%s) programming language!", user, env.Version))
	messages = append(messages, "Type 'quit' when you're done, 'help' if you get lost!")

	// check for new version about 10% of the time,
	// to avoid too many hangups
	if r, e := rand.Int(rand.Reader, big.NewInt(100)); e == nil && r.Int64() < 10 {
		if newver, update := util.UpdateAvailable(env.Version); update {
			msg := fmt.Sprintf(
				"\n*** Update available: %s (your version is %s) ***",
				newver,
				env.Version,
			)
			messages = append(messages, lipgloss.NewStyle().Faint(true).Render(msg))
		}
	}

	return Model{
		signals:         sigs,
		user:            user,
		runner:          r,
		prompt:          getPrompt,
		env:             env,
		in:              in,
		history:         history,
		historyPoint:    len(history),
		historyFile:     historyFile,
		historyMaxLInes: maxLines,
		dirty:           "",
		messages:        messages,
		err:             nil,
	}
}

type Model struct {
	signals         signals
	user            string
	runner          Runner
	env             *object.Environment
	prompt          func(*object.Environment) string
	history         []string
	historyPoint    int
	historyFile     string
	historyMaxLInes int
	dirty           string
	messages        []string
	in              textinput.Model
	err             error
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("abs-repl"), textarea.Blink)
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
			if m.in.Value() == "quit" {
				return m.Quit()
			}

			return m.Eval()
		case tea.KeyTab:
			return m.EngagePlaceholder()
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
	if len(m.messages) == 0 {
		return m.in.View()
	}

	return fmt.Sprintf(
		"%s\n%s",
		strings.Join(m.messages, "\n"),
		m.in.View(),
	)
}

func (m Model) Clear() (Model, tea.Cmd) {
	m.messages = []string{}
	m.in.Placeholder = ""

	return m, tea.ClearScreen
}

func (m Model) Quit() (Model, tea.Cmd) {
	saveHistory(m.historyFile, m.historyMaxLInes, m.history)

	return m, tea.Quit
}

func (m Model) EngagePlaceholder() (Model, tea.Cmd) {
	if m.in.Placeholder != "" {
		m.in.SetValue(m.in.Placeholder)
	}

	return m, nil
}

type evalCmd struct {
	code   string
	result string
}

var ch = make(chan evalCmd)

func (m Model) Eval() (Model, tea.Cmd) {
	m.in.Placeholder = ""
	m.dirty = ""

	if m.in.Value() == "" {
		return m.Print(evalCmd{})
	}

	m.history = append(m.history, m.in.Value())
	m.historyPoint = len(m.history)

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
	m.messages = append(m.messages, m.prompt(m.env)+msg.code)
	s := msg.result

	if s != "" {
		m.messages = append(m.messages, s)
	}

	m.in.Prompt = m.prompt(m.env)
	m.in.Placeholder = ""
	m.in.Reset()

	return m, nil
}

func (m Model) Interrupt() (Model, tea.Cmd) {
	m.messages = append(m.messages, m.prompt(m.env)+m.in.Value())
	m.in.Reset()

	return m, nil
}
