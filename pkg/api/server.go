package api

import (
	"dpkgraph/pkg/graph"
	"dpkgraph/pkg/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	graph   *graph.Graph
	storage *storage.Storage
}

func NewServer(graph *graph.Graph, storage *storage.Storage) *Server {
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
	router.GET("/edges/:id", s.GetEdgesFromNode)

	return router.Run(":" + port)
}
