package api

import (
	"net/http"
)

// requireGET requires requests to be an HTTP GET. Otherwise, it responds with
// a 405 status code.
func requireGET(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			http.Error(w, "only HTTP GET is supported", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// logRequests logs HTTP requests.
func logRequests(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		log.Debugf("HTTP %s %v", req.Method, req.URL)
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
