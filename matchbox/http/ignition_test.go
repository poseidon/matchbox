package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/server"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestIgnitionHandler_V2JSON(t *testing.T) {
	content := `{"ignition":{"version":"2.0.0","config":{}},"storage":{},"systemd":{"units":[{"name":"etcd2.service","enable":true}]},"networkd":{},"passwd":{}}`
	profile := &storagepb.Profile{
		Id:         fake.Group.Profile,
		IgnitionId: "file.ign",
	}
	store := &fake.FixedStore{
		Profiles:        map[string]*storagepb.Profile{fake.Group.Profile: profile},
		IgnitionConfigs: map[string]string{"file.ign": content},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: store})
	h := srv.ignitionHandler(c)
	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - raw Ignition config served directly
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, content, w.Body.String())
}

func TestIgnitionHandler_V2YAML(t *testing.T) {
	// exercise templating features, not a realistic Container Linux Config template
	content := `
systemd:
  units:
    - name: {{.service_name}}.service
      enable: true
    - name: {{.uuid}}.service
      enable: true
    - name: {{.request.query.foo}}.service
      enable: true
      contents: {{.request.raw_query}}
`
	expectedIgnitionV2 := `{"ignition":{"version":"2.0.0","config":{}},"storage":{},"systemd":{"units":[{"name":"etcd2.service","enable":true},{"name":"a1b2c3d4.service","enable":true},{"name":"some-param.service","enable":true,"contents":"foo=some-param\u0026bar=b"}]},"networkd":{},"passwd":{}}`
	store := &fake.FixedStore{
		Profiles:        map[string]*storagepb.Profile{fake.Group.Profile: testProfileIgnitionYAML},
		IgnitionConfigs: map[string]string{testProfileIgnitionYAML.IgnitionId: content},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: store})
	h := srv.ignitionHandler(c)
	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?foo=some-param&bar=b", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - Container Linux Config template rendered with Group selectors, metadata, and query variables
	// - Transformed to an Ignition config (JSON)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedIgnitionV2, w.Body.String())
}

func TestIgnitionHandler_MissingCtxProfile(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := srv.ignitionHandler(c)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIgnitionHandler_MissingIgnitionConfig(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := srv.ignitionHandler(c)
	ctx := withProfile(context.Background(), fake.Profile)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIgnitionHandler_MissingTemplateMetadata(t *testing.T) {
	content := `
ignition_version: 1
systemd:
  units:
    - name: {{.missing_key}}
      enable: true
`
	store := &fake.FixedStore{
		Profiles:        map[string]*storagepb.Profile{fake.Group.Profile: fake.Profile},
		IgnitionConfigs: map[string]string{fake.Profile.IgnitionId: content},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: store})
	h := srv.ignitionHandler(c)
	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - Ignition template rendering errors because "missing_key" is not
	// present in the template variables
	assert.Equal(t, http.StatusNotFound, w.Code)
}
