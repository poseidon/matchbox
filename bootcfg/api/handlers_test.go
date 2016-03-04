package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireGET(t *testing.T) {
	next := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "next")
	}
	h := requireGET(http.HandlerFunc(next))
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "next", w.Body.String())
}

func TestRequireGET_WrongMethod(t *testing.T) {
	next := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "next")
	}
	h := requireGET(http.HandlerFunc(next))
	req, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "only HTTP GET is supported\n", w.Body.String())
}
