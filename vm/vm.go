package vm

import (
	"fmt"
	"github.com/Neal-C/compiler-in-go/code"
	"github.com/Neal-C/compiler-in-go/compiler"
	"github.com/Neal-C/compiler-in-go/object"
)

const StackSize = 2048

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := self.executeBinaryOperation(op)

			if err != nil {
				return err
			}
		case code.OpTrue:

			err := self.push(True)

			if err != nil {
				return err
			}

		case code.OpFalse:

			err := self.push(False)

			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := self.executeComparison(op)

			if err != nil {
				return err
			}

		case code.OpPop:
			self.pop()
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

func (self *VM) LastPoppedStackElement() object.Object {
	return self.stack[self.stackPointer]
}

func (self *VM) executeBinaryOperation(op code.Opcode) error {
	right := self.pop()
	left := self.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return self.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)

}

func (self *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	rightValue := right.(*object.Integer).Value
	leftValue := left.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operation: %d", op)
	}
	return self.push(&object.Integer{Value: result})
}

func (self *VM) executeComparison(op code.Opcode) error {
	rightHandSign := self.pop()
	leftHandSign := self.pop()

	if leftHandSign.Type() == object.INTEGER_OBJ && rightHandSign.Type() == object.INTEGER_OBJ {
		return self.executeIntegerComparison(op, leftHandSign, rightHandSign)
	}

	switch op {
	case code.OpEqual:
		return self.push(nativeBoolToBooleanObject(leftHandSign == rightHandSign))
	case code.OpNotEqual:
		return self.push(nativeBoolToBooleanObject(leftHandSign != rightHandSign))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, leftHandSign.Type(), rightHandSign.Type())
	}
}

func (self *VM) executeIntegerComparison(op code.Opcode, leftHandSign object.Object, rightHandSign object.Object) error {
	leftValue := leftHandSign.(*object.Integer).Value
	rightValue := rightHandSign.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return self.push(nativeBoolToBooleanObject(leftValue == rightValue))
	case code.OpNotEqual:
		return self.push(nativeBoolToBooleanObject(leftValue != rightValue))
	case code.OpGreaterThan:
		return self.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unkown op: %d", op)

	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	} else {
		return False
	}
}
