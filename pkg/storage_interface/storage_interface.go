package storage_interface

type Storage interface {
	SaveNode(node *Node) error
	SaveEdge(edge *Edge) error
	LoadNodes() ([]*Node, error)
	LoadEdges() ([]*Edge, error)
	DeleteNode(id string) error
	DeleteEdge(sourceID, targetID, label string) error
	Close() error
}
