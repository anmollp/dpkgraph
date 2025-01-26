package graph

import (
	"testing"
)

func TestFindShortestPath(t *testing.T) {
	testStorage := setupTestStorage(t)
	g := NewGraph(testStorage)

	// Add nodes
	g.AddNode("1", "Person", nil)
	g.AddNode("2", "Person", nil)
	g.AddNode("3", "Person", nil)
	g.AddNode("4", "Person", nil)

	// Add edges
	g.AddEdge("1", "2", "", 1, map[string]interface{}{})
	g.AddEdge("2", "3", "", 2, map[string]interface{}{})
	g.AddEdge("1", "3", "", 4, map[string]interface{}{})
	g.AddEdge("3", "4", "", 1, map[string]interface{}{})

	// Test cases
	tests := []struct {
		name      string
		from      string
		to        string
		wantPath  []string
		expectErr bool
	}{
		{
			name:     "Direct Path",
			from:     "1",
			to:       "2",
			wantPath: []string{"1", "2"},
		},
		{
			name:     "Shortest Path via Intermediate Node",
			from:     "1",
			to:       "3",
			wantPath: []string{"1", "2", "3"},
		},
		{
			name:     "Path to Another Node",
			from:     "1",
			to:       "4",
			wantPath: []string{"1", "2", "3", "4"},
		},
		{
			name:      "Unreachable Node",
			from:      "1",
			to:        "5", // Node does not exist
			expectErr: true,
		},
		{
			name:      "Nonexistent Start Node",
			from:      "0", // Node does not exist
			to:        "3",
			expectErr: true,
		},
		{
			name:      "Nonexistent End Node",
			from:      "1",
			to:        "6", // Node does not exist
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := g.FindShortestPathWeighted(tt.from, tt.to)

			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}

			// Check the path if no error was expected
			if !tt.expectErr && !equalPaths(path, tt.wantPath) {
				t.Fatalf("expected path: %v, got: %v", tt.wantPath, path)
			}
		})
	}
}

// Helper function to compare two paths
func equalPaths(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
