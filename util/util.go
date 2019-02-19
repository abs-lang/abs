package util

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/abs-lang/abs/object"
)

// Checks whether the element e is in the
// list of strings s
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func IsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)

	return err == nil
}

// ExpandPath (path) resolves leading "~/" to user's HomeDir
// returns expanded path, error
func ExpandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

// GetEnvVar (varName, defaultVal)
// Return the varName value from the ABS env, or OS env, or default value in that order
func GetEnvVar(env *object.Environment, varName, defaultVal string) string {
	var ok bool
	var value string
	valueObj, ok := env.Get(varName)
	if ok {
		value = valueObj.Inspect()
	} else {
		value = os.Getenv(varName)
		if len(value) == 0 {
			value = defaultVal
		}
	}
	return value
}
