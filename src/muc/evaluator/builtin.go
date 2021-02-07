package evaluator

import (
	"muc/object"
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
	"len": object.GetBuiltinByName("len"),
	"puts": object.GetBuiltinByName("puts"),
	"first": object.GetBuiltinByName("first"),
	// "first": &object.Builtin{Fn: _first},
	// "print": &object.Builtin{Fn: _print},
}