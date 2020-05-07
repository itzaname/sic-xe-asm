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
	case 1: // Basic format 1
		buffer := bytes.Buffer{}
		writer := bitio.NewWriter(&buffer)
		err := writer.WriteByte(node.Instruction.OpCode)
		writer.Close()
		return buffer.Bytes(), err
	case 2: // Register funcs
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
	startAddr := -1
	buffer := bytes.Buffer{}
	bufferSize := 0

	// Helper function to avoid having the code in the loop
	// manages the creation of text records
	writeTextRecord := func() {
		if bufferSize == 0 {
			return
		}
		asm.obj.Text = append(asm.obj.Text, machine.TextRecord{
			Start:  startAddr,
			Length: bufferSize,
			Object: strings.ToUpper(hex.EncodeToString(buffer.Bytes())),
		})
		buffer.Reset()
		startAddr = -1
		bufferSize = 0
		asm.flags.endModule = false
	}

	iter := asm.graph.Iterator()
	for iter.Next() {
		// 28 bytes is 56 columns
		if bufferSize >= 28 || asm.flags.endModule {
			writeTextRecord()
		}
		switch iter.Node().(type) {
		case *graph.DirectiveNode: // Directive code gen
			node := iter.Node().(*graph.DirectiveNode)
			data, write, err := asm.handleDirective(node)
			if err != nil {
				return fmt.Errorf("line %d: %s", node.Debug.Line, err.Error())
			}
			if write {
				n, err := buffer.Write(data)
				if err != nil {
					return fmt.Errorf("line %d: failed to write buffer: %s", node.Debug.Line, err.Error())
				}
				bufferSize = bufferSize + n
			}
			fmt.Printf("%.6X: %20s\n", node.Address(), node.Debug.Source)
			break
		case *graph.InstructionNode: // Instruction code gen
			node := iter.Node().(*graph.InstructionNode)
			if startAddr < 0 {
				startAddr = node.Address()
			}
			data, err := asm.instructionBytes(node)
			if err != nil {
				return fmt.Errorf("line %d: %s", node.Debug.Line, err.Error())
			}
			n, err := buffer.Write(data)
			if err != nil {
				return fmt.Errorf("line %d: failed to write buffer: %s", node.Debug.Line, err.Error())
			}
			bufferSize = bufferSize + n
			fmt.Printf("%.6X: %20s    ->   %10s\n", node.Address(), node.Debug.Source, strings.ToUpper(hex.EncodeToString(data)))
			break
		default:
			return fmt.Errorf("invalid node '#%d' at 0x%X", iter.Index(), iter.Address())
		}
	}

	writeTextRecord()

	return nil
}
