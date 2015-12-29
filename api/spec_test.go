package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// testSpec specifies a named group of configs.
	testSpec = &Spec{
		ID: "g1h2i3j4",
		BootConfig: &BootConfig{
			Kernel: "fake-kernel",
			Initrd: []string{"fake-initrd"},
			Cmdline: map[string]interface{}{
				"a": "b",
				"c": "",
			},
		},
	}
	expectedSpecJSON = `{"id":"g1h2i3j4","boot":{"kernel":"fake-kernel","initrd":["fake-initrd"],"cmdline":{"a":"b","c":""}}}`
)

func TestSpecHandler(t *testing.T) {
	store := &fixedStore{
		Specs: map[string]*Spec{"g1h2i3j4": testSpec},
	}
	h := specResource{store: store}
	req, _ := http.NewRequest("GET", "/g1h2i3j4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	// assert that:
	// - spec is rendered as JSON
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedSpecJSON, w.Body.String())
}

func TestSpecHandler_MissingConfig(t *testing.T) {
	store := &emptyStore{}
	h := specResource{store}
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
