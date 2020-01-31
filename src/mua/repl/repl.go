package repl

import (
	"bufio"
	"fmt"
	"io"
	"mua/lexer"
	"mua/token"
	"mua/parser"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
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