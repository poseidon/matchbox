package http

import (
	"bytes"
	"net/http"
	"strings"

	ct "github.com/coreos/container-linux-config-transpiler/config"
	ignition "github.com/coreos/ignition/config/v2_2"
	"github.com/sirupsen/logrus"

	"github.com/poseidon/matchbox/matchbox/server"
	pb "github.com/poseidon/matchbox/matchbox/server/serverpb"
)

// ignitionHandler returns a handler that responds with the Ignition config
// matching the request. The Ignition file referenced in the Profile is parsed
// as raw Ignition (for .ign/.ignition) or rendered from a Container Linux
// Config (YAML) and converted to Ignition. Ignition configs are served as HTTP
// JSON responses.
func (s *Server) ignitionHandler(core server.Server) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
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

		contents, err := core.IgnitionGet(ctx, &pb.IgnitionGetRequest{Name: profile.IgnitionId})
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"labels":     labelsFromRequest(nil, req),
				"group":      group.Id,
				"group_name": group.Name,
				"profile":    group.Profile,
			}).Infof("No Ignition or Container Linux Config template named: %s", profile.IgnitionId)
			http.NotFound(w, req)
			return
		}

		// match was successful
		s.logger.WithFields(logrus.Fields{
			"labels":  labelsFromRequest(nil, req),
			"group":   group.Id,
			"profile": profile.Id,
		}).Debug("Matched an Ignition or Container Linux Config template")

		// Skip rendering if raw Ignition JSON is provided
		if isIgnition(profile.IgnitionId) {
			_, report, err := ignition.Parse([]byte(contents))
			if err != nil {
				s.logger.Warningf("warning parsing Ignition JSON: %s", report.String())
			}
			s.writeJSON(w, []byte(contents))
			return
		}

		// Container Linux Config template

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

		// Parse bytes into a Container Linux Config
		config, ast, report := ct.Parse(buf.Bytes())
		if report.IsFatal() {
			s.logger.Errorf("error parsing Container Linux config: %s", report.String())
			http.NotFound(w, req)
			return
		}

		// Convert Container Linux Config into an Ignition Config
		ign, report := ct.Convert(config, "", ast)
		if report.IsFatal() {
			s.logger.Errorf("error converting Container Linux config: %s", report.String())
			http.NotFound(w, req)
			return
		}

		s.renderJSON(w, ign)
		return
	}
	return http.HandlerFunc(fn)
}

// isIgnition returns true if the file should be treated as plain Ignition.
func isIgnition(filename string) bool {
	return strings.HasSuffix(filename, ".ign") || strings.HasSuffix(filename, ".ignition")
}
