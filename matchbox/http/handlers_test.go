package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/coreos/coreos-baremetal/matchbox/server"
	"github.com/coreos/coreos-baremetal/matchbox/storage/storagepb"
	fake "github.com/coreos/coreos-baremetal/matchbox/storage/testfakes"
)

func TestSelectGroup(t *testing.T) {
	store := &fake.FixedStore{
		Groups: map[string]*storagepb.Group{fake.Group.Id: fake.Group},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: store})
	next := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, fake.Group, group)
		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - query params are used to match uuid=a1b2c3d4 to fake.Group
	// - the fake.Group is added to the context
	// - next handler is called
	h := srv.selectGroup(c, ContextHandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}

func TestSelectProfile(t *testing.T) {
	store := &fake.FixedStore{
		Groups:   map[string]*storagepb.Group{fake.Group.Id: fake.Group},
		Profiles: map[string]*storagepb.Profile{fake.Group.Profile: fake.Profile},
	}
	logger, _ := logtest.NewNullLogger()
	srv := NewServer(&Config{Logger: logger})
	c := server.NewServer(&server.Config{Store: store})
	next := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		profile, err := profileFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, fake.Profile, profile)
		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - query params are used to match uuid=a1b2c3d4 to fake.Group's fakeProfile
	// - the fake.Profile is added to the context
	// - next handler is called
	h := srv.selectProfile(c, ContextHandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}
