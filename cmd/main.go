package main

import (
	"dpkgraph/pkg/api"
	"dpkgraph/pkg/graph"
	"dpkgraph/pkg/storage"
	"log"
)

func main() {
	g := graph.NewGraph()
	s, err := storage.NewStorage("graph.db")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer s.Close()

	loadedGraph, err := s.LoadGraph()
	if loadedGraph != nil && err != nil {
		g = loadedGraph
	} else {
		log.Println("Starting with an empty graph")
	}

	server := api.NewServer(g, s)
	log.Fatal(server.Start("8080"))
}
