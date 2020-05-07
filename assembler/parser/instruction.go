package parser

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"fmt"
	"strings"
)

func (p *Parser) instructionFromToken(token []string) (*graph.InstructionNode, error) {
	node := graph.InstructionNode{}
	baseSize := 0
	tokenSize := 0

	// Read label and instruction
	if ins, extended, err := machine.InstructionByName(strings.ToUpper(token[0])); err == nil {
		node.Instruction = ins
		node.Flags.E = extended
		baseSize = 1
		tokenSize = baseSize
	} else {
		ins, extended, err := machine.InstructionByName(strings.ToUpper(token[1]))
		if err != nil {
			return nil, err
		}
		node.Name = token[0]
		node.Instruction = ins
		node.Flags.E = extended
		baseSize = 2
		tokenSize = baseSize
	}

	// Format 2
	if node.Instruction.Format == 2 {
		if len(token) < baseSize+1 {
			return nil, fmt.Errorf("expected argument")
		}
		args := strings.SplitN(token[baseSize], ",", 2)
		node.Debug.Args = len(args)

		// Operand 1
		reg, err := machine.RegisterByName(strings.ToUpper(args[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid register '%s'", args[0])
		}
		node.Operands = append(node.Operands, graph.Operand{
			Type: 0,
			Data: reg,
		})
		tokenSize++

		if len(args) < 2 {
			node.Debug.Source = strings.Join(token[:tokenSize], " ")
			node.Debug.Tokens = len(node.Debug.Source)
			return &node, nil
		}

		opr, err := p.readOperand(args[1])
		if err != nil {
			return nil, err
		}
		node.Operands = append(node.Operands, opr)
	}

	// Format 3
	if node.Instruction.Format == 3 && !node.Instruction.Special {
		args := strings.SplitN(token[baseSize], ",", 2)
		node.Debug.Args = len(args)

		// Operand 1
		opr, err := p.readOperand(args[0])
		if err != nil {
			return nil, err
		}
		node.Operands = append(node.Operands, opr)

		tokenSize++

		if len(args) < 2 {
			node.Debug.Source = strings.Join(token[:tokenSize], " ")
			node.Debug.Tokens = len(node.Debug.Source)
			return &node, nil
		}

		opr, err = p.readOperand(args[1])
		if err != nil {
			return nil, err
		}
		node.Operands = append(node.Operands, opr)
		tokenSize++
	}

	node.Debug.Source = strings.Join(token[:tokenSize], " ")
	node.Debug.Tokens = len(node.Debug.Source)

	return &node, nil
}
