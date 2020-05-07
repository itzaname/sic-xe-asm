package assembler

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
)

type Assembler struct {
	obj   ObjectFile
	flags asmFlags
	graph *graph.Graph
}

type ObjectFile struct {
	Header        machine.HeaderRecord
	End           machine.EndRecord
	Modifications []machine.ModificationRecord
	Text          []machine.TextRecord
}

type asmFlags struct {
	base      bool
	baseAddr  int
	endModule bool
}

func GetObject(file string) (string, error) {
	p, err := parser.New(file)
	if err != nil {
		return "", err
	}

	asm := Assembler{
		graph: p.Graph(),
	}

	if err := asm.generateObjectItems(); err != nil {
		return "", err
	}

	asm.obj.Header.Length = asm.obj.Text[len(asm.obj.Text)-1].Length + asm.obj.Text[len(asm.obj.Text)-1].Start
	buffer := ""
	buffer = buffer + asm.obj.Header.String() + "\n"
	for i := 0; i < len(asm.obj.Text); i++ {
		buffer = buffer + asm.obj.Text[i].String() + "\n"
	}
	for i := 0; i < len(asm.obj.Modifications); i++ {
		buffer = buffer + asm.obj.Modifications[i].String() + "\n"
	}
	buffer = buffer + asm.obj.End.String() + "\n"

	return buffer, nil
}
