package parser

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func (p *Parser) equHandler(directive *machine.Directive, token []string) (*graph.DirectiveNode, error) {
	node := graph.DirectiveNode{}
	node.Directive = directive
	node.Debug.Source = strings.Join(token[:3], " ")
	node.Debug.Tokens = len(node.Debug.Source)
	node.Name = token[0]
	node.Graph = p.nodeGraph

	args := []string{}
	wasSplit := false
	isDelimiter := func(char uint8) bool {
		return char == '+' || char == '-' || char == '*' || char == '/'
	}

	// ITEM(SEP)ITEM
	// *-BASE
	if len(token[2]) > 1 {
		for i := 0; i < len(token[2]); i++ {
			if isDelimiter(token[2][i]) {
				wasSplit = true
				args = append(args, token[2][:i])
				args = append(args, string(token[2][i]))
				args = append(args, token[2][i+1:])
			}
		}
	}
	if !wasSplit {
		args = append(args, token[2])
	}

	if len(args) != 1 && len(args) != 3 {
		return nil, fmt.Errorf("needed 1 or 3 args got %d", len(args))
	}

	node.Data = args

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
			Size: val * 3,
			Data: nil,
		}, nil
	case "WORD":
		// Single
		if val, err := strconv.Atoi(data); err == nil {
			return &graph.Storage{
				Type: 2,
				Size: 3,
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

var directiveList map[string]func(directive *machine.Directive, token []string) (*graph.DirectiveNode, error)

func (p *Parser) directiveFromToken(token []string) (*graph.DirectiveNode, error) {
	// First run
	if directiveList == nil {
		directiveList = map[string]func(directive *machine.Directive, token []string) (*graph.DirectiveNode, error){}

		// Things will get painful to do all these special cases for so it's better
		// to just tackle it in the parser with special functions
		directiveList["EQU"] = p.equHandler
		// I never did anything else with this
		// I had big plans
	}

	// Do work
	node := graph.DirectiveNode{}

	// Read no label directive
	if directive, err := machine.DirectiveByName(token[0]); err == nil {
		if f, ok := directiveList[directive.Name]; ok {
			return f(directive, token)
		}
		node.Directive = directive
		if len(token) == 1 {
			node.Debug.Source = strings.Join(token[:1], " ")
			node.Debug.Tokens = len(node.Debug.Source)
			return &node, nil
		}
		node.Data = token[1]
		node.Debug.Source = strings.Join(token[:2], " ")
		node.Debug.Tokens = len(node.Debug.Source)
		return &node, nil
	}

	// Read label directive
	if directive, err := machine.DirectiveByName(token[1]); err == nil {
		if f, ok := directiveList[directive.Name]; ok {
			return f(directive, token)
		}
		node.Directive = directive
		node.Debug.Source = strings.Join(token[:3], " ")
		node.Debug.Tokens = len(node.Debug.Source)
		node.Name = token[0]
		if directive.Storage { // Read storage directives separately
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
