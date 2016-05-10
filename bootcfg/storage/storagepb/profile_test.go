package storagepb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testProfile = &Profile{
		Id:         "id",
		CloudId:    "cloud.yaml",
		IgnitionId: "ignition.json",
	}
)

func TestProfileParse(t *testing.T) {
	cases := []struct {
		json    string
		profile *Profile
	}{
		{`{"id": "id", "cloud_id": "cloud.yaml", "ignition_id": "ignition.json"}`, testProfile},
	}
	for _, c := range cases {
		profile, _ := ParseProfile([]byte(c.json))
		assert.Equal(t, c.profile, profile)
	}
}

func TestProfileValidate(t *testing.T) {
	cases := []struct {
		profile *Profile
		valid   bool
	}{
		{testProfile, true},
		{&Profile{Id: "a1b2c3d4"}, true},
		{&Profile{}, false},
	}
	for _, c := range cases {
		valid := c.profile.AssertValid() == nil
		assert.Equal(t, c.valid, valid)
	}
}

func TestProfileCopy(t *testing.T) {
	profile := &Profile{
		Id:         "id",
		CloudId:    "cloudy.tmpl",
		IgnitionId: "ignition.tmpl",
		Boot: &NetBoot{
			Kernel:  "/image/kernel",
			Initrd:  []string{"/image/initrd_a"},
			Cmdline: map[string]string{"a": "b"},
		},
	}
	copy := profile.Copy()
	// assert that:
	// - Profile fields are copied
	// - mutation of the copy does not affect the original
	assert.Equal(t, profile.Id, copy.Id)
	assert.Equal(t, profile.Name, copy.Name)
	assert.Equal(t, profile.IgnitionId, copy.IgnitionId)
	assert.Equal(t, profile.CloudId, copy.CloudId)
	assert.Equal(t, profile.Boot, copy.Boot)

	copy.Id = "a-copy"
	copy.Boot.Initrd = []string{"/image/initrd_b"}
	copy.Boot.Cmdline["c"] = "d"
	assert.NotEqual(t, profile.Id, copy.Id)
	assert.NotEqual(t, profile.Boot.Initrd, copy.Boot.Initrd)
	assert.NotEqual(t, profile.Boot.Cmdline, copy.Boot.Cmdline)
}
