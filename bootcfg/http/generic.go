package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/coreos-baremetal/bootcfg/server"
	pb "github.com/coreos/coreos-baremetal/bootcfg/server/serverpb"
)

// genericHandler returns a handler that responds with generic file for
// the requester.
func (s *Server) genericHandler(core server.Server) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil || group.Profile == "" {
			http.NotFound(w, req)
			return
		}
		profile, err := core.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile})
		if err != nil || profile.GenericId == "" {
			http.NotFound(w, req)
			return
		}
		contents, err := core.GenericGet(ctx, profile.GenericId)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// collect data for rendering
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

		// render the template of a generic config with data
		var buf bytes.Buffer
		err = renderTemplate(&buf, data, contents)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		config := buf.String()
		http.ServeContent(w, req, "", time.Time{}, strings.NewReader(config))
	}
	return ContextHandlerFunc(fn)
}
