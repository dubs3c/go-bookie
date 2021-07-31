package gobookie

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	goose "github.com/advancedlogic/GoOse"
	"github.com/go-chi/chi/v5"
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
	var pageNumber int
	var bookmarkLimit int
	var total int
	var err error

	if page := r.URL.Query().Get("page"); page == "" {
		pageNumber = 1
	} else {
		if pageNumber, err = strconv.Atoi(page); err != nil {
			RespondWithError(w, 400, "page should be integer")
			return
		}
	}

	if limit := r.URL.Query().Get("limit"); limit == "" {
		bookmarkLimit = 10
	} else {
		if bookmarkLimit, err = strconv.Atoi(limit); err != nil {
			RespondWithError(w, 400, "limit should be integer")
			return
		}
	}

	bookmarks, err := s.BookmarkRepositoryGetAllBookmarks(pageNumber, bookmarkLimit)

	if err != nil {
		RespondWithError(w, 500, "Something went wrong while fetching bookmarks")
		log.Println(err)
		return
	}

	bookmarksAmount, _ := s.BookmarkRepositoryCount()

	if bookmarksAmount > bookmarkLimit {
		total = (bookmarksAmount + bookmarkLimit - 1) / bookmarkLimit
	} else {
		total = (bookmarkLimit + bookmarksAmount - 1) / bookmarksAmount
	}

	paginatedResult := PaginatedBookmarks{
		Page:       pageNumber,
		TotalPages: total,
		Limit:      bookmarkLimit,
		Data:       bookmarks,
	}

	RespondWithJSON(w, 200, paginatedResult)
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

// DeleteBookmark Delete a specific bookmark by its ID
func (s *Server) MoveBookmarkToTrash(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	bookmark, err := s.BookmarkRepositoryTrashBookmarkByID(bookmarkID)

	if err != nil {
		RespondWithError(w, 500, "Could not mark bookmark as deleted")
		return
	}
	RespondWithJSON(w, 200, bookmark)
}

// DeleteBookmark Delete a specific bookmark by its ID
func (s *Server) CorsOptions(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, 200, "")
}

// UpdateBookmark Update a bookmark by its ID
func (s *Server) UpdateBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := chi.URLParam(r, "bookmarkID")

	var data UpdateBookmarkRequest
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		RespondWithError(w, 400, "Could not read json")
		return
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
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
		err := patch(data, &bookmark)
		if err != nil {
			RespondWithError(w, 500, "Something went wrong when trying to patch bookmark...")
			return
		}
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

// Error - A function that returns a custom error message
func Error(msg string) error {
	return errors.New(msg)
}

// Patch - Update a destination struct with non empty values from a src struct
// This function looks at src and takes the non empty values
// and assigns it to dst on the correct field
func patch(v interface{}, bookmark *Bookmark) error {
	src := reflect.ValueOf(v)
	if !src.IsValid() {
		return Error("src reflect not valid")
	}

	if src.Type().Kind() != reflect.Struct {
		return Error("src parameter is not a struct")
	}

	dst := reflect.ValueOf(bookmark).Elem()
	if !dst.IsValid() {
		return Error("dst reflect not valid")
	}

	for i := 0; i < src.NumField(); i++ {
		if src.Field(i).Kind() == reflect.String {
			if src.Field(i).String() != "" {
				fieldName := src.Type().Field(i).Name
				x := dst.FieldByName(fieldName)
				if x.IsValid() {
					if !x.CanSet() {
						return Error("ERROR: Can not update destintion struct")
					}
					x.SetString(src.Field(i).String())
				}
			}
		}

		if src.Field(i).Kind() == reflect.Ptr {
			l := src.Field(i).Elem()
			if l.IsValid() {
				fieldPointer := reflect.Indirect(src.Field(i))
				if fieldPointer.Type().Kind() == reflect.Bool {
					fieldName := src.Type().Field(i).Name
					x := dst.FieldByName(fieldName)

					if !x.CanSet() {
						return Error("ERROR: Can not update destintion struct")
					}

					x.SetBool(fieldPointer.Bool())
				}
			}
		}
	}
	return nil
}
