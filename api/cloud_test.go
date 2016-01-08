package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudHandler(t *testing.T) {
	cloudcfg := &CloudConfig{
		Content: "#cloud-config",
	}
	store := &fixedStore{
		Groups:       []Group{testGroup},
		Specs:        map[string]*Spec{testGroup.Spec: testSpec},
		CloudConfigs: map[string]*CloudConfig{testSpec.CloudConfig: cloudcfg},
	}
	h := cloudHandler(store)
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	// assert that:
	// - match parameters to a Spec
	// - render the Spec's cloud config
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, cloudcfg.Content, w.Body.String())
}

func TestCloudHandler_NoMatchingSpec(t *testing.T) {
	store := &emptyStore{}
	h := cloudHandler(store)
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCloudHandler_MissingCloudConfig(t *testing.T) {
	store := &fixedStore{
		Groups: []Group{testGroup},
		Specs:  map[string]*Spec{testGroup.Spec: testSpec},
	}
	h := cloudHandler(store)
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
