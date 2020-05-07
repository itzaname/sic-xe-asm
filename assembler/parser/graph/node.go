package graph

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"fmt"
	"strconv"
)

type Node interface {
	Label() string
	Size() int
	Valid() bool
	Get() interface{}
	Address() int
}

type DebugData struct {
	Line   int
	Tokens int
	Args   int
	Source string
}

type InstructionNode struct {
	Name        string
	Instruction *machine.Instruction
	Operands    []Operand
	Flags       machine.Flags
	Debug       DebugData
	Addr        int
	*Graph
}

type DirectiveNode struct {
	Name      string
	Directive *machine.Directive
	Data      interface{}
	Debug     DebugData
	Addr      int
	*Graph
}

///////////////////////////////////////////////////////
// InstructionNode
///////////////////////////////////////////////////////

type Storage struct {
	// 1 - byte
	// 2 - word
	Type uint8
	Size int
	Data interface{}
}

type Operand struct {
	// 0 Register
	// 1 Label
	// 2 Value
	// 3 Literal
	Type uint8
	// 0 Direct
	// 1 Indirect
	// 2 Immediate
	Addressing uint8
	Data       interface{}
}

func (node *InstructionNode) Label() string {
	return node.Name
}

func (node *InstructionNode) Size() int {
	if node.Instruction.Format == 3 && node.Flags.E == 1 {
		return 4
	}

	return int(node.Instruction.Format)
}

func (node *InstructionNode) Valid() bool {
	if node.Graph == nil || node.Instruction == nil {
		return false
	}

	switch node.Instruction.Format {
	case 1:
		return len(node.Operands) == 0
	case 2:
		if len(node.Operands) == 2 {
			return node.Operands[0].Type == 0 && node.Operands[1].Type == 0
		}
		return false
	case 3:
		if len(node.Operands) == 1 {
			return node.Operands[0].Type != 0
		}
		return false
	}

	return false
}

func (node *InstructionNode) Get() interface{} {
	return node
}

func (node *InstructionNode) Address() int {
	return node.Addr
}

///////////////////////////////////////////////////////
// DirectiveNode
///////////////////////////////////////////////////////

func (node *DirectiveNode) Label() string {
	return node.Name
}

func (node *DirectiveNode) Size() int {
	if !node.Directive.Storage {
		return 0
	}
	if store, ok := node.Data.(*Storage); ok {
		return store.Size
	}
	return 0
}

func (node *DirectiveNode) Valid() bool {
	if node.Directive != nil {
		if node.Directive.Storage && node.Data == nil {
			return false
		}
		return true
	}

	return false
}

func (node *DirectiveNode) Get() interface{} {
	return node
}

func (node *DirectiveNode) Address() int {
	// This is a hack and I hate it
	if node.Directive.Name == "EQU" {
		readSingle := func(item string) int {
			// If requesting our address
			if item == "*" {
				return node.Addr
			}
			// Arg is number
			if val, err := strconv.Atoi(item); err == nil {
				return val
			}
			// Is another label
			if val, ok := node.SymTable[item]; ok {
				return val.Address()
			}

			panic(fmt.Sprintln("FAILED TO RESOLVE EQU EXPRESSION", item))
		}
		args := node.Data.([]string)
		if len(args) == 1 {
			return readSingle(args[0])
		} else {
			opr1 := readSingle(args[0])
			opr2 := readSingle(args[2])

			switch args[1] {
			case "+":
				return opr1 + opr2
			case "-":
				return opr1 - opr2
			case "/":
				return opr1 / opr2
			}
		}

		panic(fmt.Sprintln("FAILED TO RESOLVE EQU EXPRESSION", args))
	}
	return node.Addr
}
