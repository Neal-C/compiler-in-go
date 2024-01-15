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

		err = myVM.Run()

		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElement := myVM.LastPoppedStackElement()

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
	case bool:
		err := testBooleanObject(bool(expected), actual)

		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%v)", actual, actual)
		}
	}

}

func TestIntegerArithmetic(t *testing.T) {
	testTable := []VmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, testTable)
}

func TestBooleanExpression(t *testing.T) {
	testTable := []VmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
	}

	runVmTests(t, testTable)
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)

	if !ok {
		return fmt.Errorf("actual is not a *object.Bookean. got = %T (%v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("actual has the wrong value. got = %t, want = %t", result.Value, expected)
	}

	return nil
}

func TestConditionals(t *testing.T) {
	testTable := []VmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if ( 1 > 2 ) { 10 }", Null},
		{"if (false) { 10 }", Null},
	}

	runVmTests(t, testTable)
}
