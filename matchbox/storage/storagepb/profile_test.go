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
			Kernel: "/image/kernel",
			Initrd: []string{"/image/initrd_a"},
			Args:   []string{"a=b"},
		},
	}
	clone := profile.Copy()
	// assert that:
	// - Profile fields are copied to the clone
	// - Mutation of the clone does not affect the original
	assert.Equal(t, profile.Id, clone.Id)
	assert.Equal(t, profile.Name, clone.Name)
	assert.Equal(t, profile.IgnitionId, clone.IgnitionId)
	assert.Equal(t, profile.CloudId, clone.CloudId)
	assert.Equal(t, profile.Boot, clone.Boot)

	// mutate the NetBoot struct
	clone.Boot.Initrd = []string{"/image/initrd_b"}
	clone.Boot.Args = []string{"console=ttyS0"}
	assert.NotEqual(t, profile.Boot.Initrd, clone.Boot.Initrd)
	assert.NotEqual(t, profile.Boot.Args, clone.Boot.Args)
}

func TestNetBootCopy(t *testing.T) {
	boot := &NetBoot{
		Kernel: "/image/kernel",
		Initrd: []string{"/image/initrd_a"},
		Args:   []string{"a=b"},
	}

	clone := boot.Copy()
	// assert that:
	// - NetBoot fields are copied to the clone
	// - Mutation of the clone does not affect the original
	assert.Equal(t, boot.Kernel, clone.Kernel)
	assert.Equal(t, boot.Initrd, clone.Initrd)
	assert.Equal(t, boot.Args, clone.Args)

	// mutate the clone's slice field contents
	extra := []string{"extra"}
	copy(clone.Initrd, extra)
	copy(clone.Args, extra)
	assert.NotEqual(t, boot.Initrd, clone.Initrd)
	assert.NotEqual(t, boot.Args, clone.Args)
}
