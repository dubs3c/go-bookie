package gobookie

import (
	"testing"
)

func TestPatchUtil(t *testing.T) {

	trueValue := true

	test1 := UpdateBookmarkRequest{
		URL:     "www.coolurl.com",
		Title:   "PATCHED",
		Deleted: &trueValue,
	}

	bookmark := Bookmark{
		ID:          1,
		Title:       "Cool website title",
		Description: "A description",
		URL:         "www.booringurl.com",
		Image:       "image",
		Archived:    false,
		Deleted:     false,
	}

	patch(test1, &bookmark)

	if bookmark.Title != test1.Title {
		t.Errorf("Expected bookmark.Title to be '%s', got %s", test1.Title, bookmark.Title)
	}

	if bookmark.URL != test1.URL {
		t.Errorf("Expected bookmark.URL to be '%s', got %s", test1.URL, bookmark.URL)
	}

	if bookmark.Deleted != *test1.Deleted {
		t.Errorf("Expected bookmark.Deleted to be %t, got %t", *test1.Deleted, bookmark.Deleted)
	}

}
