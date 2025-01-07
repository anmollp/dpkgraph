package graph

import "fmt"

type Node struct {
	ID         string
	Label      string
	Properties map[string]interface{}
}

type Edge struct {
	From       string
	To         string
	Label      string
	Properties map[string]interface{}
}
type Graph struct {
	Nodes map[string]*Node
	Edges map[string][]*Edge
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string][]*Edge),
	}
}

func (g *Graph) AddNode(id, label string, properties map[string]interface{}) error {
	if _, exists := g.Nodes[id]; exists {
		return fmt.Errorf("node with ID %s already exists", id)
	}
	g.Nodes[id] = &Node{
		ID:         id,
		Label:      label,
		Properties: properties,
	}
	return nil
}

func (g *Graph) GetNode(id string) (*Node, error) {
	node, exists := g.Nodes[id]
	if !exists {
		return nil, fmt.Errorf("node with ID %s not found", id)
	}
	return node, nil
}

func (g *Graph) AddEdge(from, to, label string, properties map[string]interface{}) error {
	if _, exists := g.Nodes[from]; !exists {
		return fmt.Errorf("source node %s not found", from)
	}

	if _, exists := g.Nodes[to]; !exists {
		return fmt.Errorf("destination node %s not found", from)
	}

	g.Edges[from] = append(g.Edges[from], &Edge{
		From:       from,
		To:         to,
		Label:      label,
		Properties: properties,
	})
	return nil
}

func (g *Graph) GetEdgesFromNode(id string) ([]*Edge, error) {
	edges, exists := g.Edges[id]
	if !exists {
		return nil, fmt.Errorf("no edges found for node %s", id)
	}
	return edges, nil
}

func (g *Graph) DeleteNode(id string) error {
	if _, exists := g.Nodes[id]; !exists {
		return fmt.Errorf("node with ID %s not found", id)
	}
	delete(g.Nodes, id)
	delete(g.Edges, id)

	for from, edges := range g.Edges {
		var updatedEdges []*Edge
		for _, edge := range edges {
			if edge.From != id {
				updatedEdges = append(updatedEdges, edge)
			}
		}
		g.Edges[from] = updatedEdges
	}
	return nil
}
