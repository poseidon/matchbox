package api

import (
	"sort"
	"strings"
)

// RequirementSet is a map of key:value equality requirements which
// match against any Labels which are supersets.
type RequirementSet map[string]string

// Matches returns true if the given labels satisfy all the requirements,
// false otherwise.
func (r RequirementSet) Matches(labels map[string]string) bool {
	for key, val := range r {
		if labels == nil || labels[key] != val {
			return false
		}
	}
	return true
}

func (r RequirementSet) String() string {
	requirements := make([]string, 0, len(r))
	for key, value := range r {
		requirements = append(requirements, key+"="+value)
	}
	// sort by "key=value" pairs for a deterministic ordering
	sort.StringSlice(requirements).Sort()
	return strings.Join(requirements, ",")
}
