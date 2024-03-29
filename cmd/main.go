package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	gobookie "github.com/dubs3c/go-bookie"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func banner() {
	fmt.Println("")
	fmt.Println(" ,---,.                           ,-.            					")
	fmt.Println(",'  .'  \\                      ,--/ /|   ,--,              		")
	fmt.Println(",---.' .' |   ,---.     ,---.  ,--. :/ | ,--.'|              		")
	fmt.Println("|   |  |: |  '   ,'\\   '   ,'\\ :  : ' /  |  |,               	")
	fmt.Println(":   :  :  / /   /   | /   /   ||  '  /   `--'_       ,---.   		")
	fmt.Println(":   |    ; .   ; ,. :.   ; ,. :'  |  :   ,' ,'|     /     \\  		")
	fmt.Println("|   :     \\'   | |: :'   | |: :|  |   \\  '  | |    /    /  | 	")
	fmt.Println("|   |   . |'   | .; :'   | .; :'  : |. \\ |  | :   .    ' / | 		")
	fmt.Println("'   :  '; ||   :    ||   :    ||  | ' \\ \\'  : |__ '   ;   /| 	")
	fmt.Println("|   |  | ;  \\  \\  /  \\   \\  / '  : |--' |  | '.'|'   |  / | 	")
	fmt.Println("|   :   /    `----'    `----'  ;  |,'    ;  :    ;|   :    | 		")
	fmt.Println("|   | ,'                       '--'      |  ,   /  \\   \\  /  	")
	fmt.Println("`----'                                    ---`-'    `----'   		")
	fmt.Println("")
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Cookies", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	pgxpool, err := gobookie.DBInit()

	if err != nil {
		os.Exit(1)
	}

	s := &gobookie.Server{
		DB:     pgxpool,
		Router: r,
	}

	r.Route("/v1", func(r chi.Router) {
		r.Use(gobookie.Authentication(s))

		r.Route("/bookmarks", func(r chi.Router) {
			r.Get("/", s.ListBookmarks)
			r.Post("/", s.CreateBookmark)
			r.Route("/{bookmarkID}", func(r chi.Router) {
				r.Get("/", s.GetBookmark)
				r.Put("/", s.UpdateBookmark)
				r.Patch("/", s.UpdateBookmark)
				r.Delete("/", s.DeleteBookmark)
			})
		})

		r.Route("/tags", func(r chi.Router) {
			r.Get("/", s.ListTags)
			r.Post("/", s.CreateTag)
			r.Delete("/", s.DeleteTag)
			r.Put("/", s.UpdateTag)
		})
		r.Route("/users", func(r chi.Router) {
			r.Use(gobookie.IsAdmin)
			r.Get("/", s.ListUsers)
			r.Post("/", s.CreateUser)
			r.Delete("/", s.DeleteUser)
			r.Put("/", s.UpdateUser)
		})

		r.Post("/logout", s.UserLogout)
	})

	// These endpoints do not require auth
	r.Post("/login", s.UserLogin)

	HTTPServer := &http.Server{
		Addr:           "0.0.0.0:8080",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   7 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	banner()
	log.Printf("[+] Starting server on %s...", HTTPServer.Addr)
	HTTPServer.ListenAndServe()
}
