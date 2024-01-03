package vm

import (
	"github.com/Neal-C/compiler-in-go/code"
	"github.com/Neal-C/compiler-in-go/compiler"
	"github.com/Neal-C/compiler-in-go/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	stackPointer int // stackPointer always points to the next value, top of stack is (stackPointer - 1)
}

func New(bytecode *compiler.ByteCode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]object.Object, StackSize),
		stackPointer: 0,
	}
}
