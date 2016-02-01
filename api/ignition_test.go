package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestIgnitionHandler(t *testing.T) {
	ignitioncfg := `{"ignitionVersion": 1}`
	store := &fixedStore{
		Specs:           map[string]*Spec{testGroup.Spec: testSpec},
		IgnitionConfigs: map[string]string{testSpec.IgnitionConfig: ignitioncfg},
	}
	h := ignitionHandler(store)
	ctx := withGroup(context.Background(), &testGroup)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - the Spec's ignition config is rendered
	expectedJSON := `{"ignitionVersion":1,"storage":{},"systemd":{},"networkd":{},"passwd":{}}`
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
