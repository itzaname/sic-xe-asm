package parser

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"fmt"
	"strconv"
	"strings"
)

func (p *Parser) readOperand(opr string) (graph.Operand, error) {
	// Indirect addr
	if opr[0] == '#' {
		// Direct number
		if val, err := strconv.Atoi(opr[1:]); err == nil {
			return graph.Operand{
				Type: 2,
				Data: val,
			}, nil
		}

		return graph.Operand{}, fmt.Errorf("unkown data type: '%s'", opr)
	}

	// Try register
	if reg, err := machine.RegisterByName(strings.ToUpper(opr)); err == nil {
		return graph.Operand{
			Type: 0,
			Data: reg,
		}, nil
	}

	// Assume its a label
	return graph.Operand{
		Type: 1,
		Data: opr,
	}, nil
}

func (p *Parser) instructionFromToken(token []string) (*graph.InstructionNode, error) {
	node := graph.InstructionNode{}
	node.Debug.Tokens = len(token)
	baseSize := 0

	// Read label and instruction
	if ins, extended, err := machine.InstructionByName(strings.ToUpper(token[0])); err == nil {
		node.Instruction = ins
		node.Flags.E = extended
		baseSize = 1
	} else {
		ins, extended, err := machine.InstructionByName(strings.ToUpper(token[1]))
		if err != nil {
			return nil, err
		}
		node.Name = token[0]
		node.Instruction = ins
		node.Flags.E = extended
		baseSize = 2
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

		if len(args) < 2 {
			return &node, nil
		}

		opr, err := p.readOperand(args[1])
		if err != nil {
			return nil, err
		}
		node.Operands = append(node.Operands, opr)
	}

	// Format 3
	if node.Instruction.Format == 3 {
		args := strings.SplitN(token[baseSize], ",", 2)
		node.Debug.Args = len(args)

		// Operand 1
		opr, err := p.readOperand(args[0])
		if err != nil {
			return nil, err
		}
		node.Operands = append(node.Operands, opr)

		if len(args) < 2 {
			return &node, nil
		}

		opr, err = p.readOperand(args[1])
		if err != nil {
			return nil, err
		}
		node.Operands = append(node.Operands, opr)
	}

	return &node, nil
}

func (p *Parser) readStorageItem(directive *machine.Directive, data string) (*graph.Storage, error) {
	switch directive.Name {
	case "RESB":
	case "RESW":
		dataInput, err := p.readArrayStorageInput(data)
		if err != nil {
			return nil, err
		}

		switch dataInput.(type) {
		case string:
			return &graph.Storage{
				Type: 2,
				Data: dataInput,
			}, nil
		case int:
			return &graph.Storage{
				Type: 1,
				Data: dataInput,
			}, nil
		}
	case "WORD":

	case "BYTE":
	default:
		return nil, fmt.Errorf("non storage directive '%s'", directive.Name)
	}

	return nil, fmt.Errorf("non storage directive '%s'", directive.Name)
}

func (p *Parser) directiveFromToken(token []string) (*graph.DirectiveNode, error) {
	if len(token) < 3 {
		return nil, fmt.Errorf("expected 3 tokens got %d", len(token))
	}

	node := graph.DirectiveNode{}
	node.Name = token[0]

	directive, err := machine.DirectiveByName(token[1])
	if err != nil {
		return nil, err
	}
	node.Directive = directive

	if directive.Storage {
		item, err := p.readStorageItem(directive, token[2])
		if err != nil {
			return nil, err
		}
		node.Data = item
		return &node, nil
	}

	return &node, nil
}

func (p *Parser) nodeFromToken(token []string) (graph.Node, error) {
	if len(token) >= 1 {
		if _, _, err := machine.InstructionByName(strings.ToUpper(token[0])); err == nil {
			return p.instructionFromToken(token)
		}
	}
	if len(token) >= 2 {
		if _, _, err := machine.InstructionByName(strings.ToUpper(token[1])); err == nil {
			return p.instructionFromToken(token)
		}
		if _, err := machine.DirectiveByName(strings.ToUpper(token[1])); err == nil {
			return p.directiveFromToken(token)
		}
	}
	return nil, fmt.Errorf("expected instruction or directive")
}
