package gobookie

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

// GetAccessToken Looks up access token in DB and returns corresponding user id
func (s *Server) GetAccessToken(token string) (int, error) {
	var userID int = 0
	err := s.DB.QueryRow(context.Background(), `
		SELECT user_fk
		FROM access_tokens
		WHERE token = $1
		`, token).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, err
}

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
func (s *Server) BookmarkRepositoryGetAllBookmarks(page int, limit int, archived bool, deleted bool, tags []string) ([]*BookmarkList, error) {

	var (
		bm     []*BookmarkList
		offset int
		rows   pgx.Rows
		err    error
	)

	bm = make([]*BookmarkList, 0)

	if page <= 1 {
		offset = 0
	} else {
		offset = limit*page - limit
	}

	sqlSelect := `
	SELECT b.id, b.title, b.description, b.image, b.url, b.archived, b.deleted, string_agg(COALESCE(t.name, '')::text, ',') AS tags
	FROM bookmarks AS b
	LEFT JOIN bookmark_has_tags AS bht ON bht.bookmark_fk = b.id
	LEFT JOIN tags AS t ON t.id = bht.tag_fk `

	sqlWhere := "WHERE b.deleted = $1 AND b.archived = $2 "

	if len(tags) >= 1 {
		sqlWhere += "AND t.name = ANY($5) "
	}

	sqlGroup := `
	GROUP BY b.id
	ORDER BY b.id desc
	LIMIT $3 OFFSET $4`

	sqlQuery := sqlSelect + sqlWhere + sqlGroup

	if len(tags) >= 1 {
		rows, err = s.DB.Query(context.Background(), sqlQuery, deleted, archived, limit, offset, tags)
	} else {
		rows, err = s.DB.Query(context.Background(), sqlQuery, deleted, archived, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		n := new(BookmarkList)
		err = rows.Scan(&n.ID, &n.Title, &n.Description, &n.Image, &n.URL, &n.Archived, &n.Deleted, &n.Tags)
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

	// TODO - Consider creating two queries so we can extract all tag data and return a list of tags
	row := s.DB.QueryRow(context.Background(),
		`SELECT bm.id, bm.title, bm.description, bm.body, bm.image,
		bm.url, bm.archived, bm.deleted, string_agg(COALESCE(t.name, '')::text, ',') AS tags, bm.created_at
	FROM bookmarks AS bm
	LEFT JOIN bookmark_has_tags AS bht
	ON bht.bookmark_fk = bm.id
	LEFT JOIN tags AS t
	ON t.id = bht.tag_fk
	WHERE bm.id=$1
	GROUP BY bm.id
	`, bookmarkID)

	err := row.Scan(&bookmark.ID, &bookmark.Title, &bookmark.Description, &bookmark.Body, &bookmark.Image, &bookmark.URL, &bookmark.Archived, &bookmark.Deleted, &bookmark.Tags, &bookmark.CreatedAt)

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

func (s *Server) TagsRepositoryGetTags() (Tags, error) {
	tags := Tags{Tags: []Tag{}}

	rows, err := s.DB.Query(context.Background(),
		"SELECT id, name FROM tags")

	if err != nil {
		return tags, err
	}

	defer rows.Close()

	for rows.Next() {
		t := new(Tag)
		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			log.Println("Error scanning tags:", err)
		}
		tags.Tags = append(tags.Tags, *t)
	}

	return tags, nil
}

func (s *Server) TagsRepositoryGetTagsByBookmarkID(bookmarkID int) (Tags, error) {

	tags := Tags{Tags: []Tag{}}

	rows, err := s.DB.Query(context.Background(),
		`SELECT t.id, t.name FROM tags AS t
		LEFT JOIN bookmark_has_tags AS bt
		ON t.id = bt.id
		WHERE bt.id = $1
		`, bookmarkID)

	if err != nil {
		return tags, err
	}

	defer rows.Close()

	for rows.Next() {
		t := new(Tag)
		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			log.Println("Error scanning tags by bookmark id:", err)
		}
		tags.Tags = append(tags.Tags, *t)
	}

	return tags, nil
}

func (s *Server) TagsRepositoryCreateTag(bookmarkID int, TagName string) error {
	var tagID int = 0
	// TODO - Make a SQL function instead
	err := s.DB.QueryRow(context.Background(), "INSERT INTO tags(name) values($1) RETURNING id", TagName).Scan(&tagID)

	if err != nil {
		return err
	}

	_, err = s.DB.Query(context.Background(), "INSERT INTO bookmark_has_tags(bookmark_fk, tag_fk) values($1, $2)", bookmarkID, tagID)

	return err
}

func (s *Server) TagsRepositoryUpdateTagByID(tagID int, TagNameUpdate string) error {
	_, err := s.DB.Query(context.Background(), "UPDATE tags SET name = $1 WHERE id = $2", TagNameUpdate, tagID)
	return err
}

func (s *Server) TagsRepositoryDeleteTagByBookmarkIDAndTagID(bookmarkID int, tagName string) error {
	_, err := s.DB.Query(context.Background(),
		`DELETE FROM bookmark_has_tags AS bht
		WHERE bht.bookmark_fk = $1 AND tag_fk IN (
			SELECT t.id FROM tags AS t
			WHERE t.name = $2
		)`, bookmarkID, tagName)
	return err
}
