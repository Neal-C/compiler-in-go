package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name          string
	OperandsWidth []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
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
	return ""
}
