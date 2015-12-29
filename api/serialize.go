package api

import (
	"encoding/json"
	"net/http"
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json"
)

// renderJSON encodes structs to JSON, writes the response to the
// ResponseWriter, and logs encoding errors.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		log.Errorf("error JSON encoding: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonContentType)
	_, err = w.Write(js)
	if err != nil {
		log.Errorf("error writing to response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
