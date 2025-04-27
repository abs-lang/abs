package terminal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/util"
)

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
	ABS_MAX_HISTORY_LINES = "10000"
)

// Expand full path to ABS_HISTORY_FILE for current user and get ABS_MAX_HISTORY_LINES
// 1) we look in the ABS global environment as these vars can be set by the ABS init file
// 2) we look in the OS environment
// 3) we use the constant defaults
func getHistoryConfiguration(env *object.Environment) (string, int) {
	// obtain any ABS global environment vars or OS environment vars
	// ABS_MAX_HISTORY_LINES
	maxHistoryLines := util.GetEnvVar(env, "ABS_MAX_HISTORY_LINES", ABS_MAX_HISTORY_LINES)
	maxLines, err := strconv.Atoi(maxHistoryLines)
	if err != nil {
		maxLines, _ = strconv.Atoi(ABS_MAX_HISTORY_LINES)
		fmt.Printf("ABS_MAX_HISTORY_LINES must be an integer: %s; using default: %d\n", maxHistoryLines, maxLines)
	}
	// ABS_HISTORY_FILE
	historyFile := util.GetEnvVar(env, "ABS_HISTORY_FILE", ABS_HISTORY_FILE)
	if maxLines > 0 {
		// expand the ABS_HISTORY_FILE to the user's HomeDir
		filePath, err := util.ExpandPath(historyFile)
		if err != nil {
			fmt.Printf("Unable to expand ABS history file path: %s\nError: %s\n", historyFile, err.Error())
			os.Exit(99)
		}
		historyFile = filePath
	}
	return historyFile, maxLines
}

// getHistory - read the history file and split it into the local history[...] slice
func getHistory(historyFile string, maxLines int) []string {
	var history []string
	if maxLines == 0 {
		// do not open a history file for zero max lines
		return history
	}
	// verify the expanded historyFile exists, if not create it now
	fd, ok := os.OpenFile(historyFile, os.O_RDONLY|os.O_CREATE, 0666)
	if ok != nil {
		fmt.Printf("Cannot create or read ABS history file: %s\nError: %s\n", historyFile, ok.Error())
		os.Exit(99)
	}
	fd.Close()
	// read the file and split the lines into history[...]
	bytes, err := os.ReadFile(historyFile)
	if err != nil {
		return history
	}
	// fill the local history from the file
	if len(bytes) > 0 {
		history = strings.Split(string(bytes), "\n")
	}
	return history
}

// addToHistory - append unique next line to local history[...]
func addToHistory(history []string, maxLines int, line string) []string {
	if maxLines == 0 {
		// do not save history for zero max lines
		return []string{}
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
	return history
}

// saveHistory - save the local history containing maxLines to historyFile
func saveHistory(historyFile string, maxLines int, history []string) error {
	if maxLines == 0 {
		// do not save a history file for zero max lines
		return nil
	}
	if len(history) > maxLines {
		// remove the excess lines from the front of the history slice
		history = history[len(history)-maxLines:]
	}
	// write the augmented local history back out to the file
	historyStr := strings.Join(history, "\n")
	return os.WriteFile(historyFile, []byte(historyStr), 0664)
}
