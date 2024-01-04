package vm

import (
	"fmt"
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

func (self *VM) StackTop() object.Object {
	if self.stackPointer == 0 {
		return nil
	}

	return self.stack[self.stackPointer-1]
}

func (self *VM) Run() error {

	for indexPointer := 0; indexPointer < len(self.instructions); indexPointer++ {

		// type coercion
		op := code.Opcode(self.instructions[indexPointer])

		switch op {
		case code.OpConstant:
			operandIndex := indexPointer + 1
			constIndex := code.ReadUint16(self.instructions[operandIndex:])

			// increment indexPointer by the bytes size (16 = 2 x bytes )
			indexPointer += 2

			err := self.push(self.constants[constIndex])

			if err != nil {
				return err
			}
		case code.OpAdd:
			right := self.pop()
			left := self.pop()

			rightValue := right.(*object.Integer).Value
			leftValue := left.(*object.Integer).Value

			result := leftValue + rightValue
			_ = self.push(&object.Integer{Value: result})
		}
	}

	return nil
}

func (self *VM) push(obj object.Object) error {

	if self.stackPointer >= StackSize {
		return fmt.Errorf("stack overflow : https://stackoverflow.com/")
	}

	self.stack[self.stackPointer] = obj
	self.stackPointer++

	return nil
}

func (self *VM) pop() object.Object {
	returnedObj := self.stack[self.stackPointer-1]
	self.stackPointer--
	return returnedObj
}
