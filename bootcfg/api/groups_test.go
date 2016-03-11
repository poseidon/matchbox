package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/coreos-baremetal/bootcfg/storage"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNewGroupsResource(t *testing.T) {
	store := &fixedStore{}
	gr := newGroupsResource(store)
	assert.Equal(t, store, gr.store)
}

func TestGroupsResource_MatchProfileHandler(t *testing.T) {
	store := &fixedStore{
		Groups:   map[string]*storagepb.Group{testGroup.Id: testGroup},
		Profiles: map[string]*storagepb.Profile{testGroup.Profile: testProfile},
	}
	gr := newGroupsResource(store)
	next := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		profile, err := profileFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, testProfile, profile)
		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - request arguments are used to match uuid=a1b2c3d4 -> testGroup
	// - the group's Profile is found by id and added to the context
	// - next handler is called
	h := gr.matchProfileHandler(ContextHandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}

func TestGroupsResource_MatchGroupHandler(t *testing.T) {
	store := &fixedStore{
		Groups: map[string]*storagepb.Group{testGroup.Id: testGroup},
	}
	gr := newGroupsResource(store)
	next := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, testGroup, group)
		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - request arguments are used to match uuid=a1b2c3d4 -> testGroup
	// - next handler is called
	h := gr.matchGroupHandler(ContextHandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}

func TestGroupsResource_FindMatch(t *testing.T) {
	store := &fixedStore{
		Groups: map[string]*storagepb.Group{testGroup.Id: testGroup},
	}
	cases := []struct {
		store         storage.Store
		labels        map[string]string
		expectedGroup *storagepb.Group
		expectedErr   error
	}{
		{store, map[string]string{"uuid": "a1b2c3d4"}, testGroup, nil},
		{store, nil, nil, errNoMatchingGroup},
		// no groups in the store
		{&emptyStore{}, map[string]string{"a": "b"}, nil, errNoMatchingGroup},
	}

	for _, c := range cases {
		gr := newGroupsResource(c.store)
		group, err := gr.findMatch(c.labels)
		assert.Equal(t, c.expectedGroup, group)
		assert.Equal(t, c.expectedErr, err)
	}
}
