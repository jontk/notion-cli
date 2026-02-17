package models

type Post struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Content         string   `json:"content,omitempty"`
	Status          string   `json:"status"`
	Week            int      `json:"week,omitempty"`
	Pillar          string   `json:"pillar,omitempty"`
	PublishDate     string   `json:"publish_date,omitempty"`
	PublishedDate   string   `json:"published_date,omitempty"`
	BlogURL         string   `json:"blog_url,omitempty"`
	DistributedTo   []string `json:"distributed_to,omitempty"`
	DistributedDate string   `json:"distributed_date,omitempty"`
	LinkedInDraft   string   `json:"linkedin_draft,omitempty"`
	TwitterThread   string   `json:"twitter_thread,omitempty"`
	HNTitle         string   `json:"hn_title,omitempty"`
	RedditTitle     string   `json:"reddit_title,omitempty"`
	Hashtags        []string `json:"hashtags,omitempty"`
	URL             string   `json:"url"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

type PostInput struct {
	Title           string   `json:"title,omitempty"`
	Content         string   `json:"content,omitempty"`
	Status          string   `json:"status,omitempty"`
	Week            int      `json:"week,omitempty"`
	Pillar          string   `json:"pillar,omitempty"`
	PublishDate     string   `json:"publish_date,omitempty"`
	PublishedDate   string   `json:"published_date,omitempty"`
	BlogURL         string   `json:"blog_url,omitempty"`
	DistributedTo   []string `json:"distributed_to,omitempty"`
	DistributedDate string   `json:"distributed_date,omitempty"`
	LinkedInDraft   string   `json:"linkedin_draft,omitempty"`
	TwitterThread   string   `json:"twitter_thread,omitempty"`
	HNTitle         string   `json:"hn_title,omitempty"`
	RedditTitle     string   `json:"reddit_title,omitempty"`
	Hashtags        []string `json:"hashtags,omitempty"`
}
