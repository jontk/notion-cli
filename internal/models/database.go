package models

type DatabaseInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type PropertyInfo struct {
	Type    string         `json:"type"`
	Options map[string]any `json:"options,omitempty"`
}

type Schema struct {
	Properties map[string]PropertyInfo `json:"properties"`
}
