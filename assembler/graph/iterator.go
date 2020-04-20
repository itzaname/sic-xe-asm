package graph

type Iterator struct {
	index   int
	address int
	graph   *Graph
}

func (itr *Iterator) Address() int {
	return itr.address
}

func (itr *Iterator) Index() int {
	return itr.index
}

func (itr *Iterator) Node() *Node {
	return &itr.graph.Nodes[itr.index]
}

func (itr *Iterator) Next() bool {
	itr.index += 1
	if itr.index-1 >= len(itr.graph.Nodes) {
		return false
	}

	if itr.index-1 > 0 {
		itr.address += itr.graph.Nodes[itr.index-1].Size()
	}

	return true
}
