package api

import (
	"net"
	"net/http"
	"path/filepath"
)

// MachineAttrs collects machine identifiers and attributes.
type MachineAttrs struct {
	UUID string
	MAC  net.HardwareAddr
}

// Machine defines the configuration for a specific machine.
type Machine struct {
	// machine identifier
	ID string `json:"id"`
	// boot kernel, initrd, and kernel options
	BootConfig *BootConfig `json:"boot"`
	// reference a Spec
	SpecID string `json:"spec_id"`
}

// machineResource serves the configuration for a specific machine.
type machineResource struct {
	store Store
}

func newMachineResource(mux *http.ServeMux, pattern string, store Store) {
	mr := &machineResource{
		store: store,
	}
	mux.Handle(pattern, logRequests(requireGET(mr)))
}

func (r *machineResource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := filepath.Base(req.URL.Path)
	machine, err := r.store.Machine(id)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	if machine.BootConfig == nil && machine.SpecID != "" {
		// machine references a Spec, attempt to add Spec properties
		spec, err := r.store.Spec(machine.SpecID)
		if err == nil {
			machine.BootConfig = spec.BootConfig
		}
	}
	renderJSON(w, machine)
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
