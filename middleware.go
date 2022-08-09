package gobookie

import (
	"context"
	"log"
	"net/http"
	"time"
)

// HTTP middleware setting a value on the request context
func Authentication(s *Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			token, err := r.Cookie("token")
			if err != nil {
				log.Println("Error occurred while reading cookie", err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			u, ok := s.UserRepositoryIsAuthenticatedByToken(token.Value)

			if ok {
				// Check if access token has expired
				if time.Now().After(u.LoggedInAt.Add(time.Second * time.Duration(604800))) {
					log.Println("after time check")
					err := s.UserRepositoryDeleteAccessToken(u.AccessToken)
					if err != nil {
						log.Println(err)
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}
					log.Println("Access token expired for ", u.Email)
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), User{}, &u)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				log.Println("Token does not exist")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		return http.HandlerFunc(fn)
	}

}

// HTTP middleware setting a value on the request context
func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(User{}).(User)

		if !user.IsAdmin {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
