package api

import (
	"bytes"
	"net/http"
	"text/template"

	ignition "github.com/coreos/ignition/src/config"
	"golang.org/x/net/context"
)

// ignitionHandler returns a handler that responds with the ignition config
// for the requester.
func ignitionHandler(store Store) ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		group, err := groupFromContext(ctx)
		if err != nil || group.Spec == "" {
			http.NotFound(w, req)
			return
		}
		spec, err := store.Spec(group.Spec)
		if err != nil || spec.IgnitionConfig == "" {
			http.NotFound(w, req)
			return
		}
		contents, err := store.IgnitionConfig(spec.IgnitionConfig)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// collect data for rendering Ignition Config
		data := make(map[string]string)
		for k := range group.Metadata {
			data[k] = group.Metadata[k]
		}
		data["query"] = req.URL.RawQuery

		// render the Ignition Config with context and group data
		tmpl, err := template.New("ignition").Parse(contents)
		if err != nil {
			log.Errorf("error parsing template: %v", err)
			http.NotFound(w, req)
			return
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			log.Errorf("error rendering template: %v", err)
			http.NotFound(w, req)
			return
		}

		// validate the Ignition Config
		config, err := ignition.Parse(buf.Bytes())
		if err != nil {
			log.Errorf("error parsing ignition config: %v", err)
			http.NotFound(w, req)
			return
		}
		renderJSON(w, config)
	}
	return ContextHandlerFunc(fn)
}
