package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

// When the VM executes OpConstant it retrieves the constant
// using the operand as an index and pushes it on to the stack.

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
	OpNull
	OpSetGlobal
	OpGetGlobal
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturnValue
)

type Definition struct {
	Name          string
	OperandsWidth []int // number of bytes (1 x 8, 2 x 8, ...)
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
	OpNull:          {"OpNull", []int{}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpArray:         {"OpArray", []int{2}},
	OpHash:          {"OpHash", []int{2}},
	OpIndex:         {"OpIndex", []int{}},
	OpCall:          {"OpCall", []int{}},
	OpReturnValue:   {"OpReturnValue", []int{}},
}

func LookUp(op byte) (*Definition, error) {
	definition, ok := definitions[Opcode(op)]

	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return definition, nil
}

// That's how you make bytecode

func Make(op Opcode, operands ...int) []byte {
	definition, ok := definitions[op]

	if !ok {
		return []byte{}
	}

	instructionLen := 1

	for _, width := range definition.OperandsWidth {
		instructionLen += width
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	for index, operand := range operands {

		width := definition.OperandsWidth[index]

		offset := 1

		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}

		offset += width
	}

	return instruction
}

func (self Instructions) String() string {
	var out bytes.Buffer

	index := 0

	for index < len(self) {
		definition, err := LookUp(self[index])

		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(definition, self[index+1:])

		fmt.Fprintf(&out, "%04d %s\n", index, self.fmtInstruction(definition, operands))

		index += (1 + read)
	}

	return out.String()

}

func (self Instructions) fmtInstruction(definition *Definition, operands []int) string {
	operandCount := len(definition.OperandsWidth)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return definition.Name
	case 1:
		return fmt.Sprintf("%s %d", definition.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", definition.Name)
}

func ReadOperands(definition *Definition, instructions Instructions) ([]int, int) {
	operands := make([]int, len(definition.OperandsWidth))

	offset := 0

	for index, width := range definition.OperandsWidth {
		switch width {
		case 2:
			operands[index] = int(ReadUint16(instructions[offset:]))
		}
		offset += width
	}

	return operands, offset
}

func ReadUint16(instruction Instructions) uint16 {
	return binary.BigEndian.Uint16(instruction)
}
