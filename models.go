package gobookie

// Bookmark Models
type Bookmark struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	Archived    bool   `json:"archived"`
	Deleted     bool   `json:"deleted"`
}

// CreateBookmarkRequest - Used for creating a bookmark
type CreateBookmarkRequest struct {
	URL string `json:"url"`
}

// UpdateBookmarkRequest - PUT/PATCH object
type UpdateBookmarkRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Body        string `json:"body,omitempty"`
	Image       string `json:"image,omitempty"`
	URL         string `json:"url,omitempty"`
	Archived    *bool  `json:"archived,omitempty"`
	Deleted     *bool  `json:"deleted,omitempty"`
}

// Bookmarks - A collection of bookmarks
type Bookmarks struct {
	Bookmarks []Bookmark
}
