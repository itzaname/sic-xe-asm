package machine

import (
	"fmt"
	"strings"
)

type Instruction struct {
	Name   string
	OpCode uint8
	Format uint8
}

// This was a pain to make
var insList = []Instruction{
	{Name: "ADD", OpCode: 0x18, Format: 3},
	{Name: "ADDF", OpCode: 0x58, Format: 3},
	{Name: "ADDR", OpCode: 0x90, Format: 2},
	{Name: "AND", OpCode: 0x40, Format: 3},
	{Name: "CLEAR", OpCode: 0x04, Format: 2},
	{Name: "COMP", OpCode: 0x28, Format: 3},
	{Name: "COMPF", OpCode: 0x88, Format: 3},
	{Name: "COMPR", OpCode: 0xA0, Format: 2},
	{Name: "DIV", OpCode: 0x24, Format: 3},
	{Name: "DIVF", OpCode: 0x64, Format: 3},
	{Name: "DIVR", OpCode: 0x9C, Format: 1},
	{Name: "FIX", OpCode: 0xC4, Format: 1},
	{Name: "FLOAT", OpCode: 0xC0, Format: 1},
	{Name: "HIO", OpCode: 0xF4, Format: 1},
	{Name: "J", OpCode: 0x3C, Format: 3},
	{Name: "JEQ", OpCode: 0x30, Format: 3},
	{Name: "JGT", OpCode: 0x34, Format: 3},
	{Name: "JLT", OpCode: 0x38, Format: 3},
	{Name: "JSUB", OpCode: 0x48, Format: 3},
	{Name: "LDA", OpCode: 0x00, Format: 3},
	{Name: "LDB", OpCode: 0x68, Format: 3},
	{Name: "LDCH", OpCode: 0x50, Format: 3},
	{Name: "LDF", OpCode: 0x70, Format: 3},
	{Name: "LDL", OpCode: 0x08, Format: 3},
	{Name: "LDS", OpCode: 0x6C, Format: 3},
	{Name: "LDT", OpCode: 0x74, Format: 3},
	{Name: "LDX", OpCode: 0x04, Format: 3},
	{Name: "LPS", OpCode: 0xD0, Format: 3},
	{Name: "MUL", OpCode: 0x20, Format: 3},
	{Name: "MULF", OpCode: 0x70, Format: 3},
	{Name: "MULR", OpCode: 0x98, Format: 2},
	{Name: "NORM", OpCode: 0xC8, Format: 1},
	{Name: "OR", OpCode: 0x44, Format: 3},
	{Name: "RD", OpCode: 0xD8, Format: 3},
	{Name: "RMO", OpCode: 0xAC, Format: 2},
	{Name: "RSUB", OpCode: 0x4C, Format: 3},
	{Name: "SHIFTL", OpCode: 0xA4, Format: 2},
	{Name: "SHIFTR", OpCode: 0xA8, Format: 2},
	{Name: "SIO", OpCode: 0xF0, Format: 1},
	{Name: "SSK", OpCode: 0xEC, Format: 3},
	{Name: "STA", OpCode: 0x0C, Format: 3},
	{Name: "STB", OpCode: 0x78, Format: 3},
	{Name: "STCH", OpCode: 0x54, Format: 3},
	{Name: "STF", OpCode: 0x80, Format: 3},
	{Name: "STI", OpCode: 0xD4, Format: 3},
	{Name: "STL", OpCode: 0x14, Format: 3},
	{Name: "STS", OpCode: 0x7C, Format: 3},
	{Name: "STSW", OpCode: 0xE8, Format: 3},
	{Name: "STT", OpCode: 0x84, Format: 3},
	{Name: "STX", OpCode: 0x10, Format: 3},
	{Name: "SUB", OpCode: 0x1C, Format: 3},
	{Name: "SUBF", OpCode: 0x5C, Format: 3},
	{Name: "SUBR", OpCode: 0x94, Format: 2},
	{Name: "SVC", OpCode: 0xB0, Format: 2},
	{Name: "TD", OpCode: 0xE0, Format: 3},
	{Name: "TIO", OpCode: 0xF8, Format: 1},
	{Name: "TIX", OpCode: 0x2C, Format: 3},
	{Name: "TIXR", OpCode: 0xB8, Format: 2},
	{Name: "WD", OpCode: 0xDC, Format: 3},
}

var insNameCache = map[string]*Instruction{}
var insCodeCache = map[uint8]*Instruction{}

func InstructionByOpCode(code uint8) (*Instruction, error) {
	if ins, ok := insCodeCache[code]; ok {
		return ins, nil
	}

	for i := 0; i < len(insList); i++ {
		if insList[i].OpCode == code {
			insCodeCache[insList[i].OpCode] = &insList[i]
			insNameCache[insList[i].Name] = &insList[i]
			return &insList[i], nil
		}
	}

	return nil, fmt.Errorf("invalid opcode value: '0x%X'", code)
}

func InstructionByName(name string) (*Instruction, bool, error) {
	name = strings.ToUpper(name)
	extended := false
	if name[0] == '+' {
		name = name[1:]
		extended = true
	}

	if ins, ok := insNameCache[name]; ok {
		return ins, extended, nil
	}

	for i := 0; i < len(insList); i++ {
		if insList[i].Name == name {
			insCodeCache[insList[i].OpCode] = &insList[i]
			insNameCache[insList[i].Name] = &insList[i]
			return &insList[i], extended, nil
		}
	}

	return nil, extended, fmt.Errorf("invalid opcode name: '%s'", name)
}
