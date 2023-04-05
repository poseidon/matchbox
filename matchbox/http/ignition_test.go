package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/poseidon/matchbox/matchbox/server"
	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
	fake "github.com/poseidon/matchbox/matchbox/storage/testfakes"
)

func TestIgnitionHandler_v3_4(t *testing.T) {
	const content = `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.4.0"},"kernelArguments":{},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":false,"name":"docker.service"}]}}`
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
	core := server.NewServer(&server.Config{Store: store})
	h := srv.ignitionHandler(core)

	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	h.ServeHTTP(w, req)
	// assert that:
	// - serve Ignition JSON
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.Header().Get(contentType))
	assert.Equal(t, content, w.Body.String())
}

func TestIgnitionHandler_v3_1(t *testing.T) {
	const ign31 = `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.1.0"},"kernelArguments":{},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":false,"name":"docker.service"}]}}`
	const ign34 = `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.4.0"},"kernelArguments":{},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":false,"name":"docker.service"}]}}`
	profile := &storagepb.Profile{
		Id:         fake.Group.Profile,
		IgnitionId: "file.ign",
	}
	store := &fake.FixedStore{
		Profiles: map[string]*storagepb.Profile{
			fake.Group.Profile: profile,
		},
		IgnitionConfigs: map[string]string{
			"file.ign": ign31,
		},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	core := server.NewServer(&server.Config{Store: store})
	h := srv.ignitionHandler(core)

	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	h.ServeHTTP(w, req)
	// assert that:
	// - older Ignition v3.x converted to compatible latest version (e.g. v3.3)
	// - serve Ignition JSON
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.Header().Get(contentType))
	assert.Equal(t, ign34, w.Body.String())
}

func TestIgnitionHandler_MissingIgnition(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	core := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := srv.ignitionHandler(core)

	ctx := withProfile(context.Background(), fake.Profile)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIgnitionHandler_Butane(t *testing.T) {
	// exercise templating features, not a realistic Butane template
	butane := `
variant: flatcar
version: 1.1.0
systemd:
  units:
    - name: {{.uuid}}.service
      enabled: true
      contents: {{.pod_network}}
    - name: {{.request.query.foo}}.service
      enabled: true
      contents: {{.request.raw_query}}
`
	expectedIgnition := `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.4.0"},"kernelArguments":{},"passwd":{},"storage":{},"systemd":{"units":[{"contents":"10.2.0.0/16","enabled":true,"name":"a1b2c3d4.service"},{"contents":"foo=some-param\u0026bar=b","enabled":true,"name":"some-param.service"}]}}`

	store := &fake.FixedStore{
		Profiles: map[string]*storagepb.Profile{
			fake.Group.Profile: testProfileWithButane,
		},
		IgnitionConfigs: map[string]string{
			testProfileWithButane.IgnitionId: butane,
		},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	core := server.NewServer(&server.Config{Store: store})
	h := srv.ignitionHandler(core)

	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/?foo=some-param&bar=b", nil)
	h.ServeHTTP(w, req)
	// assert that:
	// - Template rendered with Group selectors, metadata, and query variables
	// - Butane translated to an Ignition config
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.Header().Get(contentType))
	assert.Equal(t, expectedIgnition, w.Body.String())
}

func TestIgnitionHandler_MissingCtxProfile(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	core := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := srv.ignitionHandler(core)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIgnitionHandler_MissingTemplateMetadata(t *testing.T) {
	butane := `
variant: flatcar
version: 1.1.0
systemd:
  units:
    - name: {{.missing_key}}
      enabled: true
`
	store := &fake.FixedStore{
		Profiles: map[string]*storagepb.Profile{
			fake.Group.Profile: fake.Profile,
		},
		IgnitionConfigs: map[string]string{
			fake.Profile.IgnitionId: butane,
		},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	core := server.NewServer(&server.Config{Store: store})
	h := srv.ignitionHandler(core)

	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	h.ServeHTTP(w, req)
	// assert that:
	// - Template rendering errors because "missing_key" is not present
	assert.Equal(t, http.StatusNotFound, w.Code)
}
