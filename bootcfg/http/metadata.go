package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

const plainContentType = "plain/text"

// genericHandler returns a handler that responds with the metadata env file
// matching the request.
func (s *Server) metadataHandler() ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"labels": labelsFromRequest(nil, req),
			}).Infof("No matching group")
			http.NotFound(w, req)
			return
		}

		// match was successful
		s.logger.WithFields(logrus.Fields{
			"labels": labelsFromRequest(nil, req),
			"group":  group.Id,
		}).Debug("Matched group metadata")

		// collect data for response
		data := make(map[string]interface{})
		if group.Metadata != nil {
			err = json.Unmarshal(group.Metadata, &data)
			if err != nil {
				s.logger.Errorf("error unmarshalling metadata: %v", err)
				http.NotFound(w, req)
				return
			}
		}
		for key, value := range group.Selector {
			data[key] = value
		}

		w.Header().Set(contentType, plainContentType)
		for key, value := range data {
			fmt.Fprintf(w, "%s=%v\n", strings.ToUpper(key), value)
		}
	}
	return ContextHandlerFunc(fn)
}
