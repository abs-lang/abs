package repl

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

var reverseSearchStr string
var lastSearchPosition int

func clearReverseSearch() {
	reverseSearchStr = ""
	initReverseSearch()
}

func initReverseSearch() {
	lastSearchPosition = len(history)
}

// Trigger a reverse search in the shell
// on Ctrl + R: we will go through the history
// and find a match for our search text
func reverseSearch() prompt.KeyBind {
	return prompt.KeyBind{
		Key: prompt.ControlR,
		Fn: func(buf *prompt.Buffer) {
			// Get the text on the repl
			text := buf.Text()

			// If there's text and the search string
			// is not set, let's define our search string (ie. "curl")
			if text != "" && reverseSearchStr == "" {
				reverseSearchStr = text
			}

			// If our last search position is 0,
			// it means we went through the entire
			// history. At this point, it's useless
			// to search more.
			if lastSearchPosition == 0 {
				return
			}

			// Let's go from the last entry in the history
			// to the first: if we find an entry that matches
			// our text, let's replace the input in the repl
			// with the history, and break out. Let's also set
			// the lastSearchPosition so that, if the user
			// presses Ctrl + R once again, we start our search
			// from the last position onward.
			for i := lastSearchPosition - 1; i >= 0; i-- {
				lastSearchPosition = i

				if strings.Contains(history[i], reverseSearchStr) && text != history[i] {
					buf.CursorRight(len(text))
					buf.DeleteBeforeCursor(len(text))
					buf.InsertText(history[i], false, false)
					break
				}
			}
			return
		},
	}
}
