package repl

import (
	"bufio"
	"fmt"
	"github.com/Neal-C/interpreter-in-go/evaluator"
	"github.com/Neal-C/interpreter-in-go/lexer"
	"github.com/Neal-C/interpreter-in-go/object"
	"github.com/Neal-C/interpreter-in-go/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

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

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
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
