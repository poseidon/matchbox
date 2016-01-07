package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGroupConfig(t *testing.T) {
	validData := `
api_version: v1alpha1
groups:
  - name: node1
    spec: worker-central
    require:
      role: worker
      region: us-central1-a
`
	validConfig := &GroupConfig{
		APIVersion: "v1alpha1",
		Groups: []Group{
			Group{
				Name:          "node1",
				Specification: "worker-central",
				Matcher: RequirementSet(map[string]string{
					"role":   "worker",
					"region": "us-central1-a",
				}),
			},
		},
	}
	wrongVersion := `api_version:`

	cases := []struct {
		data           string
		expectedConfig *GroupConfig
		expectedErr    error
	}{
		{validData, validConfig, nil},
		{wrongVersion, nil, ErrInvalidVersion},
	}
	for _, c := range cases {
		config, err := ParseGroupConfig([]byte(c.data))
		assert.Equal(t, c.expectedConfig, config)
		assert.Equal(t, c.expectedErr, err)
	}
}
