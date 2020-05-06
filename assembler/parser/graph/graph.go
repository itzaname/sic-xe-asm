package graph

import "fmt"

type Graph struct {
	Nodes    []Node
	SymTable map[string]*Node
}

func New() Graph {
	return Graph{
		Nodes:    []Node{},
		SymTable: map[string]*Node{},
	}
}

func (graph *Graph) Iterator() Iterator {
	return Iterator{graph: graph}
}

func (graph *Graph) Append(node Node) int {
	graph.Nodes = append(graph.Nodes, node)

	if node.Label() != "" {
		graph.SymTable[node.Label()] = &graph.Nodes[len(graph.Nodes)-1]
	}
	return len(graph.Nodes)
}

func (graph *Graph) Insert(node Node, i int) error {
	if i > len(graph.Nodes) {
		return fmt.Errorf("out of range: max index %d", len(graph.Nodes))
	}

	tmp := append([]Node{}, graph.Nodes[i:]...)
	graph.Nodes = append(graph.Nodes[0:i], node)
	graph.Nodes = append(graph.Nodes, tmp...)

	return nil
}

func (graph *Graph) UpdateAddr() {
	addr := 0
	for i := 0; i < len(graph.Nodes); i++ {
		if node, ok := graph.Nodes[i].(*InstructionNode); ok {
			node.Addr = addr
			graph.Nodes[i] = node
			addr = addr + graph.Nodes[i].Size()
			continue
		}
		if node, ok := graph.Nodes[i].(*DirectiveNode); ok {
			node.Addr = addr
			graph.Nodes[i] = node
			addr = addr + graph.Nodes[i].Size()
			continue
		}
	}
}

func (graph *Graph) LinkNodes() (int, error) {
	counter := 0

	for i := 0; i < len(graph.Nodes); i++ {
		if val, ok := graph.Nodes[i].(*InstructionNode); ok {
			for x := 0; x < len(val.Operands); x++ {
				if val.Operands[x].Type == 1 {
					if node, ok := graph.SymTable[val.Operands[x].Data.(string)]; ok {
						val.Operands[x].Data = node
						graph.Nodes[i] = val
						counter++
					} else {
						return counter, fmt.Errorf("unresolved symbol '%s'", val.Operands[x].Data.(string))
					}
				}
			}
		}
	}

	return counter, nil
}
