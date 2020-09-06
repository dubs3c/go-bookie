package gobookie

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
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
	data := &CreateBookmarkRequest{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		RespondWithError(w, 400, "Could not decode json")
	}

	if data.URL == "" {
		RespondWithError(w, 400, "You need to enter a URL")
		return
	}

	LastInsertedID, err := s.BookmarkRepositoryInsert(data)

	go func() {
		err := getarticle(context.Background(), s.DB, LastInsertedID, data.URL)
		if err != nil {
			log.Printf("ERROR: Could not fetch website %s: %s", data.URL, err.Error())
		}
	}()

	if err != nil {
		RespondWithError(w, 500, "Could not create bookmark")
		return
	}

	if LastInsertedID == 0 {
		RespondWithError(w, 500, "Bookmark was not created")
		return
	}

	RespondWithStatusCode(w, 201)
}

// ListBookmarks Return all bookmarks
func (s *Server) ListBookmarks(w http.ResponseWriter, r *http.Request) {

	bookmarks, err := s.BookmarkRepositoryGetAllBookmarks()

	if err != nil {
		RespondWithError(w, 500, "Something went wrong while fetching bookmarks")
		return
	}

	RespondWithJSON(w, 200, bookmarks)
}

// GetBookmark Get a specific bookmark by its ID
func (s *Server) GetBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	bookmark, err := s.BookmarkRepositoryGetBookmarkByID(bookmarkID)

	if bookmark == (Bookmark{}) {
		RespondWithStatusCode(w, 404)
		return
	}

	if err != nil {
		log.Println(err)
		RespondWithError(w, 500, "Could not fetch bookmark")
		return
	}

	RespondWithJSON(w, 200, bookmark)
}

// DeleteBookmark Delete a specific bookmark by its ID
func (s *Server) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	bookmark, err := s.BookmarkRepositoryDeleteBookmarkByID(bookmarkID)

	if err != nil {
		RespondWithError(w, 500, "Could not delete bookmark")
		return
	}

	RespondWithJSON(w, 200, bookmark)
}

// UpdateBookmark Update a bookmark by its ID
func (s *Server) UpdateBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	var data UpdateBookmarkRequest
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	log.Println(string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Println(data)
		log.Println(err.Error())
		RespondWithError(w, 400, "Could not decode json")
		return
	}

	bookmark, err := s.BookmarkRepositoryGetBookmarkByID(bookmarkID)

	if err != nil {
		RespondWithError(w, 500, "Could not fetch bookmark")
		return
	}

	// Because PUT and PATCH performs identical operations,
	// this is an easy way to perform different database
	// operations based on the method.
	switch r.Method {
	case "PUT":

		if data.Archived != nil {
			bookmark.Archived = *data.Archived
		}
		if data.Deleted != nil {
			bookmark.Deleted = *data.Deleted
		}
		bookmark.Body = data.Body
		bookmark.Description = data.Description
		bookmark.Image = data.Image
		bookmark.Title = data.Title
		bookmark.URL = data.URL

		if err = s.BookmarkRepositoryUpdateBookmark(bookmark); err != nil {
			RespondWithError(w, 500, "Could not update bookmark")
			return
		}
	case "PATCH":
		log.Println(bookmark)
		fuck(data, &bookmark)
		log.Println(bookmark)
		if err = s.BookmarkRepositoryPatchBookmark(bookmark); err != nil {
			RespondWithError(w, 500, "Could not patch bookmark")
			return
		}
	default:
		RespondWithStatusCode(w, 400)
		return
	}

	RespondWithStatusCode(w, 200)
}

func fuck(v interface{}, bookmark *Bookmark) {
	t := reflect.ValueOf(v)
	if !t.IsValid() {
		log.Println("src reflect not valid")
		return
	}

	bm := reflect.ValueOf(bookmark).Elem()
	if !bm.IsValid() {
		log.Println("src reflect not valid")
		return
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Kind() == reflect.String {
			if t.Field(i).String() != "" {
				fieldName := t.Type().Field(i).Name
				x := bm.FieldByName(fieldName)
				if x.IsValid() {
					if !x.CanSet() {
						log.Println("ERROR: Can not update destintion struct")
						return
					}
					x.SetString(t.Field(i).String())
				}

			}
		}
	}
}
