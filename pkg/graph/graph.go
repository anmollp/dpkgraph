package graph

import (
	"dpkgraph/pkg/storage_interface"
	"fmt"
	"regexp"
	"strings"
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

func (g *Graph) DeleteNode(id string) error {
	if _, exists := g.Nodes[id]; !exists {
		return fmt.Errorf("node with ID %s not found", id)
	}
	outEdgePattern := fmt.Sprintf("%s->*:*", id)
	inEdgePattern := fmt.Sprintf("*->%s:*", id)
	if err := g.RemoveEdges(outEdgePattern); err != nil {
		return err
	}
	if err := g.RemoveEdges(inEdgePattern); err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.Nodes, id)
	return g.Storage.DeleteNode(id)
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

func (g *Graph) AddEdge(from, to, label string, properties map[string]interface{}) error {
	edgePattern := fmt.Sprintf("%s->%s:%s", from, to, label)
	if edges, _ := g.SearchEdges(edgePattern); len(edges) != 0 {
		return fmt.Errorf("edge already exists")
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.Nodes[from]; !exists {
		return fmt.Errorf("source node %s not found", from)
	}

	if _, exists := g.Nodes[to]; !exists {
		return fmt.Errorf("destination node %s not found", to)
	}

	edge := &storage_interface.Edge{
		From:       from,
		To:         to,
		Label:      label,
		Properties: properties,
	}

	edgeKey := edge.GetKey()
	g.Edges[edgeKey] = append(g.Edges[edgeKey], edge)
	if g.Storage != nil {
		if err := g.Storage.SaveEdge(edge); err != nil {
			return fmt.Errorf("failed to save edge: %w", err)
		}
	}
	return nil
}

func (g *Graph) SearchEdges(pattern string) ([]*storage_interface.Edge, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var matchedEdges []*storage_interface.Edge
	regexPattern := "^" + strings.ReplaceAll(regexp.QuoteMeta(pattern), "\\*", ".*") + "$"
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}
	for key, edge := range g.Edges {
		if re.Match([]byte(key)) {
			matchedEdges = append(matchedEdges, edge...)
		}
	}
	return matchedEdges, nil
}

func (g *Graph) DeleteEdge(from, to, label string) error {
	edgePattern := fmt.Sprintf("%s->%s:%s", from, to, label)
	err := g.RemoveEdges(edgePattern)
	if err != nil {
		return err
	}
	return g.Storage.DeleteEdge(from, to, label)
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
		edgeKey := edge.GetKey()
		g.Edges[edgeKey] = append(g.Edges[edgeKey], edge)
	}
	return nil
}

func (g *Graph) RemoveEdges(pattern string) error {
	edgesToRemove, err := g.SearchEdges(pattern)
	if err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, edge := range edgesToRemove {
		delete(g.Edges, edge.GetKey())
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
