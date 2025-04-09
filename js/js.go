package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/abs-lang/abs/evaluator"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"
)

// Version of the ABS interpreter
var Version = "dev"

// This function takes ABS code
// and evaluates it, using a buffer
// to store it's output.
// Once the code is evaluated, both
// the output and the return value of
// the script are returned to js in the
// form of {out, result}.
func runCode(this js.Value, i []js.Value) interface{} {
	m := make(map[string]interface{})
	var buf bytes.Buffer
	// the first argument to our function
	code := i[0].String()
	env := object.NewEnvironment(&buf, "", Version, true)
	lex := lexer.New(code)
	p := parser.New(lex)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors(), buf)
		m["out"] = buf.String()
		return js.ValueOf(m)
	}

	result := evaluator.BeginEval(program, env, lex)
	m["out"] = buf.String()
	m["result"] = result.Inspect()

	return js.ValueOf(m)
}

func printParserErrors(errors []string, buf bytes.Buffer) {
	fmt.Fprintf(&buf, "%s", " parser errors:\n")
	for _, msg := range errors {
		fmt.Fprintf(&buf, "%s", "\t"+msg+"\n")
	}
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("abs_run_code", js.FuncOf(runCode))
	<-c
}
