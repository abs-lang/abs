package terminal

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/runner"
)

func TestAssignStatements(t *testing.T) {
	for _, stmt := range exampleStatements {
		discard := bufio.NewReadWriter(bufio.NewReader(strings.NewReader("")), bufio.NewWriter(io.Discard))
		stdio := &object.Stdio{Stdin: discard, Stdout: discard, Stderr: discard}
		_, ok, errs := runner.Run(stmt, object.NewEnvironment(stdio, ".", "test", false))

		if !ok {
			t.Fatalf("%s (code evaluated: %s)", errs[0], stmt)
		}
	}
}
