package repl

import (
	"bufio"
	"fmt"
	"io"
	"mua/ast"
	"mua/evaluator"
	"mua/lexer"
	"mua/object"
	"mua/parser"
	"mua/token"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		// TODO: Handle the Direction Keyboards

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect() + "\n")
		}
	}
}

func printToken(out io.Writer, l lexer.Lexer) {
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		io.WriteString(out, fmt.Sprintf("%+v\n", tok))
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg + "\n")
	}
}

func printStatements(out io.Writer, program *ast.Program) {
	io.WriteString(out, program.String() + "\n")
}