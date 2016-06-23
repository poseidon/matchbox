package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	fuze "github.com/coreos/fuze/config"
	ignition "github.com/coreos/ignition/config"
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// ignitionHandler returns a handler that responds with the Ignition config
// for the requester. The Ignition file referenced in the Profile is parsed
// as raw Ignition (for .ign/.ignition) or rendered to a Fuze config (YAML)
// and converted to Ignition. Ignition configs are served as HTTP JSON
// responses.
func ignitionHandler(srv server.Server) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil || group.Profile == "" {
			http.NotFound(w, req)
			return
		}
		profile, err := srv.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile})
		if err != nil || profile.IgnitionId == "" {
			http.NotFound(w, req)
			return
		}
		contents, err := srv.IgnitionGet(ctx, profile.IgnitionId)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// Skip rendering if raw Ignition JSON is provided
		if isIgnition(profile.IgnitionId) {
			_, err := ignition.Parse([]byte(contents))
			if err != nil {
				log.Warningf("warning parsing Ignition JSON: %v", err)
			}
			writeJSON(w, []byte(contents))
			return
		}

		// Fuze Config template

		// collect data for rendering Ignition Config
		data := make(map[string]interface{})
		if group.Metadata != nil {
			err = json.Unmarshal(group.Metadata, &data)
			if err != nil {
				log.Errorf("error unmarshalling metadata: %v", err)
				http.NotFound(w, req)
				return
			}
		}
		data["query"] = req.URL.RawQuery
		for key, value := range group.Selector {
			data[strings.ToLower(key)] = value
		}

		// render the template for an Ignition config with data
		var buf bytes.Buffer
		err = renderTemplate(&buf, data, contents)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// Parse fuze config into an Ignition config
		config, err := fuze.ParseAsV2_0_0(buf.Bytes())
		if err == nil {
			renderJSON(w, config)
			return
		}

		log.Errorf("error parsing Ignition config: %v", err)
		http.NotFound(w, req)
		return
	}
	return ContextHandlerFunc(fn)
}

// isIgnition returns true if the file should be treated as plain Ignition.
func isIgnition(filename string) bool {
	return strings.HasSuffix(filename, ".ign") || strings.HasSuffix(filename, ".ignition")
}
