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
	return nil, fmt.Errorf("no matching boot configuration")
}

// AddUUID adds a BootConfig for the machine with the given UUID.
func (a *MapBootAdapter) AddUUID(uuid string, config *BootConfig) {
	a.uuids[uuid] = config
}

// AddMAC adds a BootConfig for the machine with NIC with the given MAC
// address.
func (a *MapBootAdapter) AddMAC(mac net.HardwareAddr, config *BootConfig) {
	a.macs[mac.String()] = config
}

// SetDefault sets the default or fallback BootConfig to use if no machine
// attribute matches are found.
func (a *MapBootAdapter) SetDefault(config *BootConfig) {
	a.fallback = config
}
