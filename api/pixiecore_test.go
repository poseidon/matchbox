package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validMAC = "52:54:00:89:d8:10"

func TestPixiecoreHandler(t *testing.T) {
	bootcfg := &BootConfig{
		Kernel: "/images/kernel",
		Initrd: []string{"/images/initrd_a", "/images/initrd_b"},
		Cmdline: map[string]interface{}{
			"a": "b",
			"c": "",
		},
	}
	store := &fixedStore{BootCfg: bootcfg}
	h := pixiecoreHandler(store)
	req, _ := http.NewRequest("GET", "/"+validMAC, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	expectedJSON := `{"kernel":"/images/kernel","initrd":["/images/initrd_a","/images/initrd_b"],"cmdline":{"a":"b","c":""}}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedJSON, w.Body.String())
}

func TestPixiecoreHandler_InvalidMACAddress(t *testing.T) {
	store := &emptyStore{}
	h := pixiecoreHandler(store)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid MAC address /\n", w.Body.String())
}

func TestPixiecoreHandler_MissingConfig(t *testing.T) {
	store := &emptyStore{}
	h := pixiecoreHandler(store)
	req, _ := http.NewRequest("GET", "/"+validMAC, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
