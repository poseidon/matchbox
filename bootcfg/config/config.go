package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/coreos/coreos-baremetal/bootcfg/api"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
	"github.com/coreos/pkg/capnslog"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

const (
	// APIVersion of the config types.
	APIVersion = "v1alpha1"
)

// Config parse errors.
var (
	ErrIncorrectVersion = errors.New("config: incorrect API version")
)

var log = capnslog.NewPackageLogger("github.com/coreos/coreos-baremetal/bootcfg", "config")

// Config is a user defined matching of machines to specifications.
type Config struct {
	APIVersion string `yaml:"api_version"`
	// allow YAML source for Groups
	YAMLGroups []api.Group `yaml:"groups"`
	// populate protobuf Groups at parse
	Groups []*storagepb.Group `yaml:"-"`
}

// LoadConfig opens a file and parses YAML data to returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(os.ExpandEnv(path))
	if err != nil {
		return nil, err
	}
	return ParseConfig(data)
}

// ParseConfig parses YAML data and returns a Config.
func ParseConfig(data []byte) (*Config, error) {
	config := new(Config)
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	// convert YAML Groups into protobuf Groups
	config.Groups = make([]*storagepb.Group, 0)
	for _, ygroup := range config.YAMLGroups {
		group := &storagepb.Group{
			Name:         ygroup.Name,
			Profile:      ygroup.Spec,
			Requirements: normalizeMatchers(ygroup.Matcher),
		}
		// Id: Generate a random UUID or use the name
		if ygroup.Name == "" {
			group.Id = uuid.NewV4().String()
		} else {
			group.Id = group.Name
		}
		// Metadata: go-yaml unmarshal provides Config.Metadata as a
		// map[string]interface{}, which unmarshals nested maps as
		// map[interface{}]interface{}. Walk the metadata, filtering non-string
		// keys and nested elements.
		if b, err := json.Marshal(filterValues(ygroup.Metadata)); err == nil {
			group.Metadata = b
		} else {
			return nil, fmt.Errorf("config: cannot marshal metadata %v", err)
		}
		config.Groups = append(config.Groups, group)
	}
	// validate Config and Groups
	if err := config.validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func normalizeMatchers(reqs map[string]string) map[string]string {
	for key, val := range reqs {
		switch strings.ToLower(key) {
		case "mac":
			if macAddr, err := net.ParseMAC(val); err == nil {
				// range iteration copy with mutable map
				reqs[key] = macAddr.String()
				log.Errorf("normalizing MAC address %s to %s", val, macAddr.String())
			}
		}
	}
	return reqs
}

// filterValues returns a new map, filtering out key/value pairs whose value
// is not a string, []string, or string key'd map. Recurses on interface
// values.
func filterValues(unknown map[string]interface{}) map[string]interface{} {
	// copy the map, skipping all disallowed value types
	m := make(map[string]interface{})
	for key, v := range unknown {
		switch val := v.(type) {
		case string:
			m[key] = val
		case []string:
			m[key] = val
		case []interface{}:
			m[key] = filterSlice(val)
		case map[interface{}]interface{}:
			m[key] = filterValues(filterNonStringKeys(val))
		default:
			log.Errorf("ignoring metadata value %v", val)
		}
	}
	return m
}

// filterSlice returns a new slice, filtering out elements whose keys are not
// strings.
func filterSlice(unknown []interface{}) []string {
	s := make([]string, 0, len(unknown))
	for _, e := range unknown {
		switch elem := e.(type) {
		case string:
			s = append(s, elem)
		default:
			log.Errorf("ignoring metadata elem %v", elem)
		}
	}
	return s
}

// filterNonStringKeys returns a new map, filtering out key/value pairs whose
// keys are not strings.
func filterNonStringKeys(unknown map[interface{}]interface{}) map[string]interface{} {
	// copy the map, skipping all non-string keys
	m := make(map[string]interface{})
	for k, val := range unknown {
		switch key := k.(type) {
		case string:
			m[key] = val
		default:
			log.Errorf("ignoring metadata key %v", key)
		}
	}
	return m
}

// validate the group config's API version and reserved tag matchers.
func (c *Config) validate() error {
	if c.APIVersion != APIVersion {
		return ErrIncorrectVersion
	}
	for _, group := range c.YAMLGroups {
		for key, val := range group.Matcher {
			switch strings.ToLower(key) {
			case "mac":
				macAddr, err := net.ParseMAC(val)
				if err != nil {
					return fmt.Errorf("config: invalid MAC address %s", val)
				}
				if val != macAddr.String() {
					return fmt.Errorf("config: normalize MAC address %s to %v", val, macAddr.String())
				}
			}
		}
	}
	return nil
}
