package parser

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"encoding/hex"
	"fmt"
	"strconv"
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

func (p *Parser) readStorageItem(directive *machine.Directive, data string) (*graph.Storage, error) {
	switch directive.Name {
	case "RESB":
		// size declaration
		val, err := strconv.Atoi(data)
		if err != nil {
			return nil, err
		}
		return &graph.Storage{
			Type: 2,
			Size: val,
			Data: nil,
		}, nil
	case "RESW":
		// size declaration
		val, err := strconv.Atoi(data)
		if err != nil {
			return nil, err
		}
		return &graph.Storage{
			Type: 2,
			Size: val,
			Data: nil,
		}, nil
	case "WORD":
		// Single
		if val, err := strconv.Atoi(data); err == nil {
			return &graph.Storage{
				Type: 2,
				Size: 1,
				Data: val,
			}, nil
		}
		return nil, fmt.Errorf("unkown storage directive '%s'", directive.Name)
	case "BYTE":
		input, class, err := p.readInputString(data)
		if err != nil {
			return nil, err
		}

		if class == 'X' {
			val, err := hex.DecodeString(input)
			if err != nil {
				fmt.Println(input)
				return nil, err
			}
			return &graph.Storage{
				Type: 1,
				Size: len(val),
				Data: val,
			}, nil
		}

		inputBytes := []byte(input)
		return &graph.Storage{
			Type: 2,
			Size: len(inputBytes),
			Data: inputBytes,
		}, nil
	default:
		return nil, fmt.Errorf("non storage directive '%s'", directive.Name)
	}
}

func (p *Parser) directiveFromToken(token []string) (*graph.DirectiveNode, error) {
	node := graph.DirectiveNode{}

	if directive, err := machine.DirectiveByName(token[0]); err == nil {
		node.Directive = directive
		node.Data = token[1]
		node.Debug.Source = strings.Join(token[:2], " ")
		node.Debug.Tokens = len(node.Debug.Source)
		return &node, nil
	}

	if directive, err := machine.DirectiveByName(token[1]); err == nil {
		node.Directive = directive
		node.Debug.Source = strings.Join(token[:3], " ")
		node.Debug.Tokens = len(node.Debug.Source)
		node.Name = token[0]
		if directive.Storage {
			item, err := p.readStorageItem(directive, token[2])
			if err != nil {
				return nil, err
			}
			node.Data = item
			return &node, nil
		}
		node.Data = token[2]
		return &node, nil
	}

	return nil, fmt.Errorf("invalid directive or directive format")
}

func (p *Parser) nodeFromToken(token []string) (graph.Node, error) {
	if len(token) >= 1 {
		if _, _, err := machine.InstructionByName(strings.ToUpper(token[0])); err == nil {
			return p.instructionFromToken(token)
		}
		if _, err := machine.DirectiveByName(strings.ToUpper(token[0])); err == nil {
			return p.directiveFromToken(token)
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
