package main

// Bookmark Models
type Bookmark struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	Archived    bool   `json:"archived,omitempty"`
	Deleted     bool   `json:"deleted,omitempty"`
}

type Bookmarks struct {
	Bookmarks []Bookmark
}
