package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/coreos/coreos-baremetal/bootcfg/api"
	"github.com/stretchr/testify/assert"
)

var validData = `
api_version: v1alpha1
groups:
  - name: node1
    spec: worker
    require:
      role: worker
      region: us-central1-a
`

var validConfig = &Config{
	APIVersion: "v1alpha1",
	Groups: []api.Group{
		api.Group{
			Name: "node1",
			Spec: "worker",
			Matcher: api.RequirementSet(map[string]string{
				"role":   "worker",
				"region": "us-central1-a",
			}),
		},
	},
}

func TestLoadConfig(t *testing.T) {
	f, err := ioutil.TempFile("", "config.yaml")
	assert.Nil(t, err)
	defer os.Remove(f.Name())
	f.Write([]byte(validData))

	config, err := LoadConfig(f.Name())
	assert.Equal(t, validConfig, config)
	assert.Nil(t, err)
	// read from file that does not exist
	config, err = LoadConfig("")
	assert.Nil(t, config)
	assert.NotNil(t, err)
}

func TestParseConfig(t *testing.T) {
	invalidData := `api_version:`
	invalidYAML := `	tabs	tabs	tabs`

	cases := []struct {
		data           string
		expectedConfig *Config
		expectedErr    error
	}{
		{validData, validConfig, nil},
		{invalidData, nil, ErrIncorrectVersion},
		{invalidYAML, nil, fmt.Errorf("yaml: found character that cannot start any token")},
	}
	for _, c := range cases {
		config, err := ParseConfig([]byte(c.data))
		assert.Equal(t, c.expectedConfig, config)
		assert.Equal(t, c.expectedErr, err)
	}
}

func TestValidate(t *testing.T) {
	incorrectVersion := &Config{
		APIVersion: "v1wrong",
	}
	invalidMAC := &Config{
		APIVersion: "v1alpha1",
		Groups: []api.Group{
			api.Group{
				Matcher: api.RequirementSet(map[string]string{
					"mac": "?:?:?:?",
				}),
			},
		},
	}
	nonNormalizedMAC := &Config{
		APIVersion: "v1alpha1",
		Groups: []api.Group{
			api.Group{
				Matcher: api.RequirementSet(map[string]string{
					"mac": "aB:Ab:3d:45:cD:10",
				}),
			},
		},
	}

	cases := []struct {
		config      *Config
		expectedErr error
	}{
		{validConfig, nil},
		{incorrectVersion, ErrIncorrectVersion},
		{invalidMAC, fmt.Errorf("config: invalid MAC address ?:?:?:?")},
		{nonNormalizedMAC, fmt.Errorf("config: normalize MAC address aB:Ab:3d:45:cD:10 to ab:ab:3d:45:cd:10")},
	}
	for _, c := range cases {
		assert.Equal(t, c.expectedErr, c.config.validate())
	}
}
