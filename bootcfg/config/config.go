package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/coreos/coreos-baremetal/bootcfg/api"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

const (
	// APIVersion of the config types.
	APIVersion = "v1alpha1"
)

// Config parse errors
var (
	ErrIncorrectVersion = errors.New("api: incorrect API version")
)

// Config is a user defined matching of machines to specifications.
type Config struct {
	APIVersion string      `yaml:"api_version"`
	Groups     []api.Group `yaml:"groups"`
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
	if err := config.validate(); err != nil {
		return nil, err
	}
	return config, nil
}

// validate the group config's API version and reserved tag matchers.
func (c *Config) validate() error {
	if c.APIVersion != APIVersion {
		return ErrIncorrectVersion
	}
	for _, group := range c.Groups {
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

// PBGroups returns the parsed storagepb.Group slice.
func (c *Config) PBGroups() []*storagepb.Group {
	groups := make([]*storagepb.Group, len(c.Groups))
	i := 0
	for _, g := range c.Groups {
		group := &storagepb.Group{
			Id:           uuid.NewV4().String(),
			Name:         g.Name,
			Profile:      g.Spec,
			Metadata:     make(map[string]string),
			Requirements: g.Matcher,
		}
		// gRPC message fields must have concrete types.
		// Limit YAML metadata nesting to a depth of 1 for now.
		for key, unknown := range g.Metadata {
			switch val := unknown.(type) {
			case string:
				group.Metadata[key] = val
			default:
				// skip subtree
			}
		}
		groups[i] = group
		i++
	}
	return groups
}
