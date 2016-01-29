package api

import (
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
		for key, value := range group.Metadata {
			fmt.Fprintf(w, "%s=%s\n", strings.ToUpper(key), value)
		}
		attrs := labelsFromRequest(req)
		for key, value := range attrs {
			fmt.Fprintf(w, "%s=%s\n", strings.ToUpper(key), value)
		}
	}
	return ContextHandlerFunc(fn)
}
