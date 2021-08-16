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
			token       UserIDContext
		)

		if authHeader = r.Header.Get("Authorization"); authHeader == "" {
			w.WriteHeader(401)
			return
		}

		if bearerValue = strings.Split(authHeader, " "); len(bearerValue) != 2 {
			w.WriteHeader(401)
			return
		}

		token = UserIDContext(bearerValue[1])

		if userID, err = s.GetAccessToken(token.String()); err != nil {
			w.WriteHeader(401)
			return
		}

		ctx := context.WithValue(r.Context(), token, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)

}
