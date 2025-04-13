package terminal

import (
	"os"
	"os/user"
	"strings"

	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const ABS_DEFAULT_PROMPT = "> "

var exampleStatements = []string{
	"`ls -la`",
	"`cat /etc/hosts`",
	"['a', 'b', 'c'].map(f(l) {l.upper()})",
	"1..10",
	"1 in [0,1,2,3,4]",
	"'string' ~ 'sTrINg'",
	"true || sleep(1000)",
	"true && sleep(1000)",
	"one, two, three = [1, 2, 3]",
	"!!true",
	"10 ** 2",
	"6 <=> 5",
	"{'x': 1}?.x?.x",
	"defer echo(1); echo(2)",
	"\"hello world\"[-2]",
	"\"hello world\"[:5]",
	"\"hello %s\".fmt(\"world\")",
	"`cat /etc/hosts`.lines()",
	"10.3.ceil()",
	"[1, 2] + [3]",
	"[{'name': 'Lebron', 'age': 40}, {'name': 'Michael', 'age': 'older...'}].tsv()",
	"f greeter(greeting = 'hello'){ '%s world'.fmt(greeting) }",
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

type Lines []string

func (ls *Lines) Add(l string) {
	*ls = append(*ls, l)
}

func (ls *Lines) Dump() tea.Cmd {
	return tea.Println(ls.Join())
}

func (ls *Lines) Join() string {
	return strings.Join(*ls, "\n")
}

func applySuggestion(s, textToReplace, suggestion string) string {
	wstart := strings.LastIndex(s, textToReplace)
	wend := wstart + len(textToReplace) - 1
	head := s[:wstart]
	tail := ""

	if len(s) > wend {
		tail = s[wend+1:]
	}

	return head + suggestion + tail
}
