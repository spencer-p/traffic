package traffic

type Node interface {
	Name() string
	Edges() []Edge
	AddEdge(e Edge)
}

type node struct {
	name  string
	edges []Edge
}

func (n *node) Name() string {
	return n.name
}

func (n *node) Edges() []Edge {
	return n.edges
}

func (n *node) AddEdge(e Edge) {
	n.edges = append(n.edges, e)
}
