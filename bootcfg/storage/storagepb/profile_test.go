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
		{&Profile{Id: "a1b2c3d4"}, true},
		{&Profile{}, false},
	}
	for _, c := range cases {
		valid := c.profile.AssertValid() == nil
		assert.Equal(t, c.valid, valid)
	}
}
