package compiler

import (
	"fmt"
	"github.com/Neal-C/compiler-in-go/ast"
	"github.com/Neal-C/compiler-in-go/code"
	"github.com/Neal-C/compiler-in-go/lexer"
	"github.com/Neal-C/compiler-in-go/object"
	"github.com/Neal-C/compiler-in-go/parser"
	"testing"
)

type CompilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tableTests := []CompilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1; 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input:             " 1 - 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             " 1 * 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             " 2 / 1",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tableTests)
}

func runCompilerTests(t *testing.T, tests []CompilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		myCompiler := New()

		err := myCompiler.Compile(program)

		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := myCompiler.ByteCode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)

		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(tt.expectedConstants, bytecode.Constants)

		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	myLexer := lexer.New(input)
	myParser := parser.New(myLexer)
	program := myParser.ParseProgram()

	return program
}

func testInstructions(expectedInstructions []code.Instructions, actualInstructions code.Instructions) error {
	concatted := concatInstructions(expectedInstructions)

	if len(actualInstructions) != len(concatted) {
		return fmt.Errorf("wrong instructions length. \n want = %q \n got = %q", concatted, actualInstructions)
	}

	for index, instruction := range concatted {
		if actualInstructions[index] != instruction {
			return fmt.Errorf("wrong instruction at position %d.\nwant = %q\n got= %q", index, concatted, actualInstructions[index])
		}
	}

	return nil
}

func concatInstructions(instructions []code.Instructions) code.Instructions {
	var out code.Instructions

	for _, instruction := range instructions {
		out = append(out, instruction...)
	}

	return out
}

func testConstants(expected []any, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got = %d , want = %d", len(actual), len(expected))
	}

	for index, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[index])

			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed : %s ", index, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)

	if !ok {
		return fmt.Errorf("object is not *object.Integer. got = %T (%v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value, got = %d , want = %d", result.Value, expected)
	}

	return nil

}
