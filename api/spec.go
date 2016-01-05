package api

import (
	"fmt"
	"net/http"
	"path/filepath"
)

// Spec is a named group of configs.
type Spec struct {
	// spec identifier
	ID string `json:"id"`
	// boot kernel, initrd, and kernel options
	BootConfig *BootConfig `json:"boot"`
	// cloud config id
	CloudConfig string `json:"cloud_id"`
	// ignition config id
	IgnitionConfig string `json:"ignition_id"`
}

// specResource serves Spec resources by id.
type specResource struct {
	store Store
}

func newSpecResource(mux *http.ServeMux, pattern string, store Store) {
	gr := &specResource{
		store: store,
	}
	mux.Handle(pattern, logRequests(requireGET(gr)))
}

func (r *specResource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := filepath.Base(req.URL.Path)
	spec, err := r.store.Spec(id)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	renderJSON(w, spec)
}

// getMatchingSpec returns the Spec matching the given attributes. Attributes
// are matched in priority order (UUID, MAC, default).
func getMatchingSpec(store Store, attrs MachineAttrs) (*Spec, error) {
	if machine, err := store.Machine(attrs.UUID); err == nil && machine.Spec != nil {
		return machine.Spec, nil
	}
	if machine, err := store.Machine(attrs.MAC.String()); err == nil && machine.Spec != nil {
		return machine.Spec, nil
	}
	if machine, err := store.Machine("default"); err == nil && machine.Spec != nil {
		return machine.Spec, nil
	}
	return nil, fmt.Errorf("no spec matching %v", attrs)
}
