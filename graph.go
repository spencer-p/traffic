package traffic

import (
	"errors"
	"log"
	"math"
	"sync"
)

// Graph usage or internal errors
var (
	ErrExistingNode  = errors.New("traffic: node already exists")
	ErrMissingNode   = errors.New("traffic: node does not exist")
	ErrDisconnected  = errors.New("traffic: no connection exists")
	ErrMissingSearch = errors.New("traffic: no path yet found, haven't searched")
)

type Graph struct {
	nodes     map[string]Node
	trees     map[string]SpanningTree
	treeMutex sync.Mutex
}

type SpanningTree struct {
	visited   map[string]bool
	tree      map[string]Node
	distances map[string]float64
	edgeTree  map[string]Edge
}

type Step struct {
	node Node
	edge Edge
}

func NewGraph() *Graph {
	return &Graph{nodes: make(map[string]Node), trees: make(map[string]SpanningTree)}
}

func (g *Graph) Dijkstra(start string) (map[string]Node, error) {
	// First check if the node exists
	if g.nodes[start] == nil {
		return nil, ErrMissingNode
	}

	span := SpanningTree{}

	log.Println("Finding spanning tree for", start)

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

	for len(span.visited) != len(g.nodes) {
		log.Printf("Visited %d of %d\n", len(span.visited), len(g.nodes))

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

		log.Printf("Currently looking at node %v\n", current)

		// Mark this as visited
		span.visited[currentName] = true

		// Update all the connected nodes
		for _, edge := range current.Edges() {
			// If distance[current] + edge.Weight() < distance[edge.To()]
			if span.distances[currentName]+edge.Weight() < span.distances[edge.To()] {
				log.Printf("Updating adjacent node %s\n", edge.To())
				span.distances[edge.To()] = span.distances[currentName] + edge.Weight()
				span.tree[edge.To()] = current
				span.edgeTree[edge.To()] = edge
			}
		}
	}

	g.treeMutex.Lock()
	g.trees[start] = span
	g.treeMutex.Unlock()
	return span.tree, nil
}

func (g *Graph) Path(start, destination string) ([]Step, error) {
	// Check nodes are in the graph
	if g.nodes[start] == nil || g.nodes[destination] == nil {
		return nil, ErrMissingNode
	}

	g.treeMutex.Lock()
	if _, ok := g.trees[start]; !ok {
		return nil, ErrMissingSearch
	}

	tree := g.trees[start].tree
	edgeTree := g.trees[start].edgeTree
	g.treeMutex.Unlock()

	// Fail if no path
	if tree[destination] == nil {
		return nil, ErrDisconnected
	}

	// Write out the path
	path := make([]Step, 1)
	path[0] = Step{node: g.nodes[destination], edge: edgeTree[g.nodes[destination].Name()]}
	for walk := tree[destination]; walk != g.nodes[start]; walk = tree[walk.Name()] {
		//log.Println("Walking at", walk.Name())
		path = append(path, Step{node: walk, edge: edgeTree[walk.Name()]})
	}

	return path, nil
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
