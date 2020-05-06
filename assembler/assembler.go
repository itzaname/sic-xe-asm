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
	base bool
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

	return "", nil
}
