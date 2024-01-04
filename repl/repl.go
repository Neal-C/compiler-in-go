package repl

import (
	"bufio"
	"fmt"
	"github.com/Neal-C/compiler-in-go/compiler"
	"github.com/Neal-C/compiler-in-go/lexer"
	"github.com/Neal-C/compiler-in-go/parser"
	"github.com/Neal-C/compiler-in-go/vm"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		monkeyLexer := lexer.New(line)
		monkeyParser := parser.New(monkeyLexer)
		program := monkeyParser.ParseProgram()

		if len(monkeyParser.Errors()) != 0 {
			printParseErrors(out, monkeyParser.Errors())
			continue
		}

		myCompiler := compiler.New()
		err := myCompiler.Compile(program)

		if err != nil {
			fmt.Fprintf(out, "Whoops! compilation failed:\n %s\n", err)
			continue
		}

		machine := vm.New(myCompiler.ByteCode())

		err = machine.Run()

		if err != nil {
			fmt.Fprintf(out, "Whoops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		stackTop := machine.StackTop()

		if stackTop != nil {
			io.WriteString(out, stackTop.Inspect())
			io.WriteString(out, "\n")
		}
	}

}

func printParseErrors(writer io.Writer, errors []string) {
	for _, error := range errors {
		_, _ = io.WriteString(writer, "\t"+error+"\n")
		// no error handling apparently
	}
}
