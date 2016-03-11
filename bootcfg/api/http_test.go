package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestLabelsFromRequest(t *testing.T) {
	emptyMap := map[string]string{}
	cases := []struct {
		urlString string
		labels    map[string]string
	}{
		{"http://a.io", emptyMap},
		{"http://a.io?uuid=a1b2c3", map[string]string{"uuid": "a1b2c3"}},
		{"http://a.io?uuid=a1b2c3", map[string]string{"uuid": "a1b2c3"}},
		{"http://a.io?mac=52:da:00:89:d8:10", map[string]string{"mac": validMACStr}},
		{"http://a.io?mac=52-da-00-89-d8-10", map[string]string{"mac": validMACStr}},
		{"http://a.io?uuid=a1b2c3&mac=52:da:00:89:d8:10", map[string]string{"uuid": "a1b2c3", "mac": validMACStr}},
		// parse and set MAC regardless of query argument case
		{"http://a.io?UUID=a1b2c3&MAC=52:DA:00:89:d8:10", map[string]string{"UUID": "a1b2c3", "MAC": validMACStr}},
		// ignore MAC addresses which do not parse
		{"http://a.io?mac=x:x:x:x:x:x", emptyMap},
	}
	for _, c := range cases {
		req, err := http.NewRequest("GET", c.urlString, nil)
		assert.Nil(t, err)
		assert.Equal(t, c.labels, labelsFromRequest(req))
	}
}

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
