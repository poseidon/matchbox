package api

import (
	"errors"
	"net/http"
	"sort"

	"github.com/coreos/coreos-baremetal/bootcfg/storage"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
	"golang.org/x/net/context"
)

var (
	errNoMatchingGroup = errors.New("api: No matching Group")
)

type groupsResource struct {
	store storage.Store
}

func newGroupsResource(store storage.Store) *groupsResource {
	return &groupsResource{
		store: store,
	}
}

// matchGroupHandler returns a ContextHandler that matches machine requests to
// a Group, adds the Group to the ctx, and calls the next handler. The next
// handler should handle the case that no matching Group is found.
func (gr *groupsResource) matchGroupHandler(next ContextHandler) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		attrs := labelsFromRequest(req)
		// match machine request
		group, err := gr.findMatch(attrs)
		if err == nil {
			// add the Group to the ctx for next handler
			ctx = withGroup(ctx, group)
		}
		next.ServeHTTP(ctx, w, req)
	}
	return ContextHandlerFunc(fn)
}

// matchProfileHandler returns a ContextHandler that matches machine requests
// to a Profile, adds Profile to the ctx, and calls the next handler. The
// next handler should handle the case that no matching Profile is found.
func (gr *groupsResource) matchProfileHandler(next ContextHandler) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		attrs := labelsFromRequest(req)
		// match machine request
		group, err := gr.findMatch(attrs)
		if err == nil {
			// lookup Profile by id
			profile, err := gr.store.ProfileGet(group.Profile)
			if err == nil {
				// add the Profile to the ctx for next handler
				ctx = withProfile(ctx, profile)
			}
		}
		next.ServeHTTP(ctx, w, req)
	}
	return ContextHandlerFunc(fn)
}

// findMatch returns the first Group whose Matcher is satisfied by the given
// labels. Groups are attempted in sorted order, preferring those with
// more matcher conditions, alphabetically.
func (gr *groupsResource) findMatch(labels map[string]string) (*storagepb.Group, error) {
	groups, err := gr.store.GroupList()
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(storagepb.ByReqs(groups)))
	for _, group := range groups {
		if group.Matches(labels) {
			return group, nil
		}
	}
	return nil, errNoMatchingGroup
}
