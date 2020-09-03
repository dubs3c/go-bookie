package gobookie

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

func (s *Server) BookmarkRepositoryInsert(bookmark *Bookmark) (int64, error) {
	res, err := s.DB.Exec(context.Background(), "INSERT INTO bookmarks(title,description,body,image,url,archived,deleted) values($1,$2,$3,$4,$5,$6,$7)", bookmark.Title, bookmark.Description, bookmark.Body, bookmark.Image, bookmark.URL, bookmark.Archived, bookmark.Deleted)
	return res.RowsAffected(), err
}

func (s *Server) BookmarkRepositoryGetAllBookmarks() ([]*Bookmark, error) {
	var bm []*Bookmark
	rows, err := s.DB.Query(context.Background(), "select id, title, description, body, image, url, archived, deleted from bookmarks")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		n := new(Bookmark)
		err = rows.Scan(&n.ID, &n.Title, &n.Description, &n.Body, &n.Image, &n.URL, &n.Archived, &n.Deleted)
		if err != nil {
			fmt.Println(err)
		}
		bm = append(bm, n)
	}
	return bm, nil
}

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

func (s *Server) BookmarkRepositoryDeleteBookmarkByID(bookmarkID string) (int64, error) {
	rows, err := s.DB.Query(context.Background(), "DElETE FROM bookmarks WHERE id=$1", bookmarkID)
	defer rows.Close()

	return rows.CommandTag().RowsAffected(), err
}

func (s *Server) BookmarkRepositoryUpdateBookmark(bookmark Bookmark) error {
	res, err := s.DB.Exec(context.Background(), "UPDATE bookmarks SET title=$1, description=$2, body=$3, image=$4, url=$5, archived=$6, deleted=$7 WHERE id=$8", bookmark.Title, bookmark.Description, bookmark.Body, bookmark.Image, bookmark.URL, bookmark.Archived, bookmark.Deleted, bookmark.ID)

	if res.RowsAffected() != 1 {
		return err
	}

	return nil
}
