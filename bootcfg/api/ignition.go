package api

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
	"strings"

	"github.com/coreos/coreos-baremetal/bootcfg/storage"
	ignition "github.com/coreos/ignition/src/config"
	"golang.org/x/net/context"
)

// ignitionHandler returns a handler that responds with the Ignition config
// for the requester. The Ignition file referenced in the Profile is rendered
// with metadata and parsed and validated as either YAML or JSON based on the
// extension. The Ignition config is served as an HTTP JSON response.
func ignitionHandler(store storage.Store) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil || group.Profile == "" {
			http.NotFound(w, req)
			return
		}
		profile, err := store.ProfileGet(group.Profile)
		if err != nil || profile.IgnitionId == "" {
			http.NotFound(w, req)
			return
		}
		contents, err := store.IgnitionGet(profile.IgnitionId)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// collect data for rendering Ignition Config
		var data map[string]interface{}
		err = json.Unmarshal(group.Metadata, &data)
		if err != nil {
			log.Errorf("error unmarshalling metadata: %v", err)
			http.NotFound(w, req)
			return
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
		if isYAML(profile.IgnitionId) {
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
