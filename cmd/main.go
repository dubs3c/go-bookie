package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	gobookie "github.com/dubs3c/go-bookie"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	pgxpool, err := gobookie.DBInit()

	if err != nil {
		os.Exit(1)
	}

	s := &gobookie.Server{
		DB:     pgxpool,
		Router: r,
	}

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1/bookmarks", func(r chi.Router) {
			r.Get("/", s.ListBookmarks)
			r.Post("/", s.CreateBookmark)
			r.Route("/{bookmarkID}", func(r chi.Router) {
				r.Get("/", s.GetBookmark)
				r.Put("/", s.UpdateBookmark)
				r.Patch("/", s.UpdateBookmark)
				r.Delete("/", s.DeleteBookmark)
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

	banner()
	log.Printf("[+] Starting server on %s...", HTTPServer.Addr)
	HTTPServer.ListenAndServe()
}
