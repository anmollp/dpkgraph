package api

import (
	"fmt"
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

func (s *Server) GetEdgesByQuery(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	label := c.Query("label")

	if from == "" && to == "" && label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one query parameter (from, to, or label) is required"})
		return
	}
	if from == "" {
		from = "*"
	}
	if to == "" {
		to = "*"
	}
	if label == "" {
		label = "*"
	}

	pattern := fmt.Sprintf("%s->%s:%s", from, to, label)

	edges, err := s.graph.SearchEdges(pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if edges == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No edges found"})
		return
	}
	c.JSON(http.StatusOK, edges)
}

func (s *Server) DeleteEdge(c *gin.Context) {
	var edge EdgeRequest
	if err := c.ShouldBindJSON(&edge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := s.graph.DeleteEdge(edge.From, edge.To, edge.Label)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Edge deleted successfully"})
}
