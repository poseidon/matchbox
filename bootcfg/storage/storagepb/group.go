package storagepb

import (
  "sort"
  "strings"
)

// Labels present key to value mappings, independent of their storage.
type Labels interface {
  // Get returns the value for the given label.
  Get(label string) string
}

// LabelSet is a map of key:value labels.
type LabelSet map[string]string

// Get returns the value for the given label.
func (ls LabelSet) Get(label string) string {
  return ls[label]
}

// Matches returns true if the given labels satisfy all the requirements,
// false otherwise.
func (g *Group) Matches(labels Labels) bool {
  for key, val := range g.Requirements {
    if labels == nil || labels.Get(key) != val {
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

// byMatcher defines a collection of Group structs which have a deterministic
// order by increasing number of Requirements, then by sorted key/value
// strings. For example, a Group with Requirements {a:b, c:d} should be ordered
// before one with {a:d, c:d} and {a:b, d:e}.
type ByMatcher []Group

func (groups ByMatcher) Len() int {
  return len(groups)
}

func (groups ByMatcher) Swap(i, j int) {
  groups[i], groups[j] = groups[j], groups[i]
}

func (groups ByMatcher) Less(i, j int) bool {
  if len(groups[i].Requirements) == len(groups[j].Requirements) {
    return groups[i].requirementString() < groups[j].requirementString()
  }
  return len(groups[i].Requirements) < len(groups[j].Requirements)
}
