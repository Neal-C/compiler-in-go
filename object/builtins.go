package object

import (
	"fmt"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		Name: "len",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to len not supported, got %s", args[0].Type())
				}

			},
		},
	},
	{
		Name: "puts",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return nil
			},
		},
	},
	{
		Name: "first",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to first must be an ARRAY, got %s", args[0].Type())
				}

				myArray := args[0].(*Array)

				if len(myArray.Elements) > 0 {
					return myArray.Elements[0]
				}

				return nil
			},
		},
	},
	{
		Name: "last",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to last must be an ARRAY, got %s", args[0].Type())
				}

				myArray := args[0].(*Array)
				length := len(myArray.Elements)
				if len(myArray.Elements) > 0 {
					return myArray.Elements[length-1]
				}

				return nil
			},
		},
	},
	{
		Name: "rest",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to rest must be an ARRAY, got %s", args[0].Type())
				}

				myArray := args[0].(*Array)
				length := len(myArray.Elements)
				if len(myArray.Elements) > 0 {
					newElements := make([]Object, length-1)
					copy(newElements, myArray.Elements[1:length])
					return &Array{Elements: newElements}
				}

				return nil
			},
		},
	},
	{
		Name: "push",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to push must be an ARRAY, got %s", args[0].Type())
				}

				myArray := args[0].(*Array)
				length := len(myArray.Elements)

				newElements := make([]Object, length+1)
				copy(newElements, myArray.Elements)
				newElements[length] = args[1]
				return &Array{Elements: newElements}

			},
		},
	},
}

func newError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {

	for _, definition := range Builtins {
		if definition.Name == name {
			return definition.Builtin
		}
	}

	return nil

}
