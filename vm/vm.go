package vm

import (
	"fmt"
	"github.com/Neal-C/compiler-in-go/code"
	"github.com/Neal-C/compiler-in-go/compiler"
	"github.com/Neal-C/compiler-in-go/object"
)

const StackSize = 2048
const GlobalSize = 65536
const MaxFrames = 1024

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants    []object.Object
	stack        []object.Object
	stackPointer int // stackPointer always points to the next value, top of stack is (stackPointer - 1)
	globals      []object.Object
	frames       []*Frame
	framesIndex  int
}

type Frame struct {
	fn           *object.CompiledFunction
	indexPointer int
	basePointer  int
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{
		fn:           fn,
		indexPointer: -1,
		basePointer:  basePointer,
	}
}

func (self *Frame) Instructions() code.Instructions {
	return self.fn.Instructions
}

func New(bytecode *compiler.ByteCode) *VM {

	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		stackPointer: 0,
		globals:      make([]object.Object, GlobalSize),
		frames:       frames,
		framesIndex:  1,
	}
}

func NewWithGlobalStore(bytecode *compiler.ByteCode, globals []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = globals
	return vm
}

func (self *VM) StackTop() object.Object {
	if self.stackPointer == 0 {
		return nil
	}

	return self.stack[self.stackPointer-1]
}

func (self *VM) Run() error {

	var indexPointer int
	var instructions code.Instructions
	var op code.Opcode

	for self.currentFrame().indexPointer < len(self.currentFrame().Instructions())-1 {

		self.currentFrame().indexPointer++

		indexPointer = self.currentFrame().indexPointer

		instructions = self.currentFrame().Instructions()
		// type coercion
		op = code.Opcode(instructions[indexPointer])

		switch op {
		case code.OpConstant:
			operandIndex := indexPointer + 1
			constIndex := code.ReadUint16(instructions[operandIndex:])

			// increment indexPointer by the bytes size (16 = 2 x bytes )
			self.currentFrame().indexPointer += 2

			err := self.push(self.constants[constIndex])

			if err != nil {
				return err
			}

		case code.OpPop:

			self.pop()

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

		case code.OpBang:

			err := self.executeBangOperator()

			if err != nil {
				return err
			}

		case code.OpMinus:

			err := self.executeMinusOperator()

			if err != nil {
				return err
			}

		case code.OpJump:

			position := int(code.ReadUint16(instructions[indexPointer+1:]))
			self.currentFrame().indexPointer = position - 1
		case code.OpJumpNotTruthy:

			position := int(code.ReadUint16(instructions[indexPointer+1:]))
			self.currentFrame().indexPointer += 2

			condition := self.pop()

			if !isTruthy(condition) {
				self.currentFrame().indexPointer = position - 1
			}

		case code.OpNull:

			err := self.push(Null)

			if err != nil {
				return err
			}
		case code.OpSetGlobal:

			globalIndex := code.ReadUint16(instructions[indexPointer+1:])

			self.currentFrame().indexPointer += 2

			self.globals[globalIndex] = self.pop()

		case code.OpGetGlobal:

			globalIndex := code.ReadUint16(instructions[indexPointer+1:])

			self.currentFrame().indexPointer += 2

			resolvedValue := self.globals[globalIndex]

			err := self.push(resolvedValue)

			if err != nil {
				return err
			}

		case code.OpArray:

			numberOfElements := int(code.ReadUint16(instructions[indexPointer+1:]))

			self.currentFrame().indexPointer += 2

			array := self.buildArray(self.stackPointer-numberOfElements, self.stackPointer)

			self.stackPointer = self.stackPointer - numberOfElements

			err := self.push(array)

			if err != nil {
				return err
			}

		case code.OpHash:

			numberOfElements := int(code.ReadUint16(instructions[indexPointer+1:]))

			self.currentFrame().indexPointer += 2

			hash, err := self.buildHash(self.stackPointer-numberOfElements, self.stackPointer)

			if err != nil {
				return err
			}

			self.stackPointer = self.stackPointer - numberOfElements

			err = self.push(hash)

			if err != nil {
				return err
			}
		case code.OpIndex:

			index := self.pop()
			left := self.pop()

			err := self.executeIndexOperation(left, index)

			if err != nil {
				return err
			}

		case code.OpCall:

			numberOfArguments := code.ReadUint8(instructions[indexPointer+1:])

			self.currentFrame().indexPointer += 1

			fn, ok := self.stack[self.stackPointer-1-int(numberOfArguments)].(*object.CompiledFunction)

			if !ok {
				return fmt.Errorf("calling a non-function")
			}

			newFrame := NewFrame(fn, self.stackPointer)

			self.pushFrame(newFrame)
			self.stackPointer = newFrame.basePointer + fn.NumberOfLocals
		case code.OpReturnValue:

			returnValue := self.pop()

			frame := self.popFrame()
			self.stackPointer = frame.basePointer - 1

			err := self.push(returnValue)

			if err != nil {
				return err
			}

		case code.OpReturn:

			frame := self.popFrame()
			self.stackPointer = frame.basePointer - 1

			err := self.push(Null)

			if err != nil {
				return err
			}

		case code.OpSetLocal:

			localIndex := code.ReadUint8(instructions[indexPointer+1:])

			self.currentFrame().indexPointer += 1

			frame := self.currentFrame()

			self.stack[frame.basePointer+int(localIndex)] = self.pop()

		case code.OpGetLocal:

			localIndex := code.ReadUint8(instructions[indexPointer+1:])

			self.currentFrame().indexPointer += 1

			frame := self.currentFrame()

			localBinding := self.stack[frame.basePointer+int(localIndex)]

			err := self.push(localBinding)

			if err != nil {
				return err
			}

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

	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return self.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return self.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)

	}

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

func (self *VM) executeBangOperator() error {
	operandee := self.pop()

	switch operandee {
	case True:
		return self.push(False)
	case False:
		return self.push(True)
	case Null:
		return self.push(True)
	default:
		return self.push(False)
	}
}

func (self *VM) executeMinusOperator() error {
	operandee := self.pop()

	if operandee.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operandee.Type())
	}

	value := operandee.(*object.Integer).Value

	return self.push(&object.Integer{Value: -value})
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false

	default:
		return true
	}
}

func (self *VM) executeBinaryStringOperation(op code.Opcode, left object.Object, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return self.push(&object.String{Value: leftValue + rightValue})
}

func (self *VM) buildArray(startIndex int, endIndex int) object.Object {

	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = self.stack[i]
	}
	return &object.Array{Elements: elements}
}

func (self *VM) buildHash(startIndex int, endIndex int) (object.Object, error) {

	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := self.stack[i]
		value := self.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}

		hashkey, ok := key.(object.Hashable)

		if !ok {
			return nil, fmt.Errorf("unusable as a hash key: %s", key.Type())
		}

		hashedPairs[hashkey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil

}

func (self *VM) executeIndexOperation(left object.Object, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return self.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return self.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported for : %s", left.Type())
	}
}

func (self *VM) executeArrayIndex(array object.Object, index object.Object) error {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value

	maxBoundary := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > maxBoundary {
		return self.push(Null)
	}

	return self.push(arrayObject.Elements[idx])
}

func (self *VM) executeHashIndex(hash object.Object, index object.Object) error {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)

	if !ok {
		return fmt.Errorf("unusable as a key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]

	if !ok {
		return self.push(Null)
	}

	return self.push(pair.Value)

}

func (self *VM) currentFrame() *Frame {
	return self.frames[self.framesIndex-1]
}

func (self *VM) pushFrame(frame *Frame) {
	self.frames[self.framesIndex] = frame
	self.framesIndex++
}

func (self *VM) popFrame() *Frame {
	self.framesIndex--
	return self.frames[self.framesIndex]
}
