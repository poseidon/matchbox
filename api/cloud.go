package api

import (
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
		config, err := store.CloudConfig(attrs)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		http.ServeContent(w, req, "", time.Time{}, strings.NewReader(config.Content))
	}
	return http.HandlerFunc(fn)
}
