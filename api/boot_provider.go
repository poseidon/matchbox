package api

import (
	"fmt"
)

// BootConfig describes the boot configuration of a host.
type BootConfig struct {
	// the URL of the kernel boot image
	Kernel string `json:"kernel"`
	// the initrd URLs which will be flattened into a single filesystem
	Initrd []string `json:"initrd"`
	// command line arguments to the kernel
	Cmdline map[string]interface{} `json:"cmdline"`
}

// A BootConfigProvider provides a mapping from MAC addresses to BootConfigs.
type BootConfigProvider interface {
	Add(addr string, config *BootConfig)
	Get(addr string) (*BootConfig, error)
}

// NewBootConfig returns a new memory map BootConfigProvider.
func NewBootConfigProvider() BootConfigProvider {
	return &bootConfigProvider{
		mac2boot: make(map[string]*BootConfig),
	}
}

const DefaultAddr = "default"

// bootConfigProvider implements a MAC address to Boot config map in memory.
type bootConfigProvider struct {
	mac2boot map[string]*BootConfig
}

func (p *bootConfigProvider) Add(addr string, config *BootConfig) {
	p.mac2boot[addr] = config
}

func (p *bootConfigProvider) Get(addr string) (*BootConfig, error) {
	if config, ok := p.mac2boot[addr]; ok {
		return config, nil
	} else if config, ok := p.mac2boot[DefaultAddr]; ok {
		return config, nil
	}
	return nil, fmt.Errorf("no boot config for %s", addr)
}
