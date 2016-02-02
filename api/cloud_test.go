package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestCloudHandler(t *testing.T) {
	cloudContent := "#cloud-config"
	store := &fixedStore{
		Specs:        map[string]*Spec{testGroup.Spec: testSpec},
		CloudConfigs: map[string]string{testSpec.CloudConfig: cloudContent},
	}
	h := cloudHandler(store)
	ctx := withGroup(context.Background(), &testGroup)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - the Spec's cloud config is served
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, cloudContent, w.Body.String())
}

func TestCloudHandler_MissingCtxSpec(t *testing.T) {
	h := cloudHandler(&emptyStore{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCloudHandler_MissingCloudConfig(t *testing.T) {
	h := cloudHandler(&emptyStore{})
	ctx := withSpec(context.Background(), testSpec)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
