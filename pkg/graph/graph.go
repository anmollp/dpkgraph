package graph

import (
	"dpkgraph/pkg/storage_interface"
	"fmt"
	"sync"
)

type Graph struct {
	Nodes   map[string]*storage_interface.Node
	Edges   map[string][]*storage_interface.Edge
	Storage storage_interface.Storage
	mu      sync.RWMutex
}

func NewGraph(storage storage_interface.Storage) *Graph {
	return &Graph{
		Nodes:   make(map[string]*storage_interface.Node),
		Edges:   make(map[string][]*storage_interface.Edge),
		Storage: storage,
	}
}

func (g *Graph) AddNode(id, label string, properties map[string]interface{}) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.Nodes[id]; exists {
		return fmt.Errorf("node with ID %s already exists", id)
	}
	node := &storage_interface.Node{
		ID:         id,
		Label:      label,
		Properties: properties,
	}
	g.Nodes[id] = node
	if g.Storage != nil {
		if err := g.Storage.SaveNode(node); err != nil {
			return fmt.Errorf("failed to save node: %w", err)
		}
	}
	return nil
}

func (g *Graph) GetNode(id string) (*storage_interface.Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	node, exists := g.Nodes[id]
	if !exists {
		return nil, fmt.Errorf("node with ID %s not found", id)
	}
	return node, nil
}

func (g *Graph) DeleteNode(nodeId string) error {
	node, exists := g.Nodes[nodeId]
	if !exists {
		return fmt.Errorf("node with ID %s not found", nodeId)
	}

	inDegree, err := g.InDegree(nodeId)
	if err != nil {
		return err
	}

	if inDegree > 0 {
		return fmt.Errorf("node %s has incoming edges and cannot be deleted", nodeId)
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	outgoingEdges := g.Edges[nodeId]

	if err := g.Storage.DeleteEdges(outgoingEdges); err != nil {
		return fmt.Errorf("failed to delete outgoing edges for node %s: %w", nodeId, err)
	}
	delete(g.Edges, nodeId)

	if err := g.Storage.DeleteNodes([]*storage_interface.Node{node}); err != nil {
		return fmt.Errorf("failed to delete node %s from storage: %w", nodeId, err)
	}
	delete(g.Nodes, nodeId)
	return nil
}

func (g *Graph) LoadNodes() error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.Storage == nil {
		return nil
	}
	nodes, err := g.Storage.LoadNodes()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}
	for _, node := range nodes {
		g.Nodes[node.ID] = node
	}
	return nil
}

func (g *Graph) AddEdge(from, to, label string, weight float64, properties map[string]interface{}) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.Nodes[from]; !exists {
		return fmt.Errorf("source node %s not found", from)
	}

	if _, exists := g.Nodes[to]; !exists {
		return fmt.Errorf("destination node %s not found", to)
	}

	edges, _ := g.Edges[from]
	for _, edge := range edges {
		if edge.From == from && edge.To == to && edge.Label == edge.Label {
			return fmt.Errorf("edge from %v to %v with label %v already exists", from, to, label)
		}
	}

	edge := &storage_interface.Edge{
		From:       from,
		To:         to,
		Label:      label,
		Weight:     weight,
		Properties: properties,
	}

	g.Edges[from] = append(g.Edges[from], edge)
	if g.Storage != nil {
		if err := g.Storage.SaveEdge(edge); err != nil {
			return fmt.Errorf("failed to save edge: %w", err)
		}
	}
	return nil
}

func (g *Graph) GetEdge(from, to, label string) ([]*storage_interface.Edge, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.Nodes[from]; !exists {
		return nil, fmt.Errorf("source node %s not found", from)
	}

	var matchedEdges []*storage_interface.Edge
	for _, edge := range g.Edges[from] {
		toMatch := to == "" || edge.To == to
		labelMatch := label == "" || edge.Label == label
		if toMatch && labelMatch {
			matchedEdges = append(matchedEdges, edge)
		}
	}
	return matchedEdges, nil
}

func (g *Graph) DeleteEdge(from, to, label string) error {
	if _, exists := g.Nodes[from]; !exists {
		return fmt.Errorf("source node %s not found", from)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	var matchedEdges []*storage_interface.Edge
	for _, edge := range g.Edges[from] {
		toMatch := to == "" || edge.To == to
		labelMatch := label == "" || edge.Label == label
		if toMatch && labelMatch {
			matchedEdges = append(matchedEdges, edge)

		}
	}
	if err := g.Storage.DeleteEdges(matchedEdges); err != nil {
		return err
	}
	for _, edge := range matchedEdges {
		delete(g.Edges, edge.From)
	}
	return nil
}

func (g *Graph) LoadEdges() error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.Storage == nil {
		return nil
	}
	edges, err := g.Storage.LoadEdges()
	if err != nil {
		return fmt.Errorf("failed to load edges: %w", err)
	}
	for _, edge := range edges {
		g.Edges[edge.From] = append(g.Edges[edge.From], edge)
	}
	return nil
}

func (g *Graph) FindNodesByProperties(properties map[string][]interface{}) ([]*storage_interface.Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var matchedNodes []*storage_interface.Node
	for _, node := range g.Nodes {
		if matchProperties(node.Properties, properties) {
			matchedNodes = append(matchedNodes, node)
		}
	}
	return matchedNodes, nil
}

func matchProperties(nodeProps map[string]interface{}, searchProps map[string][]interface{}) bool {
	for key, values := range searchProps {
		nodeValue, exists := nodeProps[key]
		if !exists {
			return false
		}

		matchFound := false
		for _, value := range values {
			if nodeValue == value {
				matchFound = true
				break
			}
		}

		if !matchFound {
			return false
		}
	}
	return true
}

func (g *Graph) InDegree(nodeId string) (int, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.Nodes[nodeId]; !exists {
		return 0, fmt.Errorf("node %s not found", nodeId)
	}

	count := 0
	for _, edges := range g.Edges {
		for _, edge := range edges {
			if edge.To == nodeId {
				count++
			}
		}
	}
	return count, nil
}

func (g *Graph) OutDegree(nodeId string) (int, error) {
	g.mu.RLock()
	g.mu.RUnlock()

	if _, exists := g.Nodes[nodeId]; !exists {
		return 0, fmt.Errorf("node %s not found", nodeId)
	}
	return len(g.Edges[nodeId]), nil
}
