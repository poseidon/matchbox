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
		attrs := MachineAttrs{MAC: macAddr}
		machine, err := store.Machine(attrs.MAC.String())
		if err != nil {
			http.NotFound(w, req)
			return
		}
		renderJSON(w, machine.BootConfig)
	}
	return http.HandlerFunc(fn)
}
