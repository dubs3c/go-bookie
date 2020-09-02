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

// BookmarkRequest the req
type BookmarkRequest struct {
	*Bookmark

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

type BookmarkResponse struct {
	*Bookmark

	//User *UserPayload `json:"user,omitempty"`

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

// ErrResponse err
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
