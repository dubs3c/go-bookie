package gobookie

import (
	"encoding/json"
	"net/http"
)

// RespondWithError - Return an error
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"error": msg})
}

// RespondWithJSON - Respond with a json formatted string
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return
}

// RespondWithStatusCode - Respond with a status code without setting a message
func RespondWithStatusCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	return
}
