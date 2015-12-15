package api

import (
	"encoding/json"
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
		log.Infof("pixiecore boot config request for %+v", attrs)

		config, err := store.BootConfig(attrs)
		if err != nil {
			http.NotFound(w, req)
			return
		}
		if err := json.NewEncoder(w).Encode(config); err != nil {
			log.Infof("error writing to response, %s", err)
		}
	}
	return http.HandlerFunc(fn)
}
