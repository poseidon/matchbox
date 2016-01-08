package api

import (
	"errors"
	"fmt"
	"sort"

	"gopkg.in/yaml.v2"
)

// GroupConfig parser errors
var (
	ErrInvalidVersion = errors.New("api: mismatched API version")
)

// Group associates matcher conditions with a Specification identifier. The
// Matcher conditions may be satisfied by zero or more machines.
type Group struct {
	// Human readable name (optional)
	Name string `yaml:"name"`
	// Spec identifier
	Spec string `yaml:"spec"`
	// matcher conditions
	Matcher RequirementSet `yaml:"require"`
}

// GroupConfig define an group import structure.
type GroupConfig struct {
	APIVersion string  `yaml:"api_version"`
	Groups     []Group `yaml:"groups"`
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

// ParseGroupConfig parses a YAML group config and returns a GroupConfig.
func ParseGroupConfig(data []byte) (*GroupConfig, error) {
	config := new(GroupConfig)
	err := yaml.Unmarshal(data, config)
	if err == nil && config.APIVersion != APIVersion {
		return nil, ErrInvalidVersion
	}
	return config, err
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
func (r *groupsResource) listGroups() ([]Group, error) {
	return r.store.ListGroups()
}

// select returns the first Group whose Matcher is satisfied by the given
// labels. Groups are attempted in sorted order, preferring those with
// more matcher conditions.
func (r *groupsResource) findMatch(labels Labels) (*Group, error) {
	groups, err := r.store.ListGroups()
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(byMatcher(groups)))
	for _, group := range groups {
		if group.Matcher.Matches(labels) {
			return &group, nil
		}
	}
	return nil, fmt.Errorf("no Group matching %v", labels)
}
