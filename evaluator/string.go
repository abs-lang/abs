package evaluator

import (
	"abs/ast"
	"abs/lexer"
	"abs/object"
	"abs/parser"
	"abs/util"
	"strings"
)

func getStringFns() map[string]*object.Builtin {
	return map[string]*object.Builtin{
		"split": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				// TODO we're passing the wrong amount of args to these functions
				if len(args) != 2 {
					return util.NewError("wrong number of arguments. got=%d, want=1",
						len(args))
				}
				if args[0].Type() != object.STRING_OBJ {
					return util.NewError("argument to `split` must be STRING, got %s", args[0].Type())
				}

				s := args[0].(*object.String)
				sep := args[1].(*object.String)
				parts := strings.Split(s.Value, sep.Value)
				length := len(parts)
				elements := make([]object.Object, length, length)

				for k, v := range parts {
					elements[k] = &object.String{Value: v}
				}

				return &object.Array{Elements: elements}
			},
		},
		"ok": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				s := args[0].(*object.String)

				return &object.Boolean{Value: s.Ok}
			},
		},
		// Converts a valid JSON document to an ABS hash.
		//
		// One interesting thing here is that we're creating
		// a new environment from scratch, whereas it might
		// be interesting to use the existing one. That would
		// allow to do things like:
		//
		// x = 10
		// '{"key": x}'.json()["key"] // 10
		//
		// Also, we're instantiating a new lexer & parser from
		// scratch, so this is a tad slow.
		//
		// This method is incomplete as it currently does not
		// support most JSON types, but rather just objects,
		// ie. "[1, 2, 3]".json() won't work.
		"json": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != object.STRING_OBJ {
					return util.NewError("argument to `split` must be STRING, got %s", args[0].Type())
				}

				s := args[0].(*object.String)
				env := object.NewEnvironment()
				l := lexer.New(s.Value)
				p := parser.New(l)
				hl, ok := p.ParseHashLiteral().(*ast.HashLiteral)

				if ok {
					return evalHashLiteral(hl, env)
				}

				return util.NewError("argument to `json` must be a valid JSON object, got '%s'", s.Value)
			},
		},
	}
}
