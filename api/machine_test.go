package api

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// testMachine specifies its configuration.
	testMachine = &Machine{
		ID: "a1b2c3d4",
		BootConfig: &BootConfig{
			Kernel: "fake-kernel",
			Initrd: []string{"fake-initrd"},
			Cmdline: map[string]interface{}{
				"a": "b",
				"c": "",
			},
		},
	}
	// testSharedSpecMachine references a Spec configuration.
	testSharedSpecMachine = &Machine{
		ID:     "a1b2c3d4",
		SpecID: "g1h2i3j4",
	}
	expectedMachineJSON           = `{"id":"a1b2c3d4","boot":{"kernel":"fake-kernel","initrd":["fake-initrd"],"cmdline":{"a":"b","c":""}},"spec_id":""}`
	expectedSharedSpecMachineJSON = `{"id":"a1b2c3d4","boot":{"kernel":"fake-kernel","initrd":["fake-initrd"],"cmdline":{"a":"b","c":""}},"spec_id":"g1h2i3j4"}`
)

func TestMachineHandler(t *testing.T) {
	store := &fixedStore{
		Machines: map[string]*Machine{"a1b2c3d4": testMachine},
	}
	h := machineResource{store: store}
	req, _ := http.NewRequest("GET", "/a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	// assert that:
	// - machine config is rendered as JSON
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedMachineJSON, w.Body.String())
}

func TestMachineHandler_SharedSpec(t *testing.T) {
	store := &fixedStore{
		Machines: map[string]*Machine{"a1b2c3d4": testSharedSpecMachine},
		Specs:    map[string]*Spec{"g1h2i3j4": testSpec},
	}
	h := machineResource{store: store}
	req, _ := http.NewRequest("GET", "/a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	// assert that:
	// - machine config is rendered as JSON
	// - boot config values from the spec are merged into the JSON response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedSharedSpecMachineJSON, w.Body.String())
}

func TestMachineHandler_MissingConfig(t *testing.T) {
	store := &emptyStore{}
	h := machineResource{store}
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAttrsFromRequest(t *testing.T) {
	hwAddr, err := net.ParseMAC("52:54:00:89:d8:10")
	assert.Nil(t, err)
	cases := []struct {
		urlString string
		attrs     MachineAttrs
	}{
		{"http://a.io", MachineAttrs{}},
		{"http://a.io?uuid=a1b2c3", MachineAttrs{UUID: "a1b2c3"}},
		{"http://a.io?mac=52:54:00:89:d8:10", MachineAttrs{MAC: hwAddr}},
		{"http://a.io?mac=52-54-00-89-d8-10", MachineAttrs{MAC: hwAddr}},
		{"http://a.io?uuid=a1b2c3&mac=52:54:00:89:d8:10", MachineAttrs{UUID: "a1b2c3", MAC: hwAddr}},
		// leave MAC nil if it does not parse
		{"http://a.io?mac=?:?:?", MachineAttrs{}},
	}
	for _, c := range cases {
		req, err := http.NewRequest("GET", c.urlString, nil)
		assert.Nil(t, err)
		assert.Equal(t, attrsFromRequest(req), c.attrs)
	}
}
