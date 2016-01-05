package api

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// testSpec specifies a named group of configs for testing purposes.
	testSpec = &Spec{
		ID: "g1h2i3j4",
		BootConfig: &BootConfig{
			Kernel: "/image/kernel",
			Initrd: []string{"/image/initrd_a", "/image/initrd_b"},
			Cmdline: map[string]interface{}{
				"a": "b",
				"c": "",
			},
		},
		CloudConfig:    "cloud-config.yml",
		IgnitionConfig: "ignition.json",
	}
	emptySpec = &Spec{
		ID: "empty",
	}
	expectedSpecJSON = `{"id":"g1h2i3j4","boot":{"kernel":"/image/kernel","initrd":["/image/initrd_a","/image/initrd_b"],"cmdline":{"a":"b","c":""}},"cloud_id":"cloud-config.yml","ignition_id":"ignition.json"}`
)

func TestSpecHandler(t *testing.T) {
	store := &fixedStore{
		Specs: map[string]*Spec{"g1h2i3j4": testSpec},
	}
	h := specResource{store: store}
	req, _ := http.NewRequest("GET", "/g1h2i3j4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	// assert that:
	// - spec is rendered as JSON
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedSpecJSON, w.Body.String())
}

func TestSpecHandler_MissingConfig(t *testing.T) {
	store := &emptyStore{}
	h := specResource{store}
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetMatchingSpec_ByUUID(t *testing.T) {
	hwAddr, err := net.ParseMAC(validMACStr)
	assert.Nil(t, err)
	store := &fixedStore{
		Machines: map[string]*Machine{
			"a1b2c3d4":      testMachine,
			hwAddr.String(): testMachineEmptySpec,
		},
		Specs: map[string]*Spec{"g1h2i3j4": testSpec},
	}
	spec, err := getMatchingSpec(store, MachineAttrs{UUID: "a1b2c3d4", MAC: hwAddr})
	// assert that:
	// - attributes match the testMachine's Spec
	// - UUID attribute takes priority over MAC address
	assert.Nil(t, err)
	assert.Equal(t, testSpec, spec)
}

func TestGetMatchingSpec_ByMAC(t *testing.T) {
	hwAddr, err := net.ParseMAC(validMACStr)
	assert.Nil(t, err)
	store := &fixedStore{
		Machines: map[string]*Machine{
			hwAddr.String(): testMachine,
			"default":       testMachineEmptySpec,
		},
		Specs: map[string]*Spec{"g1h2i3j4": testSpec},
	}
	spec, err := getMatchingSpec(store, MachineAttrs{MAC: hwAddr})
	// assert that:
	// - attributes match the testMachine's Spec
	// - MAC address attribute takes priority over default
	assert.Nil(t, err)
	assert.Equal(t, testSpec, spec)
}

func TestGetMatchingSpec_Default(t *testing.T) {
	store := &fixedStore{
		Machines: map[string]*Machine{"default": testMachine},
		Specs:    map[string]*Spec{"g1h2i3j4": testSpec},
	}
	spec, err := getMatchingSpec(store, MachineAttrs{UUID: "any-uuid"})
	assert.Nil(t, err)
	assert.Equal(t, testSpec, spec)
}

func TestGetMatchingSpec_NoMatchingSpec(t *testing.T) {
	store := &emptyStore{}
	spec, err := getMatchingSpec(store, MachineAttrs{})
	assert.Nil(t, spec)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), "no spec matching { }")
	}
}
