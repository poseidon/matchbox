package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CloudConfig defines the cloud-init config to initialize a client machine.
type CloudConfig struct {
	Content string
}

// cloudHandler returns a handler that responds with the cloud config the
// client machine should use.
func cloudHandler(store Store) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		attrs := attrsFromRequest(req)
		spec, err := getMatchingSpec(store, attrs)
		if err != nil || spec.CloudConfig == "" {
			http.NotFound(w, req)
			return
		}

		config, err := store.CloudConfig(spec.CloudConfig)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		http.ServeContent(w, req, "", time.Time{}, strings.NewReader(config.Content))
	}
	return http.HandlerFunc(fn)
}

// getMatchingSpec returns the Spec matching the given attributes.
func getMatchingSpec(store Store, attrs MachineAttrs) (*Spec, error) {
	if machine, err := store.Machine(attrs.UUID); err == nil && machine.Spec != nil {
		return machine.Spec, nil
	}
	if machine, err := store.Machine(attrs.MAC.String()); err == nil && machine.Spec != nil {
		return machine.Spec, nil
	}
	return nil, fmt.Errorf("no spec matching %v", attrs)
}
