package gobookie

import "time"

type UserIDContext string

func (c UserIDContext) String() string {
	return string(c)
}

// Bookmark Models
type Bookmark struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	Image       string    `json:"image"`
	URL         string    `json:"url"`
	Archived    bool      `json:"archived"`
	Deleted     bool      `json:"deleted"`
	Tags        string    `json:"tags"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Bookmark list model
type BookmarkList struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	Archived    bool   `json:"archived"`
	Deleted     bool   `json:"deleted"`
	Tags        string `json:"tags"`
}

type PaginatedBookmarks struct {
	Page       int             `json:"page"`
	TotalPages int             `json:"totalPages"`
	Limit      int             `json:"limit"`
	Data       []*BookmarkList `json:"data"`
}

// CreateBookmarkRequest - Used for creating a bookmark
type CreateBookmarkRequest struct {
	URL string `json:"url"`
}

type CreateTagRequest struct {
	BookmarkID int    `json:"bookmarkID"`
	TagName    string `json:"tagName"`
}

type UpdateTagRequest struct {
	TagName string `json:"tagName"`
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

type Tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Tags struct {
	Tags []Tag
}
