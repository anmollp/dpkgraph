package storage

import (
	"dpkgraph/pkg/graph"
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
)

type Storage struct {
	db *bbolt.DB
}

func NewStorage(filePath string) (*Storage, error) {
	db, err := bbolt.Open(filePath, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveGraph(g *graph.Graph) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("graph"))
		if err != nil {
			return err
		}
		data, err := json.Marshal(g)
		if err != nil {
			return err
		}
		return bucket.Put([]byte("graph_data"), data)
	})
}

func (s *Storage) LoadGraph() (*graph.Graph, error) {
	var g *graph.Graph
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("graph"))
		if bucket == nil {
			return fmt.Errorf("graph bucket not found")
		}
		data := bucket.Get([]byte("graph_data"))
		if data == nil {
			return fmt.Errorf("no graph data found")
		}
		return json.Unmarshal(data, &g)
	})
	return g, err
}

func (s *Storage) Close() error {
	return s.db.Close()
}
