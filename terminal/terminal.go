package terminal

import (
	"crypto/rand"
	"fmt"
	"math/big"
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
// placeholder
// remove dependencies
// worth renaming repl to runner? and maybe terminal back to repl

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
		tea.NewProgram(initialConsole(user, version, env, runner)),
	}
}

type model struct {
	user            string
	version         string
	runner          Runner
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

func initialConsole(user string, version string, env *object.Environment, r Runner) model {
	in := textinput.New()
	in.Prompt = getPromptPrefix(env)
	in.Placeholder = "`date`"
	historyFile, maxLines := getHistoryConfiguration(env)
	history := getHistory(historyFile, maxLines)
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
		runner:          r,
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

			res := m.runner(m.in.Value())
			m.history = append(m.history, m.in.Value())
			m.historyPoint = len(m.history)

			if res != "" {
				m.messages = append(m.messages, res)
			}
			m.in.Reset()
		case tea.KeyCtrlL:

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

func getPromptPrefix(env *object.Environment) string {
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

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func changeLivePrefix() (string, bool) {
	livePrefix := formatLivePrefix(LivePrefixState.LivePrefix)
	return livePrefix, LivePrefixState.IsEnable
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
const ABS_PROMPT_PREFIX = "⧐  "
