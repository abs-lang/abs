package runner

import (
	"github.com/abs-lang/abs/evaluator"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"
)

// Run, well, runs an abs program.
// It returns output of the program,
// whether it encountered an error
// and parsing errors, so that we
// can print helpful error locations
// for you to fix he code.
func Run(code string, env *object.Environment) (out object.Object, ok bool, parseErrors []string) {
	lex := lexer.New(code)
	p := parser.New(lex)

	program := p.ParseProgram()
	parseErrors = p.Errors()

	if len(parseErrors) != 0 {
		return object.NULL, false, parseErrors
	}

	// invoke BeginEval() passing in the program, env, and lexer for error position
	// NB. Eval(node, env) is recursive so we can't call it directly
	evaluated := evaluator.BeginEval(program, env, lex)

	if evaluated == nil {
		return object.NULL, false, []string{}
	}

	return evaluated, evaluated.Type() != object.ERROR_OBJ, parseErrors
}
