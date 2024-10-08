package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}}, //Init test with opconstant. 0XFE not 0xFF (to check endian) encode 255,254 to bytes = 0xFE
	}

	for _, tt := range tests {
		//Pass opcode and [65534] to make
		instruction := Make(tt.op, tt.operands...)

		//Check if byteslice is correct size
		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. Wanted=%v, Got=%v", len(tt.expected), len(instruction))
		}

		//Check endian
		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("Wrong byte at pos %d. Wanted=%d, Got=%d", i, b, instruction[i])
			}
		}
	}
}
