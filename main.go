package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

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
	r.Route("/{bookmarkID}", func(r chi.Router) {
		r.Get("/", getBookmark)
		// r.Put("/", updateBookmark)
		// r.Delete("/", deleteBookmark)
	})
	return r
}

func listBookmarks(w http.ResponseWriter, r *http.Request) {

}
