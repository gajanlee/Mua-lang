package evaluator

import (
	"muc/ast"
	"muc/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, statement := range program.Statements {
		// Find all Macro Definitions and record
		if isMacroDefinition(statement) {
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		// remove the Macro Definitions
		definitionIndex := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

func isMacroDefinition(node ast.Statement) bool {
	letStatement, ok := node.(*ast.LetStatement)
	if !ok { return false }
	_, ok = letStatement.Value.(*ast.MacroLiteral)
	if !ok { return false }

	return true
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	letStatement, _ := stmt.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)

	macro := &object.Macro {
		Parameters: macroLiteral.Parameters,
		Env:        env,
		Body:       macroLiteral.Body,
	}

	env.Set(letStatement.Name.Value, macro)
}

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok { return node }
		macro, ok := isMacroCall(callExpression, env)
		if !ok { return node }

		// convert to QUOTE, prevent evaluating
		args := quoteArgs(callExpression)
		evalEnv := extendMacroEnv(macro, args)
		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros")
		}

		return quote.Node
	})
}

func isMacroCall(expr *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	identifier, ok := expr.Function.(*ast.Identifier)
	if !ok { return nil, false }
	obj, ok := env.Get(identifier.Value)
	if !ok { return nil, false }

	macro, ok := obj.(*object.Macro)
	if !ok { return nil, false }
	return macro, true
}

func quoteArgs(expr *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}
	for _, arg := range expr.Arguments {
		args = append(args, &object.Quote{Node: arg})
	}
	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := object.NewEnclosedEnvironment(macro.Env)

	for paramIdx, param := range macro.Parameters {
		extended.Set(param.Value, args[paramIdx])
	}
	return extended
}