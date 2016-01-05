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
	store := &fixedStore{
		Machines: map[string]*Machine{"a1b2c3d4": testMachine},
	}
	h := ipxeHandler(store)
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	// assert that:
	// - machine config is rendered as an iPXE script
	expectedScript := `#!ipxe
kernel /image/kernel a=b c
initrd /image/initrd_a /image/initrd_b 
boot
`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedScript, w.Body.String())
}

func TestIPXEHandler_NoMatchingSpec(t *testing.T) {
	store := &emptyStore{}
	h := ipxeHandler(store)
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIPXEHandler_RenderTemplateError(t *testing.T) {
	// machine spec has a nil BootConfig so template.Execute errors
	store := &fixedStore{
		Machines: map[string]*Machine{
			"a1b2c3d4": &Machine{
				ID:   "a1b2c3d4",
				Spec: &Spec{BootConfig: nil},
			},
		},
	}
	h := ipxeHandler(store)
	req, _ := http.NewRequest("GET", "/?uuid=a1b2c3d4", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIPXEHandler_WriteError(t *testing.T) {
	store := &fixedStore{
		Machines: map[string]*Machine{"a1b2c3d4": testMachine},
	}
	h := ipxeHandler(store)
	req, _ := http.NewRequest("GET", "/?uuid=a1b2c3d4", nil)
	w := NewUnwriteableResponseWriter()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Empty(t, w.Body.String())
}
