package graph

import (
	"dpkgraph/pkg/storage"
	"dpkgraph/pkg/storage_interface"
	"os"
	"testing"
)

func setupTestStorage(t *testing.T) storage_interface.Storage {
	dbPath := "test_graph.db"
	t.Cleanup(func() { os.RemoveAll(dbPath) })

	boltStorage, err := storage.NewBoltStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize BoltStorage: %v", err)
	}
	t.Cleanup(func() { boltStorage.Close() })

	return boltStorage
}

func TestAddNode(t *testing.T) {
	testStorage := setupTestStorage(t)
	g := NewGraph(testStorage)

	err := g.AddNode("1", "Person", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("failed to add node: %v", err)
	}

	node, err := g.GetNode("1")
	if err != nil {
		t.Fatalf("failed to retrieve node: %v", err)
	}
	if node.Label != "Person" || node.Properties["name"] != "Alice" {
		t.Fatalf("unexpected node data: %+v", node)
	}
}

func TestAddEdge(t *testing.T) {
	testStorage := setupTestStorage(t)
	g := NewGraph(testStorage)

	// Add nodes for the edge
	err := g.AddNode("1", "Person", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("failed to add source node: %v", err)
	}
	err = g.AddNode("2", "Person", map[string]interface{}{"name": "Bob"})
	if err != nil {
		t.Fatalf("failed to add target node: %v", err)
	}

	// Add edge
	err = g.AddEdge("1", "2", "knows", nil)
	if err != nil {
		t.Fatalf("failed to add edge: %v", err)
	}

	edges, err := g.SearchEdges("1->*:*")
	if err != nil {
		t.Fatalf("failed to retrieve edges: %v", err)
	}
	if len(edges) != 1 || edges[0].To != "2" || edges[0].Label != "knows" {
		t.Fatalf("unexpected edge data: %+v", edges[0])
	}
}

func TestAddDuplicateEdge(t *testing.T) {
	testStorage := setupTestStorage(t)
	g := NewGraph(testStorage)

	// Add nodes for the edge
	err := g.AddNode("1", "Person", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("failed to add source node: %v", err)
	}
	err = g.AddNode("2", "Person", map[string]interface{}{"name": "Bob"})
	if err != nil {
		t.Fatalf("failed to add target node: %v", err)
	}

	// Add edge
	err = g.AddEdge("1", "2", "knows", nil)
	if err != nil {
		t.Fatalf("failed to add edge: %v", err)
	}

	err = g.AddEdge("1", "2", "knows", nil)
	if err == nil {
		t.Fatalf("failed to check duplicate edge: %v", err)
	}
}

func TestDeleteNode(t *testing.T) {
	testStorage := setupTestStorage(t)
	g := NewGraph(testStorage)

	err := g.AddNode("1", "Person", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("failed to add source node: %v", err)
	}

	err = g.DeleteNode("1")
	if err != nil {
		t.Fatalf("failed to delete node: %v", err)
	}

	node, err := g.GetNode("1")
	if err == nil && node != nil {
		t.Fatalf("unexpected node data after reload: %+v", node)
	}
}

func TestDeleteEdge(t *testing.T) {
	testStorage := setupTestStorage(t)
	g := NewGraph(testStorage)

	err := g.AddNode("1", "Person", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("failed to add source node: %v", err)
	}
	err = g.AddNode("2", "Person", map[string]interface{}{"name": "Bob"})
	if err != nil {
		t.Fatalf("failed to add target node: %v", err)
	}

	err = g.AddEdge("1", "2", "knows", nil)
	if err != nil {
		t.Fatalf("failed to add edge: %v", err)
	}

	err = g.DeleteEdge("1", "2", "knows")
	if err != nil {
		t.Fatalf("failed to delete edge: %v", err)
	}

	edges, err := g.SearchEdges("1->*:*")
	if err == nil && edges != nil {
		t.Fatalf("unexpected edge data after reload: %+v", edges[0])
	}
}

func TestPersistence(t *testing.T) {
	testStorage := setupTestStorage(t)

	// Create a graph and add nodes and edges
	g := NewGraph(testStorage)
	err := g.AddNode("1", "Person", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Fatalf("failed to add node: %v", err)
	}
	err = g.AddNode("2", "Person", map[string]interface{}{"name": "Bob"})
	if err != nil {
		t.Fatalf("failed to add node: %v", err)
	}
	err = g.AddEdge("1", "2", "knows", nil)
	if err != nil {
		t.Fatalf("failed to add edge: %v", err)
	}

	// Restart graph to test persistence
	g = NewGraph(testStorage)
	err = g.LoadNodes()
	if err != nil {
		t.Fatalf("failed to load nodes: %v", err)
	}
	err = g.LoadEdges()
	if err != nil {
		t.Fatalf("failed to load edges: %v", err)
	}

	// Verify nodes
	node, err := g.GetNode("1")
	if err != nil {
		t.Fatalf("failed to retrieve node after reload: %v", err)
	}
	if node.Label != "Person" || node.Properties["name"] != "Alice" {
		t.Fatalf("unexpected node data after reload: %+v", node)
	}

	// Verify edges
	edges, err := g.SearchEdges("1->2:knows")
	if err != nil {
		t.Fatalf("failed to retrieve edges after reload: %v", err)
	}
	if len(edges) == 0 {
		t.Fatalf("No such edge exists")
	}
	if len(edges) != 1 || edges[0].To != "2" || edges[0].Label != "knows" {
		t.Fatalf("unexpected edge data after reload: %+v", edges[0])
	}

	err = g.DeleteEdge("1", "2", "knows")
	if err != nil {
		t.Fatalf("failed to delete edge: %v", err)
	}

	err = g.DeleteNode("1")
	if err != nil {
		t.Fatalf("failed to delete node: %v", err)
	}

	// Restart graph to test persistence
	g = NewGraph(testStorage)
	err = g.LoadNodes()
	if err != nil {
		t.Fatalf("failed to load nodes: %v", err)
	}
	err = g.LoadEdges()
	if err != nil {
		t.Fatalf("failed to load edges: %v", err)
	}

	// Verify nodes
	node, err = g.GetNode("1")
	if err == nil && node != nil {
		t.Fatalf("unexpected node data after reload: %+v", node)
	}

	// Verify edges
	edges, err = g.SearchEdges("1->*:*")
	if err == nil && edges != nil {
		t.Fatalf("unexpected edge data after reload: %+v", edges[0])
	}
}
