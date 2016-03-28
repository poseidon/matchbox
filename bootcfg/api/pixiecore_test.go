package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
	fake "github.com/coreos/coreos-baremetal/bootcfg/storage/testfakes"
)

func TestPixiecoreHandler(t *testing.T) {
	store := &fake.FixedStore{
		Groups:   map[string]*storagepb.Group{testGroupWithMAC.Id: testGroupWithMAC},
		Profiles: map[string]*storagepb.Profile{testGroupWithMAC.Profile: testProfile},
	}
	srv := server.NewServer(&server.Config{Store: store})
	h := pixiecoreHandler(srv)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/"+validMACStr, nil)
	h.ServeHTTP(context.Background(), w, req)
	// assert that:
	// - MAC address parameter is used for Group matching
	// - the Profile's NetBoot config is rendered as Pixiecore JSON
	expectedJSON := `{"kernel":"/image/kernel","initrd":["/image/initrd_a","/image/initrd_b"],"cmdline":{"a":"b","c":""}}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonContentType, w.HeaderMap.Get(contentType))
	assert.Equal(t, expectedJSON, w.Body.String())
}

func TestPixiecoreHandler_InvalidMACAddress(t *testing.T) {
	srv := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := pixiecoreHandler(srv)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid MAC address /\n", w.Body.String())
}

func TestPixiecoreHandler_NoMatchingGroup(t *testing.T) {
	srv := server.NewServer(&server.Config{Store: &fake.EmptyStore{}})
	h := pixiecoreHandler(srv)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/"+validMACStr, nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPixiecoreHandler_NoMatchingProfile(t *testing.T) {
	store := &fake.FixedStore{
		Groups: map[string]*storagepb.Group{testGroup.Id: testGroup},
	}
	srv := server.NewServer(&server.Config{Store: store})
	h := pixiecoreHandler(srv)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/"+validMACStr, nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
