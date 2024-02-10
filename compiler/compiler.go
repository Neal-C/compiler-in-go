package compiler

import (
	"fmt"
	"github.com/Neal-C/compiler-in-go/ast"
	"github.com/Neal-C/compiler-in-go/code"
	"github.com/Neal-C/compiler-in-go/object"
	"sort"
)

type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable
	scopes      []CompilationScope
	scopeIndex  int
}

type EmittedInstruction struct {
	OpCode   code.Opcode
	Position int
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

func New() *Compiler {

	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()

	for index, value := range object.Builtins {
		symbolTable.DefineBuiltin(index, value.Name)
	}

	return &Compiler{
		constants:   []object.Object{},
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
		symbolTable: symbolTable,
	}
}

func NewWithState(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = symbolTable
	compiler.constants = constants
	return compiler
}

func (self *Compiler) currentInstructions() code.Instructions {
	return self.scopes[self.scopeIndex].instructions
}

func (self *Compiler) Compile(node ast.Node) error {

	switch node := node.(type) {
	case *ast.Program:

		for _, stmt := range node.Statements {
			err := self.Compile(stmt)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:

		err := self.Compile(node.Expression)

		if err != nil {
			return err
		}

		self.emit(code.OpPop)

	case *ast.InfixExpression:

		// Compiler magic, right there
		if node.Operator == "<" {
			err := self.Compile(node.Right)
			if err != nil {
				return err
			}

			err = self.Compile(node.Left)
			if err != nil {
				return err
			}

			self.emit(code.OpGreaterThan)
			return nil

		}

		err := self.Compile(node.Left)
		if err != nil {
			return err
		}

		err = self.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			self.emit(code.OpAdd)
		case "-":
			self.emit(code.OpSub)
		case "*":
			self.emit(code.OpMul)
		case "/":
			self.emit(code.OpDiv)
		case ">":
			self.emit(code.OpGreaterThan)
		case "==":
			self.emit(code.OpEqual)
		case "!=":
			self.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator : %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		self.emit(code.OpConstant, self.addConstants(integer))
	case *ast.Boolean:
		if node.Value {
			self.emit(code.OpTrue)
		} else {
			self.emit(code.OpFalse)
		}
	case *ast.PrefixExpression:
		err := self.Compile(node.Right)

		if err != nil {
			return err
		}

		switch node.Operator {
		case "-":
			self.emit(code.OpMinus)
		case "!":
			self.emit(code.OpBang)
		default:
			return fmt.Errorf("unkknown operator: %s", node.Operator)
		}

	case *ast.IfExpression:

		err := self.Compile(node.Condition)

		if err != nil {
			return err
		}

		// Emit with a bogus value that gets back-patched later
		jumpNotTruthyPosition := self.emit(code.OpJumpNotTruthy, 9999)

		err = self.Compile(node.Consequence)

		if err != nil {
			return err
		}

		if self.lastInstructionIs(code.OpPop) {
			self.removeLastPop()
		}

		jumpOverAlternativePosition := self.emit(code.OpJump, 9999)

		afterConsequencePos := len(self.currentInstructions())
		self.changeOperand(jumpNotTruthyPosition, afterConsequencePos)

		if node.Alternative == nil {
			self.emit(code.OpNull)
		} else {

			err = self.Compile(node.Alternative)

			if err != nil {
				return err
			}

			if self.lastInstructionIs(code.OpPop) {
				self.removeLastPop()
			}
		}

		afterAlternativePosition := len(self.currentInstructions())
		self.changeOperand(jumpOverAlternativePosition, afterAlternativePosition)

	case *ast.BlockStatement:

		for _, stmt := range node.Statements {
			err := self.Compile(stmt)

			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:

		symbol := self.symbolTable.Define(node.Name.Value)
		err := self.Compile(node.Value)

		if err != nil {
			return err
		}

		if symbol.Scope == GlobalScope {

			self.emit(code.OpSetGlobal, symbol.Index)

		} else {
			self.emit(code.OpSetLocal, symbol.Index)
		}
	case *ast.Identifier:

		symbol, ok := self.symbolTable.Resolve(node.Value)

		if !ok {
			// Compile time errors !!
			return fmt.Errorf("undefined variable : %s", node.Value)
		}

		self.loadSymbol(symbol)
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}

		constantPoolIndex := self.addConstants(str)

		self.emit(code.OpConstant, constantPoolIndex)

	case *ast.ArrayLiteral:

		for _, element := range node.Elements {
			err := self.Compile(element)

			if err != nil {
				return err
			}
		}

		self.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		var keys []ast.Expression

		for key := range node.Pairs {
			keys = append(keys, key)
		}

		sort.Slice(keys, func(i int, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {

			err := self.Compile(k)

			if err != nil {
				return err
			}

			err = self.Compile(node.Pairs[k])

			if err != nil {
				return err
			}
		}

		self.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.IndexExpression:

		err := self.Compile(node.Left)

		if err != nil {
			return err
		}

		err = self.Compile(node.Index)

		if err != nil {
			return err
		}

		self.emit(code.OpIndex)

	case *ast.FunctionLiteral:

		self.enterScope()

		if node.Name != "" {
			self.symbolTable.DefineFunctionName(node.Name)
		}

		for _, param := range node.Parameters {
			self.symbolTable.Define(param.Value)
		}

		err := self.Compile(node.Body)

		if err != nil {
			return err
		}

		if self.lastInstructionIs(code.OpPop) {
			self.replaceLastPopWithReturn()
		}

		if !self.lastInstructionIs(code.OpReturnValue) {
			self.emit(code.OpReturn)
		}

		freeSymbols := self.symbolTable.FreeSymbols
		numberOfLocals := self.symbolTable.numberOfDefinitions
		instructions := self.leaveScope()

		for _, symbol := range freeSymbols {
			self.loadSymbol(symbol)
		}

		compiledFn := &object.CompiledFunction{
			Instructions:       instructions,
			NumberOfLocals:     numberOfLocals,
			NumberOfParameters: len(node.Parameters),
		}

		fnIndex := self.addConstants(compiledFn)

		self.emit(code.OpClosure, fnIndex, len(freeSymbols))

	case *ast.ReturnStatement:

		err := self.Compile(node.ReturnValue)

		if err != nil {
			return err
		}

		self.emit(code.OpReturnValue)
	case *ast.CallExpression:

		err := self.Compile(node.Function)

		if err != nil {
			return err
		}

		for _, arg := range node.Arguments {
			err := self.Compile(arg)

			if err != nil {
				return err
			}
		}

		self.emit(code.OpCall, len(node.Arguments))

	}

	return nil
}

func (self *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: self.currentInstructions(),
		Constants:    self.constants,
	}
}

func (self *Compiler) addConstants(obj object.Object) int {
	self.constants = append(self.constants, obj)

	return len(self.constants) - 1
}

func (self *Compiler) addInstruction(instructions code.Instructions) int {
	posNewInstruction := len(self.currentInstructions())
	updatedInstructions := append(self.currentInstructions(), instructions...)

	self.scopes[self.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

func (self *Compiler) emit(op code.Opcode, operands ...int) int {

	instruction := code.Make(op, operands...)
	position := self.addInstruction(instruction)

	self.setLastInstruction(op, position)

	return position
}

func (self *Compiler) setLastInstruction(op code.Opcode, position int) {

	previous := self.scopes[self.scopeIndex].lastInstruction
	last := EmittedInstruction{OpCode: op, Position: position}

	self.scopes[self.scopeIndex].previousInstruction = previous
	self.scopes[self.scopeIndex].lastInstruction = last
}

func (self *Compiler) lastInstructionIs(op code.Opcode) bool {

	if len(self.currentInstructions()) == 0 {
		return false
	}
	return self.scopes[self.scopeIndex].lastInstruction.OpCode == op
}

func (self *Compiler) removeLastPop() {

	last := self.scopes[self.scopeIndex].lastInstruction
	previous := self.scopes[self.scopeIndex].previousInstruction

	oldInstructions := self.currentInstructions()
	newInstructions := oldInstructions[:last.Position]

	self.scopes[self.scopeIndex].instructions = newInstructions
	self.scopes[self.scopeIndex].lastInstruction = previous

}

func (self *Compiler) replaceInstruction(pos int, newInstruction []byte) {

	instructions := self.currentInstructions()

	for index := 0; index < len(newInstruction); index++ {
		instructions[pos+index] = newInstruction[index]
	}
}

func (self *Compiler) changeOperand(operationPosition int, operand int) {

	instructions := self.currentInstructions()

	// operationPosition is where we were
	op := code.Opcode(instructions[operationPosition])
	// operand is 2 bytes
	newInstruction := code.Make(op, operand)

	self.replaceInstruction(operationPosition, newInstruction)
}

func (self *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	self.scopes = append(self.scopes, scope)
	self.scopeIndex++

	self.symbolTable = NewEnclosedSymbolTable(self.symbolTable)
}

func (self *Compiler) leaveScope() code.Instructions {
	instructions := self.currentInstructions()

	self.scopes = self.scopes[:len(self.scopes)-1]
	self.scopeIndex--

	self.symbolTable = self.symbolTable.OuterTable

	return instructions
}

func (self *Compiler) replaceLastPopWithReturn() {

	lastInstructionPosition := self.scopes[self.scopeIndex].lastInstruction.Position

	self.replaceInstruction(lastInstructionPosition, code.Make(code.OpReturnValue))

	self.scopes[self.scopeIndex].lastInstruction.OpCode = code.OpReturnValue

}

func (self *Compiler) loadSymbol(symbl Symbol) {
	switch symbl.Scope {
	case GlobalScope:
		self.emit(code.OpGetGlobal, symbl.Index)
	case LocalScope:
		self.emit(code.OpGetLocal, symbl.Index)
	case BuiltinScope:
		self.emit(code.OpGetBuiltin, symbl.Index)
	case FreeScope:
		self.emit(code.OpGetFree, symbl.Index)
	case FunctionScope:
		self.emit(code.OpCurrentClosure)
	default:
		panic("[*Compiler::loadSymbol] : unhandled case")
	}
}
