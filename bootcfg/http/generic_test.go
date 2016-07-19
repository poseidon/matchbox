package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/mikeynap/coreos-baremetal/bootcfg/server"
	"github.com/mikeynap/coreos-baremetal/bootcfg/storage/storagepb"
	fake "github.com/mikeynap/coreos-baremetal/bootcfg/storage/testfakes"
)

func TestGenericHandler(t *testing.T) {
	content := `#foo-bar-baz template
UUID={{.uuid}}
SERVICE={{.service_name}}
`
	expected := `#foo-bar-baz template
UUID=a1b2c3d4
SERVICE=etcd2
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
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - Generic config is rendered with Group metadata and selectors
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
	h.ServeHTTP(context.Background(), w, req)
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
	h.ServeHTTP(ctx, w, req)
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
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - Generic template rendering errors because "missing_key" is not
	// present in the Group metadata
	assert.Equal(t, http.StatusNotFound, w.Code)
}
