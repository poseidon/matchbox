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

func TestGenericHandler(t *testing.T) {
	content := `#foo-bar-baz template
UUID={{.uuid}}
SERVICE={{.service_name}}
FOO={{.request.query.foo}}
`
	expected := `#foo-bar-baz template
UUID=a1b2c3d4
SERVICE=etcd2
FOO=some-param
`
	store := &fake.FixedStore{
		Profiles:       map[string]*storagepb.Profile{fake.Group.Profile: fake.Profile},
		GenericConfigs: map[string]string{fake.Profile.GenericId: content},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: store})
	h := srv.genericHandler(c)
	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?foo=some-param", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - Generic config is rendered with Group selectors, metadata, and query variables
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestGenericHandler_MissingCtxProfile(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := srv.genericHandler(c)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGenericHandler_MissingCloudConfig(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := srv.genericHandler(c)
	ctx := withProfile(context.Background(), fake.Profile)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGenericHandler_MissingTemplateMetadata(t *testing.T) {
	content := `#foo-bar-baz template
KEY={{.missing_key}}
`
	store := &fake.FixedStore{
		Profiles:       map[string]*storagepb.Profile{fake.Group.Profile: fake.Profile},
		GenericConfigs: map[string]string{fake.Profile.GenericId: content},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: store})
	h := srv.cloudHandler(c)
	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - Generic template rendering errors because "missing_key" is not
	// present in the template variables
	assert.Equal(t, http.StatusNotFound, w.Code)
}
