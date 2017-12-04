package traffic

type SpanningTree struct {
	visited            map[string]bool
	tree               map[string]Node
	distances          map[string]float64
	edgeTree           map[string]Edge
	start, destination Node
}

type Step struct {
	node Node
	edge Edge
}

func (st *SpanningTree) Path() ([]Step, error) {
	// Check nodes are in the spanning tree
	if st.visited[st.start.Name()] == false || st.visited[st.destination.Name()] == false {
		return nil, ErrMissingNode
	}

	// Write out the path
	path := make([]Step, 1)
	path[0] = Step{node: st.destination, edge: st.edgeTree[st.destination.Name()]}
	for walk := st.tree[st.destination.Name()]; walk != st.start; walk = st.tree[walk.Name()] {
		//log.Println("Walking at", walk.Name())
		path = append(path, Step{node: walk, edge: st.edgeTree[walk.Name()]})
	}

	return path, nil
}
