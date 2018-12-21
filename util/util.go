package util

import (
	"abs/object"
	"fmt"
)

func NewError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

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
