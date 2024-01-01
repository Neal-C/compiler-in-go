package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tableTests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
	}

	for _, tt := range tableTests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length, want = %d, got = %d", len(tt.expected), len(instruction))
		}

		for index, instructionByte := range tt.expected {
			if instruction[index] != tt.expected[index] {
				t.Errorf("wrong instructionBye at position %d, want %d , got %d", index, instructionByte, instruction[index])
			}
		}
	}
}

func TestInstructionString(t *testing.T) {
	instructions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65_535),
	}

	expected := `000 OpConstant 1 003 OpConstant 2 006 OpConstant 65535`

	var concatted Instructions

	for _, instruction := range instructions {
		concatted = append(concatted, instruction...)
	}

	if concatted.String() != expected {
		t.Errorf("instruction wrongly formatted.\nwant = %q\ngot = %q", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tableTests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65_535}, 2},
	}

	for _, tt := range tableTests {
		instruction := Make(tt.op, tt.operands...)

		definition, err := LookUp(byte(tt.op))

		if err != nil {
			t.Fatalf("definition not found: %q \n", err)
		}

		operandsRead, numberOfBytesRead := ReadOperands(definition, instruction[1:])

		if numberOfBytesRead != tt.bytesRead {
			t.Fatalf("numberOfBytesRead wrong. want = %d . got = %d", tt.bytesRead, numberOfBytesRead)
		}

		for index, want := range tt.operands {
			if operandsRead[index] != want {
				t.Errorf("operand wrong. want = %d. got = %d", want, operandsRead[index])
			}
		}
	}
}
