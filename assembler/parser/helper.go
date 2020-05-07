package parser

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"fmt"
	"strconv"
	"strings"
)

func (p *Parser) readInputString(input string) (string, uint8, error) {
	if len(input) < 4 {
		return "", ' ', fmt.Errorf("expected data got: %s", input)
	}

	// X'DEADBEEF'
	if input[0] == 'X' && input[1] == '\'' {
		for i := 2; i < len(input); i++ {
			if input[i] == '\'' && input[i-1] != '\\' {
				return input[2:i], 'X', nil
			}
		}
	}

	// C'BIG STRING'
	if input[0] == 'C' && input[1] == '\'' {
		for i := 2; i < len(input); i++ {
			if input[i] == '\'' && input[i-1] != '\\' {
				return input[2:i], 'C', nil
			}
		}
	}

	return "nil", ' ', fmt.Errorf("could not read input: %s", input)
}

func (p *Parser) readOperand(opr string) (graph.Operand, error) {
	// Immediate addressing
	if opr[0] == '#' {
		if num, err := strconv.Atoi(opr[1:]); err == nil {
			return graph.Operand{
				Type:       2,
				Addressing: 2,
				Data:       num,
			}, nil
		}

		return graph.Operand{
			Type:       1,
			Addressing: 2,
			Data:       opr[1:],
		}, nil
	}

	// Indirect addressing
	if opr[0] == '@' {
		return graph.Operand{
			Type:       1,
			Addressing: 1,
			Data:       opr[1:],
		}, nil
	}

	// Literal
	if opr[0] == '=' {
		input, _, err := p.readInputString(opr[1:])
		if err != nil {
			return graph.Operand{}, err
		}
		return graph.Operand{
			Type:       3,
			Addressing: 0,
			Data:       input,
		}, nil
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
		Type:       1,
		Addressing: 0,
		Data:       opr,
	}, nil
}
