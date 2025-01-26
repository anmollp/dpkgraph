package storage_interface

type Storage interface {
	SaveNode(node *Node) error
	SaveEdge(edge *Edge) error
	LoadNodes() ([]*Node, error)
	LoadEdges() ([]*Edge, error)
	DeleteNodes([]*Node) error
	DeleteEdges([]*Edge) error
	Close() error
}
