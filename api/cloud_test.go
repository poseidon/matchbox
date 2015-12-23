package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudHandler(t *testing.T) {
	cloudcfg := &CloudConfig{
		Content: "#cloud-config",
	}
	store := &fixedStore{
		CloudCfg: cloudcfg,
	}
	h := cloudHandler(store)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, cloudcfg.Content, w.Body.String())
}

func TestCloudHandler_MissingConfig(t *testing.T) {
	store := &emptyStore{}
	h := cloudHandler(store)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
