package main

import (
	"dpkgraph/pkg/api"
	"dpkgraph/pkg/graph"
	"dpkgraph/pkg/storage"
	"log"
)

func main() {
	boltStorage, err := storage.NewBoltStorage("dpkgraph.db")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer boltStorage.Close()

	g := graph.NewGraph(boltStorage)

	if err := g.LoadNodes(); err != nil {
		log.Fatalf("Failed to load nodes: %v", err)
	}

	if err := g.LoadEdges(); err != nil {
		log.Fatalf("Failed to load edges: %v", err)
	}

	server := api.NewServer(g, boltStorage)
	if err := server.Start("8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
