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

	case string:
		err := testStringObject(expected, actual)

		if err != nil {
			t.Errorf("testStringObject failed : %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%v)", actual, actual)
		}

	case []int:
		array, ok := actual.(*object.Array)

		if !ok {
			t.Errorf("object is not array : %T (%v)", actual, actual)
			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong number of elements. want = %d, got = %d ", len(expected), len(array.Elements))
			return
		}

		for index, expectedElement := range expected {
			err := testIntegerObject(int64(expectedElement), array.Elements[index])

			if err != nil {
				t.Errorf("testIntegerObject failed : %s", err)
			}
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
		{"if ((if (false) { 10 } )) { 10 } else { 20 }", 20},
	}

	runVmTests(t, testTable)
}

func TestGlobalLetStatements(t *testing.T) {
	testTable := []VmTestCase{
		{"let one = 1; one;", 1},
		{"let one = 1; let two = 2; one + two;", 3},
		{"let one = 1; let two = one + one; one + two;", 3},
	}

	runVmTests(t, testTable)
}

func TestStringExpressions(t *testing.T) {
	testTable := []VmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVmTests(t, testTable)
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)

	if !ok {
		return fmt.Errorf("object is not String. got = %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got = %q , want = %q", result.Value, expected)
	}

	return nil
}

func TestArrayLiterals(t *testing.T) {
	testTable := []VmTestCase{
		{input: "[]", expected: []int{}},
		{input: "[1, 2, 3]", expected: []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, testTable)
}

func TestHashLiterals(t *testing.T) {
	testable := []VmTestCase{
		{
			input:    "{}",
			expected: map[object.HashKey]int64{},
		},
		{
			input: "{1: 2, 2: 3}",
			expected: map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			input: "{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			expected: map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	runVmTests(t, testable)
}
