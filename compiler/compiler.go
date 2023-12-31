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

func (self *Compiler) Compile(code ast.Node) error {
	return nil
}

func (self *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: self.instructions,
		Constants:    self.constants,
	}
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
