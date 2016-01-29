package api

import (
	"errors"
	"net/http"
	"sort"

	"golang.org/x/net/context"
)

var (
	errNoMatchingGroup = errors.New("api: No matching Group")
)

// Group associates matcher conditions with a Specification identifier. The
// Matcher conditions may be satisfied by zero or more machines.
type Group struct {
	// Human readable name (optional)
	Name string `yaml:"name"`
	// Spec identifier
	Spec string `yaml:"spec"`
	// Custom Metadata
	Metadata map[string]string `yaml:"metadata"`
	// matcher conditions
	Matcher RequirementSet `yaml:"require"`
}

// byMatcher defines a collection of Group structs which have a deterministic
// order in increasing number of Matchers, then by sorted key/value pair strings.
// For example, Matcher {a:b, c:d} is ordered before {a:d, c:d} and {a:b, d:e}.
type byMatcher []Group

func (g byMatcher) Len() int {
	return len(g)
}

func (g byMatcher) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g byMatcher) Less(i, j int) bool {
	if len(g[i].Matcher) == len(g[j].Matcher) {
		return g[i].Matcher.String() < g[j].Matcher.String()
	}
	return len(g[i].Matcher) < len(g[j].Matcher)
}

type groupsResource struct {
	store Store
}

func newGroupsResource(store Store) *groupsResource {
	res := &groupsResource{
		store: store,
	}
	return res
}

// listGroups lists all Group resources.
func (gr *groupsResource) listGroups() ([]Group, error) {
	return gr.store.ListGroups()
}

// matchSpecHandler returns a ContextHandler that matches machine requests
// to a Spec and adds the Spec to the ctx and calls the next handler. The
// next handler should handle the case that no matching Spec is found.
func (gr *groupsResource) matchSpecHandler(next ContextHandler) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		attrs := labelsFromRequest(req)
		// match machine request
		group, err := gr.findMatch(attrs)
		if err == nil {
			// lookup Spec by id
			spec, err := gr.store.Spec(group.Spec)
			if err == nil {
				// add the Spec to the ctx for next handler
				ctx = withSpec(ctx, spec)
			}
		}
		next.ServeHTTP(ctx, w, req)
	}
	return ContextHandlerFunc(fn)
}

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

// findMatch returns the first Group whose Matcher is satisfied by the given
// labels. Groups are attempted in sorted order, preferring those with
// more matcher conditions, alphabetically.
func (gr *groupsResource) findMatch(labels Labels) (*Group, error) {
	groups, err := gr.store.ListGroups()
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(byMatcher(groups)))
	for _, group := range groups {
		if group.Matcher.Matches(labels) {
			return &group, nil
		}
	}
	return nil, errNoMatchingGroup
}
