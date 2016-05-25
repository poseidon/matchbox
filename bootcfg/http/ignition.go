package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	ignition "github.com/coreos/ignition/config"
	ignitionTypes "github.com/coreos/ignition/config/types"
	ignitionV1 "github.com/coreos/ignition/config/v1"
	ignitionV1Types "github.com/coreos/ignition/config/v1/types"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// ignitionHandler returns a handler that responds with the Ignition config
// for the requester. The Ignition file referenced in the Profile is rendered
// with metadata and parsed and validated as either YAML or JSON based on the
// extension. The Ignition config is served as an HTTP JSON response.
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

		// Unmarshal YAML or JSON to Ignition V2
		cfg, err := parseToV2(buf.Bytes())
		if err == nil {
			renderJSON(w, cfg)
			return
		}

		// Unmarshal YAML or JSON to Ignition V1
		oldCfg, err := parseToV1(buf.Bytes())
		if err == nil {
			renderJSON(w, oldCfg)
			return
		}

		log.Errorf("error parsing Ignition config: %v", err)
		http.NotFound(w, req)
		return
	}
	return ContextHandlerFunc(fn)
}

// parseToV2 parses raw JSON in Ignition v2 format and returns the
// Ignition v2 Config struct.
func parseToV2(data []byte) (cfg ignitionTypes.Config, err error) {
	// parse JSON v2 to Ignition
	cfg, err = ignition.ParseFromLatest(data)
	if err == nil {
		return cfg, nil
	}
	if majorVersion(data) == 2 {
		err = yaml.Unmarshal(data, &cfg)
	}
	return cfg, err
}

// parseToV1 parses raw JSON or YAML in Ignition v1 format and returns the
// Ignition v1 Config struct.
func parseToV1(data []byte) (cfg ignitionV1Types.Config, err error) {
	// parse JSON v1 to Ignition
	cfg, err = ignitionV1.Parse(data)
	if err == nil {
		return cfg, nil
	}
	// unmarshal YAML v1 to Ignition
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func majorVersion(data []byte) int64 {
	var composite struct {
		Version  *int `json:"ignitionVersion" yaml:"ignition_version"`
		Ignition struct {
			Version *string `json:"version" yaml:"version"`
		} `json:"ignition" yaml:"ignition"`
	}
	if yaml.Unmarshal(data, &composite) != nil {
		return 0
	}
	var major int64
	if composite.Ignition.Version != nil && *composite.Ignition.Version == "2.0.0" {
		major = 2
	} else if composite.Version != nil {
		major = int64(*composite.Version)
	}
	return major
}
