package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

const plainContentType = "plain/text"

func (s *Server) metadataHandler() ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		w.Header().Set(contentType, plainContentType)

		data := make(map[string]interface{})
		if group.Metadata != nil {
			err = json.Unmarshal(group.Metadata, &data)
			if err != nil {
				log.Errorf("error unmarshalling metadata: %v", err)
				http.NotFound(w, req)
				return
			}
		}
		for key, value := range group.Selector {
			data[key] = value
		}

		for key, value := range data {
			fmt.Fprintf(w, "%s=%v\n", strings.ToUpper(key), value)
		}
	}
	return ContextHandlerFunc(fn)
}
