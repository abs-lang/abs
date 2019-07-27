package util

import (
	"os"
	"os/user"
	"path/filepath"
	"regexp"
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

// InterpolateStringVars (str, env)
// return input string with $vars interpolated from environment
func InterpolateStringVars(str string, env *object.Environment) string {
	// Match all strings preceded by
	// a $ or a \$
	re := regexp.MustCompile("(\\\\)?\\$([a-zA-Z_0-9]{1,})")
	str = re.ReplaceAllStringFunc(str, func(m string) string {
		// If the string starts with a backslash,
		// that's an escape, so we should replace
		// it with the remaining portion of the match.
		// \$VAR becomes $VAR
		if string(m[0]) == "\\" {
			return m[1:]
		}
		// If the string starts with $, then
		// it's an interpolation. Let's
		// replace $VAR with the variable
		// named VAR in the ABS' environment.
		// If the variable is not found, we
		// just dump an empty string
		v, ok := env.Get(m[1:])

		if !ok {
			return ""
		}
		return v.Inspect()
	})
	return str
}
