package traffic

import (
	"fmt"
	"testing"
)

func ExampleMinQueue() {
	// New queue
	q := newMinQueue(5)

	// Initialize nodes named 0-4
	var nodes [5]Node
	for i := 0; i < len(nodes); i++ {
		nodes[i] = &node{name: fmt.Sprint(i)}
		q.push(nodes[i])
	}

	// Give em some priorities
	q.update(nodes[0], 3.14)
	q.update(nodes[1], 0.2)
	q.update(nodes[2], 1)
	q.update(nodes[3], 0.1)

	// Modify a priority
	q.update(nodes[0], 0.5)

	for q.Len() > 0 {
		n := q.pop()
		_ = n
		fmt.Println(n.Name())
	}

	// Output:
	// 3
	// 1
	// 0
	// 2
	// 4
}

func TestMinQueueSetup(t *testing.T) {
	// New queue
	q := newMinQueue(5)

	t.Run("heap length is zero", func(t *testing.T) {
		if len(q.heap) != 0 {
			t.Fail()
		}
	})
	t.Run("heap capacity is five", func(t *testing.T) {
		if cap(q.heap) != 5 {
			t.Fail()
		}
	})
	t.Run("lookup length is zero", func(t *testing.T) {
		if len(q.lookup) != 0 {
			t.Fail()
		}
	})
}
