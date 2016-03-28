package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

func TestRequireGET(t *testing.T) {
	next := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "next")
	}
	h := requireGET(http.HandlerFunc(next))
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "next", w.Body.String())
}

func TestRequireGET_WrongMethod(t *testing.T) {
	next := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "next")
	}
	h := requireGET(http.HandlerFunc(next))
	req, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "only HTTP GET is supported\n", w.Body.String())
}

func TestSelectGroup(t *testing.T) {
	store := &fixedStore{
		Groups: map[string]*storagepb.Group{testGroup.Id: testGroup},
	}
	srv := server.NewServer(&server.Config{Store: store})
	next := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, testGroup, group)
		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - query params are used to match uuid=a1b2c3d4 to testGroup
	// - the testGroup is added to the context
	// - next handler is called
	h := selectGroup(srv, ContextHandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}

func TestSelectProfile(t *testing.T) {
	store := &fixedStore{
		Groups:   map[string]*storagepb.Group{testGroup.Id: testGroup},
		Profiles: map[string]*storagepb.Profile{testGroup.Profile: testProfile},
	}
	srv := server.NewServer(&server.Config{Store: store})
	next := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		profile, err := profileFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, testProfile, profile)
		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - query params are used to match uuid=a1b2c3d4 to testGroup's testProfile
	// - the testProfile is added to the context
	// - next handler is called
	h := selectProfile(srv, ContextHandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}
