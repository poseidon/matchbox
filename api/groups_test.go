package api

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validMACStr = "52:54:00:89:d8:10"
	testGroup   = Group{
		Name:    "test group",
		Spec:    "g1h2i3j4",
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
