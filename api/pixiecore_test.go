package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPixiecoreHandler(t *testing.T) {
	store := &fixedStore{
		Groups: []Group{testGroupWithMAC},
		Specs:  map[string]*Spec{testGroupWithMAC.Spec: testSpec},
	}
	h := pixiecoreHandler(newGroupsResource(store), store)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/"+validMACStr, nil)
	h.ServeHTTP(w, req)
	// assert that:
	// - MAC address argument is used for Spec matching
	// - the Spec's boot config is rendered as Pixiecore JSON
	expectedJSON := `{"kernel":"/image/kernel","initrd":["/image/initrd_a","/image/initrd_b"],"cmdline":{"a":"b","c":""}}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedJSON, w.Body.String())
}

func TestPixiecoreHandler_InvalidMACAddress(t *testing.T) {
	h := pixiecoreHandler(&groupsResource{}, &emptyStore{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid MAC address /\n", w.Body.String())
}

func TestPixiecoreHandler_NoMatchingGroup(t *testing.T) {
	h := pixiecoreHandler(newGroupsResource(&emptyStore{}), &emptyStore{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/"+validMACStr, nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPixiecoreHandler_NoMatchingSpec(t *testing.T) {
	store := &fixedStore{
		Groups: []Group{testGroupWithMAC},
	}
	h := pixiecoreHandler(newGroupsResource(store), &emptyStore{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/"+validMACStr, nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
