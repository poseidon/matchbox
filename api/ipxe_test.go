package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPXEInspect(t *testing.T) {
	h := ipxeInspect()
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, ipxeBootstrap, w.Body.String())
}

func TestIPXEHandler(t *testing.T) {
	bootcfg := &BootConfig{
		Kernel: "/images/kernel",
		Initrd: []string{"/images/initrd_a", "/images/initrd_b"},
		Cmdline: map[string]interface{}{
			"a": "b",
			"c": "",
		},
	}
	expected := `#!ipxe
kernel /images/kernel a=b c
initrd /images/initrd_a /images/initrd_b 
boot
`
	store := &fixedStore{
		BootCfg: bootcfg,
	}
	h := ipxeHandler(store)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestIPXEHandler_MissingConfig(t *testing.T) {
	store := &emptyStore{}
	h := ipxeHandler(store)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
