package models

type Task struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Status    string   `json:"status"`
	Priority  string   `json:"priority,omitempty"`
	DueDate   string   `json:"due_date,omitempty"`
	Category  string   `json:"category,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	URL       string   `json:"url"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type TaskInput struct {
	Title    string   `json:"title"`
	Status   string   `json:"status,omitempty"`
	Priority string   `json:"priority,omitempty"`
	DueDate  string   `json:"due_date,omitempty"`
	Category string   `json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Notes    string   `json:"notes,omitempty"`
}
