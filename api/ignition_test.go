package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ignition "github.com/coreos/ignition/src/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestIgnitionHandler(t *testing.T) {
	ignitioncfg := &ignition.Config{}
	store := &fixedStore{
		IgnitionConfigs: map[string]*ignition.Config{testSpec.IgnitionConfig: ignitioncfg},
	}
	h := ignitionHandler(store)
	ctx := withSpec(context.Background(), testSpec)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - the Spec's ignition config is rendered
	expectedJSON := `{"ignitionVersion":0,"storage":{},"systemd":{},"networkd":{},"passwd":{}}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedJSON, w.Body.String())
}

func TestIgnitionHandler_MissingCtxSpec(t *testing.T) {
	h := ignitionHandler(&emptyStore{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIgnitionHandler_MissingIgnitionConfig(t *testing.T) {
	h := ignitionHandler(&emptyStore{})
	ctx := withSpec(context.Background(), testSpec)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
