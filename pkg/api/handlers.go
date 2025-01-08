package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) AddNode(c *gin.Context) {
	var nodeData NodeRequest
	if err := c.ShouldBindJSON(&nodeData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.graph.AddNode(nodeData.ID, nodeData.Label, nodeData.Properties)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "node added"})
}

func (s *Server) GetNode(c *gin.Context) {
	id := c.Param("id")
	node, err := s.graph.GetNode(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, node)
}

func (s *Server) DeleteNode(c *gin.Context) {
	id := c.Param("id")
	err := s.graph.DeleteNode(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "node deleted"})
}

func (s *Server) AddEdge(c *gin.Context) {
	var edgeData EdgeRequest
	if err := c.ShouldBindJSON(&edgeData); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	err := s.graph.AddEdge(edgeData.From, edgeData.To, edgeData.Label, edgeData.Properties)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "edge added"})
}

func (s *Server) GetEdgesFromNode(c *gin.Context) {
	id := c.Param("id")
	edges, err := s.graph.GetEdgesFromNode(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, edges)
}
