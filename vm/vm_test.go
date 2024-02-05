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

type vmTestCase struct {
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

func runVmTests(t *testing.T, tableTests []vmTestCase) {
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

	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)

		if !ok {
			t.Errorf("object is not a Hash. got = %T (%+v)", actual, actual)
			return
		}

		if len(hash.Pairs) != len(expected) {
			t.Errorf("hash has wrong number of pairs. want = %d , got = %d", len(expected), len(hash.Pairs))
			return
		}

		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]

			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}

			err := testIntegerObject(expectedValue, pair.Value)

			if err != nil {
				t.Errorf("testIntegerObject failed : %s", err)
			}
		}
	}

}

func TestIntegerArithmetic(t *testing.T) {
	testTable := []vmTestCase{
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
	testTable := []vmTestCase{
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
	testTable := []vmTestCase{
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
	testTable := []vmTestCase{
		{"let one = 1; one;", 1},
		{"let one = 1; let two = 2; one + two;", 3},
		{"let one = 1; let two = one + one; one + two;", 3},
	}

	runVmTests(t, testTable)
}

func TestStringExpressions(t *testing.T) {
	testTable := []vmTestCase{
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
	testTable := []vmTestCase{
		{input: "[]", expected: []int{}},
		{input: "[1, 2, 3]", expected: []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, testTable)
}

func TestHashLiterals(t *testing.T) {
	testable := []vmTestCase{
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

func TestIndexExpressions(t *testing.T) {
	testTable := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Null},
		{"[1, 2, 3][999]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}

	runVmTests(t, testTable)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	testTable := []vmTestCase{
		{
			input: `
					let fivePlusTen = fn() { return 5 + 10 };
					fivePlusTen()
					`,
			expected: 15,
		},
		{
			input: `
				let one = fn() { 1; };
				let two = fn() { 2; };
				one() + two()
				`,
			expected: 3,
		},
		{
			input: `
				let a = fn() { 1 };
				let b = fn() { a() + 1 };
				let c = fn() { b() + 1 };
				c();
				`,
			expected: 3,
		},
	}

	runVmTests(t, testTable)
}

func TestFunctionsWithReturnStatements(t *testing.T) {
	testTable := []vmTestCase{
		{
			input: `
				let earlyExit = fn() { return 99; 100; };
				earlyExit();
				`,
			expected: 99,
		},
		{
			input: `
				let earlyExit = fn() { return 99; return 100; };
				earlyExit();
				`,
			expected: 99,
		},
	}

	runVmTests(t, testTable)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	testTable := []vmTestCase{
		{
			input: `
				let noReturn = fn() { };
				noReturn();
				`,
			expected: Null,
		},
		{
			input: `
				let noReturn = fn() { };
				let noReturnTwo = fn() { noReturn(); };
				noReturn();
				noReturnTwo();
				`,
			expected: Null,
		},
	}

	runVmTests(t, testTable)
}

func TestFirstClassCitizenFunctions(t *testing.T) {
	testTable := []vmTestCase{
		{
			input: `
				let returnsOne = fn() { 1; };
				let returnsOneReturner = fn() { returnsOne; };
				returnsOneReturner()();
				`,
			expected: 1,
		},
		{
			input: `
				let returnsOneReturner = fn() {
				let returnsOne = fn() { 1; };
				returnsOne;
				};
				returnsOneReturner()();
				`,
			expected: 1,
		},
	}

	runVmTests(t, testTable)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	testTable := []vmTestCase{
		{
			input: `
				let one = fn() { let one = 1; one };
				one();
				`,
			expected: 1,
		},
		{
			input: `
				let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
				oneAndTwo();
				`,
			expected: 3,
		},
		{
			input: `
				let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
				let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
				oneAndTwo() + threeAndFour();
				`,
			expected: 10,
		},
		{
			input: `
				let firstFoobar = fn() { let foobar = 50; foobar; };
				let secondFoobar = fn() { let foobar = 100; foobar; };
				firstFoobar() + secondFoobar();
				`,
			expected: 150,
		},
		{
			input: `
				let globalSeed = 50;
				let minusOne = fn() {
				let num = 1;
				globalSeed - num;
				}
				let minusTwo = fn() {
				let num = 2;
				globalSeed - num;
				}
				minusOne() + minusTwo();
				`,
			expected: 97,
		},
	}

	runVmTests(t, testTable)

}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	testTable := []vmTestCase{
		{
			input: `
				let identity = fn(a) { a; }
				identity(4);
				`,
			expected: 4,
		},
		{
			input: `
				let sum = fn(a,b) { a + b; }
				sum(1,2);
					`,
			expected: 3,
		},
		{
			input: `
					let sum = fn(a,b){
							let c = a + b;
							c
								};
					sum(1,2);
					`,
			expected: 3,
		},
		{
			input: `
					let sum = fn(a,b){
							let c = a + b;
							c
								};
					sum(1,2) + sum(3,4);
					`,
			expected: 10,
		},
		{
			input: `
					let sum = fn(a,b){
							let c = a + b;
							c
								};
					let outerFn = fn() {
					sum(1,2) + sum(3,4);
					}
					outerFn();
					`,
			expected: 10,
		},
		{
			input: `
				let globalNum = 10;
				let sum = fn(a, b) {
				let c = a + b;
				c + globalNum;
				};
				let outer = fn() {
				sum(1, 2) + sum(3, 4) + globalNum;
				};
				outer() + globalNum;
				`,
			expected: 50,
		},
	}

	runVmTests(t, testTable)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	testTable := []vmTestCase{
		{
			input:    `fn() { 1; }(1);`,
			expected: `wrong number of arguments: want=0, got=1`,
		},
		{
			input:    `fn(a) { a; }();`,
			expected: `wrong number of arguments: want=1, got=0`,
		},
		{
			input:    `fn(a, b) { a + b; }(1);`,
			expected: `wrong number of arguments: want=2, got=1`,
		},
	}

	for _, tt := range testTable {
		program := parse(tt.input)

		myCompiler := compiler.New()

		err := myCompiler.Compile(program)

		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		myVM := New(myCompiler.ByteCode())

		err = myVM.Run()

		if err == nil {
			t.Fatalf("expected VM error but resulted in none.")
		}

		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error: want=%q, got=%q", tt.expected, err)
		}
	}

}

func TestBuiltinFunctions(t *testing.T) {
	testTable := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{
			`len(1)`,
			&object.Error{
				Message: "argument to `len` not supported, got INTEGER",
			},
		},
		{`len("one", "two")`,
			&object.Error{
				Message: "wrong number of arguments. got=2, want=1",
			},
		},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, Null},
		{`first([1, 2, 3])`, 1},
		{`first([])`, Null},
		{`first(1)`,
			&object.Error{
				Message: "argument to `first` must be ARRAY, got INTEGER",
			},
		},
		{`last([1, 2, 3])`, 3},
		{`last([])`, Null},
		{`last(1)`,
			&object.Error{
				Message: "argument to `last` must be ARRAY, got INTEGER",
			},
		},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, Null},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`,
			&object.Error{
				Message: "argument to `push` must be ARRAY, got INTEGER",
			},
		},
	}

	runVmTests(t, testTable)
}
