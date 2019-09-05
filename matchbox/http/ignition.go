package http

import (
	"bytes"
	"encoding/json"
	"github.com/coreos/fcct/config/common"
	"net/http"
	"strings"

	ct "github.com/coreos/container-linux-config-transpiler/config"
	fcct "github.com/coreos/fcct/config"

	ignition "github.com/coreos/ignition/config/v2_2"
	ignitionV2 "github.com/coreos/ignition/v2/config/v3_0"
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
			_, reportV2, errV2 := ignitionV2.Parse([]byte(contents))
			if err != nil || errV2 != nil {
				if err != nil {
					s.logger.Warningf("warning parsing Ignition JSON: %s", report.String())
				}
				if errV2 != nil {
					s.logger.Warningf("warning parsing Ignition JSON: %s", reportV2.String())
				}
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

		ignJSON := transpileConfig(s, buf.Bytes())

		if ignJSON == nil {
			http.NotFound(w, req)
			return
		}

		s.writeJSON(w, ignJSON)
		return
	}
	return http.HandlerFunc(fn)
}
func transpileConfig(s *Server, input []byte) []byte {
	fcctOptions := common.TranslateOptions{Strict: false, Pretty: false}

	// Convert fcc config into an Ignition Config
	ignV2json, errorFcct := fcct.Translate(input, fcctOptions)
	if errorFcct != nil {
		s.logger.Errorf("Parsing Container Linux with v3 schema failed, with try with v2.2: %s", errorFcct)
		// Parse bytes into a Container Linux Config
		config, ast, report := ct.Parse(input)
		if report.IsFatal() {
			s.logger.Errorf("error parsing Container Linux config: %s", report.String())
			return nil
		}

		// Convert Container Linux Config into an Ignition Config
		ign, report := ct.Convert(config, "", ast)
		if report.IsFatal() {
			s.logger.Errorf("error converting Container Linux config: %s", report.String())
			return nil
		}
		ignJSON, err := json.Marshal(ign)
		if err != nil {
			s.logger.Errorf("error JSON encoding: %v", err)
			return nil
		}
		return ignJSON
	}
	return ignV2json
}

// isIgnition returns true if the file should be treated as plain Ignition.
func isIgnition(filename string) bool {
	return strings.HasSuffix(filename, ".ign") || strings.HasSuffix(filename, ".ignition")
}
