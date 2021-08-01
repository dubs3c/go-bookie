package gobookie

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

// BookmarkRepositoryInsert - Insert bookmark into database
func (s *Server) BookmarkRepositoryInsert(bookmark *CreateBookmarkRequest) (int, error) {
	var last int = 0
	err := s.DB.QueryRow(context.Background(), "INSERT INTO bookmarks(url) values($1) RETURNING id", bookmark.URL).Scan(&last)
	return last, err
}

// BookmarkRepositoryCount Get total amount of bookmarks
func (s *Server) BookmarkRepositoryCount() (int, error) {
	var count int
	rows, err := s.DB.Query(context.Background(), "SELECT count(*) as count FROM bookmarks")

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

// BookmarkRepositoryGetAllBookmarks - Get all bookmarks from database
func (s *Server) BookmarkRepositoryGetAllBookmarks(page int, limit int, archived bool, deleted bool) ([]*BookmarkList, error) {
	var bm []*BookmarkList
	var offset int

	if page <= 1 {
		offset = 0
	} else {
		offset = limit*page - limit
	}

	rows, err := s.DB.Query(context.Background(),
		"select id, title, description, image, url, archived, deleted from bookmarks where deleted=$4 and archived=$3 order by created_at desc LIMIT $1 OFFSET $2", limit, offset, archived, deleted)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		n := new(BookmarkList)
		err = rows.Scan(&n.ID, &n.Title, &n.Description, &n.Image, &n.URL, &n.Archived, &n.Deleted)
		if err != nil {
			log.Println("Error scanning bookmarkList:", err)
		}
		bm = append(bm, n)
	}

	return bm, nil
}

// BookmarkRepositoryGetBook.markByID - Get a specifc bookmark by its database id
func (s *Server) BookmarkRepositoryGetBookmarkByID(bookmarkID string) (Bookmark, error) {
	var bookmark Bookmark

	row := s.DB.QueryRow(context.Background(), "SELECT id, title, description, body, image, url, archived, deleted FROM bookmarks WHERE id=$1", bookmarkID)

	err := row.Scan(&bookmark.ID, &bookmark.Title, &bookmark.Description, &bookmark.Body, &bookmark.Image, &bookmark.URL, &bookmark.Archived, &bookmark.Deleted)

	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		log.Println(err)
		return bookmark, err
	}

	return bookmark, nil
}

// BookmarkRepositoryDeleteBookmarkByID - Delete a specifc bookmark by its database id
func (s *Server) BookmarkRepositoryDeleteBookmarkByID(bookmarkID string) (int64, error) {
	rows, err := s.DB.Query(context.Background(), "DElETE FROM bookmarks WHERE id=$1", bookmarkID)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	return rows.CommandTag().RowsAffected(), err
}

// BookmarkRepositoryTrashBookmarkByID - Mark bookmark as deleted
func (s *Server) BookmarkRepositoryTrashBookmarkByID(bookmarkID string) (int64, error) {
	rows, err := s.DB.Query(context.Background(), "UPDATE bookmarks SET deleted=true WHERE id=$1", bookmarkID)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	return rows.CommandTag().RowsAffected(), err
}

// BookmarkRepositoryUpdateBookmark - Update a specifc bookmark by its database id
func (s *Server) BookmarkRepositoryUpdateBookmark(bookmark Bookmark) error {
	res, err := s.DB.Exec(context.Background(), "UPDATE bookmarks SET title=$1, description=$2, body=$3, image=$4, url=$5, archived=$6, deleted=$7 WHERE id=$8", bookmark.Title, bookmark.Description, bookmark.Body, bookmark.Image, bookmark.URL, bookmark.Archived, bookmark.Deleted, bookmark.ID)

	if res.RowsAffected() != 1 {
		return err
	}

	return nil
}

// BookmarkRepositoryPatchBookmark - Patch a specifc bookmark by its database id
func (s *Server) BookmarkRepositoryPatchBookmark(bookmark Bookmark) error {
	sql := `UPDATE bookmarks SET
  			title = COALESCE($1, title),
			description = COALESCE($2, description),
  			body = COALESCE($3, body),
  			image = COALESCE($4, image),
			url = COALESCE($5, url),
			archived = COALESCE($6, archived),
  			deleted = COALESCE($7, deleted)
			WHERE id = $8;`
	res, err := s.DB.Exec(context.Background(), sql, bookmark.Title, bookmark.Description, bookmark.Body, bookmark.Image, bookmark.URL, bookmark.Archived, bookmark.Deleted, bookmark.ID)

	if res.RowsAffected() != 1 {
		return err
	}

	return nil
}
