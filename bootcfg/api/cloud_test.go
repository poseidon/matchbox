package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestCloudHandler(t *testing.T) {
	content := "#cloud-config"
	store := &fixedStore{
		Profiles:     map[string]*storagepb.Profile{testGroup.Profile: testProfile},
		CloudConfigs: map[string]string{testProfile.CloudId: content},
	}
	h := cloudHandler(store)
	ctx := withGroup(context.Background(), testGroup)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - Cloud config is rendered with Group metadata
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, content, w.Body.String())
}

func TestCloudHandler_MissingCtxProfile(t *testing.T) {
	h := cloudHandler(&emptyStore{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCloudHandler_MissingCloudConfig(t *testing.T) {
	h := cloudHandler(&emptyStore{})
	ctx := withProfile(context.Background(), testProfile)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
