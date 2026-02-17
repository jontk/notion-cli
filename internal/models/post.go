package models

type Post struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content,omitempty"`
	Status      string   `json:"status"`
	Platforms   []string `json:"platforms,omitempty"`
	PublishDate string   `json:"publish_date,omitempty"`
	URL         string   `json:"url"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type PostInput struct {
	Title       string   `json:"title"`
	Content     string   `json:"content,omitempty"`
	Status      string   `json:"status,omitempty"`
	Platforms   []string `json:"platforms,omitempty"`
	PublishDate string   `json:"publish_date,omitempty"`
}
