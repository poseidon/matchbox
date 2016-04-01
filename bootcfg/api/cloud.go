package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	cloudinit "github.com/coreos/coreos-cloudinit/config"
	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// CloudConfig defines a cloud-init config.
type CloudConfig struct {
	Content string
}

// cloudHandler returns a handler that responds with the cloud config for the
// requester.
func cloudHandler(srv server.Server) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil || group.Profile == "" {
			http.NotFound(w, req)
			return
		}
		resp, err := srv.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile})
		if err != nil || resp.Profile.CloudId == "" {
			http.NotFound(w, req)
			return
		}
		contents, err := srv.CloudGet(ctx, resp.Profile.CloudId)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// collect data for rendering
		var data map[string]interface{}
		err = json.Unmarshal(group.Metadata, &data)
		if err != nil {
			log.Error("error unmarshalling metadata")
			http.NotFound(w, req)
			return
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
			log.Error("error parsing user-data")
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
