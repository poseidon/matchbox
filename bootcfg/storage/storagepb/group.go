package storagepb

import (
	"sort"
	"strings"
)

// Matches returns true if the given labels satisfy all the requirements,
// false otherwise.
func (g *Group) Matches(labels map[string]string) bool {
	for key, val := range g.Requirements {
		if labels == nil || labels[key] != val {
			return false
		}
	}
	return true
}

// requirementString returns Group requirements as a string of sorted key
// value pairs for comparisons.
func (g *Group) requirementString() string {
	reqs := make([]string, 0, len(g.Requirements))
	for key, value := range g.Requirements {
		reqs = append(reqs, key+"="+value)
	}
	// sort by "key=value" pairs for a deterministic ordering
	sort.StringSlice(reqs).Sort()
	return strings.Join(reqs, ",")
}

// ByReqs defines a collection of Group structs which have a deterministic
// sorted order by increasing number of Requirements, then by sorted key/value
// strings. For example, a Group with Requirements {a:b, c:d} should be ordered
// after one with {a:b} and before one with {a:d, c:d}.
type ByReqs []*Group

func (groups ByReqs) Len() int {
	return len(groups)
}

func (groups ByReqs) Swap(i, j int) {
	groups[i], groups[j] = groups[j], groups[i]
}

func (groups ByReqs) Less(i, j int) bool {
	if len(groups[i].Requirements) == len(groups[j].Requirements) {
		return groups[i].requirementString() < groups[j].requirementString()
	}
	return len(groups[i].Requirements) < len(groups[j].Requirements)
}
