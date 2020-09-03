package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	jsoniter "github.com/json-iterator/go"
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

	HTTPServer := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("[+] Starting server on %s...", HTTPServer.Addr)
	HTTPServer.ListenAndServe()
}

// CreateBookmark api
func (s *Server) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	data := &Bookmark{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		respondWithError(w, 400, "Could not decode json")
	}

	res, err := s.DB.Exec(context.Background(), "INSERT INTO bookmarks(title,description,body,image,url,archived,deleted) values($1,$2,$3,$4,$5,$6,$7)", data.Title, data.Description, data.Body, data.Image, data.URL, data.Archived, data.Deleted)

	if err != nil {
		respondWithError(w, 500, "Could not create bookmark")
	}

	if res.RowsAffected() != 1 {
		respondWithError(w, 500, "Bookmark was not")
	}

	respondWithStatusCode(w, 201)
}

// ListBookmarks Return all bookmarks
func (s *Server) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	var bookmarks []*Bookmark
	rows, err := s.DB.Query(context.Background(), "select id, title, description, body, image, url, archived, deleted from bookmarks")
	defer rows.Close()

	if err != nil {
		respondWithError(w, 500, "Something went wrong while fetching bookmarks")
	}

	for rows.Next() {
		n := new(Bookmark)
		err = rows.Scan(&n.ID, &n.Title, &n.Description, &n.Body, &n.Image, &n.URL, &n.Archived, &n.Deleted)
		if err != nil {
			fmt.Println(err)
		}
		bookmarks = append(bookmarks, n)
	}

	respondWithJson(w, 200, bookmarks)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithStatusCode(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return
}
