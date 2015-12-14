package api

import (
	"fmt"
	"net"
)

// MapBootAdapter maps MachineAttrs to BootConfigs using an in-memory map.
type MapBootAdapter struct {
	uuids    map[string]*BootConfig
	macs     map[string]*BootConfig
	fallback *BootConfig
}

// NewMapBootAdapter returns a new in-memory BootAdapter.
func NewMapBootAdapter() *MapBootAdapter {
	return &MapBootAdapter{
		uuids: make(map[string]*BootConfig),
		macs:  make(map[string]*BootConfig),
	}
}

// Get returns the BootConfig for the machine with the given attributes.
// Matches are searched in priority order: UUID, MAC address, default.
func (a *MapBootAdapter) Get(attrs MachineAttrs) (*BootConfig, error) {
	if config, ok := a.uuids[attrs.UUID]; ok {
		return config, nil
	}
	if config, ok := a.macs[attrs.MAC.String()]; ok {
		return config, nil
	}
	if a.fallback != nil {
		return a.fallback, nil
	}
	log.Infof("No boot config found for %+v", attrs)
	return nil, fmt.Errorf("no matching boot configuration")
}

// SetUUID sets the BootConfig for the machine with the given UUID.
func (a *MapBootAdapter) SetUUID(uuid string, config *BootConfig) {
	a.uuids[uuid] = config
}

// SetMAC sets the BootConfig for the NIC with the given MAC address.
func (a *MapBootAdapter) SetMAC(mac net.HardwareAddr, config *BootConfig) {
	a.macs[mac.String()] = config
}

// SetDefault sets the default BootConfig if no machine attributes match.
func (a *MapBootAdapter) SetDefault(config *BootConfig) {
	a.fallback = config
}
