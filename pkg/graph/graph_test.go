package graph

import (
	"testing"
)

func TestGraph(t *testing.T) {
	g := NewGraph()

	err := g.AddNode("1", "Person", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	err = g.AddNode("2", "Person", map[string]interface{}{"name": "Bob"})
	if err != nil {
		t.Fatal(err)
	}

	err = g.AddEdge("1", "2", "knows", nil)
	if err != nil {
		t.Fatal(err)
	}

	node, err := g.GetNode("1")
	if err != nil || node.Label != "Person" {
		t.Fatal("failed to retrieve node")
	}

	edges, err := g.GetEdgesFromNode("1")
	if err != nil || len(edges) != 1 {
		t.Fatal("failed to retrieve edges")
	}
}
