package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
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
func (s *Server) ignitionHandler(core server.Server) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"labels": labelsFromRequest(nil, req),
			}).Infof("No matching group")
			http.NotFound(w, req)
			return
		}

		profile, err := core.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile})
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"labels":     labelsFromRequest(nil, req),
				"group":      group.Id,
				"group_name": group.Name,
			}).Infof("No profile named: %s", group.Profile)
			http.NotFound(w, req)
			return
		}

		contents, err := core.IgnitionGet(ctx, profile.IgnitionId)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"labels":     labelsFromRequest(nil, req),
				"group":      group.Id,
				"group_name": group.Name,
				"profile":    group.Profile,
			}).Infof("No Ignition or Fuze template named: %s", profile.IgnitionId)
			http.NotFound(w, req)
			return
		}

		// match was successful
		s.logger.WithFields(logrus.Fields{
			"labels":  labelsFromRequest(nil, req),
			"group":   group.Id,
			"profile": profile.Id,
		}).Debug("Matched an Ignition or Fuze template")

		// Skip rendering if raw Ignition JSON is provided
		if isIgnition(profile.IgnitionId) {
			_, err := ignition.Parse([]byte(contents))
			if err != nil {
				s.logger.Warningf("warning parsing Ignition JSON: %v", err)
			}
			s.writeJSON(w, []byte(contents))
			return
		}

		// Fuze Config template

		// collect data for rendering
		data := make(map[string]interface{})
		if group.Metadata != nil {
			err = json.Unmarshal(group.Metadata, &data)
			if err != nil {
				s.logger.Errorf("error unmarshalling metadata: %v", err)
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
		err = s.renderTemplate(&buf, data, contents)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// Parse fuze config into an Ignition config
		config, err := fuze.ParseAsV2_0_0(buf.Bytes())
		if err == nil {
			s.renderJSON(w, config)
			return
		}

		s.logger.Errorf("error parsing Ignition config: %v", err)
		http.NotFound(w, req)
		return
	}
	return ContextHandlerFunc(fn)
}

// isIgnition returns true if the file should be treated as plain Ignition.
func isIgnition(filename string) bool {
	return strings.HasSuffix(filename, ".ign") || strings.HasSuffix(filename, ".ignition")
}
