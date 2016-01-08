package api

import (
	"net/http"
	"path/filepath"
)

// pixiecoreHandler returns a handler that renders Boot Configs as JSON to
// implement the Pixiecore API specification.
// https://github.com/danderson/pixiecore/blob/master/README.api.md
func pixiecoreHandler(store Store) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		macAddr, err := parseMAC(filepath.Base(req.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// pixiecore only provides MAC addresses
		attrs := LabelSet(map[string]string{"mac": macAddr.String()})
		spec, err := getMatchingSpec(store, attrs)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		renderJSON(w, spec.BootConfig)
	}
	return http.HandlerFunc(fn)
}
