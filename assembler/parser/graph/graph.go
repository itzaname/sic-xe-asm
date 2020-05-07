package graph

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/machine"
	"fmt"
	"math/rand"
	"strconv"
)

type Graph struct {
	Nodes    []Node
	SymTable map[string]Node
}

func New() Graph {
	return Graph{
		Nodes:    []Node{},
		SymTable: map[string]Node{},
	}
}

func (graph *Graph) Iterator() Iterator {
	return Iterator{graph: graph}
}

func (graph *Graph) Append(node Node) int {
	graph.Nodes = append(graph.Nodes, node)

	if node.Label() != "" {
		graph.SymTable[node.Label()] = graph.Nodes[len(graph.Nodes)-1]
	}
	return len(graph.Nodes)
}

func (graph *Graph) Insert(node Node, i int) error {
	if i > len(graph.Nodes) {
		return fmt.Errorf("out of range: max index %d", len(graph.Nodes))
	}

	if node.Label() != "" {
		graph.SymTable[node.Label()] = graph.Nodes[len(graph.Nodes)-1]
	}

	/*tmp := append([]Node{}, graph.Nodes[i:]...)
	graph.Nodes = append(graph.Nodes[0:i], node)
	graph.Nodes = append(graph.Nodes, tmp...)*/

	graph.Nodes = append(graph.Nodes[:i], append([]Node{node}, graph.Nodes[i:]...)...)

	return nil
}

func (graph *Graph) UpdateAddr() error {
	addr := 0
	// Handle start address
	if node, ok := graph.Nodes[0].(*DirectiveNode); ok {
		if node.Directive.Name == "START" {
			base, err := strconv.Atoi(node.Data.(string))
			if err != nil {
				return fmt.Errorf("could not read start address: %s", err.Error())
			}
			addr = base
		}
	}

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

	return nil
}

func (graph *Graph) UpdateSymtable() int {
	counter := 0

	for i := 0; i < len(graph.Nodes); i++ {
		if graph.Nodes[i].Label() != "" {
			graph.SymTable[graph.Nodes[i].Label()] = graph.Nodes[i]
			counter++
		}
	}

	return counter
}

func (graph *Graph) ResolveLiterals() (int, error) {
	counter := 0
	insertIndex := 0
	insertLine := 0
	tmpList := []DirectiveNode{}
	// Find the LTORG
	for i := 0; i < len(graph.Nodes); i++ {
		if val, ok := graph.Nodes[i].(*DirectiveNode); ok {
			if val.Directive.Name == "LTORG" {
				insertIndex = i + 1
				insertLine = val.Debug.Line
				break
			}
		}
	}

	for i := 0; i < len(graph.Nodes); i++ {
		if val, ok := graph.Nodes[i].(*InstructionNode); ok {
			for x := 0; x < len(val.Operands); x++ {
				if val.Operands[x].Type == 3 {
					if data, ok := val.Operands[x].Data.(string); ok {
						// At this point we have an instruction with a literal that we will pull out and remap elsewhere
						name := strconv.Itoa(rand.Int())
						node := DirectiveNode{
							Name:      name,
							Directive: &machine.Directive{Name: "BYTE", Storage: true},
							Data: &Storage{
								Type: 1,
								Size: len(data),
								Data: []byte(data),
							},
							Debug: val.Debug,
							Graph: val.Graph,
						}
						node.Debug.Source = "LITERAL -> " + data
						tmpList = append(tmpList, node)
						val.Operands[x] = Operand{
							Type:       1,
							Addressing: 0,
							Data:       name,
						}
						graph.Nodes[i] = val
						counter++
					} else {
						return counter, fmt.Errorf("line %d: resolver: bad operand '%s'", val.Debug.Line, val.Operands[x].Type)
					}
				}
			}
		}
	}

	// Insert the nodes
	for i := 0; i < len(tmpList); i++ {
		if insertIndex > 0 && tmpList[i].Debug.Line < insertLine {
			if err := graph.Insert(&tmpList[i], insertIndex); err != nil {
				return counter, err
			}
		} else {
			// Insert before end directive
			if node, ok := graph.Nodes[len(graph.Nodes)-1].(*DirectiveNode); ok {
				if node.Directive.Name == "END" {
					if err := graph.Insert(&tmpList[i], len(graph.Nodes)-1); err != nil {
						return counter, err
					}
				}
			} else {
				graph.Append(&tmpList[i])
			}
		}
	}

	// We just destroyed the symbol table pointers
	graph.UpdateSymtable()
	return counter, nil
}

func (graph *Graph) LinkNodes() (int, error) {
	counter := 0

	for i := 0; i < len(graph.Nodes); i++ {
		if val, ok := graph.Nodes[i].(*DirectiveNode); ok {
			if val.Directive.Resolved {
				if node, ok := graph.SymTable[val.Data.(string)]; ok {
					val.Data = node
					graph.Nodes[i] = val
					counter++
				} else {
					return counter, fmt.Errorf("unresolved symbol '%s'", val.Data.(string))
				}
			}
		}
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
