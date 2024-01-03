package vm

import (
	"fmt"
	"github.com/Neal-C/compiler-in-go/ast"
	"github.com/Neal-C/compiler-in-go/compiler"
	"testing"

	"github.com/Neal-C/compiler-in-go/lexer"
	"github.com/Neal-C/compiler-in-go/object"
	"github.com/Neal-C/compiler-in-go/parser"
)

type VmTestCase struct {
	input    string
	expected any
}

func parse(input string) *ast.Program {
	myLexer := lexer.New(input)
	myParser := parser.New(myLexer)
	program := myParser.ParseProgram()

	return program
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

func runVmTests(t *testing.T, tableTests []VmTestCase) {
	t.Helper()

	for _, tt := range tableTests {
		program := parse(tt.input)

		myCompiler := compiler.New()

		err := myCompiler.Compile(program)

		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		myVM := New(myCompiler.ByteCode())

		err := myVM.Run()

		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElement := myVM.StackTop()

		testExpectedObject(t, tt.expected, stackElement)

	}
}

func testExpectedObject(t *testing.T, expected any, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)

		if err != nil {
			t.Errorf("testIntegerObject failed %s", err)
		}
	}

}

func TestIntegerArithmetic(t *testing.T) {
	testTable := []VmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 2}, // TODO: FIXME in the actual book
	}

	runVmTests(t, testTable)
}
