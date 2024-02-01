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

func testConstants(expectedConstants []any, actualObjects []object.Object) error {
	if len(expectedConstants) != len(actualObjects) {
		return fmt.Errorf("wrong number of constants. got = %d , want = %d", len(actualObjects), len(expectedConstants))
	}

	for index, constant := range expectedConstants {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actualObjects[index])

			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed : %s ", index, err)
			}

		case string:
			err := testStringObject(constant, actualObjects[index])

			if err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", index, err)
			}

		case []code.Instructions:
			fn, ok := actualObjects[index].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d - not a function: %T",
					index, actualObjects[index])
			}
			err := testInstructions(constant, fn.Instructions)
			if err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s",
					index, err)
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
			input: `if (true) { 10 } else { 20 }; 3333;`,
			// expected constants, in the constants pool/ data section / static
			expectedConstants: []any{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				// jump to 13th index
				code.Make(code.OpJump, 13),
				// 0010
				code.Make(code.OpConstant, 1),
				// 0013
				code.Make(code.OpPop),
				// 0014
				code.Make(code.OpConstant, 2),
				// 0017
				code.Make(code.OpPop),
			},
		},
		{
			input:             "if (true) { 10 }; 3333;",
			expectedConstants: []any{10, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 007
				code.Make(code.OpJump, 11),
				// 0010
				code.Make(code.OpNull),
				// 0011
				code.Make(code.OpPop),
				// 0012
				code.Make(code.OpConstant, 1),
				// 0015
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}

func TestGlobalLetStatemets(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input: `
				let one = 1;
				let two = 2;
				`,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpConstant, 0),
				// 0003
				code.Make(code.OpSetGlobal, 0),
				// 0006
				code.Make(code.OpConstant, 1),
				// 0009
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
				let one = 1;
				one;
				`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpConstant, 0),
				// 0003
				code.Make(code.OpSetGlobal, 0),
				// 0006
				code.Make(code.OpGetGlobal, 0),
				// 0009
				code.Make(code.OpPop),
			},
		},
		{
			input: `
				let one = 1;
				let two = one;
				two;
				`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpConstant, 0),
				// 0003
				code.Make(code.OpSetGlobal, 0),
				// 0006
				code.Make(code.OpGetGlobal, 0),
				// 0009
				code.Make(code.OpSetGlobal, 1),
				// 0012
				code.Make(code.OpGetGlobal, 1),
				// 0015
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}

func TestStringExpressions(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input:             `"monkey"`,
			expectedConstants: []any{"monkey"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `"mon" +  "key"`,
			expectedConstants: []any{"mon", "key"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)

	if !ok {
		return fmt.Errorf("object is not a String. got = %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got = %q , want = %q", result.Value, expected)
	}

	return nil
}

func TestArrayLiterals(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input:             "[]",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "[1 , 2 , 3]",
			expectedConstants: []any{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		}, {
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}

func TestHashLiterals(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input:             "{}",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{ 1: 2, 3: 4, 5: 6}",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpHash, 6),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2 + 3, 4: 5 * 6}",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}

func TestIndexExpressions(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []any{1, 2, 3, 1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2}[2 - 1]",
			expectedConstants: []any{1, 2, 2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}

func TestFunctions(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input: `fn() { return 5 + 10 }`,
			expectedConstants: []any{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() { 5 + 10 }`,
			expectedConstants: []any{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() { 1; 2 }`,
			expectedConstants: []any{
				1,
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, testTable)
}

func TestCompilerScopes(t *testing.T) {
	myCompiler := New()

	if myCompiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got = %d, want = %d", myCompiler.scopeIndex, 0)
	}

	myCompiler.emit(code.OpMul)

	myCompiler.enterScope()

	if myCompiler.scopeIndex != 1 {
		t.Errorf("scopeIndex wrong. got = %d, want = %d", myCompiler.scopeIndex, 1)
	}

	myCompiler.emit(code.OpSub)

	if len(myCompiler.scopes[myCompiler.scopeIndex].instructions) != 1 {
		t.Errorf("instructions length wrong. got=%d", len(myCompiler.scopes[myCompiler.scopeIndex].instructions))
	}

	lastInstruction := myCompiler.scopes[myCompiler.scopeIndex].lastInstruction

	if lastInstruction.OpCode != code.OpSub {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d", lastInstruction.OpCode, code.OpSub)
	}

	myCompiler.leaveScope()

	if myCompiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got = %d, want = %d", myCompiler.scopeIndex, 0)
	}

	myCompiler.emit(code.OpAdd)

	if len(myCompiler.scopes[myCompiler.scopeIndex].instructions) != 2 {
		t.Errorf("instructions length wrong. got=%d", len(myCompiler.scopes[myCompiler.scopeIndex].instructions))
	}

	lastInstruction = myCompiler.scopes[myCompiler.scopeIndex].lastInstruction

	if lastInstruction.OpCode != code.OpAdd {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d", lastInstruction.OpCode, code.OpAdd)
	}

	previousInstruction := myCompiler.scopes[myCompiler.scopeIndex].previousInstruction

	if previousInstruction.OpCode != code.OpMul {
		t.Errorf("previousInstruction.Opcode wrong. got=%d, want=%d", previousInstruction.OpCode, code.OpMul)

	}

}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input: `fn() { }`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, testTable)
}

func TestFunctionsCalls(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input: `fn(){ 24 }();`,
			expectedConstants: []any{
				24,
				[]code.Instructions{
					code.Make(code.OpConstant, 0), // the literal 24
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 1), // the compiled Fn
				code.Make(code.OpCall),
				code.Make(code.OpPop),
			},
		},
		{
			input: `let noArg = fn() { 24 };
					noArg();
					`,
			expectedConstants: []any{
				24,
				[]code.Instructions{
					code.Make(code.OpConstant, 0), // the literal 24
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 1), // the compiled Fn
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}

func TestLetStatementsScopes(t *testing.T) {
	testTable := []CompilerTestCase{
		{
			input: `
				let num = 55;
				fn() { num };
					`,
			expectedConstants: []any{
				55,
				[]code.Instructions{
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
					fn() {
					let num = 55;
					num
					}
					`,
			expectedConstants: []any{
				55,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
					fn() {
					let a = 55;
					let b = 77;
					a + b
					}
					`,
			expectedConstants: []any{
				55,
				77,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetLocal, 1),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testTable)
}
