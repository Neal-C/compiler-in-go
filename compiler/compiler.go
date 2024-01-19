package compiler

import (
	"fmt"
	"github.com/Neal-C/compiler-in-go/ast"
	"github.com/Neal-C/compiler-in-go/code"
	"github.com/Neal-C/compiler-in-go/object"
)

type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
}

type EmittedInstruction struct {
	OpCode   code.Opcode
	Position int
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
	}
}

func NewWithState(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = symbolTable
	compiler.constants = constants
	return compiler
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

		if self.lastInstructionIsPop() {
			self.removeLastPop()
		}

		jumpOverAlternativePosition := self.emit(code.OpJump, 9999)

		afterConsequencePos := len(self.instructions)
		self.changeOperand(jumpNotTruthyPosition, afterConsequencePos)

		if node.Alternative == nil {
			self.emit(code.OpNull)
		} else {

			err = self.Compile(node.Alternative)

			if err != nil {
				return err
			}

			if self.lastInstructionIsPop() {
				self.removeLastPop()
			}
		}

		afterAlternativePosition := len(self.instructions)
		self.changeOperand(jumpOverAlternativePosition, afterAlternativePosition)

	case *ast.BlockStatement:

		for _, stmt := range node.Statements {
			err := self.Compile(stmt)

			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:

		err := self.Compile(node.Value)

		if err != nil {
			return err
		}

		symbol := self.symbolTable.Define(node.Name.Value)
		self.emit(code.OpSetGlobal, symbol.Index)
	case *ast.Identifier:

		symbol, ok := self.symbolTable.Resolve(node.Value)

		if !ok {
			// Compile time errors !!
			return fmt.Errorf("undefined variable : %s", node.Value)
		}

		self.emit(code.OpGetGlobal, symbol.Index)
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
	}

	return nil
}

func (self *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: self.instructions,
		Constants:    self.constants,
	}
}

func (self *Compiler) addConstants(obj object.Object) int {
	self.constants = append(self.constants, obj)

	return len(self.constants) - 1
}

func (self *Compiler) addInstruction(instructions code.Instructions) int {
	posNewInstruction := len(self.instructions)
	self.instructions = append(self.instructions, instructions...)

	return posNewInstruction
}

func (self *Compiler) emit(op code.Opcode, operands ...int) int {

	instruction := code.Make(op, operands...)
	position := self.addInstruction(instruction)

	self.setLastInstruction(op, position)

	return position
}

func (self *Compiler) setLastInstruction(op code.Opcode, position int) {

	previous := self.lastInstruction
	last := EmittedInstruction{OpCode: op, Position: position}

	self.previousInstruction = previous
	self.lastInstruction = last
}

func (self *Compiler) lastInstructionIsPop() bool {
	return self.lastInstruction.OpCode == code.OpPop
}

func (self *Compiler) removeLastPop() {

	self.instructions = self.instructions[:self.lastInstruction.Position]

	self.lastInstruction = self.previousInstruction
}

func (self *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for index := 0; index < len(newInstruction); index++ {
		self.instructions[pos+index] = newInstruction[index]
	}
}

func (self *Compiler) changeOperand(operationPosition int, operand int) {

	// operationPosition is where we were
	op := code.Opcode(self.instructions[operationPosition])
	// operand is 2 bytes
	newInstruction := code.Make(op, operand)

	self.replaceInstruction(operationPosition, newInstruction)
}
