package repl

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/abs-lang/abs/util"
)

// support for ABS init file
const ABS_INIT_FILE = "~/.absrc"

func getAbsInitFile(interactive bool) {
	// get ABS_INIT_FILE from OS environment or default
	initFile := os.Getenv("ABS_INIT_FILE")
	if len(initFile) == 0 {
		initFile = ABS_INIT_FILE
	}
	// expand the ABS_INIT_FILE to the user's HomeDir
	filePath, err := util.ExpandPath(initFile)
	if err != nil {
		fmt.Printf("Unable to expand ABS init file path: %s\nError: %s\n", initFile, err.Error())
		os.Exit(99)
	}
	initFile = filePath
	// read and eval the abs init file
	code, err := ioutil.ReadFile(initFile)
	if err != nil {
		// abs init file is optional -- nothing to do here
		return
	}
	Run(string(code), interactive)
}

// support for user config of ABS REPL prompt string
const ABS_PROMPT_PREFIX = "⧐  "

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
