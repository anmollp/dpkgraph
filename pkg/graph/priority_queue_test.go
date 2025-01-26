package graph

import (
	"container/heap"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := &PriorityQueue{}
	heap.Init(pq)

	heap.Push(pq, &Item{Value: "1", Priority: 2})
	heap.Push(pq, &Item{Value: "2", Priority: 1})

	item := heap.Pop(pq).(*Item)
	if item.Value != "2" {
		t.Fatalf("Expected '2', got '%s'", item.Value)
	}
}
