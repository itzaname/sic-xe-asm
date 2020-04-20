package main

import (
	"ci.itzana.me/itzaname/sic-xe-asm/graph"
	"ci.itzana.me/itzaname/sic-xe-asm/machine"
	"fmt"
)

func main() {
	g := graph.Graph{Nodes: []graph.Node{
		{Instruction: &machine.Instruction{Format: 3}},
		{Instruction: &machine.Instruction{Format: 3}},
		{Instruction: &machine.Instruction{Format: 3}},
		{Instruction: &machine.Instruction{Format: 3}},
		{Instruction: &machine.Instruction{Format: 3}},
		{Instruction: &machine.Instruction{Format: 3}},
	}}

	itr := g.Iterator()
	for itr.Next() {
		fmt.Println(itr.Index(), itr.Address())
	}
}
