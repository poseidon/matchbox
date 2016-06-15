package http

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
	srv := server.NewServer(&server.Config{Store: store})
	h := genericHandler(srv)
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
	srv := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := genericHandler(srv)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGenericHandler_MissingCloudConfig(t *testing.T) {
	srv := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := genericHandler(srv)
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
	srv := server.NewServer(&server.Config{Store: store})
	h := cloudHandler(srv)
	ctx := withGroup(context.Background(), fake.Group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - Generic template rendering errors because "missing_key" is not
	// present in the Group metadata
	assert.Equal(t, http.StatusNotFound, w.Code)
}
