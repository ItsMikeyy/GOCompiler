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
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	//Get opcode defenition
	def, ok := definitions[Opcode(op)]
	//Check if exists
	if !ok {
		return nil, fmt.Errorf("opcode %v is not defined", op)
	}

	return def, nil

}

func Make(op Opcode, operands ...int) []byte {
	//Get opcode defenition
	def, ok := definitions[op]

	//Check if exists
	if !ok {
		return []byte{}
	}

	//Get correct size fo byte slice
	var instructionLen int = 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	//Create byte slice
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	//Iterate over opearnds convert to binary depending on width
	//Add width to offset to get to next op
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}
