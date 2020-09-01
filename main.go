package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Server struct
type Server struct {
	Context context.Context
	DB      *pgxpool.Pool
	Router  *chi.Mux
}

func dbinit() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), "host=127.0.0.1 port=5432 user=bookie dbname=bookie password=bookie sslmode=disable")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	//defer dbpool.Close()
	return dbpool, err
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	pgxpool, err := dbinit()

	if err != nil {
		os.Exit(1)
	}

	s := &Server{
		DB:     pgxpool,
		Router: r,
	}

	s.Router.Route("/api", func(r chi.Router) {
		s.Router.Route("/v1/bookmarks", func(r chi.Router) {
			s.Router.Get("/", s.ListBookmarks)
			s.Router.Post("/", s.CreateBookmark)
			s.Router.Route("/{bookmarkID}", func(r chi.Router) {
				// r.Get("/", getBookmark)
				// r.Put("/", updateBookmark)
				// r.Delete("/", deleteBookmark)
			})
		})
	})

	http.ListenAndServe(":8000", r)
}

// CreateBookmark api
func (s *Server) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	data := &BookmarkRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	res, err := s.DB.Exec(context.Background(), "insert into bookmarks(title,description,body,image,url,archived,deleted) values($1,$2,$3,$4,$5,$6,$7)", data.Title, data.Description, data.Body, data.Image, data.URL, data.Archived, data.Deleted)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if res.RowsAffected() != 1 {
		fmt.Println("Did not insert row")
		render.Render(w, r, GeneralErrRequest("Something went wrong inserting to db"))
		return
	}
}

// ListBookmarks Return all bookmarks
func (s *Server) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	var bookmarks []Bookmark
	rows, err := s.DB.Query(context.Background(), "select id, title, description, body, image, url, archived, deleted from bookmarks")

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	defer rows.Close()

	for rows.Next() {
		n := new(Bookmark)
		err = rows.Scan(&n.ID, &n.Title, &n.Description, &n.Body, &n.Image, &n.URL, &n.Archived, &n.Deleted)
		if err != nil {
			fmt.Println(err)
		}
		bookmarks = append(bookmarks, *n)
	}
	BookmarksJSON, _ := json.Marshal(bookmarks)
	w.Header().Set("Content-Type", "application/json")
	w.Write(BookmarksJSON)
}

// Bind aaa
func (a *BookmarkRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Bookmark == nil {
		return errors.New("missing required Bookmark fields")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	a.ProtectedID = ""                                   // unset the protected ID
	a.Bookmark.Title = strings.ToLower(a.Bookmark.Title) // as an example, we down-case
	return nil
}

// Render an error response
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest aaa
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// GeneralErrRequest gen err
func GeneralErrRequest(err string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 500,
		StatusText:     "Server Error.",
		ErrorText:      err,
	}
}
