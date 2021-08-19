package gobookie

import (
	"context"
	"net/http"
	"strings"
)

// Authentication Only allow authenticated to query API
// Checks if a given bearer token exists in DB
func (s *Server) Authentication(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			authHeader  string
			bearerValue []string
			userID      int
			err         error
			token       string
		)

		if authHeader = r.Header.Get("Authorization"); authHeader == "" {
			w.WriteHeader(401)
			return
		}

		if bearerValue = strings.Split(authHeader, " "); len(bearerValue) != 2 {
			w.WriteHeader(401)
			return
		}

		token = bearerValue[1]

		// Should probably check if token is valid as well :)
		if userID, err = s.GetAccessToken(token); err != nil {
			w.WriteHeader(401)
			return
		}

		rq := &UserRequestData{
			Token:  token,
			UserID: userID,
		}

		ctx := context.WithValue(r.Context(), UserData(UserRequestData{}), rq)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)

}
