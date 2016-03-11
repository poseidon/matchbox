package api

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestMetadataHandler(t *testing.T) {
	h := metadataHandler()
	ctx := withGroup(context.Background(), testGroup)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?uuid=a1b2c3d4&mac="+validMACStr, nil)
	h.ServeHTTP(ctx, w, req)
	// assert that:
	// - the Group's custom metadata is
	// - query argument attributes are added to the metadata
	// - key names are upper case
	expectedData := map[string]string{
		"K8S_VERSION":  "v1.1.2",
		"POD_NETWORK":  "10.2.0.0/16",
		"SERVICE_NAME": "etcd2",
		"UUID":         "a1b2c3d4",
		"MAC":          validMACStr,
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedData, metadataToMap(w.Body.String()))
	assert.Equal(t, plainContentType, w.HeaderMap.Get(contentType))
}

func TestMetadataHandler_MissingCtxGroup(t *testing.T) {
	h := metadataHandler()
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
		data[pair[0]] = pair[1]
	}
	return data
}
