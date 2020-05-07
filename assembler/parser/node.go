package parser

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"fmt"
	"strings"
)

func (p *Parser) nodeFromToken(token []string) (graph.Node, error) {
	if len(token) >= 1 { // No label
		if _, _, err := machine.InstructionByName(strings.ToUpper(token[0])); err == nil {
			return p.instructionFromToken(token)
		}
		if _, err := machine.DirectiveByName(strings.ToUpper(token[0])); err == nil {
			return p.directiveFromToken(token)
		}
	}
	if len(token) >= 2 { // Has label
		if _, _, err := machine.InstructionByName(strings.ToUpper(token[1])); err == nil {
			return p.instructionFromToken(token)
		}
		if _, err := machine.DirectiveByName(strings.ToUpper(token[1])); err == nil {
			return p.directiveFromToken(token)
		}
	}
	return nil, fmt.Errorf("expected instruction or directive")
}
