package terminal

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/abs-lang/abs/evaluator"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"
)

func runCode(code string, env *object.Environment) []string {
	lex := lexer.New(code)
	p := parser.New(lex)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return p.Errors()
	}

	evaluated := evaluator.BeginEval(program, env, lex)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		return []string{evaluated.Inspect()}
	}

	return []string{}
}

func TestAssignStatements(t *testing.T) {
	for _, stmt := range exampleStatements {
		discard := bufio.NewReadWriter(bufio.NewReader(strings.NewReader("")), bufio.NewWriter(io.Discard))
		stdio := &object.Stdio{Stdin: discard, Stdout: discard, Stderr: discard}
		errs := runCode(stmt, object.NewEnvironment(stdio, ".", "test", false))

		if len(errs) > 0 {
			t.Fatalf("%s (code evaluated: %s)", errs[0], stmt)
		}
	}
}
