package evaluator

import (
	"fmt"
	"mua/object"
)

/**
* Registe BUILT-IN FUNCTIONs
* 	- len
* 	- first
* 	- last
*	- rest
*	- push
*	- pop
*	- map
*	- reduce
*
*   - puts
*/
var builtins = map[string]*object.Builtin {
	"len": &object.Builtin{Fn: _len},
	"first": &object.Builtin{Fn: _first},

	"puts": &object.Builtin{Fn: _puts},
}

func _len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	}
	return newError("argument to `len` not supported, got %s", args[0].Type())
}

func _first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.Array:
		return arg.Elements[0]
	default:
		return newError("unsupported type: %T", args[0].Type())
	}
}

func _puts(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}