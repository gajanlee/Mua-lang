package repl

import (
	"bufio"
	"fmt"
	"io"
	"muc/compiler"
	"muc/lexer"
	"muc/object"
	"muc/parser"
	"muc/vm"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Compilation failed:\n %s\n", err)
			continue
		}

		bytecode := comp.Bytecode()
		constants = bytecode.Constants

		machine := vm.NewWithGlobalsState(bytecode, globals)

		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Executing bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg + "\n")
	}
}