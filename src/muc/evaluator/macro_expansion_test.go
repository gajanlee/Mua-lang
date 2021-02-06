package evaluator

import (
	"muc/ast"
	"muc/lexer"
	"muc/object"
	"muc/parser"
	"testing"
)

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`
			let infixExpression = macro() { quote(1 + 2) }
			infixExpression()
			`,
			`(1 + 2)`,
		},
		{
			`
			let reverse = macro(a, b) { quote(unquote(b) - unquote(a)) }
			reverse(2 + 2, 10 - 5)
			`,	// let every node be QUOTE_OBJ
			`(10 - 5) - (2 + 2)`,
		},
		{
			`
			let unless = macro(condition, consequence, alternative) {
				quote(if (!unquote(condition)) {
					unquote(consequence)
				} else {
					unquote(alternative)
				})
			}
			unless(10 > 5, print("not greater"), print("greater"))
			`,
			`if (!(10 > 5)) { print("not greater") } else { print("greater")}`,
		},
	}

	for _, tt := range tests {
		expected := testParseProgram(tt.expected)
		program := testParseProgram(tt.input)

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)
		
		if expanded.String() != expected.String() {
			t.Errorf("not equal. expected=%q. got=%q",
				expected.String(), expanded.String())
		}
	}
}