package api

import (
	"net"
	"net/http"
)

// MachineAttrs collects machine identifiers and attributes.
type MachineAttrs struct {
	UUID string
	MAC  net.HardwareAddr
}

// attrsFromRequest returns MachineAttrs from request query parameters.
func attrsFromRequest(req *http.Request) MachineAttrs {
	params := req.URL.Query()
	// if MAC address is unset or fails to parse, leave it nil
	var macAddr net.HardwareAddr
	if params.Get("mac") != "" {
		macAddr, _ = parseMAC(params.Get("mac"))
	}
	return MachineAttrs{
		UUID: params.Get("uuid"),
		MAC:  macAddr,
	}
}

// parseMAC wraps net.ParseMAC with logging.
func parseMAC(s string) (net.HardwareAddr, error) {
	macAddr, err := net.ParseMAC(s)
	if err != nil {
		// invalid MAC arguments may be common
		log.Debugf("error parsing MAC address: %s", err)
		return nil, err
	}
	return macAddr, err
}
