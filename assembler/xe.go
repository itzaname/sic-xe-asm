package assembler

import (
	"bytes"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"fmt"
	"github.com/icza/bitio"
)

func calcPCDisp(pc, target int) (int, bool) {
	disp := target - pc
	if disp < -0x800 || disp > 0x7FF {
		return 0, false
	}
	if disp < 0 {
		disp = disp + 0x1000
	}
	return disp, true
}

func (asm *Assembler) extendedInstruction(node *graph.InstructionNode) ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := bitio.NewWriter(&buffer)

	if len(node.Operands) < 1 {
		return nil, fmt.Errorf("need at least 1 operand for format 3 instruction")
	}

	// Support extended size
	nodeSize := node.Size()
	pc := node.Address() + nodeSize
	formatSize := 12
	if node.Flags.E == 1 {
		formatSize = 20
	}

	// Write opcode
	writer.WriteBits(uint64(node.Instruction.OpCode), 6)

	var value int
	// Set value and basic addressing
	switch node.Operands[0].Type {
	case 1:
		if target, ok := node.Operands[0].Data.(*graph.Node); ok {
			tmp := *target
			value = tmp.Address()
			switch node.Operands[0].Addressing {
			case 0: // Direct op m
				node.Flags.N = 1
				node.Flags.I = 1
				break
			case 1: // Indirect op @m
				node.Flags.N = 1
				break
			case 2: // Immediate op #m
				node.Flags.I = 1
				break
			}
			break
		}
		return nil, fmt.Errorf("unresolved node at line #%d at address 0x%X", node.Debug.Line, node.Address())
	case 2:
		if node.Operands[0].Addressing != 2 {
			return nil, fmt.Errorf("bad adresing mode #%d at address 0x%X: got '%d' wanted '2'", node.Debug.Line, node.Address(), node.Operands[0].Addressing)
		}
		node.Flags.I = 1
		if data, ok := node.Operands[0].Data.(int); ok {
			value = data
			break
		}
		return nil, fmt.Errorf("invalid node at line #%d at address 0x%X", node.Debug.Line, node.Address())
	default:
		return nil, fmt.Errorf("unknow adressing mode for format 3 type '%d' line #%d", node.Operands[0].Type, node.Debug.Line)
	}

	// Set X flag if needed
	if len(node.Operands) > 1 {
		if node.Operands[1].Type != 0 {
			return nil, fmt.Errorf("second argument for format 3 must be X register got type '%d' line #%d", node.Operands[0].Type, node.Debug.Line)
		}
		if opr, ok := node.Operands[1].Data.(*machine.Register); ok {
			if opr.Name != "X" {
				return nil, fmt.Errorf("second argument for format 3 must be X register got '%s' line #%d", opr.Name, node.Debug.Line)
			}
			node.Flags.X = 1
		} else {
			return nil, fmt.Errorf("second argument for format 3 must be X register got invalid type line #%d", node.Debug.Line)
		}
	}

	// Fix addressing
	if node.Flags.E == 0 && node.Operands[0].Type == 1 {
		// Attempt PC
		disp, ok := calcPCDisp(pc, value)
		if ok {
			value = disp
			node.Flags.P = 1
		} else {
			return nil, fmt.Errorf("NEED BASE")
		}
	}

	if err := writer.WriteBits(uint64(node.Flags.N), 1); err != nil {
		return nil, err
	}
	if err := writer.WriteBits(uint64(node.Flags.I), 1); err != nil {
		return nil, err
	}
	if err := writer.WriteBits(uint64(node.Flags.X), 1); err != nil {
		return nil, err
	}
	if err := writer.WriteBits(uint64(node.Flags.B), 1); err != nil {
		return nil, err
	}
	if err := writer.WriteBits(uint64(node.Flags.P), 1); err != nil {
		return nil, err
	}
	if err := writer.WriteBits(uint64(node.Flags.E), 1); err != nil {
		return nil, err
	}

	if err := writer.WriteBits(uint64(value), uint8(formatSize)); err != nil {
		return nil, err
	}

	writer.Close()
	return buffer.Bytes(), nil
}
