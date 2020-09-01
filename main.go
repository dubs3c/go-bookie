package main

import (
	"net/http"
    "context"
    "os"
    "fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
    "github.com/jackc/pgx/v4/pgxpool"
)

// Bookmark Models
type Bookmark struct {
    ID uint `json:"id"`
    Title string `json:"title"`
    Description string `json:"description"`
    Body string `json:"body"`
    Image string `json:"image"`
    URL string `json:"url"`
    Archived bool `json:"archived"`
    Deleted bool `json:"deleted"`
}


func dbinit() (*Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(),"host=127.0.0.1 port=5432 user=bookie dbname=bookie password=bookie sslmode=disable")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer dbpool.Close()
	return dbpool,err
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	pgxpool, err = dbinit()

	if err != nil {
		os.Exit(1)
	}

	r.Route("/api", func(r chi.Router) {
	    r.Route("/v1", func(r chi.Router) {
            r.Mount("/bookmarks", bookmarksRouter())
			//r.Mount("/tags", tagsRouter())
			//r.Mount("/users", usersRouter())
			//r.Mount("/export", exportRouter())
			//r.Mount("/integrations", integrationsRouter())
			//r.Mount("/me", meRouter())
		})
	})

	http.ListenAndServe(":8000", r)
}


func bookmarksRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", listBookmarks)
    r.Post("/", CreateBookmark)
	r.Route("/{bookmarkID}", func(r chi.Router) {
		// r.Get("/", getBookmark)
		// r.Put("/", updateBookmark)
		// r.Delete("/", deleteBookmark)
	})
	return r
}

func CreateBookmark(w http.ResponseWriter, r *http.Request) {

}

func listBookmarks(w http.ResponseWriter, r *http.Request) {

}



