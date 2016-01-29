package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var (
	validMACStr         = "52:da:00:89:d8:10"
	nonNormalizedMACStr = "52:dA:00:89:d8:10"
	testGroup           = Group{
		Name: "test group",
		Spec: "g1h2i3j4",
		Metadata: map[string]string{
			"k8s_version": "v1.1.2",
			"pod_network": "10.2.0.0/16",
		},
		Matcher: RequirementSet(map[string]string{"uuid": "a1b2c3d4"}),
	}
	testGroupWithMAC = Group{
		Name:    "machine with a MAC",
		Spec:    "g1h2i3j4",
		Matcher: RequirementSet(map[string]string{"mac": validMACStr}),
	}
	testGroupNoSpec = &Group{
		Name:    "test group with missing spec",
		Spec:    "",
		Matcher: RequirementSet(map[string]string{"uuid": "a1b2c3d4"}),
	}
)

func TestByMatcherSort(t *testing.T) {
	oneCondition := Group{
		Name: "two matcher conditions",
		Matcher: RequirementSet(map[string]string{
			"region": "a",
		}),
	}
	twoConditions := Group{
		Name: "group with two matcher conditions",
		Matcher: RequirementSet(map[string]string{
			"region": "a",
			"zone":   "z",
		}),
	}
	dualConditions := Group{
		Name: "another group with two matcher conditions",
		Matcher: RequirementSet(map[string]string{
			"region": "b",
			"zone":   "z",
		}),
	}
	cases := []struct {
		input    []Group
		expected []Group
	}{
		{[]Group{oneCondition, dualConditions, twoConditions}, []Group{oneCondition, twoConditions, dualConditions}},
		{[]Group{twoConditions, dualConditions, oneCondition}, []Group{oneCondition, twoConditions, dualConditions}},
		{[]Group{testGroup, testGroupWithMAC, oneCondition, twoConditions, dualConditions}, []Group{testGroupWithMAC, oneCondition, testGroup, twoConditions, dualConditions}},
	}
	// assert that
	// - groups are sorted in increasing Matcher length
	// - when Matcher lengths are equal, groups are sorted by key=value strings.
	// - group ordering is deterministic
	for _, c := range cases {
		sort.Sort(byMatcher(c.input))
		assert.Equal(t, c.expected, c.input)
	}
}

func TestNewGroupsResource(t *testing.T) {
	store := &fixedStore{}
	gr := newGroupsResource(store)
	assert.Equal(t, store, gr.store)
}

func TestGroupsResource_ListGroups(t *testing.T) {
	expectedGroups := []Group{Group{Name: "test group"}}
	store := &fixedStore{
		Groups: expectedGroups,
	}
	res := newGroupsResource(store)
	groups, err := res.listGroups()
	assert.Nil(t, err)
	assert.Equal(t, expectedGroups, groups)
}

func TestGroupsResource_MatchSpecHandler(t *testing.T) {
	store := &fixedStore{
		Groups: []Group{testGroup},
		Specs:  map[string]*Spec{testGroup.Spec: testSpec},
	}
	gr := newGroupsResource(store)
	next := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		spec, err := specFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, testSpec, spec)
		fmt.Fprintf(w, "next handler called")
	}
	// assert that:
	// - request arguments are used to match uuid=a1b2c3d4 -> testGroup
	// - the group's Spec is found by id and added to the context
	// - next handler is called
	h := gr.matchSpecHandler(ContextHandlerFunc(next))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "?uuid=a1b2c3d4", nil)
	h.ServeHTTP(context.Background(), w, req)
	assert.Equal(t, "next handler called", w.Body.String())
}

func TestGroupsResource_MatchGroupHandler(t *testing.T) {
	store := &fixedStore{
		Groups: []Group{testGroup},
	}
	gr := newGroupsResource(store)
	next := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		assert.Nil(t, err)
		assert.Equal(t, &testGroup, group)
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
		Groups: []Group{testGroup},
		Specs:  map[string]*Spec{testGroup.Spec: testSpec},
	}
	uuidLabel := LabelSet(map[string]string{
		"uuid": "a1b2c3d4",
	})

	cases := []struct {
		store         Store
		labels        Labels
		expectedGroup *Group
		expectedErr   error
	}{
		{store, uuidLabel, &testGroup, nil},
		{store, nil, nil, errNoMatchingGroup},
		// no groups in the store
		{&emptyStore{}, uuidLabel, nil, errNoMatchingGroup},
	}

	for _, c := range cases {
		gr := newGroupsResource(c.store)
		group, err := gr.findMatch(c.labels)
		assert.Equal(t, c.expectedGroup, group)
		assert.Equal(t, c.expectedErr, err)
	}
}
