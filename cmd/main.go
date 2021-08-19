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

	pgxpool, err := gobookie.DBInit()

	if err != nil {
		os.Exit(1)
	}

	s := &gobookie.Server{
		DB:     pgxpool,
		Router: r,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/login", s.UserLogin)
		r.With(s.Authentication).Get("/logout", s.UserLogout)
		r.Post("/register", s.UserRegister)
		r.Group(func(r chi.Router) {
			r.Use(s.Authentication)
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
				r.Use(s.Authentication)
				r.Get("/", s.ListTags)
				r.Post("/", s.CreateTag)
				r.Delete("/", s.DeleteTag)
				r.Put("/", s.UpdateTag)
			})
		})

	})

	HTTPServer := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   7 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	banner()
	log.Printf("[+] Starting server on %s...", HTTPServer.Addr)
	HTTPServer.ListenAndServe()
}
