package traffic

import (
	"errors"
	"log"
	"math"
)

// Graph usage or internal errors
var (
	ErrExistingNode = errors.New("traffic: node already exists")
	ErrMissingNode  = errors.New("traffic: node does not exist")
	ErrDisconnected = errors.New("traffic: no connection exists")
)

type Graph interface {
	Dijkstra(start string) (map[string]Node, error)
	Path(start, destination string) ([]string, error)
	AddNode(toadd Node) error
	AddEdge(toadd Edge)
}

type graph struct {
	nodes map[string]Node
}

func NewGraph() Graph {
	return &graph{nodes: make(map[string]Node)}
}

func (g *graph) Dijkstra(start string) (map[string]Node, error) {
	// First check if the node exists
	if g.nodes[start] == nil {
		return nil, ErrMissingNode
	}

	log.Println("Finding spanning tree for", start)

	// The map of already visited nodes
	visited := make(map[string]bool)

	// The final parent graph that will be generated
	tree := make(map[string]Node)

	// Distance values
	distances := make(map[string]float64)
	for name := range g.nodes {
		// All non start values are an infinite distance
		if name != start {
			distances[name] = math.Inf(0)
		}
	}

	for len(visited) != len(g.nodes) {
		log.Printf("Visited %d of %d\n", len(visited), len(g.nodes))

		// Find the current node to update around
		var currentName string
		if len(visited) == 0 {
			currentName = start
		} else {
			min := math.Inf(0)
			for name, d := range distances {
				if visited[name] == false && d <= min {
					currentName = name
					min = d
				}
			}
		}
		current := g.nodes[currentName]

		log.Printf("Currently looking at node %v\n", current)

		// Mark this as visited
		visited[currentName] = true

		// Update all the connected nodes
		for _, edge := range current.Edges() {
			// If distance[current] + edge.Weight() < distance[edge.To()]
			if distances[currentName]+edge.Weight() < distances[edge.To()] {
				log.Printf("Updating adjacent node %s\n", edge.To())
				distances[edge.To()] = distances[currentName] + edge.Weight()
				tree[edge.To()] = current
			}
		}
	}

	return tree, nil
}

func (g *graph) Path(start, destination string) ([]string, error) {
	// Check nodes are in the graph
	if g.nodes[start] == nil || g.nodes[destination] == nil {
		return nil, ErrMissingNode
	}

	// Find the spanning tree
	tree, err := g.Dijkstra(start)
	if err != nil {
		return nil, err
	}

	// Fail if no path
	if tree[destination] == nil {
		return nil, ErrDisconnected
	}

	// Write out the path
	path := make([]string, 1)
	path[0] = destination
	for walk := tree[destination]; walk != nil; walk = tree[walk.Name()] {
		log.Println("Walking at", walk.Name())
		path = append(path, walk.Name())
	}

	return path, nil
}

func (g *graph) AddNode(toadd Node) error {
	if existing := g.nodes[toadd.Name()]; existing != nil && existing != toadd {
		return ErrExistingNode
	}
	g.nodes[toadd.Name()] = toadd
	return nil
}

func (g *graph) AddEdge(toadd Edge) {
	log.Printf("Adding edge %v\n", toadd)

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
