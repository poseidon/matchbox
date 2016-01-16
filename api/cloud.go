package api

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// CloudConfig defines a cloud-init config.
type CloudConfig struct {
	Content string
}

// cloudHandler returns a handler that responds with the cloud config for the
// requester.
func cloudHandler(store Store) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		spec, err := specFromContext(ctx)
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
	return ContextHandlerFunc(fn)
}
