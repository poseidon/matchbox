package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
	fake "github.com/coreos/coreos-baremetal/bootcfg/storage/testfakes"
)

var expectedIgnition = `{"ignitionVersion":1,"storage":{},"systemd":{"units":[{"name":"etcd2.service","enable":true}]},"networkd":{},"passwd":{}}`

func TestIgnitionHandler(t *testing.T) {
	content := `{"ignitionVersion": 1,"systemd":{"units":[{"name":"{{.service_name}}.service","enable":true}]}}`
	store := &fake.FixedStore{
		Profiles:        map[string]*storagepb.Profile{fake.Group.Profile: fake.Profile},
		IgnitionConfigs: map[string]string{fake.Profile.IgnitionId: content},
	}
	srv := server.NewServer(&server.Config{Store: store})
	h := ignitionHandler(srv)
	ctx := withGroup(context.Background(), fake.Group)
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
	store := &fake.FixedStore{
		Profiles:        map[string]*storagepb.Profile{fake.Group.Profile: testProfileIgnitionYAML},
		IgnitionConfigs: map[string]string{testProfileIgnitionYAML.IgnitionId: content},
	}
	srv := server.NewServer(&server.Config{Store: store})
	h := ignitionHandler(srv)
	ctx := withGroup(context.Background(), fake.Group)
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

func TestIgnitionHandler_MissingCtxProfile(t *testing.T) {
	srv := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := ignitionHandler(srv)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIgnitionHandler_MissingIgnitionConfig(t *testing.T) {
	srv := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := ignitionHandler(srv)
	ctx := withProfile(context.Background(), fake.Profile)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
