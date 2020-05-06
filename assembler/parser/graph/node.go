package graph

import "ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"

type Node interface {
	Label() string
	Size() int
	Valid() bool
	Get() interface{}
}

type DebugData struct {
	Line   int
	Tokens int
	Args   int
}

type InstructionNode struct {
	Name        string
	Instruction *machine.Instruction
	Operands    []Operand
	Flags       machine.Flags
	Debug       DebugData
	*Graph
}

type DirectiveNode struct {
	Name      string
	Directive *machine.Directive
	Data      interface{}
	Debug     DebugData
	*Graph
}

///////////////////////////////////////////////////////
// InstructionNode
///////////////////////////////////////////////////////

type Storage struct {
	// 1 - size
	// 2 - data
	Type uint8
	Data interface{}
}

type Operand struct {
	// 0 Register
	// 1 Node
	// 2 Immediate
	Type uint8
	Data interface{}
}

func (node *InstructionNode) Label() string {
	return node.Name
}

func (node *InstructionNode) Size() int {
	if node.Instruction.Format == 3 && node.Flags.E {
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

	switch node.Data.(type) {
	case int: // WORD
		return 3
	case byte:
		return 1
	case []byte:
		return len(node.Data.([]byte))
	}

	return -1
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
