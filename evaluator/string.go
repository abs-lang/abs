package evaluator

import (
	"abs/object"
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
					return util.NewError("argument to `split` must be STRING, got %s",
						args[0].Type())
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
	}
}
