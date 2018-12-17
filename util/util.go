package util

import (
	"abs/object"
	"fmt"
)

func NewError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
