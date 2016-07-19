package http

import (
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

		// collect data for rendering
		data, err := collectVariables(req, group)
		if err != nil {
			s.logger.Errorf("error collecting variables: %v", err)
			http.NotFound(w, req)
			return
		}

		w.Header().Set(contentType, plainContentType)
		for key, value := range data {
			fmt.Fprintf(w, "%s=%v\n", strings.ToUpper(key), value)
		}
	}
	return ContextHandlerFunc(fn)
}
