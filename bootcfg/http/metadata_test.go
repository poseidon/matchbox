package http

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
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
	req, _ := http.NewRequest("GET", "/?mac=52-54-00-a1-9c-ae", nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - the Group's custom metadata and selectors are served
	// - key names are upper case
	expectedData := map[string]string{
		// group metadata
		"META": "data",
		"ETCD": "map[name:node1]",
		"SOME": "map[nested:map[data:some-value]]",
		// group selector
		"MAC": "52:54:00:a1:9c:ae",
		// HACK(dghubble): Not testing query params until #84
	}
	assert.Equal(t, http.StatusOK, w.Code)
	// convert response (random order) to map (tests compare in order)
	assert.Equal(t, expectedData, metadataToMap(w.Body.String()))
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
		// Issue #84 - improve list and map printouts
		{&storagepb.Group{Metadata: []byte(`{"list":["3","d"]}`)}, "LIST=[3 d]\n"},
	}
	for _, c := range cases {
		ctx := withGroup(context.Background(), c.group)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		h.ServeHTTP(ctx, w, req)
		// assert that each Group's metadata is formatted:
		// - key names are upper case
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
	h.ServeHTTP(context.Background(), w, req)
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
		// HACK(dghubble) - Skip map unwinding until #84
		if pair[0] == "REQUEST" {
			continue
		}
		data[pair[0]] = pair[1]
	}
	return data
}
