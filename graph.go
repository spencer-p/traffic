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

type graph struct {
	nodes map[string]Node
}

func NewGraph() *graph {
	return &graph{nodes: make(map[string]Node)}
}

func (g *graph) Dijkstra(start, destination string, agent Agent) (*SpanningTree, error) {
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

	// And a priority queue for the next node to visit
	priorityQueue := newMinQueue(len(g.nodes))
	for i := range g.nodes {
		priorityQueue.push(g.nodes[i])
	}
	priorityQueue.update(g.nodes[start], 0)

	for priorityQueue.Len() > 0 && !span.visited[destination] {

		// Find the current node to update around
		current := priorityQueue.pop()

		// Mark this as visited
		span.visited[current.Name()] = true

		// Update all the connected nodes
		for _, edge := range current.Edges() {
			// If distance[current] + edge.Weight() < distance[edge.To()]
			if span.distances[current.Name()]+edge.Weight(agent) < span.distances[edge.To()] {
				span.distances[edge.To()] = span.distances[current.Name()] + edge.Weight(agent)
				span.tree[edge.To()] = current
				span.edgeTree[edge.To()] = edge

				// Update in priority queue
				priorityQueue.update(g.nodes[edge.To()], span.distances[edge.To()])
			}
		}
	}

	return &span, nil
}

func (g *graph) AddEdge(toadd Edge) {
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
