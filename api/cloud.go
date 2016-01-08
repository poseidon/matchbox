package api

import (
	"net/http"
	"strings"
	"time"
)

// CloudConfig defines a cloud-init config.
type CloudConfig struct {
	Content string
}

// cloudHandler returns a handler that responds with the cloud config for the
// requester.
func cloudHandler(store Store) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		attrs := labelsFromRequest(req)
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
