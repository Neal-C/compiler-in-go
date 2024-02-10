package main

import (
	"flag"
	"fmt"
	"github.com/Neal-C/compiler-in-go/compiler"
	"github.com/Neal-C/compiler-in-go/evaluator"
	"github.com/Neal-C/compiler-in-go/lexer"
	"github.com/Neal-C/compiler-in-go/object"
	"github.com/Neal-C/compiler-in-go/parser"
	"github.com/Neal-C/compiler-in-go/vm"
	"time"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")
var input = `
let fibonacci = fn(x) {
if (x == 0) {
0
} else {
if (x == 1) {
return 1;
} else {
fibonacci(x - 1) + fibonacci(x - 2);
}
}
};
fibonacci(35);
`

func main() {
	flag.Parse()
	var duration time.Duration
	var result object.Object

	myLexer := lexer.New(input)

	myParser := parser.New(myLexer)

	program := myParser.ParseProgram()

	if *engine == "vm" {
		myCompiler := compiler.New()
		err := myCompiler.Compile(program)
		if err != nil {
			fmt.Printf("compiler error: %s", err)
			return
		}
		myVM := vm.New(myCompiler.ByteCode())
		start := time.Now()
		err = myVM.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)
			return
		}
		duration = time.Since(start)
		result = myVM.LastPoppedStackElement()
	} else {
		env := object.NewEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	}
	fmt.Printf(
		"engine=%s, result=%s, duration=%s\n",
		*engine,
		result.Inspect(),
		duration)
}
