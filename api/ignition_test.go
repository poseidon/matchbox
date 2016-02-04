package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var expectedIgnition = `{"ignitionVersion":1,"storage":{},"systemd":{"units":[{"name":"etcd2.service","enable":true}]},"networkd":{},"passwd":{}}`

func TestIgnitionHandler(t *testing.T) {
	content := `{"ignitionVersion": 1,"systemd":{"units":[{"name":"{{.service_name}}.service","enable":true}]}}`
	store := &fixedStore{
		Specs:           map[string]*Spec{testGroup.Spec: testSpec},
		IgnitionConfigs: map[string]string{testSpec.IgnitionConfig: content},
	}
	h := ignitionHandler(store)
	ctx := withGroup(context.Background(), &testGroup)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - Ignition template is rendered with Group metadata
	// - Rendered Ignition template is parsed as JSON
	// - Ignition Config served as JSON
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedIgnition, w.Body.String())
}

func TestIgnitionHandler_YAMLIgnition(t *testing.T) {
	content := `
ignition_version: 1
systemd:
  units:
    - name: {{.service_name}}.service
      enable: true
`
	store := &fixedStore{
		Specs:           map[string]*Spec{testGroup.Spec: testSpecWithIgnitionYAML},
		IgnitionConfigs: map[string]string{testSpecWithIgnitionYAML.IgnitionConfig: content},
	}
	h := ignitionHandler(store)
	ctx := withGroup(context.Background(), &testGroup)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - Ignition template is rendered with Group metadata
	// - Rendered Ignition template ending in .yaml is parsed as YAML
	// - Ignition Config served as JSON
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedIgnition, w.Body.String())
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
