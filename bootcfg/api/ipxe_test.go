package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestIPXEInspect(t *testing.T) {
	h := ipxeInspect()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, ipxeBootstrap, w.Body.String())
}

func TestIPXEHandler(t *testing.T) {
	h := ipxeHandler()
	ctx := withSpec(context.Background(), testSpec)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - the Spec's boot config is rendered as an iPXE script
	expectedScript := `#!ipxe
kernel /image/kernel a=b c
initrd /image/initrd_a /image/initrd_b 
boot
`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedScript, w.Body.String())
}

func TestIPXEHandler_MissingCtxSpec(t *testing.T) {
	h := ipxeHandler()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIPXEHandler_RenderTemplateError(t *testing.T) {
	h := ipxeHandler()
	// a Spec with nil BootConfig forces a template.Execute error
	ctx := withSpec(context.Background(), &Spec{BootConfig: nil})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestIPXEHandler_WriteError(t *testing.T) {
	h := ipxeHandler()
	ctx := withSpec(context.Background(), testSpec)
	w := NewUnwriteableResponseWriter()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(ctx, w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Empty(t, w.Body.String())
}
