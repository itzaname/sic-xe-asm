package main

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/graph"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"fmt"
)

func main() {
	g := graph.Graph{Nodes: []graph.Node{
		graph.Node(&graph.InstructionNode{Instruction: &machine.Instruction{Format: 3}}),
		graph.Node(&graph.InstructionNode{Instruction: &machine.Instruction{Format: 3}}),
		graph.Node(&graph.InstructionNode{Instruction: &machine.Instruction{Format: 3}}),
		graph.Node(&graph.InstructionNode{Instruction: &machine.Instruction{Format: 3}}),
	}}

	itr := g.Iterator()
	for itr.Next() {
		fmt.Println(itr.Index(), itr.Address(), itr.Node().Get())
	}
}
