package storage_interface

import "fmt"

type Edge struct {
	From       string
	To         string
	Label      string
	Properties map[string]interface{}
}

func (e Edge) GetKey() string {
	return fmt.Sprintf("%s->%s:%s", e.From, e.To, e.Label)
}
