package api

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"net/http"
	"strings"

	ignition "github.com/coreos/ignition/src/config"
	"golang.org/x/net/context"
)

// ignitionHandler returns a handler that responds with the Ignition config
// for the requester. The Ignition file referenced in the Spec is rendered
// with metadata and parsed and validated as either YAML or JSON based on the
// extension. The Ignition config is served as an HTTP JSON response.
func ignitionHandler(store Store) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil || group.Spec == "" {
			http.NotFound(w, req)
			return
		}
		spec, err := store.Spec(group.Spec)
		if err != nil || spec.IgnitionConfig == "" {
			http.NotFound(w, req)
			return
		}
		contents, err := store.IgnitionConfig(spec.IgnitionConfig)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// collect data for rendering Ignition Config
		data := make(map[string]interface{})
		for k := range group.Metadata {
			data[k] = group.Metadata[k]
		}
		data["query"] = req.URL.RawQuery

		// render the template for an Ignition config with data
		var buf bytes.Buffer
		err = renderTemplate(&buf, data, contents)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// Unmarshal YAML or JSON Ignition config
		var cfg ignition.Config
		if isYAML(spec.IgnitionConfig) {
			if err := yaml.Unmarshal(buf.Bytes(), &cfg); err != nil {
				log.Errorf("error parsing YAML Ignition config: %v", err)
				http.NotFound(w, req)
				return
			}
		} else {
			cfg, err = ignition.Parse(buf.Bytes())
			if err != nil {
				log.Errorf("error parsing JSON Ignition config: %v", err)
				http.NotFound(w, req)
				return
			}
		}
		// Marshal Ignition config as JSON HTTP response
		renderJSON(w, cfg)
	}
	return ContextHandlerFunc(fn)
}

func isYAML(filename string) bool {
	return strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml")
}
