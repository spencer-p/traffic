package traffic

import (
	"container/heap"
	"math"
)

// minQueueItem wraps a node with its priority and heap index.
type minQueueItem struct {
	value    Node
	priority float64
	index    int
}

// minQueue satisfies heap.Interface to provide a minimum priority queue.
type minQueue struct {
	heap   []*minQueueItem
	lookup map[string]int
}

func newMinQueue(length int) minQueue {
	return minQueue{
		heap:   make([]*minQueueItem, 0, length),
		lookup: make(map[string]int, length),
	}
}

func (m *minQueue) Len() int {
	return len(m.heap)
}

func (m *minQueue) Less(i, j int) bool {
	return m.heap[i].priority < m.heap[j].priority
}

func (m *minQueue) Swap(i, j int) {
	m.heap[i], m.heap[j] = m.heap[j], m.heap[i]
	m.heap[i].index = i
	m.heap[j].index = j
	m.lookup[m.heap[i].value.Name()] = i
	m.lookup[m.heap[j].value.Name()] = j
}

func (m *minQueue) Push(x interface{}) {
	// Recieved interface must be of type Node
	n := x.(Node)

	// Create the new element
	item := &minQueueItem{
		value:    n,
		priority: math.Inf(0),
		index:    len(m.heap),
	}

	// Update minQueue fields
	m.lookup[n.Name()] = item.index
	m.heap = append(m.heap, item)
}

func (m *minQueue) Pop() interface{} {
	n := len(m.heap)
	item := m.heap[n-1]
	item.index = -1 // just in case
	m.heap = m.heap[:n-1]
	delete(m.lookup, item.value.Name())
	return item.value
}

func (m *minQueue) update(x Node, priority float64) {
	item := m.heap[m.lookup[x.Name()]]
	item.priority = priority
	heap.Fix(m, item.index)
}

func (m *minQueue) push(x Node) {
	heap.Push(m, x)
}

func (m *minQueue) pop() Node {
	return heap.Pop(m).(Node)
}
