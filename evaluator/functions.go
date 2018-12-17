package evaluator

import (
	"abs/object"
	"abs/util"
	"fmt"
	"os"
)

func getFns() map[string]*object.Builtin {
	return map[string]*object.Builtin{
		"len": &object.Builtin{Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return util.NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return util.NewError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
		},
		"echo": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return NULL
			},
		},
		"env": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				switch arg := args[0].(type) {
				case *object.String:
					return &object.String{Value: os.Getenv(arg.Value)}
				default:
					return util.NewError("argument to `env` not supported, got %s",
						args[0].Type())
				}
			},
		},
	}

}
