package http

import (
	"bytes"
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
// matching the request. The Ignition file referenced in the Profile is parsed
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
			_, report, err := ignition.Parse([]byte(contents))
			if err != nil {
				s.logger.Warningf("warning parsing Ignition JSON: %s", report.String())
			}
			s.writeJSON(w, []byte(contents))
			return
		}

		// Fuze Config template

		// collect data for rendering
		data, err := collectVariables(req, group)
		if err != nil {
			s.logger.Errorf("error collecting variables: %v", err)
			http.NotFound(w, req)
			return
		}

		// render the template for an Ignition config with data
		var buf bytes.Buffer
		err = s.renderTemplate(&buf, data, contents)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// Parse bytes into a Fuze Config
		config, report := fuze.Parse(buf.Bytes())
		if report.IsFatal() {
			s.logger.Errorf("error parsing Fuze config: %s", report.String())
			http.NotFound(w, req)
			return
		}

		// Convert Fuze Config into an Ignition Config
		ign, report := fuze.ConvertAs2_0_0(config)
		if report.IsFatal() {
			s.logger.Errorf("error converting Fuze config: %s", report.String())
			http.NotFound(w, req)
			return
		}

		s.renderJSON(w, ign)
		return
	}
	return ContextHandlerFunc(fn)
}

// isIgnition returns true if the file should be treated as plain Ignition.
func isIgnition(filename string) bool {
	return strings.HasSuffix(filename, ".ign") || strings.HasSuffix(filename, ".ignition")
}
