package models

type Event struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	Type      string   `json:"type,omitempty"`
	Location  string   `json:"location,omitempty"`
	Attendees []string `json:"attendees,omitempty"`
	Status    string   `json:"status,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	URL       string   `json:"url"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type EventInput struct {
	Title     string   `json:"title"`
	Date      string   `json:"date,omitempty"`
	Type      string   `json:"type,omitempty"`
	Location  string   `json:"location,omitempty"`
	Attendees []string `json:"attendees,omitempty"`
	Status    string   `json:"status,omitempty"`
	Notes     string   `json:"notes,omitempty"`
}
