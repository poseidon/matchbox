package storagepb

import (
	"net"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testGroup = &Group{
		Id:      "node1",
		Name:    "test group",
		Profile: "g1h2i3j4",
		Selector: map[string]string{
			"uuid": "a1b2c3d4",
			"mac":  "52:da:00:89:d8:10",
		},
		Metadata: []byte(`{"some-key":"some-val"}`),
	}
	testGroupWithoutProfile = &Group{
		Name:     "test group without profile",
		Profile:  "",
		Selector: map[string]string{"uuid": "a1b2c3d4"},
	}
)

func TestGroupParse(t *testing.T) {
	cases := []struct {
		json  string
		group *Group
	}{
		{`{"id":"node1","name":"test group","profile":"g1h2i3j4","selector":{"uuid":"a1b2c3d4","mac":"52:da:00:89:d8:10"},"metadata":{"some-key":"some-val"}}`, testGroup},
	}
	for _, c := range cases {
		group, _ := ParseGroup([]byte(c.json))
		assert.Equal(t, c.group, group)
	}
}

func TestGroupCopy(t *testing.T) {
	copy := testGroup.Copy()
	// assert that:
	// - Group fields are copied
	// - mutation of the copy does not affect the original
	assert.Equal(t, testGroup.Id, copy.Id)
	assert.Equal(t, testGroup.Name, copy.Name)
	assert.Equal(t, testGroup.Profile, copy.Profile)
	assert.Equal(t, testGroup.Selector, copy.Selector)
	assert.Equal(t, testGroup.Metadata, copy.Metadata)

	copy.Id = "a-copy"
	copy.Selector["region"] = "us-west"
	assert.NotEqual(t, testGroup.Id, copy.Id)
	assert.NotEqual(t, testGroup.Selector, copy.Selector)
}

func TestGroupMatches(t *testing.T) {
	cases := []struct {
		labels    map[string]string
		selectors map[string]string
		expected  bool
	}{
		{map[string]string{"a": "b"}, map[string]string{"a": "b"}, true},
		{map[string]string{"a": "b"}, map[string]string{"a": "c"}, false},
		{map[string]string{"uuid": "a", "mac": "b"}, map[string]string{"uuid": "a"}, true},
		{map[string]string{"uuid": "a"}, map[string]string{"uuid": "a", "mac": "b"}, false},
	}
	// assert that:
	// - Group selectors must be satisfied for a match
	// - labels may provide additional key/value pairs
	for _, c := range cases {
		group := &Group{Selector: c.selectors}
		assert.Equal(t, c.expected, group.Matches(c.labels))
	}
}

func TestNormalize(t *testing.T) {
	expectedInvalidMAC := &net.AddrError{Err: "invalid MAC address", Addr: "not-a-mac"}
	cases := []struct {
		selectors  map[string]string
		normalized map[string]string
		err        error
	}{
		{map[string]string{"platform": "metal"}, map[string]string{"platform": "metal"}, nil},
		{map[string]string{"mac": "52-da-00-89-d8-10"}, map[string]string{"mac": "52:da:00:89:d8:10"}, nil},
		{map[string]string{"MAC": "52-da-00-89-d8-10"}, map[string]string{"MAC": "52:da:00:89:d8:10"}, nil},
		// un-normalized MAC address should be normalized
		{map[string]string{"mac": "52-DA-00-89-D8-10"}, map[string]string{"mac": "52:da:00:89:d8:10"}, nil},
		{map[string]string{"MAC": "52-DA-00-89-D8-10"}, map[string]string{"MAC": "52:da:00:89:d8:10"}, nil},
		// invalid MAC address should be rejected
		{map[string]string{"mac": "not-a-mac"}, map[string]string{"mac": "not-a-mac"}, expectedInvalidMAC},
	}
	for _, c := range cases {
		group := &Group{Id: "id", Selector: c.selectors}
		err := group.Normalize()
		// assert that:
		// - Group selectors (MAC addresses) are normalized
		// - Invalid MAC addresses cause a normalization error
		assert.Equal(t, c.err, err)
		assert.Equal(t, c.normalized, group.Selector)
	}
}

func TestGroupValidate(t *testing.T) {
	cases := []struct {
		group *Group
		valid bool
	}{
		{&Group{Id: "node1", Profile: "k8s-controller"}, true},
		{testGroupWithoutProfile, false},
		{&Group{Id: "node1"}, false},
		{&Group{}, false},
	}
	for _, c := range cases {
		valid := c.group.AssertValid() == nil
		assert.Equal(t, c.valid, valid)
	}
}

func TestSelectorString(t *testing.T) {
	group := Group{
		Selector: map[string]string{
			"a": "b",
			"c": "d",
		},
	}
	expected := "a=b,c=d"
	assert.Equal(t, expected, group.selectorString())
}

func TestGroupSort(t *testing.T) {
	oneCondition := &Group{
		Name: "group with one selector",
		Selector: map[string]string{
			"region": "a",
		},
	}
	twoConditions := &Group{
		Name: "group with two selectors",
		Selector: map[string]string{
			"region": "a",
			"zone":   "z",
		},
	}
	dualConditions := &Group{
		Name: "group with two selectors",
		Selector: map[string]string{
			"region": "b",
			"zone":   "z",
		},
	}
	cases := []struct {
		input    []*Group
		expected []*Group
	}{
		{[]*Group{oneCondition, dualConditions, twoConditions}, []*Group{oneCondition, twoConditions, dualConditions}},
		{[]*Group{twoConditions, dualConditions, oneCondition}, []*Group{oneCondition, twoConditions, dualConditions}},
		{[]*Group{testGroup, testGroupWithoutProfile, oneCondition, twoConditions, dualConditions}, []*Group{oneCondition, testGroupWithoutProfile, testGroup, twoConditions, dualConditions}},
	}
	// assert that
	// - Group ordering is deterministic
	// - Groups are sorted by increasing Selector length
	// - when Selectors are equal in length, sort by key=value strings.
	for _, c := range cases {
		sort.Sort(ByReqs(c.input))
		assert.Equal(t, c.expected, c.input)
	}
}
