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

func TestApplySuggestions(t *testing.T) {
	tests := [][]string{
		{"int", "int", "intersect", "intersect"},
		{"f(input)", "input", "TTT", "f(TTT)"},
	}
	for _, tt := range tests {
		res := applySuggestion(tt[0], tt[1], tt[2])

		if res != tt[3] {
			t.Fatalf("got %s exp %s", res, tt[3])
		}
	}
}
