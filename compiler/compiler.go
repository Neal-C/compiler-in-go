package compiler

import (
	"github.com/Neal-C/compiler-in-go/ast"
	"github.com/Neal-C/compiler-in-go/code"
	"github.com/Neal-C/compiler-in-go/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
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
	case *ast.InfixExpression:
		err := self.Compile(node.Left)
		if err != nil {
			return err
		}

		err = self.Compile(node.Right)
		if err != nil {
			return err
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		self.emit(code.OpConstant, self.addConstants(integer))
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

	return position
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
