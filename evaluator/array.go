package evaluator

import (
	"abs/object"
	"abs/util"
)

func getArrayFns() map[string]*object.Builtin {
	return map[string]*object.Builtin{
		"first": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1",
						len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return util.NewError("argument to `first` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}

				return NULL
			},
		},
		"map": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return util.NewError("wrong number of arguments. got=%d, want=2", len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return util.NewError("argument to `map` must be ARRAY, got %s",
						args[0].Type())
				}
				if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
					return util.NewError("argument to `map` must be FUNCTION, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				newElements := make([]object.Object, length, length)
				copy(newElements, arr.Elements)

				for k, v := range arr.Elements {
					newElements[k] = applyFunction(args[1], []object.Object{v})
				}

				return &object.Array{Elements: newElements}
			},
		},
	}
}
