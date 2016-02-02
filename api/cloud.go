package api

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	cloudinit "github.com/coreos/coreos-cloudinit/config"
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
		group, err := groupFromContext(ctx)
		if err != nil || group.Spec == "" {
			http.NotFound(w, req)
			return
		}
		spec, err := store.Spec(group.Spec)
		if err != nil || spec.CloudConfig == "" {
			http.NotFound(w, req)
			return
		}
		contents, err := store.CloudConfig(spec.CloudConfig)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// collect data for rendering
		data := make(map[string]string)
		for k := range group.Metadata {
			data[k] = group.Metadata[k]
		}

		// render the template of a cloud config with data
		var buf bytes.Buffer
		err = renderTemplate(&buf, data, contents)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		config := buf.String()
		if !cloudinit.IsCloudConfig(config) && !cloudinit.IsScript(config) {
			log.Errorf("error parsing user-data")
			http.NotFound(w, req)
			return
		}

		if cloudinit.IsCloudConfig(config) {
			if _, err = cloudinit.NewCloudConfig(config); err != nil {
				log.Errorf("error parsing cloud config: %v", err)
				http.NotFound(w, req)
				return
			}
		}
		http.ServeContent(w, req, "", time.Time{}, strings.NewReader(config))
	}
	return ContextHandlerFunc(fn)
}
