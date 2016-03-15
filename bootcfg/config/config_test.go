package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
	"github.com/stretchr/testify/assert"
)

var validData = `
api_version: v1alpha1
groups:
  - name: node1
    profile: worker
    require:
      role: worker
      region: us-central1-a
      mac: aB:Ab:3d:45:cD:10
    metadata:
      a: b
      c:
        - d
        - e
      f:
        g:
          h
`

var validGroups = []*storagepb.Group{
	&storagepb.Group{
		Id:      "node1",
		Name:    "node1",
		Profile: "worker",
		Requirements: map[string]string{
			"role":   "worker",
			"region": "us-central1-a",
			"mac":    "ab:ab:3d:45:cd:10",
		},
		Metadata: []byte(`{"a":"b","c":["d","e"],"f":{"g":"h"}}`),
	},
}

func TestLoadConfig(t *testing.T) {
	f, err := ioutil.TempFile("", "config.yaml")
	assert.Nil(t, err)
	defer os.Remove(f.Name())
	f.Write([]byte(validData))

	config, err := LoadConfig(f.Name())
	assert.Nil(t, err)
	assert.Equal(t, "v1alpha1", config.APIVersion)
	assert.Equal(t, validGroups, config.Groups)
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
		expectedGroups []*storagepb.Group
		expectedErr    error
	}{
		{validData, validGroups, nil},
		{invalidData, nil, ErrIncorrectVersion},
		{invalidYAML, nil, fmt.Errorf("yaml: found character that cannot start any token")},
	}
	for _, c := range cases {
		config, err := ParseConfig([]byte(c.data))
		assert.Equal(t, c.expectedErr, err)
		if c.expectedErr == nil {
			assert.Equal(t, c.expectedGroups, config.Groups)
		}
	}
}

func TestValidate(t *testing.T) {
	incorrectVersion := &Config{
		APIVersion: "v1wrong",
	}
	invalidMAC := &Config{
		APIVersion: "v1alpha1",
		YAMLGroups: []Group{
			Group{
				Requirements: map[string]string{
					"mac": "?:?:?:?",
				},
			},
		},
	}
	nonNormalizedMAC := &Config{
		APIVersion: "v1alpha1",
		YAMLGroups: []Group{
			Group{
				Requirements: map[string]string{
					"mac": "aB:Ab:3d:45:cD:10",
				},
			},
		},
	}
	validConfig := &Config{
		APIVersion: "v1alpha1",
		YAMLGroups: []Group{
			Group{
				Name:    "node1",
				Profile: "worker",
				Requirements: map[string]string{
					"role":   "worker",
					"region": "us-central1-a",
					"mac":    "ab:ab:3d:45:cd:10",
				},
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

func TestFilterSlice(t *testing.T) {
	s := []interface{}{"a", 3.14, "b", "c"}
	expected := []string{"a", "b", "c"}
	filtered := filterSlice(s)
	assert.Equal(t, expected, filtered)
}

func TestFilterKeys(t *testing.T) {
	m := map[interface{}]interface{}{
		"a":  "b",
		3.14: "c",
	}
	expected := map[string]interface{}{
		"a": "b",
	}
	filtered := filterNonStringKeys(m)
	assert.Equal(t, expected, filtered)
}

func TestFilterValues(t *testing.T) {
	m := map[string]interface{}{
		"a": "b",
		"c": 3.14,
		"d": map[interface{}]interface{}{
			"e":  "f",
			3.14: "g",
			"h":  []interface{}{"i", "j", "k", 3.14},
		},
		"l": true,
		"m": "true",
	}
	expected := map[string]interface{}{
		"a": "b",
		"d": map[string]interface{}{
			"e": "f",
			"h": []string{"i", "j", "k"},
		},
		"m": "true",
	}
	filtered := filterValues(m)
	assert.Equal(t, expected, filtered)
}
