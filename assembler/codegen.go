package assembler

import (
	"bytes"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"encoding/hex"
	"fmt"
	"github.com/icza/bitio"
	"strings"
)

func (asm *Assembler) instructionBytes(node *graph.InstructionNode) ([]byte, error) {
	switch node.Instruction.Format {
	case 1:
		buffer := bytes.Buffer{}
		writer := bitio.NewWriter(&buffer)
		err := writer.WriteByte(node.Instruction.OpCode)
		writer.Close()
		return buffer.Bytes(), err
	case 2:
		buffer := bytes.Buffer{}
		writer := bitio.NewWriter(&buffer)
		defer writer.Close()
		// Always do the instruction
		if err := writer.WriteByte(node.Instruction.OpCode); err != nil {
			return nil, err
		}

		// If > 0 args do this
		if len(node.Operands) > 0 {
			// Op 1
			if node.Operands[0].Type != 0 {
				return nil, fmt.Errorf("invalid operand 1 type: '%d'", node.Operands[0].Type)
			}
			if err := writer.WriteBits(uint64(node.Operands[0].Data.(*machine.Register).Number), 4); err != nil {
				return nil, err
			}
			// Op 2
			if len(node.Operands) > 1 {
				if node.Operands[1].Type != 0 {
					return nil, fmt.Errorf("invalid operand 2 type: '%d'", node.Operands[1].Type)
				}
				if err := writer.WriteBits(uint64(node.Operands[1].Data.(*machine.Register).Number), 4); err != nil {
					return nil, err
				}
				writer.Close()
				return buffer.Bytes(), nil
			} else {
				err := writer.WriteBits(0, 4)
				writer.Close()
				return buffer.Bytes(), err
			}

		}
		err := writer.WriteByte(0)
		writer.Close()
		return buffer.Bytes(), err
	case 3:
		return asm.extendedInstruction(node)
	}

	return nil, fmt.Errorf("invalid instruction format: '%d'", node.Instruction.Format)
}

func (asm *Assembler) generateObjectItems() error {
	asm.graph.UpdateAddr()
	iter := asm.graph.Iterator()
	for iter.Next() {
		switch iter.Node().(type) {
		case *graph.DirectiveNode:
			node := iter.Node().(*graph.DirectiveNode)
			fmt.Println("ITEM:", node.Address(), node.Debug.Source)
			break
		case *graph.InstructionNode:
			node := iter.Node().(*graph.InstructionNode)
			data, err := asm.instructionBytes(node)
			if err != nil {
				return err
			}
			fmt.Println("ITEM:", node.Address(), node.Debug.Source, strings.ToUpper(hex.EncodeToString(data)), "FLAGS", node.Flags)

			break
		default:
			return fmt.Errorf("invalid node '#%d' at 0x%X", iter.Index(), iter.Address())
		}
	}

	return nil
}
