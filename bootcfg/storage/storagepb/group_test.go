package storagepb

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validMACStr = "52:da:00:89:d8:10"
	testGroup   = Group{
		Name:    "test group",
		Profile: "g1h2i3j4",
		Requirements: map[string]string{
			"uuid": "a1b2c3d4",
			"mac":  validMACStr,
		},
	}
	testGroupWithoutProfile = Group{
		Name:         "test group without profile",
		Profile:      "",
		Requirements: map[string]string{"uuid": "a1b2c3d4"},
	}
)

func TestGroupMatches(t *testing.T) {
	cases := []struct {
		labels   map[string]string
		reqs     map[string]string
		expected bool
	}{
		{map[string]string{"a": "b"}, map[string]string{"a": "b"}, true},
		{map[string]string{"a": "b"}, map[string]string{"a": "c"}, false},
		{map[string]string{"uuid": "a", "mac": "b"}, map[string]string{"uuid": "a"}, true},
		{map[string]string{"uuid": "a"}, map[string]string{"uuid": "a", "mac": "b"}, false},
	}
	// assert that:
	// - Group requirements are satisfied in order to be a match
	// - labels may provide additional key/value pairs
	for _, c := range cases {
		group := &Group{Requirements: c.reqs}
		assert.Equal(t, c.expected, group.Matches(c.labels))
	}
}

func TestRequirementString(t *testing.T) {
	group := Group{
		Requirements: map[string]string{
			"a": "b",
			"c": "d",
		},
	}
	expected := "a=b,c=d"
	assert.Equal(t, expected, group.requirementString())
}

func TestGroupSort(t *testing.T) {
	oneCondition := Group{
		Name: "group with one requirement",
		Requirements: map[string]string{
			"region": "a",
		},
	}
	twoConditions := Group{
		Name: "group with two requirements",
		Requirements: map[string]string{
			"region": "a",
			"zone":   "z",
		},
	}
	dualConditions := Group{
		Name: "group with two requirements",
		Requirements: map[string]string{
			"region": "b",
			"zone":   "z",
		},
	}
	cases := []struct {
		input    []Group
		expected []Group
	}{
		{[]Group{oneCondition, dualConditions, twoConditions}, []Group{oneCondition, twoConditions, dualConditions}},
		{[]Group{twoConditions, dualConditions, oneCondition}, []Group{oneCondition, twoConditions, dualConditions}},
		{[]Group{testGroup, testGroupWithoutProfile, oneCondition, twoConditions, dualConditions}, []Group{oneCondition, testGroupWithoutProfile, testGroup, twoConditions, dualConditions}},
	}
	// assert that
	// - Group ordering is deterministic
	// - Groups are sorted by increasing Requirements length
	// - when Requirements are equal in length, sort by key=value strings.
	for _, c := range cases {
		sort.Sort(ByReqs(c.input))
		assert.Equal(t, c.expected, c.input)
	}
}
