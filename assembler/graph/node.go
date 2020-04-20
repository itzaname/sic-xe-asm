package graph

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
)

type Node struct {
	Label       string
	Instruction *machine.Instruction
	Operands    []Operand
	Flags       machine.Flags
	*Graph
}

type Operand struct {
	// 0 Register
	// 1 Node
	// 2 Immediate
	Type uint8
	Data interface{}
}

func (node *Node) Address() int {
	count := 0
	for i := 0; i < len(node.Graph.Nodes); i++ {
		count += node.Graph.Nodes[i].Size()
		if &node.Graph.Nodes[i] == node {
			return count
		}
	}

	return -1
}

func (node *Node) Valid() bool {
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

func (node *Node) Size() int {
	if node.Instruction.Format == 3 && node.Flags.E {
		return 4
	}

	return int(node.Instruction.Format)
}
