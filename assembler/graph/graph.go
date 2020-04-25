package graph

import "fmt"

type Graph struct {
	Nodes    []Node
	SymTable []map[string]Node
}

func New() Graph {
	return Graph{}
}

func (graph *Graph) Iterator() Iterator {
	return Iterator{graph: graph}
}

func (graph *Graph) Append(node Node) int {
	graph.Nodes = append(graph.Nodes, node)
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
