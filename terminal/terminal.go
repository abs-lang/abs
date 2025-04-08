package terminal

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"os"
	"os/user"
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
// quit command
// help
// autocompleter
// navigate commands up
// print parse errors / generale errors correctly
// > noexist (for example)
// stdin() not working
// sleep blocks everything
// maybe only save incrementally in history https://stackoverflow.com/questions/7151261/append-to-a-file-in-go
// remove deprecated ioutil methods
// remove dependencies
// worth renaming repl to runner? and maybe terminal back to repl
// add prompt formatting tests
// more example statements

type (
	errMsg error
)

type Terminal struct {
	program *tea.Program
}

func (t *Terminal) Run() error {
	_, err := t.program.Run()

	return err
}

type Runner func(string) string

func New(user string, version string, env *object.Environment, runner Runner) *Terminal {
	return &Terminal{
		tea.NewProgram(getInitialState(user, version, env, runner)),
	}
}

type Model struct {
	user            string
	version         string
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

func (m Model) Clear() (Model, tea.Cmd) {
	m.messages = []string{}
	m.in.Placeholder = ""
	return m, tea.ClearScreen
}

func (m Model) Quit() (Model, tea.Cmd) {
	saveHistory(m.historyFile, m.historyMaxLInes, m.history)

	return m, tea.Quit
}

func (m Model) Eval() (Model, tea.Cmd) {
	m.in.Placeholder = ""
	m.dirty = ""
	m.messages = append(m.messages, m.prompt(m.env)+m.in.Value())

	if m.in.Value() == "" {
		return m, nil
	}

	res := m.runner(m.in.Value())
	m.history = append(m.history, m.in.Value())
	m.historyPoint = len(m.history)

	if res != "" {
		m.messages = append(m.messages, res)
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

var exampleStatements = []string{
	"`ls -la`",
	"`cat /etc/hosts`",
	"['a', 'b', 'c'].map(f(l) {l.upper()})",
	"1..10",
	"1 in [0,1,2,3,4]",
	"'string' ~ 'sTrINg'",
	"true || sleep(1000)",
	"true && sleep(1000)",
}

func getInitialState(user string, version string, env *object.Environment, r Runner) Model {
	in := textinput.New()
	in.Prompt = getPrompt(env)
	in.Placeholder = exampleStatements[mrand.Intn(len(exampleStatements))] + " # just something you can run..."
	historyFile, maxLines := getHistoryConfiguration(env)
	history := getHistory(historyFile, maxLines)
	in.Focus()
	messages := []string{}
	messages = append(messages, fmt.Sprintf("Hello %s, welcome to the ABS (%s) programming language!", user, version))
	messages = append(messages, "Type 'quit' when you're done, 'help' if you get lost!")

	// check for new version about 10% of the time,
	// to avoid too many hangups
	if r, e := rand.Int(rand.Reader, big.NewInt(100)); e == nil && r.Int64() < 10 {
		if newver, update := util.UpdateAvailable(version); update {
			msg := fmt.Sprintf(
				"\n*** Update available: %s (your version is %s) ***",
				newver,
				version,
			)
			messages = append(messages, lipgloss.NewStyle().Faint(true).Render(msg))
		}
	}

	return Model{
		user:            user,
		version:         version,
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

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("abs-repl"), textarea.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)

	m.in, tiCmd = m.in.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlD:
			return m.Quit()
		case tea.KeyCtrlC:
			return m.Interrupt()
		case tea.KeyEnter:
			return m.Eval()
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

	case errMsg:
		m.err = msg
		return m, nil
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

func getPrompt(env *object.Environment) string {
	prompt := util.GetEnvVar(env, "ABS_PROMPT_PREFIX", ABS_DEFAULT_PROMPT)
	prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#4287f5")).Render(prompt)
	livePrompt := util.GetEnvVar(env, "ABS_PROMPT_LIVE_PREFIX", "false")

	if livePrompt == "true" {
		return formatLivePrefix(prompt)
	}

	return prompt
}

// format ABS_PROMPT_PREFIX = "{user}@{host}:{dir} $"
func formatLivePrefix(prefix string) string {
	livePrefix := prefix
	if strings.Contains(prefix, "{") {
		userInfo, _ := user.Current()
		user := userInfo.Username
		host, _ := os.Hostname()
		dir, _ := os.Getwd()
		// shorten homedir to ~/
		homeDir := userInfo.HomeDir
		dir = strings.Replace(dir, homeDir, "~", 1)
		// format the livePrefix
		livePrefix = strings.Replace(livePrefix, "{user}", user, 1)
		livePrefix = strings.Replace(livePrefix, "{host}", host, 1)
		livePrefix = strings.Replace(livePrefix, "{dir}", dir, 1)
	}
	return livePrefix
}

// support for user config of ABS REPL prompt string
var ABS_DEFAULT_PROMPT = "> "
