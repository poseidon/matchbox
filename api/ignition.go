package api

import (
	"net/http"

	"golang.org/x/net/context"
)

// ignitionHandler returns a handler that responds with the ignition config
// for the requester.
func ignitionHandler(store Store) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		spec, err := specFromContext(ctx)
		if err != nil || spec.IgnitionConfig == "" {
			http.NotFound(w, req)
			return
		}
		config, err := store.IgnitionConfig(spec.IgnitionConfig)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		renderJSON(w, config)
	}
	return ContextHandlerFunc(fn)
}
