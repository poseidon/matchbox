package storagepb

import (
	"encoding/json"
	"net"
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

// ToRichGroup converts a Group into a RichGroup suitable for writing and
// user manipulation.
func (g *Group) ToRichGroup() (*RichGroup, error) {
	metadata := make(map[string]interface{})
	if g.Metadata != nil {
		err := json.Unmarshal(g.Metadata, &metadata)
		if err != nil {
			return nil, err
		}
	}
	return &RichGroup{
		Id:           g.Id,
		Name:         g.Name,
		Profile:      g.Profile,
		Requirements: g.Requirements,
		Metadata:     metadata,
	}, nil
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

// RichGroup is a user provided Group definition.
type RichGroup struct {
	// machine readable Id
	Id string `json:"id,omitempty"`
	// Human readable name
	Name string `json:"name,omitempty"`
	// Profile id
	Profile string `json:"profile,omitempty"`
	// tags required to match the group
	Requirements map[string]string `json:"requirements,omitempty"`
	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ToGroup converts a user provided RichGroup into a Group which can be
// serialized as a protocol buffer.
func (rg *RichGroup) ToGroup() (*Group, error) {
	var metadata []byte
	if rg.Metadata != nil {
		var err error
		metadata, err = json.Marshal(rg.Metadata)
		if err != nil {
			return nil, err
		}
	}
	return &Group{
		Id:           rg.Id,
		Name:         rg.Name,
		Profile:      rg.Profile,
		Requirements: normalizeSelectors(rg.Requirements),
		Metadata:     metadata,
	}, nil
}

func normalizeSelectors(selectors map[string]string) map[string]string {
	for key, val := range selectors {
		switch strings.ToLower(key) {
		case "mac":
			if macAddr, err := net.ParseMAC(val); err == nil {
				// range iteration copy with mutable map
				selectors[key] = macAddr.String()
			}
		}
	}
	return selectors
}
