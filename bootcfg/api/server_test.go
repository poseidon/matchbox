package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerSpecRoute(t *testing.T) {
	store := &fixedStore{
		Specs: map[string]*Spec{"g1h2i3j4": testSpec},
	}
	h := NewServer(&Config{Store: store}).HTTPHandler()
	req, _ := http.NewRequest("GET", "/spec/g1h2i3j4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	// assert that:
	// - spec is rendered as JSON
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedSpecJSON, w.Body.String())
}

func TestServerSpecRoute_WrongMethod(t *testing.T) {
	h := NewServer(&Config{}).HTTPHandler()
	req, _ := http.NewRequest("POST", "/spec/g1h2i3j4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "only HTTP GET is supported\n", w.Body.String())
}
