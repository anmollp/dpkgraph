package api

type NodeRequest struct {
	ID         string                 `json:"id"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties"`
}

type EdgeRequest struct {
	From       string                 `json:"from"`
	To         string                 `json:"to"`
	Label      string                 `json:"label"`
	Weight     float64                `json:"weight"`
	Properties map[string]interface{} `json:"properties"`
}
