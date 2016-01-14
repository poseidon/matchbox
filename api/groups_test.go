package api

import (
	"fmt"
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

// Test parsing group config YAML data.
func TestParseGroupConfig(t *testing.T) {
	validData := `
api_version: v1alpha1
groups:
  - name: node1
    spec: worker
    require:
      role: worker
      region: us-central1-a
`
	validConfig := &GroupConfig{
		APIVersion: "v1alpha1",
		Groups: []Group{
			Group{
				Name: "node1",
				Spec: "worker",
				Matcher: RequirementSet(map[string]string{
					"role":   "worker",
					"region": "us-central1-a",
				}),
			},
		},
	}
	wrongVersion := `api_version:`
	invalidMAC := `
api_version: v1alpha1
groups:
  - name: group
    require:
      mac: ?:?:?:?
`
	nonNormalizedMAC := `
api_version: v1alpha1
groups:
  - name: group
    require:
      mac: aB:Ab:3d:45:cD:10
`

	cases := []struct {
		data           string
		expectedConfig *GroupConfig
		expectedErr    error
	}{
		{validData, validConfig, nil},
		{wrongVersion, nil, ErrInvalidVersion},
		{invalidMAC, nil, fmt.Errorf("api: invalid MAC address ?:?:?:?")},
		{nonNormalizedMAC, nil, fmt.Errorf("api: normalize MAC address aB:Ab:3d:45:cD:10 to ab:ab:3d:45:cd:10")},
	}
	for _, c := range cases {
		config, err := ParseGroupConfig([]byte(c.data))
		assert.Equal(t, c.expectedConfig, config)
		assert.Equal(t, c.expectedErr, err)
	}
}
