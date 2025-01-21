package storage

import (
	"dpkgraph/pkg/storage_interface"
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
)

type BoltStorage struct {
	db *bbolt.DB
}

const (
	nodesBucket = "nodes"
	edgesBucket = "edges"
)

func NewBoltStorage(dbPath string) (*BoltStorage, error) {
	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(nodesBucket))
		if err != nil {
			return fmt.Errorf("create bucket %s: %w", nodesBucket, err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(edgesBucket))
		if err != nil {
			return fmt.Errorf("create bucket %s: %w", edgesBucket, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set up buckets: %w", err)
	}
	return &BoltStorage{db: db}, nil
}

func (b *BoltStorage) SaveNode(node *storage_interface.Node) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(nodesBucket))
		data, err := json.Marshal(node)
		if err != nil {
			return fmt.Errorf("failed to serialize node: %w", err)
		}
		return bucket.Put([]byte(node.ID), data)
	})
}

func (b *BoltStorage) LoadNodes() ([]*storage_interface.Node, error) {
	var nodes []*storage_interface.Node
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(nodesBucket))
		return bucket.ForEach(func(k, v []byte) error {
			var node storage_interface.Node
			if err := json.Unmarshal(v, &node); err != nil {
				return fmt.Errorf("failed to deserialize node: %w", err)
			}
			nodes = append(nodes, &node)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (b *BoltStorage) SaveEdge(edge *storage_interface.Edge) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(edgesBucket))
		data, err := json.Marshal(edge)
		if err != nil {
			return fmt.Errorf("failed to serialize edge: %w", err)
		}
		key := fmt.Sprintf("%s->%s:%s", edge.From, edge.To, edge.Label)
		return bucket.Put([]byte(key), data)
	})
}

func (b *BoltStorage) LoadEdges() ([]*storage_interface.Edge, error) {
	var edges []*storage_interface.Edge
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(edgesBucket))
		return bucket.ForEach(func(k, v []byte) error {
			var edge storage_interface.Edge
			if err := json.Unmarshal(v, &edge); err != nil {
				return fmt.Errorf("failed to unmarshal edge: %w", err)
			}
			edges = append(edges, &edge)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return edges, nil
}

func (b *BoltStorage) DeleteNode(id string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		nodesBucket := tx.Bucket([]byte(nodesBucket))
		if nodesBucket == nil {
			return fmt.Errorf("nodes bucket does not exist")
		}
		return nodesBucket.Delete([]byte(id))
	})
}

func (b *BoltStorage) DeleteEdge(from, to, label string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		edgesBucket := tx.Bucket([]byte(edgesBucket))
		if edgesBucket == nil {
			return fmt.Errorf("edges bucket does not exist")
		}
		edgeKey := fmt.Sprintf("%s->%s:%s", from, to, label)

		if edgesBucket.Get([]byte(edgeKey)) == nil {
			return fmt.Errorf("edge %s does not exist", edgeKey)
		}

		return edgesBucket.Delete([]byte(edgeKey))
	})
}

func (b *BoltStorage) Close() error {
	return b.db.Close()
}
