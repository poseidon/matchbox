package api

import (
	"net"
)

// MachineAttrs collects machine identifiers and attributes.
type MachineAttrs struct {
	UUID string
	MAC  net.HardwareAddr
}
