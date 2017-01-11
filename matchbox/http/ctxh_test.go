package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNewHandler(t *testing.T) {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ContextHandler called")
	}
	h := NewHandler(ContextHandlerFunc(fn))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, "ContextHandler called", w.Body.String())
}
