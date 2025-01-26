package graph

import (
	"container/heap"
	"fmt"
	"math"
	"slices"
)

func (g *Graph) FindShortestPathWeighted(from, to string) ([]string, error) {
	if from == to {
		return []string{from}, nil
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.Nodes[from]; !exists {
		return nil, fmt.Errorf("source node %s not found", from)
	}
	if _, exists := g.Nodes[to]; !exists {
		return nil, fmt.Errorf("target node %s not found", to)
	}

	dist := make(map[string]float64)
	prev := make(map[string]string)
	for node := range g.Nodes {
		dist[node] = math.Inf(1)
	}
	dist[from] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{Value: from, Priority: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item).Value
		if current == to {
			break
		}
		for _, edge := range g.Edges[current] {
			neighbor := edge.To
			if edge.Weight < 0 {
				return nil, fmt.Errorf("negative weight edge detected from %s to %s", current, neighbor)
			}
			alt := dist[current] + edge.Weight
			if alt < dist[neighbor] {
				dist[neighbor] = alt
				prev[neighbor] = current
				heap.Push(pq, &Item{Value: neighbor, Priority: alt})
			}
		}
	}

	if dist[to] == math.Inf(1) {
		return nil, fmt.Errorf("no path found from %s to %s", from, to)
	}

	var path []string
	for node := to; node != ""; node = prev[node] {
		path = append(path, node)
	}

	slices.Reverse(path)
	return path, nil
}
