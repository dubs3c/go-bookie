package gobookie

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	goose "github.com/advancedlogic/GoOse"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	jsoniter "github.com/json-iterator/go"
)

// Server struct
type Server struct {
	DB     *pgxpool.Pool
	Router *chi.Mux
}

func getarticle(ctx context.Context, db *pgxpool.Pool, id int, url string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()
	g := goose.New()
	article, err := g.ExtractFromURL(url)
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, "UPDATE bookmarks SET body=$2, title=$3, description=$4, image=$5 WHERE id=$1", id, article.RawHTML, article.Title, article.MetaDescription, article.TopImage)
	if err != nil {
		cancel()
		return err
	}
	return nil
}

// CreateBookmark Create bookmark
func (s *Server) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	data := &Bookmark{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		respondWithError(w, 400, "Could not decode json")
	}

	LastInsertedID, err := s.BookmarkRepositoryInsert(data)

	go func() {
		err := getarticle(context.Background(), s.DB, LastInsertedID, data.URL)
		if err != nil {
			log.Printf("ERROR: Could not fetch website %s: %s", data.URL, err.Error())
		}
	}()

	if err != nil {
		respondWithError(w, 500, "Could not create bookmark")
		return
	}

	if LastInsertedID == 0 {
		respondWithError(w, 500, "Bookmark was not created")
		return
	}

	respondWithStatusCode(w, 201)
}

// ListBookmarks Return all bookmarks
func (s *Server) ListBookmarks(w http.ResponseWriter, r *http.Request) {

	bookmarks, err := s.BookmarkRepositoryGetAllBookmarks()

	if err != nil {
		respondWithError(w, 500, "Something went wrong while fetching bookmarks")
		return
	}

	respondWithJSON(w, 200, bookmarks)
}

// GetBookmark Get a specific bookmark by its ID
func (s *Server) GetBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	bookmark, err := s.BookmarkRepositoryGetBookmarkByID(bookmarkID)

	if bookmark == (Bookmark{}) {
		respondWithStatusCode(w, 404)
		return
	}

	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Could not fetch bookmark")
		return
	}

	respondWithJSON(w, 200, bookmark)
}

// DeleteBookmark Delete a specific bookmark by its ID
func (s *Server) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	bookmark, err := s.BookmarkRepositoryDeleteBookmarkByID(bookmarkID)

	if err != nil {
		respondWithError(w, 500, "Could not delete bookmark")
		return
	}

	respondWithJSON(w, 200, bookmark)
}

// UpdateBookmark Update a bookmark by its ID
func (s *Server) UpdateBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	data := &Bookmark{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		respondWithError(w, 400, "Could not decode json")
		return
	}

	bookmark, err := s.BookmarkRepositoryGetBookmarkByID(bookmarkID)

	if err != nil {
		respondWithError(w, 500, "Could not fetch bookmark")
		return
	}

	bookmark.Archived = data.Archived
	bookmark.Body = data.Body
	bookmark.Deleted = data.Deleted
	bookmark.Description = data.Description
	bookmark.Image = data.Image
	bookmark.Title = data.Title
	bookmark.URL = data.URL

	if err = s.BookmarkRepositoryUpdateBookmark(bookmark); err != nil {
		respondWithError(w, 500, "Could not update bookmark")
		return
	}

	respondWithStatusCode(w, 200)
}

/* ----------------- Response methods ----------------- */

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return
}

func respondWithStatusCode(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return
}
