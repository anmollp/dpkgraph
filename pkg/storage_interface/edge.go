package storage_interface

type Edge struct {
	From       string
	To         string
	Label      string
	Weight     float64
	Properties map[string]interface{}
}
