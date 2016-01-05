package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ignition "github.com/coreos/ignition/src/config"
	"github.com/stretchr/testify/assert"
)

func TestIgnitionHandler(t *testing.T) {
	ignitioncfg := &ignition.Config{}
	store := &fixedStore{
		Machines:        map[string]*Machine{"a1b2c3d4": testMachine},
		IgnitionConfigs: map[string]*ignition.Config{"ignition.json": ignitioncfg},
	}
	h := ignitionHandler(store)
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	// assert that:
	// - match parameters to a Spec
	// - render the Spec's ignition config
	expectedJSON := `{"ignitionVersion":0,"storage":{},"systemd":{},"networkd":{},"passwd":{}}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedJSON, w.Body.String())
}

func TestIgnitionHandler_NoMatchingSpec(t *testing.T) {
	store := &emptyStore{}
	h := ignitionHandler(store)
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIgnitionHandler_MissingIgnitionConfig(t *testing.T) {
	store := &fixedStore{
		Machines: map[string]*Machine{"a1b2c3d4": testMachine},
	}
	h := ignitionHandler(store)
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
