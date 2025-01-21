package api

import (
	"dpkgraph/pkg/graph"
	"dpkgraph/pkg/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	graph   *graph.Graph
	storage *storage.BoltStorage
}

func NewServer(graph *graph.Graph, storage *storage.BoltStorage) *Server {
	return &Server{
		graph:   graph,
		storage: storage,
	}
}

func (s *Server) Start(port string) error {
	router := gin.Default()

	router.POST("/nodes", s.AddNode)
	router.GET("/nodes/:id", s.GetNode)
	router.DELETE("/nodes/:id", s.DeleteNode)

	router.POST("/edges", s.AddEdge)
	router.GET("/edges", s.GetEdgesByQuery)
	router.DELETE("/edges", s.DeleteEdge)

	return router.Run(":" + port)
}
