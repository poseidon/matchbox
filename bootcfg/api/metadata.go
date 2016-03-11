package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

const plainContentType = "plain/text"

func metadataHandler() ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		w.Header().Set(contentType, plainContentType)

		var data map[string]interface{}
		err = json.Unmarshal(group.Metadata, &data)
		if err != nil {
			log.Error("error unmarshalling metadata")
			http.NotFound(w, req)
			return
		}
		for key, value := range data {
			fmt.Fprintf(w, "%s=%s\n", strings.ToUpper(key), value)
		}
		attrs := labelsFromRequest(req)
		for key, value := range attrs {
			fmt.Fprintf(w, "%s=%s\n", strings.ToUpper(key), value)
		}
	}
	return ContextHandlerFunc(fn)
}
