package traffic

import (
	"errors"
	"math"
)

// Graph usage or internal errors
var (
	ErrExistingNode  = errors.New("traffic: node already exists")
	ErrMissingNode   = errors.New("traffic: node does not exist")
	ErrDisconnected  = errors.New("traffic: no connection exists")
	ErrMissingSearch = errors.New("traffic: no path yet found, haven't searched")
)

type Graph struct {
	nodes map[string]Node
}

func NewGraph() *Graph {
	return &Graph{nodes: make(map[string]Node)}
}

func (g *Graph) Dijkstra(start, destination string, agent Agent) (*SpanningTree, error) {
	// First check if the nodes exist
	if g.nodes[start] == nil || g.nodes[destination] == nil {
		return nil, ErrMissingNode
	}

	span := SpanningTree{start: g.nodes[start], destination: g.nodes[destination]}

	// The map of already visited nodes
	span.visited = make(map[string]bool)

	// The final parent graph that will be generated
	span.tree = make(map[string]Node)

	// Parent edges for pathing
	span.edgeTree = make(map[string]Edge)

	// Distance values
	span.distances = make(map[string]float64)
	for name := range g.nodes {
		// All non start values are an infinite distance
		if name != start {
			span.distances[name] = math.Inf(0)
		}
	}

	for len(span.visited) != len(g.nodes) && !span.visited[destination] {

		// Find the current node to update around
		var currentName string
		if len(span.visited) == 0 {
			currentName = start
		} else {
			min := math.Inf(0)
			for name, d := range span.distances {
				if span.visited[name] == false && d <= min {
					currentName = name
					min = d
				}
			}
		}
		current := g.nodes[currentName]

		// Mark this as visited
		span.visited[currentName] = true

		// Update all the connected nodes
		for _, edge := range current.Edges() {
			// If distance[current] + edge.Weight() < distance[edge.To()]
			if span.distances[currentName]+edge.Weight(agent) < span.distances[edge.To()] {
				span.distances[edge.To()] = span.distances[currentName] + edge.Weight(agent)
				span.tree[edge.To()] = current
				span.edgeTree[edge.To()] = edge
			}
		}
	}

	return &span, nil
}

func (g *Graph) AddNode(toadd Node) error {
	if existing := g.nodes[toadd.Name()]; existing != nil && existing != toadd {
		return ErrExistingNode
	}
	g.nodes[toadd.Name()] = toadd
	return nil
}

func (g *Graph) AddEdge(toadd Edge) {
	//log.Printf("Adding edge %v\n", toadd)

	// Add the edge to the From() node
	if existing := g.nodes[toadd.From()]; existing == nil {
		// If the node doesn't exist, make it
		edges := make([]Edge, 1)
		edges[0] = toadd
		g.nodes[toadd.From()] = &node{name: toadd.From(), edges: edges}
	} else {
		// Otherwise simple append
		existing.AddEdge(toadd)
	}

	// Make sure the destination node exists
	if existing := g.nodes[toadd.To()]; existing == nil {
		g.nodes[toadd.To()] = &node{name: toadd.To(), edges: make([]Edge, 0)}
	}
}
