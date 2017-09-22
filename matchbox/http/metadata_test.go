package http

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"context"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

func TestMetadataHandler(t *testing.T) {
	group := &storagepb.Group{
		Id:       "test-group",
		Selector: map[string]string{"mac": "52:54:00:a1:9c:ae"},
		Metadata: []byte(`{"meta":"data", "etcd":{"name":"node1"},"some":{"nested":{"data":"some-value"}}}`),
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	h := srv.metadataHandler()
	ctx := withGroup(context.Background(), group)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?mac=52-54-00-a1-9c-ae&foo=bar&count=3&gate=true", nil)
	h.ServeHTTP(w, req.WithContext(ctx))
	// assert that:
	// - Group selectors, metadata, and query variables are formatted
	// - nested metadata are namespaced
	// - key names are upper case
	// - key/value pairs are newline separated
	expectedLines := map[string]string{
		// group metadata
		"META":             "data",
		"ETCD_NAME":        "node1",
		"SOME_NESTED_DATA": "some-value",
		// group selector
		"MAC": "52:54:00:a1:9c:ae",
		// request
		"REQUEST_QUERY_MAC":   "52:54:00:a1:9c:ae",
		"REQUEST_QUERY_FOO":   "bar",
		"REQUEST_QUERY_COUNT": "3",
		"REQUEST_QUERY_GATE":  "true",
		"REQUEST_RAW_QUERY":   "mac=52-54-00-a1-9c-ae&foo=bar&count=3&gate=true",
	}
	assert.Equal(t, http.StatusOK, w.Code)
	// convert response (random order) to map (tests compare in order)
	assert.Equal(t, expectedLines, metadataToMap(w.Body.String()))
	assert.Equal(t, plainContentType, w.HeaderMap.Get(contentType))
}

func TestMetadataHandler_MetadataEdgeCases(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	h := srv.metadataHandler()
	// groups with different metadata
	cases := []struct {
		group    *storagepb.Group
		expected string
	}{
		{&storagepb.Group{Metadata: []byte(`{"num":3}`)}, "NUM=3\n"},
		{&storagepb.Group{Metadata: []byte(`{"yes":true}`)}, "YES=true\n"},
		{&storagepb.Group{Metadata: []byte(`{"no":false}`)}, "NO=false\n"},
	}
	for _, c := range cases {
		ctx := withGroup(context.Background(), c.group)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		h.ServeHTTP(w, req.WithContext(ctx))
		// assert that:
		// - Group metadata key names are upper case
		// - key/value pairs are newline separated
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), c.expected)
		assert.Equal(t, plainContentType, w.HeaderMap.Get(contentType))
	}
}

func TestMetadataHandler_MissingCtxGroup(t *testing.T) {
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	h := srv.metadataHandler()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// metadataToMap converts a KEY=val\nKEY=val ResponseWriter body to a map for
// testing purposes.
func metadataToMap(metadata string) map[string]string {
	scanner := bufio.NewScanner(strings.NewReader(metadata))
	data := make(map[string]string)
	for scanner.Scan() {
		token := scanner.Text()
		pair := strings.SplitN(token, "=", 2)
		if len(pair) != 2 {
			continue
		}
		data[pair[0]] = pair[1]
	}
	return data
}
