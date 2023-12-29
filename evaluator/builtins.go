package evaluator

import (
	"fmt"
	"github.com/Neal-C/compiler-in-go/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to len not supported, got %s", args[0].Type())
			}

		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be an ARRAY, got %s", args[0].Type())
			}

			myArray := args[0].(*object.Array)

			if len(myArray.Elements) > 0 {
				return myArray.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to last must be an ARRAY, got %s", args[0].Type())
			}

			myArray := args[0].(*object.Array)
			length := len(myArray.Elements)
			if len(myArray.Elements) > 0 {
				return myArray.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to rest must be an ARRAY, got %s", args[0].Type())
			}

			myArray := args[0].(*object.Array)
			length := len(myArray.Elements)
			if len(myArray.Elements) > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, myArray.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to push must be an ARRAY, got %s", args[0].Type())
			}

			myArray := args[0].(*object.Array)
			length := len(myArray.Elements)

			newElements := make([]object.Object, length+1)
			copy(newElements, myArray.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}

		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
