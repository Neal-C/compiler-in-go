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
		{
			input:             "-1",
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
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
			return fmt.Errorf("wrong instruction at position %d.\nwant = %q\n got= %q", index, concatted, actualInstructions)
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

func TestBooleanExpressions(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input:             "true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		}, {
			input:             "false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{input: "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}

func TestConditionals(t *testing.T) {

	// What we care about are the jump instructions the compiler emits and
	// that they have correct operands.

	testTable := []CompilerTestCase{
		{
			input: `if (true) { 10 }; 3333;`,
			// expected constants, in the constants pool/ data section / static
			expectedConstants: []any{10, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				// we're sitting on 0000
				code.Make(code.OpTrue),
				// 0001
				// size of his operands is 2, so that takes 0002 0003, so we go to 0004
				code.Make(code.OpJumpNotTruthy, 7),
				// 0004
				// we're sitting on 0004
				// size of his operands is 2, so that takes 0005 and 0006
				code.Make(code.OpConstant, 0),
				// 0007
				// no operands
				code.Make(code.OpPop),
				// 0008
				code.Make(code.OpConstant, 1),
				// 0011
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}
